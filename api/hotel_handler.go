package api

import (
	"github.com/gadisamenu/hotel-reservation/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HotelHandler struct {
	store *db.Store
}

func NewHotelHandler(store *db.Store) *HotelHandler {
	return &HotelHandler{
		store: store,
	}
}

func (h *HotelHandler) HandleGetHotels(c *fiber.Ctx) error {

	filter := bson.M{}
	hotels, err := h.store.Hotel.GetAll(c.Context(), filter)

	if err != nil {
		return ErrNotFound("hotels")
	}
	return c.JSON(hotels)
}

func (h *HotelHandler) HandleGetRooms(c *fiber.Ctx) error {
	id := c.Params("id")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return ErrInvalidId()
	}

	filter := bson.M{"hotelId": oid}
	rooms, err := h.store.Room.GetList(c.Context(), filter)
	if err != nil {
		return ErrNotFound("rooms")
	}
	return c.JSON(rooms)
}

func (h *HotelHandler) HandleGetHotelById(c *fiber.Ctx) error {
	id := c.Params("id")

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return ErrInvalidId()
	}

	hotel, err := h.store.Hotel.GetById(c.Context(), oid)
	if err != nil {
		return ErrNotFound("hotel")
	}

	return c.JSON(hotel)
}
