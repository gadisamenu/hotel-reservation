package api

import (
	"errors"
	"fmt"
	"time"

	"github.com/gadisamenu/hotel-reservation/config"
	"github.com/gadisamenu/hotel-reservation/db"
	"github.com/gadisamenu/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	userstore db.UserStore
}

func NewAuthHandler(userStore db.UserStore) *AuthHandler {
	return &AuthHandler{
		userstore: userStore,
	}
}

type AuthParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	User  *types.User `json:"user"`
	Token string      `json:"token"`
}

type genericResp struct {
	Type string `json:"type"`
	Msg  string `json:"msg"`
}

func (h *AuthHandler) HandleAuthenticate(c *fiber.Ctx) error {
	var params AuthParams

	if err := c.BodyParser(&params); err != nil {
		return err
	}

	user, err := h.userstore.GetUserByEmail(c.Context(), params.Email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ErrInvalidCredentials()
		}
		return err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.EncryptedPassword), []byte(params.Password)); err != nil {
		return ErrInvalidCredentials()
	}

	resp := AuthResponse{
		User:  user,
		Token: createTokenFromUser(user),
	}

	return c.JSON(resp)
}

func createTokenFromUser(user *types.User) string {
	now := time.Now()
	expires := now.Add(time.Hour * 4)
	claims := jwt.MapClaims{
		"id":      user.Id,
		"email":   user.Email,
		"expires": expires,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := config.JWT_SECRET

	tokenStr, err := token.SignedString([]byte(secret))
	if err != nil {
		fmt.Println("failed to sign token with secret", err)
	}
	return tokenStr

}
