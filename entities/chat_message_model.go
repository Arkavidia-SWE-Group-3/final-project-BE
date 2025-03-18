package entities

import "github.com/google/uuid"

type ChatMessage struct {
	ID      uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	RoomID  uuid.UUID `json:"room_id"`
	UserID  uuid.UUID `json:"user_id"`
	Message string    `json:"message"`

	User *User     `gorm:"foreignKey:UserID"`
	Room *ChatRoom `gorm:"foreignKey:RoomID"`

	Timestamp
}
