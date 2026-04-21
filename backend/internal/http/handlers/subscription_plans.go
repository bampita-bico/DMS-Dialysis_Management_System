package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/dmsafrica/dms/internal/db/sqlc"
	"github.com/dmsafrica/dms/internal/http/middleware"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SubscriptionPlansHandler struct {
	pool *pgxpool.Pool
}

func NewSubscriptionPlansHandler(pool *pgxpool.Pool) *SubscriptionPlansHandler {
	return &SubscriptionPlansHandler{pool: pool}
}

// GetCurrentPlan returns the hospital's current subscription plan
// GET /api/v1/subscription/plan
func (h *SubscriptionPlansHandler) GetCurrentPlan(c *gin.Context) {
	hospitalIDStr := c.GetString(middleware.CtxHospitalID)
	hospitalID, err := uuid.Parse(hospitalIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid hospital ID"})
		return
	}

	queries := sqlc.New(h.pool)
	planInfo, err := queries.GetHospitalPlan(c.Request.Context(), hospitalID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get subscription plan"})
		return
	}

	// Parse enabled_modules JSONB
	var enabledModules map[string]bool
	if err := json.Unmarshal(planInfo.EnabledModules, &enabledModules); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse modules"})
		return
	}

	// Get tier capabilities
	capabilities := getTierCapabilities(planInfo.SubscriptionPlan)

	c.JSON(http.StatusOK, gin.H{
		"subscription_plan": planInfo.SubscriptionPlan,
		"enabled_modules":   enabledModules,
		"capabilities":      capabilities,
		"upgrade_available": planInfo.SubscriptionPlan != "enterprise",
	})
}

// UpdatePlan updates the hospital's subscription plan
// PUT /api/v1/subscription/plan
// Admin/System endpoint only
func (h *SubscriptionPlansHandler) UpdatePlan(c *gin.Context) {
	hospitalIDStr := c.GetString(middleware.CtxHospitalID)
	hospitalID, err := uuid.Parse(hospitalIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid hospital ID"})
		return
	}

	var req struct {
		SubscriptionPlan string `json:"subscription_plan" binding:"required,oneof=basic standard enterprise"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	queries := sqlc.New(h.pool)
	err = queries.UpdateHospitalPlan(c.Request.Context(), sqlc.UpdateHospitalPlanParams{
		ID:               hospitalID,
		SubscriptionPlan: req.SubscriptionPlan,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update plan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":           "Subscription plan updated successfully",
		"subscription_plan": req.SubscriptionPlan,
	})
}

// UpdateModules updates the hospital's enabled modules
// PUT /api/v1/subscription/modules
// Admin endpoint only
func (h *SubscriptionPlansHandler) UpdateModules(c *gin.Context) {
	hospitalIDStr := c.GetString(middleware.CtxHospitalID)
	hospitalID, err := uuid.Parse(hospitalIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid hospital ID"})
		return
	}

	var req map[string]bool
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate module names
	validModules := map[string]bool{
		"lab_management":      true,
		"full_pharmacy":       true,
		"hr_management":       true,
		"inventory_tracking":  true,
		"advanced_billing":    true,
		"offline_sync":        true,
		"chw_program":         true,
		"imaging_integration": true,
		"outcomes_reporting":  true,
	}

	for moduleName := range req {
		if !validModules[moduleName] {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":  "Invalid module name",
				"module": moduleName,
			})
			return
		}
	}

	// Convert to JSONB
	modulesJSON, err := json.Marshal(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encode modules"})
		return
	}

	queries := sqlc.New(h.pool)
	err = queries.UpdateEnabledModules(c.Request.Context(), sqlc.UpdateEnabledModulesParams{
		ID:             hospitalID,
		EnabledModules: modulesJSON,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update modules"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":         "Modules updated successfully",
		"enabled_modules": req,
	})
}

// ListPlans returns all available subscription plans with details
// GET /api/v1/subscription/plans
func (h *SubscriptionPlansHandler) ListPlans(c *gin.Context) {
	plans := []gin.H{
		{
			"tier":        "basic",
			"name":        "Basic Plan",
			"price":       1200000,
			"currency":    "UGX",
			"interval":    "month",
			"description": "Essential dialysis operations for District Hospitals (5-10 beds)",
			"features": []string{
				"Core patient management (unlimited)",
				"Dialysis session tracking",
				"Vitals monitoring",
				"Vascular access management",
				"Basic billing (invoices & payments)",
				"Clinical outcomes reporting",
				"Water treatment logs",
				"Basic lab integration",
			},
			"table_count": 48,
			"max_patients": 0,
			"government_price": 15000000,
			"government_interval": "one-time (MOH purchase)",
		},
		{
			"tier":        "standard",
			"name":        "Standard Plan",
			"price":       2800000,
			"currency":    "UGX",
			"interval":    "month",
			"description": "Comprehensive features for Regional Referral Hospitals (10-30 beds)",
			"features": []string{
				"Everything in Basic",
				"Full lab management",
				"Complete pharmacy module",
				"HR & staff scheduling",
				"Inventory tracking",
				"Advanced billing & insurance claims",
				"Equipment maintenance tracking",
				"Staff performance tracking",
				"Audit logs & notifications",
			},
			"table_count":  68,
			"max_patients": 0,
			"recommended":  true,
			"government_price": 35000000,
			"government_interval": "one-time (MOH purchase)",
		},
		{
			"tier":        "enterprise",
			"name":        "Enterprise Plan",
			"price":       6000000,
			"currency":    "UGX",
			"interval":    "month",
			"description": "Full feature set for National Referral Hospitals (30+ beds)",
			"features": []string{
				"Everything in Standard",
				"Offline mobile sync",
				"Community Health Workers (CHW) module",
				"Imaging integration (PACS)",
				"Advanced analytics & dashboards",
				"Multi-site management",
				"API access for integrations",
				"24/7 priority support",
				"Dedicated account manager",
			},
			"table_count":  93,
			"max_patients": 0, // Unlimited
			"government_price": 75000000,
			"government_interval": "one-time (MOH purchase)",
			"custom_pricing": true,
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"plans": plans,
		"note": "Government hospitals: MOH one-time purchase. Private hospitals: Monthly subscription (10% discount on annual).",
	})
}

// getTierCapabilities returns the capabilities for each tier
func getTierCapabilities(tier string) gin.H {
	baseCapabilities := gin.H{
		"patients":         true,
		"sessions":         true,
		"vitals":           true,
		"vascular_access":  true,
		"medications":      true,
		"outcomes":         true,
		"billing":          true,
		"water_treatment":  true,
	}

	standardExtras := gin.H{
		"lab_alerts":            true,
		"equipment_maintenance": true,
		"staff_schedules":       true,
		"insurance_claims":      true,
		"audit_logs":            true,
		"notifications":         true,
	}

	enterpriseExtras := gin.H{
		"lab_management":      true,
		"full_pharmacy":       true,
		"hr_management":       true,
		"inventory_tracking":  true,
		"advanced_billing":    true,
		"imaging_integration": true,
		"chw_program":         true,
		"offline_sync":        true,
	}

	switch tier {
	case "basic":
		return baseCapabilities
	case "standard":
		for k, v := range standardExtras {
			baseCapabilities[k] = v
		}
		return baseCapabilities
	case "enterprise":
		for k, v := range standardExtras {
			baseCapabilities[k] = v
		}
		for k, v := range enterpriseExtras {
			baseCapabilities[k] = v
		}
		return baseCapabilities
	default:
		return baseCapabilities
	}
}
