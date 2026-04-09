package server

import (
	"net/http"
	"time"

	"github.com/dmsafrica/dms/internal/config"
	"github.com/dmsafrica/dms/internal/http/middleware"
	"github.com/dmsafrica/dms/internal/http/routes"
	"github.com/dmsafrica/dms/internal/security"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func New(cfg config.Config, pool *pgxpool.Pool) *http.Server {
	// Optimized for low-resource: avoid default logger; add structured logging later.
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.RequestID())

	jwtSvc := security.NewJWTService(cfg.JWTSecret)
	routes.Register(r, jwtSvc, pool)

	return &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
}
