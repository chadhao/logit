package api

import (
	"github.com/chadhao/logit/modules/record"
	"github.com/labstack/echo/v4"
)

// AddRecord 添加一条新的记录
func AddRecord(c echo.Context) error {
	rar := new(record.RequestAddRecord)
	if err := c.Bind(rar); err != nil {
		return err
	}
	// r, err := rar.ConstructToRecord(userID)
	// if err != nil {
	// 	return err
	// }
	// if err = r.Add(); err != nil {
	// 	return
	// }
	return nil
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
func AddNote(c echo.Context) (err error) {
	// recordID, err := primitive.ObjectIDFromHex(c.Param("id"))
	// if err != nil {
	// 	return err
	// }
	// r, err := record.GetRecord(recordID)
	// if err != nil {
	// 	return err
	// }
	// if r.UserID != userID {
	// 	return
	// }
	var note record.INote
	noteType := record.NoteType(c.FormValue("noteType"))
	switch noteType {
	case record.SYSTEMNOTE:
		request := new(record.RequestAddSystemNote)
		if err = c.Bind(request); err != nil {
			return err
		}
		note, err = request.ConstructToSystemNote()
		if err != nil {
			return err
		}
	case record.OTHERWORKNOTE:
		request := new(record.RequestAddOtherWorkNote)
		if err = c.Bind(request); err != nil {
			return err
		}
		note, err = request.ConstructToOtherWorkNote()
		if err != nil {
			return err
		}
		// case record.MODIFICATIONNOTE:
		// 	request := new(record.RequestAddModificationNote)
		// 	if err = c.Bind(request); err != nil {
		// 		return err
		// 	}
		// 	note, err = request.ConstructToModificationNote(userID)
		// 	if err != nil {
		// 		return err
		// 	}
	}

	if err = note.Add(); err != nil {
		return err
	}

	return nil
}

// OfflineSyncRecords 离线返回在线状态后记录同步
func OfflineSyncRecords(c echo.Context) {}
