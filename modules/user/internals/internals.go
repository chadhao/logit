package internals

import (
	"github.com/chadhao/logit/modules/user/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// HasAccessTo 判断A用户是否有权限查询B用户
func HasAccessTo(adminIDPlus, driverID, transportOperatorID primitive.ObjectID) bool {
	return model.HasAccessTo(adminIDPlus, driverID, transportOperatorID)
}
