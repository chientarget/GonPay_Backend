// internal/delivery/http/notification_handler.go
package http

import (
	"GonPay_Backend/internal/domain"
	"GonPay_Backend/internal/usecase"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type NotificationHandler struct {
	notificationUseCase *usecase.NotificationUseCase
}

func NewNotificationHandler(notificationUseCase *usecase.NotificationUseCase) *NotificationHandler {
	return &NotificationHandler{
		notificationUseCase: notificationUseCase,
	}
}

func (h *NotificationHandler) GetNotifications(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int64)

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 || limit > 100 {
		limit = 10
	}

	notifications, err := h.notificationUseCase.GetUserNotifications(userID, page, limit)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, notifications)
}

func (h *NotificationHandler) GetUnreadCount(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int64)

	count, err := h.notificationUseCase.GetUnreadCount(userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]int{"unread_count": count})
}

func (h *NotificationHandler) MarkAsRead(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid notification ID")
		return
	}

	userID := r.Context().Value("user_id").(int64)

	if err := h.notificationUseCase.MarkAsRead(id, userID); err != nil {
		switch err {
		case domain.ErrInvalidOperation:
			respondWithError(w, http.StatusForbidden, "Cannot modify this notification")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Notification marked as read"})
}

func (h *NotificationHandler) MarkAllAsRead(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int64)

	if err := h.notificationUseCase.MarkAllAsRead(userID); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "All notifications marked as read"})
}
