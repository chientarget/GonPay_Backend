// internal/domain/notification.go
package domain

import (
	"time"
)

type NotificationType string

const (
	NotificationTypeTransaction NotificationType = "TRANSACTION"
	NotificationTypeSecurity    NotificationType = "SECURITY"
	NotificationTypeLimit       NotificationType = "LIMIT"
	NotificationTypePromotion   NotificationType = "PROMOTION"
	NotificationTypeAccount     NotificationType = "ACCOUNT"
	NotificationTypeSystem      NotificationType = "SYSTEM"
)

type Notification struct {
	ID               int64            `json:"id"`
	UserID           int64            `json:"user_id"`
	Title            string           `json:"title"`
	Content          string           `json:"content"`
	NotificationType NotificationType `json:"notification_type"`
	IsRead           bool             `json:"is_read"`
	CreatedAt        time.Time        `json:"created_at"`
}

type NotificationRepository interface {
	Create(notification *Notification) error
	GetByID(id int64) (*Notification, error)
	GetByUserID(userID int64, limit, offset int) ([]*Notification, error)
	GetUnreadCount(userID int64) (int, error)
	MarkAsRead(id int64, userID int64) error
	MarkAllAsRead(userID int64) error
}
