package main

import (
	"GonPay_Backend/internal/config"
	httpDelivery "GonPay_Backend/internal/delivery/http"
	"GonPay_Backend/internal/delivery/middleware"
	"GonPay_Backend/internal/repository"
	"GonPay_Backend/internal/usecase"
	"GonPay_Backend/pkg/logger"
	"GonPay_Backend/pkg/validator"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Cannot load config:", err)
	}

	// Initialize logger
	logger := logger.NewLogger()

	// Initialize validator
	validator := validator.NewValidator()
	dbConnStr := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.DBName,
		cfg.Database.SSLMode,
	)

	db, err := repository.NewPostgresDB(dbConnStr)
	if err != nil {
		logger.Error("Cannot connect to database", "error", err)
		os.Exit(1)
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	walletRepo := repository.NewWalletRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)

	// Initialize use cases
	userUseCase := usecase.NewUserUseCase(userRepo, validator, cfg.JWT.Secret, cfg.JWT.TTL)
	walletUseCase := usecase.NewWalletUseCase(walletRepo, transactionRepo)
	transactionUseCase := usecase.NewTransactionUseCase(transactionRepo)

	// Initialize payment method repository and usecase
	paymentMethodRepo := repository.NewPaymentMethodRepository(db)
	paymentMethodUseCase := usecase.NewPaymentMethodUseCase(paymentMethodRepo)
	paymentMethodHandler := httpDelivery.NewPaymentMethodHandler(paymentMethodUseCase)

	// Initialize beneficiary repository and usecase
	beneficiaryRepo := repository.NewBeneficiaryRepository(db)
	beneficiaryUseCase := usecase.NewBeneficiaryUseCase(beneficiaryRepo)
	beneficiaryHandler := httpDelivery.NewBeneficiaryHandler(beneficiaryUseCase)

	// Initialize handlers
	userHandler := httpDelivery.NewUserHandler(userUseCase)
	walletHandler := httpDelivery.NewWalletHandler(walletUseCase)
	transactionHandler := httpDelivery.NewTransactionHandler(transactionUseCase)

	// Initialize middleware
	mid := middleware.NewMiddleware(logger, cfg.JWT.Secret)

	// Initialize router
	router := mux.NewRouter()

	// Public routes
	router.HandleFunc("/api/register", userHandler.Register).Methods("POST")
	router.HandleFunc("/api/login", userHandler.Login).Methods("POST")

	// Protected routes
	api := router.PathPrefix("/api").Subrouter()
	api.Use(mid.AuthMiddleware)
	api.Use(mid.LoggingMiddleware)

	// User routes
	api.HandleFunc("/users/profile", userHandler.GetProfile).Methods("GET")
	api.HandleFunc("/users/profile", userHandler.UpdateProfile).Methods("PUT")
	api.HandleFunc("/users/password", userHandler.ChangePassword).Methods("PUT")

	// Wallet routes
	api.HandleFunc("/wallets", walletHandler.CreateWallet).Methods("POST")
	api.HandleFunc("/wallets", walletHandler.GetUserWallets).Methods("GET")
	api.HandleFunc("/wallets/{id}", walletHandler.GetWallet).Methods("GET")
	api.HandleFunc("/wallets/{id}/deactivate", walletHandler.DeactivateWallet).Methods("POST")
	api.HandleFunc("/wallets/transfer", walletHandler.Transfer).Methods("POST")
	api.HandleFunc("/wallets/{id}/deposit", walletHandler.Deposit).Methods("POST")
	api.HandleFunc("/wallets/{id}/withdraw", walletHandler.Withdraw).Methods("POST")

	// Transaction routes
	api.HandleFunc("/transactions", transactionHandler.GetUserTransactions).Methods("GET")

	// Payment Method routes
	api.HandleFunc("/payment-methods", paymentMethodHandler.GetUserPaymentMethods).Methods("GET")
	api.HandleFunc("/payment-methods", paymentMethodHandler.CreatePaymentMethod).Methods("POST")
	api.HandleFunc("/payment-methods/{id}", paymentMethodHandler.GetPaymentMethod).Methods("GET")
	api.HandleFunc("/payment-methods/{id}", paymentMethodHandler.UpdatePaymentMethod).Methods("PUT")
	api.HandleFunc("/payment-methods/{id}", paymentMethodHandler.DeletePaymentMethod).Methods("DELETE")
	api.HandleFunc("/payment-methods/{id}/set-default", paymentMethodHandler.SetDefaultPaymentMethod).Methods("PUT")

	// Beneficiary routes
	api.HandleFunc("/beneficiaries", beneficiaryHandler.GetUserBeneficiaries).Methods("GET")
	api.HandleFunc("/beneficiaries", beneficiaryHandler.CreateBeneficiary).Methods("POST")
	api.HandleFunc("/beneficiaries/{id}", beneficiaryHandler.GetBeneficiary).Methods("GET")
	api.HandleFunc("/beneficiaries/{id}", beneficiaryHandler.UpdateBeneficiary).Methods("PUT")
	api.HandleFunc("/beneficiaries/{id}", beneficiaryHandler.DeleteBeneficiary).Methods("DELETE")

	// Initialize repositories
	auditRepo := repository.NewAuditRepository(db)
	notificationRepo := repository.NewNotificationRepository(db)
	transactionLimitRepo := repository.NewTransactionLimitRepository(db)

	// Initialize use cases
	auditUseCase := usecase.NewAuditUseCase(auditRepo)
	notificationUseCase := usecase.NewNotificationUseCase(notificationRepo)
	transactionLimitUseCase := usecase.NewTransactionLimitUseCase(transactionLimitRepo)

	// Initialize handlers
	auditHandler := httpDelivery.NewAuditHandler(auditUseCase)
	notificationHandler := httpDelivery.NewNotificationHandler(notificationUseCase)
	transactionLimitHandler := httpDelivery.NewTransactionLimitHandler(transactionLimitUseCase)

	// Transaction Limits routes
	api.HandleFunc("/limits", transactionLimitHandler.SetLimit).Methods("POST")
	api.HandleFunc("/limits", transactionLimitHandler.GetLimits).Methods("GET")

	// Notifications routes
	api.HandleFunc("/notifications", notificationHandler.GetNotifications).Methods("GET")
	api.HandleFunc("/notifications/unread/count", notificationHandler.GetUnreadCount).Methods("GET")
	api.HandleFunc("/notifications/{id}/read", notificationHandler.MarkAsRead).Methods("PUT")
	api.HandleFunc("/notifications/read/all", notificationHandler.MarkAllAsRead).Methods("PUT")

	// Audit routes (with admin middleware)
	adminApi := router.PathPrefix("/api/admin").Subrouter()
	adminApi.Use(mid.AuthMiddleware)
	adminApi.Use(mid.AdminMiddleware)
	adminApi.Use(mid.LoggingMiddleware)

	adminApi.HandleFunc("/audit/logs", auditHandler.GetDateRangeLogs).Methods("GET")
	adminApi.HandleFunc("/audit/logs/action", auditHandler.GetActionLogs).Methods("GET")
	adminApi.HandleFunc("/audit/logs/entity", auditHandler.GetEntityLogs).Methods("GET")

	// User-specific audit logs are available through the regular API
	api.HandleFunc("/audit/logs", auditHandler.GetUserAuditLogs).Methods("GET")

	// Create server
	srv := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: router,
	}

	// Start server
	go func() {
		logger.Info("Server starting on port " + cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Server failed to start", "error", err)
			os.Exit(1)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown", "error", err)
	}

	logger.Info("Server stopped")
}
