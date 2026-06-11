// Package database wires up the PostgreSQL (GORM) and Redis connections and
// runs schema migrations.
package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/instaagrammeta/rentacar/backend-go/internal/config"
	"github.com/instaagrammeta/rentacar/backend-go/internal/models"
)

// Connect opens a GORM connection to PostgreSQL with sensible pool settings and
// retries, since the database container may still be starting up.
func Connect(cfg *config.Config) (*gorm.DB, error) {
	gormCfg := &gorm.Config{}
	if !cfg.Debug {
		gormCfg.Logger = logger.Default.LogMode(logger.Warn)
	}

	var db *gorm.DB
	var err error
	for attempt := 1; attempt <= 10; attempt++ {
		db, err = gorm.Open(postgres.Open(cfg.DatabaseURL), gormCfg)
		if err == nil {
			break
		}
		log.Printf("PostgreSQL not ready (attempt %d/10): %v", attempt, err)
		time.Sleep(3 * time.Second)
	}
	if err != nil {
		return nil, fmt.Errorf("could not connect to PostgreSQL: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}

// Migrate runs GORM AutoMigrate for every model.
func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(models.AllModels()...); err != nil {
		return fmt.Errorf("auto-migration failed: %w", err)
	}
	return nil
}

// ConnectRedis opens a Redis client and verifies connectivity. When Redis is
// disabled in the configuration it returns (nil, nil) and the application
// degrades gracefully (no caching / token blacklist).
func ConnectRedis(cfg *config.Config) (*redis.Client, error) {
	if !cfg.RedisEnabled {
		log.Println("Redis disabled via configuration")
		return nil, nil
	}

	client := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var err error
	for attempt := 1; attempt <= 5; attempt++ {
		if err = client.Ping(ctx).Err(); err == nil {
			return client, nil
		}
		log.Printf("Redis not ready (attempt %d/5): %v", attempt, err)
		time.Sleep(2 * time.Second)
	}
	return nil, fmt.Errorf("could not connect to Redis: %w", err)
}
