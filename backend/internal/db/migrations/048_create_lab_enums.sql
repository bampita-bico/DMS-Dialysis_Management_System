-- +goose Up
CREATE TYPE lab_priority AS ENUM ('routine','urgent','stat');
CREATE TYPE lab_status AS ENUM ('pending','collected','processing','completed','cancelled','rejected');
CREATE TYPE specimen_type AS ENUM ('whole_blood','serum','plasma','urine','stool','csf','sputum','swab','tissue','other');
CREATE TYPE imaging_modality AS ENUM ('xray','ultrasound','echo','ct','mri','fistulogram','angiogram');
CREATE TYPE result_status AS ENUM ('pending','preliminary','final','corrected','cancelled');

-- +goose Down
DROP TYPE IF EXISTS result_status;
DROP TYPE IF EXISTS imaging_modality;
DROP TYPE IF EXISTS specimen_type;
DROP TYPE IF EXISTS lab_status;
DROP TYPE IF EXISTS lab_priority;
