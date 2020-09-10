package internals

import (
	"github.com/chadhao/logit/modules/user/model"
	"github.com/chadhao/logit/modules/user/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CanUserOperatorDriver User用户拥有对Driver用户操作的权限; 即User是改TO组织的ADMIN以上权限，Driver属于此TO组织
func CanUserOperatorDriver(userID, driverID, transportOperatorID primitive.ObjectID) bool {
	return service.UserOperatorDriver(userID, driverID, transportOperatorID)
}

// GetVehicleMapByIDs 通过vehicleIDs获取vehicleID对应的其具体信息的vehicleMap
func GetVehicleMapByIDs(ids []primitive.ObjectID) (map[primitive.ObjectID]*model.Vehicle, error) {
	m := make(map[primitive.ObjectID]*model.Vehicle)
	vehicles, err := model.FindVehicles(model.FindVehiclesOpt{VehicleIDs: ids})
	if err != nil {
		return nil, err
	}
	for _, v := range vehicles {
		m[v.ID] = v
	}
	return m, nil
}
