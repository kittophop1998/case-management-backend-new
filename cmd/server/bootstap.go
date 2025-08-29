package main

import (
	"case-management/infrastructure/config"
	"case-management/infrastructure/logger"
	"case-management/infrastructure/seed"
	"case-management/internal/app/handler/http"
	"case-management/internal/app/usecase"
	"case-management/internal/platform/api"
	"case-management/internal/platform/database"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func initializeApp(cfg *config.Config, appLogger *zap.SugaredLogger) (*gin.Engine, error) {
	appLogger.Info("Initializing application components...")

	// Setup monitoring with Prometheus
	// prom := monitoring.NewPrometheus("template_go_bff")

	// Setup database connection and run migrations
	db, err := setupDatabase(cfg, appLogger)
	if err != nil {
		return nil, fmt.Errorf("database setup failed: %w", err)
	}

	auditLogger := database.NewAsyncAuditLogger(db, 100)
	gracefulExit(auditLogger)

	// ##### Application Layer: Use Cases #####
	// Log repository
	logRepo := database.NewLogPg(db)
	logUsecase := usecase.NewLogUseCase(logRepo)

	// User repository
	userDBRepo := database.NewUserPg(db)
	userUsecase := usecase.NewUserUseCase(userDBRepo)

	// Auth repository
	authUsecase := usecase.NewAuthUseCase(logUsecase, userDBRepo)

	// Master data repository
	masterDataRepo := database.NewMasterDataPg(db)
	masterDataUsecase := usecase.NewMasterDataUseCase(masterDataRepo)

	// Permission repository
	permissionRepo := database.NewPermissionPg(db)
	permissionUsecase := usecase.NewPermissionUseCase(permissionRepo)

	// Case repository
	caseRepo := database.NewCasePg(db)
	caseUsecase := usecase.NewCaseUseCase(caseRepo)

	// Customer repository
	customerRepo := database.NewCustomerPg(db)
	customerUsecase := usecase.NewCustomerUseCase(customerRepo)

	// Dashboard repository
	dashboardAPIClient := api.NewDashboardAPIClient(cfg.Services.ConnectorAPI.BaseURL)
	dashboardUsecase := usecase.NewDashboardUseCase(dashboardAPIClient)

	// Queue repository
	queueRepo := database.NewQueuePg(db)
	queueUsecase := usecase.NewQueueUsecase(auditLogger, queueRepo)

	// ##### Application Layer: Handlers #####

	// Application Layer: HTTP handlers
	http.InitHandlers(http.HandlerDeps{
		UserUC:       userUsecase,
		MasterDataUC: masterDataUsecase,
		AuthUC:       authUsecase,
		PermissionUC: permissionUsecase,
		LogUC:        logUsecase,
		CaseUC:       caseUsecase,
		CustomerUC:   customerUsecase,
		DashboardUC:  dashboardUsecase,
		QueueUC:      queueUsecase,
	})

	// Setup Gin engine and middlewares
	gin.SetMode(cfg.Server.GinMode)
	router := gin.New()

	// Set Cors
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Length", "Content-Type", "Authorization", "X-Time-Zone"},
		MaxAge:       12 * time.Hour,
	}))

	router.Use(gin.Recovery())
	router.Use(logger.GinLogger(appLogger))
	router.RemoveExtraSlash = true

	// Register all HTTP routes
	http.SetupRoutes(router)
	appLogger.Info("HTTP routes configured")

	return router, nil
}

// setupDatabase creates a new database connection and runs auto-migrations
func setupDatabase(cfg *config.Config, logger *zap.SugaredLogger) (*gorm.DB, error) {
	// Build DSN string config
	var dsn string
	if cfg.Server.GinMode == "release" {
		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
			cfg.Database.Host, cfg.Database.User, cfg.Database.Password, cfg.Database.DBName, cfg.Database.Port, cfg.Database.SSLMode)
	} else {
		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
			cfg.Database.Host, cfg.Database.User, cfg.Database.Password, cfg.Database.DBName, cfg.Database.Port, cfg.Database.SSLMode)
	}

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

func gracefulExit(logger *database.AsyncAuditLogger) {
	go func() {
		// จับ signal (Ctrl+C หรือ SIGTERM)
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

		<-c
		log.Println("[INFO] Shutting down audit logger...")
		logger.Shutdown()
		os.Exit(0)
	}()
}
