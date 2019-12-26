package api

import (
	"github.com/chadhao/logit/modules/record/model"
)

// respRecord 返回记录结构
type respRecord struct {
	model.Record `json:",inline"`
	Notes        []model.INote `json:"notes"`
}
