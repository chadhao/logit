package internal

import (
	"time"

	"github.com/chadhao/logit/modules/location/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ReqAddDrivingLoc 添加行驶信息请求结构
type ReqAddDrivingLoc struct {
	model.Location `json:",inline"`
	CreatedAt      time.Time `json:"createdAt,omitempty"`
}

func (loc *ReqAddDrivingLoc) constructToDrivingLoc(userID primitive.ObjectID) (*model.DrivingLoc, error) {
	// 如果location中的coors为空，则需要请求获取
	if loc.Location.EmptyCoors() {
		if coors, err := model.GetCoorsFromAddr(loc.Address); err == nil {
			loc.Coors = coors
		}
	}
	drivingLoc := &model.DrivingLoc{
		UserID:    userID,
		Location:  loc.Location,
		CreatedAt: loc.CreatedAt,
	}
	return drivingLoc, nil
}
