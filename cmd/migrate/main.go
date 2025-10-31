package main

import (
	"fmt"
	"log"
	"os"

	"github.com/alexanderbkl/golang-hexagonal-boilerplate/pkg/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Create migration instance
	m, err := migrate.New(
		"file://migrations",
		cfg.Database.GetDSN(),
	)
	if err != nil {
		log.Fatalf("Failed to create migrate instance: %v", err)
	}

	// Run migrations based on command
	if len(os.Args) < 2 {
		fmt.Println("Usage: migrate [up|down|version]")
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Failed to run migrations: %v", err)
		}
		fmt.Println("Migrations applied successfully")
	case "down":
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Failed to rollback migrations: %v", err)
		}
		fmt.Println("Migrations rolled back successfully")
	case "version":
		version, dirty, err := m.Version()
		if err != nil {
			log.Fatalf("Failed to get migration version: %v", err)
		}
		fmt.Printf("Version: %d, Dirty: %v\n", version, dirty)
	default:
		fmt.Println("Unknown command. Use: up, down, or version")
		os.Exit(1)
	}
}
