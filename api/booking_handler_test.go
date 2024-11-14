package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gadisamenu/hotel-reservation/api/middleware"
	"github.com/gadisamenu/hotel-reservation/db/fixtures"
	"github.com/gadisamenu/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
)

func TestUserGetBooking(t *testing.T) {
	db := setup(t)
	defer db.teardown(t)
	nonAuthUser := fixtures.AddUser(db.Store, "nonjames", "nonfoo", false)
	user := fixtures.AddUser(db.Store, "james", "foo", false)
	hotel := fixtures.AddHotel(db.Store, "Miracle", "New York", 5, nil)
	room := fixtures.AddRoom(db.Store, hotel.Id, "small", 2, 100)

	from := time.Now()
	till := from.AddDate(0, 0, 2)
	booking := fixtures.AddBooking(db.Store, user.Id, room.Id, 2, from, till)
	_ = booking

	app := fiber.New()
	bookingHandler := NewBookingHandler(db.Store)
	auth := app.Group("/", middleware.JWTAuthentication(db.User))

	auth.Get("/:id", bookingHandler.GetBookingById)

	req := httptest.NewRequest("GET", fmt.Sprintf("/%s", booking.Id.Hex()), nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Api-Token", createTokenFromUser(user))

	resp, err := app.Test(req)

	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 but got %d", resp.StatusCode)
	}

	var respBooking *types.Booking

	if err = json.NewDecoder(resp.Body).Decode(&respBooking); err != nil {
		t.Fatal(err)
	}

	if respBooking.Id != booking.Id {
		t.Fatalf("expected %s got %s ", booking.Id, respBooking.Id)
	}

	if respBooking.UserId != booking.UserId {
		t.Fatalf("expected of user  %s got %s ", booking.Id, respBooking.Id)
	}

	// test get bookings using unauthorized  user
	req = httptest.NewRequest("GET", "/", nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Api-Token", createTokenFromUser(nonAuthUser))

	resp, err = app.Test(req)

	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode == http.StatusOK {
		t.Fatalf("expected non 200 status code but got %d", resp.StatusCode)
	}

}

func TestAdminGetBookings(t *testing.T) {
	db := setup(t)
	defer db.teardown(t)
	adminUser := fixtures.AddUser(db.Store, "admin", "admin", true)
	user := fixtures.AddUser(db.Store, "james", "foo", false)
	hotel := fixtures.AddHotel(db.Store, "Miracle", "New York", 5, nil)
	room := fixtures.AddRoom(db.Store, hotel.Id, "small", 2, 100)

	from := time.Now()
	till := from.AddDate(0, 0, 2)
	booking := fixtures.AddBooking(db.Store, user.Id, room.Id, 2, from, till)
	_ = booking

	app := fiber.New()
	bookingHandler := NewBookingHandler(db.Store)
	admin := app.Group("/", middleware.JWTAuthentication(db.User), middleware.IsAdmin)
	admin.Get("/", bookingHandler.GetBookings)

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Api-Token", createTokenFromUser(adminUser))

	resp, err := app.Test(req)

	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 but got %d", resp.StatusCode)
	}

	var bookings []*types.Booking

	if err = json.NewDecoder(resp.Body).Decode(&bookings); err != nil {
		t.Fatal(err)
	}

	if len(bookings) < 1 {
		t.Fatalf("expected 1 booking but got %d ", len(bookings))
	}

	have := bookings[0]

	if have.Id != booking.Id {
		t.Fatalf("expected %s got %s ", booking.Id, have.Id)
	}

	if have.UserId != booking.UserId {
		t.Fatalf("expected of user  %s got %s ", booking.Id, have.Id)
	}

	// test get bookings using non admin user
	req = httptest.NewRequest("GET", "/", nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Api-Token", createTokenFromUser(user))

	resp, err = app.Test(req)

	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode == http.StatusOK {
		t.Fatalf("expected non 200 status code but got %d", resp.StatusCode)
	}
}
