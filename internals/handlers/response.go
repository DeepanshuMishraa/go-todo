package handlers

import "github.com/gofiber/fiber/v2"

type ErrorResponse struct {
	Error string `json:"error"`
}

type SuccessResponse struct {
	Data interface{} `json:"data"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

func RespondWithError(c *fiber.Ctx, code int, message string) error {
	return c.Status(code).JSON(ErrorResponse{Error: message})
}

func RespondWithJSON(c *fiber.Ctx, code int, payload interface{}) error {
	return c.Status(code).JSON(payload)
}
