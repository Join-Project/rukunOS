-- Migration: Create Units Table (Phase 1)
-- Description: Create units table for unit management
-- Date: 2025-01

CREATE TABLE IF NOT EXISTS units (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    code VARCHAR(50) NOT NULL,
    type VARCHAR(50) NOT NULL,
    owner_name VARCHAR(255),
    owner_phone VARCHAR(50),
    owner_email VARCHAR(255),
    address TEXT,
    status VARCHAR(20) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    UNIQUE(tenant_id, code)
);

CREATE INDEX IF NOT EXISTS idx_units_tenant_id ON units(tenant_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_units_type ON units(type) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_units_code ON units(code) WHERE deleted_at IS NULL;

-- Add foreign key constraint for unit_id in tenant_users
ALTER TABLE tenant_users 
    ADD CONSTRAINT IF NOT EXISTS fk_tenant_users_unit_id 
    FOREIGN KEY (unit_id) REFERENCES units(id) ON DELETE SET NULL;

-- Trigger untuk auto-update updated_at
DROP TRIGGER IF EXISTS update_units_updated_at ON units;
CREATE TRIGGER update_units_updated_at 
    BEFORE UPDATE ON units
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

