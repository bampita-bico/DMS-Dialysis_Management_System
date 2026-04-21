package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dmsafrica/dms/internal/config"
	"github.com/dmsafrica/dms/internal/db/sqlc"
	"github.com/dmsafrica/dms/internal/services"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	log.Println("Starting DMS Sync Worker...")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to database
	pool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	log.Println("Connected to database")

	// Create sync service
	syncService := services.NewSyncService(pool)

	// Setup graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Channel to listen for OS signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Start worker goroutine
	go runSyncWorker(ctx, pool, syncService)

	// Wait for termination signal
	<-sigChan
	log.Println("Shutdown signal received, stopping sync worker...")
	cancel()

	// Give worker time to finish current batch
	time.Sleep(2 * time.Second)
	log.Println("Sync worker stopped")
}

// runSyncWorker processes sync queue periodically
func runSyncWorker(ctx context.Context, pool *pgxpool.Pool, syncService *services.SyncService) {
	queries := sqlc.New(pool)

	// Sync every 10 seconds
	syncTicker := time.NewTicker(10 * time.Second)
	defer syncTicker.Stop()

	// Requeue failed items every 5 minutes
	requeueTicker := time.NewTicker(5 * time.Minute)
	defer requeueTicker.Stop()

	// Cleanup old items every hour
	cleanupTicker := time.NewTicker(1 * time.Hour)
	defer cleanupTicker.Stop()

	log.Println("Sync worker started - processing every 10 seconds")

	for {
		select {
		case <-ctx.Done():
			log.Println("Sync worker context cancelled")
			return

		case <-syncTicker.C:
			// Process sync queue for all hospitals
			processAllHospitals(ctx, queries, syncService)

		case <-requeueTicker.C:
			// Requeue failed items
			log.Println("Requeuing failed sync items...")
			requeueFailedForAllHospitals(ctx, queries, syncService)

		case <-cleanupTicker.C:
			// Cleanup old synced items (older than 7 days)
			log.Println("Cleaning up old synced items...")
			if err := syncService.CleanupOldSyncedItems(ctx, 7); err != nil {
				log.Printf("Cleanup failed: %v", err)
			}
		}
	}
}

// processAllHospitals processes sync queue for all active hospitals
func processAllHospitals(ctx context.Context, queries *sqlc.Queries, syncService *services.SyncService) {
	// Get all active hospitals
	hospitals, err := queries.ListHospitals(ctx, sqlc.ListHospitalsParams{
		Limit:  100,
		Offset: 0,
	})
	if err != nil {
		log.Printf("Failed to list hospitals: %v", err)
		return
	}

	if len(hospitals) == 0 {
		return // No hospitals to process
	}

	totalProcessed := 0
	for _, hospital := range hospitals {
		processed, err := syncService.ProcessSyncQueue(ctx, hospital.ID, 50)
		if err != nil {
			log.Printf("Failed to process sync queue for hospital %s: %v", hospital.ID, err)
			continue
		}

		if processed > 0 {
			totalProcessed += processed
			log.Printf("Hospital %s (%s): processed %d items", hospital.Name, hospital.ID, processed)
		}
	}

	if totalProcessed > 0 {
		log.Printf("Total items processed across all hospitals: %d", totalProcessed)
	}
}

// requeueFailedForAllHospitals requeues failed sync items for all hospitals
func requeueFailedForAllHospitals(ctx context.Context, queries *sqlc.Queries, syncService *services.SyncService) {
	hospitals, err := queries.ListHospitals(ctx, sqlc.ListHospitalsParams{
		Limit:  100,
		Offset: 0,
	})
	if err != nil {
		log.Printf("Failed to list hospitals for requeue: %v", err)
		return
	}

	totalRequeued := 0
	for _, hospital := range hospitals {
		count, err := syncService.RequeueFailedSyncs(ctx, hospital.ID)
		if err != nil {
			log.Printf("Failed to requeue for hospital %s: %v", hospital.ID, err)
			continue
		}
		totalRequeued += count
	}

	if totalRequeued > 0 {
		log.Printf("Total failed items requeued: %d", totalRequeued)
	}
}
