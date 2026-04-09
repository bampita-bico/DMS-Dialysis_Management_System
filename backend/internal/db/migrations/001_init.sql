-- +goose Up
-- DMS base initialization: extensions + tenant context helper.

CREATE EXTENSION IF NOT EXISTS pgcrypto;
CREATE EXTENSION IF NOT EXISTS citext;

-- The tenant context is provided by the API per transaction:
--   SET LOCAL app.hospital_id = '<uuid>'
-- RLS policies will use current_setting('app.hospital_id', true).

CREATE OR REPLACE FUNCTION dms_current_hospital_id() RETURNS uuid LANGUAGE sql STABLE AS $func$ SELECT NULLIF(current_setting('app.hospital_id', true), '')::uuid; $func$;

-- Common helper: update updated_at on write
CREATE OR REPLACE FUNCTION dms_set_updated_at() RETURNS trigger LANGUAGE plpgsql AS $func$ BEGIN NEW.updated_at = now(); RETURN NEW; END; $func$;

-- +goose Down
DROP FUNCTION IF EXISTS dms_set_updated_at();
DROP FUNCTION IF EXISTS dms_current_hospital_id();
DROP EXTENSION IF EXISTS citext;
DROP EXTENSION IF EXISTS pgcrypto;
