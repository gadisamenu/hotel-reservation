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

// const dburi = "mongodb://admin:pass@localhost:27017?authSource=admin&retryWrites=true&w=majority"
const dburi = "mongodb://admin:pass@localhost:27017/"

var config = fiber.Config{
	ErrorHandler: func(c *fiber.Ctx, err error) error {
		return c.JSON(map[string]string{"error": err.Error()})

	},
}

func main() {
	listenAddr := flag.String("listenAddr", ":5000", "The listen address for api server")
	flag.Parse()

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dburi))
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println(client)

	app := fiber.New(config)
	appv1 := app.Group("/api/v1")

	// handlers initialization
	userHanlder := api.NewUserHandler(db.NewMongoUserStore(client))

	appv1.Post("/users", userHanlder.HandleInsertUser)
	appv1.Get("/users", userHanlder.HandleGetUsers)
	appv1.Get("/users/:id", userHanlder.HandleGetUser)

	app.Listen(*listenAddr)
}
