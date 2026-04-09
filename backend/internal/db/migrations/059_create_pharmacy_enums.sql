-- +goose Up
CREATE TYPE medication_form AS ENUM ('tablet','capsule','syrup','injection','infusion','powder','cream','ointment','inhaler','drops','patch','suppository');
CREATE TYPE medication_route AS ENUM ('oral','iv','im','sc','topical','inhaled','rectal','sublingual','buccal','transdermal','intraperitoneal');
CREATE TYPE prescription_status AS ENUM ('active','completed','cancelled','discontinued','on_hold');
CREATE TYPE stock_movement_type AS ENUM ('purchase','dispensed','returned','expired','damaged','adjustment','transfer_in','transfer_out');

-- +goose Down
DROP TYPE IF EXISTS stock_movement_type;
DROP TYPE IF EXISTS prescription_status;
DROP TYPE IF EXISTS medication_route;
DROP TYPE IF EXISTS medication_form;
