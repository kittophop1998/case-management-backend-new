package main

import (
	"case-management/internal/app/handler/http"
	"case-management/internal/app/usecase"
	"case-management/internal/platform/database"
	"case-management/pkg/config"
	"case-management/pkg/logger"
	"fmt"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func initializeApp(cfg *config.Config, appLogger *zap.SugaredLogger) (*gin.Engine, error) {
	appLogger.Info("Initializing application components...")

	// Setup monitoring with Prometheus
	// prom := monitoring.NewPrometheus("template_go_bff")

	// Setup database connection and run migrations
	db, err := setupDatabase(cfg.Database, appLogger)
	if err != nil {
		return nil, fmt.Errorf("database setup failed: %w", err)
	}

	// --- Dependency Injection (from innermost to outermost layer) ---
	// Platform Layer: Concrete implementations for database and external clients

	// Domain Layer: Repositories
	// User repository
	userDBRepo := database.NewUserPg(db)
	userUsecase := usecase.NewUserUseCase(userDBRepo)

	// Master data repository
	masterDataRepo := database.NewMasterDataPg(db)
	masterDataUsecase := usecase.NewMasterDataUseCase(masterDataRepo)

	// Auth repository
	authRepo := database.NewAuthPg(db)
	authUsecase := usecase.NewAuthUseCase(userDBRepo, authRepo)

	// Application Layer: HTTP handlers
	http.InitHandlers(userUsecase, masterDataUsecase, authUsecase)

	// Setup Gin engine and middlewares
	gin.SetMode(cfg.Server.GinMode)
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(logger.GinLogger(appLogger))

	// Register all HTTP routes
	http.SetupRoutes(router)
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
