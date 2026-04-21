package errors

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorCode identifies the specific error type for machine-readable frontend consumption.
type ErrorCode string

// Validation errors (400)
const (
	ErrInvalidInput    ErrorCode = "INVALID_INPUT"
	ErrInvalidID       ErrorCode = "INVALID_ID"
	ErrInvalidDate     ErrorCode = "INVALID_DATE"
	ErrInvalidStatus   ErrorCode = "INVALID_STATUS"
	ErrMissingField    ErrorCode = "MISSING_FIELD"
)

// Clinical safety errors (409 Conflict — the request conflicts with clinical rules)
const (
	ErrNoActiveAccess         ErrorCode = "NO_ACTIVE_VASCULAR_ACCESS"
	ErrNoDryWeight            ErrorCode = "NO_RECENT_DRY_WEIGHT"
	ErrConflictingSession     ErrorCode = "CONFLICTING_ACTIVE_SESSION"
	ErrMachineUnavailable     ErrorCode = "MACHINE_UNAVAILABLE"
	ErrMachineMaintenanceDue  ErrorCode = "MACHINE_MAINTENANCE_DUE"
	ErrNoVitalsRecorded       ErrorCode = "NO_VITALS_RECORDED"
	ErrAllergyConflict        ErrorCode = "DRUG_ALLERGY_CONFLICT"
	ErrDrugInteraction        ErrorCode = "DRUG_INTERACTION_DETECTED"
	ErrLabCriticalValue       ErrorCode = "LAB_CRITICAL_VALUE"
)

// Scheduling errors (409)
const (
	ErrShiftConflict     ErrorCode = "SHIFT_OVERLAP_CONFLICT"
	ErrLeaveConflict     ErrorCode = "STAFF_ON_LEAVE"
	ErrStaffUnavailable  ErrorCode = "STAFF_UNAVAILABLE"
)

// Inventory errors (409)
const (
	ErrInsufficientStock ErrorCode = "INSUFFICIENT_STOCK"
	ErrLowStock          ErrorCode = "LOW_STOCK_WARNING"
)

// Auth errors (401/403)
const (
	ErrUnauthorized  ErrorCode = "UNAUTHORIZED"
	ErrForbidden     ErrorCode = "FORBIDDEN"
	ErrInactiveUser  ErrorCode = "ACCOUNT_INACTIVE"
	ErrInvalidCreds  ErrorCode = "INVALID_CREDENTIALS"
	ErrTokenExpired  ErrorCode = "TOKEN_EXPIRED"
	ErrRateLimited   ErrorCode = "RATE_LIMITED"
)

// Resource errors (404)
const (
	ErrNotFound ErrorCode = "NOT_FOUND"
)

// Server errors (500)
const (
	ErrInternal    ErrorCode = "INTERNAL_ERROR"
	ErrDatabase    ErrorCode = "DATABASE_ERROR"
	ErrTransaction ErrorCode = "TRANSACTION_ERROR"
)

// ErrorResponse is the standard error JSON returned by all endpoints.
type ErrorResponse struct {
	Error ErrorBody `json:"error"`
}

// ErrorBody contains the structured error details.
type ErrorBody struct {
	Code    ErrorCode   `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// WarningResponse returns a successful response but with attached warnings
// (e.g., drug interaction detected but severity is low).
type WarningResponse struct {
	Data     interface{} `json:"data"`
	Warnings []Warning   `json:"warnings,omitempty"`
}

// Warning represents a non-blocking clinical alert.
type Warning struct {
	Code    ErrorCode   `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// --- Helper functions for common responses ---

// BadRequest sends a 400 response with the given error code and message.
func BadRequest(c *gin.Context, code ErrorCode, message string) {
	c.JSON(http.StatusBadRequest, ErrorResponse{
		Error: ErrorBody{Code: code, Message: message},
	})
}

// BadRequestWithDetails sends a 400 response with additional details.
func BadRequestWithDetails(c *gin.Context, code ErrorCode, message string, details interface{}) {
	c.JSON(http.StatusBadRequest, ErrorResponse{
		Error: ErrorBody{Code: code, Message: message, Details: details},
	})
}

// NotFound sends a 404 response.
func NotFound(c *gin.Context, resource string) {
	c.JSON(http.StatusNotFound, ErrorResponse{
		Error: ErrorBody{Code: ErrNotFound, Message: resource + " not found"},
	})
}

// Conflict sends a 409 response for clinical safety / scheduling / inventory violations.
func Conflict(c *gin.Context, code ErrorCode, message string) {
	c.JSON(http.StatusConflict, ErrorResponse{
		Error: ErrorBody{Code: code, Message: message},
	})
}

// ConflictWithDetails sends a 409 response with additional details.
func ConflictWithDetails(c *gin.Context, code ErrorCode, message string, details interface{}) {
	c.JSON(http.StatusConflict, ErrorResponse{
		Error: ErrorBody{Code: code, Message: message, Details: details},
	})
}

// Unauthorized sends a 401 response.
func Unauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, ErrorResponse{
		Error: ErrorBody{Code: ErrUnauthorized, Message: message},
	})
}

// Forbidden sends a 403 response.
func Forbidden(c *gin.Context, message string) {
	c.JSON(http.StatusForbidden, ErrorResponse{
		Error: ErrorBody{Code: ErrForbidden, Message: message},
	})
}

// Internal sends a 500 response. The message should be user-safe (no internals).
func Internal(c *gin.Context, message string) {
	c.JSON(http.StatusInternalServerError, ErrorResponse{
		Error: ErrorBody{Code: ErrInternal, Message: message},
	})
}

// InternalWithCode sends a 500 response with a specific error code.
func InternalWithCode(c *gin.Context, code ErrorCode, message string) {
	c.JSON(http.StatusInternalServerError, ErrorResponse{
		Error: ErrorBody{Code: code, Message: message},
	})
}

// WithWarnings sends a 200 response with data and attached warnings.
func WithWarnings(c *gin.Context, data interface{}, warnings []Warning) {
	c.JSON(http.StatusOK, WarningResponse{Data: data, Warnings: warnings})
}

// CreatedWithWarnings sends a 201 response with data and attached warnings.
func CreatedWithWarnings(c *gin.Context, data interface{}, warnings []Warning) {
	c.JSON(http.StatusCreated, WarningResponse{Data: data, Warnings: warnings})
}
