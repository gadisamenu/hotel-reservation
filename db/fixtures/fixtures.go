package fixtures

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gadisamenu/hotel-reservation/db"
	"github.com/gadisamenu/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddUser(store *db.Store, fname, lname string, isAdmin bool) *types.User {
	param := &types.CreateUserParam{
		FirstName: fname,
		LastName:  lname,
		Email:     fmt.Sprintf("%s@%s.com", fname, lname),
		Password:  fmt.Sprintf("%s_%s", fname, lname),
	}

	user, err := param.MapToUser()
	if err != nil {
		log.Fatal(err)
	}
	user.IsAdmin = isAdmin

	inserted, err := store.User.InsertUser(context.Background(), user)
	if err != nil {
		log.Fatal(err)
	}
	return inserted
}

func AddHotel(store *db.Store, name, loc string, rating int, roomsId []primitive.ObjectID) *types.Hotel {
	roomsOId := roomsId

	if roomsOId == nil {
		roomsOId = []primitive.ObjectID{}
	}

	hotel := &types.Hotel{
		Name:     name,
		Location: loc,
		Rooms:    roomsOId,
		Rating:   rating,
	}

	insertedHotel, err := store.Hotel.Insert(context.Background(), hotel)

	if err != nil {
		log.Fatal(err)
	}

	return insertedHotel
}

func AddRoom(store *db.Store, hotelId primitive.ObjectID, size string, numBed int, price float64) *types.Room {
	room := &types.Room{
		HotelId: hotelId,
		Size:    size,
		NumBed:  numBed,
		Price:   price,
	}
	inserted, err := store.Room.Insert(context.Background(), room)

	if err != nil {
		log.Fatal(err)
	}

	return inserted
}

func AddBooking(store *db.Store, userId, roomId primitive.ObjectID, numPerson int, from, till time.Time) *types.Booking {
	booking := &types.Booking{
		RoomId:    roomId,
		UserId:    userId,
		FromDate:  from,
		ToDate:    till,
		NumPerson: numPerson,
	}

	inserted, err := store.Booking.InsertBooking(context.Background(), booking)

	if err != nil {
		log.Fatal(err)
	}

	return inserted
}
