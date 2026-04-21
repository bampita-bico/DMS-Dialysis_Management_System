package handlers

import (
	"net/http"
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
	Email    string `json:"email" binding:"required,email"`
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

	// Find user by email
	user, err := queries.GetUserForLogin(ctx, req.Email)
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

	// Return token and user info
	c.JSON(http.StatusOK, LoginResponse{
		Token: token,
		User: gin.H{
			"id":          user.ID,
			"email":       user.Email,
			"full_name":   user.FullName,
			"hospital_id": user.HospitalID,
		},
	})
}
