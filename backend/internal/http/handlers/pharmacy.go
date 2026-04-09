package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/dmsafrica/dms/internal/db/sqlc"
	"github.com/dmsafrica/dms/internal/db/tenant"
	"github.com/dmsafrica/dms/internal/http/middleware"
)

type PharmacyHandler struct {
	pool *pgxpool.Pool
}

func NewPharmacyHandler(pool *pgxpool.Pool) *PharmacyHandler {
	return &PharmacyHandler{pool: pool}
}

// Medication catalog endpoints
func (h *PharmacyHandler) ListMedications(c *gin.Context) {
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
	medications, err := queries.ListActiveMedications(ctx, hospitalUUID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list medications"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, medications)
}

func (h *PharmacyHandler) SearchMedications(c *gin.Context) {
	hospitalID := c.GetString(middleware.CtxHospitalID)
	searchTerm := c.Query("q")

	if searchTerm == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Search term required"})
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
	medications, err := queries.SearchMedications(ctx, sqlc.SearchMedicationsParams{
		HospitalID: hospitalUUID,
		Column2:    pgtype.Text{String: searchTerm, Valid: true},
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search medications"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, medications)
}

func (h *PharmacyHandler) GetMedication(c *gin.Context) {
	hospitalID := c.GetString(middleware.CtxHospitalID)
	medicationID := c.Param("id")

	medicationUUID, err := uuid.Parse(medicationID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid medication ID"})
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
	medication, err := queries.GetMedication(ctx, medicationUUID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Medication not found"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, medication)
}

// Stock management endpoints
func (h *PharmacyHandler) GetStockLevels(c *gin.Context) {
	hospitalID := c.GetString(middleware.CtxHospitalID)
	medicationID := c.Param("medication_id")

	medicationUUID, err := uuid.Parse(medicationID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid medication ID"})
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
	stock, err := queries.ListStockByMedication(ctx, medicationUUID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get stock levels"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, stock)
}

func (h *PharmacyHandler) ListLowStock(c *gin.Context) {
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
	stock, err := queries.ListLowStock(ctx, hospitalUUID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list low stock"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, stock)
}

// Drug interactions endpoint
func (h *PharmacyHandler) CheckDrugInteraction(c *gin.Context) {
	hospitalID := c.GetString(middleware.CtxHospitalID)

	var req struct {
		MedicationAID string `json:"medication_a_id" binding:"required,uuid"`
		MedicationBID string `json:"medication_b_id" binding:"required,uuid"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	medAUUID, _ := uuid.Parse(req.MedicationAID)
	medBUUID, _ := uuid.Parse(req.MedicationBID)

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
	interaction, err := queries.CheckInteraction(ctx, sqlc.CheckInteractionParams{
		MedicationAID: medAUUID,
		MedicationBID: medBUUID,
	})

	if err != nil {
		// No interaction found
		c.JSON(http.StatusOK, gin.H{"has_interaction": false})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"has_interaction": true,
		"interaction":     interaction,
	})
}

func (h *PharmacyHandler) ListDrugInteractions(c *gin.Context) {
	hospitalID := c.GetString(middleware.CtxHospitalID)
	medicationID := c.Param("medication_id")

	medicationUUID, err := uuid.Parse(medicationID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid medication ID"})
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
	interactions, err := queries.ListInteractionsForMedication(ctx, medicationUUID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list interactions"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, interactions)
}
