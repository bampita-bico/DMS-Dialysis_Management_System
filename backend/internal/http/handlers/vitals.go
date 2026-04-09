package handlers

import (
	"net/http"
	"time"

	"github.com/dmsafrica/dms/internal/db/sqlc"
	"github.com/dmsafrica/dms/internal/db/tenant"
	"github.com/dmsafrica/dms/internal/http/middleware"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type VitalsHandler struct {
	pool *pgxpool.Pool
}

func NewVitalsHandler(pool *pgxpool.Pool) *VitalsHandler {
	return &VitalsHandler{pool: pool}
}

func (h *VitalsHandler) RecordVitals(c *gin.Context) {
	var req struct {
		SessionID             uuid.UUID  `json:"session_id" binding:"required"`
		PatientID             uuid.UUID  `json:"patient_id" binding:"required"`
		TimeOnDialysisMins    *int32     `json:"time_on_dialysis_mins"`
		BPSystolic            *int32     `json:"bp_systolic"`
		BPDiastolic           *int32     `json:"bp_diastolic"`
		HeartRate             *int32     `json:"heart_rate"`
		Temperature           *float64   `json:"temperature"`
		SpO2                  *int32     `json:"spo2"`
		RespiratoryRate       *int32     `json:"respiratory_rate"`
		BloodFlowActual       *int32     `json:"blood_flow_actual"`
		DialysateFlowActual   *int32     `json:"dialysate_flow_actual"`
		VenousPressure        *int32     `json:"venous_pressure"`
		ArterialPressure      *int32     `json:"arterial_pressure"`
		TMP                   *int32     `json:"tmp"`
		UFRemovedSoFar        *float64   `json:"uf_removed_so_far"`
		ConductivityActual    *float64   `json:"conductivity_actual"`
		HasHypotensionAlert   bool       `json:"has_hypotension_alert"`
		HasHypertensionAlert  bool       `json:"has_hypertension_alert"`
		HasTachycardiaAlert   bool       `json:"has_tachycardia_alert"`
		Notes                 *string    `json:"notes"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hospitalIDStr, _ := c.Get(middleware.CtxHospitalID)
	userIDStr, _ := c.Get(middleware.CtxUserID)
	hospitalID, _ := uuid.Parse(hospitalIDStr.(string))
	recordedBy, _ := uuid.Parse(userIDStr.(string))

	ctx := c.Request.Context()
	tx, err := h.pool.Begin(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}
	defer tx.Rollback(ctx)

	if err := tenant.SetLocalHospitalID(ctx, tx, hospitalID.String()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set tenant context"})
		return
	}

	queries := sqlc.New(tx)

	vital, err := queries.CreateSessionVital(ctx, sqlc.CreateSessionVitalParams{
		HospitalID:           hospitalID,
		SessionID:            req.SessionID,
		PatientID:            req.PatientID,
		RecordedBy:           recordedBy,
		RecordedAt:           pgtype.Timestamptz{Time: time.Now(), Valid: true},
		TimeOnDialysisMins:   pgtype.Int4{Int32: *req.TimeOnDialysisMins, Valid: req.TimeOnDialysisMins != nil},
		BpSystolic:           pgtype.Int4{Int32: *req.BPSystolic, Valid: req.BPSystolic != nil},
		BpDiastolic:          pgtype.Int4{Int32: *req.BPDiastolic, Valid: req.BPDiastolic != nil},
		HeartRate:            pgtype.Int4{Int32: *req.HeartRate, Valid: req.HeartRate != nil},
		Temperature:          pgtype.Numeric{Int: nil, Exp: 0, NaN: false, InfinityModifier: 0, Valid: req.Temperature != nil},
		Spo2:                 pgtype.Int4{Int32: *req.SpO2, Valid: req.SpO2 != nil},
		RespiratoryRate:      pgtype.Int4{Int32: *req.RespiratoryRate, Valid: req.RespiratoryRate != nil},
		BloodFlowActual:      pgtype.Int4{Int32: *req.BloodFlowActual, Valid: req.BloodFlowActual != nil},
		DialysateFlowActual:  pgtype.Int4{Int32: *req.DialysateFlowActual, Valid: req.DialysateFlowActual != nil},
		VenousPressure:       pgtype.Int4{Int32: *req.VenousPressure, Valid: req.VenousPressure != nil},
		ArterialPressure:     pgtype.Int4{Int32: *req.ArterialPressure, Valid: req.ArterialPressure != nil},
		Tmp:                  pgtype.Int4{Int32: *req.TMP, Valid: req.TMP != nil},
		UfRemovedSoFar:       pgtype.Numeric{Int: nil, Exp: 0, NaN: false, InfinityModifier: 0, Valid: req.UFRemovedSoFar != nil},
		ConductivityActual:   pgtype.Numeric{Int: nil, Exp: 0, NaN: false, InfinityModifier: 0, Valid: req.ConductivityActual != nil},
		HasHypotensionAlert:  req.HasHypotensionAlert,
		HasHypertensionAlert: req.HasHypertensionAlert,
		HasTachycardiaAlert:  req.HasTachycardiaAlert,
		Notes:                pgtype.Text{String: *req.Notes, Valid: req.Notes != nil},
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record vitals"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusCreated, vital)
}

func (h *VitalsHandler) ListVitalsBySession(c *gin.Context) {
	sessionID, err := uuid.Parse(c.Param("session_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
		return
	}

	hospitalIDStr, _ := c.Get(middleware.CtxHospitalID)
	hospitalID, _ := uuid.Parse(hospitalIDStr.(string))

	ctx := c.Request.Context()
	tx, err := h.pool.Begin(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}
	defer tx.Rollback(ctx)

	if err := tenant.SetLocalHospitalID(ctx, tx, hospitalID.String()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set tenant context"})
		return
	}

	queries := sqlc.New(tx)
	vitals, err := queries.ListVitalsBySession(ctx, sessionID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list vitals"})
		return
	}

	tx.Commit(ctx)
	c.JSON(http.StatusOK, vitals)
}

func (h *VitalsHandler) ListVitalsWithAlerts(c *gin.Context) {
	sessionID, err := uuid.Parse(c.Param("session_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
		return
	}

	hospitalIDStr, _ := c.Get(middleware.CtxHospitalID)
	hospitalID, _ := uuid.Parse(hospitalIDStr.(string))

	ctx := c.Request.Context()
	tx, err := h.pool.Begin(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}
	defer tx.Rollback(ctx)

	if err := tenant.SetLocalHospitalID(ctx, tx, hospitalID.String()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set tenant context"})
		return
	}

	queries := sqlc.New(tx)
	vitals, err := queries.ListVitalsWithAlerts(ctx, sessionID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list vitals with alerts"})
		return
	}

	tx.Commit(ctx)
	c.JSON(http.StatusOK, vitals)
}

func (h *VitalsHandler) AcknowledgeAlert(c *gin.Context) {
	vitalID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid vital ID"})
		return
	}

	hospitalIDStr, _ := c.Get(middleware.CtxHospitalID)
	userIDStr, _ := c.Get(middleware.CtxUserID)
	hospitalID, _ := uuid.Parse(hospitalIDStr.(string))
	acknowledgedBy, _ := uuid.Parse(userIDStr.(string))

	ctx := c.Request.Context()
	tx, err := h.pool.Begin(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}
	defer tx.Rollback(ctx)

	if err := tenant.SetLocalHospitalID(ctx, tx, hospitalID.String()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set tenant context"})
		return
	}

	queries := sqlc.New(tx)
	vital, err := queries.AcknowledgeAlert(ctx, sqlc.AcknowledgeAlertParams{
		ID:                    vitalID,
		AlertAcknowledgedBy:   pgtype.UUID{Bytes: acknowledgedBy, Valid: true},
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to acknowledge alert"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, vital)
}
