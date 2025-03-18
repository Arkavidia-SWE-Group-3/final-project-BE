package entities

import "github.com/google/uuid"

type ChatRoom struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	FirstUserID  uuid.UUID `json:"first_user_id"`
	SecondUserID uuid.UUID `json:"second_user_id"`

	FirstUser  *User `gorm:"foreignKey:FirstUserID"`
	SecondUser *User `gorm:"foreignKey:SecondUserID"`

	Timestamp
}
