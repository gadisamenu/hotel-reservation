package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gadisamenu/hotel-reservation/db"
	"github.com/gadisamenu/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BookingParams struct {
	FromDate  time.Time `json:"fromDate"`
	ToDate    time.Time `json:"toDate"`
	NumPerson int       `json:"numPerson"`
}

func (bp *BookingParams) validate() error {
	now := time.Now()

	if now.After(bp.FromDate) || now.After(bp.ToDate) {
		return fmt.Errorf("you can't book for the past")
	}

	if bp.FromDate.After(bp.ToDate) {
		return fmt.Errorf("toDate can't be before fromDate")
	}

	return nil
}

type RoomHandler struct {
	store *db.Store
}

func NewRoomHandler(store *db.Store) *RoomHandler {
	return &RoomHandler{
		store: store,
	}
}

func (h *RoomHandler) HandleGetRooms(c *fiber.Ctx) error {

	bookings, err := h.store.Booking.GetBookings(c.Context(), db.MapStr{})
	if err != nil {
		return ErrNotFound("rooms")
	}

	return c.JSON(bookings)
}

func (h *RoomHandler) HandleBookRooms(c *fiber.Ctx) error {
	roomId, err := primitive.ObjectIDFromHex(c.Params("id"))

	if err != nil {
		return ErrInvalidId()
	}

	var params BookingParams
	if err := c.BodyParser(&params); err != nil {
		return ErrBadRequest()
	}

	if err := params.validate(); err != nil {
		return err
	}

	user, err := getAuthUser(c)
	if err != nil {
		return err
	}

	filter := db.MapStr{
		"fromDate": bson.M{
			"$gte": params.FromDate,
		},
		"toDate": bson.M{
			"$lte": params.ToDate,
		},
		"roomId": roomId,
	}

	bookings, err := h.store.Booking.GetBookings(c.Context(), filter)
	if err != nil {
		return err
	}

	if len(bookings) > 0 {
		return NewError(http.StatusBadRequest, fmt.Sprintf("the room %s is already booked", roomId.Hex()))
	}

	booking := types.Booking{
		UserId:    user.Id,
		RoomId:    roomId,
		FromDate:  params.FromDate,
		ToDate:    params.ToDate,
		NumPerson: params.NumPerson,
	}

	response, err := h.store.Booking.InsertBooking(c.Context(), &booking)

	if err != nil {
		return err
	}

	return c.JSON(response)
}
