package model

import (
	"context"
	"errors"
	"time"

	"github.com/chadhao/logit/modules/user/constant"
	"github.com/chadhao/logit/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	// TOIdentity transport operator identity
	TOIdentity string
)

const (
	// TO_ADMIN admin
	TO_ADMIN TOIdentity = "to_admin"
	// TO_SUPER super
	TO_SUPER TOIdentity = "to_super"
	// TO_DRIVER driver
	TO_DRIVER TOIdentity = "to_driver"
)

// GetRole 获取role
func (t TOIdentity) GetRole() int {
	identity := -1
	switch {
	case t == TO_SUPER:
		identity = constant.ROLE_TO_SUPER
	case t == TO_ADMIN:
		identity = constant.ROLE_TO_ADMIN
	case t == TO_DRIVER:
		identity = constant.ROLE_DRIVER
	}
	return identity
}

type (
	// TransportOperator TO组织信息
	TransportOperator struct {
		ID            primitive.ObjectID `json:"id" bson:"_id"`
		LicenseNumber string             `json:"licenseNumber" bson:"licenseNumber"`
		Name          string             `json:"name" bson:"name"`
		IsVerified    bool               `json:"isVerified" bson:"isVerified"`
		IsCompany     bool               `json:"isCompany" bson:"isCompany"`
		CreatedAt     time.Time          `json:"createdAt" bson:"createdAt"`
	}

	// TransportOperatorIdentity TO用户身份信息
	TransportOperatorIdentity struct {
		ID                  primitive.ObjectID `json:"id" bson:"_id"`
		UserID              primitive.ObjectID `json:"userID" bson:"userID"`
		TransportOperatorID primitive.ObjectID `json:"transportOperatorID" bson:"transportOperatorID"`
		Identity            TOIdentity         `json:"identity" bson:"identity"`
		Contact             *string            `json:"contact" bson:"contact"`
		CreatedAt           time.Time          `json:"createdAt" bson:"createdAt"`
	}

	// TransportOperatorIdentityDetail TO用户身份信息和其所属组织信息
	TransportOperatorIdentityDetail struct {
		ID                  primitive.ObjectID `json:"id" bson:"_id"`
		UserID              primitive.ObjectID `json:"userID" bson:"userID"`
		TransportOperatorID primitive.ObjectID `json:"transportOperatorID" bson:"transportOperatorID"`
		Identity            TOIdentity         `json:"identity" bson:"identity"`
		Contact             *string            `json:"contact" bson:"contact"`
		CreatedAt           time.Time          `json:"createdAt" bson:"createdAt"`
		TransportOperator   *TransportOperator `json:"transportOperator" bson:"transportOperator"`
	}
)

// Create 创建TO组织信息
func (t *TransportOperator) Create(user *User) error {
	db.Client().UseSession(context.TODO(), func(sessionContext mongo.SessionContext) error {
		// 使用事务
		if err := sessionContext.StartTransaction(); err != nil {
			return err
		}
		// 创建TO组织信息
		if _, err := toCollection.InsertOne(context.TODO(), t); err != nil {
			return err
		}
		// 该TO创建super身份
		superIdentity := &TransportOperatorIdentity{
			ID:                  primitive.NewObjectID(),
			UserID:              user.ID,
			TransportOperatorID: t.ID,
			Identity:            TO_SUPER,
			CreatedAt:           time.Now(),
		}
		if err := superIdentity.Create(user); err != nil {
			sessionContext.AbortTransaction(sessionContext)
			return err
		}
		// 自雇性质时，当user已经有driver信息时，自动添加为TO下driver身份
		if !t.IsCompany {
			if _, err := FindDriver(FindDriverOpt{ID: user.ID}); err == nil {
				driverIdentity := &TransportOperatorIdentity{
					ID:                  primitive.NewObjectID(),
					UserID:              user.ID,
					TransportOperatorID: t.ID,
					Identity:            TO_DRIVER,
					CreatedAt:           time.Now(),
				}
				if err := driverIdentity.Create(user); err != nil {
					sessionContext.AbortTransaction(sessionContext)
					return err
				}
			}
		}
		return sessionContext.CommitTransaction(sessionContext)
	})
	return nil
}

// TransportOperatorFindOpt 查询TO组织选项
type TransportOperatorFindOpt struct {
	ID            primitive.ObjectID
	LicenseNumber string
}

// TransportOperatorFind 查询TO组织
func TransportOperatorFind(opt TransportOperatorFindOpt) (*TransportOperator, error) {
	conditions := bson.D{}
	if !opt.ID.IsZero() {
		conditions = append(conditions, primitive.E{Key: "_id", Value: opt.ID})
	}
	if len(opt.LicenseNumber) > 0 {
		conditions = append(conditions, primitive.E{Key: "licenseNumber", Value: opt.LicenseNumber})
	}
	query := bson.D{primitive.E{Key: "$or", Value: conditions}}

	to := &TransportOperator{}
	err := toCollection.FindOne(context.TODO(), query).Decode(to)
	return to, err
}

// Update 更新TO组织信息
func (t *TransportOperator) Update() error {
	filter := bson.D{primitive.E{Key: "_id", Value: t.ID}}
	if result, _ := toCollection.ReplaceOne(context.TODO(), filter, t); result.MatchedCount != 1 {
		return errors.New("transportOperator not updated")
	}
	return nil
}

// Delete 删除TO组织
func (t *TransportOperator) Delete() error {
	filter := bson.D{primitive.E{Key: "_id", Value: t.ID}}
	_, err := toCollection.DeleteOne(context.TODO(), filter)
	return err
}

// TransportOperatorFilterOpt TO组织过滤选项
type TransportOperatorFilterOpt struct {
	IsVerified    *bool
	IsCompany     *bool
	LicenseNumber string
	Name          string
}

// TransportOperatorFilter TO组织过滤
func TransportOperatorFilter(opt TransportOperatorFilterOpt) ([]*TransportOperator, error) {

	tos := []*TransportOperator{}

	filter := bson.M{}
	if opt.IsVerified != nil {
		filter["isVerified"] = *opt.IsVerified
	}
	if opt.IsCompany != nil {
		filter["isCompany"] = *opt.IsCompany
	}
	if len(opt.LicenseNumber) > 0 {
		filter["licenseNumber"] = primitive.Regex{Pattern: opt.LicenseNumber, Options: "i"}
	}
	if len(opt.Name) > 0 {
		filter["name"] = primitive.Regex{Pattern: opt.Name, Options: "i"}
	}

	cursor, err := toCollection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	if err = cursor.All(context.TODO(), &tos); err != nil {
		return nil, err
	}
	return tos, nil
}

// TransportOperatorIdentityFilterOpt TO组织下的用户身份信息过滤选项
type TransportOperatorIdentityFilterOpt struct {
	UserID              primitive.ObjectID
	TransportOperatorID primitive.ObjectID
	Identity            TOIdentity
}

// TransportOperatorIdentityFilter TO组织下的用户身份信息过滤
func TransportOperatorIdentityFilter(opt TransportOperatorIdentityFilterOpt) ([]*TransportOperatorIdentityDetail, error) {
	// 检索条件
	filter := bson.M{}
	if !opt.UserID.IsZero() {
		filter["userID"] = opt.UserID
	}
	if !opt.TransportOperatorID.IsZero() {
		filter["transportOperatorID"] = opt.TransportOperatorID
	}
	if opt.Identity != "" {
		filter["identity"] = opt.Identity
	}
	// 按条件查询用户身份并获取对应的组织信息
	query := mongo.Pipeline{
		bson.D{
			primitive.E{Key: "$lookup", Value: bson.D{
				primitive.E{Key: "from", Value: "transportOperator"},
				primitive.E{Key: "localField", Value: "transportOperatorID"},
				primitive.E{Key: "foreignField", Value: "_id"},
				primitive.E{Key: "as", Value: "transportOperator"},
			}},
		},
		bson.D{
			primitive.E{Key: "$match", Value: filter},
		},
		bson.D{
			primitive.E{Key: "$unwind", Value: "$transportOperator"},
		},
	}

	cursor, err := toICollection.Aggregate(context.TODO(), query)
	if err != nil {
		return nil, err
	}

	tos := []*TransportOperatorIdentityDetail{}
	if err = cursor.All(context.TODO(), &tos); err != nil {
		return nil, err
	}
	return tos, nil
}

// Create 创建TO用户身份
func (t *TransportOperatorIdentity) Create(user *User) error {

	db.Client().UseSession(context.TODO(), func(sessionContext mongo.SessionContext) error {
		// 使用事务
		if err := sessionContext.StartTransaction(); err != nil {
			return err
		}
		// 创建TO用户身份
		_, err := toICollection.InsertOne(context.TODO(), t)
		if err != nil {
			return err
		}
		// 更新用户信息
		roleLen, roles := len(user.RoleIDs), utils.RolesAssert(user.RoleIDs)
		switch {
		case t.Identity == TO_SUPER && !roles.Is(constant.ROLE_TO_SUPER):
			user.RoleIDs = append(user.RoleIDs, constant.ROLE_TO_SUPER)
		case t.Identity == TO_ADMIN && !roles.Is(constant.ROLE_TO_ADMIN):
			user.RoleIDs = append(user.RoleIDs, constant.ROLE_TO_ADMIN)
		case t.Identity == TO_DRIVER && !roles.Is(constant.ROLE_DRIVER):
			user.RoleIDs = append(user.RoleIDs, constant.ROLE_DRIVER)
		}

		if roleLen != len(user.RoleIDs) {
			if err := user.Update(); err != nil {
				sessionContext.AbortTransaction(sessionContext)
				return err
			}
		}
		return sessionContext.CommitTransaction(sessionContext)
	})
	return nil
}

// Delete 删除TO用户身份
func (t *TransportOperatorIdentity) Delete(user *User) error {

	db.Client().UseSession(context.TODO(), func(sessionContext mongo.SessionContext) error {
		// 使用事务
		if err := sessionContext.StartTransaction(); err != nil {
			return err
		}
		if _, err := toICollection.DeleteOne(context.TODO(), bson.M{"_id": t.ID}); err != nil {
			return err
		}
		// 删除后需要检查角色的role是否需要更新
		toIs, err := TransportOperatorIdentityFilter(TransportOperatorIdentityFilterOpt{
			UserID:   t.UserID,
			Identity: t.Identity,
		})
		if err != nil {
			sessionContext.AbortTransaction(sessionContext)
			return err
		}
		if len(toIs) == 0 {
			roles := utils.RolesAssert(user.RoleIDs)
			identity := t.Identity.GetRole()
			for i := 0; i < len(roles); i++ {
				if roles[i] == identity {
					roles = append(roles[:i], roles[i+1:]...)
					user.RoleIDs = roles
					if err := user.Update(); err != nil {
						sessionContext.AbortTransaction(sessionContext)
						return err
					}
					break
				}
			}
		}

		return sessionContext.CommitTransaction(sessionContext)
	})
	return nil
}

// TransportOperatorIdentityFind 查询TO组织的角色信息
func TransportOperatorIdentityFind(TransportOperatorIdentityID primitive.ObjectID) (*TransportOperatorIdentity, error) {
	out := &TransportOperatorIdentity{}
	err := toICollection.FindOne(context.TODO(), bson.M{"_id": TransportOperatorIdentityID}).Decode(out)
	return out, err
}

// UpdateContact 更新联系方式
func (t *TransportOperatorIdentity) UpdateContact(contact string) error {
	update := bson.M{"$set": bson.M{"contact": contact}}
	_, err := db.Collection("transportOperatorIdentity").UpdateOne(context.TODO(), bson.M{"_id": t.ID}, update)
	return err
}

// TransportOperatorIdentityExists 用户在改TO组织下的身份是否存在
type TransportOperatorIdentityExists struct {
	UserID              primitive.ObjectID
	TransportOperatorID primitive.ObjectID
	Identity            []TOIdentity // 如果此项长度大于1，则用户身份属于此种任意一中，都返回True
}

// IsTransportOperatorIdentityExists 检查用户在改TO组织下的身份是否存在
func IsTransportOperatorIdentityExists(opt TransportOperatorIdentityExists) bool {

	conditions := primitive.A{
		bson.D{primitive.E{Key: "userID", Value: opt.UserID}},
		bson.D{primitive.E{Key: "transportOperatorID", Value: opt.TransportOperatorID}},
		bson.D{primitive.E{Key: "identity", Value: primitive.E{Key: "$in", Value: opt.Identity}}},
		bson.D{primitive.E{Key: "deletedAt", Value: nil}},
	}

	filter := bson.D{primitive.E{Key: "$and", Value: conditions}}

	count, _ := toICollection.CountDocuments(context.TODO(), filter)
	return count > 0
}

// GetIdentitiesByUserIDs 批量通过userIDs获取他们相关的identity信息
func GetIdentitiesByUserIDs(uids []primitive.ObjectID) (map[primitive.ObjectID][]*TransportOperatorIdentityDetail, error) {
	m := make(map[primitive.ObjectID][]*TransportOperatorIdentityDetail)
	var tois []*TransportOperatorIdentityDetail
	filter := bson.M{
		"userID": bson.M{"$in": uids},
	}

	query := mongo.Pipeline{
		bson.D{
			primitive.E{Key: "$lookup", Value: bson.D{
				primitive.E{Key: "from", Value: "transportOperator"},
				primitive.E{Key: "localField", Value: "transportOperatorID"},
				primitive.E{Key: "foreignField", Value: "_id"},
				primitive.E{Key: "as", Value: "transportOperator"},
			}},
		},
		bson.D{
			primitive.E{Key: "$match", Value: filter},
		},
		bson.D{
			primitive.E{Key: "$unwind", Value: "$transportOperator"},
		},
	}
	cursor, err := toICollection.Aggregate(context.TODO(), query)
	if err != nil {
		return nil, err
	}
	if err = cursor.All(context.TODO(), &tois); err != nil {
		return nil, err
	}
	for _, toi := range tois {
		m[toi.UserID] = append(m[toi.UserID], toi)
	}
	return m, nil
}
