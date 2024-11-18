// internal/delivery/http/audit_handler.go
package http

import (
	"GonPay_Backend/internal/domain"
	"GonPay_Backend/internal/usecase"
	"net/http"
	"strconv"
	"time"
)

type AuditHandler struct {
	auditUseCase *usecase.AuditUseCase
}

func NewAuditHandler(auditUseCase *usecase.AuditUseCase) *AuditHandler {
	return &AuditHandler{
		auditUseCase: auditUseCase,
	}
}

// GetUserAuditLogs returns audit logs for the authenticated user
func (h *AuditHandler) GetUserAuditLogs(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int64)

	page, limit := getPaginationParams(r)

	logs, err := h.auditUseCase.GetUserAuditLogs(userID, page, limit)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, logs)
}

// GetActionLogs returns audit logs filtered by action (Admin only)
func (h *AuditHandler) GetActionLogs(w http.ResponseWriter, r *http.Request) {
	action := domain.AuditAction(r.URL.Query().Get("action"))
	if action == "" {
		respondWithError(w, http.StatusBadRequest, "action parameter is required")
		return
	}

	page, limit := getPaginationParams(r)

	logs, err := h.auditUseCase.GetActionLogs(action, page, limit)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, logs)
}

// GetDateRangeLogs returns audit logs within a date range (Admin only)
func (h *AuditHandler) GetDateRangeLogs(w http.ResponseWriter, r *http.Request) {
	startDate, endDate, err := getDateRangeParams(r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	page, limit := getPaginationParams(r)

	logs, err := h.auditUseCase.GetDateRangeLogs(startDate, endDate, page, limit)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, logs)
}

// GetEntityLogs returns audit logs for a specific entity (Admin only)
func (h *AuditHandler) GetEntityLogs(w http.ResponseWriter, r *http.Request) {
	entityType := r.URL.Query().Get("entity_type")
	if entityType == "" {
		respondWithError(w, http.StatusBadRequest, "entity_type parameter is required")
		return
	}

	entityID, err := strconv.ParseInt(r.URL.Query().Get("entity_id"), 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid entity_id")
		return
	}

	logs, err := h.auditUseCase.GetEntityLogs(entityType, entityID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, logs)
}

// Helper functions

func getPaginationParams(r *http.Request) (page, limit int) {
	page, _ = strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	limit, _ = strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 || limit > 100 {
		limit = 10
	}

	return page, limit
}

func getDateRangeParams(r *http.Request) (startDate, endDate time.Time, err error) {
	startDateStr := r.URL.Query().Get("start_date")
	if startDateStr == "" {
		startDate = time.Now().AddDate(0, -1, 0) // Default to 1 month ago
	} else {
		startDate, err = time.Parse(time.RFC3339, startDateStr)
		if err != nil {
			return time.Time{}, time.Time{}, err
		}
	}

	endDateStr := r.URL.Query().Get("end_date")
	if endDateStr == "" {
		endDate = time.Now()
	} else {
		endDate, err = time.Parse(time.RFC3339, endDateStr)
		if err != nil {
			return time.Time{}, time.Time{}, err
		}
	}

	return startDate, endDate, nil
}
