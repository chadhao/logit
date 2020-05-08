package api

import (
	"net/http"
	"sort"

	logInternals "github.com/chadhao/logit/modules/log/internals"

	"github.com/chadhao/logit/modules/record/model"
	userApi "github.com/chadhao/logit/modules/user/internals"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// addRecord 添加一条新的记录
func addRecord(c echo.Context) error {

	uid, _ := c.Get("user").(primitive.ObjectID)

	req := new(reqAddRecord)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	r, err := req.constructToRecord(uid)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	if err = r.Add(); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, r)
}

// getLatestRecord 获取上一条记录
func getLatestRecord(c echo.Context) error {

	uid, _ := c.Get("user").(primitive.ObjectID)

	r, err := model.GetLastestRecord(uid)
	if err != nil {
		return err
	}

	notesMap, err := model.GetNotesByRecordIDs([]primitive.ObjectID{r.ID})
	if err != nil {
		return err
	}

	respRecord := &respRecord{
		Record: *r,
		Notes:  notesMap[r.ID],
	}

	return c.JSON(http.StatusOK, respRecord)
}

// deleteLatestRecord 删除上一条记录
func deleteLatestRecord(c echo.Context) error {

	uid, _ := c.Get("user").(primitive.ObjectID)

	req := new(reqRecord)
	var err error
	if req.ID, err = primitive.ObjectIDFromHex(c.Param("id")); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if err := req.deleteRecord(uid); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, "success")
}

// getRecords 获取记录
func getRecords(c echo.Context) error {

	req := new(reqRecords)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	uid, _ := c.Get("user").(primitive.ObjectID)

	if req.DriverID != uid.Hex() {
		return c.JSON(http.StatusUnauthorized, "driver has no authorization")
	}

	records, err := req.getRecords()
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, records)
}

// addNote 为记录添加笔记
func addNote(c echo.Context) error {

	req := new(reqAddNote)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	uid, _ := c.Get("user").(primitive.ObjectID)
	if !req.isDriversRecord(uid) {
		return c.JSON(http.StatusUnauthorized, "driver has no authorization")
	}

	var (
		note model.INote
		err  error
	)
	switch req.NoteType {
	case model.OTHERWORKNOTE:
		note, err = req.constructToOtherWorkNote()
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
	case model.MODIFICATIONNOTE:
		note, err = req.constructToModificationNote(uid)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
	case model.TRIPNOTE:
		note, err = req.constructToTripNote()
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
	default:
		return c.JSON(http.StatusBadRequest, "no match noteType")
	}

	if err = note.Add(); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, note)
}

// offlineSyncRecords 离线返回在线状态后记录同步
// 1. 对records按照时间排序，检查相邻两条之间的时间位置是否准确
// 2. 批量更新入数据库
func offlineSyncRecords(c echo.Context) error {
	uid, _ := c.Get("user").(primitive.ObjectID)

	reqs := []reqAddRecord{}
	log := &logInternals.ReqAddLog{
		Type:    "error",
		FromFun: "offlineSyncRecords",
		From:    &uid,
	}

	if err := c.Bind(&reqs); err != nil {
		e := err.Error()
		log.Message = &e
		go logInternals.AddLog(log)
		return c.JSON(http.StatusBadRequest, e)
	}
	l := len(reqs)
	if l == 0 {
		e := "no records obtained"
		log.Message = &e
		go logInternals.AddLog(log)
		return c.JSON(http.StatusBadRequest, e)
	}

	sort.Slice(reqs, func(a, b int) bool {
		return reqs[a].Time.Before(reqs[b].Time)
	})

	records := model.Records{}
	for i := 0; i < l; i++ {
		r, err := reqs[i].constructToSyncRecord(uid)
		if err != nil {
			e := err.Error()
			log.Message = &e
			go logInternals.AddLog(log)
			return err
		}
		records = append(records, *r)
	}

	if err := records.SyncAdd(); err != nil {
		e := err.Error()
		log.Message = &e
		go logInternals.AddLog(log)
		return c.JSON(http.StatusBadRequest, e)
	}
	return c.JSON(http.StatusOK, records)

}

// Transport operator 权限下的record操作
// toGetRecords 获取记录
func toGetRecords(c echo.Context) error {

	req := new(reqRecords)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	var (
		uid primitive.ObjectID
		did primitive.ObjectID
		tid primitive.ObjectID
		err error
	)
	uid, _ = c.Get("user").(primitive.ObjectID)
	did, err = primitive.ObjectIDFromHex(req.DriverID)
	tid, err = primitive.ObjectIDFromHex(req.TransportOperatorID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if !userApi.HasAccessTo(uid, did, tid) {
		return c.JSON(http.StatusUnauthorized, "to has no authorization")
	}

	records, err := req.getRecords()
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, records)
}
