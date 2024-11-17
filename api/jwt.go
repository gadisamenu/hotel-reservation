package api

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gadisamenu/hotel-reservation/db"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuthentication(userStore db.UserStore) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token, ok := c.GetReqHeaders()["X-Api-Token"]

		if !ok {
			fmt.Println("unauthorized")
			return ErrUnAuthorized()
		}
		claims, err := validateToken(token[0])

		if err != nil {
			fmt.Println(err)
			return ErrUnAuthorized()
		}

		expiration := claims["expires"].(string)
		layout := "2006-01-02T15:04:05.999999999-07:00"
		expTime, err := time.Parse(layout, expiration)

		if err != nil {
			fmt.Println(err)
			return ErrUnAuthorized()
		}

		if time.Now().After(expTime) {
			return NewError(http.StatusUnauthorized, "token expired")
		}

		userId := claims["id"].(string)

		user, err := userStore.GetUserByID(c.Context(), userId)
		if err != nil {
			fmt.Println(err)
			return ErrUnAuthorized()
		}

		c.Context().SetUserValue("user", user)
		return c.Next()
	}
}

func validateToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {

		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Printf("Unexpected  signing method %v", t.Header["alg"])
			return nil, ErrUnAuthorized()
		}

		secret := os.Getenv(("JWT_SECRET"))

		return []byte(secret), nil

	})

	if err != nil {
		log.Error(err)
		return nil, ErrUnAuthorized()

	}

	if !token.Valid {
		return nil, ErrUnAuthorized()
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return nil, ErrUnAuthorized()
	}

	return claims, nil

}
