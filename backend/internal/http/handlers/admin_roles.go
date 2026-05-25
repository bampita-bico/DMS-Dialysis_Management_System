package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/dmsafrica/dms/internal/db/sqlc"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type defaultRole struct {
	Name        string
	Description string
	Permissions []string
	IsSystem    bool
}

var defaultRoles = []defaultRole{
	{
		Name:        "super_admin",
		Description: "System administrator with full access to all features",
		Permissions: []string{"*"},
		IsSystem:    true,
	},
	{
		Name:        "admin",
		Description: "Hospital administrator with access to most features",
		Permissions: []string{
			"users:*", "departments:*", "settings:*",
			"patients:*", "sessions:*", "reports:*",
			"inventory:*", "billing:*",
		},
		IsSystem: true,
	},
	{
		Name:        "doctor",
		Description: "Medical doctor with full clinical access",
		Permissions: []string{
			"patients:read", "patients:write",
			"sessions:read", "sessions:write",
			"prescriptions:*", "lab_orders:*",
			"diagnoses:*", "medical_records:*",
		},
		IsSystem: true,
	},
	{
		Name:        "nurse",
		Description: "Nurse with session monitoring and vitals access",
		Permissions: []string{
			"patients:read",
			"sessions:read", "sessions:write",
			"vitals:*", "nursing_notes:*", "complications:*",
		},
		IsSystem: true,
	},
	{
		Name:        "lab_technician",
		Description: "Laboratory technician",
		Permissions: []string{
			"patients:read",
			"lab_orders:read", "lab_results:*", "specimens:*",
		},
		IsSystem: true,
	},
	{
		Name:        "pharmacist",
		Description: "Pharmacist with dispensing and stock access",
		Permissions: []string{
			"patients:read",
			"prescriptions:read", "prescriptions:verify",
			"pharmacy_stock:*", "medications:dispense",
		},
		IsSystem: true,
	},
	{
		Name:        "receptionist",
		Description: "Reception and registration staff",
		Permissions: []string{
			"patients:read", "patients:write",
			"appointments:*", "check_in:*", "patient_contacts:*",
		},
		IsSystem: true,
	},
}

func ensureDefaultRoles(ctx context.Context, pool *pgxpool.Pool, hospitalID uuid.UUID) error {
	for _, role := range defaultRoles {
		permissionsJSON, err := json.Marshal(role.Permissions)
		if err != nil {
			return err
		}

		const q = `
			INSERT INTO roles (hospital_id, name, description, permissions, is_system)
			VALUES ($1, $2, $3, $4, $5)
			ON CONFLICT (hospital_id, name) DO UPDATE
			SET description = EXCLUDED.description,
			    permissions = EXCLUDED.permissions,
			    is_system = EXCLUDED.is_system,
			    deleted_at = NULL
		`

		if _, err := pool.Exec(ctx, q, hospitalID, role.Name, role.Description, permissionsJSON, role.IsSystem); err != nil {
			return err
		}
	}

	return nil
}

func findRoleByName(ctx context.Context, queries *sqlc.Queries, hospitalID uuid.UUID, roleName string) (sqlc.Role, error) {
	roles, err := queries.ListRoles(ctx, hospitalID)
	if err != nil {
		return sqlc.Role{}, err
	}

	for _, role := range roles {
		if strings.EqualFold(role.Name, roleName) {
			return role, nil
		}
	}

	return sqlc.Role{}, fmt.Errorf("role %q not found", roleName)
}

func textOrNull(value string) pgtype.Text {
	return pgtype.Text{String: value, Valid: value != ""}
}
