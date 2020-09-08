package model

import (
	"context"
	"errors"

	"googlemaps.github.io/maps"
)

// Address 位置信息
type Address string

// GetCoorsFromAddr 通过位置获取坐标信息
func (addr Address) GetCoorsFromAddr() (coors Coors, err error) {
	req := &maps.GeocodingRequest{
		Address: string(addr),
	}
	resp, err := mapClient.Geocode(context.TODO(), req)
	if err != nil {
		return
	}
	if len(resp) == 0 {
		err = errors.New("can not find any match address")
		return
	}
	coors = Coors{Lat: resp[0].Geometry.Location.Lat, Lng: resp[0].Geometry.Location.Lng}
	return
}
