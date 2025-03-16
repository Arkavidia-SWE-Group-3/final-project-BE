package entities

import "github.com/google/uuid"

type Companies struct {
	ID       uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Name     string    `json:"name"`
	Slug     string    `json:"slug"`
	About    string    `json:"about"`
	Industry string    `json:"industry"`
	UserID   uuid.UUID `json:"user_id"`

	User *User `gorm:"foreignKey:UserID"`

	Timestamp
}
