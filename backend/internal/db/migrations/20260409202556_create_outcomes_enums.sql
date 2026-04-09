-- +goose Up
CREATE TYPE outcome_trend AS ENUM (
    'improving',
    'stable',
    'declining',
    'critical'
);

CREATE TYPE death_setting AS ENUM (
    'during_dialysis',
    'hospital',
    'home',
    'transit',
    'other'
);

CREATE TYPE hospitalization_outcome AS ENUM (
    'discharged',
    'transferred',
    'deceased',
    'absconded'
);

CREATE TYPE report_type AS ENUM (
    'monthly',
    'quarterly',
    'annual',
    'donor',
    'moh',
    'regulatory'
);

-- +goose Down
DROP TYPE IF EXISTS report_type;
DROP TYPE IF EXISTS hospitalization_outcome;
DROP TYPE IF EXISTS death_setting;
DROP TYPE IF EXISTS outcome_trend;
