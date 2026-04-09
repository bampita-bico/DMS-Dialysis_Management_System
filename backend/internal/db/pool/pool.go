package pool

import (
	"context"
	"fmt"
	"time"

	"github.com/dmsafrica/dms/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func New(ctx context.Context, cfg config.Config) (*pgxpool.Pool, error) {
	url := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
		cfg.DBSSLMode,
	)

	pcfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, err
	}

	// Keep pool small for low-resource machines; scale via env.
	pcfg.MaxConns = cfg.DBMaxConns
	pcfg.MinConns = 0
	pcfg.MaxConnLifetime = 30 * time.Minute
	pcfg.MaxConnIdleTime = 5 * time.Minute
	pcfg.HealthCheckPeriod = 30 * time.Second

	p, err := pgxpool.NewWithConfig(ctx, pcfg)
	if err != nil {
		return nil, err
	}

	ctxPing, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := p.Ping(ctxPing); err != nil {
		p.Close()
		return nil, err
	}

	return p, nil
}
