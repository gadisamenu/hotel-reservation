package api

import (
	"errors"

	"github.com/gadisamenu/hotel-reservation/db"
	"github.com/gadisamenu/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserHandler struct {
	userstore db.UserStore
}

func NewUserHandler(userstore db.UserStore) *UserHandler {
	return &UserHandler{
		userstore: userstore,
	}
}

func (h *UserHandler) HandleDeleteUser(c *fiber.Ctx) error {
	var id = c.Params("id")

	if err := h.userstore.DeleteUser(c.Context(), id); err != nil {
		return err
	}
	return c.JSON(map[string]string{"deleted": id})
}

func (h *UserHandler) HandleUpdateUser(c *fiber.Ctx) error {
	var (
		id     = c.Params("id")
		params types.UpdateUserParam
	)

	if err := c.BodyParser(&params); err != nil {
		return ErrBadRequest()
	}

	filter := db.MapStr{"_id": id}

	if err := h.userstore.UpdateUser(c.Context(), filter, params); err != nil {
		return err
	}

	return c.JSON(map[string]string{"updated": id})
}

func (h *UserHandler) HandleCreateUser(c *fiber.Ctx) error {
	var params types.CreateUserParam
	if err := c.BodyParser(&params); err != nil {
		return ErrBadRequest()
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
		if errors.Is(err, mongo.ErrNoDocuments) {
			return c.JSON(map[string]string{"error": "user not found"})
		}
		return err
	}

	return c.JSON(user)
}
