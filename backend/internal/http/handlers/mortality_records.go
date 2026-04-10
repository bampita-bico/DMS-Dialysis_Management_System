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

type MortalityRecordsHandler struct {
	pool *pgxpool.Pool
}

func NewMortalityRecordsHandler(pool *pgxpool.Pool) *MortalityRecordsHandler {
	return &MortalityRecordsHandler{pool: pool}
}

// Create creates a new mortality record
// POST /api/v1/mortality-records
func (h *MortalityRecordsHandler) Create(c *gin.Context) {
	var req struct {
		PatientID              string  `json:"patient_id" binding:"required"`
		SessionID              *string `json:"session_id"`
		DateOfDeath            string  `json:"date_of_death" binding:"required"`
		TimeOfDeath            string  `json:"time_of_death"`
		DeathSetting           string  `json:"death_setting" binding:"required"`
		SessionRelated         bool    `json:"session_related"`
		PrimaryCauseOfDeath    string  `json:"primary_cause_of_death" binding:"required"`
		ContributingFactors    string  `json:"contributing_factors"`
		Icd10Code              string  `json:"icd10_code"`
		AutopsyPerformed       bool    `json:"autopsy_performed"`
		AutopsyFindings        string  `json:"autopsy_findings"`
		CertifiedBy            *string `json:"certified_by"`
		DeathCertificateNumber string  `json:"death_certificate_number"`
		Notes                  string  `json:"notes"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hospitalIDStr := c.GetString(middleware.CtxHospitalID)
	userIDStr := c.GetString(middleware.CtxUserID)
	hospitalID, _ := uuid.Parse(hospitalIDStr)
	userID, _ := uuid.Parse(userIDStr)

	patientID, err := uuid.Parse(req.PatientID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid patient_id"})
		return
	}

	dateOfDeath, err := time.Parse("2006-01-02", req.DateOfDeath)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date_of_death format"})
		return
	}

	// Start transaction with RLS
	ctx := c.Request.Context()
	tx, err := h.pool.Begin(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to begin transaction"})
		return
	}
	defer tx.Rollback(ctx)

	if err := tenant.SetLocalHospitalID(ctx, tx, hospitalIDStr); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to set tenant context"})
		return
	}

	// Prepare parameters
	var sessionID pgtype.UUID
	if req.SessionID != nil {
		sessID, err := uuid.Parse(*req.SessionID)
		if err == nil {
			sessionID = pgtype.UUID{Bytes: sessID, Valid: true}
		}
	}

	var timeOfDeath pgtype.Time
	if req.TimeOfDeath != "" {
		t, err := time.Parse("15:04:05", req.TimeOfDeath)
		if err == nil {
			timeOfDeath = pgtype.Time{
				Microseconds: int64(t.Hour()*3600+t.Minute()*60+t.Second()) * 1000000,
				Valid:        true,
			}
		}
	}

	var certifiedBy pgtype.UUID
	if req.CertifiedBy != nil {
		certID, err := uuid.Parse(*req.CertifiedBy)
		if err == nil {
			certifiedBy = pgtype.UUID{Bytes: certID, Valid: true}
		}
	}

	queries := sqlc.New(tx)
	record, err := queries.CreateMortalityRecord(ctx, sqlc.CreateMortalityRecordParams{
		HospitalID:             hospitalID,
		PatientID:              patientID,
		SessionID:              sessionID,
		DateOfDeath:            pgtype.Date{Time: dateOfDeath, Valid: true},
		TimeOfDeath:            timeOfDeath,
		DeathSetting:           sqlc.DeathSetting(req.DeathSetting),
		SessionRelated:         req.SessionRelated,
		PrimaryCauseOfDeath:    req.PrimaryCauseOfDeath,
		ContributingFactors:    pgtype.Text{String: req.ContributingFactors, Valid: req.ContributingFactors != ""},
		Icd10Code:              pgtype.Text{String: req.Icd10Code, Valid: req.Icd10Code != ""},
		AutopsyPerformed:       req.AutopsyPerformed,
		AutopsyFindings:        pgtype.Text{String: req.AutopsyFindings, Valid: req.AutopsyFindings != ""},
		ReportedBy:             userID,
		CertifiedBy:            certifiedBy,
		DeathCertificateNumber: pgtype.Text{String: req.DeathCertificateNumber, Valid: req.DeathCertificateNumber != ""},
		Notes:                  pgtype.Text{String: req.Notes, Valid: req.Notes != ""},
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create mortality record", "details": err.Error()})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusCreated, record)
}

// Get retrieves a specific mortality record by ID
// GET /api/v1/mortality-records/:id
func (h *MortalityRecordsHandler) Get(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid mortality record ID"})
		return
	}

	hospitalIDStr := c.GetString(middleware.CtxHospitalID)
	ctx := c.Request.Context()
	tx, err := h.pool.Begin(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to begin transaction"})
		return
	}
	defer tx.Rollback(ctx)

	if err := tenant.SetLocalHospitalID(ctx, tx, hospitalIDStr); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to set tenant context"})
		return
	}

	queries := sqlc.New(tx)
	record, err := queries.GetMortalityRecord(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "mortality record not found"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, record)
}

// GetByPatient retrieves the mortality record for a patient
// GET /api/v1/patients/:patient_id/mortality-record
func (h *MortalityRecordsHandler) GetByPatient(c *gin.Context) {
	patientIDStr := c.Param("patient_id")
	patientID, err := uuid.Parse(patientIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid patient_id"})
		return
	}

	hospitalIDStr := c.GetString(middleware.CtxHospitalID)
	ctx := c.Request.Context()
	tx, err := h.pool.Begin(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to begin transaction"})
		return
	}
	defer tx.Rollback(ctx)

	if err := tenant.SetLocalHospitalID(ctx, tx, hospitalIDStr); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to set tenant context"})
		return
	}

	queries := sqlc.New(tx)
	record, err := queries.GetMortalityByPatient(ctx, patientID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "mortality record not found"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, record)
}

// List lists all mortality records for the hospital
// GET /api/v1/mortality-records
func (h *MortalityRecordsHandler) List(c *gin.Context) {
	hospitalIDStr := c.GetString(middleware.CtxHospitalID)
	hospitalID, _ := uuid.Parse(hospitalIDStr)

	ctx := c.Request.Context()
	tx, err := h.pool.Begin(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to begin transaction"})
		return
	}
	defer tx.Rollback(ctx)

	if err := tenant.SetLocalHospitalID(ctx, tx, hospitalIDStr); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to set tenant context"})
		return
	}

	queries := sqlc.New(tx)
	records, err := queries.ListMortalitiesByHospital(ctx, hospitalID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list mortality records"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, records)
}

// ListByPeriod lists mortality records within a date range
// GET /api/v1/mortality-records/period?start_date=2024-01-01&end_date=2024-12-31
func (h *MortalityRecordsHandler) ListByPeriod(c *gin.Context) {
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if startDateStr == "" || endDateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start_date and end_date are required"})
		return
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_date format"})
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end_date format"})
		return
	}

	hospitalIDStr := c.GetString(middleware.CtxHospitalID)
	hospitalID, _ := uuid.Parse(hospitalIDStr)

	ctx := c.Request.Context()
	tx, err := h.pool.Begin(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to begin transaction"})
		return
	}
	defer tx.Rollback(ctx)

	if err := tenant.SetLocalHospitalID(ctx, tx, hospitalIDStr); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to set tenant context"})
		return
	}

	queries := sqlc.New(tx)
	records, err := queries.ListMortalitiesByPeriod(ctx, sqlc.ListMortalitiesByPeriodParams{
		HospitalID:  hospitalID,
		DateOfDeath: pgtype.Date{Time: startDate, Valid: true},
		DateOfDeath_2: pgtype.Date{Time: endDate, Valid: true},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list mortality records"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, records)
}

// ListSessionRelated lists all session-related deaths
// GET /api/v1/mortality-records/session-related
func (h *MortalityRecordsHandler) ListSessionRelated(c *gin.Context) {
	hospitalIDStr := c.GetString(middleware.CtxHospitalID)
	hospitalID, _ := uuid.Parse(hospitalIDStr)

	ctx := c.Request.Context()
	tx, err := h.pool.Begin(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to begin transaction"})
		return
	}
	defer tx.Rollback(ctx)

	if err := tenant.SetLocalHospitalID(ctx, tx, hospitalIDStr); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to set tenant context"})
		return
	}

	queries := sqlc.New(tx)
	records, err := queries.ListSessionRelatedDeaths(ctx, hospitalID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list session-related deaths"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, records)
}

// ListBySetting lists mortality records filtered by death setting
// GET /api/v1/mortality-records/setting/:setting
func (h *MortalityRecordsHandler) ListBySetting(c *gin.Context) {
	setting := c.Param("setting")
	if setting == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "death setting is required"})
		return
	}

	hospitalIDStr := c.GetString(middleware.CtxHospitalID)
	hospitalID, _ := uuid.Parse(hospitalIDStr)

	ctx := c.Request.Context()
	tx, err := h.pool.Begin(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to begin transaction"})
		return
	}
	defer tx.Rollback(ctx)

	if err := tenant.SetLocalHospitalID(ctx, tx, hospitalIDStr); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to set tenant context"})
		return
	}

	queries := sqlc.New(tx)
	records, err := queries.ListDeathsBySetting(ctx, sqlc.ListDeathsBySettingParams{
		HospitalID:   hospitalID,
		DeathSetting: sqlc.DeathSetting(setting),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list deaths by setting"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, records)
}

// Certify certifies a death and adds certificate number
// POST /api/v1/mortality-records/:id/certify
func (h *MortalityRecordsHandler) Certify(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid mortality record ID"})
		return
	}

	var req struct {
		CertifiedBy            *string `json:"certified_by"`
		DeathCertificateNumber string  `json:"death_certificate_number" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hospitalIDStr := c.GetString(middleware.CtxHospitalID)
	userIDStr := c.GetString(middleware.CtxUserID)
	userID, _ := uuid.Parse(userIDStr)

	ctx := c.Request.Context()
	tx, err := h.pool.Begin(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to begin transaction"})
		return
	}
	defer tx.Rollback(ctx)

	if err := tenant.SetLocalHospitalID(ctx, tx, hospitalIDStr); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to set tenant context"})
		return
	}

	// Use provided certified_by or default to current user
	var certifiedBy pgtype.UUID
	if req.CertifiedBy != nil {
		certID, err := uuid.Parse(*req.CertifiedBy)
		if err == nil {
			certifiedBy = pgtype.UUID{Bytes: certID, Valid: true}
		}
	} else {
		certifiedBy = pgtype.UUID{Bytes: userID, Valid: true}
	}

	queries := sqlc.New(tx)
	record, err := queries.CertifyDeath(ctx, sqlc.CertifyDeathParams{
		ID:                     id,
		CertifiedBy:            certifiedBy,
		DeathCertificateNumber: pgtype.Text{String: req.DeathCertificateNumber, Valid: true},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to certify death"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, record)
}
