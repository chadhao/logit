package api

import (
	"errors"
	"net/http"
	"time"

	"github.com/chadhao/logit/modules/record/model"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// LoadRoutes 路由添加
func LoadRoutes(e *echo.Echo) {
	e.POST("/record", addRecord)
	e.DELETE("/record/:id", deleteLastestRecord)
	e.GET("/records", getRecords)
	e.POST("/record/note", addNote)
}

// addRecord 添加一条新的记录
func addRecord(c echo.Context) error {
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

// deleteLastestRecord 删除上一条记录
func deleteLastestRecord(c echo.Context) error {
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

	if err := reqR.deleteRecord(userID); err != nil {
		return err
	}
	return nil
}

// getRecords 获取记录
func getRecords(c echo.Context) (err error) {
	reqR := new(reqRecords)
	reqR.From, _ = time.Parse(time.RFC3339, c.QueryParam("from"))
	reqR.To, _ = time.Parse(time.RFC3339, c.QueryParam("to"))

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

// addNote 为记录添加笔记
func addNote(c echo.Context) (err error) {
	userID, err := primitive.ObjectIDFromHex(c.Request().Header.Get("userID"))
	if err != nil {
		return err
	}
	req := new(reqAddNote)
	if err := c.Bind(req); err != nil {
		return err
	}
	if !req.isUsersRecord(userID) {
		return errors.New("no authorization")
	}

	var note model.INote
	switch req.NoteType {
	case model.OTHERWORKNOTE:
		note, err = req.constructToOtherWorkNote()
		if err != nil {
			return err
		}
	case model.MODIFICATIONNOTE:
		note, err = req.constructToModificationNote(userID)
		if err != nil {
			return err
		}
	default:
		return errors.New("no match noteType")
	}

	if err = note.Add(); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, note)
}

// offlineSyncRecords 离线返回在线状态后记录同步
func offlineSyncRecords(c echo.Context) {}
