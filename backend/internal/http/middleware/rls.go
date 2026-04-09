package middleware

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

// RLSMiddleware sets the PostgreSQL session variables for Row Level Security
// This MUST be called after JWTAuth middleware so hospital_id and user_id are in context
func RLSMiddleware(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get hospital_id and user_id from JWT claims (set by JWTAuth middleware)
		hospitalID, exists := c.Get(CtxHospitalID)
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing hospital_id in context"})
			return
		}

		userID, exists := c.Get(CtxUserID)
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing user_id in context"})
			return
		}

		// Acquire a connection from the pool
		conn, err := pool.Acquire(c.Request.Context())
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "database connection failed"})
			return
		}
		defer conn.Release()

		// Set PostgreSQL session variables for RLS
		// These are used by RLS policies: current_setting('app.current_hospital_id')::UUID
		_, err = conn.Exec(
			context.Background(),
			`SELECT
				set_config('app.current_hospital_id', $1, false),
				set_config('app.current_user_id', $2, false)`,
			hospitalID,
			userID,
		)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to set RLS context"})
			return
		}

		// Store the connection in context so handlers can use it
		c.Set("db_conn", conn)

		c.Next()
	}
}
