package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/steenfuentes/portfolio-backend/internal/api"
	"github.com/steenfuentes/portfolio-backend/internal/database"
)

func main() {
	// Get data directory
	dataDir := os.Getenv("DATA_DIR")
	if dataDir == "" {
		dataDir = "./data"
	}

	// Check if backup is enabled
	bucketName := os.Getenv("BACKUP_BUCKET")
	if bucketName != "" {
		log.Printf("Backup enabled with bucket: %s", bucketName)
		
		// Initialize empty database first
		db, err := database.New(dataDir)
		if err != nil {
			log.Fatalf("Failed to initialize database: %v", err)
		}
		db.Close()
		
		// Restore from backup
		dbPath := filepath.Join(dataDir, "portfolio.db")
		backupManager, err := database.NewBackupManager(nil, bucketName, dbPath)
		if err != nil {
			log.Printf("Failed to create backup manager: %v", err)
		} else {
			ctx := context.Background()
			if err := backupManager.RestoreLatest(ctx); err != nil {
				log.Printf("Failed to restore backup: %v", err)
			}
			defer backupManager.Close()
		}
	}

	// Initialize database
	db, err := database.New(dataDir)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Run migrations
	migrationsPath := os.Getenv("MIGRATIONS_PATH")
	if migrationsPath == "" {
		migrationsPath = "./migrations"
	}

	if err := db.Migrate(migrationsPath); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	
	// Start periodic backups if enabled
	if bucketName != "" {
		dbPath := filepath.Join(dataDir, "portfolio.db")
		backupManager, err := database.NewBackupManager(db, bucketName, dbPath)
		if err != nil {
			log.Printf("Warning: Failed to create backup manager: %v", err)
		} else {
			// Backup every 6 hours
			backupInterval := 6 * time.Hour
			if interval := os.Getenv("BACKUP_INTERVAL"); interval != "" {
				if d, err := time.ParseDuration(interval); err == nil {
					backupInterval = d
				}
			}
			
			log.Printf("Starting periodic backups every %v", backupInterval)
			backupManager.StartPeriodicBackups(backupInterval)
			
			// Create initial backup
			ctx := context.Background()
			if err := backupManager.Backup(ctx); err != nil {
				log.Printf("Initial backup failed: %v", err)
			}
		}
	}

	// Create server
	server := api.NewServer(db)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, server); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}