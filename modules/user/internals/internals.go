package internals

import (
	"github.com/chadhao/logit/modules/user/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// HasAccessTo 判断userID用户是拥有TO中的ADMIN以上权限，并且driverID用户也从属于这个TO
func HasAccessTo(userID, driverID, transportOperatorID primitive.ObjectID) bool {
	return model.HasAccessTo(userID, driverID, transportOperatorID)
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
