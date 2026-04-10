package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/dmsafrica/dms/internal/db/sqlc"
	"github.com/dmsafrica/dms/internal/db/tenant"
	"github.com/dmsafrica/dms/internal/http/middleware"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
)

type MedicationsHandler struct {
	pool *pgxpool.Pool
}

func NewMedicationsHandler(pool *pgxpool.Pool) *MedicationsHandler {
	return &MedicationsHandler{pool: pool}
}

// Create creates a new medication in the catalog
// POST /api/v1/medications
func (h *MedicationsHandler) Create(c *gin.Context) {
	var req struct {
		GenericName          string              `json:"generic_name" binding:"required"`
		BrandNames           []string            `json:"brand_names"`
		DrugClass            string              `json:"drug_class"`
		Form                 sqlc.MedicationForm `json:"form" binding:"required"`
		Strength             string              `json:"strength"`
		Unit                 string              `json:"unit"`
		IsControlled         bool                `json:"is_controlled"`
		RequiresPrescription bool                `json:"requires_prescription"`
		IsEssentialWho       bool                `json:"is_essential_who"`
		StorageConditions    string              `json:"storage_conditions"`
		CostPerUnit          *float64            `json:"cost_per_unit"`
		ReorderLevel         *int32              `json:"reorder_level"`
		Notes                string              `json:"notes"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hospitalID := c.GetString(middleware.CtxHospitalID)
	ctx := c.Request.Context()

	// Marshal brand names to JSONB
	brandNamesJSON, err := json.Marshal(req.BrandNames)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid brand_names format"})
		return
	}

	// Handle optional fields
	var drugClass pgtype.Text
	if req.DrugClass != "" {
		drugClass = pgtype.Text{String: req.DrugClass, Valid: true}
	}

	var strength pgtype.Text
	if req.Strength != "" {
		strength = pgtype.Text{String: req.Strength, Valid: true}
	}

	var unit pgtype.Text
	if req.Unit != "" {
		unit = pgtype.Text{String: req.Unit, Valid: true}
	}

	var storageConditions pgtype.Text
	if req.StorageConditions != "" {
		storageConditions = pgtype.Text{String: req.StorageConditions, Valid: true}
	}

	var costPerUnit pgtype.Numeric
	if req.CostPerUnit != nil {
		dec := decimal.NewFromFloat(*req.CostPerUnit)
		costPerUnit = pgtype.Numeric{
			Int:              dec.BigInt(),
			Exp:              dec.Exponent(),
			NaN:              false,
			InfinityModifier: 0,
			Valid:            true,
		}
	}

	var reorderLevel pgtype.Int4
	if req.ReorderLevel != nil {
		reorderLevel = pgtype.Int4{Int32: *req.ReorderLevel, Valid: true}
	}

	var notes pgtype.Text
	if req.Notes != "" {
		notes = pgtype.Text{String: req.Notes, Valid: true}
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

	params := sqlc.CreateMedicationParams{
		HospitalID:           uuid.MustParse(hospitalID),
		GenericName:          req.GenericName,
		BrandNames:           brandNamesJSON,
		DrugClass:            drugClass,
		Form:                 req.Form,
		Strength:             strength,
		Unit:                 unit,
		IsControlled:         req.IsControlled,
		RequiresPrescription: req.RequiresPrescription,
		IsEssentialWho:       req.IsEssentialWho,
		StorageConditions:    storageConditions,
		CostPerUnit:          costPerUnit,
		ReorderLevel:         reorderLevel,
		Notes:                notes,
	}

	medication, err := queries.CreateMedication(ctx, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create medication"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusCreated, medication)
}

// Get retrieves a medication by ID
// GET /api/v1/medications/:id
func (h *MedicationsHandler) Get(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
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
	medication, err := queries.GetMedication(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "medication not found"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, medication)
}

// List lists all medications for the hospital
// GET /api/v1/medications
func (h *MedicationsHandler) List(c *gin.Context) {
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
	medications, err := queries.ListMedicationsByHospital(ctx, uuid.MustParse(hospitalID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list medications"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, medications)
}

// ListActive lists all active medications
// GET /api/v1/medications/active
func (h *MedicationsHandler) ListActive(c *gin.Context) {
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
	medications, err := queries.ListActiveMedications(ctx, uuid.MustParse(hospitalID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list active medications"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, medications)
}

// ListByClass lists medications by drug class
// GET /api/v1/medications/class/:class
func (h *MedicationsHandler) ListByClass(c *gin.Context) {
	drugClass := c.Param("class")
	if drugClass == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "drug class is required"})
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
	medications, err := queries.ListMedicationsByClass(ctx, sqlc.ListMedicationsByClassParams{
		HospitalID: uuid.MustParse(hospitalID),
		DrugClass:  pgtype.Text{String: drugClass, Valid: true},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list medications by class"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, medications)
}

// Search searches medications by name or class
// GET /api/v1/medications/search?query=aspirin
func (h *MedicationsHandler) Search(c *gin.Context) {
	query := c.Query("query")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query parameter is required"})
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
	medications, err := queries.SearchMedications(ctx, sqlc.SearchMedicationsParams{
		HospitalID: uuid.MustParse(hospitalID),
		Column2:    pgtype.Text{String: query, Valid: true},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to search medications"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, medications)
}

// Update updates a medication
// PATCH /api/v1/medications/:id
func (h *MedicationsHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid medication ID format"})
		return
	}

	var req struct {
		GenericName          string              `json:"generic_name" binding:"required"`
		BrandNames           []string            `json:"brand_names"`
		DrugClass            string              `json:"drug_class"`
		Form                 sqlc.MedicationForm `json:"form" binding:"required"`
		Strength             string              `json:"strength"`
		Unit                 string              `json:"unit"`
		IsControlled         bool                `json:"is_controlled"`
		RequiresPrescription bool                `json:"requires_prescription"`
		IsEssentialWho       bool                `json:"is_essential_who"`
		StorageConditions    string              `json:"storage_conditions"`
		CostPerUnit          *float64            `json:"cost_per_unit"`
		ReorderLevel         *int32              `json:"reorder_level"`
		IsActive             bool                `json:"is_active"`
		Notes                string              `json:"notes"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hospitalID := c.GetString(middleware.CtxHospitalID)
	ctx := c.Request.Context()

	// Marshal brand names to JSONB
	brandNamesJSON, err := json.Marshal(req.BrandNames)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid brand_names format"})
		return
	}

	// Handle optional fields
	var drugClass pgtype.Text
	if req.DrugClass != "" {
		drugClass = pgtype.Text{String: req.DrugClass, Valid: true}
	}

	var strength pgtype.Text
	if req.Strength != "" {
		strength = pgtype.Text{String: req.Strength, Valid: true}
	}

	var unit pgtype.Text
	if req.Unit != "" {
		unit = pgtype.Text{String: req.Unit, Valid: true}
	}

	var storageConditions pgtype.Text
	if req.StorageConditions != "" {
		storageConditions = pgtype.Text{String: req.StorageConditions, Valid: true}
	}

	var costPerUnit pgtype.Numeric
	if req.CostPerUnit != nil {
		dec := decimal.NewFromFloat(*req.CostPerUnit)
		costPerUnit = pgtype.Numeric{
			Int:              dec.BigInt(),
			Exp:              dec.Exponent(),
			NaN:              false,
			InfinityModifier: 0,
			Valid:            true,
		}
	}

	var reorderLevel pgtype.Int4
	if req.ReorderLevel != nil {
		reorderLevel = pgtype.Int4{Int32: *req.ReorderLevel, Valid: true}
	}

	var notes pgtype.Text
	if req.Notes != "" {
		notes = pgtype.Text{String: req.Notes, Valid: true}
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

	params := sqlc.UpdateMedicationParams{
		ID:                   id,
		GenericName:          req.GenericName,
		BrandNames:           brandNamesJSON,
		DrugClass:            drugClass,
		Form:                 req.Form,
		Strength:             strength,
		Unit:                 unit,
		IsControlled:         req.IsControlled,
		RequiresPrescription: req.RequiresPrescription,
		IsEssentialWho:       req.IsEssentialWho,
		StorageConditions:    storageConditions,
		CostPerUnit:          costPerUnit,
		ReorderLevel:         reorderLevel,
		IsActive:             req.IsActive,
		Notes:                notes,
	}

	medication, err := queries.UpdateMedication(ctx, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update medication"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, medication)
}

// Delete soft deletes a medication
// DELETE /api/v1/medications/:id
func (h *MedicationsHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
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
	err = queries.DeleteMedication(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete medication"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "medication deleted successfully"})
}
