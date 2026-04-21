package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/dmsafrica/dms/internal/db/sqlc"
	"github.com/dmsafrica/dms/internal/db/tenant"
	apierr "github.com/dmsafrica/dms/internal/http/errors"
	"github.com/dmsafrica/dms/internal/http/middleware"
)

type PrescriptionHandler struct {
	pool *pgxpool.Pool
}

func NewPrescriptionHandler(pool *pgxpool.Pool) *PrescriptionHandler {
	return &PrescriptionHandler{pool: pool}
}

func (h *PrescriptionHandler) CreatePrescription(c *gin.Context) {
	hospitalID := c.GetString(middleware.CtxHospitalID)

	var req struct {
		PatientID      uuid.UUID  `json:"patient_id" binding:"required"`
		SessionID      *uuid.UUID `json:"session_id"`
		PrescribedBy   uuid.UUID  `json:"prescribed_by" binding:"required"`
		PrescribedDate time.Time  `json:"prescribed_date" binding:"required"`
		PrescribedTime time.Time  `json:"prescribed_time" binding:"required"`
		ValidFrom      time.Time  `json:"valid_from" binding:"required"`
		ValidUntil     *time.Time `json:"valid_until"`
		Diagnosis      *string    `json:"diagnosis"`
		ClinicalNotes  *string    `json:"clinical_notes"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	tx, err := h.pool.Begin(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to begin transaction"})
		return
	}
	defer tx.Rollback(ctx)

	if err := tenant.SetLocalHospitalID(ctx, tx, hospitalID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set tenant context"})
		return
	}

	queries := sqlc.New(tx)

	var sessionID pgtype.UUID
	if req.SessionID != nil {
		sessionID = pgtype.UUID{Bytes: *req.SessionID, Valid: true}
	}

	var validUntil pgtype.Date
	if req.ValidUntil != nil {
		validUntil = pgtype.Date{Time: *req.ValidUntil, Valid: true}
	}

	prescription, err := queries.CreatePrescription(ctx, sqlc.CreatePrescriptionParams{
		HospitalID:     uuid.MustParse(hospitalID),
		PatientID:      req.PatientID,
		SessionID:      sessionID,
		PrescribedBy:   req.PrescribedBy,
		PrescribedDate: pgtype.Date{Time: req.PrescribedDate, Valid: true},
		PrescribedTime: pgtype.Time{Microseconds: int64(req.PrescribedTime.Hour()*3600+req.PrescribedTime.Minute()*60+req.PrescribedTime.Second()) * 1000000, Valid: true},
		ValidFrom:      pgtype.Date{Time: req.ValidFrom, Valid: true},
		ValidUntil:     validUntil,
		Diagnosis:      pgtype.Text{String: *req.Diagnosis, Valid: req.Diagnosis != nil},
		ClinicalNotes:  pgtype.Text{String: *req.ClinicalNotes, Valid: req.ClinicalNotes != nil},
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create prescription"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusCreated, prescription)
}

func (h *PrescriptionHandler) GetPrescription(c *gin.Context) {
	hospitalID := c.GetString(middleware.CtxHospitalID)
	prescriptionID := c.Param("id")

	prescriptionUUID, err := uuid.Parse(prescriptionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid prescription ID"})
		return
	}

	ctx := c.Request.Context()
	tx, err := h.pool.Begin(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to begin transaction"})
		return
	}
	defer tx.Rollback(ctx)

	if err := tenant.SetLocalHospitalID(ctx, tx, hospitalID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set tenant context"})
		return
	}

	queries := sqlc.New(tx)
	prescription, err := queries.GetPrescription(ctx, prescriptionUUID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Prescription not found"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, prescription)
}

func (h *PrescriptionHandler) ListPrescriptionsByPatient(c *gin.Context) {
	hospitalID := c.GetString(middleware.CtxHospitalID)
	patientID := c.Param("patient_id")

	patientUUID, err := uuid.Parse(patientID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid patient ID"})
		return
	}

	limit := int32(50)
	offset := int32(0)

	ctx := c.Request.Context()
	tx, err := h.pool.Begin(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to begin transaction"})
		return
	}
	defer tx.Rollback(ctx)

	if err := tenant.SetLocalHospitalID(ctx, tx, hospitalID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set tenant context"})
		return
	}

	queries := sqlc.New(tx)
	prescriptions, err := queries.ListMedicationPrescriptionsByPatient(ctx, sqlc.ListMedicationPrescriptionsByPatientParams{
		PatientID: patientUUID,
		Limit:     limit,
		Offset:    offset,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list prescriptions"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, prescriptions)
}

func (h *PrescriptionHandler) VerifyPrescription(c *gin.Context) {
	hospitalID := c.GetString(middleware.CtxHospitalID)
	staffID := c.GetString(middleware.CtxUserID)
	prescriptionID := c.Param("id")

	prescriptionUUID, err := uuid.Parse(prescriptionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid prescription ID"})
		return
	}

	staffUUID, _ := uuid.Parse(staffID)

	ctx := c.Request.Context()
	tx, err := h.pool.Begin(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to begin transaction"})
		return
	}
	defer tx.Rollback(ctx)

	if err := tenant.SetLocalHospitalID(ctx, tx, hospitalID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set tenant context"})
		return
	}

	queries := sqlc.New(tx)
	prescription, err := queries.VerifyPrescription(ctx, sqlc.VerifyPrescriptionParams{
		ID:                   prescriptionUUID,
		PharmacistVerifiedBy: pgtype.UUID{Bytes: staffUUID, Valid: true},
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify prescription"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, prescription)
}

func (h *PrescriptionHandler) DispensePrescription(c *gin.Context) {
	hospitalID := c.GetString(middleware.CtxHospitalID)
	staffID := c.GetString(middleware.CtxUserID)
	prescriptionID := c.Param("id")

	prescriptionUUID, err := uuid.Parse(prescriptionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid prescription ID"})
		return
	}

	staffUUID, _ := uuid.Parse(staffID)

	ctx := c.Request.Context()
	tx, err := h.pool.Begin(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to begin transaction"})
		return
	}
	defer tx.Rollback(ctx)

	if err := tenant.SetLocalHospitalID(ctx, tx, hospitalID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set tenant context"})
		return
	}

	queries := sqlc.New(tx)
	prescription, err := queries.DispensePrescription(ctx, sqlc.DispensePrescriptionParams{
		ID:          prescriptionUUID,
		DispensedBy: pgtype.UUID{Bytes: staffUUID, Valid: true},
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to dispense prescription"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, prescription)
}

func (h *PrescriptionHandler) CancelPrescription(c *gin.Context) {
	hospitalID := c.GetString(middleware.CtxHospitalID)
	staffID := c.GetString(middleware.CtxUserID)
	prescriptionID := c.Param("id")

	var req struct {
		CancellationReason *string `json:"cancellation_reason"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	prescriptionUUID, err := uuid.Parse(prescriptionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid prescription ID"})
		return
	}

	staffUUID, _ := uuid.Parse(staffID)

	ctx := c.Request.Context()
	tx, err := h.pool.Begin(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to begin transaction"})
		return
	}
	defer tx.Rollback(ctx)

	if err := tenant.SetLocalHospitalID(ctx, tx, hospitalID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set tenant context"})
		return
	}

	queries := sqlc.New(tx)
	prescription, err := queries.CancelPrescription(ctx, sqlc.CancelPrescriptionParams{
		ID:                 prescriptionUUID,
		CancelledBy:        pgtype.UUID{Bytes: staffUUID, Valid: true},
		CancellationReason: pgtype.Text{String: *req.CancellationReason, Valid: req.CancellationReason != nil},
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel prescription"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, prescription)
}

// CreatePrescriptionItem adds a medication line to a prescription with safety checks.
// POST /api/v1/prescriptions/:id/items
func (h *PrescriptionHandler) CreatePrescriptionItem(c *gin.Context) {
	hospitalID := c.GetString(middleware.CtxHospitalID)
	prescriptionID := c.Param("id")

	prescriptionUUID, err := uuid.Parse(prescriptionID)
	if err != nil {
		apierr.BadRequest(c, apierr.ErrInvalidID, "Invalid prescription ID")
		return
	}

	var req struct {
		MedicationID       uuid.UUID  `json:"medication_id" binding:"required"`
		Dose               string     `json:"dose" binding:"required"`
		Frequency          string     `json:"frequency" binding:"required"`
		Route              string     `json:"route" binding:"required"`
		DurationDays       *int32     `json:"duration_days"`
		QuantityPrescribed *int32     `json:"quantity_prescribed"`
		Instructions       *string    `json:"instructions"`
		StartDate          time.Time  `json:"start_date" binding:"required"`
		EndDate            *time.Time `json:"end_date"`
		IsPrn              bool       `json:"is_prn"`
		PrnIndication      *string    `json:"prn_indication"`
		IsStat             bool       `json:"is_stat"`
		Notes              *string    `json:"notes"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		apierr.BadRequest(c, apierr.ErrInvalidInput, err.Error())
		return
	}

	ctx := c.Request.Context()
	tx, err := h.pool.Begin(ctx)
	if err != nil {
		apierr.InternalWithCode(c, apierr.ErrDatabase, "Failed to begin transaction")
		return
	}
	defer tx.Rollback(ctx)

	if err := tenant.SetLocalHospitalID(ctx, tx, hospitalID); err != nil {
		apierr.InternalWithCode(c, apierr.ErrTransaction, "Failed to set tenant context")
		return
	}

	queries := sqlc.New(tx)

	// Fetch the prescription to get the patient_id
	prescription, err := queries.GetPrescription(ctx, prescriptionUUID)
	if err != nil {
		if err == pgx.ErrNoRows {
			apierr.NotFound(c, "Prescription")
			return
		}
		apierr.Internal(c, "Failed to retrieve prescription")
		return
	}

	// Fetch the medication to get its name for allergy checking
	medication, err := queries.GetMedication(ctx, req.MedicationID)
	if err != nil {
		if err == pgx.ErrNoRows {
			apierr.NotFound(c, "Medication")
			return
		}
		apierr.Internal(c, "Failed to retrieve medication")
		return
	}

	var warnings []apierr.Warning

	// Safety Check 1: Drug-allergy cross-check
	allergy, err := queries.CheckDrugAllergy(ctx, sqlc.CheckDrugAllergyParams{
		PatientID: prescription.PatientID,
		Allergen:  medication.GenericName,
	})
	if err == nil {
		// Allergy found — hard block
		apierr.ConflictWithDetails(c, apierr.ErrAllergyConflict,
			"Patient is allergic to "+medication.GenericName+" (severity: "+string(allergy.Severity)+")",
			gin.H{
				"allergy_id":    allergy.ID,
				"allergen":      allergy.Allergen,
				"severity":      allergy.Severity,
				"reaction":      allergy.Reaction,
				"medication_id": req.MedicationID,
			})
		return
	}

	// Also check drug class
	if medication.DrugClass.Valid && medication.DrugClass.String != "" {
		classAllergy, err := queries.CheckDrugAllergy(ctx, sqlc.CheckDrugAllergyParams{
			PatientID: prescription.PatientID,
			Allergen:  medication.DrugClass.String,
		})
		if err == nil {
			apierr.ConflictWithDetails(c, apierr.ErrAllergyConflict,
				"Patient is allergic to drug class "+medication.DrugClass.String+" ("+medication.GenericName+" belongs to this class)",
				gin.H{
					"allergy_id":    classAllergy.ID,
					"allergen":      classAllergy.Allergen,
					"severity":      classAllergy.Severity,
					"drug_class":    medication.DrugClass.String,
					"medication_id": req.MedicationID,
				})
			return
		}
	}

	// Safety Check 2: Drug-drug interaction check
	activeMedIDs, err := queries.ListActiveMedicationIDsForPatient(ctx, prescription.PatientID)
	if err != nil {
		apierr.Internal(c, "Failed to check active medications")
		return
	}

	for _, existingMedID := range activeMedIDs {
		if existingMedID == req.MedicationID {
			continue
		}
		interaction, err := queries.CheckInteraction(ctx, sqlc.CheckInteractionParams{
			MedicationAID: req.MedicationID,
			MedicationBID: existingMedID,
		})
		if err != nil {
			continue // no interaction for this pair
		}

		if interaction.Severity == "severe" || interaction.Severity == "contraindicated" {
			apierr.ConflictWithDetails(c, apierr.ErrDrugInteraction,
				"Severe drug interaction detected: "+interaction.Description,
				gin.H{
					"interaction_id":  interaction.ID,
					"severity":        interaction.Severity,
					"clinical_effect": interaction.ClinicalEffect,
					"recommendation":  interaction.ManagementRecommendation,
					"medication_a_id": interaction.MedicationAID,
					"medication_b_id": interaction.MedicationBID,
				})
			return
		}

		warnings = append(warnings, apierr.Warning{
			Code:    apierr.ErrDrugInteraction,
			Message: "Drug interaction (" + interaction.Severity + "): " + interaction.Description,
			Details: gin.H{
				"interaction_id":  interaction.ID,
				"severity":        interaction.Severity,
				"clinical_effect": interaction.ClinicalEffect,
				"recommendation":  interaction.ManagementRecommendation,
			},
		})
	}

	// All checks passed — create the item
	var durationDays pgtype.Int4
	if req.DurationDays != nil {
		durationDays = pgtype.Int4{Int32: *req.DurationDays, Valid: true}
	}
	var quantityPrescribed int32
	if req.QuantityPrescribed != nil {
		quantityPrescribed = *req.QuantityPrescribed
	}
	var endDate pgtype.Date
	if req.EndDate != nil {
		endDate = pgtype.Date{Time: *req.EndDate, Valid: true}
	}

	item, err := queries.CreatePrescriptionItem(ctx, sqlc.CreatePrescriptionItemParams{
		HospitalID:         uuid.MustParse(hospitalID),
		PrescriptionID:     prescriptionUUID,
		MedicationID:       req.MedicationID,
		Dose:               req.Dose,
		Frequency:          req.Frequency,
		Route:              sqlc.MedicationRoute(req.Route),
		DurationDays:       durationDays,
		QuantityPrescribed: quantityPrescribed,
		Instructions:       pgtype.Text{String: derefStr(req.Instructions), Valid: req.Instructions != nil},
		StartDate:          pgtype.Date{Time: req.StartDate, Valid: true},
		EndDate:            endDate,
		IsPrn:              req.IsPrn,
		PrnIndication:      pgtype.Text{String: derefStr(req.PrnIndication), Valid: req.PrnIndication != nil},
		IsStat:             req.IsStat,
		Notes:              pgtype.Text{String: derefStr(req.Notes), Valid: req.Notes != nil},
	})
	if err != nil {
		apierr.Internal(c, "Failed to create prescription item")
		return
	}

	if err := tx.Commit(ctx); err != nil {
		apierr.InternalWithCode(c, apierr.ErrTransaction, "Failed to commit transaction")
		return
	}

	if len(warnings) > 0 {
		apierr.CreatedWithWarnings(c, item, warnings)
		return
	}
	c.JSON(http.StatusCreated, item)
}

func derefStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
