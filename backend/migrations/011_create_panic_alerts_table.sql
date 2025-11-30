-- Migration: Create Panic Alerts Table
-- Description: Create panic_alerts table for panic button functionality
-- Date: 2025-01

CREATE TABLE IF NOT EXISTS panic_alerts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    unit_id UUID REFERENCES units(id) ON DELETE SET NULL,
    location JSONB, -- {lat, lng} jika available
    status VARCHAR(20) DEFAULT 'active', -- active, responded, resolved
    responded_by UUID REFERENCES users(id) ON DELETE SET NULL,
    responded_at TIMESTAMP,
    resolved_at TIMESTAMP,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_panic_alerts_tenant_id ON panic_alerts(tenant_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_panic_alerts_status ON panic_alerts(status) WHERE deleted_at IS NULL AND status = 'active';
CREATE INDEX IF NOT EXISTS idx_panic_alerts_created_at ON panic_alerts(created_at DESC) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_panic_alerts_user_id ON panic_alerts(user_id) WHERE deleted_at IS NULL;

-- Trigger untuk auto-update updated_at
DROP TRIGGER IF EXISTS update_panic_alerts_updated_at ON panic_alerts;
CREATE TRIGGER update_panic_alerts_updated_at 
    BEFORE UPDATE ON panic_alerts
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();









