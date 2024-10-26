package main

import (
	"context"
	"log"

	"github.com/gadisamenu/hotel-reservation/db"
	"github.com/gadisamenu/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client     *mongo.Client
	hotelStore db.HotelStore
	roomStore  db.RoomStore
	ctx        = context.Background()
)

func seedHotel(name string, location string, rating int) {

	hotel := types.Hotel{
		Name:     name,
		Location: location,
		Rooms:    []primitive.ObjectID{},
		Rating:   rating,
	}

	insertedHotel, _ := hotelStore.Insert(ctx, &hotel)

	rooms := []types.Room{
		{
			Size:    "small",
			Price:   100,
			Number:  1,
			HotelId: insertedHotel.Id,
		},
		{
			Size:    "normal",
			Price:   130,
			Number:  2,
			HotelId: insertedHotel.Id,
		},
		{
			Size:    "kingsize",
			Price:   200,
			Number:  3,
			HotelId: insertedHotel.Id,
		},
	}

	for _, room := range rooms {
		_, err := roomStore.Insert(ctx, &room)

		if err != nil {
			log.Fatal(err)
		}
	}
}

func main() {

	seedHotel("5start", "Addis", 3)
	seedHotel("Nothing", "New York", 2)
	seedHotel("Nothing", "Sky", 4)

}

func init() {
	var err error
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DbUri))
	if err != nil {
		log.Fatal(err)
	}

	err = client.Database(db.MongoDbname).Drop(ctx)

	if err != nil {
		log.Fatal(err)
	}

	hotelStore = db.NewMongoHotelStore(client)
	roomStore = db.NewMongoRoomStore(client, hotelStore)
}
