package api

import (
	"errors"
	"time"

	"github.com/chadhao/logit/modules/location/model"
	"go.mongodb.org/mongo-driver/bson/primitive"

	valid "github.com/asaskevich/govalidator"
)

// reqAddDrivingLoc 添加行驶信息请求结构
type reqAddDrivingLoc struct {
	model.Coors `json:",inline"`
	CreatedAt   time.Time `json:"createdAt,omitempty"`
}

func (reqAdd *reqAddDrivingLoc) constructToDrivingLoc(driverID primitive.ObjectID) (*model.DrivingLoc, error) {
	// 如果location中的coors为空，则需要请求获取
	drivingLoc := &model.DrivingLoc{
		DriverID:  driverID,
		Coors:     reqAdd.Coors,
		CreatedAt: reqAdd.CreatedAt,
	}
	return drivingLoc, nil
}

// reqDrivingLocs 行驶信息请求结构
type reqDrivingLocs struct {
	DriverID string    `json:"driverID" query:"driverID" valid:"required"`
	From     time.Time `json:"from" query:"from" valid:"required"`
	To       time.Time `json:"to" query:"to" valid:"optional"`
}

func (req *reqDrivingLocs) valid() error {
	if _, err := valid.ValidateStruct(req); err != nil {
		return err
	}
	if req.To.IsZero() {
		req.To = time.Now()
	}
	if req.From.After(req.To) {
		return errors.New("times order is wrong")
	}
	return nil
}

func (req *reqDrivingLocs) getDrivingLocs() ([]model.DrivingLoc, error) {
	if err := req.valid(); err != nil {
		return nil, err
	}
	driverID, err := primitive.ObjectIDFromHex(req.DriverID)
	if err != nil {
		return nil, err
	}
	drivingLocs, err := model.GetDrivingLocs(driverID, req.From, req.To)
	if err != nil {
		return nil, err
	}
	return drivingLocs, err
}
