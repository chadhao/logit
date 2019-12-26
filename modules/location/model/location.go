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
		Lat float64 `bson:"lat" json:"lat" valid:"required"`
		Lng float64 `bson:"lng" json:"lng" valid:"required"`
	}

	// DrivingLoc 司机在行驶过程中的位置信息
	DrivingLoc struct {
		ID        primitive.ObjectID `bson:"_id" json:"id" valid:"-"`
		UserID    primitive.ObjectID `bson:"userID" json:"userID" valid:"required"`
		CreatedAt time.Time          `bson:"createdAt" json:"createdAt" valid:"required"`
		Coors     Coors              `bson:"coors" json:"coors" valid:"required"`
	}

	// Address 位置信息
	Address string
)

// EmptyCoors 判断coors是否为空
func (dLoc *DrivingLoc) emptyCoors() bool {
	return dLoc.Coors.Lat == 0
}

// Save 保存司机行驶的位置信息到数据库
func (dLoc *DrivingLoc) Save() error {
	if dLoc.emptyCoors() {
		return errors.New("coors cannot be null")
	}
	if dLoc.CreatedAt.IsZero() {
		dLoc.CreatedAt = time.Now()
	}
	if _, err := valid.ValidateStruct(dLoc); err != nil {
		return err
	}
	_, err := drivingLocCol.InsertOne(context.TODO(), dLoc)
	return err
}

// GetCoorsFromAddr 通过位置获取坐标信息
func (addr Address) GetCoorsFromAddr() (coors Coors, err error) {
	req := &maps.GeocodingRequest{
		Address: string(addr),
	}
	resp, err := mapClient.Geocode(context.TODO(), req)
	if err != nil {
		return
	}
	if len(resp) == 0 {
		err = errors.New("can not find any match address")
		return
	}
	coors = Coors{Lat: resp[0].Geometry.Location.Lat, Lng: resp[0].Geometry.Location.Lng}
	return
}

// GetAddrFromCoors 通过坐标信息获取位置
func (coors Coors) GetAddrFromCoors() (addr Address, err error) {
	req := &maps.GeocodingRequest{
		LatLng: &maps.LatLng{Lat: coors.Lat, Lng: coors.Lng},
	}
	resp, err := mapClient.Geocode(context.TODO(), req)
	if err != nil {
		return
	}
	if len(resp) == 0 {
		err = errors.New("can not find any match address")
		return
	}
	addr = Address(resp[0].FormattedAddress)
	return
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
