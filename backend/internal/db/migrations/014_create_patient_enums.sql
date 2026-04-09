-- +goose Up
CREATE TYPE sex_type AS ENUM ('male','female','intersex','unknown');
CREATE TYPE blood_type AS ENUM ('A+','A-','B+','B-','AB+','AB-','O+','O-','unknown');
CREATE TYPE marital_status AS ENUM ('single','married','divorced','widowed','unknown');
CREATE TYPE severity_level AS ENUM ('mild','moderate','severe','life_threatening');
CREATE TYPE id_type AS ENUM ('national_id','passport','nhif','nhis','refugee_id','hospital_mrn','other');
CREATE TYPE contact_type AS ENUM ('phone','email','address','whatsapp');
CREATE TYPE flag_type AS ENUM ('high_risk','infectious','hiv_positive','hepatitis_b','hepatitis_c','non_compliant','vip','deceased','allergy_alert');

-- +goose Down
DROP TYPE IF EXISTS flag_type;
DROP TYPE IF EXISTS contact_type;
DROP TYPE IF EXISTS id_type;
DROP TYPE IF EXISTS severity_level;
DROP TYPE IF EXISTS marital_status;
DROP TYPE IF EXISTS blood_type;
DROP TYPE IF EXISTS sex_type;
