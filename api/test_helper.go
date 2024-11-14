package api

import (
	"context"
	"testing"

	"github.com/gadisamenu/hotel-reservation/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type testdb struct {
	client *mongo.Client
	*db.Store
}

func (tdb *testdb) teardown(t *testing.T) {
	if err := tdb.client.Database(db.MongoDbname).Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
}

func setup(t *testing.T) *testdb {

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DbUri))
	if err != nil {
		t.Fatal(err)
	}

	hotelStore := db.NewMongoHotelStore(client)

	return &testdb{
		client: client,
		Store: &db.Store{
			Hotel:   hotelStore,
			Room:    db.NewMongoRoomStore(client, hotelStore),
			User:    db.NewMongoUserStore(client),
			Booking: db.NewMongoBookingStore(client),
		},
	}
}
