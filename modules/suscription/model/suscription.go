package model

import (
	"context"
	"errors"
	"time"

	valid "github.com/asaskevich/govalidator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Suscription 用户的订阅状态
type Suscription struct {
	DriverID    primitive.ObjectID `bson:"_id" json:"driverID" valid:"required"`
	Renew       bool               `bson:"renew" json:"renew" valid:"required"`
	ExpiredDate string             `bson:"expiredDate" json:"expiredDate" valid:"required"`
	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt" valid:"required"`
}

// Add 创建新的用户订阅状态
func (s *Suscription) Add() (err error) {
	s.ExpiredDate = time.Now().In(loc).Format("2006-01-02")
	if _, err = valid.ValidateStruct(s); err != nil {
		return
	}
	_, err = suscriptionCollection.InsertOne(context.TODO(), s)
	return
}

func (s *Suscription) changeExpiredDate(expire string) (err error) {
	update := bson.M{"$set": bson.M{"expiredAt": expire}}
	_, err = suscriptionCollection.UpdateOne(context.TODO(), bson.M{"_id": s.DriverID}, update)
	return
}

// IsExpired 用户订阅状态是否过期
func (s *Suscription) IsExpired() bool {
	t, _ := time.ParseInLocation("2006-01-02", s.ExpiredDate, loc)
	return t.AddDate(0, 0, 1).Before(time.Now())
}

// Record 订阅记录
type Record struct {
	ID        primitive.ObjectID `bson:"_id" json:"id" valid:"required"`
	DriverID  primitive.ObjectID `bson:"driverID" json:"driverID" valid:"required"`
	Charge    int64              `bson:"charge" json:"charge" valid:"required"`
	Refund    int64              `bson:"refund" json:"refund" valid:"-"`
	StartDate string             `bson:"startDate" json:"startDate" valid:"required"`
	EndDate   string             `bson:"endDate" json:"endDate" valid:"required"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt" valid:"required"`
}

// Add 创建一条订阅记录
func (r *Record) Add(days int) (err error) {

	// 获取订阅状态，如果过期时间已过期，则初始日期为当前日期;
	// 如果过期时间未过期，则初始日期为过期日期
	// 结束日期为初始日期加上添加天数
	s, err := GetSuscription(r.DriverID)
	if err != nil {
		return
	}

	expire, err := time.ParseInLocation("2006-01-02", s.ExpiredDate, loc)
	if err != nil {
		return
	}
	r.StartDate = s.ExpiredDate
	if expire.Before(time.Now()) {
		r.StartDate = time.Now().In(loc).Format("2006-01-02")
	}

	start, err := time.Parse("2006-01-02", r.StartDate)
	if err != nil {
		return
	}
	r.EndDate = start.AddDate(0, 0, days).In(loc).Format("2006-01-02")

	// 验证结构并存入数据库
	if _, err = valid.ValidateStruct(r); err != nil {
		return
	}
	if _, err = recordCollection.InsertOne(context.TODO(), r); err != nil {
		return
	}

	// 改变该用户订阅状态的过期时间，为结束时间
	if err = s.changeExpiredDate(r.EndDate); err != nil {
		return
	}
	return
}

// MakeRefund 订阅记录退款
func (r *Record) MakeRefund(refund int64) (err error) {

	if refund > r.Charge {
		err = errors.New("cannot refund more than charge")
		return
	}

	update := bson.M{"$set": bson.M{"refund": refund}}
	_, err = recordCollection.UpdateOne(context.TODO(), bson.M{"_id": r.ID}, update)
	return

}

// GetSuscription 获取用户的订阅状态
func GetSuscription(driverID primitive.ObjectID) (*Suscription, error) {
	s := &Suscription{}
	err := suscriptionCollection.FindOne(context.TODO(), bson.M{"_id": driverID}).Decode(s)
	return s, err
}

// GetRecord 通过id获取用户的订阅记录
func GetRecord(id primitive.ObjectID) (*Record, error) {
	r := &Record{}
	err := recordCollection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(r)
	return r, err
}

// GetRecords 获取用户的订阅记录
func GetRecords(driverID primitive.ObjectID) ([]Record, error) {
	records := []Record{}

	cursor, err := recordCollection.Find(context.TODO(), bson.M{"driverID": driverID})
	if err != nil {
		return nil, err
	}
	if err = cursor.All(context.TODO(), &records); err != nil {
		return nil, err
	}
	return records, nil
}
