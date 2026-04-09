package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/dmsafrica/dms/internal/db/sqlc"
	"github.com/dmsafrica/dms/internal/db/tenant"
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
