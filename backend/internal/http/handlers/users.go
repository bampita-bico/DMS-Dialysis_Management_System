package handlers

import (
	"net/http"
	"strings"

	"github.com/dmsafrica/dms/internal/db/sqlc"
	"github.com/dmsafrica/dms/internal/http/middleware"
	"github.com/dmsafrica/dms/internal/security"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UsersHandler struct {
	pool *pgxpool.Pool
}

func NewUsersHandler(pool *pgxpool.Pool) *UsersHandler {
	return &UsersHandler{pool: pool}
}

func (h *UsersHandler) Create(c *gin.Context) {
	var req struct {
		HospitalID string `json:"hospital_id"`
		Email      string `json:"email" binding:"required,email"`
		Phone      string `json:"phone"`
		Password   string `json:"password" binding:"required,min=8"`
		FullName   string `json:"full_name" binding:"required"`
		IsActive   bool   `json:"is_active"`
		RoleName   string `json:"role_name"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get hospital_id from JWT context
	hospitalIDStr, _ := c.Get(middleware.CtxHospitalID)
	hospitalID, _ := uuid.Parse(hospitalIDStr.(string))
	if req.HospitalID != "" {
		var err error
		hospitalID, err = uuid.Parse(req.HospitalID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid hospital ID"})
			return
		}
	}

	// Hash password
	hashedPassword, err := security.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	queries := sqlc.New(h.pool)
	if err := ensureDefaultRoles(c.Request.Context(), h.pool, hospitalID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to prepare roles"})
		return
	}

	user, err := queries.CreateUser(c.Request.Context(), sqlc.CreateUserParams{
		HospitalID:   hospitalID,
		Email:        req.Email,
		Phone:        pgtype.Text{String: req.Phone, Valid: req.Phone != ""},
		PasswordHash: hashedPassword,
		FullName:     req.FullName,
		IsActive:     req.IsActive,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	roleName := req.RoleName
	if roleName == "" {
		roleName = "doctor"
	}
	role, err := findRoleByName(c.Request.Context(), queries, hospitalID, roleName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to resolve role"})
		return
	}

	assignerID := c.GetString(middleware.CtxUserID)
	assignerUUID, err := uuid.Parse(assignerID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid requester"})
		return
	}

	_, err = queries.AssignUserRole(c.Request.Context(), sqlc.AssignUserRoleParams{
		HospitalID:   hospitalID,
		UserID:       user.ID,
		RoleID:       role.ID,
		DepartmentID: pgtype.UUID{},
		AssignedBy:   pgtype.UUID{Bytes: assignerUUID, Valid: true},
		ExpiresAt:    pgtype.Timestamptz{},
	})
	if err != nil && !strings.Contains(strings.ToLower(err.Error()), "duplicate") {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to assign role"})
		return
	}

	// Don't return password_hash
	user.PasswordHash = ""

	c.JSON(http.StatusCreated, gin.H{
		"user":      user,
		"role_name": role.Name,
	})
}

func (h *UsersHandler) Get(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	queries := sqlc.New(h.pool)
	user, err := queries.GetUser(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	// Don't return password_hash
	user.PasswordHash = ""

	c.JSON(http.StatusOK, user)
}

func (h *UsersHandler) List(c *gin.Context) {
	hospitalIDStr := c.Query("hospital_id")
	if hospitalIDStr == "" {
		ctxHospitalID, _ := c.Get(middleware.CtxHospitalID)
		hospitalIDStr = ctxHospitalID.(string)
	}
	hospitalID, err := uuid.Parse(hospitalIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid hospital ID"})
		return
	}

	queries := sqlc.New(h.pool)
	users, err := queries.ListActiveUsers(c.Request.Context(), hospitalID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list users"})
		return
	}

	// Don't return password hashes
	for i := range users {
		users[i].PasswordHash = ""
	}

	c.JSON(http.StatusOK, users)
}

func (h *UsersHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Phone    string `json:"phone"`
		FullName string `json:"full_name" binding:"required"`
		IsActive bool   `json:"is_active"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	queries := sqlc.New(h.pool)
	user, err := queries.UpdateUser(c.Request.Context(), sqlc.UpdateUserParams{
		ID:       id,
		Email:    req.Email,
		Phone:    pgtype.Text{String: req.Phone, Valid: req.Phone != ""},
		FullName: req.FullName,
		IsActive: req.IsActive,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user"})
		return
	}

	user.PasswordHash = ""
	c.JSON(http.StatusOK, user)
}

func (h *UsersHandler) ResetPassword(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	var req struct {
		Password string `json:"password" binding:"required,min=8"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := security.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	queries := sqlc.New(h.pool)
	if err := queries.UpdateUserPassword(c.Request.Context(), sqlc.UpdateUserPasswordParams{
		ID:           id,
		PasswordHash: hashedPassword,
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to reset password"})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *UsersHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	queries := sqlc.New(h.pool)
	err = queries.SoftDeleteUser(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete user"})
		return
	}

	c.Status(http.StatusNoContent)
}
