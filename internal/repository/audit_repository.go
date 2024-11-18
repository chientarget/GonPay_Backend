// internal/repository/audit_repository.go
package repository

import (
	"GonPay_Backend/internal/domain"
	"database/sql"
	_ "database/sql"
	"net"
	"time"
)

type auditRepository struct {
	db *PostgresDB
}

func NewAuditRepository(db *PostgresDB) domain.AuditRepository {
	return &auditRepository{db: db}
}

func (r *auditRepository) Create(log *domain.AuditLog) error {
	query := `
        INSERT INTO audit_logs (user_id, action, entity_type, entity_id, old_value, new_value, ip_address, user_agent)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
        RETURNING log_id, created_at`

	return r.db.DB.QueryRow(
		query,
		log.UserID,
		log.Action,
		log.EntityType,
		log.EntityID,
		log.OldValue,
		log.NewValue,
		log.IPAddress.String(),
		log.UserAgent,
	).Scan(&log.ID, &log.CreatedAt)
}

func (r *auditRepository) GetByUserID(userID int64, limit, offset int) ([]*domain.AuditLog, error) {
	query := `
        SELECT 
            log_id, user_id, action, entity_type, entity_id,
            old_value, new_value, ip_address, user_agent, created_at
        FROM audit_logs 
        WHERE user_id = $1
        ORDER BY created_at DESC
        LIMIT $2 OFFSET $3`

	rows, err := r.db.DB.Query(query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanAuditLogs(rows)
}

func (r *auditRepository) GetByAction(action domain.AuditAction, limit, offset int) ([]*domain.AuditLog, error) {
	query := `
        SELECT 
            log_id, user_id, action, entity_type, entity_id,
            old_value, new_value, ip_address, user_agent, created_at
        FROM audit_logs 
        WHERE action = $1
        ORDER BY created_at DESC
        LIMIT $2 OFFSET $3`

	rows, err := r.db.DB.Query(query, action, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanAuditLogs(rows)
}

func (r *auditRepository) GetByDateRange(startDate, endDate time.Time, limit, offset int) ([]*domain.AuditLog, error) {
	query := `
        SELECT 
            log_id, user_id, action, entity_type, entity_id,
            old_value, new_value, ip_address, user_agent, created_at
        FROM audit_logs 
        WHERE created_at BETWEEN $1 AND $2
        ORDER BY created_at DESC
        LIMIT $3 OFFSET $4`

	rows, err := r.db.DB.Query(query, startDate, endDate, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanAuditLogs(rows)
}

func (r *auditRepository) GetByEntityTypeAndID(entityType string, entityID int64) ([]*domain.AuditLog, error) {
	query := `
        SELECT 
            log_id, user_id, action, entity_type, entity_id, 
            old_value, new_value, ip_address, user_agent, created_at
        FROM audit_logs 
        WHERE entity_type = $1 AND entity_id = $2
        ORDER BY created_at DESC`

	rows, err := r.db.DB.Query(query, entityType, entityID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanAuditLogs(rows)
}

func (r *auditRepository) queryAuditLogs(query string, args ...interface{}) ([]*domain.AuditLog, error) {
	rows, err := r.db.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []*domain.AuditLog
	for rows.Next() {
		log := &domain.AuditLog{}
		var ipStr string
		err := rows.Scan(
			&log.ID,
			&log.UserID,
			&log.Action,
			&log.EntityType,
			&log.EntityID,
			&log.OldValue,
			&log.NewValue,
			&ipStr,
			&log.UserAgent,
			&log.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		log.IPAddress = net.ParseIP(ipStr)
		logs = append(logs, log)
	}

	return logs, nil
}

func (r *auditRepository) scanAuditLogs(rows *sql.Rows) ([]*domain.AuditLog, error) {
	var logs []*domain.AuditLog

	for rows.Next() {
		log := &domain.AuditLog{}
		var ipStr, userAgentStr sql.NullString

		err := rows.Scan(
			&log.ID,
			&log.UserID,
			&log.Action,
			&log.EntityType,
			&log.EntityID,
			&log.OldValue,
			&log.NewValue,
			&ipStr,
			&userAgentStr,
			&log.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		if ipStr.Valid {
			log.IPAddress = net.ParseIP(ipStr.String)
		} else {
			log.IPAddress = nil
		}

		if userAgentStr.Valid {
			log.UserAgent = userAgentStr.String
		} else {
			log.UserAgent = ""
		}

		logs = append(logs, log)
	}

	return logs, nil
}
