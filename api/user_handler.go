package api

import (
	"github.com/gadisamenu/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
)

func HandleGetUsers(c *fiber.Ctx) error {
	user := types.User{
		FirstName: "GD",
		LastName:  "AM",
	}
	return c.JSON(user)
}

func HandleGetUser(c *fiber.Ctx) error {
	return c.JSON("Gadisa")
}
