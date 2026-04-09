package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/dmsafrica/dms/internal/config"
	"github.com/dmsafrica/dms/internal/db/pool"
	"github.com/google/uuid"
)

// DefaultRole represents a system role with permissions
type DefaultRole struct {
	Name        string
	Description string
	Permissions []string
	IsSystem    bool
}

var defaultRoles = []DefaultRole{
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
		IsSystem:    true,
	},
	{
		Name:        "doctor",
		Description: "Medical doctor - can prescribe, view patient records, manage sessions",
		Permissions: []string{
			"patients:read", "patients:write",
			"sessions:read", "sessions:write",
			"prescriptions:*", "lab_orders:*",
			"diagnoses:*", "medical_records:*",
		},
		IsSystem:    true,
	},
	{
		Name:        "nurse",
		Description: "Nurse - can administer medication, monitor sessions, record vitals",
		Permissions: []string{
			"patients:read",
			"sessions:read", "sessions:write",
			"vitals:*", "medications:administer",
			"nursing_notes:*", "complications:*",
		},
		IsSystem:    true,
	},
	{
		Name:        "lab_technician",
		Description: "Laboratory technician - processes lab orders and enters results",
		Permissions: []string{
			"patients:read",
			"lab_orders:read", "lab_results:*",
			"specimens:*",
		},
		IsSystem:    true,
	},
	{
		Name:        "pharmacist",
		Description: "Pharmacist - manages medication inventory and dispensing",
		Permissions: []string{
			"patients:read",
			"prescriptions:read", "prescriptions:verify",
			"pharmacy_stock:*", "medications:dispense",
		},
		IsSystem:    true,
	},
	{
		Name:        "receptionist",
		Description: "Receptionist - patient registration, appointments, check-in",
		Permissions: []string{
			"patients:read", "patients:write",
			"appointments:*", "check_in:*",
			"patient_contacts:*",
		},
		IsSystem:    true,
	},
	{
		Name:        "biomedical_engineer",
		Description: "Biomedical engineer - equipment maintenance and calibration",
		Permissions: []string{
			"equipment:*", "machines:*",
			"maintenance:*", "calibration:*",
			"water_treatment:*",
		},
		IsSystem:    true,
	},
	{
		Name:        "finance_officer",
		Description: "Finance officer - billing, payments, financial reports",
		Permissions: []string{
			"patients:read",
			"invoices:*", "payments:*",
			"financial_reports:*", "billing:*",
			"insurance_claims:*",
		},
		IsSystem:    true,
	},
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run seed_roles.go <hospital_id>")
	}

	hospitalIDStr := os.Args[1]
	hospitalID, err := uuid.Parse(hospitalIDStr)
	if err != nil {
		log.Fatalf("Invalid hospital_id: %v", err)
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	ctx := context.Background()
	pg, err := pool.New(ctx, cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pg.Close()

	fmt.Printf("Seeding default roles for hospital: %s\n", hospitalID)

	for _, role := range defaultRoles {
		permissionsJSON, err := json.Marshal(role.Permissions)
		if err != nil {
			log.Printf("Failed to marshal permissions for %s: %v", role.Name, err)
			continue
		}

		var roleID uuid.UUID
		err = pg.QueryRow(
			ctx,
			`INSERT INTO roles (hospital_id, name, description, permissions, is_system)
			VALUES ($1, $2, $3, $4, $5)
			ON CONFLICT (hospital_id, name) DO UPDATE
			SET description = EXCLUDED.description,
			    permissions = EXCLUDED.permissions
			RETURNING id`,
			hospitalID,
			role.Name,
			role.Description,
			permissionsJSON,
			role.IsSystem,
		).Scan(&roleID)

		if err != nil {
			log.Printf("Failed to insert role %s: %v", role.Name, err)
			continue
		}

		fmt.Printf("✓ Created/Updated role: %s (ID: %s)\n", role.Name, roleID)
	}

	fmt.Println("\nRole seeding complete!")
	fmt.Println("\nDefault roles:")
	for _, role := range defaultRoles {
		fmt.Printf("  - %s: %s\n", role.Name, role.Description)
	}
}
