package api

import (
	"errors"
	"net/http"
	"sort"

	"github.com/chadhao/logit/modules/record/model"
	"github.com/chadhao/logit/modules/user/constant"
	"github.com/chadhao/logit/utils"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// addRecord 添加一条新的记录
func addRecord(c echo.Context) error {

	roles := utils.RolesAssert(c.Get("roles"))
	if !roles.Is(constant.ROLE_DRIVER) {
		return errors.New("not driver")
	}

	uid, _ := c.Get("user").(primitive.ObjectID)

	req := new(reqAddRecord)
	if err := c.Bind(req); err != nil {
		return err
	}

	r, err := req.constructToRecord(uid)
	if err != nil {
		return err
	}
	if err = r.Add(); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, r)
}

// getLatestRecord 获取上一条记录
func getLatestRecord(c echo.Context) error {

	roles := utils.RolesAssert(c.Get("roles"))
	if !roles.Is(constant.ROLE_DRIVER) {
		return errors.New("not driver")
	}

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

	roles := utils.RolesAssert(c.Get("roles"))
	if !roles.Is(constant.ROLE_DRIVER) {
		return errors.New("not driver")
	}

	uid, _ := c.Get("user").(primitive.ObjectID)

	req := new(reqRecord)
	var err error
	if req.ID, err = primitive.ObjectIDFromHex(c.Param("id")); err != nil {
		return err
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
		return err
	}
	uid, _ := c.Get("user").(primitive.ObjectID)

	roles := utils.RolesAssert(c.Get("roles"))
	switch {
	case roles.Is(constant.ROLE_ADMIN):
	case roles.Is(constant.ROLE_DRIVER):
		if req.DriverID != uid.Hex() {
			return errors.New("no authorization")
		}
	default:
		return errors.New("not allowed")
	}

	records, err := req.getRecords()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, records)
}

// addNote 为记录添加笔记
func addNote(c echo.Context) error {

	req := new(reqAddNote)
	if err := c.Bind(req); err != nil {
		return err
	}

	uid, _ := c.Get("user").(primitive.ObjectID)

	roles := utils.RolesAssert(c.Get("roles"))
	switch {
	case roles.Is(constant.ROLE_ADMIN):
	case roles.Is(constant.ROLE_DRIVER):
		if !req.isDriversRecord(uid) {
			return errors.New("no authorization")
		}
	default:
		return errors.New("not allowed")
	}

	var (
		note model.INote
		err  error
	)
	switch req.NoteType {
	case model.OTHERWORKNOTE:
		note, err = req.constructToOtherWorkNote()
		if err != nil {
			return err
		}
	case model.MODIFICATIONNOTE:
		note, err = req.constructToModificationNote(uid)
		if err != nil {
			return err
		}
	case model.TRIPNOTE:
		note, err = req.constructToTripNote()
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
// 1. 对records按照时间排序，检查相邻两条之间的时间位置是否准确
// 2. 批量更新入数据库
func offlineSyncRecords(c echo.Context) error {
	roles := utils.RolesAssert(c.Get("roles"))
	if !roles.Is(constant.ROLE_DRIVER) {
		return errors.New("not driver")
	}

	uid, _ := c.Get("user").(primitive.ObjectID)

	reqs := []reqAddRecord{}
	if err := c.Bind(&reqs); err != nil {
		return err
	}
	l := len(reqs)
	if l == 0 {
		return errors.New("no records obtained")
	}

	sort.Slice(reqs, func(a, b int) bool {
		return reqs[a].Time.Before(reqs[b].Time)
	})

	records := model.Records{}
	for i := 0; i < l; i++ {
		r, err := reqs[i].constructToSyncRecord(uid)
		if err != nil {
			return err
		}
		records = append(records, *r)
	}

	if err := records.SyncAdd(); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, records)

}
