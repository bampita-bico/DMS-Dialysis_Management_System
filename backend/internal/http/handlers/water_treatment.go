package handlers

import (
	"net/http"
	"time"

	"github.com/dmsafrica/dms/internal/db/sqlc"
	"github.com/dmsafrica/dms/internal/db/tenant"
	"github.com/dmsafrica/dms/internal/http/middleware"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type WaterTreatmentHandler struct {
	pool *pgxpool.Pool
}

func NewWaterTreatmentHandler(pool *pgxpool.Pool) *WaterTreatmentHandler {
	return &WaterTreatmentHandler{pool: pool}
}

func (h *WaterTreatmentHandler) LogWaterTest(c *gin.Context) {
	var req struct {
		TestDate             string   `json:"test_date" binding:"required"`
		TestTime             string   `json:"test_time" binding:"required"`
		SampleLocation       string   `json:"sample_location" binding:"required"`
		BacterialCountCfuMl  *float64 `json:"bacterial_count_cfu_ml"`
		EndotoxinLevelEuMl   *float64 `json:"endotoxin_level_eu_ml"`
		ChlorineLevelPpm     *float64 `json:"chlorine_level_ppm"`
		ChloramineLevelPpm   *float64 `json:"chloramine_level_ppm"`
		PHLevel              *float64 `json:"ph_level"`
		ConductivityUsCm     *float64 `json:"conductivity_us_cm"`
		HardnessMgL          *float64 `json:"hardness_mg_l"`
		BacteriaResult       string   `json:"bacteria_result"`
		EndotoxinResult      string   `json:"endotoxin_result"`
		ChlorineResult       string   `json:"chlorine_result"`
		OverallResult        string   `json:"overall_result"`
		OutOfSpecParameters  *string  `json:"out_of_spec_parameters"`
		CorrectiveActionTaken *string `json:"corrective_action_taken"`
		CorrectiveActionBy   *uuid.UUID `json:"corrective_action_by"`
		RetestRequired       bool     `json:"retest_required"`
		RetestDate           *string  `json:"retest_date"`
		SystemsShutDown      bool     `json:"systems_shut_down"`
		Notes                *string  `json:"notes"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hospitalIDStr, _ := c.Get(middleware.CtxHospitalID)
	userIDStr, _ := c.Get(middleware.CtxUserID)
	hospitalID, _ := uuid.Parse(hospitalIDStr.(string))
	testedBy, _ := uuid.Parse(userIDStr.(string))

	ctx := c.Request.Context()
	tx, err := h.pool.Begin(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}
	defer tx.Rollback(ctx)

	if err := tenant.SetLocalHospitalID(ctx, tx, hospitalID.String()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set tenant context"})
		return
	}

	queries := sqlc.New(tx)

	testDate, _ := time.Parse("2006-01-02", req.TestDate)
	testTime, _ := time.Parse("15:04", req.TestTime)

	var retestDate pgtype.Date
	if req.RetestDate != nil {
		t, _ := time.Parse("2006-01-02", *req.RetestDate)
		retestDate = pgtype.Date{Time: t, Valid: true}
	}

	log, err := queries.CreateWaterTreatmentLog(ctx, sqlc.CreateWaterTreatmentLogParams{
		HospitalID:            hospitalID,
		TestDate:              pgtype.Date{Time: testDate, Valid: true},
		TestTime:              pgtype.Time{Microseconds: int64(testTime.Hour()*3600+testTime.Minute()*60) * 1000000, Valid: true},
		TestedBy:              testedBy,
		SampleLocation:        req.SampleLocation,
		BacterialCountCfuMl:   pgtype.Numeric{Int: nil, Exp: 0, NaN: false, InfinityModifier: 0, Valid: req.BacterialCountCfuMl != nil},
		EndotoxinLevelEuMl:    pgtype.Numeric{Int: nil, Exp: 0, NaN: false, InfinityModifier: 0, Valid: req.EndotoxinLevelEuMl != nil},
		ChlorineLevelPpm:      pgtype.Numeric{Int: nil, Exp: 0, NaN: false, InfinityModifier: 0, Valid: req.ChlorineLevelPpm != nil},
		ChloramineLevelPpm:    pgtype.Numeric{Int: nil, Exp: 0, NaN: false, InfinityModifier: 0, Valid: req.ChloramineLevelPpm != nil},
		PhLevel:               pgtype.Numeric{Int: nil, Exp: 0, NaN: false, InfinityModifier: 0, Valid: req.PHLevel != nil},
		ConductivityUsCm:      pgtype.Numeric{Int: nil, Exp: 0, NaN: false, InfinityModifier: 0, Valid: req.ConductivityUsCm != nil},
		HardnessMgL:           pgtype.Numeric{Int: nil, Exp: 0, NaN: false, InfinityModifier: 0, Valid: req.HardnessMgL != nil},
		BacteriaResult:        sqlc.WaterTestResult(req.BacteriaResult),
		EndotoxinResult:       sqlc.WaterTestResult(req.EndotoxinResult),
		ChlorineResult:        sqlc.WaterTestResult(req.ChlorineResult),
		OverallResult:         sqlc.WaterTestResult(req.OverallResult),
		OutOfSpecParameters:   pgtype.Text{String: *req.OutOfSpecParameters, Valid: req.OutOfSpecParameters != nil},
		CorrectiveActionTaken: pgtype.Text{String: *req.CorrectiveActionTaken, Valid: req.CorrectiveActionTaken != nil},
		CorrectiveActionBy: func() pgtype.UUID {
			if req.CorrectiveActionBy != nil {
				return pgtype.UUID{Bytes: *req.CorrectiveActionBy, Valid: true}
			}
			return pgtype.UUID{}
		}(),
		RetestRequired:        req.RetestRequired,
		RetestDate:            retestDate,
		SystemsShutDown:       req.SystemsShutDown,
		ShutdownTime:          pgtype.Timestamptz{Time: time.Time{}, Valid: false},
		ResumedTime:           pgtype.Timestamptz{Time: time.Time{}, Valid: false},
		Notes:                 pgtype.Text{String: *req.Notes, Valid: req.Notes != nil},
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to log water test"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusCreated, log)
}

func (h *WaterTreatmentHandler) ListWaterTests(c *gin.Context) {
	dateStr := c.Query("date")
	if dateStr == "" {
		dateStr = time.Now().Format("2006-01-02")
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
		return
	}

	hospitalIDStr, _ := c.Get(middleware.CtxHospitalID)
	hospitalID, _ := uuid.Parse(hospitalIDStr.(string))

	ctx := c.Request.Context()
	tx, err := h.pool.Begin(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}
	defer tx.Rollback(ctx)

	if err := tenant.SetLocalHospitalID(ctx, tx, hospitalID.String()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set tenant context"})
		return
	}

	queries := sqlc.New(tx)
	tests, err := queries.ListWaterTestsByDate(ctx, sqlc.ListWaterTestsByDateParams{
		HospitalID: hospitalID,
		TestDate:   pgtype.Date{Time: date, Valid: true},
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list water tests"})
		return
	}

	tx.Commit(ctx)
	c.JSON(http.StatusOK, tests)
}

func (h *WaterTreatmentHandler) ListFailedWaterTests(c *gin.Context) {
	hospitalIDStr, _ := c.Get(middleware.CtxHospitalID)
	hospitalID, _ := uuid.Parse(hospitalIDStr.(string))

	ctx := c.Request.Context()
	tx, err := h.pool.Begin(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}
	defer tx.Rollback(ctx)

	if err := tenant.SetLocalHospitalID(ctx, tx, hospitalID.String()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set tenant context"})
		return
	}

	queries := sqlc.New(tx)
	tests, err := queries.ListFailedWaterTests(ctx, hospitalID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list failed water tests"})
		return
	}

	tx.Commit(ctx)
	c.JSON(http.StatusOK, tests)
}
