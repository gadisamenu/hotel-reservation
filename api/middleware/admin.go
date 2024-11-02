package middleware

import (
	"fmt"

	"github.com/gadisamenu/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
)

func IsAdmin(c *fiber.Ctx) error {

	user, ok := c.Context().UserValue("user").(*types.User)

	if !ok || !user.IsAdmin {
		return fmt.Errorf("unauthorized")
	}

	return c.Next()
}
