package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dmsafrica/dms/internal/config"
	"github.com/dmsafrica/dms/internal/db/pool"
	"github.com/dmsafrica/dms/internal/db/sqlc"
	"github.com/dmsafrica/dms/internal/http/server"
	"github.com/dmsafrica/dms/internal/services"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	ctx := context.Background()
	pg, err := pool.New(ctx, cfg)
	if err != nil {
		log.Fatalf("db: %v", err)
	}
	defer pg.Close()

	httpSrv := server.New(cfg, pg)

	// Start HTTP server
	go func() {
		log.Printf("http listening on %s", cfg.HTTPAddr)
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("http: %v", err)
		}
	}()

	// Start sync worker in background
	syncService := services.NewSyncService(pg)
	syncCtx, syncCancel := context.WithCancel(context.Background())
	defer syncCancel()

	go runSyncWorker(syncCtx, pg, syncService)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Println("Shutting down...")
	syncCancel() // Stop sync worker

	ctxShutdown, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = httpSrv.Shutdown(ctxShutdown)
}

// runSyncWorker processes sync queue periodically
func runSyncWorker(ctx context.Context, pg *pgxpool.Pool, syncService *services.SyncService) {
	queries := sqlc.New(pg)

	syncTicker := time.NewTicker(10 * time.Second)
	defer syncTicker.Stop()

	requeueTicker := time.NewTicker(5 * time.Minute)
	defer requeueTicker.Stop()

	cleanupTicker := time.NewTicker(1 * time.Hour)
	defer cleanupTicker.Stop()

	log.Println("Sync worker started (background goroutine)")

	for {
		select {
		case <-ctx.Done():
			log.Println("Sync worker stopped")
			return

		case <-syncTicker.C:
			// Process sync queue for all hospitals
			hospitals, err := queries.ListHospitals(ctx)
			if err != nil {
				continue
			}

			for _, hospital := range hospitals {
				processed, _ := syncService.ProcessSyncQueue(ctx, hospital.ID, 50)
				if processed > 0 {
					log.Printf("Synced %d items for hospital %s", processed, hospital.Name)
				}
			}

		case <-requeueTicker.C:
			// Requeue failed items
			hospitals, err := queries.ListHospitals(ctx)
			if err != nil {
				continue
			}

			for _, hospital := range hospitals {
				_ = syncService.RequeueFailedSyncs(ctx, hospital.ID)
			}

		case <-cleanupTicker.C:
			// Cleanup old synced items
			_ = syncService.CleanupOldSyncedItems(ctx, 7)
		}
	}
}
