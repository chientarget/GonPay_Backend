// internal/usecase/notification_usecase.go
package usecase

import (
	"GonPay_Backend/internal/domain"
)

type NotificationUseCase struct {
	notificationRepo domain.NotificationRepository
}

func NewNotificationUseCase(notificationRepo domain.NotificationRepository) *NotificationUseCase {
	return &NotificationUseCase{
		notificationRepo: notificationRepo,
	}
}

func (u *NotificationUseCase) CreateNotification(userID int64, title string, content string, notificationType domain.NotificationType) (*domain.Notification, error) {
	notification := &domain.Notification{
		UserID:           userID,
		Title:            title,
		Content:          content,
		NotificationType: notificationType,
		IsRead:           false,
	}

	if err := u.notificationRepo.Create(notification); err != nil {
		return nil, err
	}

	return notification, nil
}

func (u *NotificationUseCase) GetUserNotifications(userID int64, page, limit int) ([]*domain.Notification, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit
	return u.notificationRepo.GetByUserID(userID, limit, offset)
}

func (u *NotificationUseCase) GetNotification(id int64, userID int64) (*domain.Notification, error) {
	notification, err := u.notificationRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if notification.UserID != userID {
		return nil, domain.ErrInvalidOperation
	}

	return notification, nil
}

func (u *NotificationUseCase) GetUnreadCount(userID int64) (int, error) {
	return u.notificationRepo.GetUnreadCount(userID)
}

func (u *NotificationUseCase) MarkAsRead(id int64, userID int64) error {
	return u.notificationRepo.MarkAsRead(id, userID)
}

func (u *NotificationUseCase) MarkAllAsRead(userID int64) error {
	return u.notificationRepo.MarkAllAsRead(userID)
}
