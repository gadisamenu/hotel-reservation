package db

import "context"

const MongoDbname = "hotel-reservation"
const TestMongoDbname = "hotel-reservation-test"

// const dburi = "mongodb://admin:pass@localhost:27017?authSource=admin&retryWrites=true&w=majority"
const DbUri = "mongodb://admin:pass@localhost:27017/"

type Dropper interface {
	Drop(context.Context) error
}

type Store struct {
	Hotel   HotelStore
	User    UserStore
	Room    RoomStore
	Booking BookingStore
}
