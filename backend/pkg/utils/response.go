package utils

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/labstack/echo/v4"
)

type ErrorResponseBody struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Errors  interface{} `json:"errors,omitempty"`
}

func ErrorResponse(ctx echo.Context, code int, message string, err error) error {
	response := ErrorResponseBody{
		Success: false,
		Message: message,
	}

	// Handle validation errors specifically
	if e, ok := err.(validation.Errors); ok {
		response.Errors = e
	} else if err != nil {
		response.Errors = err.Error()
	}

	return ctx.JSON(code, response)
}

type SuccessResponseBody struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func SuccessResponse(ctx echo.Context, code int, message string, data interface{}) error {
	return ctx.JSON(code, SuccessResponseBody{
		Success: true,
		Message: message,
		Data:    data,
	})
}
