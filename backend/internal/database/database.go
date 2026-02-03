package database

import (
	"fmt"
	"log"
	"time"

	"docker-heatmap/internal/config"
	"docker-heatmap/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect() error {
	var err error

	logLevel := logger.Silent
	if config.AppConfig.Environment == "development" {
		logLevel = logger.Info
	}

	DB, err = gorm.Open(postgres.Open(config.AppConfig.DatabaseURL), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return err
	}

	// Get underlying SQL DB
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}

	// Connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("Database connected successfully")
	return nil
}

func Migrate() error {
	log.Println("Running database migrations...")

	// Drop existing tables if they have wrong schema (development only)
	if config.AppConfig.Environment == "development" {
		if err := fixSchemaIfNeeded(); err != nil {
			log.Printf("Schema fix warning: %v", err)
		}
	}

	return DB.AutoMigrate(
		&models.User{},
		&models.DockerAccount{},
		&models.ActivityEvent{},
	)
}

// fixSchemaIfNeeded checks for column naming issues and fixes them
func fixSchemaIfNeeded() error {
	// Check if old column exists
	var count int64
	err := DB.Raw(`
		SELECT COUNT(*) FROM information_schema.columns 
		WHERE table_name = 'users' AND column_name = 'git_hub_id'
	`).Scan(&count).Error

	if err != nil {
		return fmt.Errorf("failed to check schema: %w", err)
	}

	if count > 0 {
		log.Println("Fixing schema: renaming columns with wrong naming convention...")

		// Rename columns in users table
		migrations := []string{
			`ALTER TABLE users RENAME COLUMN git_hub_id TO github_id`,
			`ALTER TABLE users RENAME COLUMN git_hub_username TO github_username`,
			`ALTER TABLE users RENAME COLUMN git_hub_email TO github_email`,
		}

		for _, migration := range migrations {
			if err := DB.Exec(migration).Error; err != nil {
				log.Printf("Migration warning (may already be correct): %v", err)
			}
		}

		// Drop and recreate the unique index with correct name
		DB.Exec(`DROP INDEX IF EXISTS idx_users_git_hub_id`)
		DB.Exec(`CREATE UNIQUE INDEX IF NOT EXISTS idx_users_github_id ON users(github_id)`)

		log.Println("Schema fix completed")
	}

	// Drop raw_data column from activity_events if it exists (causes JSON insert issues)
	var rawDataExists int64
	DB.Raw(`
		SELECT COUNT(*) FROM information_schema.columns 
		WHERE table_name = 'activity_events' AND column_name = 'raw_data'
	`).Scan(&rawDataExists)

	if rawDataExists > 0 {
		log.Println("Dropping raw_data column from activity_events...")
		if err := DB.Exec(`ALTER TABLE activity_events DROP COLUMN raw_data`).Error; err != nil {
			log.Printf("Failed to drop raw_data column: %v", err)
		} else {
			log.Println("Dropped raw_data column successfully")
		}
	}

	return nil
}

func Close() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
