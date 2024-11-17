package api

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type Error struct {
	Code int    `json:"code"`
	Err  string `json:"error"`
}

// Error implements error.
func (e Error) Error() string {
	return e.Err
}

func ErrorHandler(c *fiber.Ctx, err error) error {
	if apiErr, ok := err.(Error); ok {
		return c.Status(apiErr.Code).JSON(apiErr)
	}

	apiErr := NewError(http.StatusInternalServerError, err.Error())
	return c.Status(apiErr.Code).JSON(apiErr)
}

func NewError(code int, err string) Error {
	return Error{
		Code: code,
		Err:  err,
	}
}

func ErrInvalidId() Error {
	return Error{
		Code: http.StatusBadRequest,
		Err:  "invalid id given",
	}
}

func ErrUnAuthorized() Error {
	return Error{
		Code: http.StatusUnauthorized,
		Err:  "unauthorized",
	}
}

func ErrInvalidCredentials() Error {
	return Error{
		Code: http.StatusBadRequest,
		Err:  "invalid credentials",
	}
}

func ErrBadRequest() Error {
	return Error{
		Code: http.StatusBadRequest,
		Err:  "invalid JSON request",
	}

}

func ErrNotFound(res string) Error {
	return Error{
		Code: http.StatusNotFound,
		Err:  res + " not found",
	}
}
