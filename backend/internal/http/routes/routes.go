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
	hospitalsHandler := handlers.NewHospitalsHandler(pool)
	usersHandler := handlers.NewUsersHandler(pool)
	subscriptionHandler := handlers.NewSubscriptionPlansHandler(pool)
	patientsHandler := handlers.NewPatientsHandler(pool)
	medicalHistoryHandler := handlers.NewMedicalHistoryHandler(pool)
	vascularAccessHandler := handlers.NewVascularAccessHandler(pool)
	clinicalOutcomesHandler := handlers.NewClinicalOutcomesHandler(pool)
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
	invoicesHandler := handlers.NewInvoicesHandler(pool)
	paymentsHandler := handlers.NewPaymentsHandler(pool)
	billingAccountsHandler := handlers.NewBillingAccountsHandler(pool)
	insuranceClaimsHandler := handlers.NewInsuranceClaimsHandler(pool)
	staffProfilesHandler := handlers.NewStaffProfilesHandler(pool)
	shiftAssignmentsHandler := handlers.NewShiftAssignmentsHandler(pool)
	leaveRecordsHandler := handlers.NewLeaveRecordsHandler(pool)
	mortalityRecordsHandler := handlers.NewMortalityRecordsHandler(pool)
	hospitalizationsHandler := handlers.NewHospitalizationsHandler(pool)
	dialysisSessionsHandler := handlers.NewDialysisSessionsHandler(pool)
	sessionComplicationsHandler := handlers.NewSessionComplicationsHandler(pool)
	sessionFluidBalanceHandler := handlers.NewSessionFluidBalanceHandler(pool)
	labResultsHandler := handlers.NewLabResultsHandler(pool)
	labCriticalAlertsHandler := handlers.NewLabCriticalAlertsHandler(pool)

	// Protected group
	auth := r.Group("/api/v1")
	auth.Use(middleware.JWTAuth(jwtSvc))
	{
		auth.GET("/me", func(c *gin.Context) {
			c.JSON(200, gin.H{"hospital_id": c.GetString(middleware.CtxHospitalID)})
		})

		// Hospital management endpoints
		auth.POST("/hospitals", hospitalsHandler.Create)
		auth.GET("/hospitals", hospitalsHandler.List)
		auth.GET("/hospitals/:id", hospitalsHandler.Get)
		auth.PATCH("/hospitals/:id", hospitalsHandler.Update)
		auth.DELETE("/hospitals/:id", hospitalsHandler.Delete)

		// User management endpoints
		auth.POST("/users", usersHandler.Create)
		auth.GET("/users", usersHandler.List)
		auth.GET("/users/:id", usersHandler.Get)
		auth.PATCH("/users/:id", usersHandler.Update)
		auth.DELETE("/users/:id", usersHandler.Delete)

		// Subscription management endpoints
		auth.GET("/subscription/plan", subscriptionHandler.GetCurrentPlan)
		auth.PUT("/subscription/plan", subscriptionHandler.UpdatePlan)
		auth.PUT("/subscription/modules", subscriptionHandler.UpdateModules)
		auth.GET("/subscription/plans", subscriptionHandler.ListPlans)

		// Patient endpoints
		auth.POST("/patients", patientsHandler.Create)
		auth.GET("/patients", patientsHandler.List)
		auth.GET("/patients/search", patientsHandler.Search)
		auth.GET("/patients/:id", patientsHandler.Get)
		auth.DELETE("/patients/:id", patientsHandler.Delete)

		// Patient Medical History - Diagnoses
		auth.POST("/patients/:patient_id/diagnoses", medicalHistoryHandler.CreateDiagnosis)
		auth.GET("/patients/:patient_id/diagnoses", medicalHistoryHandler.ListDiagnosesByPatient)
		auth.GET("/patients/:patient_id/diagnoses/primary", medicalHistoryHandler.GetPrimaryDiagnosis)

		// Patient Medical History - Comorbidities
		auth.POST("/patients/:patient_id/comorbidities", medicalHistoryHandler.CreateComorbidity)
		auth.GET("/patients/:patient_id/comorbidities", medicalHistoryHandler.ListComorbiditiesByPatient)
		auth.PATCH("/comorbidities/:id/status", medicalHistoryHandler.UpdateComorbidityStatus)

		// Patient Medical History - Allergies
		auth.POST("/patients/:patient_id/allergies", medicalHistoryHandler.CreateAllergy)
		auth.GET("/patients/:patient_id/allergies", medicalHistoryHandler.GetActiveAllergies)
		auth.GET("/patients/:patient_id/allergies/check", medicalHistoryHandler.CheckDrugAllergy)

		// Vascular Access Management
		auth.POST("/vascular-access", vascularAccessHandler.Create)
		auth.GET("/vascular-access/:id", vascularAccessHandler.Get)
		auth.PATCH("/vascular-access/:id", vascularAccessHandler.Update)
		auth.POST("/vascular-access/:id/abandon", vascularAccessHandler.Abandon)
		auth.GET("/patients/:patient_id/vascular-access", vascularAccessHandler.ListByPatient)
		auth.GET("/patients/:patient_id/vascular-access/primary", vascularAccessHandler.GetPrimary)

		// Clinical Outcomes
		auth.POST("/clinical-outcomes", clinicalOutcomesHandler.Create)
		auth.GET("/patients/:patient_id/clinical-outcomes", clinicalOutcomesHandler.ListByPatient)
		auth.GET("/patients/:patient_id/clinical-outcomes/latest", clinicalOutcomesHandler.GetLatest)
		auth.GET("/clinical-outcomes/declining", clinicalOutcomesHandler.ListDeclining)
		auth.GET("/clinical-outcomes/by-trend", clinicalOutcomesHandler.ListByTrend)

		// Dialysis session endpoints
		auth.POST("/sessions", sessionsHandler.CreateSession)
		auth.GET("/sessions/:id", sessionsHandler.GetSession)
		auth.GET("/sessions", sessionsHandler.ListSessionsByDate)
		auth.POST("/sessions/:id/start", sessionsHandler.StartSession)
		auth.POST("/sessions/:id/complete", sessionsHandler.CompleteSession)
		auth.POST("/sessions/:id/abort", sessionsHandler.AbortSession)

		// Dialysis Sessions (enhanced) endpoints
		auth.POST("/dialysis-sessions", dialysisSessionsHandler.Create)
		auth.GET("/dialysis-sessions/:id", dialysisSessionsHandler.Get)
		auth.GET("/patients/:patient_id/dialysis-sessions", dialysisSessionsHandler.ListByPatient)
		auth.GET("/dialysis-sessions/date/:scheduled_date", dialysisSessionsHandler.ListByDate)
		auth.GET("/machines/:machine_id/active-sessions", dialysisSessionsHandler.ListActiveByMachine)
		auth.GET("/dialysis-sessions/active", dialysisSessionsHandler.ListActive)
		auth.POST("/dialysis-sessions/:id/start", dialysisSessionsHandler.Start)
		auth.POST("/dialysis-sessions/:id/complete", dialysisSessionsHandler.Complete)
		auth.POST("/dialysis-sessions/:id/abort", dialysisSessionsHandler.Abort)
		auth.PATCH("/dialysis-sessions/:id/status", dialysisSessionsHandler.UpdateStatus)

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

		// Lab Results endpoints
		auth.POST("/lab-results", labResultsHandler.Create)
		auth.GET("/lab-results/:id", labResultsHandler.Get)
		auth.GET("/lab-order-items/:order_item_id/result", labResultsHandler.GetByOrderItem)
		auth.GET("/lab-orders/:order_id/results", labResultsHandler.ListByOrder)
		auth.GET("/lab-results/pending-verification", labResultsHandler.ListPendingVerification)
		auth.GET("/lab-results/critical", labResultsHandler.ListCritical)
		auth.POST("/lab-results/:id/verify", labResultsHandler.Verify)
		auth.PATCH("/lab-results/:id", labResultsHandler.Update)
		auth.DELETE("/lab-results/:id", labResultsHandler.Delete)

		// Lab Critical Alerts endpoints
		auth.POST("/lab-critical-alerts", labCriticalAlertsHandler.Create)
		auth.GET("/lab-critical-alerts/:id", labCriticalAlertsHandler.Get)
		auth.GET("/patients/:patient_id/critical-alerts", labCriticalAlertsHandler.ListByPatient)
		auth.GET("/lab-critical-alerts/unacknowledged", labCriticalAlertsHandler.ListUnacknowledged)
		auth.GET("/lab-critical-alerts", labCriticalAlertsHandler.ListByDateRange)
		auth.POST("/lab-critical-alerts/:id/acknowledge", labCriticalAlertsHandler.Acknowledge)
		auth.POST("/lab-critical-alerts/:id/notify-doctor", labCriticalAlertsHandler.NotifyDoctor)
		auth.DELETE("/lab-critical-alerts/:id", labCriticalAlertsHandler.Delete)

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

		// Equipment endpoints
		auth.POST("/equipment", equipmentHandler.CreateEquipment)
		auth.GET("/equipment", equipmentHandler.ListEquipment)
		auth.PATCH("/equipment/:id/status", equipmentHandler.UpdateEquipmentStatus)
		auth.POST("/equipment/faults", equipmentHandler.ReportFault)
		auth.GET("/equipment/faults/unresolved", equipmentHandler.ListUnresolvedFaults)

		// Consumables endpoints
		auth.GET("/consumables", consumablesHandler.ListConsumables)
		auth.GET("/consumables/:consumable_id/inventory", consumablesHandler.GetInventoryLevels)
		auth.GET("/consumables/inventory/low-stock", consumablesHandler.ListLowStock)
		auth.GET("/consumables/inventory/expiring", consumablesHandler.ListExpiringStock)
		auth.POST("/consumables/usage", consumablesHandler.RecordUsage)
		auth.GET("/sessions/:session_id/consumables", consumablesHandler.ListSessionUsage)

		// Invoice endpoints
		auth.POST("/invoices", invoicesHandler.Create)
		auth.GET("/invoices/:id", invoicesHandler.Get)
		auth.GET("/invoices/number/:invoice_number", invoicesHandler.GetByNumber)
		auth.GET("/patients/:patient_id/invoices", invoicesHandler.ListByPatient)
		auth.GET("/billing-accounts/:account_id/invoices", invoicesHandler.ListByAccount)
		auth.GET("/invoices", invoicesHandler.ListByStatus)
		auth.GET("/invoices/overdue", invoicesHandler.ListOverdue)
		auth.PATCH("/invoices/:id/status", invoicesHandler.UpdateStatus)
		auth.PATCH("/invoices/:id/payment", invoicesHandler.UpdatePayment)

		// Payment endpoints
		auth.POST("/payments", paymentsHandler.Create)
		auth.GET("/payments/:id", paymentsHandler.Get)
		auth.GET("/invoices/:invoice_id/payments", paymentsHandler.ListByInvoice)
		auth.GET("/patients/:patient_id/payments", paymentsHandler.ListByPatient)
		auth.GET("/payments", paymentsHandler.ListByDateRange)
		auth.GET("/payments/method/:method", paymentsHandler.ListByMethod)
		auth.GET("/payments/total", paymentsHandler.GetTotal)

		// Billing Account endpoints
		auth.POST("/billing-accounts", billingAccountsHandler.Create)
		auth.GET("/billing-accounts/:id", billingAccountsHandler.Get)
		auth.GET("/patients/:patient_id/billing-account", billingAccountsHandler.GetByPatient)
		auth.GET("/billing-accounts", billingAccountsHandler.List)
		auth.PATCH("/billing-accounts/:id/balance", billingAccountsHandler.UpdateBalance)
		auth.PATCH("/billing-accounts/:id/status", billingAccountsHandler.UpdateStatus)

		// Insurance Claims endpoints
		auth.POST("/insurance-claims", insuranceClaimsHandler.Create)
		auth.GET("/insurance-claims/:id", insuranceClaimsHandler.Get)
		auth.GET("/insurance-claims/number/:claim_number", insuranceClaimsHandler.GetByNumber)
		auth.GET("/invoices/:invoice_id/claims", insuranceClaimsHandler.ListByInvoice)
		auth.GET("/insurance-schemes/:scheme_id/claims", insuranceClaimsHandler.ListByScheme)
		auth.GET("/insurance-claims", insuranceClaimsHandler.ListByStatus)
		auth.GET("/insurance-claims/pending", insuranceClaimsHandler.ListPending)
		auth.POST("/insurance-claims/:id/submit", insuranceClaimsHandler.Submit)
		auth.POST("/insurance-claims/:id/approve", insuranceClaimsHandler.Approve)
		auth.POST("/insurance-claims/:id/reject", insuranceClaimsHandler.Reject)

		// Staff Profiles endpoints
		auth.POST("/staff-profiles", staffProfilesHandler.Create)
		auth.GET("/staff-profiles/:id", staffProfilesHandler.Get)
		auth.GET("/users/:user_id/staff-profile", staffProfilesHandler.GetByUser)
		auth.GET("/staff-profiles", staffProfilesHandler.List)
		auth.GET("/staff-profiles/cadre/:cadre", staffProfilesHandler.ListByCadre)
		auth.GET("/staff-profiles/active", staffProfilesHandler.ListActive)
		auth.GET("/departments/:department_id/staff", staffProfilesHandler.ListByDepartment)
		auth.GET("/staff-profiles/expiring-licenses", staffProfilesHandler.ListExpiringLicenses)
		auth.PATCH("/staff-profiles/:id", staffProfilesHandler.Update)

		// Shift Assignments endpoints
		auth.POST("/shift-assignments", shiftAssignmentsHandler.Create)
		auth.GET("/shift-assignments/:id", shiftAssignmentsHandler.Get)
		auth.GET("/shift-assignments/date/:shift_date", shiftAssignmentsHandler.ListByDate)
		auth.GET("/staff/:staff_id/shifts", shiftAssignmentsHandler.ListByStaff)
		auth.GET("/shift-assignments/unconfirmed", shiftAssignmentsHandler.ListUnconfirmed)
		auth.POST("/shift-assignments/:id/confirm", shiftAssignmentsHandler.Confirm)
		auth.POST("/shift-assignments/:id/clock-in", shiftAssignmentsHandler.ClockIn)
		auth.POST("/shift-assignments/:id/clock-out", shiftAssignmentsHandler.ClockOut)

		// Leave Records endpoints
		auth.POST("/leave-records", leaveRecordsHandler.Create)
		auth.GET("/leave-records/:id", leaveRecordsHandler.Get)
		auth.GET("/staff/:staff_id/leave", leaveRecordsHandler.ListByStaff)
		auth.GET("/leave-records/pending", leaveRecordsHandler.ListPending)
		auth.GET("/leave-records", leaveRecordsHandler.ListByDateRange)
		auth.POST("/leave-records/:id/approve", leaveRecordsHandler.Approve)
		auth.POST("/leave-records/:id/reject", leaveRecordsHandler.Reject)

		// Mortality Records endpoints
		auth.POST("/mortality-records", mortalityRecordsHandler.Create)
		auth.GET("/mortality-records/:id", mortalityRecordsHandler.Get)
		auth.GET("/patients/:patient_id/mortality-record", mortalityRecordsHandler.GetByPatient)
		auth.GET("/mortality-records", mortalityRecordsHandler.List)
		auth.GET("/mortality-records/period", mortalityRecordsHandler.ListByPeriod)
		auth.GET("/mortality-records/session-related", mortalityRecordsHandler.ListSessionRelated)
		auth.GET("/mortality-records/setting/:setting", mortalityRecordsHandler.ListBySetting)
		auth.POST("/mortality-records/:id/certify", mortalityRecordsHandler.Certify)

		// Hospitalizations endpoints
		auth.POST("/hospitalizations", hospitalizationsHandler.Create)
		auth.GET("/hospitalizations/:id", hospitalizationsHandler.Get)
		auth.GET("/patients/:patient_id/hospitalizations", hospitalizationsHandler.ListByPatient)
		auth.GET("/hospitalizations", hospitalizationsHandler.ListByPeriod)
		auth.GET("/hospitalizations/dialysis-related", hospitalizationsHandler.ListDialysisRelated)
		auth.GET("/hospitalizations/access-related", hospitalizationsHandler.ListAccessRelated)
		auth.PATCH("/hospitalizations/:id/discharge", hospitalizationsHandler.UpdateDischarge)

		// Session Complications endpoints
		auth.POST("/session-complications", sessionComplicationsHandler.Create)
		auth.GET("/session-complications/:id", sessionComplicationsHandler.Get)
		auth.GET("/sessions/:session_id/complications", sessionComplicationsHandler.ListBySession)
		auth.GET("/patients/:patient_id/complications", sessionComplicationsHandler.ListByPatient)
		auth.GET("/session-complications/severe", sessionComplicationsHandler.ListSevere)
		auth.PATCH("/session-complications/:id", sessionComplicationsHandler.Update)
		auth.DELETE("/session-complications/:id", sessionComplicationsHandler.Delete)

		// Session Fluid Balance endpoints
		auth.POST("/session-fluid-balance", sessionFluidBalanceHandler.Create)
		auth.GET("/session-fluid-balance/:id", sessionFluidBalanceHandler.Get)
		auth.GET("/sessions/:session_id/fluid-balance", sessionFluidBalanceHandler.GetBySession)
		auth.GET("/patients/:patient_id/fluid-balance", sessionFluidBalanceHandler.ListByPatient)
		auth.PATCH("/session-fluid-balance/:id", sessionFluidBalanceHandler.Update)
		auth.DELETE("/session-fluid-balance/:id", sessionFluidBalanceHandler.Delete)
	}
}
