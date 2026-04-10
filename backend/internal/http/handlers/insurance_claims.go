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

type InsuranceClaimsHandler struct {
	pool *pgxpool.Pool
}

func NewInsuranceClaimsHandler(pool *pgxpool.Pool) *InsuranceClaimsHandler {
	return &InsuranceClaimsHandler{pool: pool}
}

// Create creates a new insurance claim
// POST /api/v1/insurance-claims
func (h *InsuranceClaimsHandler) Create(c *gin.Context) {
	var req struct {
		InvoiceID     string  `json:"invoice_id" binding:"required"`
		SchemeID      string  `json:"scheme_id" binding:"required"`
		PatientID     string  `json:"patient_id" binding:"required"`
		ClaimNumber   string  `json:"claim_number" binding:"required"`
		ClaimDate     string  `json:"claim_date" binding:"required"`
		ClaimedAmount float64 `json:"claimed_amount" binding:"required"`
		Status        string  `json:"status"`
		Notes         string  `json:"notes"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hospitalIDStr := c.GetString(middleware.CtxHospitalID)
	hospitalID, _ := uuid.Parse(hospitalIDStr)

	invoiceID, err := uuid.Parse(req.InvoiceID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid invoice_id"})
		return
	}

	schemeID, err := uuid.Parse(req.SchemeID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid scheme_id"})
		return
	}

	patientID, err := uuid.Parse(req.PatientID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid patient_id"})
		return
	}

	claimDate, err := time.Parse("2006-01-02", req.ClaimDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid claim_date format"})
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

	status := sqlc.ClaimStatus("draft")
	if req.Status != "" {
		status = sqlc.ClaimStatus(req.Status)
	}

	var claimedAmount pgtype.Numeric
	claimedAmount.Scan(req.ClaimedAmount)

	queries := sqlc.New(tx)
	claim, err := queries.CreateInsuranceClaim(ctx, sqlc.CreateInsuranceClaimParams{
		HospitalID:    hospitalID,
		InvoiceID:     invoiceID,
		SchemeID:      schemeID,
		PatientID:     patientID,
		ClaimNumber:   req.ClaimNumber,
		ClaimDate:     pgtype.Date{Time: claimDate, Valid: true},
		ClaimedAmount: claimedAmount,
		Status:        status,
		Notes:         pgtype.Text{String: req.Notes, Valid: req.Notes != ""},
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create insurance claim", "details": err.Error()})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusCreated, claim)
}

// Get retrieves a specific insurance claim by ID
// GET /api/v1/insurance-claims/:id
func (h *InsuranceClaimsHandler) Get(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid claim ID"})
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
	claim, err := queries.GetInsuranceClaim(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "insurance claim not found"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, claim)
}

// GetByNumber retrieves an insurance claim by claim number
// GET /api/v1/insurance-claims/number/:claim_number
func (h *InsuranceClaimsHandler) GetByNumber(c *gin.Context) {
	claimNumber := c.Param("claim_number")
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
	claim, err := queries.GetClaimByNumber(ctx, sqlc.GetClaimByNumberParams{
		HospitalID:  hospitalID,
		ClaimNumber: claimNumber,
	})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "insurance claim not found"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, claim)
}

// ListByInvoice lists insurance claims for an invoice
// GET /api/v1/invoices/:invoice_id/claims
func (h *InsuranceClaimsHandler) ListByInvoice(c *gin.Context) {
	invoiceIDStr := c.Param("invoice_id")
	invoiceID, err := uuid.Parse(invoiceIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid invoice_id"})
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
	claims, err := queries.ListClaimsByInvoice(ctx, invoiceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list insurance claims"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, claims)
}

// ListByScheme lists insurance claims for a scheme
// GET /api/v1/insurance-schemes/:scheme_id/claims
func (h *InsuranceClaimsHandler) ListByScheme(c *gin.Context) {
	schemeIDStr := c.Param("scheme_id")
	schemeID, err := uuid.Parse(schemeIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid scheme_id"})
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
	claims, err := queries.ListClaimsByScheme(ctx, schemeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list insurance claims"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, claims)
}

// ListByStatus lists insurance claims by status
// GET /api/v1/insurance-claims?status=submitted
func (h *InsuranceClaimsHandler) ListByStatus(c *gin.Context) {
	status := c.Query("status")
	if status == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "status parameter is required"})
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
	claims, err := queries.ListClaimsByStatus(ctx, sqlc.ListClaimsByStatusParams{
		HospitalID: hospitalID,
		Status:     sqlc.ClaimStatus(status),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list insurance claims"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, claims)
}

// ListPending lists pending insurance claims
// GET /api/v1/insurance-claims/pending
func (h *InsuranceClaimsHandler) ListPending(c *gin.Context) {
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
	claims, err := queries.ListPendingClaims(ctx, hospitalID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list pending claims"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, claims)
}

// Submit submits a claim to the insurance company
// POST /api/v1/insurance-claims/:id/submit
func (h *InsuranceClaimsHandler) Submit(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid claim ID"})
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

	queries := sqlc.New(tx)
	claim, err := queries.SubmitClaim(ctx, sqlc.SubmitClaimParams{
		ID:          id,
		SubmittedBy: pgtype.UUID{Bytes: userID, Valid: true},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to submit claim"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, claim)
}

// Approve approves an insurance claim
// POST /api/v1/insurance-claims/:id/approve
func (h *InsuranceClaimsHandler) Approve(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid claim ID"})
		return
	}

	var req struct {
		ApprovedAmount    float64 `json:"approved_amount" binding:"required"`
		ApprovedByInsurer string  `json:"approved_by_insurer" binding:"required"`
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

	var approvedAmount pgtype.Numeric
	approvedAmount.Scan(req.ApprovedAmount)

	queries := sqlc.New(tx)
	claim, err := queries.ApproveClaim(ctx, sqlc.ApproveClaimParams{
		ID:                id,
		ApprovedAmount:    approvedAmount,
		ApprovedByInsurer: pgtype.Text{String: req.ApprovedByInsurer, Valid: true},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to approve claim"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, claim)
}

// Reject rejects an insurance claim
// POST /api/v1/insurance-claims/:id/reject
func (h *InsuranceClaimsHandler) Reject(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid claim ID"})
		return
	}

	var req struct {
		RejectionReason string `json:"rejection_reason" binding:"required"`
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
	claim, err := queries.RejectClaim(ctx, sqlc.RejectClaimParams{
		ID:              id,
		RejectionReason: pgtype.Text{String: req.RejectionReason, Valid: true},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to reject claim"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, claim)
}
