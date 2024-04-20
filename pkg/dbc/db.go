package dbc

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log/slog"
	"time"
)

var DB *gorm.DB

func DatabaseConnection() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("posts.db"), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
	if err != nil {
		return nil, err
	}

	// Get the generic database interface of *gorm.DB
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}
	sqlDB.SetMaxOpenConns(1) // Only one open connection at a time
	sqlDB.SetConnMaxLifetime(time.Minute * 20)
	sqlDB.SetMaxIdleConns(1) // Reduce the number of idle connections

	return db, nil
}

func GetDB() *gorm.DB {
	// Check if DB is initialized
	if DB == nil {
		// Try to establish a new connection if DB is nil
		return establishNewConnection()
	}

	// Get the underlying sql.DB object from the GORM DB instance
	sqlDB, err := DB.DB()
	if err != nil {
		// If error in getting the sql.DB, try to establish a new connection
		return establishNewConnection()
	}

	// Ping the database to check the connection
	if err := sqlDB.Ping(); err != nil {
		// If ping fails, try to establish a new connection
		return establishNewConnection()
	}

	// Return the existing DB instance if all checks pass
	return DB
}

func establishNewConnection() *gorm.DB {
	newDB, err := DatabaseConnection()
	if err != nil {
		slog.Error("Error establishing new database connection", "error", err.Error())
		return nil
	}
	DB = newDB
	return DB
}

func CloseDB(db *gorm.DB) {
	dbInstance, _ := db.DB()
	err := dbInstance.Close()
	if err != nil {
		slog.Error("Error while closing DB connection. Not a problem actually", "error", err.Error())
	} else {
		slog.Info("DB connection is closed successfully")
	}
}
