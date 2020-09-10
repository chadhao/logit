package api

import (
	"net/http"
	"sort"
	"time"

	"github.com/chadhao/logit/modules/record/model"
	"github.com/chadhao/logit/modules/record/service"
	userApi "github.com/chadhao/logit/modules/user/internals"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

/************************************************************************************************/
/********************************** Driver 权限下的record操作 *************************************/
/************************************************************************************************/

type addRecordRequest struct {
	Type          model.Type         `json:"type" valid:"required"`
	Time          time.Time          `json:"time" valid:"required"`
	Duration      string             `json:"duration" valid:"required"`
	StartLocation model.Location     `json:"startLocation" valid:"required"`
	EndLocation   model.Location     `json:"endLocation" valid:"required"`
	VehicleID     primitive.ObjectID `json:"vehicleID" valid:"-"`
	StartMileAge  *float64           `json:"startDistance,omitempty" valid:"-"`
	EndMileAge    *float64           `json:"endDistance,omitempty" valid:"-"`
	ClientTime    *time.Time         `json:"clientTime,omitempty" valid:"-"`
}

func (r *addRecordRequest) toCreateRecordInput(driverID primitive.ObjectID) (*service.CreateRecordInput, error) {
	duration, err := time.ParseDuration(r.Duration)
	if err != nil {
		return nil, err
	}
	out := &service.CreateRecordInput{
		DriverID:      driverID,
		Type:          r.Type,
		Time:          r.Time,
		Duration:      duration,
		StartLocation: r.StartLocation,
		EndLocation:   r.EndLocation,
		VehicleID:     r.VehicleID,
		StartMileAge:  r.StartMileAge,
		EndMileAge:    r.EndMileAge,
		ClientTime:    r.ClientTime,
	}
	return out, nil
}

// addRecord 添加一条新的记录
func addRecord(c echo.Context) error {

	uid, _ := c.Get("user").(primitive.ObjectID)

	req := new(addRecordRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	in, err := req.toCreateRecordInput(uid)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	resp, err := service.CreateRecord(in)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, resp)
}

// offlineSyncRecords 离线返回在线状态后记录同步
// 1. 对records按照时间排序，检查相邻两条之间的时间位置是否准确
// 2. 批量更新入数据库
func offlineSyncRecords(c echo.Context) error {

	uid, _ := c.Get("user").(primitive.ObjectID)

	reqs := []*addRecordRequest{}
	if err := c.Bind(&reqs); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	sort.Slice(reqs, func(a, b int) bool {
		return reqs[a].Time.Before(reqs[b].Time)
	})

	l := len(reqs)
	inputs := make([]*service.CreateRecordInput, l)
	for i := 0; i < l; i++ {
		input, err := reqs[i].toCreateRecordInput(uid)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		inputs = append(inputs, input)
	}

	resp, err := service.SyncRecords(inputs)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, resp)

}

// getLatestRecord 获取上一条记录
func getLatestRecord(c echo.Context) error {
	uid, _ := c.Get("user").(primitive.ObjectID)
	resp, err := service.GetLastestRecord(uid)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, resp)
}

// deleteLatestRecord 删除上一条记录
func deleteLatestRecord(c echo.Context) error {

	uid, _ := c.Get("user").(primitive.ObjectID)
	recordID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if err := service.DeleteLatestRecord(&service.DeleteLatestRecordInput{
		RecordID: recordID,
		DriverID: uid,
	}); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, "success")
}

// GetRecordsRequest 获取记录请求参数
type GetRecordsRequest struct {
	DriverID string    `query:"driverID"`
	From     time.Time `query:"from"`
	To       time.Time `query:"to"`
}

// getRecords 获取记录
func getRecords(c echo.Context) error {

	req := new(GetRecordsRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	uid, _ := c.Get("user").(primitive.ObjectID)
	if req.DriverID != uid.Hex() {
		return c.JSON(http.StatusUnauthorized, "driver has no authorization")
	}

	records, err := service.GetRecords(&service.GetRecordsInput{
		DriverID: req.DriverID,
		From:     req.From,
		To:       req.To,
	})
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, records)
}

type addNoteRequest struct {
	service.AddNoteInput `json:",inline"`
}

// addNote 司机为记录添加笔记
func addNote(c echo.Context) error {
	req := new(addNoteRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	uid, _ := c.Get("user").(primitive.ObjectID)
	req.By = uid

	if !req.AddNoteInput.IsDriversRecord(uid) {
		return c.JSON(http.StatusBadRequest, "driver has no authorization")
	}

	note, err := service.AddNote(&req.AddNoteInput)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, note)
}

/************************************************************************************************/
/************************** Transport operator 权限下的record操作 *********************************/
/************************************************************************************************/
type toGetRecordsRequest struct {
	service.GetRecordsInput
}

// toGetRecords 获取记录
func toGetRecords(c echo.Context) error {

	req := new(toGetRecordsRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	var (
		uid primitive.ObjectID // userID
		did primitive.ObjectID // driverID
		tid primitive.ObjectID // toID
		err error
	)
	uid, _ = c.Get("user").(primitive.ObjectID)
	did, err = primitive.ObjectIDFromHex(req.DriverID)
	tid, err = primitive.ObjectIDFromHex(req.TransportOperatorID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if !userApi.CanUserOperatorDriver(uid, did, tid) {
		return c.JSON(http.StatusUnauthorized, "to has no authorization")
	}

	records, err := service.GetRecords(&req.GetRecordsInput)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, records)
}
