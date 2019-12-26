package model

import (
	"context"
	"errors"
	"time"

	valid "github.com/asaskevich/govalidator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"googlemaps.github.io/maps"
)

type (
	// Coors represents a location on the Earth.
	Coors struct {
		Lat float64 `bson:"lat,omitempty" json:"lat,omitempty" valid:"required"`
		Lng float64 `bson:"lng,omitempty" json:"lng,omitempty" valid:"required"`
	}
	// Location 位置信息
	Location struct {
		Address string `bson:"address" json:"address" valid:"required"`
		Coors   Coors  `bson:"coors,omitempty" json:"coors,omitempty" valid:"required"`
	}
	// DrivingLoc 司机在行驶过程中的位置信息
	DrivingLoc struct {
		ID        primitive.ObjectID `bson:"_id" json:"id" valid:"-"`
		UserID    primitive.ObjectID `bson:"userID" json:"userID" valid:"required"`
		CreatedAt time.Time          `bson:"createdAt" json:"createdAt" valid:"required"`
		Location  `bson:",inline" json:",inline"`
	}
)

// EmptyCoors 判断coors是否为空
func (l *Location) EmptyCoors() bool {
	return l.Coors.Lat == 0
}

// Save 保存司机行驶的位置信息到数据库
func (dLoc *DrivingLoc) Save() error {
	if dLoc.CreatedAt.IsZero() {
		dLoc.CreatedAt = time.Now()
	}
	if _, err := valid.ValidateStruct((dLoc)); err != nil {
		return err
	}
	_, err := drivingLocCol.InsertOne(context.TODO(), dLoc)
	return err
}

// GetDrivingLocs 通过userID和指定时间段返回司机行驶位置信息
func GetDrivingLocs(userID primitive.ObjectID, from, to time.Time) ([]DrivingLoc, error) {
	drivingLocs := []DrivingLoc{}
	query := bson.M{
		"userID": userID,
		"gte":    bson.M{"createdAt": from},
		"lte":    bson.M{"createdAt": to},
	}
	cursor, err := drivingLocCol.Find(context.TODO(), query)
	if err != nil {
		return nil, err
	}
	if err = cursor.All(context.TODO(), &drivingLocs); err != nil {
		return nil, err
	}
	return drivingLocs, nil
}

// GetCoorsFromAddr 通过位置获取坐标信息
func GetCoorsFromAddr(addr string) (coors Coors, err error) {
	findPlaceReq := &maps.FindPlaceFromTextRequest{
		Input:     addr,
		InputType: maps.FindPlaceFromTextInputTypeTextQuery,
	}
	findPlaceResp, err := mapClient.FindPlaceFromText(context.TODO(), findPlaceReq)
	if err != nil {
		return
	}
	if len(findPlaceResp.Candidates) == 0 {
		err = errors.New("can not find any match address")
		return
	}
	latlng := findPlaceResp.Candidates[0].Geometry.Location
	coors = Coors{Lat: latlng.Lat, Lng: latlng.Lng}
	return
}
