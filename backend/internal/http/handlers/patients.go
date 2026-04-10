package handlers

import (
	"net/http"
	"time"

	"github.com/dmsafrica/dms/internal/db/sqlc"
	"github.com/dmsafrica/dms/internal/http/middleware"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PatientsHandler struct {
	pool *pgxpool.Pool
}

func NewPatientsHandler(pool *pgxpool.Pool) *PatientsHandler {
	return &PatientsHandler{pool: pool}
}

// Blood type mapping from display format to database format
var bloodTypeMapping = map[string]string{
	"A+":      "a_positive",
	"A-":      "a_negative",
	"B+":      "b_positive",
	"B-":      "b_negative",
	"AB+":     "ab_positive",
	"AB-":     "ab_negative",
	"O+":      "o_positive",
	"O-":      "o_negative",
	"unknown": "unknown",
	// Also accept the database format
	"a_positive":  "a_positive",
	"a_negative":  "a_negative",
	"b_positive":  "b_positive",
	"b_negative":  "b_negative",
	"ab_positive": "ab_positive",
	"ab_negative": "ab_negative",
	"o_positive":  "o_positive",
	"o_negative":  "o_negative",
}

func mapBloodType(input string) (string, bool) {
	mapped, ok := bloodTypeMapping[input]
	return mapped, ok
}

func (h *PatientsHandler) Create(c *gin.Context) {
	var req struct {
		MRN                string `json:"mrn" binding:"required"`
		NationalID         string `json:"national_id"`
		FullName           string `json:"full_name" binding:"required"`
		PreferredName      string `json:"preferred_name"`
		DateOfBirth        string `json:"date_of_birth" binding:"required"`
		Sex                string `json:"sex" binding:"required"`
		BloodType          string `json:"blood_type"`
		MaritalStatus      string `json:"marital_status"`
		Nationality        string `json:"nationality"`
		Religion           string `json:"religion"`
		Occupation         string `json:"occupation"`
		EducationLevel     string `json:"education_level"`
		PrimaryLanguage    string `json:"primary_language"`
		InterpreterNeeded  bool   `json:"interpreter_needed"`
		PrimaryDoctorID    string `json:"primary_doctor_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate and map blood type
	var mappedBloodType string
	if req.BloodType != "" {
		var ok bool
		mappedBloodType, ok = mapBloodType(req.BloodType)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid blood_type",
				"valid_values": []string{"A+", "A-", "B+", "B-", "AB+", "AB-", "O+", "O-", "unknown"},
			})
			return
		}
	} else {
		mappedBloodType = "unknown"
	}

	hospitalIDStr, _ := c.Get(middleware.CtxHospitalID)
	userIDStr, _ := c.Get(middleware.CtxUserID)
	hospitalID, _ := uuid.Parse(hospitalIDStr.(string))
	userID, _ := uuid.Parse(userIDStr.(string))

	dob, err := time.Parse("2006-01-02", req.DateOfBirth)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date_of_birth format"})
		return
	}

	var primaryDoctorID pgtype.UUID
	if req.PrimaryDoctorID != "" {
		doctorUUID, err := uuid.Parse(req.PrimaryDoctorID)
		if err == nil {
			primaryDoctorID = pgtype.UUID{Bytes: doctorUUID, Valid: true}
		}
	}

	queries := sqlc.New(h.pool)
	patient, err := queries.CreatePatient(c.Request.Context(), sqlc.CreatePatientParams{
		HospitalID:        hospitalID,
		Mrn:               req.MRN,
		NationalID:        pgtype.Text{String: req.NationalID, Valid: req.NationalID != ""},
		FullName:          req.FullName,
		PreferredName:     pgtype.Text{String: req.PreferredName, Valid: req.PreferredName != ""},
		DateOfBirth:       pgtype.Date{Time: dob, Valid: true},
		Sex:               sqlc.SexType(req.Sex),
		BloodType:         sqlc.BloodType(mappedBloodType),
		MaritalStatus:     sqlc.NullMaritalStatus{MaritalStatus: sqlc.MaritalStatus(req.MaritalStatus), Valid: req.MaritalStatus != ""},
		Nationality:       pgtype.Text{String: req.Nationality, Valid: req.Nationality != ""},
		Religion:          pgtype.Text{String: req.Religion, Valid: req.Religion != ""},
		Occupation:        pgtype.Text{String: req.Occupation, Valid: req.Occupation != ""},
		EducationLevel:    pgtype.Text{String: req.EducationLevel, Valid: req.EducationLevel != ""},
		PrimaryLanguage:   pgtype.Text{String: req.PrimaryLanguage, Valid: req.PrimaryLanguage != ""},
		InterpreterNeeded: req.InterpreterNeeded,
		RegistrationDate:  pgtype.Date{Time: time.Now(), Valid: true},
		RegisteredBy:      userID,
		PrimaryDoctorID:   primaryDoctorID,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create patient", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, patient)
}

func (h *PatientsHandler) Get(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid patient ID"})
		return
	}

	queries := sqlc.New(h.pool)
	patient, err := queries.GetPatient(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "patient not found"})
		return
	}

	c.JSON(http.StatusOK, patient)
}

func (h *PatientsHandler) List(c *gin.Context) {
	hospitalIDStr, _ := c.Get(middleware.CtxHospitalID)
	hospitalID, _ := uuid.Parse(hospitalIDStr.(string))

	queries := sqlc.New(h.pool)
	patients, err := queries.ListActivePatients(c.Request.Context(), sqlc.ListActivePatientsParams{
		HospitalID: hospitalID,
		Limit:      100,
		Offset:     0,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list patients"})
		return
	}

	c.JSON(http.StatusOK, patients)
}

func (h *PatientsHandler) Search(c *gin.Context) {
	query := c.Query("q")
	searchType := c.Query("type")

	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query parameter 'q' required"})
		return
	}

	hospitalIDStr, _ := c.Get(middleware.CtxHospitalID)
	hospitalID, _ := uuid.Parse(hospitalIDStr.(string))

	queries := sqlc.New(h.pool)

	switch searchType {
	case "mrn":
		patient, err := queries.GetPatientByMRN(c.Request.Context(), sqlc.GetPatientByMRNParams{
			HospitalID: hospitalID,
			Mrn:        query,
		})
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "patient not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"results": []sqlc.Patient{patient}})

	case "national_id":
		patient, err := queries.GetPatientByNationalID(c.Request.Context(), pgtype.Text{String: query, Valid: true})
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "patient not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"results": []sqlc.Patient{patient}})

	case "name":
		patients, err := queries.SearchPatientsByName(c.Request.Context(), sqlc.SearchPatientsByNameParams{
			HospitalID: hospitalID,
			FullName:   "%" + query + "%",
			Limit:      50,
			Offset:     0,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "search failed"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"results": patients})

	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid search type, use: mrn, national_id, name"})
	}
}

func (h *PatientsHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid patient ID"})
		return
	}

	queries := sqlc.New(h.pool)
	err = queries.SoftDeletePatient(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete patient"})
		return
	}

	c.Status(http.StatusNoContent)
}
