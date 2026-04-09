-- +goose Up
CREATE TYPE equipment_category AS ENUM (
    'dialysis_machine',
    'ro_unit',
    'water_softener',
    'bp_monitor',
    'weighing_scale',
    'ecg_machine',
    'pulse_oximeter',
    'glucometer',
    'thermometer',
    'wheelchair',
    'stretcher',
    'oxygen_concentrator',
    'defibrillator',
    'other'
);

CREATE TYPE equipment_status AS ENUM (
    'operational',
    'in_use',
    'maintenance',
    'faulty',
    'under_repair',
    'decommissioned',
    'retired'
);

CREATE TYPE maintenance_type AS ENUM (
    'scheduled',
    'corrective',
    'preventive',
    'calibration',
    'inspection',
    'repair',
    'replacement'
);

CREATE TYPE fault_severity AS ENUM (
    'minor',
    'moderate',
    'severe',
    'critical'
);

CREATE TYPE consumable_category AS ENUM (
    'dialyzer',
    'bloodline',
    'av_needle',
    'syringe',
    'gauze',
    'gloves',
    'mask',
    'disinfectant',
    'saline',
    'heparin',
    'other'
);

CREATE TYPE certification_type AS ENUM (
    'regulatory',
    'calibration',
    'safety_inspection',
    'quality_assurance',
    'warranty',
    'insurance',
    'other'
);

-- +goose Down
DROP TYPE IF EXISTS certification_type;
DROP TYPE IF EXISTS consumable_category;
DROP TYPE IF EXISTS fault_severity;
DROP TYPE IF EXISTS maintenance_type;
DROP TYPE IF EXISTS equipment_status;
DROP TYPE IF EXISTS equipment_category;
