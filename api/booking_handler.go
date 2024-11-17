package api

import (
	"errors"

	"github.com/gadisamenu/hotel-reservation/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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
		return ErrNotFound("bookings")
	}

	return c.JSON(bookings)
}

func (h *BookingHandler) CancelBooking(c *fiber.Ctx) error {
	booking, err := h.store.Booking.GetBookingById(c.Context(), c.Params("id"))

	if err != nil {
		return ErrNotFound("booking")
	}

	user, err := getAuthUser(c)
	if err != nil {
		return err
	}

	if booking.UserId != user.Id {
		return ErrUnAuthorized()
	}

	err = h.store.Booking.UpdateById(c.Context(), booking.Id.String(), bson.M{"canceled": true})
	if err != nil {
		return err
	}

	return c.JSON(&genericResp{
		Type: "success",
		Msg:  "updated",
	})
}

func (h *BookingHandler) GetBookingById(c *fiber.Ctx) error {
	id := c.Params("id")

	user, err := getAuthUser(c)
	if err != nil {
		return err
	}

	booking, err := h.store.Booking.GetBookingById(c.Context(), id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ErrNotFound("booking")
		}
		return err
	}

	if booking.UserId != user.Id {
		return ErrUnAuthorized()
	}

	return c.JSON(booking)
}
