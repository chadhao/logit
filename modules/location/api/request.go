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
	model.Location `json:",inline"`
	CreatedAt      time.Time `json:"createdAt,omitempty"`
}

func (reqAdd *reqAddDrivingLoc) constructToDrivingLoc(userID primitive.ObjectID) (*model.DrivingLoc, error) {
	// 如果location中的coors为空，则需要请求获取
	if reqAdd.Location.EmptyCoors() {
		if coors, err := model.GetCoorsFromAddr(reqAdd.Address); err == nil {
			reqAdd.Coors = coors
		}
	}
	drivingLoc := &model.DrivingLoc{
		UserID:    userID,
		Location:  reqAdd.Location,
		CreatedAt: reqAdd.CreatedAt,
	}
	return drivingLoc, nil
}

// reqDrivingLocs 行驶信息请求结构
type reqDrivingLocs struct {
	From time.Time `json:"from" valid:"required"`
	To   time.Time `json:"to" valid:"-"`
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

func (req *reqDrivingLocs) getDrivingLocs(userID primitive.ObjectID) ([]model.DrivingLoc, error) {
	if err := req.valid(); err != nil {
		return nil, err
	}
	drivingLocs, err := model.GetDrivingLocs(userID, req.From, req.To)
	if err != nil {
		return nil, err
	}
	return drivingLocs, err
}
