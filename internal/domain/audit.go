// internal/domain/audit.go
package domain

import (
	"net"
	"time"
)

type AuditAction string

const (
	AuditActionLogin            AuditAction = "LOGIN"
	AuditActionLogout           AuditAction = "LOGOUT"
	AuditActionTransfer         AuditAction = "TRANSFER"
	AuditActionUpdateProfile    AuditAction = "UPDATE_PROFILE"
	AuditActionChangePassword   AuditAction = "CHANGE_PASSWORD"
	AuditActionEnable2FA        AuditAction = "ENABLE_2FA"
	AuditActionDisable2FA       AuditAction = "DISABLE_2FA"
	AuditActionFailedLogin      AuditAction = "FAILED_LOGIN"
	AuditActionAddPaymentMethod AuditAction = "ADD_PAYMENT_METHOD"
	AuditActionUpdateLimits     AuditAction = "UPDATE_LIMITS"
)

type AuditLog struct {
	ID         int64       `json:"id"`
	UserID     int64       `json:"user_id"`
	Action     AuditAction `json:"action"`
	EntityType string      `json:"entity_type"`
	EntityID   int64       `json:"entity_id"`
	OldValue   []byte      `json:"old_value,omitempty"`
	NewValue   []byte      `json:"new_value,omitempty"`
	IPAddress  net.IP      `json:"ip_address"`
	UserAgent  string      `json:"user_agent"`
	CreatedAt  time.Time   `json:"created_at"`
}

type AuditRepository interface {
	Create(log *AuditLog) error
	GetByUserID(userID int64, limit, offset int) ([]*AuditLog, error)
	GetByAction(action AuditAction, limit, offset int) ([]*AuditLog, error)
	GetByDateRange(startDate, endDate time.Time, limit, offset int) ([]*AuditLog, error)
	GetByEntityTypeAndID(entityType string, entityID int64) ([]*AuditLog, error)
}
