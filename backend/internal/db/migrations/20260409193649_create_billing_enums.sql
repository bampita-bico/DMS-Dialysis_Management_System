-- +goose Up
CREATE TYPE account_status AS ENUM (
    'active',
    'suspended',
    'closed',
    'defaulted'
);

CREATE TYPE invoice_status AS ENUM (
    'draft',
    'issued',
    'partially_paid',
    'paid',
    'overdue',
    'cancelled',
    'written_off'
);

CREATE TYPE payment_method AS ENUM (
    'cash',
    'mobile_money',
    'bank_transfer',
    'card',
    'cheque',
    'insurance',
    'waiver',
    'other'
);

CREATE TYPE claim_status AS ENUM (
    'draft',
    'submitted',
    'under_review',
    'approved',
    'partially_approved',
    'rejected',
    'paid'
);

CREATE TYPE payment_plan_status AS ENUM (
    'active',
    'completed',
    'defaulted',
    'cancelled'
);

CREATE TYPE waiver_status AS ENUM (
    'pending',
    'approved',
    'rejected'
);

-- +goose Down
DROP TYPE IF EXISTS waiver_status;
DROP TYPE IF EXISTS payment_plan_status;
DROP TYPE IF EXISTS claim_status;
DROP TYPE IF EXISTS payment_method;
DROP TYPE IF EXISTS invoice_status;
DROP TYPE IF EXISTS account_status;
