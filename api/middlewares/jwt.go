package middlewares

import (
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuthentication(c *fiber.Ctx) error {
	token, ok := c.GetReqHeaders()["X-Api-Token"]

	if !ok {
		return fmt.Errorf("unauthorized")
	}
	claims, err := validateToken(token[0])

	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("unauthorized")
	}

	expiration := claims["expires"].(string)
	layout := "2006-01-02T15:04:05.999999999-07:00"
	expTime, err := time.Parse(layout, expiration)

	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("unauthorized")
	}

	if time.Now().After(expTime) {
		return fmt.Errorf("token expired")
	}

	return c.Next()
}

func validateToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {

		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Printf("Unexpected  signing method %v", t.Header["alg"])
			return nil, fmt.Errorf("unauthorized")
		}

		secret := os.Getenv(("JWT_SECRET"))

		return []byte(secret), nil

	})

	if err != nil {
		log.Error(err)
		return nil, fmt.Errorf("unauthorized")

	}

	if !token.Valid {
		return nil, fmt.Errorf("unauthorized")
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return nil, fmt.Errorf("unauthorized")
	}

	return claims, nil

}
