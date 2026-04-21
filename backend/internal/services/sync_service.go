package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/dmsafrica/dms/internal/db/sqlc"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

// SyncService handles data synchronization and conflict detection
type SyncService struct {
	pool    *pgxpool.Pool
	queries *sqlc.Queries
}

// NewSyncService creates a new sync service instance
func NewSyncService(pool *pgxpool.Pool) *SyncService {
	return &SyncService{
		pool:    pool,
		queries: sqlc.New(pool),
	}
}

// ProcessSyncQueue processes pending sync items for a hospital
// Returns number of items processed and any error
func (s *SyncService) ProcessSyncQueue(ctx context.Context, hospitalID uuid.UUID, limit int32) (int, error) {
	// Get pending sync items
	items, err := s.queries.GetPendingSyncItems(ctx, sqlc.GetPendingSyncItemsParams{
		HospitalID: hospitalID,
		Limit:      limit,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to get pending sync items: %w", err)
	}

	if len(items) == 0 {
		return 0, nil // Nothing to process
	}

	log.Printf("Processing %d pending sync items for hospital %s", len(items), hospitalID)

	processed := 0
	for _, item := range items {
		if err := s.ProcessSyncItem(ctx, item); err != nil {
			log.Printf("Failed to process sync item %s: %v", item.ID, err)
			// Mark as failed but continue with other items
			errMsg := pgtype.Text{String: err.Error(), Valid: true}
			_ = s.queries.MarkSyncFailed(ctx, sqlc.MarkSyncFailedParams{
				ID:           item.ID,
				ErrorMessage: errMsg,
			})
			continue
		}

		// Mark as synced
		if err := s.queries.MarkSyncSynced(ctx, item.ID); err != nil {
			log.Printf("Failed to mark sync item %s as synced: %v", item.ID, err)
		}

		processed++
	}

	log.Printf("Successfully processed %d/%d sync items", processed, len(items))
	return processed, nil
}

// ProcessSyncItem processes a single sync item
func (s *SyncService) ProcessSyncItem(ctx context.Context, item sqlc.SyncQueue) error {
	// Check for conflicts
	conflict, err := s.DetectConflict(ctx, item)
	if err != nil {
		return fmt.Errorf("conflict detection failed: %w", err)
	}

	if conflict != nil {
		// Create conflict record
		log.Printf("Conflict detected for %s:%s", item.EntityType, item.EntityID)
		return s.CreateConflict(ctx, item, conflict)
	}

	// No conflict - item is already in database via API
	// Just validate and mark as processed
	return nil
}

// ConflictData represents conflicting versions of data
type ConflictData struct {
	ServerVersion map[string]interface{}
	ClientVersion map[string]interface{}
	UpdatedAt     time.Time
}

// DetectConflict checks if a sync item conflicts with server data
// Returns conflict data if found, nil if no conflict
func (s *SyncService) DetectConflict(ctx context.Context, item sqlc.SyncQueue) (*ConflictData, error) {
	// Parse payload
	var clientData map[string]interface{}
	if err := json.Unmarshal(item.Payload, &clientData); err != nil {
		return nil, fmt.Errorf("failed to parse payload: %w", err)
	}

	// Get client's updated_at timestamp
	clientUpdatedAtStr, ok := clientData["updated_at"].(string)
	if !ok {
		// No timestamp in payload, can't detect conflicts
		return nil, nil
	}

	clientUpdatedAt, err := time.Parse(time.RFC3339, clientUpdatedAtStr)
	if err != nil {
		log.Printf("Warning: invalid updated_at format: %v", err)
		return nil, nil
	}

	// Query current server data
	serverData, serverUpdatedAt, err := s.GetEntityData(ctx, item.EntityType, item.EntityID)
	if err != nil {
		// Entity doesn't exist on server (new creation), no conflict
		if err.Error() == "entity not found" {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get server data: %w", err)
	}

	// Compare timestamps
	// Conflict if server version is newer than client's base version
	if serverUpdatedAt.After(clientUpdatedAt) {
		return &ConflictData{
			ServerVersion: serverData,
			ClientVersion: clientData,
			UpdatedAt:     serverUpdatedAt,
		}, nil
	}

	// No conflict
	return nil, nil
}

// GetEntityData retrieves current server data for an entity
func (s *SyncService) GetEntityData(ctx context.Context, entityType string, entityID uuid.UUID) (map[string]interface{}, time.Time, error) {
	var query string
	var updatedAt time.Time

	// Map entity types to queries
	switch entityType {
	case "patients":
		query = `SELECT row_to_json(p.*) as data, updated_at FROM patients p WHERE id = $1 AND deleted_at IS NULL`
	case "dialysis_sessions":
		query = `SELECT row_to_json(d.*) as data, updated_at FROM dialysis_sessions d WHERE id = $1 AND deleted_at IS NULL`
	case "lab_orders":
		query = `SELECT row_to_json(l.*) as data, updated_at FROM lab_orders l WHERE id = $1 AND deleted_at IS NULL`
	case "prescriptions":
		query = `SELECT row_to_json(p.*) as data, updated_at FROM prescriptions p WHERE id = $1 AND deleted_at IS NULL`
	case "invoices":
		query = `SELECT row_to_json(i.*) as data, updated_at FROM invoices i WHERE id = $1 AND deleted_at IS NULL`
	default:
		return nil, time.Time{}, fmt.Errorf("unsupported entity type: %s", entityType)
	}

	var dataJSON []byte
	err := s.pool.QueryRow(ctx, query, entityID).Scan(&dataJSON, &updatedAt)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, time.Time{}, fmt.Errorf("entity not found")
		}
		return nil, time.Time{}, err
	}

	var data map[string]interface{}
	if err := json.Unmarshal(dataJSON, &data); err != nil {
		return nil, time.Time{}, fmt.Errorf("failed to parse entity data: %w", err)
	}

	return data, updatedAt, nil
}

// CreateConflict records a sync conflict
func (s *SyncService) CreateConflict(ctx context.Context, item sqlc.SyncQueue, conflict *ConflictData) error {
	serverDataJSON, err := json.Marshal(conflict.ServerVersion)
	if err != nil {
		return fmt.Errorf("failed to marshal server data: %w", err)
	}

	clientDataJSON := item.Payload

	_, err = s.queries.CreateSyncConflict(ctx, sqlc.CreateSyncConflictParams{
		HospitalID: item.HospitalID,
		QueueID:    item.ID,
		EntityType: item.EntityType,
		EntityID:   item.EntityID,
		LocalData:  clientDataJSON,
		ServerData: serverDataJSON,
	})

	if err != nil {
		return fmt.Errorf("failed to create conflict record: %w", err)
	}

	log.Printf("Created conflict record for %s:%s", item.EntityType, item.EntityID)
	return nil
}

// RequeueFailedSyncs resets failed sync items for retry (attempts < 3)
func (s *SyncService) RequeueFailedSyncs(ctx context.Context, hospitalID uuid.UUID) error {
	err := s.queries.RequeueFailedSyncs(ctx)
	if err != nil {
		return fmt.Errorf("failed to requeue failed syncs: %w", err)
	}

	log.Printf("Requeued failed sync items for hospital %s", hospitalID)
	return nil
}

// CleanupOldSyncedItems removes synced items older than N days
func (s *SyncService) CleanupOldSyncedItems(ctx context.Context, daysOld int) error {
	query := `
		DELETE FROM sync_queue
		WHERE synced_at IS NOT NULL
		  AND synced_at < NOW() - INTERVAL '1 day' * $1
	`

	result, err := s.pool.Exec(ctx, query, daysOld)
	if err != nil {
		return fmt.Errorf("failed to cleanup old sync items: %w", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected > 0 {
		log.Printf("Cleaned up %d old synced items (older than %d days)", rowsAffected, daysOld)
	}

	return nil
}

// GetSyncStats returns sync queue statistics for a hospital
func (s *SyncService) GetSyncStats(ctx context.Context, hospitalID uuid.UUID) (map[string]int, error) {
	query := `
		SELECT
			COUNT(*) FILTER (WHERE status = 'pending') as pending,
			COUNT(*) FILTER (WHERE status = 'synced') as synced,
			COUNT(*) FILTER (WHERE status = 'failed') as failed,
			COUNT(*) FILTER (WHERE attempts >= 3) as permanently_failed
		FROM sync_queue
		WHERE hospital_id = $1
	`

	var pending, synced, failed, permFailed int
	err := s.pool.QueryRow(ctx, query, hospitalID).Scan(&pending, &synced, &failed, &permFailed)
	if err != nil {
		return nil, fmt.Errorf("failed to get sync stats: %w", err)
	}

	return map[string]int{
		"pending":             pending,
		"synced":              synced,
		"failed":              failed,
		"permanently_failed":  permFailed,
		"total":               pending + synced + failed,
	}, nil
}
