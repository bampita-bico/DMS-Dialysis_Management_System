package middleware

import (
	"net/http"
	"strings"

	"github.com/dmsafrica/dms/internal/security"
	"github.com/gin-gonic/gin"
)

const CtxHospitalID = "hospital_id"
const CtxUserID = "user_id"
const CtxEmail = "email"

func JWTAuth(jwtSvc *security.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		if h == "" || !strings.HasPrefix(strings.ToLower(h), "bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing bearer token"})
			return
		}

		token := strings.TrimSpace(h[len("Bearer "):])
		claims, err := jwtSvc.Parse(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		// Validate required claims
		if claims.HospitalID == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing hospital_id claim"})
			return
		}

		if claims.UserID == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing user_id claim"})
			return
		}

		// Store claims in context for downstream middleware and handlers
		c.Set(CtxHospitalID, claims.HospitalID)
		c.Set(CtxUserID, claims.UserID)
		c.Set(CtxEmail, claims.Email)
		c.Set("jwt_claims", claims)

		c.Next()
	}
}
