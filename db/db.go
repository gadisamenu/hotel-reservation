package db

import (
	"context"
)

type Dropper interface {
	Drop(context.Context) error
}

type Pagination struct {
	Page  int64
	Limit int64
}

type ResourceResp struct {
	Results int `json:"results"`
	Data    any `json:"data"`
	Page    int `json:"page"`
}

type Store struct {
	Hotel   HotelStore
	User    UserStore
	Room    RoomStore
	Booking BookingStore
}
