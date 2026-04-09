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
	labOrdersHandler := handlers.NewLabOrdersHandler(pool)
	labCatalogHandler := handlers.NewLabCatalogHandler(pool)
	imagingHandler := handlers.NewImagingHandler(pool)
	prescriptionHandler := handlers.NewPrescriptionHandler(pool)
	pharmacyHandler := handlers.NewPharmacyHandler(pool)

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

		// Lab endpoints
		auth.POST("/lab/orders", labOrdersHandler.CreateOrder)
		auth.GET("/lab/orders/:id", labOrdersHandler.GetOrder)
		auth.GET("/lab/orders/pending", labOrdersHandler.ListPendingOrders)
		auth.POST("/lab/orders/items/:item_id/collect", labOrdersHandler.CollectSpecimen)
		auth.POST("/lab/orders/items/:item_id/results", labOrdersHandler.AddResult)
		auth.POST("/lab/results/:id/verify", labOrdersHandler.VerifyResult)
		auth.GET("/lab/critical-alerts", labOrdersHandler.ListCriticalAlerts)

		// Lab catalog endpoints
		auth.GET("/lab/tests", labCatalogHandler.ListTests)
		auth.GET("/lab/tests/:id", labCatalogHandler.GetTest)
		auth.GET("/lab/panels", labCatalogHandler.ListPanels)
		auth.GET("/lab/panels/:id", labCatalogHandler.GetPanel)

		// Imaging endpoints
		auth.POST("/imaging/orders", imagingHandler.CreateOrder)
		auth.GET("/imaging/orders/:id", imagingHandler.GetOrder)
		auth.GET("/imaging/orders", imagingHandler.ListOrders)
		auth.POST("/imaging/orders/:id/report", imagingHandler.AddReport)

		// Prescription endpoints
		auth.POST("/prescriptions", prescriptionHandler.CreatePrescription)
		auth.GET("/prescriptions/:id", prescriptionHandler.GetPrescription)
		auth.GET("/patients/:patient_id/prescriptions", prescriptionHandler.ListPrescriptionsByPatient)
		auth.POST("/prescriptions/:id/verify", prescriptionHandler.VerifyPrescription)
		auth.POST("/prescriptions/:id/dispense", prescriptionHandler.DispensePrescription)
		auth.POST("/prescriptions/:id/cancel", prescriptionHandler.CancelPrescription)

		// Pharmacy endpoints
		auth.GET("/medications", pharmacyHandler.ListMedications)
		auth.GET("/medications/search", pharmacyHandler.SearchMedications)
		auth.GET("/medications/:id", pharmacyHandler.GetMedication)
		auth.GET("/medications/:medication_id/stock", pharmacyHandler.GetStockLevels)
		auth.GET("/pharmacy/low-stock", pharmacyHandler.ListLowStock)
		auth.POST("/pharmacy/check-interaction", pharmacyHandler.CheckDrugInteraction)
		auth.GET("/medications/:medication_id/interactions", pharmacyHandler.ListDrugInteractions)
	}
}
