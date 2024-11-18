package db

import (
	"context"
	"fmt"

	"github.com/gadisamenu/hotel-reservation/config"
	"github.com/gadisamenu/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const bookingColl = "bookings"

type BookingStore interface {
	InsertBooking(context.Context, *types.Booking) (*types.Booking, error)
	GetBookings(context.Context, MapStr) ([]*types.Booking, error)
	GetBookingById(context.Context, string) (*types.Booking, error)
	UpdateById(context.Context, string, MapStr) error
}

type MongoBookingStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoBookingStore(client *mongo.Client) *MongoBookingStore {
	return &MongoBookingStore{
		client: client,
		coll:   client.Database(config.DB_NAME).Collection(bookingColl),
	}
}

func (s *MongoBookingStore) UpdateById(ctx context.Context, id string, data MapStr) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	update := bson.M{
		"$set": data,
	}

	updated, err := s.coll.UpdateByID(ctx, oid, update)
	if err != nil {
		return err
	}
	if updated.ModifiedCount == 0 {
		return fmt.Errorf("not modified")
	}

	return nil
}

func (s *MongoBookingStore) InsertBooking(ctx context.Context, booking *types.Booking) (*types.Booking, error) {

	res, err := s.coll.InsertOne(ctx, booking)
	if err != nil {
		return nil, err
	}

	booking.Id = res.InsertedID.(primitive.ObjectID)

	return booking, nil
}

func (s *MongoBookingStore) GetBookings(ctx context.Context, filter MapStr) ([]*types.Booking, error) {
	res, err := s.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var bookings []*types.Booking

	if err = res.All(ctx, &bookings); err != nil {
		return nil, err
	}
	return bookings, nil
}

func (s *MongoBookingStore) GetBookingById(ctx context.Context, id string) (*types.Booking, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	res := s.coll.FindOne(ctx, bson.M{"_id": oid})

	var booking *types.Booking

	if err = res.Decode(&booking); err != nil {
		return nil, err
	}
	return booking, nil
}
