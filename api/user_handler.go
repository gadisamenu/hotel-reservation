package api

import (
	"github.com/gadisamenu/hotel-reservation/db"
	"github.com/gadisamenu/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	userstore db.UserStore
}

func NewUserHandler(userstore db.UserStore) *UserHandler {
	return &UserHandler{
		userstore: userstore,
	}
}

func (h *UserHandler) HandleInsertUser(c *fiber.Ctx) error {
	var params types.CreateUserParam
	if err := c.BodyParser(&params); err != nil {
		return err
	}
	if err := params.Validate(); len(err) > 0 {
		return c.JSON(err)
	}

	user, err := params.MapToUser()
	if err != nil {
		return err
	}

	user, err = h.userstore.InsertUser(c.Context(), user)
	if err != nil {
		return err
	}
	return c.JSON(user)
}

func (h *UserHandler) HandleGetUsers(c *fiber.Ctx) error {
	users, err := h.userstore.GetUsers(c.Context())

	if err != nil {
		return err
	}
	return c.JSON(users)
}

func (h *UserHandler) HandleGetUser(c *fiber.Ctx) error {
	id := c.Params("id")

	user, err := h.userstore.GetUserByID(c.Context(), id)
	if err != nil {
		return err
	}

	return c.JSON(user)
}
