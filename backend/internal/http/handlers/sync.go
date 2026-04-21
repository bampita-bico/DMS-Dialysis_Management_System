package handlers

import (
	"net/http"

	"github.com/dmsafrica/dms/internal/db/sqlc"
	"github.com/dmsafrica/dms/internal/http/middleware"
	"github.com/dmsafrica/dms/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SyncHandler struct {
	pool        *pgxpool.Pool
	syncService *services.SyncService
}

func NewSyncHandler(pool *pgxpool.Pool) *SyncHandler {
	return &SyncHandler{
		pool:        pool,
		syncService: services.NewSyncService(pool),
	}
}

// GetSyncStatus returns current sync status for the hospital
// GET /api/v1/sync/status
func (h *SyncHandler) GetSyncStatus(c *gin.Context) {
	hospitalIDStr := c.GetString(middleware.CtxHospitalID)
	hospitalID, err := uuid.Parse(hospitalIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid hospital ID"})
		return
	}

	stats, err := h.syncService.GetSyncStats(c.Request.Context(), hospitalID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get sync stats"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"hospital_id": hospitalID,
		"stats":       stats,
		"timestamp":   "now",
	})
}

// TriggerSync manually triggers sync processing
// POST /api/v1/sync/trigger
func (h *SyncHandler) TriggerSync(c *gin.Context) {
	hospitalIDStr := c.GetString(middleware.CtxHospitalID)
	hospitalID, err := uuid.Parse(hospitalIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid hospital ID"})
		return
	}

	// Process up to 100 items
	processed, err := h.syncService.ProcessSyncQueue(c.Request.Context(), hospitalID, 100)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Sync failed", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Sync completed",
		"processed": processed,
	})
}

// ListConflicts returns unresolved sync conflicts
// GET /api/v1/sync/conflicts
func (h *SyncHandler) ListConflicts(c *gin.Context) {
	hospitalIDStr := c.GetString(middleware.CtxHospitalID)
	hospitalID, err := uuid.Parse(hospitalIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid hospital ID"})
		return
	}

	queries := sqlc.New(h.pool)
	conflicts, err := queries.ListPendingConflicts(c.Request.Context(), hospitalID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list conflicts"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"conflicts": conflicts,
		"count":     len(conflicts),
	})
}

// GetConflict returns details of a specific conflict
// GET /api/v1/sync/conflicts/:id
func (h *SyncHandler) GetConflict(c *gin.Context) {
	conflictID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid conflict ID"})
		return
	}

	queries := sqlc.New(h.pool)
	conflict, err := queries.GetSyncConflict(c.Request.Context(), conflictID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Conflict not found"})
		return
	}

	c.JSON(http.StatusOK, conflict)
}

// ResolveConflict resolves a sync conflict
// POST /api/v1/sync/conflicts/:id/resolve
func (h *SyncHandler) ResolveConflict(c *gin.Context) {
	conflictID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid conflict ID"})
		return
	}

	var req struct {
		Resolution string  `json:"resolution" binding:"required,oneof=use_server use_client"`
		ResolvedBy uuid.UUID `json:"resolved_by" binding:"required"`
		Notes      string  `json:"notes"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	queries := sqlc.New(h.pool)

	// Convert types to pgtype
	resolution := sqlc.ConflictResolution(req.Resolution)
	resolvedBy := pgtype.UUID{Bytes: req.ResolvedBy, Valid: true}
	notes := pgtype.Text{String: req.Notes, Valid: true}

	err = queries.ResolveConflict(c.Request.Context(), sqlc.ResolveConflictParams{
		ID:         conflictID,
		Resolution: resolution,
		ResolvedBy: resolvedBy,
		Notes:      notes,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to resolve conflict"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Conflict resolved",
		"conflict_id": conflictID,
		"resolution":  req.Resolution,
	})
}

// RetryFailedSyncs requeues failed sync items
// POST /api/v1/sync/retry
func (h *SyncHandler) RetryFailedSyncs(c *gin.Context) {
	hospitalIDStr := c.GetString(middleware.CtxHospitalID)
	hospitalID, err := uuid.Parse(hospitalIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid hospital ID"})
		return
	}

	err = h.syncService.RequeueFailedSyncs(c.Request.Context(), hospitalID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retry syncs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Failed syncs requeued",
	})
}
