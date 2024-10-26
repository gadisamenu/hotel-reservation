package api

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http/httptest"
	"testing"

	"github.com/gadisamenu/hotel-reservation/db"
	"github.com/gadisamenu/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type testdb struct {
	db.UserStore
}

func (tdb *testdb) teardown(t *testing.T) {
	if err := tdb.UserStore.Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
}

func setup() *testdb {

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DbUri))
	if err != nil {
		log.Fatal(err)
	}

	return &testdb{
		UserStore: db.NewMongoUserStore(client),
	}
}

func TestPostUser(t *testing.T) {

	testdb := setup()
	defer testdb.teardown(t)

	params := types.CreateUserParam{
		Email:     "test@test.com",
		Password:  "password",
		FirstName: "firstName",
		LastName:  "lastName",
	}

	app := fiber.New()
	userHandler := NewUserHandler(testdb.UserStore)
	app.Post("/", userHandler.HandleCreateUser)
	b, _ := json.Marshal(params)

	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")

	resp, _ := app.Test(req)

	var user types.User
	json.NewDecoder(resp.Body).Decode(&user)

	if params.FirstName != user.FirstName {
		t.Errorf("expected firstName to be %s but found %s", params.FirstName, user.FirstName)
	}
	if params.LastName != user.LastName {
		t.Errorf("expected lastName to be %s but found %s", params.LastName, user.LastName)
	}
	if params.Email != user.Email {
		t.Errorf("expected email to be %s but found %s", params.Email, user.Email)
	}
	if len(user.Id) == 0 {
		t.Errorf("expected Id to be set")
	}
	if len(user.EncryptedPassword) > 0 {
		t.Errorf("expected EncryptedPassword not to be included in json response")
	}

}
