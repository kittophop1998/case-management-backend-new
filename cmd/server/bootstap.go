package main

import (
	"case-management/internal/app/usecase"
	"case-management/internal/platform/database"
	"case-management/pkg/config"
	"case-management/pkg/logger"
	"case-management/pkg/monitoring"
	"fmt"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// initializeApp wires up the application's dependencies
func initializeApp(cfg *config.Config, appLogger *zap.SugaredLogger) (*gin.Engine, error) {
	appLogger.Info("Initializing application components...")

	// Setup monitoring with Prometheus
	prom := monitoring.NewPrometheus("template_go_bff")

	// Setup database connection and run migrations
	db, err := setupDatabase(cfg.Database, appLogger)
	if err != nil {
		return nil, fmt.Errorf("database setup failed: %w", err)
	}

	// --- Dependency Injection (from innermost to outermost layer) ---

	// Platform Layer: Concrete implementations for database and external clients
	userDBRepo := database.NewUserPg(db)
	// agreementAPIClient := client.NewAgreementAPIClient(cfg.Services.AgreementAPI.BaseURL)

	// Usecase Layer: Core business logic, depends on platform implementations
	userUsecase := usecase.NewUserUsecase(userDBRepo)
	// agreementUsecase := usecase.NewAgreementUsecase(agreementAPIClient)

	// Handler Layer: Handles HTTP requests, depends on use cases
	userHandler := v1.NewUserHandler(userUsecase)
	// agreementHandler := v1.NewAgreementHandler(agreementUsecase)

	// Setup Gin engine and middlewares
	gin.SetMode(cfg.Server.GinMode)
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(logger.GinLogger(appLogger))
	// router.Use(v1.RequestIDMiddleware())

	// Register all HTTP routes
	// http_handler.SetupRoutes(router, prom, customerHandler, agreementHandler)
	appLogger.Info("HTTP routes configured")

	return router, nil
}

// setupDatabase creates a new database connection and runs auto-migrations
func setupDatabase(dbConfig config.DatabaseConfig, logger *zap.SugaredLogger) (*gorm.DB, error) {
	// Build DSN string config
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		dbConfig.Host, dbConfig.User, dbConfig.Password, dbConfig.DBName, dbConfig.Port, dbConfig.SSLMode)

	// Establish connection to the database
	db, err := database.InitPostgresDBStore(dsn)
	if err != nil {
		return nil, err
	}

	// Run GORM auto-migrations for domain models
	logger.Info("Running database migrations...")
	err = database.Migrate(db)
	if err != nil {
		return nil, fmt.Errorf("database migration failed: %w", err)
	}
	logger.Info("Database migrations completed")

	return db, nil
}
