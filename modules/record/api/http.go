package api

import (
	"errors"
	"github.com/chadhao/logit/modules/record/model"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

// AddRecord 添加一条新的记录
func AddRecord(c echo.Context) error {
	userID, err := primitive.ObjectIDFromHex(c.Request().Header.Get("userID"))
	if err != nil {
		return err
	}
	rar := new(reqAddRecord)
	if err := c.Bind(rar); err != nil {
		return err
	}

	// vehicleID := user.GetVehicleID()
	vehicleID := primitive.NewObjectID()

	r, err := rar.constructToRecord(userID, vehicleID)
	if err != nil {
		return err
	}
	if err = r.Add(); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, r)
}

// DeleteLastRecord 删除上一条记录
func DeleteLastRecord(c echo.Context) error {
	recordID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		return err
	}
	userID, err := primitive.ObjectIDFromHex(c.Request().Header.Get("userID"))
	if err != nil {
		return err
	}
	reqR := &reqRecord{
		ID: recordID,
	}
	r, err := reqR.getRecord()
	if err != nil {
		return err
	}
	if r.UserID != userID {
		return errors.New("no authorization")
	}
	if err := r.Delete(); err != nil {
		return err
	}
	return nil
}

// GetRecords 获取记录
func GetRecords(c echo.Context) (err error) {
	reqR := new(reqRecords)
	if err := c.Bind(reqR); err != nil {
		return err
	}
	userID, err := primitive.ObjectIDFromHex(c.Request().Header.Get("userID"))
	if err != nil {
		return err
	}
	records, err := reqR.getRecords(userID)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, records)
}

// AddNote 为记录添加笔记
func AddNote(c echo.Context) (err error) {
	recordID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		return err
	}
	reqR := &reqRecord{
		ID: recordID,
	}
	r, err := reqR.getRecord()
	if err != nil {
		return err
	}
	userID, err := primitive.ObjectIDFromHex(c.Request().Header.Get("userID"))
	if err != nil {
		return err
	}
	if r.UserID != userID {
		return
	}
	var note model.INote
	noteType := model.NoteType(c.FormValue("noteType"))
	switch noteType {
	case model.SYSTEMNOTE:
		req := new(reqAddSystemNote)
		if err = c.Bind(req); err != nil {
			return err
		}
		note, err = req.constructToSystemNote()
		if err != nil {
			return err
		}
	case model.OTHERWORKNOTE:
		req := new(reqAddOtherWorkNote)
		if err = c.Bind(req); err != nil {
			return err
		}
		note, err = req.constructToOtherWorkNote()
		if err != nil {
			return err
		}
	case model.MODIFICATIONNOTE:
		req := new(reqAddModificationNote)
		if err = c.Bind(req); err != nil {
			return err
		}
		note, err = req.constructToModificationNote(userID)
		if err != nil {
			return err
		}
	}

	if err = note.Add(); err != nil {
		return err
	}

	return nil
}

// OfflineSyncRecords 离线返回在线状态后记录同步
func OfflineSyncRecords(c echo.Context) {}
