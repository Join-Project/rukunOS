-- Migration: Create Visitor Logs Table
-- Description: Create visitor_logs table for visitor management
-- Date: 2025-01

CREATE TABLE IF NOT EXISTS visitor_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    unit_id UUID REFERENCES units(id) ON DELETE SET NULL,
    visitor_name VARCHAR(255) NOT NULL,
    visitor_phone VARCHAR(20),
    visitor_id_number VARCHAR(50), -- KTP/SIM
    visitor_vehicle VARCHAR(100), -- Plat nomor kendaraan
    purpose TEXT, -- Tujuan kunjungan
    host_name VARCHAR(255), -- Nama yang dikunjungi
    checked_in_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    checked_out_at TIMESTAMP,
    notes TEXT,
    checked_in_by UUID REFERENCES users(id) ON DELETE SET NULL, -- Security/admin yang mencatat
    checked_out_by UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_visitor_logs_tenant_id ON visitor_logs(tenant_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_visitor_logs_unit_id ON visitor_logs(unit_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_visitor_logs_checked_in_at ON visitor_logs(checked_in_at DESC) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_visitor_logs_checked_out_at ON visitor_logs(checked_out_at) WHERE checked_out_at IS NULL AND deleted_at IS NULL;

-- Trigger untuk auto-update updated_at
DROP TRIGGER IF EXISTS update_visitor_logs_updated_at ON visitor_logs;
CREATE TRIGGER update_visitor_logs_updated_at 
    BEFORE UPDATE ON visitor_logs
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();










