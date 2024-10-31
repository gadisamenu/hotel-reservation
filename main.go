package main

import (
	"context"
	"flag"
	"log"

	"github.com/gadisamenu/hotel-reservation/api"
	"github.com/gadisamenu/hotel-reservation/api/middlewares"
	"github.com/gadisamenu/hotel-reservation/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var config = fiber.Config{
	ErrorHandler: func(c *fiber.Ctx, err error) error {
		return c.JSON(map[string]string{"error": err.Error()})

	},
}

func main() {
	listenAddr := flag.String("listenAddr", ":5000", "The listen address for api server")
	flag.Parse()

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DbUri))
	if err != nil {
		log.Fatal(err)
	}

	app := fiber.New(config)
	auth := app.Group("/api")
	appv1 := app.Group("/api/v1", middlewares.JWTAuthentication)

	// handlers initialization
	userStore := db.NewMongoUserStore(client)
	userHanlder := api.NewUserHandler(userStore)
	hotelStore := db.NewMongoHotelStore(client)
	roomStore := db.NewMongoRoomStore(client, hotelStore)
	authHandler := api.NewAuthHandler(userStore)

	store := &db.Store{
		Hotel: hotelStore,
		Room:  roomStore,
		User:  userStore,
	}
	hotelHandler := api.NewHotelHandler(store)

	// auth handlers
	auth.Post("/auth", authHandler.HandleAuthenticate)

	// user handlers
	appv1.Delete("/users/:id", userHanlder.HandleDeleteUser)
	appv1.Put("/users/:id", userHanlder.HandleUpdateUser)
	appv1.Post("/users", userHanlder.HandleCreateUser)
	appv1.Get("/users", userHanlder.HandleGetUsers)
	appv1.Get("/users/:id", userHanlder.HandleGetUser)

	//hotel handlers
	appv1.Get("hotels", hotelHandler.HandleGetHotels)
	appv1.Get("hotels/:id", hotelHandler.HandleGetHotelById)
	appv1.Get("hotels/:id/rooms", hotelHandler.HandleGetRooms)

	app.Listen(*listenAddr)
}
