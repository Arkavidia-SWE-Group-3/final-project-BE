package domain

import (
	"errors"
)

var (
	MessageFailedGetChatRoom    = "Failed to get chat room"
	MessageFailedCreateChatRoom = "Failed to create chat room"
	MessageFailedCreateMessage  = "Failed to create message"
	MessageFailedGetMessages    = "Failed to get messages"

	MessageSuccessGetChatRoom    = "Successfully get chat room"
	MessageSuccessCreateChatRoom = "Successfully create chat room"
	MessageSuccessCreateMessage  = "Successfully create message"
	MessageSuccessGetMessages    = "Successfully get messages"

	ErrFailedGetChatRoom      = errors.New("failed to get chat room")
	ErrFailedCreateChatRoom   = errors.New("failed to create chat room")
	ErrFailedCreateMessage    = errors.New("failed to create message")
	ErrFailedGetMessages      = errors.New("failed to get messages")
	ErrUserNotExistInChatRoom = errors.New("user not exist in chat room")
)

type (
	ChatRoomsResponse struct {
		ID             string `json:"id"`
		Name           string `json:"name"`
		ProfilePicture string `json:"profile_picture"`
		Type           string `json:"type"`
		Slug           string `json:"slug"`
	}

	ChatRoomResponse struct {
		ID string `json:"id"`
	}

	CreateMessageRequest struct {
		RoomID  string `json:"room_id" validate:"required,uuid4"`
		Message string `json:"message" validate:"required"`
	}

	ChatRoomMessageResponse struct {
		ID             string                `json:"id"`
		Name           string                `json:"name"`
		ProfilePicture string                `json:"profile_picture"`
		Messages       []ChatMessageResponse `json:"messages"`
	}

	ChatMessageResponse struct {
		Message        string `json:"message"`
		Sender         string `json:"sender"`
		ProfilePicture string `json:"profile_picture"`
	}
)
