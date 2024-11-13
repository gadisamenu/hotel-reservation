package main

import (
	"context"
	"log"
	"time"

	"github.com/gadisamenu/hotel-reservation/db"
	"github.com/gadisamenu/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client       *mongo.Client
	hotelStore   db.HotelStore
	roomStore    db.RoomStore
	userStore    db.UserStore
	bookingStore db.BookingStore
	ctx          = context.Background()
)

func seedUser(fname, lname, email, password string, isAdmin bool) *types.User {

	param := &types.CreateUserParam{
		FirstName: fname,
		LastName:  lname,
		Email:     email,
		Password:  password,
	}

	user, err := param.MapToUser()
	if err != nil {
		log.Fatal(err)
	}
	user.IsAdmin = isAdmin

	inserted, err := userStore.InsertUser(ctx, user)
	if err != nil {
		log.Fatal(err)
	}
	return inserted
}
func seedHotel(name string, location string, rating int) *types.Hotel {

	hotel := &types.Hotel{
		Name:     name,
		Location: location,
		Rooms:    []primitive.ObjectID{},
		Rating:   rating,
	}

	insertedHotel, err := hotelStore.Insert(ctx, hotel)

	if err != nil {
		log.Fatal(err)
	}

	return insertedHotel
}

func seedRoom(hotelId primitive.ObjectID, size string, price float64, numBed int) *types.Room {
	room := &types.Room{
		HotelId: hotelId,
		Size:    size,
		NumBed:  numBed,
		Price:   price,
	}
	inserted, err := roomStore.Insert(context.TODO(), room)

	if err != nil {
		log.Fatal(err)
	}

	return inserted
}

func seedBooking(roomId, userId primitive.ObjectID, from, to time.Time, numPerson int) *types.Booking {
	booking := &types.Booking{
		RoomId:    roomId,
		UserId:    userId,
		FromDate:  from,
		ToDate:    to,
		NumPerson: numPerson,
	}

	inserted, err := bookingStore.InsertBooking(context.TODO(), booking)

	if err != nil {
		log.Fatal(err)
	}

	return inserted
}

func main() {
	// seed hotels
	hotel1 := seedHotel("5start", "Addis", 3)
	seedHotel("Nothing", "New York", 2)
	seedHotel("Nothing", "Sky", 4)
	//seed rooms
	room1 := seedRoom(hotel1.Id, "small", 100, 1)
	seedRoom(hotel1.Id, "medium", 200, 2)
	seedRoom(hotel1.Id, "large", 300, 2)

	//seed users
	user1 := seedUser("GD", "AM", "james@foo.com", "password", false)
	seedUser("Admin", "Admin", "admin@admin.com", "password", true)

	//seed booking
	seedBooking(room1.Id, user1.Id, time.Now(), time.Now().AddDate(0, 0, 2), 2)

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
	userStore = db.NewMongoUserStore(client)
	bookingStore = db.NewMongoBookingStore(client)
}
