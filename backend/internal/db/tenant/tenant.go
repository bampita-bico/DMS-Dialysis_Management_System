package tenant

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

// SetLocalHospitalID sets the per-transaction tenant context used by RLS policies.
// Call after BEGIN and before any tenant-scoped queries.
func SetLocalHospitalID(ctx context.Context, tx pgx.Tx, hospitalID string) error {
	_, err := tx.Exec(ctx, "SET LOCAL app.hospital_id = $1", hospitalID)
	if err != nil {
		return fmt.Errorf("set tenant: %w", err)
	}
	return nil
}
