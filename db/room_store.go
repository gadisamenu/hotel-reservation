package db

import (
	"context"

	"github.com/gadisamenu/hotel-reservation/config"
	"github.com/gadisamenu/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const roomColl = "rooms"

type RoomStore interface {
	Dropper
	Insert(context.Context, *types.Room) (*types.Room, error)
	GetList(context.Context, MapStr) ([]*types.Room, error)
}

type MongoRoomStore struct {
	client *mongo.Client
	coll   *mongo.Collection
	HotelStore
}

func NewMongoRoomStore(client *mongo.Client, hotelStore HotelStore) *MongoRoomStore {
	return &MongoRoomStore{
		client:     client,
		coll:       client.Database(config.DB_NAME).Collection(roomColl),
		HotelStore: hotelStore,
	}
}

func (s *MongoRoomStore) Drop(ctx context.Context) error {
	return s.coll.Drop(ctx)
}

func (s *MongoRoomStore) Insert(ctx context.Context, room *types.Room) (*types.Room, error) {
	roomR, err := s.coll.InsertOne(ctx, room)
	if err != nil {
		return nil, err
	}

	room.Id = roomR.InsertedID.(primitive.ObjectID)
	filter := MapStr{"_id": room.HotelId}
	update := MapStr{"$push": bson.M{"rooms": room.Id}}

	err = s.HotelStore.Update(ctx, filter, update)

	if err != nil {
		return nil, err
	}

	return room, nil
}

func (s *MongoRoomStore) GetList(ctx context.Context, filter MapStr) ([]*types.Room, error) {
	resp, err := s.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var rooms []*types.Room

	if err = resp.All(ctx, &rooms); err != nil {
		return nil, err
	}

	return rooms, nil
}
