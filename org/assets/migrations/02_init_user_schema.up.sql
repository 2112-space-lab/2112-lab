-- Enable schema
CREATE SCHEMA IF NOT EXISTS user_schema;

-- Create AuditTrail table
CREATE TABLE user_schema.audit_trails (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP NULL,
    table_name VARCHAR(255) NOT NULL,
    record_id VARCHAR(255) NOT NULL,
    action VARCHAR(50) NOT NULL,
    changes_json JSON,
    performed_by VARCHAR(255) NOT NULL,
    performed_at TIMESTAMP NOT NULL
);