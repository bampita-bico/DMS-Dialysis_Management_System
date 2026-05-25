package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/dmsafrica/dms/internal/db/tenant"
	"github.com/dmsafrica/dms/internal/http/middleware"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ClinicalTrackingHandler struct {
	pool *pgxpool.Pool
}

func NewClinicalTrackingHandler(pool *pgxpool.Pool) *ClinicalTrackingHandler {
	return &ClinicalTrackingHandler{pool: pool}
}

type clinicalEntity struct {
	Table        string
	Fields       map[string]string
	DefaultOrder string
	UpsertField  string
}

var clinicalEntities = map[string]clinicalEntity{
	"patient-clinical-profiles": {
		Table:        "patient_clinical_profiles",
		DefaultOrder: "updated_at DESC",
		UpsertField:  "patient_id",
		Fields: map[string]string{
			"patient_id":                         "uuid",
			"renal_course":                       "text",
			"kidney_disease_cause":               "text",
			"dialysis_indication":                "text",
			"symptom_onset_date":                 "date",
			"diagnosis_date":                     "date",
			"dialysis_start_date":                "date",
			"aki_stage":                          "text",
			"baseline_creatinine_mg_dl":          "numeric",
			"highest_creatinine_mg_dl":           "numeric",
			"urine_output_ml_day":                "numeric",
			"residual_kidney_function":           "text",
			"reversible_cause":                   "text",
			"recovery_plan":                      "text",
			"stop_dialysis_trial_date":           "date",
			"aki_to_ckd_reassessment_due":        "date",
			"malaria_aki_phenotype":              "text",
			"malaria_symptom_onset_date":         "date",
			"fever_onset_date":                   "date",
			"dark_urine_onset_date":              "date",
			"malaria_test_date":                  "date",
			"parasitemia":                        "text",
			"first_antimalarial_dose_at":         "timestamptz",
			"definitive_antimalarial_regimen":    "text",
			"parasite_clearance_date":            "date",
			"pregnancy_related":                  "bool",
			"gestational_age_weeks":              "numeric",
			"expected_delivery_date":             "date",
			"gravida":                            "text",
			"para":                               "text",
			"anc_visits":                         "int",
			"bp_proteinuria_summary":             "text",
			"preeclampsia_severity":              "text",
			"delivery_date":                      "date",
			"fetal_outcome":                      "text",
			"magnesium_sulfate_given":            "bool",
			"pregnancy_antihypertensive_regimen": "text",
			"postpartum_renal_recovery":          "text",
			"maternal_fetal_follow_up":           "text",
			"transplant_referral_status":         "text",
			"modality_plan":                      "text",
			"conservative_care_flag":             "bool",
			"tracking_notes":                     "text",
			"created_by":                         "uuid",
			"updated_by":                         "uuid",
		},
	},
	"session-safety-checks": {
		Table:        "session_safety_checks",
		DefaultOrder: "checked_at DESC",
		Fields: map[string]string{
			"session_id":           "uuid",
			"patient_id":           "uuid",
			"checked_by":           "uuid",
			"checked_at":           "timestamptz",
			"check_status":         "text",
			"risk_score":           "int",
			"checklist":            "jsonb",
			"hard_stop_reasons":    "jsonb",
			"source_summary":       "jsonb",
			"override_required":    "bool",
			"override_reason":      "text",
			"override_approved_by": "uuid",
			"override_approved_at": "timestamptz",
			"audit_log_id":         "uuid",
			"notes":                "text",
		},
	},
	"clinical-alerts": {
		Table:        "clinical_alerts",
		DefaultOrder: "created_at DESC",
		Fields: map[string]string{
			"patient_id":       "uuid",
			"session_id":       "uuid",
			"alert_type":       "text",
			"severity":         "text",
			"title":            "text",
			"triggering_value": "text",
			"threshold":        "text",
			"source_table":     "text",
			"source_id":        "uuid",
			"patient_context":  "jsonb",
			"suggested_action": "text",
			"governance_note":  "text",
			"status":           "text",
			"acknowledged_by":  "uuid",
			"acknowledged_at":  "timestamptz",
			"override_reason":  "text",
			"created_by":       "uuid",
		},
	},
	"treatment-telemetry": {
		Table:        "treatment_telemetry_records",
		DefaultOrder: "recorded_at DESC",
		Fields: map[string]string{
			"session_id":                   "uuid",
			"patient_id":                   "uuid",
			"recorded_by":                  "uuid",
			"recorded_at":                  "timestamptz",
			"minutes_on_dialysis":          "int",
			"blood_flow_actual":            "int",
			"dialysate_flow_actual":        "int",
			"tmp_mmhg":                     "numeric",
			"venous_pressure_mmhg":         "int",
			"arterial_pressure_mmhg":       "int",
			"conductivity_ms_cm":           "numeric",
			"temperature_celsius":          "numeric",
			"blood_volume_processed_l":     "numeric",
			"access_recirculation_percent": "numeric",
			"delivered_minutes":            "int",
			"final_delivered_dose":         "text",
			"early_termination_reason":     "text",
			"alarms":                       "jsonb",
			"interruptions":                "jsonb",
			"notes":                        "text",
		},
	},
	"access-lifecycle-events": {
		Table:        "vascular_access_lifecycle_events",
		DefaultOrder: "event_date DESC, created_at DESC",
		Fields: map[string]string{
			"patient_id":                "uuid",
			"access_id":                 "uuid",
			"event_type":                "text",
			"event_date":                "date",
			"event_time":                "time",
			"operator_id":               "uuid",
			"operator_name":             "text",
			"insertion_attempts":        "int",
			"ultrasound_used":           "bool",
			"side_site":                 "text",
			"catheter_length_cm":        "numeric",
			"tip_position_confirmation": "text",
			"immediate_complications":   "jsonb",
			"long_term_complications":   "jsonb",
			"exit_site_condition":       "text",
			"tunnel_cuff_status":        "text",
			"lock_solution":             "text",
			"dressing_date":             "date",
			"hub_scrub_compliant":       "bool",
			"removal_date":              "date",
			"removal_reason":            "text",
			"replacement_reason":        "text",
			"catheter_free_plan":        "text",
			"culture_result":            "text",
			"antibiotic_course":         "text",
			"notes":                     "text",
			"created_by":                "uuid",
		},
	},
	"infection-surveillance-events": {
		Table:        "infection_surveillance_events",
		DefaultOrder: "event_date DESC, created_at DESC",
		Fields: map[string]string{
			"patient_id":                  "uuid",
			"session_id":                  "uuid",
			"access_id":                   "uuid",
			"event_type":                  "text",
			"event_date":                  "date",
			"event_time":                  "time",
			"iv_antimicrobial_started":    "bool",
			"positive_blood_culture":      "bool",
			"access_pus_redness_swelling": "bool",
			"suspected_source":            "text",
			"organism":                    "text",
			"culture_collected_at":        "timestamptz",
			"antimicrobial_started_at":    "timestamptz",
			"hospitalized":                "bool",
			"death_related":               "bool",
			"recurrence_window_notes":     "text",
			"reported_to_registry":        "bool",
			"reported_by":                 "uuid",
			"notes":                       "text",
		},
	},
	"adequacy-reviews": {
		Table:        "dialysis_adequacy_reviews",
		DefaultOrder: "review_month DESC, created_at DESC",
		Fields: map[string]string{
			"patient_id":                        "uuid",
			"session_id":                        "uuid",
			"review_month":                      "date",
			"reviewed_by":                       "uuid",
			"pre_urea_mg_dl":                    "numeric",
			"post_urea_mg_dl":                   "numeric",
			"urr_percent":                       "numeric",
			"sp_kt_v":                           "numeric",
			"target_kt_v":                       "numeric",
			"residual_kidney_function":          "text",
			"urine_volume_ml_day":               "numeric",
			"missed_treatments":                 "int",
			"shortened_treatments":              "int",
			"interdialytic_weight_gain_kg":      "numeric",
			"normalized_protein_catabolic_rate": "numeric",
			"uf_rate_ml_kg_hr":                  "numeric",
			"adequacy_status":                   "text",
			"doctor_review_required":            "bool",
			"recommendations":                   "text",
		},
	},
	"medication-reconciliation-reviews": {
		Table:        "medication_reconciliation_reviews",
		DefaultOrder: "review_date DESC, created_at DESC",
		Fields: map[string]string{
			"patient_id":            "uuid",
			"session_id":            "uuid",
			"review_type":           "text",
			"review_date":           "date",
			"reviewed_by":           "uuid",
			"medications":           "jsonb",
			"renal_dosing_flags":    "jsonb",
			"dialysis_timing_flags": "jsonb",
			"pregnancy_cautions":    "jsonb",
			"nephrotoxin_flags":     "jsonb",
			"adherence_notes":       "text",
			"regimen_change_reason": "text",
			"recommendations":       "text",
			"status":                "text",
		},
	},
	"patient-reported-events": {
		Table:        "patient_reported_events",
		DefaultOrder: "reported_at DESC",
		Fields: map[string]string{
			"patient_id":            "uuid",
			"session_id":            "uuid",
			"event_type":            "text",
			"reported_at":           "timestamptz",
			"reported_by":           "uuid",
			"symptoms":              "jsonb",
			"education_topics":      "jsonb",
			"teach_back_completed":  "bool",
			"transport_reliability": "text",
			"no_show_reason":        "text",
			"payment_barrier":       "bool",
			"food_insecurity":       "bool",
			"caregiver_support":     "text",
			"sms_whatsapp_consent":  "bool",
			"follow_up_channel":     "text",
			"follow_up_due_at":      "timestamptz",
			"notes":                 "text",
		},
	},
	"staff-attendance-verifications": {
		Table:        "staff_attendance_verifications",
		DefaultOrder: "verification_time DESC",
		Fields: map[string]string{
			"staff_id":               "uuid",
			"user_id":                "uuid",
			"shift_assignment_id":    "uuid",
			"verification_date":      "date",
			"verification_time":      "timestamptz",
			"biometric_method":       "text",
			"biometric_device_id":    "text",
			"biometric_template_ref": "text",
			"verification_result":    "text",
			"station_assignment":     "text",
			"patient_assignment":     "jsonb",
			"handover_accepted":      "bool",
			"competency_status":      "jsonb",
			"exception_reason":       "text",
		},
	},
	"unit-safety-events": {
		Table:        "unit_safety_events",
		DefaultOrder: "event_date DESC, created_at DESC",
		Fields: map[string]string{
			"event_type":        "text",
			"severity":          "text",
			"event_date":        "date",
			"event_time":        "time",
			"affected_sessions": "jsonb",
			"affected_patients": "jsonb",
			"machine_id":        "uuid",
			"equipment_id":      "uuid",
			"water_test_id":     "uuid",
			"consumable_lots":   "jsonb",
			"immediate_action":  "text",
			"root_cause":        "text",
			"corrective_action": "text",
			"closure_due_date":  "date",
			"closed_at":         "timestamptz",
			"reported_by":       "uuid",
			"assigned_to":       "uuid",
			"notes":             "text",
		},
	},
	"interoperability-exports": {
		Table:        "interoperability_exports",
		DefaultOrder: "created_at DESC",
		Fields: map[string]string{
			"patient_id":         "uuid",
			"export_type":        "text",
			"fhir_resource_type": "text",
			"local_table":        "text",
			"local_id":           "uuid",
			"coding_system":      "text",
			"coding_version":     "text",
			"code_value":         "text",
			"ucum_unit":          "text",
			"fhir_payload":       "jsonb",
			"export_status":      "text",
			"exported_by":        "uuid",
			"exported_at":        "timestamptz",
			"error_message":      "text",
		},
	},
	"ontology-relationships": {
		Table:        "ontology_relationships",
		DefaultOrder: "created_at DESC",
		Fields: map[string]string{
			"patient_id":       "uuid",
			"source_type":      "text",
			"source_id":        "uuid",
			"relation_type":    "text",
			"target_type":      "text",
			"target_id":        "uuid",
			"relation_context": "jsonb",
			"provenance":       "text",
			"confidence":       "numeric",
			"neo4j_node_ref":   "text",
			"is_active":        "bool",
			"created_by":       "uuid",
		},
	},
	"dialysate-records": {
		Table:        "dialysate_records",
		DefaultOrder: "recorded_at DESC",
		Fields: map[string]string{
			"session_id":           "uuid",
			"patient_id":           "uuid",
			"recorded_by":          "uuid",
			"recorded_at":          "timestamptz",
			"batch_number":         "text",
			"sodium_meq_l":         "numeric",
			"potassium_meq_l":      "numeric",
			"bicarbonate_meq_l":    "numeric",
			"calcium_meq_l":        "numeric",
			"magnesium_meq_l":      "numeric",
			"chloride_meq_l":       "numeric",
			"glucose_mg_dl":        "numeric",
			"acetate_meq_l":        "numeric",
			"conductivity_ms_cm":   "numeric",
			"ph_level":             "numeric",
			"temperature_celsius":  "numeric",
			"flow_rate_ml_min":     "int",
			"total_volume_liters":  "numeric",
			"composition_verified": "bool",
			"verified_by":          "uuid",
			"deviations_noted":     "text",
			"notes":                "text",
		},
	},
	"vascular-access-assessments": {
		Table:        "vascular_access_assessments",
		DefaultOrder: "assessed_at DESC",
		Fields: map[string]string{
			"access_id":              "uuid",
			"patient_id":             "uuid",
			"session_id":             "uuid",
			"assessed_by":            "uuid",
			"assessed_at":            "timestamptz",
			"has_thrill":             "bool",
			"has_bruit":              "bool",
			"has_redness":            "bool",
			"has_swelling":           "bool",
			"has_discharge":          "bool",
			"has_bleeding":           "bool",
			"has_pain":               "bool",
			"appearance_normal":      "bool",
			"flow_rate_ml_min":       "int",
			"venous_pressure_mmhg":   "int",
			"arterial_pressure_mmhg": "int",
			"recirculation_percent":  "numeric",
			"requires_intervention":  "bool",
			"intervention_type":      "text",
			"intervention_urgency":   "text",
			"notes":                  "text",
		},
	},
	"adequacy-assessments": {
		Table:        "adequacy_assessments",
		DefaultOrder: "assessment_date DESC",
		Fields: map[string]string{
			"patient_id":               "uuid",
			"session_id":               "uuid",
			"assessed_by":              "uuid",
			"assessment_date":          "date",
			"kt_v":                     "numeric",
			"kt_v_method":              "text",
			"urr_percent":              "numeric",
			"pre_bun_mg_dl":            "numeric",
			"post_bun_mg_dl":           "numeric",
			"pre_creatinine_mg_dl":     "numeric",
			"post_creatinine_mg_dl":    "numeric",
			"dialysis_duration_mins":   "int",
			"blood_flow_rate":          "int",
			"dialyzer_clearance":       "int",
			"body_water_volume_liters": "numeric",
			"is_adequate":              "bool",
			"recommendations":          "text",
			"next_assessment_date":     "date",
			"notes":                    "text",
		},
	},
	"dialysis-prescriptions": {
		Table:        "dialysis_prescriptions",
		DefaultOrder: "created_at DESC",
		Fields: map[string]string{
			"session_id":             "uuid",
			"patient_id":             "uuid",
			"prescribed_by":          "uuid",
			"modality":               "text",
			"duration_mins":          "int",
			"target_uf_ml":           "numeric",
			"blood_flow_rate":        "int",
			"dialysate_flow_rate":    "int",
			"membrane_type":          "text",
			"membrane_surface_area":  "numeric",
			"dialysate_temp":         "numeric",
			"conductivity_target":    "numeric",
			"pd_params":              "jsonb",
			"anticoagulant":          "text",
			"anticoag_route":         "text",
			"loading_dose_units":     "numeric",
			"maintenance_dose_units": "numeric",
			"maintenance_rate":       "text",
			"sodium_target":          "numeric",
			"potassium_target":       "numeric",
			"bicarbonate_target":     "numeric",
			"calcium_target":         "numeric",
			"glucose_target":         "numeric",
			"notes":                  "text",
		},
	},
}

func (h *ClinicalTrackingHandler) ListEntity(c *gin.Context) {
	entity, ok := clinicalEntityForRequest(c)
	if !ok {
		return
	}

	hospitalID := c.GetString(middleware.CtxHospitalID)
	ctx := c.Request.Context()
	tx, err := h.pool.Begin(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to begin transaction"})
		return
	}
	defer tx.Rollback(ctx)

	if err := tenant.SetLocalHospitalID(ctx, tx, hospitalID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to set tenant context"})
		return
	}

	args := []any{hospitalID}
	where := []string{"hospital_id = $1::uuid", "deleted_at IS NULL"}
	for _, field := range []string{
		"patient_id", "session_id", "access_id", "staff_id", "user_id", "status", "check_status",
		"alert_type", "severity", "event_type", "review_type", "export_type", "verification_result",
		"reported_to_registry", "is_active",
	} {
		if _, exists := entity.Fields[field]; !exists {
			continue
		}
		value := c.Query(field)
		if value == "" {
			continue
		}
		args = append(args, value)
		where = append(where, fmt.Sprintf("%s = $%d%s", field, len(args), sqlCast(entity.Fields[field])))
	}

	query := fmt.Sprintf(
		"SELECT COALESCE(json_agg(row_to_json(t)), '[]'::json) FROM (SELECT * FROM %s WHERE %s ORDER BY %s LIMIT 300) t",
		entity.Table,
		strings.Join(where, " AND "),
		entity.DefaultOrder,
	)

	var raw []byte
	if err := tx.QueryRow(ctx, query, args...).Scan(&raw); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list clinical tracking records"})
		return
	}
	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.Data(http.StatusOK, "application/json", raw)
}

func (h *ClinicalTrackingHandler) GetEntity(c *gin.Context) {
	entity, ok := clinicalEntityForRequest(c)
	if !ok {
		return
	}

	id := c.Param("id")
	hospitalID := c.GetString(middleware.CtxHospitalID)
	ctx := c.Request.Context()
	tx, err := h.pool.Begin(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to begin transaction"})
		return
	}
	defer tx.Rollback(ctx)

	if err := tenant.SetLocalHospitalID(ctx, tx, hospitalID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to set tenant context"})
		return
	}

	query := fmt.Sprintf(
		"SELECT row_to_json(t) FROM (SELECT * FROM %s WHERE id = $1::uuid AND hospital_id = $2::uuid AND deleted_at IS NULL) t",
		entity.Table,
	)

	var raw []byte
	if err := tx.QueryRow(ctx, query, id, hospitalID).Scan(&raw); err != nil {
		if err == pgx.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "record not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get clinical tracking record"})
		return
	}
	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.Data(http.StatusOK, "application/json", raw)
}

func (h *ClinicalTrackingHandler) CreateEntity(c *gin.Context) {
	entity, ok := clinicalEntityForRequest(c)
	if !ok {
		return
	}

	var input map[string]any
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.injectUserDefaults(entity, input, c.GetString(middleware.CtxUserID), false)
	raw, status, err := h.insertOrUpsertEntity(c, entity, input)
	if err != nil {
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.Data(http.StatusCreated, "application/json", raw)
}

func (h *ClinicalTrackingHandler) UpdateEntity(c *gin.Context) {
	entity, ok := clinicalEntityForRequest(c)
	if !ok {
		return
	}

	var input map[string]any
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h.injectUserDefaults(entity, input, c.GetString(middleware.CtxUserID), true)

	assignments, args, err := buildAssignments(entity, input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if len(assignments) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no supported fields supplied"})
		return
	}

	hospitalID := c.GetString(middleware.CtxHospitalID)
	id := c.Param("id")
	ctx := c.Request.Context()
	tx, err := h.pool.Begin(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to begin transaction"})
		return
	}
	defer tx.Rollback(ctx)

	if err := tenant.SetLocalHospitalID(ctx, tx, hospitalID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to set tenant context"})
		return
	}

	args = append(args, id, hospitalID)
	query := fmt.Sprintf(
		"WITH updated AS (UPDATE %s SET %s WHERE id = $%d::uuid AND hospital_id = $%d::uuid AND deleted_at IS NULL RETURNING *) SELECT row_to_json(updated) FROM updated",
		entity.Table,
		strings.Join(assignments, ", "),
		len(args)-1,
		len(args),
	)

	var raw []byte
	if err := tx.QueryRow(ctx, query, args...).Scan(&raw); err != nil {
		if err == pgx.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "record not found"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to update clinical tracking record"})
		return
	}
	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.Data(http.StatusOK, "application/json", raw)
}

func (h *ClinicalTrackingHandler) DeleteEntity(c *gin.Context) {
	entity, ok := clinicalEntityForRequest(c)
	if !ok {
		return
	}

	id := c.Param("id")
	hospitalID := c.GetString(middleware.CtxHospitalID)
	ctx := c.Request.Context()
	tx, err := h.pool.Begin(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to begin transaction"})
		return
	}
	defer tx.Rollback(ctx)

	if err := tenant.SetLocalHospitalID(ctx, tx, hospitalID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to set tenant context"})
		return
	}

	query := fmt.Sprintf("UPDATE %s SET deleted_at = now() WHERE id = $1::uuid AND hospital_id = $2::uuid", entity.Table)
	tag, err := tx.Exec(ctx, query, id, hospitalID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete clinical tracking record"})
		return
	}
	if tag.RowsAffected() == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "record not found"})
		return
	}
	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *ClinicalTrackingHandler) AcknowledgeClinicalAlert(c *gin.Context) {
	var req struct {
		Status         string `json:"status"`
		OverrideReason string `json:"override_reason"`
	}
	_ = c.ShouldBindJSON(&req)
	if req.Status == "" {
		req.Status = "acknowledged"
	}

	hospitalID := c.GetString(middleware.CtxHospitalID)
	userID := c.GetString(middleware.CtxUserID)
	ctx := c.Request.Context()
	tx, err := h.pool.Begin(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to begin transaction"})
		return
	}
	defer tx.Rollback(ctx)

	if err := tenant.SetLocalHospitalID(ctx, tx, hospitalID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to set tenant context"})
		return
	}

	var raw []byte
	err = tx.QueryRow(ctx, `
		WITH updated AS (
			UPDATE clinical_alerts
			SET status = $2,
			    acknowledged_by = $3::uuid,
			    acknowledged_at = now(),
			    override_reason = NULLIF($4, '')
			WHERE id = $1::uuid
			  AND hospital_id = $5::uuid
			  AND deleted_at IS NULL
			RETURNING *
		)
		SELECT row_to_json(updated) FROM updated
	`, c.Param("id"), req.Status, userID, req.OverrideReason, hospitalID).Scan(&raw)
	if err != nil {
		if err == pgx.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "clinical alert not found"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to acknowledge clinical alert"})
		return
	}
	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.Data(http.StatusOK, "application/json", raw)
}

func (h *ClinicalTrackingHandler) OverrideSafetyCheck(c *gin.Context) {
	var req struct {
		OverrideReason string `json:"override_reason" binding:"required"`
		CheckStatus    string `json:"check_status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.CheckStatus == "" {
		req.CheckStatus = "overridden"
	}

	hospitalID := c.GetString(middleware.CtxHospitalID)
	userID := c.GetString(middleware.CtxUserID)
	ctx := c.Request.Context()
	tx, err := h.pool.Begin(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to begin transaction"})
		return
	}
	defer tx.Rollback(ctx)

	if err := tenant.SetLocalHospitalID(ctx, tx, hospitalID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to set tenant context"})
		return
	}

	var raw []byte
	err = tx.QueryRow(ctx, `
		WITH updated AS (
			UPDATE session_safety_checks
			SET check_status = $2,
			    override_required = FALSE,
			    override_reason = $3,
			    override_approved_by = $4::uuid,
			    override_approved_at = now()
			WHERE id = $1::uuid
			  AND hospital_id = $5::uuid
			  AND deleted_at IS NULL
			RETURNING *
		)
		SELECT row_to_json(updated) FROM updated
	`, c.Param("id"), req.CheckStatus, req.OverrideReason, userID, hospitalID).Scan(&raw)
	if err != nil {
		if err == pgx.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "safety check not found"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to override safety check"})
		return
	}
	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.Data(http.StatusOK, "application/json", raw)
}

func (h *ClinicalTrackingHandler) GetCommandCenterSummary(c *gin.Context) {
	hospitalID := c.GetString(middleware.CtxHospitalID)
	ctx := c.Request.Context()
	tx, err := h.pool.Begin(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to begin transaction"})
		return
	}
	defer tx.Rollback(ctx)

	if err := tenant.SetLocalHospitalID(ctx, tx, hospitalID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to set tenant context"})
		return
	}

	var raw []byte
	err = tx.QueryRow(ctx, `
		SELECT json_build_object(
			'blocked_safety_checks', (
				SELECT count(*) FROM session_safety_checks
				WHERE hospital_id = $1::uuid AND deleted_at IS NULL
				  AND (check_status IN ('blocked', 'failed') OR override_required = TRUE)
			),
			'open_clinical_alerts', (
				SELECT count(*) FROM clinical_alerts
				WHERE hospital_id = $1::uuid AND deleted_at IS NULL AND status IN ('open', 'pending')
			),
			'active_infection_events_30d', (
				SELECT count(*) FROM infection_surveillance_events
				WHERE hospital_id = $1::uuid AND deleted_at IS NULL
				  AND event_date >= CURRENT_DATE - INTERVAL '30 days'
			),
			'adequacy_reviews_needing_doctor', (
				SELECT count(*) FROM dialysis_adequacy_reviews
				WHERE hospital_id = $1::uuid AND deleted_at IS NULL AND doctor_review_required = TRUE
			),
			'aki_reassessments_due', (
				SELECT count(*) FROM patient_clinical_profiles
				WHERE hospital_id = $1::uuid AND deleted_at IS NULL
				  AND aki_to_ckd_reassessment_due IS NOT NULL
				  AND aki_to_ckd_reassessment_due <= CURRENT_DATE
			),
			'unit_safety_events_open', (
				SELECT count(*) FROM unit_safety_events
				WHERE hospital_id = $1::uuid AND deleted_at IS NULL AND closed_at IS NULL
			),
			'staff_attendance_exceptions_today', (
				SELECT count(*) FROM staff_attendance_verifications
				WHERE hospital_id = $1::uuid AND deleted_at IS NULL
				  AND verification_date = CURRENT_DATE
				  AND verification_result NOT IN ('verified', 'accepted')
			)
		)
	`, hospitalID).Scan(&raw)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to build command-center summary"})
		return
	}
	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.Data(http.StatusOK, "application/json", raw)
}

func (h *ClinicalTrackingHandler) GetPatientFHIRSummary(c *gin.Context) {
	hospitalID := c.GetString(middleware.CtxHospitalID)
	patientID := c.Param("id")
	ctx := c.Request.Context()
	tx, err := h.pool.Begin(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to begin transaction"})
		return
	}
	defer tx.Rollback(ctx)

	if err := tenant.SetLocalHospitalID(ctx, tx, hospitalID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to set tenant context"})
		return
	}

	var raw []byte
	err = tx.QueryRow(ctx, `
		SELECT json_build_object(
			'resourceType', 'Bundle',
			'type', 'collection',
			'timestamp', now(),
			'entry', json_build_array(
				json_build_object(
					'resource', json_build_object(
						'resourceType', 'Patient',
						'id', p.id,
						'identifier', json_build_array(json_build_object('system', 'urn:dms:mrn', 'value', p.mrn)),
						'name', json_build_array(json_build_object('text', p.full_name)),
						'gender', p.sex,
						'birthDate', p.date_of_birth
					)
				),
				json_build_object(
					'resource', json_build_object(
						'resourceType', 'Observation',
						'id', concat('renal-course-', p.id),
						'status', 'final',
						'code', json_build_object('text', 'Renal course'),
						'subject', json_build_object('reference', concat('Patient/', p.id)),
						'valueString', cp.renal_course
					)
				),
				json_build_object(
					'resource', json_build_object(
						'resourceType', 'Condition',
						'id', concat('kidney-cause-', p.id),
						'clinicalStatus', json_build_object('text', 'active'),
						'code', json_build_object('text', cp.kidney_disease_cause),
						'subject', json_build_object('reference', concat('Patient/', p.id)),
						'onsetDateTime', cp.diagnosis_date
					)
				)
			)
		)
		FROM patients p
		LEFT JOIN patient_clinical_profiles cp ON cp.patient_id = p.id AND cp.deleted_at IS NULL
		WHERE p.id = $1::uuid AND p.hospital_id = $2::uuid AND p.deleted_at IS NULL
	`, patientID, hospitalID).Scan(&raw)
	if err != nil {
		if err == pgx.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "patient not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to build FHIR summary"})
		return
	}
	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.Data(http.StatusOK, "application/fhir+json", raw)
}

func (h *ClinicalTrackingHandler) insertOrUpsertEntity(c *gin.Context, entity clinicalEntity, input map[string]any) ([]byte, int, error) {
	columns, placeholders, args, err := buildInsertColumns(entity, input, c.GetString(middleware.CtxHospitalID))
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	if len(columns) <= 1 {
		return nil, http.StatusBadRequest, fmt.Errorf("no supported fields supplied")
	}

	query := ""
	if entity.UpsertField != "" {
		assignments := make([]string, 0, len(columns))
		for _, col := range columns {
			if col == "hospital_id" || col == entity.UpsertField || col == "created_by" {
				continue
			}
			assignments = append(assignments, fmt.Sprintf("%s = EXCLUDED.%s", col, col))
		}
		assignments = append(assignments, "updated_at = now()")
		query = fmt.Sprintf(
			"WITH saved AS (INSERT INTO %s (%s) VALUES (%s) ON CONFLICT (%s) DO UPDATE SET %s RETURNING *) SELECT row_to_json(saved) FROM saved",
			entity.Table,
			strings.Join(columns, ", "),
			strings.Join(placeholders, ", "),
			entity.UpsertField,
			strings.Join(assignments, ", "),
		)
	} else {
		query = fmt.Sprintf(
			"WITH saved AS (INSERT INTO %s (%s) VALUES (%s) RETURNING *) SELECT row_to_json(saved) FROM saved",
			entity.Table,
			strings.Join(columns, ", "),
			strings.Join(placeholders, ", "),
		)
	}

	hospitalID := c.GetString(middleware.CtxHospitalID)
	ctx := c.Request.Context()
	tx, err := h.pool.Begin(ctx)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("failed to begin transaction")
	}
	defer tx.Rollback(ctx)

	if err := tenant.SetLocalHospitalID(ctx, tx, hospitalID); err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("failed to set tenant context")
	}

	var raw []byte
	if err := tx.QueryRow(ctx, query, args...).Scan(&raw); err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("failed to create clinical tracking record")
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("failed to commit transaction")
	}
	return raw, http.StatusCreated, nil
}

func (h *ClinicalTrackingHandler) injectUserDefaults(entity clinicalEntity, input map[string]any, userID string, updateOnly bool) {
	if userID == "" {
		return
	}
	defaultFields := []string{"recorded_by", "checked_by", "assessed_by", "reviewed_by", "reported_by", "exported_by", "prescribed_by"}
	if updateOnly {
		defaultFields = append(defaultFields, "updated_by")
	} else {
		defaultFields = append(defaultFields, "created_by", "updated_by")
	}
	for _, field := range defaultFields {
		if _, exists := entity.Fields[field]; !exists {
			continue
		}
		if _, provided := input[field]; !provided {
			input[field] = userID
		}
	}
}

func clinicalEntityForRequest(c *gin.Context) (clinicalEntity, bool) {
	name := c.Param("entity")
	entity, ok := clinicalEntities[name]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "unknown clinical tracking entity"})
		return clinicalEntity{}, false
	}
	return entity, true
}

func buildInsertColumns(entity clinicalEntity, input map[string]any, hospitalID string) ([]string, []string, []any, error) {
	columns := []string{"hospital_id"}
	placeholders := []string{"$1::uuid"}
	args := []any{hospitalID}

	fields := sortedInputFields(entity, input)
	for _, field := range fields {
		fieldType := entity.Fields[field]
		value, ok, err := normalizeClinicalValue(input[field], fieldType)
		if err != nil {
			return nil, nil, nil, err
		}
		if !ok {
			continue
		}
		args = append(args, value)
		columns = append(columns, field)
		placeholders = append(placeholders, fmt.Sprintf("$%d%s", len(args), sqlCast(fieldType)))
	}
	return columns, placeholders, args, nil
}

func buildAssignments(entity clinicalEntity, input map[string]any) ([]string, []any, error) {
	assignments := []string{}
	args := []any{}
	fields := sortedInputFields(entity, input)
	for _, field := range fields {
		fieldType := entity.Fields[field]
		value, ok, err := normalizeClinicalValue(input[field], fieldType)
		if err != nil {
			return nil, nil, err
		}
		if !ok {
			continue
		}
		args = append(args, value)
		assignments = append(assignments, fmt.Sprintf("%s = $%d%s", field, len(args), sqlCast(fieldType)))
	}
	return assignments, args, nil
}

func sortedInputFields(entity clinicalEntity, input map[string]any) []string {
	fields := make([]string, 0, len(input))
	for field := range input {
		if _, ok := entity.Fields[field]; ok {
			fields = append(fields, field)
		}
	}
	sort.Strings(fields)
	return fields
}

func normalizeClinicalValue(value any, fieldType string) (any, bool, error) {
	if value == nil {
		return nil, false, nil
	}
	if text, ok := value.(string); ok && strings.TrimSpace(text) == "" {
		return nil, false, nil
	}
	if fieldType == "jsonb" {
		bytes, err := json.Marshal(value)
		if err != nil {
			return nil, false, fmt.Errorf("invalid json field")
		}
		return string(bytes), true, nil
	}
	return value, true, nil
}

func sqlCast(fieldType string) string {
	switch fieldType {
	case "uuid":
		return "::uuid"
	case "date":
		return "::date"
	case "time":
		return "::time"
	case "timestamptz":
		return "::timestamptz"
	case "int":
		return "::integer"
	case "numeric":
		return "::numeric"
	case "bool":
		return "::boolean"
	case "jsonb":
		return "::jsonb"
	default:
		return ""
	}
}
