package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	action := flag.String("action", "up", "Migration action (up, down)")
	flag.Parse()

	dbConfig := getDBConfig()

	if err := createDatabaseIfNotExists(dbConfig); err != nil {
		log.Fatalf("Failed to ensure database exists: %v", err)
	}

	if err := runMigrations(dbConfig, *action); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	log.Println("Migrations completed successfully")
}

func getDBConfig() DBConfig {
	return DBConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", ""),
		Name:     getEnv("DB_NAME", "effective_mobile_db"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

func (c DBConfig) AdminConnectionString() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.SSLMode)
}

func (c DBConfig) ConnectionString() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Name, c.SSLMode)
}

func createDatabaseIfNotExists(cfg DBConfig) error {

	adminDb, err := sql.Open("postgres", cfg.AdminConnectionString())
	if err != nil {
		return fmt.Errorf("failed to connect to admin database: %w", err)
	}
	defer adminDb.Close()

	var exists bool
	err = adminDb.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)", cfg.Name).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check database existence: %w", err)
	}

	if !exists {
		_, err = adminDb.Exec(fmt.Sprintf("CREATE DATABASE %s", cfg.Name))
		if err != nil {
			return fmt.Errorf("failed to create database: %w", err)
		}
		log.Printf("Database %s created\n", cfg.Name)
	}

	return nil
}

func runMigrations(cfg DBConfig, action string) error {
	db, err := sql.Open("postgres", cfg.ConnectionString())
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create migration driver: %w", err)
	}

	migrationsPath := flag.String("migrations", "internal/migrations", "Path to migrations directory")
	flag.Parse()

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+*migrationsPath,
		"postgres", driver)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	switch action {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			return fmt.Errorf("failed to apply migrations: %w", err)
		}
	case "down":
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			return fmt.Errorf("failed to rollback migrations: %w", err)
		}
	default:
		return fmt.Errorf("unknown action: %s", action)
	}

	return nil
}
