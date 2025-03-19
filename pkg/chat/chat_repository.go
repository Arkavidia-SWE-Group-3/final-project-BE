package chat

import (
	"gorm.io/gorm"

	"Go-Starter-Template/entities"
	"context"

	"github.com/google/uuid"
)

type (
	ChatRepository interface {
		GetChatRooms(ctx context.Context, userID uuid.UUID) ([]entities.ChatRoom, error)
		GetChatRoom(ctx context.Context, userID uuid.UUID, targetUserID uuid.UUID) (entities.ChatRoom, error)
		CreateChatRoom(ctx context.Context, chatRoom entities.ChatRoom) error
		CreateMessage(ctx context.Context, message entities.ChatMessage) error
		GetMessages(ctx context.Context, roomID uuid.UUID) ([]entities.ChatMessage, error)
		GetChatRoomByRoomID(ctx context.Context, roomID uuid.UUID) (entities.ChatRoom, error)
		CheckUserExistInChatRoom(ctx context.Context, roomID uuid.UUID, userID uuid.UUID) (bool, error)
	}
	chatRepository struct {
		db *gorm.DB
	}
)

func NewChatRepository(db *gorm.DB) ChatRepository {
	return &chatRepository{db: db}
}

func (r *chatRepository) GetChatRooms(ctx context.Context, userID uuid.UUID) ([]entities.ChatRoom, error) {
	var chatRooms []entities.ChatRoom
	if err := r.db.WithContext(ctx).Preload("FirstUser").Preload("SecondUser").Where("first_user_id = ? OR second_user_id = ?", userID, userID).Find(&chatRooms).Error; err != nil {
		return nil, err
	}
	return chatRooms, nil
}

func (r *chatRepository) GetChatRoom(ctx context.Context, userID uuid.UUID, targetUserID uuid.UUID) (entities.ChatRoom, error) {
	var chatRoom entities.ChatRoom
	if err := r.db.WithContext(ctx).Preload("FirstUser").Preload("SecondUser").Where("(first_user_id = ? AND second_user_id = ?) OR (first_user_id = ? AND second_user_id = ?)", userID, targetUserID, targetUserID, userID).First(&chatRoom).Error; err != nil {
		return entities.ChatRoom{}, err
	}
	return chatRoom, nil
}

func (r *chatRepository) GetChatRoomByRoomID(ctx context.Context, roomID uuid.UUID) (entities.ChatRoom, error) {
	var chatRoom entities.ChatRoom
	if err := r.db.WithContext(ctx).Preload("FirstUser").Preload("SecondUser").Where("id = ?", roomID).First(&chatRoom).Error; err != nil {
		return entities.ChatRoom{}, err
	}
	return chatRoom, nil
}

func (r *chatRepository) CreateChatRoom(ctx context.Context, chatRoom entities.ChatRoom) error {
	if err := r.db.WithContext(ctx).Create(&chatRoom).Error; err != nil {
		return err
	}
	return nil
}

func (r *chatRepository) CreateMessage(ctx context.Context, message entities.ChatMessage) error {
	if err := r.db.WithContext(ctx).Create(&message).Error; err != nil {
		return err
	}

	chatRoom, err := r.GetChatRoomByRoomID(ctx, message.RoomID)

	if err != nil {
		return err
	}
	var targetUser entities.User

	var senderUser entities.User

	if message.UserID == chatRoom.FirstUserID {
		targetUser = *chatRoom.SecondUser
		senderUser = *chatRoom.FirstUser
	} else {
		targetUser = *chatRoom.FirstUser
		senderUser = *chatRoom.SecondUser
	}

	notification := entities.Notification{
		UserID:           targetUser.ID,
		Message:          message.Message,
		Title:            "New Message from " + senderUser.Name,
		IsRead:           false,
		NotificationType: "Message",
	}

	var existingNotification entities.Notification
	if err := r.db.WithContext(ctx).Where("user_id = ? AND title = ? AND DATE(created_at) = DATE(?)", targetUser.ID, notification.Title, gorm.Expr("NOW()")).First(&existingNotification).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			if err := r.db.WithContext(ctx).Create(&notification).Error; err != nil {
				return err
			}
		} else {
			return nil
		}
	}

	return nil
}

func (r *chatRepository) GetMessages(ctx context.Context, roomID uuid.UUID) ([]entities.ChatMessage, error) {
	var messages []entities.ChatMessage
	if err := r.db.WithContext(ctx).Preload("User").Where("room_id = ?", roomID).Find(&messages).Error; err != nil {
		return nil, err
	}
	return messages, nil
}

func (r *chatRepository) CheckUserExistInChatRoom(ctx context.Context, roomID uuid.UUID, userID uuid.UUID) (bool, error) {
	var chatRoom entities.ChatRoom
	if err := r.db.WithContext(ctx).Where("id = ?", roomID).First(&chatRoom).Error; err != nil {
		return false, err
	}
	if chatRoom.FirstUserID == userID || chatRoom.SecondUserID == userID {
		return true, nil
	}
	return false, nil
}
