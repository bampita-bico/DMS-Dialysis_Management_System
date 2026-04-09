package handlers

import (
	"net/http"

	"github.com/dmsafrica/dms/internal/db/sqlc"
	"github.com/dmsafrica/dms/internal/db/tenant"
	"github.com/dmsafrica/dms/internal/http/middleware"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type LabCatalogHandler struct {
	pool *pgxpool.Pool
}

func NewLabCatalogHandler(pool *pgxpool.Pool) *LabCatalogHandler {
	return &LabCatalogHandler{pool: pool}
}

func (h *LabCatalogHandler) ListTests(c *gin.Context) {
	hospitalIDStr, _ := c.Get(middleware.CtxHospitalID)
	hospitalID, _ := uuid.Parse(hospitalIDStr.(string))

	ctx := c.Request.Context()
	tx, err := h.pool.Begin(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}
	defer tx.Rollback(ctx)

	if err := tenant.SetLocalHospitalID(ctx, tx, hospitalID.String()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set tenant context"})
		return
	}

	queries := sqlc.New(tx)
	tests, err := queries.ListActiveLabTests(ctx, hospitalID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list tests"})
		return
	}

	tx.Commit(ctx)
	c.JSON(http.StatusOK, tests)
}

func (h *LabCatalogHandler) GetTest(c *gin.Context) {
	testID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid test ID"})
		return
	}

	hospitalIDStr, _ := c.Get(middleware.CtxHospitalID)
	hospitalID, _ := uuid.Parse(hospitalIDStr.(string))

	ctx := c.Request.Context()
	tx, err := h.pool.Begin(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}
	defer tx.Rollback(ctx)

	if err := tenant.SetLocalHospitalID(ctx, tx, hospitalID.String()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set tenant context"})
		return
	}

	queries := sqlc.New(tx)
	test, err := queries.GetLabTest(ctx, testID)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Test not found"})
		return
	}

	tx.Commit(ctx)
	c.JSON(http.StatusOK, test)
}

func (h *LabCatalogHandler) ListPanels(c *gin.Context) {
	hospitalIDStr, _ := c.Get(middleware.CtxHospitalID)
	hospitalID, _ := uuid.Parse(hospitalIDStr.(string))

	ctx := c.Request.Context()
	tx, err := h.pool.Begin(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}
	defer tx.Rollback(ctx)

	if err := tenant.SetLocalHospitalID(ctx, tx, hospitalID.String()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set tenant context"})
		return
	}

	queries := sqlc.New(tx)
	panels, err := queries.ListActiveLabPanels(ctx, hospitalID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list panels"})
		return
	}

	tx.Commit(ctx)
	c.JSON(http.StatusOK, panels)
}

func (h *LabCatalogHandler) GetPanel(c *gin.Context) {
	panelID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid panel ID"})
		return
	}

	hospitalIDStr, _ := c.Get(middleware.CtxHospitalID)
	hospitalID, _ := uuid.Parse(hospitalIDStr.(string))

	ctx := c.Request.Context()
	tx, err := h.pool.Begin(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}
	defer tx.Rollback(ctx)

	if err := tenant.SetLocalHospitalID(ctx, tx, hospitalID.String()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set tenant context"})
		return
	}

	queries := sqlc.New(tx)
	panel, err := queries.GetLabPanel(ctx, panelID)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Panel not found"})
		return
	}

	tx.Commit(ctx)
	c.JSON(http.StatusOK, panel)
}
