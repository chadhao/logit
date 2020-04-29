package internals

import (
	"github.com/chadhao/logit/modules/user/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// HasAccessTo 判断A用户是否有权限查询B用户
func HasAccessTo(adminIDPlus, driverID, transportOperatorID primitive.ObjectID) bool {
	return model.HasAccessTo(adminIDPlus, driverID, transportOperatorID)
}

// GetVehicleMapByIDs 通过vehicleIDs获取vehicleID对应的其具体信息的vehicleMap
func GetVehicleMapByIDs(ids []primitive.ObjectID) (map[primitive.ObjectID]model.Vehicle, error) {
	m := make(map[primitive.ObjectID]model.Vehicle)
	vehicles, err := model.FindVehiclesByIDs(ids)
	if err != nil {
		return nil, err
	}
	for _, v := range vehicles {
		m[v.ID] = v
	}
	return m, nil
}
