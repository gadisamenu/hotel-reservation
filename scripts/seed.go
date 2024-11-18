package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/gadisamenu/hotel-reservation/config"
	"github.com/gadisamenu/hotel-reservation/db"
	"github.com/gadisamenu/hotel-reservation/db/fixtures"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MONGO_DB_URI))
	if err != nil {
		log.Fatal(err)
	}

	err = client.Database(config.DB_NAME).Drop(ctx)

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

	for i := 0; i < 200; i++ {
		name := fmt.Sprintf("hotel name %d", i)
		location := fmt.Sprintf("location %d", i)
		fixtures.AddHotel(store, name, location, rand.Intn(5)+1, nil)
	}

}
