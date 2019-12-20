package api

import (
	"github.com/chadhao/logit/modules/record"
	"github.com/labstack/echo/v4"
)

// AddRecord 添加一条新的记录
func AddRecord(c echo.Context) error {
	r := new(record.Record)
	if err := c.Bind(r); err != nil {
		return err
	}
	return nil
}

// DeleteLastRecord 删除上一条记录
func DeleteLastRecord(c echo.Context) {}

// GetRecords 获取记录
func GetRecords(c echo.Context) {}

// AddNote 为记录添加笔记
func AddNote(c echo.Context) {}

// OfflineSyncRecords 离线返回在线状态后记录同步
func OfflineSyncRecords(c echo.Context) {}
