package internal

import "go.mongodb.org/mongo-driver/bson/primitive"

// AddDrivingLoc 添加一条行驶信息
func AddDrivingLoc(userID primitive.ObjectID, req *ReqAddDrivingLoc) error {

	drivingLoc, err := req.constructToDrivingLoc(userID)
	if err != nil {
		return err
	}
	if err = drivingLoc.Save(); err != nil {
		return err
	}
	return nil
}
