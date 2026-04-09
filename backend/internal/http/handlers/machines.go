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

type MachinesHandler struct {
	pool *pgxpool.Pool
}

func NewMachinesHandler(pool *pgxpool.Pool) *MachinesHandler {
	return &MachinesHandler{pool: pool}
}

func (h *MachinesHandler) CreateMachine(c *gin.Context) {
	var req struct {
		MachineCode       string  `json:"machine_code" binding:"required"`
		SerialNumber      string  `json:"serial_number" binding:"required"`
		Model             string  `json:"model" binding:"required"`
		Manufacturer      string  `json:"manufacturer" binding:"required"`
		ManufactureYear   *int32  `json:"manufacture_year"`
		InstallationDate  *string `json:"installation_date"`
		Location          *string `json:"location"`
		IsHBVDedicated    bool    `json:"is_hbv_dedicated"`
		LastServiceDate   *string `json:"last_service_date"`
		NextServiceDate   *string `json:"next_service_date"`
		Notes             *string `json:"notes"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

	var installDate, lastService, nextService pgtype.Date
	if req.InstallationDate != nil {
		t, _ := time.Parse("2006-01-02", *req.InstallationDate)
		installDate = pgtype.Date{Time: t, Valid: true}
	}
	if req.LastServiceDate != nil {
		t, _ := time.Parse("2006-01-02", *req.LastServiceDate)
		lastService = pgtype.Date{Time: t, Valid: true}
	}
	if req.NextServiceDate != nil {
		t, _ := time.Parse("2006-01-02", *req.NextServiceDate)
		nextService = pgtype.Date{Time: t, Valid: true}
	}

	machine, err := queries.CreateDialysisMachine(ctx, sqlc.CreateDialysisMachineParams{
		HospitalID:       hospitalID,
		MachineCode:      req.MachineCode,
		SerialNumber:     req.SerialNumber,
		Model:            req.Model,
		Manufacturer:     req.Manufacturer,
		ManufactureYear:  pgtype.Int4{Int32: *req.ManufactureYear, Valid: req.ManufactureYear != nil},
		InstallationDate: installDate,
		Location:         pgtype.Text{String: *req.Location, Valid: req.Location != nil},
		Status:           sqlc.MachineStatusAvailable,
		IsHbvDedicated:   req.IsHBVDedicated,
		LastServiceDate:  lastService,
		NextServiceDate:  nextService,
		Notes:            pgtype.Text{String: *req.Notes, Valid: req.Notes != nil},
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create machine"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusCreated, machine)
}

func (h *MachinesHandler) ListMachines(c *gin.Context) {
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
	machines, err := queries.ListDialysisMachinesByHospital(ctx, hospitalID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list machines"})
		return
	}

	tx.Commit(ctx)
	c.JSON(http.StatusOK, machines)
}

func (h *MachinesHandler) ListAvailableMachines(c *gin.Context) {
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
	machines, err := queries.ListAvailableMachines(ctx, hospitalID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list available machines"})
		return
	}

	tx.Commit(ctx)
	c.JSON(http.StatusOK, machines)
}

func (h *MachinesHandler) UpdateMachineStatus(c *gin.Context) {
	machineID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid machine ID"})
		return
	}

	var req struct {
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
	machine, err := queries.UpdateMachineStatus(ctx, sqlc.UpdateMachineStatusParams{
		ID:     machineID,
		Status: sqlc.MachineStatus(req.Status),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update machine status"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, machine)
}
