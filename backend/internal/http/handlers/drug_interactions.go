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

type DrugInteractionsHandler struct {
	pool *pgxpool.Pool
}

func NewDrugInteractionsHandler(pool *pgxpool.Pool) *DrugInteractionsHandler {
	return &DrugInteractionsHandler{pool: pool}
}

// Create creates a new drug interaction record
// POST /api/v1/drug-interactions
func (h *DrugInteractionsHandler) Create(c *gin.Context) {
	var req struct {
		MedicationAID            string `json:"medication_a_id" binding:"required"`
		MedicationBID            string `json:"medication_b_id" binding:"required"`
		Severity                 string `json:"severity" binding:"required"`
		Description              string `json:"description" binding:"required"`
		ClinicalEffect           string `json:"clinical_effect"`
		ManagementRecommendation string `json:"management_recommendation"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hospitalID := c.GetString(middleware.CtxHospitalID)
	ctx := c.Request.Context()

	// Parse UUIDs
	medAID, err := uuid.Parse(req.MedicationAID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid medication_a_id format"})
		return
	}

	medBID, err := uuid.Parse(req.MedicationBID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid medication_b_id format"})
		return
	}

	// Medications should be different
	if medAID == medBID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "medication_a_id and medication_b_id must be different"})
		return
	}

	// Handle optional fields
	var clinicalEffect pgtype.Text
	if req.ClinicalEffect != "" {
		clinicalEffect = pgtype.Text{String: req.ClinicalEffect, Valid: true}
	}

	var managementRecommendation pgtype.Text
	if req.ManagementRecommendation != "" {
		managementRecommendation = pgtype.Text{String: req.ManagementRecommendation, Valid: true}
	}

	// Begin transaction
	tx, err := h.pool.Begin(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to begin transaction"})
		return
	}
	defer tx.Rollback(ctx)

	// Set tenant context
	if err := tenant.SetLocalHospitalID(ctx, tx, hospitalID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to set hospital context"})
		return
	}

	queries := sqlc.New(tx)

	params := sqlc.CreateDrugInteractionParams{
		HospitalID:               uuid.MustParse(hospitalID),
		MedicationAID:            medAID,
		MedicationBID:            medBID,
		Severity:                 req.Severity,
		Description:              req.Description,
		ClinicalEffect:           clinicalEffect,
		ManagementRecommendation: managementRecommendation,
	}

	interaction, err := queries.CreateDrugInteraction(ctx, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create drug interaction"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusCreated, interaction)
}

// Get retrieves a drug interaction by ID
// GET /api/v1/drug-interactions/:id
func (h *DrugInteractionsHandler) Get(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid interaction ID format"})
		return
	}

	hospitalID := c.GetString(middleware.CtxHospitalID)
	ctx := c.Request.Context()

	tx, err := h.pool.Begin(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to begin transaction"})
		return
	}
	defer tx.Rollback(ctx)

	if err := tenant.SetLocalHospitalID(ctx, tx, hospitalID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to set hospital context"})
		return
	}

	queries := sqlc.New(tx)
	interaction, err := queries.GetDrugInteraction(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "drug interaction not found"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, interaction)
}

// CheckInteraction checks if two medications interact
// POST /api/v1/drug-interactions/check
func (h *DrugInteractionsHandler) CheckInteraction(c *gin.Context) {
	var req struct {
		MedicationAID string `json:"medication_a_id" binding:"required"`
		MedicationBID string `json:"medication_b_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hospitalID := c.GetString(middleware.CtxHospitalID)
	ctx := c.Request.Context()

	// Parse UUIDs
	medAID, err := uuid.Parse(req.MedicationAID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid medication_a_id format"})
		return
	}

	medBID, err := uuid.Parse(req.MedicationBID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid medication_b_id format"})
		return
	}

	// Begin transaction
	tx, err := h.pool.Begin(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to begin transaction"})
		return
	}
	defer tx.Rollback(ctx)

	if err := tenant.SetLocalHospitalID(ctx, tx, hospitalID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to set hospital context"})
		return
	}

	queries := sqlc.New(tx)
	interaction, err := queries.CheckInteraction(ctx, sqlc.CheckInteractionParams{
		MedicationAID: medAID,
		MedicationBID: medBID,
	})

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	// If no interaction found, return success with no interaction
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"has_interaction": false,
			"interaction":     nil,
		})
		return
	}

	// Interaction found
	c.JSON(http.StatusOK, gin.H{
		"has_interaction": true,
		"interaction":     interaction,
	})
}

// ListForMedication lists all interactions for a specific medication
// GET /api/v1/medications/:medication_id/interactions
func (h *DrugInteractionsHandler) ListForMedication(c *gin.Context) {
	medicationIDStr := c.Param("medication_id")
	medicationID, err := uuid.Parse(medicationIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid medication ID format"})
		return
	}

	hospitalID := c.GetString(middleware.CtxHospitalID)
	ctx := c.Request.Context()

	tx, err := h.pool.Begin(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to begin transaction"})
		return
	}
	defer tx.Rollback(ctx)

	if err := tenant.SetLocalHospitalID(ctx, tx, hospitalID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to set hospital context"})
		return
	}

	queries := sqlc.New(tx)
	interactions, err := queries.ListInteractionsForMedication(ctx, medicationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list interactions"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, interactions)
}

// ListSevere lists all severe/contraindicated interactions
// GET /api/v1/drug-interactions/severe
func (h *DrugInteractionsHandler) ListSevere(c *gin.Context) {
	hospitalID := c.GetString(middleware.CtxHospitalID)
	ctx := c.Request.Context()

	tx, err := h.pool.Begin(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to begin transaction"})
		return
	}
	defer tx.Rollback(ctx)

	if err := tenant.SetLocalHospitalID(ctx, tx, hospitalID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to set hospital context"})
		return
	}

	queries := sqlc.New(tx)
	interactions, err := queries.ListSevereInteractions(ctx, uuid.MustParse(hospitalID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list severe interactions"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, interactions)
}

// Delete soft deletes a drug interaction
// DELETE /api/v1/drug-interactions/:id
func (h *DrugInteractionsHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid interaction ID format"})
		return
	}

	hospitalID := c.GetString(middleware.CtxHospitalID)
	ctx := c.Request.Context()

	tx, err := h.pool.Begin(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to begin transaction"})
		return
	}
	defer tx.Rollback(ctx)

	if err := tenant.SetLocalHospitalID(ctx, tx, hospitalID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to set hospital context"})
		return
	}

	queries := sqlc.New(tx)
	err = queries.DeleteDrugInteraction(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete drug interaction"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "drug interaction deleted successfully"})
}
