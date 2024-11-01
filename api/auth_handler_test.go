package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gadisamenu/hotel-reservation/db"
	"github.com/gadisamenu/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
)

func seedTestUser(t *testing.T, userStore db.UserStore) *types.User {
	param := &types.CreateUserParam{
		FirstName: "fTest",
		LastName:  "lTest",
		Email:     "test@gmail.com",
		Password:  "password",
	}

	user, err := param.MapToUser()
	if err != nil {
		t.Fatal(err)
	}

	_, err = userStore.InsertUser(context.TODO(), user)
	if err != nil {
		t.Fatal(err)
	}

	return user
}

func TestAuthenticateWithWrongPassword(t *testing.T) {

	testdb := setup()
	defer testdb.teardown(t)

	seedTestUser(t, testdb.UserStore)

	app := fiber.New()
	authHandler := NewAuthHandler(testdb.UserStore)

	params := &AuthParams{
		Email:    "test@gmail.com",
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

	var genResp genericResp
	if err := json.NewDecoder(resp.Body).Decode(&genResp); err != nil {
		t.Fatal(err)
	}

	if genResp.Type != "error" {
		t.Fatalf("expected generic response type error but got %s", genResp.Type)
	}
	if genResp.Msg != "invalid credentials" {
		t.Fatalf("expected generic response msg invalid credentials but got %s", genResp.Msg)
	}
}

func TestAuthenticateSuccess(t *testing.T) {

	testdb := setup()
	defer testdb.teardown(t)

	insertedUser := seedTestUser(t, testdb.UserStore)

	app := fiber.New()
	authHandler := NewAuthHandler(testdb.UserStore)

	params := &AuthParams{
		Email:    "test@gmail.com",
		Password: "password",
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
