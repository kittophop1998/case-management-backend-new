package main

import (
	"case-management/infrastructure/config"
	"case-management/infrastructure/logger"
	"case-management/infrastructure/seed"
	"case-management/internal/app/handler/http"
	"case-management/internal/app/usecase"
	"case-management/internal/platform/database"

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

	// ##### Application Layer: Use Cases #####

	// User repository
	userDBRepo := database.NewUserPg(db)
	userUsecase := usecase.NewUserUseCase(userDBRepo)

	// Master data repository
	masterDataRepo := database.NewMasterDataPg(db)
	masterDataUsecase := usecase.NewMasterDataUseCase(masterDataRepo)

	// Auth repository
	authRepo := database.NewAuthPg(db)
	authUsecase := usecase.NewAuthUseCase(userDBRepo, authRepo)

	// Permission repository
	permissionRepo := database.NewPermissionPg(db)
	permissionUsecase := usecase.NewPermissionUseCase(permissionRepo)

	// Log repository
	logRepo := database.NewLogPg(db)
	logUsecase := usecase.NewLogUseCase(logRepo)

	// Case repository
	caseRepo := database.NewCasePg(db)
	caseUsecase := usecase.NewCaseUseCase(caseRepo)

	// Customer repository
	customerRepo := database.NewCustomerPg(db)
	customerUsecase := usecase.NewCustomerUseCase(customerRepo)

	// ##### Application Layer: Handlers #####

	// Application Layer: HTTP handlers
	http.InitHandlers(
		userUsecase,
		masterDataUsecase,
		authUsecase,
		permissionUsecase,
		logUsecase,
		caseUsecase,
		customerUsecase,
	)

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
	_ = seed.SeedAllData(db)

	logger.Info("Database migrations completed")

	return db, nil
}
