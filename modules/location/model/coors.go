package model

import (
	"context"
	"errors"

	"googlemaps.github.io/maps"
)

// Coors represents a location on the Earth.
type Coors struct {
	Lat float64 `bson:"lat" json:"lat"`
	Lng float64 `bson:"lng" json:"lng"`
}

// GetAddrFromCoors 通过坐标信息获取位置
func (coors Coors) GetAddrFromCoors() (addr Address, err error) {
	req := &maps.GeocodingRequest{
		LatLng: &maps.LatLng{Lat: coors.Lat, Lng: coors.Lng},
	}
	resp, err := mapClient.Geocode(context.TODO(), req)
	if err != nil {
		return
	}
	if len(resp) == 0 {
		err = errors.New("can not find any match address")
		return
	}
	addr = Address(resp[0].FormattedAddress)
	return
}
