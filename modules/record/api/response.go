package api

import (
	"github.com/chadhao/logit/modules/record/model"
)

// ResponseRecord 返回记录结构
type ResponseRecord struct {
	model.Record `json:",inline"`
	Notes        []model.INote `json:"notes"`
}
