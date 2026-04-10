package routes

import (
	"github.com/dmsafrica/dms/internal/http/handlers"
	"github.com/dmsafrica/dms/internal/http/middleware"
	"github.com/dmsafrica/dms/internal/security"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

// RegisterWithTiers registers routes with tiered access control
// This is an EXAMPLE showing how to apply module-based middleware
func RegisterWithTiers(r *gin.Engine, jwtSvc *security.JWTService, pool *pgxpool.Pool) {
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
	equipmentHandler := handlers.NewEquipmentHandler(pool)
	consumablesHandler := handlers.NewConsumablesHandler(pool)

	// Protected group with JWT auth
	auth := r.Group("/api/v1")
	auth.Use(middleware.JWTAuth(jwtSvc))
	{
		// ====================================
		// CORE ENDPOINTS (All Tiers: Basic+)
		// ====================================

		auth.GET("/me", func(c *gin.Context) {
			c.JSON(200, gin.H{"hospital_id": c.GetString(middleware.CtxHospitalID)})
		})

		// Patient endpoints - Available in ALL tiers
		auth.POST("/patients", patientsHandler.Create)
		auth.GET("/patients", patientsHandler.List)
		auth.GET("/patients/search", patientsHandler.Search)
		auth.GET("/patients/:id", patientsHandler.Get)
		auth.DELETE("/patients/:id", patientsHandler.Delete)

		// Dialysis session endpoints - Available in ALL tiers
		auth.POST("/sessions", sessionsHandler.CreateSession)
		auth.GET("/sessions/:id", sessionsHandler.GetSession)
		auth.GET("/sessions", sessionsHandler.ListSessionsByDate)
		auth.POST("/sessions/:id/start", sessionsHandler.StartSession)
		auth.POST("/sessions/:id/complete", sessionsHandler.CompleteSession)
		auth.POST("/sessions/:id/abort", sessionsHandler.AbortSession)

		// Vitals endpoints - Available in ALL tiers
		auth.POST("/vitals", vitalsHandler.RecordVitals)
		auth.GET("/sessions/:session_id/vitals", vitalsHandler.ListVitalsBySession)
		auth.GET("/sessions/:session_id/vitals/alerts", vitalsHandler.ListVitalsWithAlerts)
		auth.POST("/vitals/:id/acknowledge-alert", vitalsHandler.AcknowledgeAlert)

		// Machine endpoints - Available in ALL tiers
		auth.POST("/machines", machinesHandler.CreateMachine)
		auth.GET("/machines", machinesHandler.ListMachines)
		auth.GET("/machines/available", machinesHandler.ListAvailableMachines)
		auth.PATCH("/machines/:id/status", machinesHandler.UpdateMachineStatus)

		// Water treatment endpoints - Available in ALL tiers (SAFETY CRITICAL)
		auth.POST("/water-tests", waterHandler.LogWaterTest)
		auth.GET("/water-tests", waterHandler.ListWaterTests)
		auth.GET("/water-tests/failed", waterHandler.ListFailedWaterTests)

		// ====================================
		// LAB MODULE (Enterprise Only)
		// ====================================

		lab := auth.Group("/lab")
		lab.Use(middleware.RequireModule(pool, "lab_management"))
		{
			// Lab orders
			lab.POST("/orders", labOrdersHandler.CreateOrder)
			lab.GET("/orders/:id", labOrdersHandler.GetOrder)
			lab.GET("/orders/pending", labOrdersHandler.ListPendingOrders)
			lab.POST("/orders/items/:item_id/collect", labOrdersHandler.CollectSpecimen)
			lab.POST("/orders/items/:item_id/results", labOrdersHandler.AddResult)
			lab.POST("/results/:id/verify", labOrdersHandler.VerifyResult)
			lab.GET("/critical-alerts", labOrdersHandler.ListCriticalAlerts)

			// Lab catalog
			lab.GET("/tests", labCatalogHandler.ListTests)
			lab.GET("/tests/:id", labCatalogHandler.GetTest)
			lab.GET("/panels", labCatalogHandler.ListPanels)
			lab.GET("/panels/:id", labCatalogHandler.GetPanel)
		}

		// ====================================
		// IMAGING MODULE (Enterprise Only)
		// ====================================

		imaging := auth.Group("/imaging")
		imaging.Use(middleware.RequireModule(pool, "imaging_integration"))
		{
			imaging.POST("/orders", imagingHandler.CreateOrder)
			imaging.GET("/orders/:id", imagingHandler.GetOrder)
			imaging.GET("/orders", imagingHandler.ListOrders)
			imaging.POST("/orders/:id/report", imagingHandler.AddReport)
		}

		// ====================================
		// PHARMACY MODULE (Enterprise Only)
		// ====================================

		pharmacy := auth.Group("/pharmacy")
		pharmacy.Use(middleware.RequireModule(pool, "full_pharmacy"))
		{
			// Prescription management
			pharmacy.POST("/prescriptions", prescriptionHandler.CreatePrescription)
			pharmacy.GET("/prescriptions/:id", prescriptionHandler.GetPrescription)
			pharmacy.GET("/patients/:patient_id/prescriptions", prescriptionHandler.ListPrescriptionsByPatient)
			pharmacy.POST("/prescriptions/:id/verify", prescriptionHandler.VerifyPrescription)
			pharmacy.POST("/prescriptions/:id/dispense", prescriptionHandler.DispensePrescription)
			pharmacy.POST("/prescriptions/:id/cancel", prescriptionHandler.CancelPrescription)

			// Medication management
			pharmacy.GET("/medications", pharmacyHandler.ListMedications)
			pharmacy.GET("/medications/search", pharmacyHandler.SearchMedications)
			pharmacy.GET("/medications/:id", pharmacyHandler.GetMedication)
			pharmacy.GET("/medications/:medication_id/stock", pharmacyHandler.GetStockLevels)
			pharmacy.GET("/low-stock", pharmacyHandler.ListLowStock)
			pharmacy.POST("/check-interaction", pharmacyHandler.CheckDrugInteraction)
			pharmacy.GET("/medications/:medication_id/interactions", pharmacyHandler.ListDrugInteractions)
		}

		// ====================================
		// EQUIPMENT & INVENTORY (Enterprise Only)
		// ====================================

		equipment := auth.Group("/equipment")
		equipment.Use(middleware.RequireModule(pool, "inventory_tracking"))
		{
			equipment.POST("", equipmentHandler.CreateEquipment)
			equipment.GET("", equipmentHandler.ListEquipment)
			equipment.PATCH("/:id/status", equipmentHandler.UpdateEquipmentStatus)
			equipment.POST("/faults", equipmentHandler.ReportFault)
			equipment.GET("/faults/unresolved", equipmentHandler.ListUnresolvedFaults)
		}

		consumables := auth.Group("/consumables")
		consumables.Use(middleware.RequireModule(pool, "inventory_tracking"))
		{
			consumables.GET("", consumablesHandler.ListConsumables)
			consumables.GET("/:consumable_id/inventory", consumablesHandler.GetInventoryLevels)
			consumables.GET("/inventory/low-stock", consumablesHandler.ListLowStock)
			consumables.GET("/inventory/expiring", consumablesHandler.ListExpiringStock)
			consumables.POST("/usage", consumablesHandler.RecordUsage)
			consumables.GET("/sessions/:session_id", consumablesHandler.ListSessionUsage)
		}
	}
}
