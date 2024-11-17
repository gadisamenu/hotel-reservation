package main

import (
	"context"
	"flag"
	"log"

	"github.com/gadisamenu/hotel-reservation/api"
	"github.com/gadisamenu/hotel-reservation/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var config = fiber.Config{
	ErrorHandler: api.ErrorHandler,
}

func main() {
	listenAddr := flag.String("listenAddr", ":5000", "The listen address for api server")
	flag.Parse()

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DbUri))
	if err != nil {
		log.Fatal(err)
	}

	userStore := db.NewMongoUserStore(client)
	hotelStore := db.NewMongoHotelStore(client)
	roomStore := db.NewMongoRoomStore(client, hotelStore)
	bookingStore := db.NewMongoBookingStore(client)

	store := &db.Store{
		Hotel:   hotelStore,
		Room:    roomStore,
		User:    userStore,
		Booking: bookingStore,
	}

	hotelHandler := api.NewHotelHandler(store)
	userHanlder := api.NewUserHandler(userStore)
	authHandler := api.NewAuthHandler(userStore)
	roomHandler := api.NewRoomHandler(store)
	bookingHandler := api.NewBookingHandler(store)

	app := fiber.New(config)
	auth := app.Group("/api")
	appv1 := app.Group("/api/v1", api.JWTAuthentication(userStore))
	admin := appv1.Group("/admin", api.IsAdmin)

	// handlers initialization

	// auth handlers
	auth.Post("/auth", authHandler.HandleAuthenticate)

	// user handlers
	appv1.Delete("/users/:id", userHanlder.HandleDeleteUser)
	appv1.Put("/users/:id", userHanlder.HandleUpdateUser)
	appv1.Post("/users", userHanlder.HandleCreateUser)
	appv1.Get("/users", userHanlder.HandleGetUsers)
	appv1.Get("/users/:id", userHanlder.HandleGetUser)

	//hotel handlers
	appv1.Get("/hotels", hotelHandler.HandleGetHotels)
	appv1.Get("/hotels/:id", hotelHandler.HandleGetHotelById)
	appv1.Get("/hotels/:id/rooms", hotelHandler.HandleGetRooms)

	//rooms handlers
	appv1.Get("/rooms", roomHandler.HandleGetRooms)
	appv1.Post("/rooms/:id/book", roomHandler.HandleBookRooms)
	appv1.Get("/rooms/:id/cancel", bookingHandler.CancelBooking)

	//booking handlers
	appv1.Get("/bookings/:id", bookingHandler.GetBookingById)

	//admin handlers
	admin.Get("/bookings", bookingHandler.GetBookings)

	app.Listen(*listenAddr)
}
