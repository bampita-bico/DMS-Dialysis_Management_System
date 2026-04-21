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

type ConsumablesHandler struct {
	pool *pgxpool.Pool
}

func NewConsumablesHandler(pool *pgxpool.Pool) *ConsumablesHandler {
	return &ConsumablesHandler{pool: pool}
}

// Consumables catalog
func (h *ConsumablesHandler) ListConsumables(c *gin.Context) {
	hospitalID := c.GetString(middleware.CtxHospitalID)
	category := c.Query("category")

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

	var consumables []sqlc.Consumable
	if category != "" {
		consumables, err = queries.ListConsumablesByCategory(ctx, sqlc.ListConsumablesByCategoryParams{
			HospitalID: hospitalUUID,
			Category:   sqlc.ConsumableCategory(category),
		})
	} else {
		consumables, err = queries.ListActiveConsumables(ctx, hospitalUUID)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list consumables"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, consumables)
}

// Inventory management
func (h *ConsumablesHandler) GetInventoryLevels(c *gin.Context) {
	hospitalID := c.GetString(middleware.CtxHospitalID)
	consumableID := c.Param("consumable_id")

	consumableUUID, err := uuid.Parse(consumableID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid consumable ID"})
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
	inventory, err := queries.ListInventoryByConsumable(ctx, consumableUUID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get inventory"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, inventory)
}

func (h *ConsumablesHandler) ListLowStock(c *gin.Context) {
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
	inventory, err := queries.ListLowStockInventory(ctx, hospitalUUID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list low stock"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, inventory)
}

func (h *ConsumablesHandler) ListExpiringStock(c *gin.Context) {
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

	expiryDate := time.Now().AddDate(0, 0, 90)
	inventory, err := queries.ListExpiringInventory(ctx, sqlc.ListExpiringInventoryParams{
		HospitalID: hospitalUUID,
		ExpiryDate: pgtype.Date{Time: expiryDate, Valid: true},
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list expiring stock"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, inventory)
}

// Usage tracking with stock deduction
func (h *ConsumablesHandler) RecordUsage(c *gin.Context) {
	hospitalID := c.GetString(middleware.CtxHospitalID)
	staffID := c.GetString(middleware.CtxUserID)

	var req struct {
		SessionID    uuid.UUID  `json:"session_id" binding:"required"`
		ConsumableID uuid.UUID  `json:"consumable_id" binding:"required"`
		InventoryID  *uuid.UUID `json:"inventory_id"`
		QuantityUsed int32      `json:"quantity_used" binding:"required"`
		ReuseNumber  *int32     `json:"reuse_number"`
		Notes        *string    `json:"notes"`
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

	// Determine which inventory batch to deduct from
	var inventoryBatchID uuid.UUID
	if req.InventoryID != nil {
		inventoryBatchID = *req.InventoryID
	} else {
		// Auto-select the best batch (FIFO by expiry date)
		batch, err := queries.GetAvailableInventoryBatch(ctx, sqlc.GetAvailableInventoryBatchParams{
			ConsumableID: req.ConsumableID,
			HospitalID:   hospitalUUID,
		})
		if err != nil {
			if err == pgx.ErrNoRows {
				apierr.Conflict(c, apierr.ErrInsufficientStock,
					"No available inventory for this consumable. All batches are depleted.")
				return
			}
			apierr.Internal(c, "Failed to find available inventory batch")
			return
		}
		inventoryBatchID = batch.ID
	}

	// Deduct stock from the inventory batch
	updatedBatch, err := queries.DeductInventory(ctx, sqlc.DeductInventoryParams{
		ID:               inventoryBatchID,
		QuantityCurrent:  req.QuantityUsed,
	})
	if err != nil {
		// The WHERE clause includes `quantity_available >= $2`, so no rows = insufficient stock
		apierr.ConflictWithDetails(c, apierr.ErrInsufficientStock,
			"Insufficient stock in the selected batch to fulfill this usage",
			gin.H{
				"inventory_id":    inventoryBatchID,
				"quantity_needed": req.QuantityUsed,
			})
		return
	}

	// Record the usage linked to the batch
	var reuseNumber pgtype.Int4
	if req.ReuseNumber != nil {
		reuseNumber = pgtype.Int4{Int32: *req.ReuseNumber, Valid: true}
	}

	var notesText pgtype.Text
	if req.Notes != nil {
		notesText = pgtype.Text{String: *req.Notes, Valid: true}
	}

	usage, err := queries.CreateConsumablesUsage(ctx, sqlc.CreateConsumablesUsageParams{
		HospitalID:   hospitalUUID,
		SessionID:    req.SessionID,
		ConsumableID: req.ConsumableID,
		InventoryID:  pgtype.UUID{Bytes: inventoryBatchID, Valid: true},
		QuantityUsed: req.QuantityUsed,
		ReuseNumber:  reuseNumber,
		RecordedBy:   pgtype.UUID{Bytes: staffUUID, Valid: true},
		Notes:        notesText,
	})
	if err != nil {
		apierr.Internal(c, "Failed to record usage")
		return
	}

	if err := tx.Commit(ctx); err != nil {
		apierr.InternalWithCode(c, apierr.ErrTransaction, "Failed to commit transaction")
		return
	}

	// Include low-stock warning if applicable
	var warnings []apierr.Warning
	if updatedBatch.IsLowStock {
		warnings = append(warnings, apierr.Warning{
			Code:    apierr.ErrLowStock,
			Message: "Stock is now below reorder level for this batch",
			Details: gin.H{
				"inventory_id":      updatedBatch.ID,
				"quantity_remaining": updatedBatch.QuantityCurrent,
			},
		})
	}

	if len(warnings) > 0 {
		apierr.CreatedWithWarnings(c, usage, warnings)
		return
	}
	c.JSON(http.StatusCreated, usage)
}

func (h *ConsumablesHandler) ListSessionUsage(c *gin.Context) {
	hospitalID := c.GetString(middleware.CtxHospitalID)
	sessionID := c.Param("session_id")

	sessionUUID, err := uuid.Parse(sessionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
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
	usage, err := queries.ListUsageBySession(ctx, sessionUUID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list usage"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, usage)
}
