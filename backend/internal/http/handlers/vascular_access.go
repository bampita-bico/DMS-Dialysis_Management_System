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

type VascularAccessHandler struct {
	pool *pgxpool.Pool
}

func NewVascularAccessHandler(pool *pgxpool.Pool) *VascularAccessHandler {
	return &VascularAccessHandler{pool: pool}
}

// Create creates a new vascular access for a patient
// POST /api/v1/vascular-access
func (h *VascularAccessHandler) Create(c *gin.Context) {
	var req struct {
		PatientID         string  `json:"patient_id" binding:"required"`
		AccessType        string  `json:"access_type" binding:"required"`
		AccessSite        string  `json:"access_site" binding:"required"`
		SiteSide          string  `json:"site_side" binding:"required"`
		InsertionDate     string  `json:"insertion_date" binding:"required"`
		InsertedBy        string  `json:"inserted_by"`
		InsertionLocation *string `json:"insertion_location"`
		IsPrimaryAccess   bool    `json:"is_primary_access"`
		Status            *string `json:"status"`
		MaturationDate    *string `json:"maturation_date"`
		CatheterType      *string `json:"catheter_type"`
		Notes             string  `json:"notes"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hospitalIDStr := c.GetString(middleware.CtxHospitalID)
	hospitalID, _ := uuid.Parse(hospitalIDStr)
	patientID, err := uuid.Parse(req.PatientID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid patient_id"})
		return
	}

	insertionDate, err := time.Parse("2006-01-02", req.InsertionDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid insertion_date format"})
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
	var insertedBy pgtype.Text
	if req.InsertedBy != "" {
		insertedBy = pgtype.Text{String: req.InsertedBy, Valid: true}
	}

	var insertionLocation pgtype.Text
	if req.InsertionLocation != nil {
		insertionLocation = pgtype.Text{String: *req.InsertionLocation, Valid: true}
	}

	var status sqlc.AccessStatus
	if req.Status != nil {
		status = sqlc.AccessStatus(*req.Status)
	} else {
		status = "active" // default status
	}

	var maturationDate pgtype.Date
	if req.MaturationDate != nil {
		matDate, err := time.Parse("2006-01-02", *req.MaturationDate)
		if err == nil {
			maturationDate = pgtype.Date{Time: matDate, Valid: true}
		}
	}

	var catheterType pgtype.Text
	if req.CatheterType != nil {
		catheterType = pgtype.Text{String: *req.CatheterType, Valid: true}
	}

	queries := sqlc.New(tx)
	access, err := queries.CreateVascularAccess(ctx, sqlc.CreateVascularAccessParams{
		HospitalID:        hospitalID,
		PatientID:         patientID,
		AccessType:        sqlc.AccessType(req.AccessType),
		AccessSite:        sqlc.AccessSite(req.AccessSite),
		SiteSide:          req.SiteSide,
		InsertionDate:     pgtype.Date{Time: insertionDate, Valid: true},
		InsertedBy:        insertedBy,
		InsertionLocation: insertionLocation,
		Status:            status,
		MaturationDate:    maturationDate,
		FirstUseDate:      pgtype.Date{},
		CatheterType:      catheterType,
		CatheterLengthCm:  pgtype.Numeric{},
		CatheterPosition:  pgtype.Text{},
		FistulaVein:       pgtype.Text{},
		FistulaArtery:     pgtype.Text{},
		GraftMaterial:     pgtype.Text{},
		IsPrimaryAccess:   req.IsPrimaryAccess,
		Notes:             pgtype.Text{String: req.Notes, Valid: req.Notes != ""},
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create vascular access", "details": err.Error()})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusCreated, access)
}

// List lists all vascular access records for a patient
// GET /api/v1/patients/:patient_id/vascular-access
func (h *VascularAccessHandler) ListByPatient(c *gin.Context) {
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
	accesses, err := queries.ListVascularAccessByPatient(ctx, patientID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list vascular access"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, accesses)
}

// Get retrieves a specific vascular access
// GET /api/v1/vascular-access/:id
func (h *VascularAccessHandler) Get(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid vascular access ID"})
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
	access, err := queries.GetVascularAccess(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "vascular access not found"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, access)
}

// GetPrimary retrieves the primary vascular access for a patient
// GET /api/v1/patients/:patient_id/vascular-access/primary
func (h *VascularAccessHandler) GetPrimary(c *gin.Context) {
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
	access, err := queries.GetPrimaryAccessForPatient(ctx, patientID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "no primary vascular access found"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, access)
}

// Update updates vascular access information
// PATCH /api/v1/vascular-access/:id
func (h *VascularAccessHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid vascular access ID"})
		return
	}

	var req struct {
		Status          *string `json:"status"`
		IsPrimaryAccess *bool   `json:"is_primary_access"`
		MaturationDate  *string `json:"maturation_date"`
		Notes           *string `json:"notes"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

	// Get current access
	access, err := queries.GetVascularAccess(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "vascular access not found"})
		return
	}

	// Prepare update parameters
	status := access.Status
	if req.Status != nil {
		status = sqlc.AccessStatus(*req.Status)
	}

	isPrimaryAccess := access.IsPrimaryAccess
	if req.IsPrimaryAccess != nil {
		isPrimaryAccess = *req.IsPrimaryAccess
	}

	maturationDate := access.MaturationDate
	if req.MaturationDate != nil {
		matDate, err := time.Parse("2006-01-02", *req.MaturationDate)
		if err == nil {
			maturationDate = pgtype.Date{Time: matDate, Valid: true}
		}
	}

	notes := access.Notes
	if req.Notes != nil {
		notes = pgtype.Text{String: *req.Notes, Valid: *req.Notes != ""}
	}

	updatedAccess, err := queries.UpdateVascularAccess(ctx, sqlc.UpdateVascularAccessParams{
		ID:                id,
		AccessType:        access.AccessType,
		AccessSite:        access.AccessSite,
		SiteSide:          access.SiteSide,
		Status:            status,
		MaturationDate:    maturationDate,
		FirstUseDate:      access.FirstUseDate,
		AbandonmentDate:   access.AbandonmentDate,
		AbandonmentReason: access.AbandonmentReason,
		CatheterType:      access.CatheterType,
		CatheterLengthCm:  access.CatheterLengthCm,
		CatheterPosition:  access.CatheterPosition,
		FistulaVein:       access.FistulaVein,
		FistulaArtery:     access.FistulaArtery,
		GraftMaterial:     access.GraftMaterial,
		IsPrimaryAccess:   isPrimaryAccess,
		Notes:             notes,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update vascular access"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, updatedAccess)
}

// Abandon marks a vascular access as abandoned
// POST /api/v1/vascular-access/:id/abandon
func (h *VascularAccessHandler) Abandon(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid vascular access ID"})
		return
	}

	var req struct {
		AbandonReason string `json:"abandon_reason" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
	_, err = queries.AbandonAccess(ctx, sqlc.AbandonAccessParams{
		ID:                id,
		AbandonmentReason: pgtype.Text{String: req.AbandonReason, Valid: true},
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to abandon vascular access"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "vascular access abandoned successfully"})
}
