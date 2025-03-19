package notification

import (
	"Go-Starter-Template/domain"
	"Go-Starter-Template/entities"
	jwtService "Go-Starter-Template/pkg/jwt"
	"context"

	"github.com/google/uuid"
)

type (
	NotificationService interface {
		GetNotification(ctx context.Context, userID string) ([]entities.Notification, error)
		ReadNotification(ctx context.Context, notificationID string, userID string) error
	}

	notificationService struct {
		notificationRepository NotificationRepository
		jwtService             jwtService.JWTService
	}
)

func NewNotificationService(notificationRepository NotificationRepository, jwtService jwtService.JWTService) NotificationService {
	return &notificationService{notificationRepository: notificationRepository, jwtService: jwtService}
}

func (s *notificationService) GetNotification(ctx context.Context, userID string) ([]entities.Notification, error) {
	res, err := s.notificationRepository.GetNotification(ctx, userID)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *notificationService) ReadNotification(ctx context.Context, notificationID string, userID string) error {

	parsedUserID, err := uuid.Parse(userID)

	if err != nil {
		return domain.ErrParseUUID
	}

	notification, err := s.notificationRepository.GetNotificationByID(ctx, notificationID)

	if err != nil || notification.UserID != parsedUserID {
		return domain.ErrorNotificationNotFound
	}

	err = s.notificationRepository.ReadNotification(ctx, notificationID)

	if err != nil {
		return err
	}

	return nil
}
