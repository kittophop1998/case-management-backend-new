package database

import (
	"case-management/internal/domain/model"
	"context"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
)

var DBStore *gorm.DB

func InitPostgresDBStore(dsn string) (*gorm.DB, error) {
	log.Println("Initial DB Store")
	log.Printf("Connecting to database with DSN: %s", dsn)

	// Parse DSN to pgx.Config
	cfg, err := pgx.ParseConfig(dsn)
	if err != nil {
		log.Printf("Failed to parse DSN: %v", err)
		return nil, err
	}

	// Open stdlib DB connection
	sqlDB := stdlib.OpenDB(*cfg)

	// Open with GORM
	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	if err != nil {
		log.Printf("Failed to connect to database via GORM: %v", err)
		return nil, err
	}

	// DB tuning
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Ping to test connection
	if err := sqlDB.PingContext(context.Background()); err != nil {
		log.Printf("Database ping failed: %v", err)
		return nil, err
	}

	log.Println("Connecting to database success")

	DBStore = db
	return db, nil
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&model.User{})
}
