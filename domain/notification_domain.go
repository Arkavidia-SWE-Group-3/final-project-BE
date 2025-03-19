package domain

import (
	"errors"
)

var (
	MessageFailedGetNotification   = "Failed to get notification"
	MessageSuccessReadNotification = "Successfully read notification"
	MessageFailedReadNotification  = "Failed to read notification"
	MessageSuccessGetNotification  = "Successfully get notification"

	ErrorNotificationNotFound = errors.New("notification not found")
)

type (
	NotificationResponse struct {
		ID      string `json:"id"`
		Title   string `json:"title"`
		Message string `json:"message"`
		IsRead  bool   `json:"is_read"`
		Type    string `json:"type"`
	}
)
