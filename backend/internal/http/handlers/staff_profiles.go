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

type StaffProfilesHandler struct {
	pool *pgxpool.Pool
}

func NewStaffProfilesHandler(pool *pgxpool.Pool) *StaffProfilesHandler {
	return &StaffProfilesHandler{pool: pool}
}

// Create creates a new staff profile
// POST /api/v1/staff-profiles
func (h *StaffProfilesHandler) Create(c *gin.Context) {
	var req struct {
		UserID                string  `json:"user_id" binding:"required"`
		DepartmentID          *string `json:"department_id"`
		Cadre                 string  `json:"cadre" binding:"required"`
		LicenseNumber         string  `json:"license_number"`
		LicenseExpiryDate     string  `json:"license_expiry_date"`
		RegistrationBody      string  `json:"registration_body"`
		Specialization        string  `json:"specialization"`
		YearsOfExperience     *int32  `json:"years_of_experience"`
		HireDate              string  `json:"hire_date"`
		EmployeeNumber        string  `json:"employee_number"`
		EmergencyContactName  string  `json:"emergency_contact_name"`
		EmergencyContactPhone string  `json:"emergency_contact_phone"`
		BloodType             string  `json:"blood_type"`
		Notes                 string  `json:"notes"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hospitalIDStr := c.GetString(middleware.CtxHospitalID)
	hospitalID, _ := uuid.Parse(hospitalIDStr)

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	// Start transaction with RLS
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

	// Prepare parameters
	var departmentID pgtype.UUID
	if req.DepartmentID != nil {
		deptID, err := uuid.Parse(*req.DepartmentID)
		if err == nil {
			departmentID = pgtype.UUID{Bytes: deptID, Valid: true}
		}
	}

	var licenseExpiryDate pgtype.Date
	if req.LicenseExpiryDate != "" {
		expDate, err := time.Parse("2006-01-02", req.LicenseExpiryDate)
		if err == nil {
			licenseExpiryDate = pgtype.Date{Time: expDate, Valid: true}
		}
	}

	var hireDate pgtype.Date
	if req.HireDate != "" {
		hDate, err := time.Parse("2006-01-02", req.HireDate)
		if err == nil {
			hireDate = pgtype.Date{Time: hDate, Valid: true}
		}
	}

	var yearsOfExperience pgtype.Int4
	if req.YearsOfExperience != nil {
		yearsOfExperience = pgtype.Int4{Int32: *req.YearsOfExperience, Valid: true}
	}

	var bloodType sqlc.NullBloodType
	if req.BloodType != "" {
		bloodType = sqlc.NullBloodType{BloodType: sqlc.BloodType(req.BloodType), Valid: true}
	}

	queries := sqlc.New(tx)
	profile, err := queries.CreateStaffProfile(ctx, sqlc.CreateStaffProfileParams{
		HospitalID:            hospitalID,
		UserID:                userID,
		DepartmentID:          departmentID,
		Cadre:                 sqlc.StaffCadre(req.Cadre),
		LicenseNumber:         pgtype.Text{String: req.LicenseNumber, Valid: req.LicenseNumber != ""},
		LicenseExpiryDate:     licenseExpiryDate,
		RegistrationBody:      pgtype.Text{String: req.RegistrationBody, Valid: req.RegistrationBody != ""},
		Specialization:        pgtype.Text{String: req.Specialization, Valid: req.Specialization != ""},
		YearsOfExperience:     yearsOfExperience,
		HireDate:              hireDate,
		EmployeeNumber:        pgtype.Text{String: req.EmployeeNumber, Valid: req.EmployeeNumber != ""},
		EmergencyContactName:  pgtype.Text{String: req.EmergencyContactName, Valid: req.EmergencyContactName != ""},
		EmergencyContactPhone: pgtype.Text{String: req.EmergencyContactPhone, Valid: req.EmergencyContactPhone != ""},
		BloodType:             bloodType,
		Notes:                 pgtype.Text{String: req.Notes, Valid: req.Notes != ""},
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create staff profile", "details": err.Error()})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusCreated, profile)
}

// Get retrieves a specific staff profile by ID
// GET /api/v1/staff-profiles/:id
func (h *StaffProfilesHandler) Get(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid staff profile ID"})
		return
	}

	hospitalIDStr := c.GetString(middleware.CtxHospitalID)
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
	profile, err := queries.GetStaffProfile(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "staff profile not found"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, profile)
}

// GetByUser retrieves a staff profile by user ID
// GET /api/v1/users/:user_id/staff-profile
func (h *StaffProfilesHandler) GetByUser(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

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
	profile, err := queries.GetStaffProfileByUser(ctx, sqlc.GetStaffProfileByUserParams{
		HospitalID: hospitalID,
		UserID:     userID,
	})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "staff profile not found"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, profile)
}

// List lists all staff profiles for the hospital
// GET /api/v1/staff-profiles
func (h *StaffProfilesHandler) List(c *gin.Context) {
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
	profiles, err := queries.ListStaffProfilesByHospital(ctx, hospitalID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list staff profiles"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, profiles)
}

// ListByCadre lists staff profiles filtered by cadre
// GET /api/v1/staff-profiles/cadre/:cadre
func (h *StaffProfilesHandler) ListByCadre(c *gin.Context) {
	cadre := c.Param("cadre")
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
	profiles, err := queries.ListStaffByCadre(ctx, sqlc.ListStaffByCadreParams{
		HospitalID: hospitalID,
		Cadre:      sqlc.StaffCadre(cadre),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list staff by cadre"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, profiles)
}

// ListActive lists only active staff profiles
// GET /api/v1/staff-profiles/active
func (h *StaffProfilesHandler) ListActive(c *gin.Context) {
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
	profiles, err := queries.ListActiveStaff(ctx, hospitalID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list active staff"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, profiles)
}

// ListByDepartment lists staff profiles for a department
// GET /api/v1/departments/:department_id/staff
func (h *StaffProfilesHandler) ListByDepartment(c *gin.Context) {
	departmentIDStr := c.Param("department_id")
	departmentID, err := uuid.Parse(departmentIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid department_id"})
		return
	}

	hospitalIDStr := c.GetString(middleware.CtxHospitalID)
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
	profiles, err := queries.ListStaffByDepartment(ctx, pgtype.UUID{Bytes: departmentID, Valid: true})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list staff by department"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, profiles)
}

// ListExpiringLicenses lists staff with licenses expiring soon
// GET /api/v1/staff-profiles/expiring-licenses?date=2024-12-31
func (h *StaffProfilesHandler) ListExpiringLicenses(c *gin.Context) {
	dateStr := c.Query("date")
	if dateStr == "" {
		// Default to 3 months from now
		dateStr = time.Now().AddDate(0, 3, 0).Format("2006-01-02")
	}

	expiryDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format"})
		return
	}

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
	profiles, err := queries.ListExpiringLicenses(ctx, sqlc.ListExpiringLicensesParams{
		HospitalID:        hospitalID,
		LicenseExpiryDate: pgtype.Date{Time: expiryDate, Valid: true},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list expiring licenses"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, profiles)
}

// Update updates a staff profile
// PATCH /api/v1/staff-profiles/:id
func (h *StaffProfilesHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid staff profile ID"})
		return
	}

	var req struct {
		Cadre               *string `json:"cadre"`
		LicenseNumber       *string `json:"license_number"`
		LicenseExpiryDate   *string `json:"license_expiry_date"`
		Specialization      *string `json:"specialization"`
		YearsOfExperience   *int32  `json:"years_of_experience"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hospitalIDStr := c.GetString(middleware.CtxHospitalID)
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

	// Get current profile to use as defaults
	currentProfile, err := queries.GetStaffProfile(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "staff profile not found"})
		return
	}

	// Prepare update params with current values as defaults
	cadre := currentProfile.Cadre
	if req.Cadre != nil {
		cadre = sqlc.StaffCadre(*req.Cadre)
	}

	licenseNumber := currentProfile.LicenseNumber
	if req.LicenseNumber != nil {
		licenseNumber = pgtype.Text{String: *req.LicenseNumber, Valid: *req.LicenseNumber != ""}
	}

	licenseExpiryDate := currentProfile.LicenseExpiryDate
	if req.LicenseExpiryDate != nil {
		expDate, err := time.Parse("2006-01-02", *req.LicenseExpiryDate)
		if err == nil {
			licenseExpiryDate = pgtype.Date{Time: expDate, Valid: true}
		}
	}

	specialization := currentProfile.Specialization
	if req.Specialization != nil {
		specialization = pgtype.Text{String: *req.Specialization, Valid: *req.Specialization != ""}
	}

	yearsOfExperience := currentProfile.YearsOfExperience
	if req.YearsOfExperience != nil {
		yearsOfExperience = pgtype.Int4{Int32: *req.YearsOfExperience, Valid: true}
	}

	profile, err := queries.UpdateStaffProfile(ctx, sqlc.UpdateStaffProfileParams{
		ID:                id,
		Cadre:             cadre,
		LicenseNumber:     licenseNumber,
		LicenseExpiryDate: licenseExpiryDate,
		Specialization:    specialization,
		YearsOfExperience: yearsOfExperience,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update staff profile"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, profile)
}
