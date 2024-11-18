// internal/delivery/middleware/admin_middleware.go
package middleware

import (
	"GonPay_Backend/internal/domain"
	"net/http"
)

func (m *Middleware) AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userRole := r.Context().Value("user_role").(string)
		if userRole != domain.RoleAdmin {
			http.Error(w, "Unauthorized: Admin access required", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
