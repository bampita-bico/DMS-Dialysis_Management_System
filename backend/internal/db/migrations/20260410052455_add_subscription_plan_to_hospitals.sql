-- +goose Up
-- Add subscription_plan column with Basic, Standard, Enterprise tiers
ALTER TABLE hospitals
ADD COLUMN subscription_plan VARCHAR(20) NOT NULL DEFAULT 'standard'
CHECK (subscription_plan IN ('basic', 'standard', 'enterprise'));

-- Add enabled_modules JSONB for feature flags
ALTER TABLE hospitals
ADD COLUMN enabled_modules JSONB NOT NULL DEFAULT '{
  "lab_management": true,
  "full_pharmacy": true,
  "hr_management": true,
  "inventory_tracking": true,
  "advanced_billing": true,
  "offline_sync": false,
  "chw_program": false,
  "imaging_integration": false,
  "outcomes_reporting": true
}'::jsonb;

-- Create index for subscription_plan lookups
CREATE INDEX idx_hospitals_subscription_plan ON hospitals(subscription_plan);

-- Create GIN index for enabled_modules JSONB queries
CREATE INDEX idx_hospitals_enabled_modules ON hospitals USING GIN(enabled_modules);

COMMENT ON COLUMN hospitals.subscription_plan IS 'Subscription tier: basic (48 tables), standard (68 tables), enterprise (93 tables)';
COMMENT ON COLUMN hospitals.enabled_modules IS 'Feature flags for optional modules';

-- +goose Down
DROP INDEX IF EXISTS idx_hospitals_enabled_modules;
DROP INDEX IF EXISTS idx_hospitals_subscription_plan;
ALTER TABLE hospitals DROP COLUMN IF EXISTS enabled_modules;
ALTER TABLE hospitals DROP COLUMN IF EXISTS subscription_plan;
