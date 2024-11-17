package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gadisamenu/hotel-reservation/db/fixtures"
	"github.com/gofiber/fiber/v2"
)

func TestAuthenticateWithWrongPassword(t *testing.T) {

	testdb := setup(t)
	defer testdb.teardown(t)

	fixtures.AddUser(testdb.Store, "james", "foo", false)

	app := fiber.New(fiber.Config{ErrorHandler: ErrorHandler})
	authHandler := NewAuthHandler(testdb.User)

	params := &AuthParams{
		Email:    "james@foo.com",
		Password: "passworddd",
	}

	app.Post("/", authHandler.HandleAuthenticate)
	b, _ := json.Marshal(params)

	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")

	resp, _ := app.Test(req)

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected response status code 400 but got %d", resp.StatusCode)
	}

	var apiErr Error
	if err := json.NewDecoder(resp.Body).Decode(&apiErr); err != nil {
		t.Fatal(err)
	}

	if apiErr.Code != http.StatusBadRequest {
		t.Fatalf("expected response status code 400 but got %d", resp.StatusCode)
	}
}

func TestAuthenticateSuccess(t *testing.T) {

	testdb := setup(t)
	defer testdb.teardown(t)

	insertedUser := fixtures.AddUser(testdb.Store, "james", "foo", false)

	app := fiber.New()
	authHandler := NewAuthHandler(testdb.User)

	params := &AuthParams{
		Email:    "james@foo.com",
		Password: "james_foo",
	}

	app.Post("/", authHandler.HandleAuthenticate)
	b, _ := json.Marshal(params)

	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")

	resp, _ := app.Test(req)

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected response status code 200 but got %d", resp.StatusCode)
	}

	var authResp AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		t.Fatal(err)
	}

	if len(authResp.Token) == 0 {
		t.Fatalf("expected token but empty")
	}

	insertedUser.EncryptedPassword = ""

	if !reflect.DeepEqual(insertedUser, authResp.User) {
		t.Fatalf("expected user %v but got %v", insertedUser, authResp.User)
	}
}
