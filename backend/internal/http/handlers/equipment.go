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

type EquipmentHandler struct {
	pool *pgxpool.Pool
}

func NewEquipmentHandler(pool *pgxpool.Pool) *EquipmentHandler {
	return &EquipmentHandler{pool: pool}
}

// Equipment CRUD
func (h *EquipmentHandler) CreateEquipment(c *gin.Context) {
	hospitalID := c.GetString(middleware.CtxHospitalID)

	var req struct {
		Name               string     `json:"name" binding:"required"`
		Category           string     `json:"category" binding:"required"`
		SerialNumber       *string    `json:"serial_number"`
		Model              *string    `json:"model"`
		Manufacturer       *string    `json:"manufacturer"`
		PurchaseDate       *time.Time `json:"purchase_date"`
		PurchaseCost       *float64   `json:"purchase_cost"`
		WarrantyExpiryDate *time.Time `json:"warranty_expiry_date"`
		Status             string     `json:"status" binding:"required"`
		Location           *string    `json:"location"`
		DepartmentID       *uuid.UUID `json:"department_id"`
		Notes              *string    `json:"notes"`
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
	hospitalUUID := uuid.MustParse(hospitalID)

	var purchaseDate, warrantyDate pgtype.Date
	if req.PurchaseDate != nil {
		purchaseDate = pgtype.Date{Time: *req.PurchaseDate, Valid: true}
	}
	if req.WarrantyExpiryDate != nil {
		warrantyDate = pgtype.Date{Time: *req.WarrantyExpiryDate, Valid: true}
	}

	var deptID pgtype.UUID
	if req.DepartmentID != nil {
		deptID = pgtype.UUID{Bytes: *req.DepartmentID, Valid: true}
	}

	var purchaseCost pgtype.Numeric
	if req.PurchaseCost != nil {
		purchaseCost = pgtype.Numeric{Valid: true}
		purchaseCost.Scan(*req.PurchaseCost)
	}

	equipment, err := queries.CreateEquipment(ctx, sqlc.CreateEquipmentParams{
		HospitalID:         hospitalUUID,
		Name:               req.Name,
		Category:           sqlc.EquipmentCategory(req.Category),
		SerialNumber:       pgtype.Text{String: *req.SerialNumber, Valid: req.SerialNumber != nil},
		Model:              pgtype.Text{String: *req.Model, Valid: req.Model != nil},
		Manufacturer:       pgtype.Text{String: *req.Manufacturer, Valid: req.Manufacturer != nil},
		PurchaseDate:       purchaseDate,
		PurchaseCost:       purchaseCost,
		WarrantyExpiryDate: warrantyDate,
		Status:             sqlc.EquipmentStatus(req.Status),
		Location:           pgtype.Text{String: *req.Location, Valid: req.Location != nil},
		DepartmentID:       deptID,
		Notes:              pgtype.Text{String: *req.Notes, Valid: req.Notes != nil},
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create equipment"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusCreated, equipment)
}

func (h *EquipmentHandler) ListEquipment(c *gin.Context) {
	hospitalID := c.GetString(middleware.CtxHospitalID)
	category := c.Query("category")
	status := c.Query("status")

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
	hospitalUUID := uuid.MustParse(hospitalID)

	var equipment []sqlc.Equipment

	if category != "" {
		equipment, err = queries.ListEquipmentByCategory(ctx, sqlc.ListEquipmentByCategoryParams{
			HospitalID: hospitalUUID,
			Category:   sqlc.EquipmentCategory(category),
		})
	} else if status != "" {
		equipment, err = queries.ListEquipmentByStatus(ctx, sqlc.ListEquipmentByStatusParams{
			HospitalID: hospitalUUID,
			Status:     sqlc.EquipmentStatus(status),
		})
	} else {
		equipment, err = queries.ListEquipmentByHospital(ctx, hospitalUUID)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list equipment"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, equipment)
}

func (h *EquipmentHandler) UpdateEquipmentStatus(c *gin.Context) {
	hospitalID := c.GetString(middleware.CtxHospitalID)
	equipmentID := c.Param("id")

	var req struct {
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	equipmentUUID, err := uuid.Parse(equipmentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid equipment ID"})
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
	equipment, err := queries.UpdateEquipmentStatus(ctx, sqlc.UpdateEquipmentStatusParams{
		ID:     equipmentUUID,
		Status: sqlc.EquipmentStatus(req.Status),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update equipment status"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, equipment)
}

// Equipment Faults
func (h *EquipmentHandler) ReportFault(c *gin.Context) {
	hospitalID := c.GetString(middleware.CtxHospitalID)
	staffID := c.GetString(middleware.CtxUserID)

	var req struct {
		EquipmentID         uuid.UUID `json:"equipment_id" binding:"required"`
		FaultDescription    string    `json:"fault_description" binding:"required"`
		Severity            string    `json:"severity" binding:"required"`
		IsEquipmentUnusable bool      `json:"is_equipment_unusable"`
		Notes               *string   `json:"notes"`
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
	hospitalUUID := uuid.MustParse(hospitalID)
	staffUUID := uuid.MustParse(staffID)

	fault, err := queries.CreateEquipmentFault(ctx, sqlc.CreateEquipmentFaultParams{
		HospitalID:          hospitalUUID,
		EquipmentID:         req.EquipmentID,
		ReportedBy:          staffUUID,
		FaultDescription:    req.FaultDescription,
		Severity:            sqlc.FaultSeverity(req.Severity),
		IsEquipmentUnusable: req.IsEquipmentUnusable,
		Notes:               pgtype.Text{String: *req.Notes, Valid: req.Notes != nil},
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to report fault"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusCreated, fault)
}

func (h *EquipmentHandler) ListUnresolvedFaults(c *gin.Context) {
	hospitalID := c.GetString(middleware.CtxHospitalID)

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
	hospitalUUID := uuid.MustParse(hospitalID)
	faults, err := queries.ListUnresolvedFaults(ctx, hospitalUUID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list faults"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, faults)
}
