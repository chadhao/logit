package model

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Record 记录
// time为生成record的标准时间
// 假如传入time是用户自定义的时间，则需要同时传入clientTime,为用户当前手机时间
type Record struct {
	ID            primitive.ObjectID `bson:"_id" json:"id" valid:"-"`
	DriverID      primitive.ObjectID `bson:"driverID" json:"driverID" valid:"required"`
	Type          Type               `bson:"type" json:"type" valid:"required"`
	Time          time.Time          `bson:"time" json:"time" valid:"required"`
	Duration      time.Duration      `bson:"duration" json:"duration" valid:"required"`
	StartLocation Location           `bson:"startLocation" json:"startLocation" valid:"required"`
	EndLocation   Location           `bson:"endLocation," json:"endLocation" valid:"required"`
	VehicleID     primitive.ObjectID `bson:"vehicleID" json:"vehicleID" valid:"required"`
	StartMileAge  *float64           `bson:"startDistance,omitempty" json:"startDistance,omitempty" valid:"-"`
	EndMileAge    *float64           `bson:"endDistance,omitempty" json:"endDistance,omitempty" valid:"-"`
	CreatedAt     time.Time          `bson:"createdAt" json:"createdAt" valid:"required"`
	ClientTime    *time.Time         `bson:"clientTime,omitempty" json:"clientTime,omitempty" valid:"-"`
	DeletedAt     *time.Time         `bson:"deletedAt,omitempty" json:"deletedAt,omitempty" valid:"-"`
	Active        *bool              `bson:"active,omitempty" json:"active,omitempty" valid:"-"`
}

// Add 记录添加
func (r *Record) Add(lastRec *Record) (err error) {
	// 数据库添加新记录，增加active; 并去除上一条record的active
	db.Client().UseSession(context.TODO(), func(sessionContext mongo.SessionContext) error {
		// 使用事务
		if err := sessionContext.StartTransaction(); err != nil {
			return err
		}

		// 添加记录
		_, err = recordCollection.InsertOne(context.TODO(), r)

		// 去除上一条record的active
		if (Record{}) != *lastRec {
			if err := lastRec.SetActiveStatus(false); err != nil {
				sessionContext.AbortTransaction(sessionContext)
				return err
			}
		}
		return sessionContext.CommitTransaction(sessionContext)
	})
	return nil
}

// SetActiveStatus 记录改变active状态
func (r *Record) SetActiveStatus(active bool) (err error) {
	update := bson.M{"$unset": bson.M{"active": ""}}
	if active {
		update = bson.M{"$set": bson.M{"active": active}}
	}

	if _, err = recordCollection.UpdateOne(context.TODO(), bson.M{"_id": r.ID}, update); err != nil {
		return
	}
	r.Active = &active
	return
}

// Delete 记录删除
func (r *Record) Delete(lastRec *Record) error {

	db.Client().UseSession(context.TODO(), func(sessionContext mongo.SessionContext) error {
		// 使用事务
		if err := sessionContext.StartTransaction(); err != nil {
			return err
		}

		update := bson.M{"$set": bson.M{"deletedAt": time.Now()}, "$unset": bson.M{"active": ""}}
		if _, err := recordCollection.UpdateOne(context.TODO(), bson.M{"_id": r.ID}, update); err != nil {
			return err
		}

		// 为上一条record添加active
		if (Record{}) != *lastRec {
			if err := lastRec.SetActiveStatus(true); err != nil {
				sessionContext.AbortTransaction(sessionContext)
				return err
			}
		}
		return sessionContext.CommitTransaction(sessionContext)
	})
	return nil
}

func (r *Record) isLatestRecord() bool {
	lastest, err := GetLastestRecord(r.DriverID)
	if err != nil {
		return false
	}
	return lastest.ID == r.ID
}

// GetLastestRecord 获取最近的一条记录
func GetLastestRecord(driverID primitive.ObjectID) (*Record, error) {

	lastRec := new(Record)
	err := recordCollection.FindOne(context.TODO(), bson.M{"driverID": driverID, "active": true}).Decode(lastRec)
	// // 按时间排序后找到最近的一条数据
	// if err == mongo.ErrNoDocuments {
	// 	opts := options.FindOne().SetSort(bson.D{{Key: "time", Value: -1}})
	// 	err = recordCollection.FindOne(context.TODO(), bson.M{"driverID": driverID, "deletedAt": nil}, opts).Decode(lastRec)
	// }
	return lastRec, err
}

// GetRecord 通过id获取记录
func GetRecord(id primitive.ObjectID) (*Record, error) {
	r := new(Record)
	err := recordCollection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(r)
	return r, err
}

// GetRecordsOpt 获取用户时间段内的记录选项
type GetRecordsOpt struct {
	From       time.Time
	To         time.Time
	GetDeleted bool
}

// GetRecords 获取用户时间段内的记录
func GetRecords(driverID primitive.ObjectID, opt ...GetRecordsOpt) ([]*Record, error) {
	query := bson.D{primitive.E{Key: "driverID", Value: driverID}}
	if len(opt) == 1 {
		if opt[0].To.IsZero() {
			opt[0].To = time.Now()
		}
		query = append(query, primitive.E{
			Key: "time",
			Value: bson.D{
				primitive.E{Key: "$gte", Value: opt[0].From},
				primitive.E{Key: "$lte", Value: opt[0].To},
			},
		})

		if !opt[0].GetDeleted {
			query = append(query, primitive.E{Key: "deletedAt", Value: nil})
		}
	}

	records := []*Record{}
	cursor, err := recordCollection.Find(context.TODO(), query)
	if err != nil {
		return nil, err
	}
	if err = cursor.All(context.TODO(), &records); err != nil {
		return nil, err
	}
	return records, nil
}

// Records .
type Records []*Record

// SyncAdd 批量上传添加
func (rs Records) SyncAdd(lastRec *Record) error {
	// 数据库添加新记录; 去除上一条record的active
	db.Client().UseSession(context.TODO(), func(sessionContext mongo.SessionContext) error {
		// 使用事务
		if err := sessionContext.StartTransaction(); err != nil {
			return err
		}

		rsI := make([]interface{}, len(rs))
		for i := range rs {
			rsI[i] = rs[i]
		}
		// 数据库添加记录
		if _, err := recordCollection.InsertMany(context.TODO(), rsI); err != nil {
			return err
		}

		// 去除上一条record的active
		if (Record{}) != *lastRec {
			if err := lastRec.SetActiveStatus(false); err != nil {
				sessionContext.AbortTransaction(sessionContext)
				return err
			}
		}
		return sessionContext.CommitTransaction(sessionContext)
	})
	return nil
}
