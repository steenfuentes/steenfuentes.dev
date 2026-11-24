package database

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

type BackupManager struct {
	db         *DB
	bucketName string
	client     *storage.Client
	dbPath     string
}

// NewBackupManager creates a new backup manager
func NewBackupManager(db *DB, bucketName, dbPath string) (*BackupManager, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create storage client: %w", err)
	}

	return &BackupManager{
		db:         db,
		bucketName: bucketName,
		client:     client,
		dbPath:     dbPath,
	}, nil
}

// RestoreLatest restores the latest backup from GCS
func (bm *BackupManager) RestoreLatest(ctx context.Context) error {
	bucket := bm.client.Bucket(bm.bucketName)
	
	// List objects to find the latest backup
	var latestBackup *storage.ObjectAttrs
	it := bucket.Objects(ctx, &storage.Query{Prefix: "backups/portfolio-"})
	
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to list backups: %w", err)
		}
		
		if latestBackup == nil || attrs.Created.After(latestBackup.Created) {
			latestBackup = attrs
		}
	}
	
	if latestBackup == nil {
		log.Println("No backup found, starting with empty database")
		return nil
	}
	
	log.Printf("Restoring from backup: %s", latestBackup.Name)
	
	// Download the backup
	rc, err := bucket.Object(latestBackup.Name).NewReader(ctx)
	if err != nil {
		return fmt.Errorf("failed to read backup: %w", err)
	}
	defer rc.Close()
	
	// Write to local file
	f, err := os.Create(bm.dbPath)
	if err != nil {
		return fmt.Errorf("failed to create local db file: %w", err)
	}
	defer f.Close()
	
	if _, err := io.Copy(f, rc); err != nil {
		return fmt.Errorf("failed to write backup to disk: %w", err)
	}
	
	log.Println("Database restored successfully")
	return nil
}

// Backup creates a backup of the current database
func (bm *BackupManager) Backup(ctx context.Context) error {
	// Create backup filename with timestamp
	filename := fmt.Sprintf("backups/portfolio-%s.db", time.Now().UTC().Format("20060102-150405"))
	
	log.Printf("Creating backup: %s", filename)
	
	// Open the database file
	f, err := os.Open(bm.dbPath)
	if err != nil {
		return fmt.Errorf("failed to open database file: %w", err)
	}
	defer f.Close()
	
	// Upload to GCS
	wc := bm.client.Bucket(bm.bucketName).Object(filename).NewWriter(ctx)
	if _, err := io.Copy(wc, f); err != nil {
		wc.Close()
		return fmt.Errorf("failed to upload backup: %w", err)
	}
	
	if err := wc.Close(); err != nil {
		return fmt.Errorf("failed to close writer: %w", err)
	}
	
	log.Printf("Backup completed: %s", filename)
	return nil
}

// StartPeriodicBackups runs backups at specified interval
func (bm *BackupManager) StartPeriodicBackups(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			ctx := context.Background()
			if err := bm.Backup(ctx); err != nil {
				log.Printf("Backup failed: %v", err)
			}
		}
	}()
}

// Close closes the storage client
func (bm *BackupManager) Close() error {
	return bm.client.Close()
}