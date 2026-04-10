package middleware

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Module tier definitions based on subscription plans
var moduleTiers = map[string][]string{
	"basic": {
		// Core 48 essential tables
		"patients", "sessions", "vitals", "vascular_access",
		"medications", "outcomes", "billing",
	},
	"standard": {
		// Essential + Recommended (68 tables)
		"patients", "sessions", "vitals", "vascular_access",
		"medications", "outcomes", "billing", "lab_alerts",
		"equipment_maintenance", "staff_schedules",
	},
	"enterprise": {
		// All 93 tables
		"patients", "sessions", "vitals", "vascular_access",
		"medications", "outcomes", "billing", "lab_management",
		"full_pharmacy", "hr_management", "inventory_tracking",
		"advanced_billing", "imaging_integration",
	},
}

// ModuleConfig represents the enabled_modules JSONB structure
type ModuleConfig struct {
	LabManagement      bool `json:"lab_management"`
	FullPharmacy       bool `json:"full_pharmacy"`
	HRManagement       bool `json:"hr_management"`
	InventoryTracking  bool `json:"inventory_tracking"`
	AdvancedBilling    bool `json:"advanced_billing"`
	OfflineSync        bool `json:"offline_sync"`
	CHWProgram         bool `json:"chw_program"`
	ImagingIntegration bool `json:"imaging_integration"`
	OutcomesReporting  bool `json:"outcomes_reporting"`
}

// RequireModule middleware checks if a module is enabled for the hospital
func RequireModule(pool *pgxpool.Pool, moduleName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		hospitalID := c.GetString(CtxHospitalID)
		if hospitalID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Hospital ID not found"})
			c.Abort()
			return
		}

		// Query hospital settings
		var subscriptionPlan string
		var enabledModulesJSON []byte

		query := `
			SELECT subscription_plan, enabled_modules
			FROM hospitals
			WHERE id = $1 AND deleted_at IS NULL
		`

		err := pool.QueryRow(context.Background(), query, hospitalID).Scan(&subscriptionPlan, &enabledModulesJSON)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check module access"})
			c.Abort()
			return
		}

		// Parse enabled_modules JSONB
		var moduleConfig ModuleConfig
		if err := json.Unmarshal(enabledModulesJSON, &moduleConfig); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid module configuration"})
			c.Abort()
			return
		}

		// Check if module is allowed in the subscription tier
		if !isModuleAllowedInTier(moduleName, subscriptionPlan) {
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "Module not available in your subscription plan",
				"module":  moduleName,
				"plan":    subscriptionPlan,
				"upgrade": "Please upgrade to access this feature",
			})
			c.Abort()
			return
		}

		// Check if module is enabled via feature flag
		if !isModuleEnabled(moduleName, moduleConfig) {
			c.JSON(http.StatusForbidden, gin.H{
				"error":  "Module is disabled for your hospital",
				"module": moduleName,
			})
			c.Abort()
			return
		}

		// Store subscription plan in context for handlers
		c.Set("subscription_plan", subscriptionPlan)
		c.Set("enabled_modules", moduleConfig)

		c.Next()
	}
}

// isModuleAllowedInTier checks if a module is available in the subscription tier
func isModuleAllowedInTier(moduleName, plan string) bool {
	// Enterprise has access to everything
	if plan == "enterprise" {
		return true
	}

	// Check if module is in the plan's allowed modules
	allowedModules, exists := moduleTiers[plan]
	if !exists {
		return false
	}

	for _, allowed := range allowedModules {
		if allowed == moduleName {
			return true
		}
	}

	return false
}

// isModuleEnabled checks if a module is enabled via feature flags
func isModuleEnabled(moduleName string, config ModuleConfig) bool {
	switch moduleName {
	case "lab_management":
		return config.LabManagement
	case "full_pharmacy":
		return config.FullPharmacy
	case "hr_management":
		return config.HRManagement
	case "inventory_tracking":
		return config.InventoryTracking
	case "advanced_billing":
		return config.AdvancedBilling
	case "offline_sync":
		return config.OfflineSync
	case "chw_program":
		return config.CHWProgram
	case "imaging_integration":
		return config.ImagingIntegration
	case "outcomes_reporting":
		return config.OutcomesReporting
	default:
		// Core modules (patients, sessions, vitals, etc.) are always enabled
		return true
	}
}

// GetHospitalPlan returns the subscription plan for a hospital
func GetHospitalPlan(pool *pgxpool.Pool, hospitalID string) (string, error) {
	var plan string
	query := `SELECT subscription_plan FROM hospitals WHERE id = $1 AND deleted_at IS NULL`
	err := pool.QueryRow(context.Background(), query, hospitalID).Scan(&plan)
	return plan, err
}
