package api

import (
	"github.com/chadhao/logit/modules/record"
	"github.com/labstack/echo/v4"
)

// AddRecord 添加一条新的记录
func AddRecord(c echo.Context) (err error) {
	rar := new(record.RequestAddRecord)
	if err = c.Bind(rar); err != nil {
		return
	}
	if err = rar.Valid(); err != nil {
		return
	}
	// r := rar.ConstructToRecord(userID)
	// if err = r.Add(); err != nil {
	// 	return
	// }
	return
}

// DeleteLastRecord 删除上一条记录
func DeleteLastRecord(c echo.Context) error {
	// recordID, err := primitive.ObjectIDFromHex(c.Param("id"))
	// if err != nil {
	// 	return err
	// }
	// userID :=
	// r, err := record.GetRecord(recordID)
	// if err != nil {
	// 	return err
	// }
	// if err := r.Delete(userID); err != nil {
	// 	return err
	// }
	return nil
}

// GetRecords 获取记录
func GetRecords(c echo.Context) {}

// AddNote 为记录添加笔记
func AddNote(c echo.Context) {}

// OfflineSyncRecords 离线返回在线状态后记录同步
func OfflineSyncRecords(c echo.Context) {}
