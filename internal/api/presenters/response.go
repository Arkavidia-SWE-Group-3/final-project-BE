package presenters

import "github.com/gofiber/fiber/v2"

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Error   interface{} `json:"error,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func SuccessResponse(ctx *fiber.Ctx, data interface{}, statusCode int, message string) error {
	resp := Response{
		Success: true,
		Message: message,
		Data:    data,
	}

	return ctx.Status(statusCode).JSON(resp)
}

func ErrorResponse(ctx *fiber.Ctx, statusCode int, message string, err error) error {
	resp := Response{
		Success: false,
		Message: message,
		Error:   err.Error(),
	}

	return ctx.Status(statusCode).JSON(resp)
}
