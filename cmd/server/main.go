// @title           Case Management API
// @version         1.0
// @description     This is a sample API documentation for our Go BFF application.
// @termsOfService  http://swagger.io/terms/

// @contact.name   SYE
// @contact.url    https://aeon.co.th
// @contact.email  sye@aeon.co.th

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host            localhost:8000
// @BasePath        /api/v1

package main

import (
	"case-management/pkg/config"
	"case-management/pkg/logger"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg, err := config.Load("./configs/config.yaml")
	if err != nil {
		log.Fatalf("FATAL: Failed to load configuration: %v", err)
	}

	appLogger := logger.New(cfg.Server.LogLevel)
	appLogger.Infof("Starting service: %s v%s", cfg.App.Name, cfg.App.Version)

	router, err := initializeApp(cfg, appLogger)
	if err != nil {
		appLogger.Fatalf("FATAL: Application initialization failed: %v", err)
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: router,
	}

	go func() {
		appLogger.Infof("Server starting on port %d", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			appLogger.Fatalf("Listen and serve error: %s", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	appLogger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		appLogger.Fatal("Server forced to shutdown:", err)
	}

	appLogger.Info("Server exiting gracefully.")
}
