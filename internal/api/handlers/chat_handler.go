package handlers

import (
	"Go-Starter-Template/pkg/chat"

	"github.com/go-playground/validator/v10"

	"github.com/gofiber/fiber/v2"

	"Go-Starter-Template/domain"
	"Go-Starter-Template/internal/api/presenters"
)

type (
	ChatHandler interface {
		GetChatRooms(c *fiber.Ctx) error
		GetChatRoom(c *fiber.Ctx) error
		SendMessage(c *fiber.Ctx) error
		GetMessages(c *fiber.Ctx) error
	}

	chatHandler struct {
		ChatService chat.ChatService
		Validator   *validator.Validate
	}
)

func NewChatHandler(chatService chat.ChatService, validator *validator.Validate) ChatHandler {
	return &chatHandler{
		ChatService: chatService,
		Validator:   validator,
	}
}

func (h *chatHandler) GetChatRooms(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	res, err := h.ChatService.GetChatRooms(c.Context(), userID)

	if err != nil {
		return presenters.ErrorResponse(c, fiber.StatusBadRequest, domain.MessageFailedGetChatRoom, err)
	}

	return presenters.SuccessResponse(c, res, fiber.StatusOK, domain.MessageSuccessGetChatRoom)

}

func (h *chatHandler) GetChatRoom(c *fiber.Ctx) error {
	targetUserID := c.Params("id")
	userID := c.Locals("user_id").(string)

	res, err := h.ChatService.GetChatRoom(c.Context(), userID, targetUserID)

	if err != nil {
		return presenters.ErrorResponse(c, fiber.StatusBadRequest, domain.MessageFailedGetChatRoom, err)
	}

	return presenters.SuccessResponse(c, res, fiber.StatusOK, domain.MessageSuccessGetChatRoom)
}

func (h *chatHandler) SendMessage(c *fiber.Ctx) error {
	var req domain.CreateMessageRequest

	if err := c.BodyParser(&req); err != nil {
		return presenters.ErrorResponse(c, fiber.StatusBadRequest, domain.MessageFailedCreateMessage, err)
	}

	if err := h.Validator.Struct(req); err != nil {
		return presenters.ErrorResponse(c, fiber.StatusBadRequest, domain.MessageFailedCreateMessage, err)
	}

	userID := c.Locals("user_id").(string)

	err := h.ChatService.SendMessage(c.Context(), req, userID)

	if err != nil {
		return presenters.ErrorResponse(c, fiber.StatusBadRequest, domain.MessageFailedCreateMessage, err)
	}

	return presenters.SuccessResponse(c, nil, fiber.StatusCreated, domain.MessageSuccessCreateMessage)
}

func (h *chatHandler) GetMessages(c *fiber.Ctx) error {
	roomID := c.Params("id")
	userID := c.Locals("user_id").(string)

	res, err := h.ChatService.GetMessages(c.Context(), userID, roomID)

	if err != nil {
		return presenters.ErrorResponse(c, fiber.StatusBadRequest, domain.MessageFailedGetMessages, err)
	}

	return presenters.SuccessResponse(c, res, fiber.StatusOK, domain.MessageSuccessGetMessages)
}
