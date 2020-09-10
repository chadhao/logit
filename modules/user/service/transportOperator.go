package service

import (
	"errors"
	"time"

	valid "github.com/asaskevich/govalidator"
	"github.com/chadhao/logit/config"
	"github.com/chadhao/logit/modules/user/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	// TransportOperatorRegisterInput TO组织注册参数
	TransportOperatorRegisterInput struct {
		Conf          config.Config
		UserID        primitive.ObjectID
		LicenseNumber string `valid:"required"`
		IsCompany     bool
		Name          string `valid:"required"`
		Contact       *string
	}
)

func (to *TransportOperatorRegisterInput) toTransportOperator() *model.TransportOperator {
	return &model.TransportOperator{
		ID:            primitive.NewObjectID(),
		LicenseNumber: to.LicenseNumber,
		Name:          to.Name,
		IsVerified:    false,
		IsCompany:     to.IsCompany,
		CreatedAt:     time.Now(),
	}
}

// TransportOperatorRegister TO组织注册
func TransportOperatorRegister(in *TransportOperatorRegisterInput) (*IssueTokenOutput, error) {
	// 参数验证
	if _, err := valid.ValidateStruct(in); err != nil {
		return nil, err
	}
	if in.IsCompany && in.Contact == nil {
		return nil, errors.New("contact is required")
	}

	// 查询用户
	user, err := model.FindUser(model.FindUserOpt{UserID: in.UserID})
	if err != nil {
		return nil, err
	}

	// TO组织注册
	to := in.toTransportOperator()
	if _, err := model.TransportOperatorFind(model.TransportOperatorFindOpt{LicenseNumber: to.LicenseNumber}); err != mongo.ErrNoDocuments {
		return nil, errors.New("license has been used")
	}

	if err := to.Create(user); err != nil {
		return nil, err
	}

	// 发放新的token
	issueTokenOutput, err := IssueToken(&IssueTokenInput{UserID: user.ID, RoleIDs: user.RoleIDs, Conf: in.Conf})
	if err != nil {
		return nil, err
	}
	return issueTokenOutput, nil
}

type (
	// TransportOperatorFindInput TO组织查询参数
	TransportOperatorFindInput struct {
		TransportOperatorID primitive.ObjectID
	}
	// TransportOperatorFindOutput TO组织查询返回参数
	TransportOperatorFindOutput struct {
		*model.TransportOperator
	}
)

// TransportOperatorFind TO组织查询
func TransportOperatorFind(in *TransportOperatorFindInput) (*TransportOperatorFindOutput, error) {
	to, err := model.TransportOperatorFind(model.TransportOperatorFindOpt{
		ID: in.TransportOperatorID,
	})
	if err != nil {
		return nil, err
	}

	return &TransportOperatorFindOutput{to}, nil
}

type (
	// TransportOperatorUpdateInput TO组织更新参数
	TransportOperatorUpdateInput struct {
		UserID              primitive.ObjectID
		TransportOperatorID primitive.ObjectID
		LicenseNumber       string
		Name                string
	}
	// TransportOperatorUpdateOutput TO组织更新返回参数
	TransportOperatorUpdateOutput struct {
		*model.TransportOperator
	}
)

// TransportOperatorUpdate TO组织更新
func TransportOperatorUpdate(in *TransportOperatorUpdateInput) (*TransportOperatorUpdateOutput, error) {
	// 检查是否有权限更新
	if !model.IsTransportOperatorIdentityExists(model.TransportOperatorIdentityExists{
		UserID:              in.UserID,
		TransportOperatorID: in.TransportOperatorID,
		Identity:            []model.TOIdentity{model.TO_SUPER},
	}) {
		return nil, errors.New("no authorization")
	}
	// 获取TO
	to, err := model.TransportOperatorFind(model.TransportOperatorFindOpt{ID: in.TransportOperatorID})
	if err != nil {
		return nil, err
	}
	// 更新TO
	if len(in.LicenseNumber) > 0 {
		to.LicenseNumber = in.LicenseNumber
	}
	if len(in.Name) > 0 {
		to.Name = in.Name
	}

	if err := to.Update(); err != nil {
		return nil, err
	}
	return &TransportOperatorUpdateOutput{to}, nil
}

type (
	// TransportOperatorsFindInput TO组织查询参数
	TransportOperatorsFindInput struct {
		IsVerified    *bool
		IsCompany     *bool
		LicenseNumber string
		Name          string
	}
	// TransportOperatorsFindOutput TO组织查询返回参数
	TransportOperatorsFindOutput struct {
		Tos []*model.TransportOperator
	}
)

// TransportOperatorsFind TO组织查询
func TransportOperatorsFind(in *TransportOperatorsFindInput) (*TransportOperatorsFindOutput, error) {
	tos, err := model.TransportOperatorFilter(model.TransportOperatorFilterOpt{
		IsCompany:     in.IsCompany,
		IsVerified:    in.IsVerified,
		LicenseNumber: in.LicenseNumber,
		Name:          in.Name,
	})
	if err != nil {
		return nil, err
	}
	return &TransportOperatorsFindOutput{Tos: tos}, nil
}

type (
	// TransportOperatorIdentityAddInput TO组织添加角色身份参数
	TransportOperatorIdentityAddInput struct {
		TransportOperatorID primitive.ObjectID
		UserID              primitive.ObjectID
		Identity            model.TOIdentity
		Contact             *string
	}
	// TransportOperatorIdentityAddOutput TO组织添加角色身份返回参数
	TransportOperatorIdentityAddOutput struct {
		*model.TransportOperatorIdentity
	}
)

// TransportOperatorIdentityAdd TO组织添加角色身份
func TransportOperatorIdentityAdd(in *TransportOperatorIdentityAddInput) (*TransportOperatorIdentityAddOutput, error) {

	toIdentity := &model.TransportOperatorIdentity{
		ID:                  primitive.NewObjectID(),
		UserID:              in.UserID,
		TransportOperatorID: in.TransportOperatorID,
		Identity:            in.Identity,
		CreatedAt:           time.Now(),
	}

	// 检查身份是否已经存在
	if model.IsTransportOperatorIdentityExists(model.TransportOperatorIdentityExists{
		UserID:              in.UserID,
		TransportOperatorID: in.TransportOperatorID,
		Identity:            []model.TOIdentity{in.Identity},
	}) {
		return nil, errors.New("identity exists")
	}

	// 获取user信息
	user, err := model.FindUser(model.FindUserOpt{UserID: in.UserID})
	if err != nil {
		return nil, err
	}

	// 获取TO组织
	to, err := model.TransportOperatorFind(model.TransportOperatorFindOpt{ID: in.TransportOperatorID})
	if err != nil {
		return nil, nil
	}

	// to super, to admin 用户在公司情况下需要填入contact信息
	if to.IsCompany {
		if in.Contact == nil && in.Identity != model.TO_DRIVER {
			return nil, errors.New("contact is required")
		}
		toIdentity.Contact = in.Contact
	} else {
		// 自雇形式下的验证
		if toIdentity.Identity == model.TO_DRIVER {
			// 添加为driver权限
			// 只能添加自己
			tos, _ := model.TransportOperatorIdentityFilter(model.TransportOperatorIdentityFilterOpt{
				UserID:              toIdentity.UserID,
				TransportOperatorID: toIdentity.TransportOperatorID,
				Identity:            model.TO_SUPER,
			})
			if len(tos) == 0 {
				return nil, errors.New("no super found")
			}
			if tos[0].UserID != toIdentity.UserID {
				return nil, errors.New("to super id and this id not match")
			}
		} else {
			return nil, errors.New("identity can only be driver")
		}
	}

	if err := toIdentity.Create(user); err != nil {
		return nil, err
	}
	return &TransportOperatorIdentityAddOutput{toIdentity}, nil
}

type (
	// TransportOperatorAssignIdentityInput TO组织Super角色添加其它super或者admin参数
	TransportOperatorAssignIdentityInput struct {
		OperatorID          primitive.ObjectID
		TransportOperatorID primitive.ObjectID
		UserID              primitive.ObjectID
		Identity            model.TOIdentity
		Contact             string `valid:"stringlength(1|50)"`
	}
	// TransportOperatorAssignIdentityOutput TO组织Super角色添加其它super或者admin返回参数
	TransportOperatorAssignIdentityOutput struct {
		*model.TransportOperatorIdentity
	}
)

// TransportOperatorAssignIdentity TO组织Super角色添加其它super或者admin
func TransportOperatorAssignIdentity(in *TransportOperatorAssignIdentityInput) (*TransportOperatorAssignIdentityOutput, error) {
	// 参数验证
	if _, err := valid.ValidateStruct(in); err != nil {
		return nil, err
	}

	// 检查操作者是否具有改TO组织的super角色
	if !model.IsTransportOperatorIdentityExists(model.TransportOperatorIdentityExists{
		UserID:              in.OperatorID,
		TransportOperatorID: in.TransportOperatorID,
		Identity:            []model.TOIdentity{model.TO_SUPER},
	}) {
		return nil, errors.New("operator has no authorization")
	}

	// 获取TO组织信息
	to, err := model.TransportOperatorFind(model.TransportOperatorFindOpt{ID: in.TransportOperatorID})
	if err != nil {
		return nil, err
	}
	if !to.IsCompany || !to.IsVerified {
		return nil, errors.New("can only add identity to verified company")
	}

	// 获取被添加者的用户信息
	user, err := model.FindUser(model.FindUserOpt{UserID: in.UserID})
	if err != nil {
		return nil, err
	}

	// 创建TO组织角色信息
	toIdentity := &model.TransportOperatorIdentity{
		ID:                  primitive.NewObjectID(),
		UserID:              in.UserID,
		TransportOperatorID: in.TransportOperatorID,
		Identity:            in.Identity,
		Contact:             &in.Contact,
		CreatedAt:           time.Now(),
	}
	if err := toIdentity.Create(user); err != nil {
		return nil, err
	}
	return &TransportOperatorAssignIdentityOutput{toIdentity}, nil
}

type (
	// TransportOperatorRemoveIdentityInput TO组织Super角色删除其它super或者admin参数
	TransportOperatorRemoveIdentityInput struct {
		OperatorID                  primitive.ObjectID
		TransportOperatorIdentityID primitive.ObjectID
	}
)

// TransportOperatorRemoveIdentity TO组织Super角色删除其它super或者admin
func TransportOperatorRemoveIdentity(in *TransportOperatorRemoveIdentityInput) error {

	// 获取TO组织角色信息
	toIdentity, err := model.TransportOperatorIdentityFind(in.TransportOperatorIdentityID)
	if err != nil {
		return err
	}
	// 检查操作者是否具有改TO组织的super角色, 并且不能删除自己
	if !model.IsTransportOperatorIdentityExists(model.TransportOperatorIdentityExists{
		UserID:              in.OperatorID,
		TransportOperatorID: toIdentity.TransportOperatorID,
		Identity:            []model.TOIdentity{model.TO_SUPER},
	}) {
		return errors.New("operator has no authorization")
	}
	if toIdentity.UserID == in.OperatorID {
		return errors.New("cannot delete yourself")
	}

	// 获取被删除者的用户信息
	user, err := model.FindUser(model.FindUserOpt{UserID: toIdentity.UserID})
	if err != nil {
		return err
	}
	if err := toIdentity.Delete(user); err != nil {
		return err
	}

	return nil
}

type (
	// TransportOperatorVerifyInput 管理人员审批TO组织参数
	TransportOperatorVerifyInput struct {
		TransportOperatorID primitive.ObjectID
	}
)

// TransportOperatorVerify 管理人员审批TO组织
func TransportOperatorVerify(in *TransportOperatorVerifyInput) error {
	// 查询TO组织
	to, err := model.TransportOperatorFind(model.TransportOperatorFindOpt{
		ID: in.TransportOperatorID,
	})
	if err != nil {
		return err
	}
	// 更新TO组织
	to.IsVerified = true
	return to.Update()
}

// UserOperatorDriver User用户拥有对Driver用户操作的权限; 即User是改TO组织的ADMIN以上权限，Driver属于此TO组织
func UserOperatorDriver(userID, DriverID, transportOperatorID primitive.ObjectID) bool {
	return model.IsTransportOperatorIdentityExists(model.TransportOperatorIdentityExists{
		UserID:              userID,
		TransportOperatorID: transportOperatorID,
		Identity:            []model.TOIdentity{model.TO_ADMIN, model.TO_SUPER},
	}) && model.IsTransportOperatorIdentityExists(model.TransportOperatorIdentityExists{
		UserID:              DriverID,
		TransportOperatorID: transportOperatorID,
		Identity:            []model.TOIdentity{model.TO_DRIVER},
	})
}
