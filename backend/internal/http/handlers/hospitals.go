package handlers

import (
	"net/http"

	"github.com/dmsafrica/dms/internal/db/sqlc"
	"github.com/dmsafrica/dms/internal/http/middleware"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type HospitalsHandler struct {
	pool *pgxpool.Pool
}

func NewHospitalsHandler(pool *pgxpool.Pool) *HospitalsHandler {
	return &HospitalsHandler{pool: pool}
}

// CreateHospital godoc
// @Summary Create a new hospital
// @Tags hospitals
// @Accept json
// @Produce json
// @Param hospital body CreateHospitalRequest true "Hospital data"
// @Success 201 {object} sqlc.Hospital
// @Router /api/v1/hospitals [post]
func (h *HospitalsHandler) Create(c *gin.Context) {
	var req struct {
		Name      string `json:"name" binding:"required"`
		ShortCode string `json:"short_code" binding:"required"`
		Tier      string `json:"tier" binding:"required"`
		Region    string `json:"region" binding:"required"`
		Country   string `json:"country"`
		Address   string `json:"address"`
		Phone     string `json:"phone"`
		Email     string `json:"email"`
		LicenseNo string `json:"license_no"`
		Settings  []byte `json:"settings"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	queries := sqlc.New(h.pool)

	hospital, err := queries.CreateHospital(c.Request.Context(), sqlc.CreateHospitalParams{
		Name:      req.Name,
		ShortCode: req.ShortCode,
		Tier:      req.Tier,
		Region:    req.Region,
		Country:   req.Country,
		Address:   pgtype.Text{String: req.Address, Valid: req.Address != ""},
		Phone:     pgtype.Text{String: req.Phone, Valid: req.Phone != ""},
		Email:     pgtype.Text{String: req.Email, Valid: req.Email != ""},
		LicenseNo: pgtype.Text{String: req.LicenseNo, Valid: req.LicenseNo != ""},
		Settings:  req.Settings,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create hospital"})
		return
	}

	c.JSON(http.StatusCreated, hospital)
}

// GetHospital godoc
// @Summary Get a hospital by ID
// @Tags hospitals
// @Produce json
// @Param id path string true "Hospital ID"
// @Success 200 {object} sqlc.Hospital
// @Router /api/v1/hospitals/{id} [get]
func (h *HospitalsHandler) Get(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid hospital ID"})
		return
	}

	queries := sqlc.New(h.pool)
	hospital, err := queries.GetHospital(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "hospital not found"})
		return
	}

	c.JSON(http.StatusOK, hospital)
}

// ListHospitals godoc
// @Summary List all hospitals
// @Tags hospitals
// @Produce json
// @Success 200 {array} sqlc.Hospital
// @Router /api/v1/hospitals [get]
func (h *HospitalsHandler) List(c *gin.Context) {
	queries := sqlc.New(h.pool)
	hospitals, err := queries.ListHospitals(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list hospitals"})
		return
	}

	c.JSON(http.StatusOK, hospitals)
}

// UpdateHospital godoc
// @Summary Update a hospital
// @Tags hospitals
// @Accept json
// @Produce json
// @Param id path string true "Hospital ID"
// @Param hospital body UpdateHospitalRequest true "Hospital data"
// @Success 200 {object} sqlc.Hospital
// @Router /api/v1/hospitals/{id} [put]
func (h *HospitalsHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid hospital ID"})
		return
	}

	var req struct {
		Name    string `json:"name" binding:"required"`
		Tier    string `json:"tier" binding:"required"`
		Region  string `json:"region" binding:"required"`
		Address string `json:"address"`
		Phone   string `json:"phone"`
		Email   string `json:"email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	queries := sqlc.New(h.pool)
	hospital, err := queries.UpdateHospital(c.Request.Context(), sqlc.UpdateHospitalParams{
		ID:      id,
		Name:    req.Name,
		Tier:    req.Tier,
		Region:  req.Region,
		Address: pgtype.Text{String: req.Address, Valid: req.Address != ""},
		Phone:   pgtype.Text{String: req.Phone, Valid: req.Phone != ""},
		Email:   pgtype.Text{String: req.Email, Valid: req.Email != ""},
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update hospital"})
		return
	}

	c.JSON(http.StatusOK, hospital)
}

// DeleteHospital godoc
// @Summary Soft delete a hospital
// @Tags hospitals
// @Param id path string true "Hospital ID"
// @Success 204
// @Router /api/v1/hospitals/{id} [delete]
func (h *HospitalsHandler) Delete(c *gin.Context) {
	// Get hospital_id from context (set by JWT middleware)
	hospitalID, exists := c.Get(middleware.CtxHospitalID)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing hospital context"})
		return
	}

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid hospital ID"})
		return
	}

	// Verify user is deleting their own hospital (or has super_admin role)
	if hospitalID.(string) != id.String() {
		c.JSON(http.StatusForbidden, gin.H{"error": "cannot delete other hospitals"})
		return
	}

	queries := sqlc.New(h.pool)
	err = queries.SoftDeleteHospital(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete hospital"})
		return
	}

	c.Status(http.StatusNoContent)
}
