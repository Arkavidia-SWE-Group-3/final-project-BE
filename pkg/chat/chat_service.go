package chat

import (
	"Go-Starter-Template/domain"
	"Go-Starter-Template/entities"
	jwtService "Go-Starter-Template/pkg/jwt"
	"Go-Starter-Template/pkg/notification"
	"context"

	"github.com/google/uuid"
)

type (
	ChatService interface {
		GetChatRooms(ctx context.Context, userID string) ([]domain.ChatRoomsResponse, error)
		GetChatRoom(ctx context.Context, userID string, targetUserID string) (domain.ChatRoomResponse, error)
		SendMessage(ctx context.Context, req domain.CreateMessageRequest, userID string) error
		GetMessages(ctx context.Context, userID string, roomID string) (domain.ChatRoomMessageResponse, error)
	}

	chatService struct {
		chatRepository         ChatRepository
		notificationRepository notification.NotificationRepository
		jwtService             jwtService.JWTService
	}
)

func NewChatService(chatRepository ChatRepository, notificationRepository notification.NotificationRepository, jwtService jwtService.JWTService) ChatService {
	return &chatService{chatRepository: chatRepository, notificationRepository: notificationRepository, jwtService: jwtService}
}

func (s *chatService) GetChatRooms(ctx context.Context, userID string) ([]domain.ChatRoomsResponse, error) {
	parsedUserID, err := uuid.Parse(userID)

	if err != nil {
		return []domain.ChatRoomsResponse{}, domain.ErrParseUUID
	}

	var chatRoomsResponse []domain.ChatRoomsResponse

	chatRooms, err := s.chatRepository.GetChatRooms(ctx, parsedUserID)

	for _, room := range chatRooms {
		if room.FirstUserID == parsedUserID {
			chatRoomsResponse = append(chatRoomsResponse, domain.ChatRoomsResponse{
				ID:             room.SecondUser.ID.String(),
				Name:           room.SecondUser.Name,
				ProfilePicture: room.SecondUser.ProfilePicture,
				Type:           room.SecondUser.Role,
				Slug:           room.SecondUser.Slug,
			})
		} else {
			chatRoomsResponse = append(chatRoomsResponse, domain.ChatRoomsResponse{
				ID:             room.FirstUser.ID.String(),
				Name:           room.FirstUser.Name,
				ProfilePicture: room.FirstUser.ProfilePicture,
				Type:           room.FirstUser.Role,
				Slug:           room.FirstUser.Slug,
			})
		}

	}

	if err != nil {
		return []domain.ChatRoomsResponse{}, domain.ErrFailedGetChatRoom
	}

	return chatRoomsResponse, nil
}

func (s *chatService) GetChatRoom(ctx context.Context, userID string, targetUserID string) (domain.ChatRoomResponse, error) {
	parsedUserID, err := uuid.Parse(userID)

	if err != nil {
		return domain.ChatRoomResponse{}, domain.ErrParseUUID
	}

	parsedTargetUserID, err := uuid.Parse(targetUserID)

	if err != nil {
		return domain.ChatRoomResponse{}, domain.ErrParseUUID
	}

	chatRoom, err := s.chatRepository.GetChatRoom(ctx, parsedUserID, parsedTargetUserID)

	if err != nil {
		var randomID uuid.UUID = uuid.New()

		err = s.chatRepository.CreateChatRoom(ctx, entities.ChatRoom{
			ID:           randomID,
			FirstUserID:  parsedUserID,
			SecondUserID: parsedTargetUserID,
		})

		if err != nil {
			return domain.ChatRoomResponse{}, domain.ErrFailedGetChatRoom
		}

		return domain.ChatRoomResponse{
			ID: randomID.String(),
		}, nil
	}

	return domain.ChatRoomResponse{
		ID: chatRoom.ID.String(),
	}, nil
}

func (s *chatService) SendMessage(ctx context.Context, req domain.CreateMessageRequest, userID string) error {
	parsedUserID, err := uuid.Parse(userID)

	if err != nil {
		return domain.ErrParseUUID
	}

	parsedChatRoomID, err := uuid.Parse(req.RoomID)

	if err != nil {
		return domain.ErrParseUUID
	}

	exist, err := s.chatRepository.CheckUserExistInChatRoom(ctx, parsedChatRoomID, parsedUserID)

	if !exist || err != nil {
		return domain.ErrUserNotExistInChatRoom
	}

	err = s.chatRepository.CreateMessage(ctx, entities.ChatMessage{
		RoomID:  parsedChatRoomID,
		UserID:  parsedUserID,
		Message: req.Message,
	})

	if err != nil {
		return domain.ErrFailedCreateMessage
	}

	chatRoom, err := s.chatRepository.GetChatRoomByRoomID(ctx, parsedChatRoomID)

	if err != nil {
		return err
	}
	var targetUser entities.User

	var senderUser entities.User

	if parsedUserID == chatRoom.FirstUserID {
		targetUser = *chatRoom.SecondUser
		senderUser = *chatRoom.FirstUser
	} else {
		targetUser = *chatRoom.FirstUser
		senderUser = *chatRoom.SecondUser
	}

	exist, err = s.notificationRepository.CheckIfSameTitleAndDateExist(ctx, targetUser.ID, "New Message from "+senderUser.Name)

	if exist || err != nil {
		return nil
	}

	err = s.notificationRepository.CreateNotification(ctx, entities.Notification{
		UserID:           targetUser.ID,
		Message:          req.Message,
		Title:            "New Message from " + senderUser.Name,
		IsRead:           false,
		NotificationType: "Message",
	})

	if err != nil {
		return err
	}

	return nil
}

func (s *chatService) GetMessages(ctx context.Context, userID string, roomID string) (domain.ChatRoomMessageResponse, error) {
	parsedUserID, err := uuid.Parse(userID)

	if err != nil {
		return domain.ChatRoomMessageResponse{}, domain.ErrParseUUID
	}

	parsedRoomID, err := uuid.Parse(roomID)

	if err != nil {
		return domain.ChatRoomMessageResponse{}, domain.ErrParseUUID
	}

	exist, err := s.chatRepository.CheckUserExistInChatRoom(ctx, parsedRoomID, parsedUserID)

	if exist == false {
		return domain.ChatRoomMessageResponse{}, domain.ErrUserNotExistInChatRoom
	}

	chatRoom, err := s.chatRepository.GetChatRoomByRoomID(ctx, parsedRoomID)

	if err != nil {
		return domain.ChatRoomMessageResponse{}, domain.ErrFailedGetChatRoom
	}

	messages, err := s.chatRepository.GetMessages(ctx, parsedRoomID)

	var chatMessageResponse []domain.ChatMessageResponse

	for _, message := range messages {
		chatMessageResponse = append(chatMessageResponse, domain.ChatMessageResponse{
			Message:        message.Message,
			Sender:         message.User.Name,
			ProfilePicture: message.User.ProfilePicture,
		})
	}

	if err != nil {
		return domain.ChatRoomMessageResponse{}, domain.ErrFailedGetMessages
	}

	if chatMessageResponse == nil {
		chatMessageResponse = []domain.ChatMessageResponse{}
	}

	if chatRoom.FirstUserID == parsedUserID {
		return domain.ChatRoomMessageResponse{
			ID:             chatRoom.ID.String(),
			Name:           chatRoom.SecondUser.Name,
			ProfilePicture: chatRoom.SecondUser.ProfilePicture,
			Messages:       chatMessageResponse,
		}, nil
	} else {
		return domain.ChatRoomMessageResponse{
			ID:             chatRoom.ID.String(),
			Name:           chatRoom.FirstUser.Name,
			ProfilePicture: chatRoom.FirstUser.ProfilePicture,
			Messages:       chatMessageResponse,
		}, nil
	}
}
