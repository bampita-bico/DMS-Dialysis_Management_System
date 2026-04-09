package routes

import (
	"github.com/dmsafrica/dms/internal/http/handlers"
	"github.com/dmsafrica/dms/internal/http/middleware"
	"github.com/dmsafrica/dms/internal/security"
	"github.com/gin-gonic/gin"
)

func Register(r *gin.Engine, jwtSvc *security.JWTService) {
	r.GET("/health", handlers.Health)

	// Protected group (placeholder)
	auth := r.Group("/api")
	auth.Use(middleware.JWTAuth(jwtSvc))
	auth.GET("/me", func(c *gin.Context) {
		c.JSON(200, gin.H{"hospital_id": c.GetString(middleware.CtxHospitalID)})
	})
}
