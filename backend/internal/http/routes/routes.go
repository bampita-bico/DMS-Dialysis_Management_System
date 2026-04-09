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
	sessionsHandler := handlers.NewSessionsHandler(pool)
	vitalsHandler := handlers.NewVitalsHandler(pool)
	machinesHandler := handlers.NewMachinesHandler(pool)
	waterHandler := handlers.NewWaterTreatmentHandler(pool)

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

		// Dialysis session endpoints
		auth.POST("/sessions", sessionsHandler.CreateSession)
		auth.GET("/sessions/:id", sessionsHandler.GetSession)
		auth.GET("/sessions", sessionsHandler.ListSessionsByDate)
		auth.POST("/sessions/:id/start", sessionsHandler.StartSession)
		auth.POST("/sessions/:id/complete", sessionsHandler.CompleteSession)
		auth.POST("/sessions/:id/abort", sessionsHandler.AbortSession)

		// Vitals endpoints
		auth.POST("/vitals", vitalsHandler.RecordVitals)
		auth.GET("/sessions/:session_id/vitals", vitalsHandler.ListVitalsBySession)
		auth.GET("/sessions/:session_id/vitals/alerts", vitalsHandler.ListVitalsWithAlerts)
		auth.POST("/vitals/:id/acknowledge-alert", vitalsHandler.AcknowledgeAlert)

		// Machine endpoints
		auth.POST("/machines", machinesHandler.CreateMachine)
		auth.GET("/machines", machinesHandler.ListMachines)
		auth.GET("/machines/available", machinesHandler.ListAvailableMachines)
		auth.PATCH("/machines/:id/status", machinesHandler.UpdateMachineStatus)

		// Water treatment endpoints
		auth.POST("/water-tests", waterHandler.LogWaterTest)
		auth.GET("/water-tests", waterHandler.ListWaterTests)
		auth.GET("/water-tests/failed", waterHandler.ListFailedWaterTests)
	}
}
