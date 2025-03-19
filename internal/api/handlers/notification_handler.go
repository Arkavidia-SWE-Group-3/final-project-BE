package handlers

import (
	"Go-Starter-Template/internal/api/presenters"
	"Go-Starter-Template/pkg/notification"
	"fmt"

	"github.com/go-playground/validator/v10"

	"github.com/gofiber/fiber/v2"

	"Go-Starter-Template/domain"
)

type (
	NotificationHandler interface {
		GetNotifications(c *fiber.Ctx) error
		ReadNotification(c *fiber.Ctx) error
	}
	notificationHandler struct {
		NotificationService notification.NotificationService
		Validator           *validator.Validate
	}
)

func NewNotificationHandler(notificationService notification.NotificationService, validator *validator.Validate) NotificationHandler {
	return &notificationHandler{
		NotificationService: notificationService,
		Validator:           validator,
	}
}

func (h *notificationHandler) GetNotifications(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	fmt.Printf("userID: %s\n", userID)

	res, err := h.NotificationService.GetNotification(c.Context(), userID)

	if err != nil {
		return presenters.ErrorResponse(c, fiber.StatusBadRequest, domain.MessageFailedGetNotification, err)
	}

	return presenters.SuccessResponse(c, res, fiber.StatusOK, domain.MessageSuccessGetNotification)
}

func (h *notificationHandler) ReadNotification(c *fiber.Ctx) error {
	notificationID := c.Params("id")
	userID := c.Locals("user_id").(string)

	err := h.NotificationService.ReadNotification(c.Context(), notificationID, userID)

	if err != nil {
		return presenters.ErrorResponse(c, fiber.StatusBadRequest, domain.MessageFailedReadNotification, err)
	}

	return presenters.SuccessResponse(c, nil, fiber.StatusOK, domain.MessageSuccessReadNotification)
}
