package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gadisamenu/hotel-reservation/db"
	"github.com/gadisamenu/hotel-reservation/db/fixtures"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {

	ctx := context.Background()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(db.DbUri))
	if err != nil {
		log.Fatal(err)
	}

	err = client.Database(db.MongoDbname).Drop(ctx)

	if err != nil {
		log.Fatal(err)
	}

	hotelStore := db.NewMongoHotelStore(client)

	store := &db.Store{
		Room:    db.NewMongoRoomStore(client, hotelStore),
		User:    db.NewMongoUserStore(client),
		Booking: db.NewMongoBookingStore(client),
		Hotel:   hotelStore,
	}

	user := fixtures.AddUser(store, "james", "foo", false)
	admin := fixtures.AddUser(store, "admin", "admin", false)
	fmt.Println("admin -> ", admin.Id)

	hotel := fixtures.AddHotel(store, "Miracle", "New York", 5, nil)
	room := fixtures.AddRoom(store, hotel.Id, "small", 2, 100)
	from := time.Now()
	till := from.AddDate(0, 0, 2)
	booking := fixtures.AddBooking(store, user.Id, room.Id, 2, from, till)

	fmt.Println("booking -> ", booking.Id)
}
