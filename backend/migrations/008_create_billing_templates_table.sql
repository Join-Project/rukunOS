-- Migration: Create Billing Templates Table
-- Description: Create billing_templates table for custom billing templates
-- Date: 2025-01

CREATE TABLE IF NOT EXISTS billing_templates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    category VARCHAR(100) NOT NULL,
    type VARCHAR(50) NOT NULL,
    description TEXT,
    amount DECIMAL(15, 2) NOT NULL,
    late_fee DECIMAL(15, 2) DEFAULT 0,
    is_system BOOLEAN DEFAULT false,
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    UNIQUE(tenant_id, name)
);

CREATE INDEX IF NOT EXISTS idx_billing_templates_tenant_id ON billing_templates(tenant_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_billing_templates_type ON billing_templates(type) WHERE deleted_at IS NULL;

-- Trigger untuk auto-update updated_at
DROP TRIGGER IF EXISTS update_billing_templates_updated_at ON billing_templates;
CREATE TRIGGER update_billing_templates_updated_at 
    BEFORE UPDATE ON billing_templates
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();










