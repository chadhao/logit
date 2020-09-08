package model

import (
	"errors"

	locModel "github.com/chadhao/logit/modules/location/model"
)

// Location 位置信息
type Location struct {
	Address locModel.Address `bson:"address,omitempty" json:"address,omitempty" valid:"-"`
	Coors   locModel.Coors   `bson:"coors,omitempty" json:"coors,omitempty" valid:"-"`
}

// Equal 判断两个位置信息是否一致
func (l *Location) Equal(o *Location) bool {
	return l.Address == o.Address
}

// FillFull 若其中一项不完整，则用另外一项查找并补完
func (l *Location) FillFull() (err error) {
	if l.Address != "" && !l.Coors.EmptyCoors() {
		return
	}
	if l.Address == "" && l.Coors.EmptyCoors() {
		return errors.New("at least one field requried")
	}
	if l.Address == "" {
		l.Address, err = l.Coors.GetAddrFromCoors()
	} else {
		l.Coors, err = l.Address.GetCoorsFromAddr()
	}
	return
}
