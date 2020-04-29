package api

import (
	"github.com/chadhao/logit/modules/record/model"
	userModel "github.com/chadhao/logit/modules/user/model"
)

// respRecord 返回记录结构
type respRecord struct {
	model.Record `json:",inline"`
	Notes        model.DifNotes    `json:"notes,omitempty"`
	Vehicle      userModel.Vehicle `json:"vehicle"`
}
