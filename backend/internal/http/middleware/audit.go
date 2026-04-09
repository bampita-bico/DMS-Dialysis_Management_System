package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// AuditMiddleware logs all state-changing requests to audit_logs table
// This should be placed after JWTAuth and RLS middleware
func AuditMiddleware(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Only audit state-changing methods
		method := c.Request.Method
		if method != http.MethodPost && method != http.MethodPut &&
			method != http.MethodPatch && method != http.MethodDelete {
			c.Next()
			return
		}

		// Get context values
		hospitalID, _ := c.Get(CtxHospitalID)
		userID, _ := c.Get(CtxUserID)
		requestID, _ := c.Get("request_id")

		// Read request body
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			// Restore body for handlers
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// Determine action from method
		action := mapMethodToAction(method)

		// Extract table name and record ID from path (simplified heuristic)
		tableName, recordID := extractPathInfo(c.Request.URL.Path)

		// Continue processing
		c.Next()

		// After request completes, write audit log asynchronously
		go func() {
			ctx := c.Request.Context()

			var newData map[string]interface{}
			if len(requestBody) > 0 {
				_ = json.Unmarshal(requestBody, &newData)
			}

			// Insert audit log
			_, err := pool.Exec(
				ctx,
				`INSERT INTO audit_logs (
					hospital_id, user_id, action, table_name, record_id,
					new_data, ip_address, user_agent, created_at
				) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
				hospitalID,
				userID,
				action,
				tableName,
				recordID,
				newData,
				c.ClientIP(),
				c.Request.UserAgent(),
				time.Now(),
			)

			if err != nil {
				// Log error but don't fail the request
				// In production, send to error tracking service
				println("audit log failed:", err.Error())
			}

			if reqID, ok := requestID.(string); ok {
				println("audit logged for request:", reqID)
			}
		}()
	}
}

func mapMethodToAction(method string) string {
	switch method {
	case http.MethodPost:
		return "CREATE"
	case http.MethodPut, http.MethodPatch:
		return "UPDATE"
	case http.MethodDelete:
		return "DELETE"
	default:
		return "UNKNOWN"
	}
}

func extractPathInfo(path string) (tableName string, recordID *uuid.UUID) {
	// Simple heuristic: /api/v1/{table}/{id}
	// In production, use proper routing introspection
	// For now, just extract from path segments

	// This is a placeholder - in real implementation,
	// handlers should set this explicitly via context
	return "unknown", nil
}
