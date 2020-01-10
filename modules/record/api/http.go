package api

import (
	"errors"
	"net/http"

	"github.com/chadhao/logit/modules/record/model"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

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

	userID, err := primitive.ObjectIDFromHex(c.Request().Header.Get("userID"))
	if err != nil {
		return err
	}

	req := new(reqRecord)
	if err := c.Bind(req); err != nil {
		return err
	}

	if err := req.deleteRecord(userID); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, "success")
}

// getRecords 获取记录
func getRecords(c echo.Context) error {

	userID, err := primitive.ObjectIDFromHex(c.Request().Header.Get("userID"))
	if err != nil {
		return err
	}

	req := new(reqRecords)
	if err := c.Bind(req); err != nil {
		return err
	}

	records, err := req.getRecords(userID)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, records)
}

// addNote 为记录添加笔记
func addNote(c echo.Context) error {
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
