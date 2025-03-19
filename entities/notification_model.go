package entities

import "github.com/google/uuid"

type Notification struct {
	ID               uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID           uuid.UUID `json:"user_id"`
	Title            string    `json:"title"`
	Message          string    `json:"message"`
	IsRead           bool      `json:"is_read"`
	NotificationType string    `json:"notification_type"`
	Timestamp
}
