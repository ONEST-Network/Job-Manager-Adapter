package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/VinVorteX/beneficiary-manager/internal/api"
	"github.com/VinVorteX/beneficiary-manager/internal/config"
	"github.com/VinVorteX/beneficiary-manager/internal/db"
	"github.com/VinVorteX/beneficiary-manager/internal/logger"
	"github.com/VinVorteX/beneficiary-manager/internal/service"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	if err := logger.Initialize(cfg.Logging.Level, cfg.Logging.PrettyPrint); err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}

	appLogger := logger.NewLogger()
	appLogger.Info("Starting Beneficiary Manager service", map[string]interface{}{
		"version": "1.0.0",
	})

	// Initialize DB with migrations
	postgres, err := db.NewPostgresDB(
		cfg.Database.GetDSN(),
		"migrations/postgres",
		appLogger,
	)
	if err != nil {
		appLogger.Error("Failed to initialize database", err, nil)
		os.Exit(1)
	}

	// Configure database connection pool
	if err := postgres.ConfigurePool(cfg.Database.MaxOpenConns, cfg.Database.MaxIdleConns); err != nil {
		appLogger.Error("Failed to configure database pool", err, nil)
		os.Exit(1)
	}

	// Initialize Service
	beneficiaryService := service.NewBeneficiaryService(postgres, appLogger)

	// Initialize Handler
	handler := api.NewHandler(beneficiaryService)

	// Create server
	server := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// Setup routes
	http.HandleFunc("/schemes", handler.GetSchemes)
	http.HandleFunc("/applications", handler.SubmitApplication)
	http.HandleFunc("/status", handler.GetApplicationStatus)

	// Start server
	go func() {
		appLogger.Info("Server starting", map[string]interface{}{
			"port": cfg.Server.Port,
		})
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			appLogger.Error("Server failed to start", err, nil)
			os.Exit(1)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	appLogger.Info("Server shutting down", nil)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		appLogger.Error("Server forced to shutdown", err, nil)
	}

	appLogger.Info("Server exited", nil)
}
