package config

import (
	"fmt"
	"os"
	"regexp"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var validDatabaseName = regexp.MustCompile(`^[A-Za-z0-9_]+$`)

func ConnectDatabase() (*gorm.DB, error) {
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	if dbName == "" {
		return nil, fmt.Errorf("DB_NAME is required")
	}

	if !validDatabaseName.MatchString(dbName) {
		return nil, fmt.Errorf("DB_NAME %q is invalid: only letters, numbers and underscores are allowed", dbName)
	}

	serverDSN := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser,
		dbPassword,
		dbHost,
		dbPort,
	)

	serverDB, err := gorm.Open(mysql.Open(serverDSN), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MySQL server; check that MySQL is running and DB_USER/DB_PASSWORD are correct: %w", err)
	}

	sqlServerDB, err := serverDB.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get MySQL server instance: %w", err)
	}
	defer sqlServerDB.Close()

	if err := sqlServerDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping MySQL server; check that MySQL is running and DB_USER/DB_PASSWORD are correct: %w", err)
	}

	createDatabaseSQL := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci", dbName)
	if err := serverDB.Exec(createDatabaseSQL).Error; err != nil {
		return nil, fmt.Errorf("failed to create database %q; check that DB_USER has CREATE privileges: %w", dbName, err)
	}

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser,
		dbPassword,
		dbHost,
		dbPort,
		dbName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to open database %q; check DB_NAME and MySQL credentials: %w", dbName, err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database %q; check that MySQL is running and credentials are correct: %w", dbName, err)
	}

	return db, nil
}
