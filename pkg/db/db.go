package db

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"time"
)

func Connect(host, port, user, password, dbName string) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbName,
	)

	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	}

	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		log.Printf("[ERROR] Failed to connect to database: %v", err)
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Printf("[ERROR] Failed to get database instance: %v", err)
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	if err := sqlDB.Ping(); err != nil {
		log.Printf("[ERROR] Database connection test failed: %v", err)
		return nil, fmt.Errorf("database connection test failed: %w", err)
	}

	log.Println("[INFO] Database connected successfully")
	return db, nil
}
