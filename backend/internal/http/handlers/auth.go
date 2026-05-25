package handlers

import (
	"net/http"
	"strings"
	"time"

	"github.com/dmsafrica/dms/internal/db/sqlc"
	"github.com/dmsafrica/dms/internal/security"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthHandler struct {
	pool   *pgxpool.Pool
	jwtSvc *security.JWTService
}

func NewAuthHandler(p *pgxpool.Pool, jwtSvc *security.JWTService) *AuthHandler {
	return &AuthHandler{
		pool:   p,
		jwtSvc: jwtSvc,
	}
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string      `json:"token"`
	User  interface{} `json:"user"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	queries := sqlc.New(h.pool)
	identifier := strings.TrimSpace(req.Email)

	// Find user by email or username alias. The JSON key remains "email" for
	// backwards compatibility with the existing frontend and API clients.
	const loginQuery = `
		SELECT u.id, u.hospital_id, u.email, u.phone, u.password_hash, u.full_name,
		       u.is_active, u.is_verified, u.last_login_at, u.password_reset_at,
		       u.created_at, u.updated_at, u.deleted_at
		FROM users u
		LEFT JOIN user_login_aliases ula
		  ON ula.user_id = u.id
		 AND ula.deleted_at IS NULL
		WHERE u.deleted_at IS NULL
		  AND (
		    lower(u.email) = lower($1)
		    OR lower(ula.username) = lower($1)
		  )
		ORDER BY CASE WHEN lower(ula.username) = lower($1) THEN 0 ELSE 1 END
		LIMIT 1
	`

	var user sqlc.User
	err := h.pool.QueryRow(ctx, loginQuery, identifier).Scan(
		&user.ID,
		&user.HospitalID,
		&user.Email,
		&user.Phone,
		&user.PasswordHash,
		&user.FullName,
		&user.IsActive,
		&user.IsVerified,
		&user.LastLoginAt,
		&user.PasswordResetAt,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to authenticate"})
		return
	}

	// Verify password
	if !security.CheckPassword(user.PasswordHash, req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Check if user is active
	if !user.IsActive {
		c.JSON(http.StatusForbidden, gin.H{"error": "Account is inactive"})
		return
	}

	// Generate JWT token (24 hour expiration)
	token, err := h.jwtSvc.Generate(
		user.HospitalID.String(),
		user.ID.String(),
		user.Email,
		24*time.Hour,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	var hospitalName string
	if err := h.pool.QueryRow(ctx, "SELECT name FROM hospitals WHERE id = $1", user.HospitalID).Scan(&hospitalName); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load hospital"})
		return
	}

	userRoles, err := queries.GetUserRoles(ctx, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load user roles"})
		return
	}

	roleNames := make([]string, 0, len(userRoles))
	isHospitalAdmin := false
	isPlatformAdmin := false
	for _, role := range userRoles {
		roleNames = append(roleNames, role.RoleName)
		if strings.EqualFold(role.RoleName, "super_admin") {
			isPlatformAdmin = true
			isHospitalAdmin = true
		}
		if strings.EqualFold(role.RoleName, "admin") {
			isHospitalAdmin = true
		}
	}

	// Return token and user info
	c.JSON(http.StatusOK, LoginResponse{
		Token: token,
		User: gin.H{
			"id":                user.ID,
			"email":             user.Email,
			"full_name":         user.FullName,
			"hospital_id":       user.HospitalID,
			"hospital_name":     hospitalName,
			"role_names":        roleNames,
			"is_admin":          isHospitalAdmin || isPlatformAdmin,
			"is_hospital_admin": isHospitalAdmin,
			"is_platform_admin": isPlatformAdmin,
		},
	})
}
