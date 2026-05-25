package handlers

import (
	"net/http"

	"github.com/dmsafrica/dms/internal/db/sqlc"
	"github.com/dmsafrica/dms/internal/http/middleware"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RolesHandler struct {
	pool *pgxpool.Pool
}

func NewRolesHandler(pool *pgxpool.Pool) *RolesHandler {
	return &RolesHandler{pool: pool}
}

func (h *RolesHandler) List(c *gin.Context) {
	hospitalIDStr := c.Query("hospital_id")
	if hospitalIDStr == "" {
		hospitalIDStr = c.GetString(middleware.CtxHospitalID)
	}

	hospitalID, err := uuid.Parse(hospitalIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid hospital ID"})
		return
	}

	if err := ensureDefaultRoles(c.Request.Context(), h.pool, hospitalID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to prepare roles"})
		return
	}

	queries := sqlc.New(h.pool)
	roles, err := queries.ListRoles(c.Request.Context(), hospitalID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list roles"})
		return
	}

	c.JSON(http.StatusOK, roles)
}

func (h *RolesHandler) ListUserRoles(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	queries := sqlc.New(h.pool)
	roles, err := queries.GetUserRoles(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list user roles"})
		return
	}

	c.JSON(http.StatusOK, roles)
}

func (h *RolesHandler) AssignUserRole(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	var req struct {
		HospitalID string `json:"hospital_id"`
		RoleID     string `json:"role_id"`
		RoleName   string `json:"role_name"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	queries := sqlc.New(h.pool)
	user, err := queries.GetUser(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	targetHospitalID := user.HospitalID
	if req.HospitalID != "" {
		hospitalID, err := uuid.Parse(req.HospitalID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid hospital ID"})
			return
		}
		if hospitalID != user.HospitalID {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user hospital does not match requested hospital"})
			return
		}
		targetHospitalID = hospitalID
	}

	if err := ensureDefaultRoles(c.Request.Context(), h.pool, targetHospitalID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to prepare roles"})
		return
	}

	var role sqlc.Role
	switch {
	case req.RoleID != "":
		roleID, err := uuid.Parse(req.RoleID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role ID"})
			return
		}
		role, err = queries.GetRole(c.Request.Context(), roleID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "role not found"})
			return
		}
	case req.RoleName != "":
		role, err = findRoleByName(c.Request.Context(), queries, targetHospitalID, req.RoleName)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "role_id or role_name is required"})
		return
	}

	assignerID := c.GetString(middleware.CtxUserID)
	assignerUUID, err := uuid.Parse(assignerID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid requester"})
		return
	}

	userRole, err := queries.AssignUserRole(c.Request.Context(), sqlc.AssignUserRoleParams{
		HospitalID:   targetHospitalID,
		UserID:       userID,
		RoleID:       role.ID,
		DepartmentID: pgtype.UUID{},
		AssignedBy:   pgtype.UUID{Bytes: assignerUUID, Valid: true},
		ExpiresAt:    pgtype.Timestamptz{},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to assign role"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"user_role": userRole,
		"role":      role,
	})
}

func (h *RolesHandler) RevokeUserRole(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	roleID, err := uuid.Parse(c.Param("role_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role ID"})
		return
	}

	result, err := h.pool.Exec(
		c.Request.Context(),
		`DELETE FROM user_roles
		 WHERE user_id = $1
		   AND role_id = $2
		   AND deleted_at IS NULL`,
		userID,
		roleID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to revoke role"})
		return
	}

	if result.RowsAffected() == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "role assignment not found"})
		return
	}

	c.Status(http.StatusNoContent)
}
