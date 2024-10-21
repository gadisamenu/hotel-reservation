package main

import (
	"flag"

	"github.com/gadisamenu/hotel-reservation/api"
	"github.com/gofiber/fiber/v2"
)

func main() {
	listenAddr := flag.String("listenAddr", ":5000", "The listen address for api server")
	flag.Parse()

	app := fiber.New()
	appv1 := app.Group("/api/v1")

	appv1.Get("/users", api.HandleGetUsers)
	appv1.Get("/users/:id", api.HandleGetUser)
	app.Listen(*listenAddr)
}
