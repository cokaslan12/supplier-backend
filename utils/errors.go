package utils

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func ErrorHandler(c *fiber.Ctx, err error) error {
	if apiError, ok := err.(Error); ok {
		return c.Status(apiError.Code).JSON(apiError)
	}
	internalErr := NewError(http.StatusInternalServerError, false, err.Error())
	return c.Status(internalErr.Code).JSON(internalErr)
}

type Error struct {
	Code    int    `json:"-"`
	Success bool   `json:"success"`
	Err     string `json:"error"`
}

// Error implements the error interface
func (e Error) Error() string {
	return e.Err
}

func NewError(code int, success bool, msg string) Error {
	return Error{
		Success: success,
		Code:    code,
		Err:     msg,
	}
}

func ErrUnAuthorized() Error {
	return Error{
		Code:    http.StatusUnauthorized,
		Success: false,
		Err:     "Unauthorized request",
	}
}

func ErrResourceNotFound(res string) Error {
	return Error{
		Code:    http.StatusNoContent,
		Success: false,
		Err:     res + " resource not found",
	}
}

func ErrBadRequest() Error {
	return Error{
		Code:    http.StatusBadRequest,
		Success: false,
		Err:     "Inavlid json request",
	}
}

func ErrInValidId() Error {
	return Error{
		Code:    http.StatusBadRequest,
		Success: false,
		Err:     "Invalid Id",
	}
}
