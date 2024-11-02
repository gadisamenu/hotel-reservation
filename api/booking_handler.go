package api

import (
	"fmt"

	"github.com/gadisamenu/hotel-reservation/db"
	"github.com/gadisamenu/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

type BookingHandler struct {
	store *db.Store
}

func NewBookingHandler(store *db.Store) *BookingHandler {
	return &BookingHandler{
		store: store,
	}
}

func (h *BookingHandler) GetBookings(c *fiber.Ctx) error {
	bookings, err := h.store.Booking.GetBookings(c.Context(), bson.M{})
	if err != nil {
		return err
	}

	return c.JSON(bookings)
}

func (h *BookingHandler) GetBookingById(c *fiber.Ctx) error {
	id := c.Params("id")
	user, ok := c.Context().UserValue("user").(*types.User)

	if !ok {
		return fmt.Errorf("unauthorized")
	}
	booking, err := h.store.Booking.GetBookingById(c.Context(), id, user.Id)
	if err != nil {
		return err
	}

	return c.JSON(booking)
}
