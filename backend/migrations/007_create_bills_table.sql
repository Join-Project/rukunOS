-- Migration: Create Bills Table
-- Description: Create bills table for billing management
-- Date: 2025-01

-- Create function for auto-update updated_at if not exists
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TABLE IF NOT EXISTS bills (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    unit_id UUID NOT NULL REFERENCES units(id) ON DELETE CASCADE,
    category VARCHAR(100) NOT NULL,
    period VARCHAR(50) NOT NULL,
    amount DECIMAL(15, 2) NOT NULL,
    late_fee DECIMAL(15, 2) DEFAULT 0,
    due_date DATE NOT NULL,
    status VARCHAR(20) DEFAULT 'pending',
    paid_at TIMESTAMP,
    payment_method VARCHAR(50),
    payment_reference VARCHAR(255),
    notes TEXT,
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_bills_tenant_id ON bills(tenant_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_bills_unit_id ON bills(unit_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_bills_status ON bills(status) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_bills_due_date ON bills(due_date) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_bills_period ON bills(period) WHERE deleted_at IS NULL;

-- Trigger untuk auto-update updated_at
DROP TRIGGER IF EXISTS update_bills_updated_at ON bills;
CREATE TRIGGER update_bills_updated_at 
    BEFORE UPDATE ON bills
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

