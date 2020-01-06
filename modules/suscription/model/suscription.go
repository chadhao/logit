package model

import (
	"context"
	"time"

	valid "github.com/asaskevich/govalidator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Suscriber 生成或取消订阅记录的接口
type Suscriber interface {
	Suscribe()
	Unsuscribe()
}

// Suscription 用户的订阅状态
type Suscription struct {
	UserID    primitive.ObjectID `bson:"_id" json:"userID" valid:"required"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt" valid:"required"`
	ExpiredAt time.Time          `bson:"expiredAt" json:"expiredAt" valid:"required"`
}

// Create 创建一条新的用户订阅状态
func (s *Suscription) Create() (err error) {

	s.CreatedAt = time.Now()
	s.ExpiredAt = time.Time{}

	if _, err = valid.ValidateStruct(s); err != nil {
		return
	}

	_, err = suscriptionCollection.InsertOne(context.TODO(), s)
	return
}

func (s *Suscription) expiredDaysAdd(days int) (err error) {
	// 添加的时间天数可为正也可为负，为负时则减去时间。
	// 如果过期时间在当前时间之前，则先将过期时间处理为当前时间。
	if s.ExpiredAt.Before(time.Now()) {
		s.ExpiredAt = time.Now()
	}
	expire := s.ExpiredAt.AddDate(0, 0, days)
	update := bson.M{"$set": bson.M{"expiredAt": expire}}
	_, err = suscriptionCollection.UpdateOne(context.TODO(), bson.M{"_id": s.UserID}, update)
	return
}

type (
	// GeneratedFrom 订阅生成
	GeneratedFrom struct {
		FromID  primitive.ObjectID `bson:"id" json:"id" valid:"required"`
		DaysAdd int                `bson:"daysAdd" json:"daysAdd" valid:"required"`
	}
	// Record 订阅记录
	Record struct {
		ID            primitive.ObjectID `bson:"_id" json:"id" valid:"-"`
		UserID        primitive.ObjectID `bson:"userID" json:"userID" valid:"required"`
		GeneratedFrom GeneratedFrom      `bson:"generatedFrom" json:"generatedFrom" valid:"required"`
		CreatedAt     time.Time          `bson:"createdAt" json:"createdAt" valid:"required"`
		UnsuscribedAt *time.Time         `bson:"unsuscribedAt,omitempty" json:"unsuscribedAt,omitempty" valid:"-"`
	}
)

// Create 创建一条订阅记录
func (r *Record) Create() error {
	// 添加订阅记录
	if _, err := valid.ValidateStruct(r); err != nil {
		return err
	}

	s, err := GetSuscription(r.UserID)
	if err != nil {
		return err
	}

	if _, err := recordCollection.InsertOne(context.TODO(), r); err != nil {
		return err
	}

	// 改变该用户订阅状态的时间
	if err := s.expiredDaysAdd(r.GeneratedFrom.DaysAdd); err != nil {
		return err
	}
	return nil
}

// Delete 取消一条订阅记录
func (r *Record) Delete() error {

	s, err := GetSuscription(r.UserID)
	if err != nil {
		return err
	}

	deletedAt := time.Now()
	r.UnsuscribedAt = &deletedAt
	update := bson.M{"$set": bson.M{"unsuscribedAt": deletedAt}}
	if _, err = recordCollection.UpdateOne(context.TODO(), bson.M{"_id": r.ID}, update); err != nil {
		return err
	}

	// 改变该用户订阅状态的时间
	if err := s.expiredDaysAdd(-r.GeneratedFrom.DaysAdd); err != nil {
		return err
	}
	return nil
}

// GetSuscription 获取用户的订阅状态
func GetSuscription(userID primitive.ObjectID) (*Suscription, error) {
	s := new(Suscription)
	err := suscriptionCollection.FindOne(context.TODO(), bson.M{"_id": userID}).Decode(s)
	return s, err
}

// GetRecord 通过id获取用户的订阅记录
func GetRecord(id primitive.ObjectID) (*Record, error) {
	r := new(Record)
	err := recordCollection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(r)
	return r, err
}

// GetRecords 获取用户的订阅记录
func GetRecords(userID primitive.ObjectID) ([]Record, error) {
	records := []Record{}

	cursor, err := recordCollection.Find(context.TODO(), bson.M{"userID": userID})
	if err != nil {
		return nil, err
	}
	if err = cursor.All(context.TODO(), &records); err != nil {
		return nil, err
	}
	return records, nil
}
