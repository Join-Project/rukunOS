-- Migration: Create Complaints Table
-- Description: Create complaints table for resident complaint management
-- Date: 2025-02-02

CREATE TABLE IF NOT EXISTS complaints (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    unit_id UUID REFERENCES units(id) ON DELETE SET NULL,
    category VARCHAR(100) NOT NULL, -- keamanan, kebersihan, fasilitas_umum, lainnya
    priority VARCHAR(20) DEFAULT 'normal', -- normal, penting, darurat
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    status VARCHAR(20) DEFAULT 'pending', -- pending, in_progress, resolved, rejected
    assigned_to UUID REFERENCES users(id) ON DELETE SET NULL,
    resolved_at TIMESTAMP,
    resolution_notes TEXT,
    attachment_urls TEXT[], -- array of attachment URLs (future)
    metadata JSONB DEFAULT '{}', -- untuk extensibility
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_complaints_tenant_id ON complaints(tenant_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_complaints_user_id ON complaints(user_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_complaints_unit_id ON complaints(unit_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_complaints_status ON complaints(status) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_complaints_category ON complaints(category) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_complaints_created_at ON complaints(created_at DESC) WHERE deleted_at IS NULL;









