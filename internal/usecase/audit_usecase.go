// internal/usecase/audit_usecase.go

package usecase

import (
	"GonPay_Backend/internal/domain"
	"net"
	"time"
)

type AuditUseCase struct {
	auditRepo domain.AuditRepository
}

func NewAuditUseCase(auditRepo domain.AuditRepository) *AuditUseCase {
	return &AuditUseCase{
		auditRepo: auditRepo,
	}
}

func (u *AuditUseCase) LogAction(
	userID int64,
	action domain.AuditAction,
	entityType string,
	entityID int64,
	oldValue, newValue []byte,
	ip net.IP,
	userAgent string,
) error {
	// Convert IP address to string
	ipString := ""
	if ip != nil {
		ipString = ip.String()
	}

	log := &domain.AuditLog{
		UserID:     userID,
		Action:     action,
		EntityType: entityType,
		EntityID:   entityID,
		OldValue:   oldValue,
		NewValue:   newValue,
		IPAddress:  net.ParseIP(ipString),
		UserAgent:  userAgent,
	}

	return u.auditRepo.Create(log)
}

func (u *AuditUseCase) GetUserAuditLogs(userID int64, page, limit int) ([]*domain.AuditLog, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit
	return u.auditRepo.GetByUserID(userID, limit, offset)
}

func (u *AuditUseCase) GetActionLogs(action domain.AuditAction, page, limit int) ([]*domain.AuditLog, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit
	return u.auditRepo.GetByAction(action, limit, offset)
}

func (u *AuditUseCase) GetDateRangeLogs(startDate, endDate time.Time, page, limit int) ([]*domain.AuditLog, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit
	return u.auditRepo.GetByDateRange(startDate, endDate, limit, offset)
}

func (u *AuditUseCase) GetEntityLogs(entityType string, entityID int64) ([]*domain.AuditLog, error) {
	return u.auditRepo.GetByEntityTypeAndID(entityType, entityID)
}
