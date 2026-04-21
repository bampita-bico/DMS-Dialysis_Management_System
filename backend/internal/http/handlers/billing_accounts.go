package handlers

import (
	"net/http"

	"github.com/dmsafrica/dms/internal/db/sqlc"
	"github.com/dmsafrica/dms/internal/db/tenant"
	"github.com/dmsafrica/dms/internal/http/middleware"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BillingAccountsHandler struct {
	pool *pgxpool.Pool
}

func NewBillingAccountsHandler(pool *pgxpool.Pool) *BillingAccountsHandler {
	return &BillingAccountsHandler{pool: pool}
}

// Create creates a new billing account
// POST /api/v1/billing-accounts
func (h *BillingAccountsHandler) Create(c *gin.Context) {
	var req struct {
		PatientID     string   `json:"patient_id" binding:"required"`
		GuarantorID   *string  `json:"guarantor_id"`
		AccountNumber string   `json:"account_number" binding:"required"`
		AccountStatus string   `json:"account_status"`
		CreditLimit   *float64 `json:"credit_limit"`
		Notes         string   `json:"notes"`
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
	var guarantorID pgtype.UUID
	if req.GuarantorID != nil {
		gID, err := uuid.Parse(*req.GuarantorID)
		if err == nil {
			guarantorID = pgtype.UUID{Bytes: gID, Valid: true}
		}
	}

	accountStatus := sqlc.AccountStatus("active")
	if req.AccountStatus != "" {
		accountStatus = sqlc.AccountStatus(req.AccountStatus)
	}

	var creditLimit pgtype.Numeric
	if req.CreditLimit != nil {
		creditLimit.Scan(*req.CreditLimit)
	}

	queries := sqlc.New(tx)
	account, err := queries.CreateBillingAccount(ctx, sqlc.CreateBillingAccountParams{
		HospitalID:    hospitalID,
		PatientID:     patientID,
		GuarantorID:   guarantorID,
		AccountNumber: req.AccountNumber,
		AccountStatus: accountStatus,
		CreditLimit:   creditLimit,
		Notes:         pgtype.Text{String: req.Notes, Valid: req.Notes != ""},
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create billing account"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusCreated, account)
}

// Get retrieves a specific billing account by ID
// GET /api/v1/billing-accounts/:id
func (h *BillingAccountsHandler) Get(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid billing account ID"})
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
	account, err := queries.GetBillingAccount(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "billing account not found"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, account)
}

// GetByPatient retrieves the billing account for a patient
// GET /api/v1/patients/:patient_id/billing-account
func (h *BillingAccountsHandler) GetByPatient(c *gin.Context) {
	patientIDStr := c.Param("patient_id")
	patientID, err := uuid.Parse(patientIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid patient_id"})
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
	account, err := queries.GetBillingAccountByPatient(ctx, sqlc.GetBillingAccountByPatientParams{
		HospitalID: hospitalID,
		PatientID:  patientID,
	})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "billing account not found"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, account)
}

// List lists all billing accounts
// GET /api/v1/billing-accounts
func (h *BillingAccountsHandler) List(c *gin.Context) {
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
	accounts, err := queries.ListBillingAccountsByHospital(ctx, hospitalID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list billing accounts"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, accounts)
}

// ListByStatus lists billing accounts by status
// GET /api/v1/billing-accounts?status=suspended
func (h *BillingAccountsHandler) ListByStatus(c *gin.Context) {
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
	accounts, err := queries.ListAccountsByStatus(ctx, sqlc.ListAccountsByStatusParams{
		HospitalID:    hospitalID,
		AccountStatus: sqlc.AccountStatus(status),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list billing accounts"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, accounts)
}

// UpdateBalance updates the balance fields of a billing account
// PATCH /api/v1/billing-accounts/:id/balance
func (h *BillingAccountsHandler) UpdateBalance(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid billing account ID"})
		return
	}

	var req struct {
		CurrentBalance float64 `json:"current_balance" binding:"required"`
		TotalBilled    float64 `json:"total_billed" binding:"required"`
		TotalPaid      float64 `json:"total_paid" binding:"required"`
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

	var currentBalance, totalBilled, totalPaid pgtype.Numeric
	currentBalance.Scan(req.CurrentBalance)
	totalBilled.Scan(req.TotalBilled)
	totalPaid.Scan(req.TotalPaid)

	queries := sqlc.New(tx)
	account, err := queries.UpdateAccountBalance(ctx, sqlc.UpdateAccountBalanceParams{
		ID:             id,
		CurrentBalance: currentBalance,
		TotalBilled:    totalBilled,
		TotalPaid:      totalPaid,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update billing account balance"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, account)
}

// UpdateStatus updates the status of a billing account
// PATCH /api/v1/billing-accounts/:id/status
func (h *BillingAccountsHandler) UpdateStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid billing account ID"})
		return
	}

	var req struct {
		AccountStatus string `json:"account_status" binding:"required"`
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
	account, err := queries.UpdateAccountStatus(ctx, sqlc.UpdateAccountStatusParams{
		ID:            id,
		AccountStatus: sqlc.AccountStatus(req.AccountStatus),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update billing account status"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, account)
}
