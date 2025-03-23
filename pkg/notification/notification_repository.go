package notification

import (
	"Go-Starter-Template/entities"
	"context"

	"gorm.io/gorm"
)

type (
	NotificationRepository interface {
		GetNotification(ctx context.Context, userID string) ([]entities.Notification, error)
		ReadNotification(ctx context.Context, notificationID string) error
		GetNotificationByID(ctx context.Context, notificationID string) (entities.Notification, error)
		CreateNotification(ctx context.Context, notification entities.Notification) error
	}

	notificationRepository struct {
		db *gorm.DB
	}
)

func NewNotificationRepository(db *gorm.DB) NotificationRepository {
	return &notificationRepository{db: db}
}

func (r *notificationRepository) ReadNotification(ctx context.Context, notificationID string) error {
	err := r.db.Model(&entities.Notification{}).Where("id = ?", notificationID).Update("is_read", true).Error

	if err != nil {
		return err
	}

	return nil
}

func (r *notificationRepository) GetNotification(ctx context.Context, userID string) ([]entities.Notification, error) {
	var notifications []entities.Notification

	err := r.db.Where("user_id = ?", userID).Find(&notifications).Error

	if err != nil {
		return nil, err
	}

	return notifications, nil
}

func (r *notificationRepository) GetNotificationByID(ctx context.Context, notificationID string) (entities.Notification, error) {
	var notification entities.Notification

	err := r.db.Where("id = ?", notificationID).First(&notification).Error

	if err != nil {
		return entities.Notification{}, err
	}

	return notification, nil
}

func (r *notificationRepository) CreateNotification(ctx context.Context, notification entities.Notification) error {
	err := r.db.WithContext(ctx).Create(&notification).Error

	if err != nil {
		return err
	}

	return nil
}
