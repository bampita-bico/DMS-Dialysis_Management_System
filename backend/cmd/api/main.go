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
	"github.com/dmsafrica/dms/internal/http/server"
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

	go func() {
		log.Printf("http listening on %s", cfg.HTTPAddr)
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("http: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctxShutdown, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = httpSrv.Shutdown(ctxShutdown)
}
