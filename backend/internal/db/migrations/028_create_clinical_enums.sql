-- +goose Up
CREATE TYPE machine_status AS ENUM ('available','in_use','maintenance','quarantine','decommissioned','offline');
CREATE TYPE access_type AS ENUM ('av_fistula','av_graft','tunnelled_cvc','non_tunnelled_cvc','peritoneal_catheter');
CREATE TYPE access_site AS ENUM ('left_forearm','right_forearm','left_upper_arm','right_upper_arm','left_femoral','right_femoral','left_internal_jugular','right_internal_jugular','left_subclavian','right_subclavian','peritoneal');
CREATE TYPE access_status AS ENUM ('active','thrombosed','infected','stenosed','abandoned','maturation');
CREATE TYPE dialysis_modality AS ENUM ('hd','hdf','pd_capd','pd_apd','crrt','sled');
CREATE TYPE session_status AS ENUM ('scheduled','confirmed','in_progress','completed','cancelled','missed','aborted');
CREATE TYPE shift_type AS ENUM ('morning','afternoon','evening','night');
CREATE TYPE complication_severity AS ENUM ('minor','moderate','severe','life_threatening');
CREATE TYPE anticoag_route AS ENUM ('systemic','regional','none');
CREATE TYPE water_test_result AS ENUM ('pass','fail','borderline','pending');

-- +goose Down
DROP TYPE IF EXISTS water_test_result;
DROP TYPE IF EXISTS anticoag_route;
DROP TYPE IF EXISTS complication_severity;
DROP TYPE IF EXISTS shift_type;
DROP TYPE IF EXISTS session_status;
DROP TYPE IF EXISTS dialysis_modality;
DROP TYPE IF EXISTS access_status;
DROP TYPE IF EXISTS access_site;
DROP TYPE IF EXISTS access_type;
DROP TYPE IF EXISTS machine_status;
