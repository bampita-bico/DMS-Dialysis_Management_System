package routes

import (
	"github.com/dmsafrica/dms/internal/http/handlers"
	"github.com/dmsafrica/dms/internal/http/middleware"
	"github.com/dmsafrica/dms/internal/security"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Register(r *gin.Engine, jwtSvc *security.JWTService, pool *pgxpool.Pool) {
	r.GET("/health", handlers.Health)

	// Initialize handlers
	patientsHandler := handlers.NewPatientsHandler(pool)

	// Protected group
	auth := r.Group("/api/v1")
	auth.Use(middleware.JWTAuth(jwtSvc))
	{
		auth.GET("/me", func(c *gin.Context) {
			c.JSON(200, gin.H{"hospital_id": c.GetString(middleware.CtxHospitalID)})
		})

		// Patient endpoints
		auth.POST("/patients", patientsHandler.Create)
		auth.GET("/patients", patientsHandler.List)
		auth.GET("/patients/search", patientsHandler.Search)
		auth.GET("/patients/:id", patientsHandler.Get)
		auth.DELETE("/patients/:id", patientsHandler.Delete)
	}
}
