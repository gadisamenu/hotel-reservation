package api

import (
	"github.com/gadisamenu/hotel-reservation/db"
	"github.com/gofiber/fiber/v2"
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

type HotelQueryParams struct {
	db.Pagination
	Rating int
}

func (h *HotelHandler) HandleGetHotels(c *fiber.Ctx) error {

	var params HotelQueryParams
	if err := c.QueryParser(&params); err != nil {
		return ErrBadRequest()
	}

	filter := db.MapStr{
		"rating": params.Rating,
	}
	hotels, err := h.store.Hotel.GetAll(c.Context(), filter, &params.Pagination)

	if err != nil {
		return ErrNotFound("hotels")
	}
	resp := db.ResourceResp{
		Data:    hotels,
		Page:    int(params.Page),
		Results: len(hotels),
	}
	return c.JSON(resp)
}

func (h *HotelHandler) HandleGetRooms(c *fiber.Ctx) error {
	id := c.Params("id")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return ErrInvalidId()
	}

	filter := db.MapStr{"hotelId": oid}
	rooms, err := h.store.Room.GetList(c.Context(), filter)
	if err != nil {
		return ErrNotFound("rooms")
	}
	return c.JSON(rooms)
}

func (h *HotelHandler) HandleGetHotelById(c *fiber.Ctx) error {
	id := c.Params("id")

	hotel, err := h.store.Hotel.GetById(c.Context(), id)
	if err != nil {
		return ErrNotFound("hotel")
	}

	return c.JSON(hotel)
}
