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

type DashboardHandler struct {
	pool *pgxpool.Pool
}

func NewDashboardHandler(pool *pgxpool.Pool) *DashboardHandler {
	return &DashboardHandler{pool: pool}
}

// DashboardStats is the response for GET /dashboard/stats
type DashboardStats struct {
	SessionsByStatus     []sqlc.CountSessionsByStatusForDateRow `json:"sessions_by_status"`
	ActiveSessionCount   int64                                  `json:"active_session_count"`
	CriticalAlertCount   int64                                  `json:"critical_alert_count"`
	OverdueInvoiceCount  int64                                  `json:"overdue_invoice_count"`
	LowStockCount        int64                                  `json:"low_stock_count"`
	StaffOnDutyCount     int64                                  `json:"staff_on_duty_count"`
	ActivePatientCount   int64                                  `json:"active_patient_count"`
	Date                 string                                 `json:"date"`
}

// GetStats returns dashboard statistics for the current hospital.
// GET /api/v1/dashboard/stats
func (h *DashboardHandler) GetStats(c *gin.Context) {
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
	today := time.Now()
	todayDate := pgtype.Date{Time: today, Valid: true}

	stats := DashboardStats{
		Date: today.Format("2006-01-02"),
	}

	// Sessions by status for today
	sessionCounts, err := queries.CountSessionsByStatusForDate(ctx, sqlc.CountSessionsByStatusForDateParams{
		HospitalID:    hospitalID,
		ScheduledDate: todayDate,
	})
	if err == nil {
		stats.SessionsByStatus = sessionCounts
	}

	// Active sessions right now
	activeCount, err := queries.CountActiveSessions(ctx, hospitalID)
	if err == nil {
		stats.ActiveSessionCount = activeCount
	}

	// Unacknowledged critical alerts
	alertCount, err := queries.CountUnacknowledgedCriticalAlerts(ctx, hospitalID)
	if err == nil {
		stats.CriticalAlertCount = alertCount
	}

	// Overdue invoices
	overdueCount, err := queries.CountOverdueInvoices(ctx, hospitalID)
	if err == nil {
		stats.OverdueInvoiceCount = overdueCount
	}

	// Low stock items
	lowStockCount, err := queries.CountLowStockItems(ctx, hospitalID)
	if err == nil {
		stats.LowStockCount = lowStockCount
	}

	// Staff on duty today
	staffCount, err := queries.CountStaffOnDutyToday(ctx, sqlc.CountStaffOnDutyTodayParams{
		HospitalID: hospitalID,
		ShiftDate:  todayDate,
	})
	if err == nil {
		stats.StaffOnDutyCount = staffCount
	}

	// Active patients
	patientCount, err := queries.CountActivePatients(ctx, hospitalID)
	if err == nil {
		stats.ActivePatientCount = patientCount
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, stats)
}
