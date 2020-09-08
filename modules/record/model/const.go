package model

type (
	// Type 记录类型
	Type string
	// NoteType 笔记类型
	NoteType string
	// HrTime 重要时间点
	HrTime float64
)

const (
	// WORK 工作记录类型
	WORK Type = "work"
	// REST 休息记录类型
	REST Type = "rest"
)
const (
	// SYSTEMNOTE 系统笔记类型
	SYSTEMNOTE NoteType = "system"
	// MODIFICATIONNOTE 人为修改笔记类型
	MODIFICATIONNOTE NoteType = "modification"
	// TRIPNOTE 行程笔记类型
	TRIPNOTE NoteType = "trip"
	// OTHERWORKNOTE 其它笔记类型
	OTHERWORKNOTE NoteType = "others"
)

const (
	// HR0D5 0.5小时
	HR0D5 HrTime = 0.5
	// HR5D5 5.5小时
	HR5D5 HrTime = 5.5
	// HR7 7.5小时
	HR7 HrTime = 7
	// HR10 10小时
	HR10 HrTime = 10
	// HR13 13小时
	HR13 HrTime = 13
	// HR24 24小时
	HR24 HrTime = 24
	// HR70 70小时
	HR70 HrTime = 70
)

func (t HrTime) getHrs() float64 {
	return float64(t)
}
