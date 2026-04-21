package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// AdminOnly ensures the requester has an admin-capable role in the current hospital.
// Allowed roles: super_admin, admin.
func AdminOnly(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		hospitalIDStr := c.GetString(CtxHospitalID)
		userIDStr := c.GetString(CtxUserID)

		hospitalID, err := uuid.Parse(hospitalIDStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid hospital context"})
			return
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid user context"})
			return
		}

		const q = `
			SELECT EXISTS (
				SELECT 1
				FROM user_roles ur
				JOIN roles r ON r.id = ur.role_id
				WHERE ur.hospital_id = $1
				  AND ur.user_id = $2
				  AND ur.deleted_at IS NULL
				  AND r.deleted_at IS NULL
				  AND r.name IN ('super_admin', 'admin')
			)
		`

		var isAdmin bool
		if err := pool.QueryRow(c.Request.Context(), q, hospitalID, userID).Scan(&isAdmin); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "authorization check failed"})
			return
		}

		if !isAdmin {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "admin privileges required"})
			return
		}

		c.Next()
	}
}
