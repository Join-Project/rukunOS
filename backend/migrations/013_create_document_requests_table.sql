-- Migration: Create Document Requests Table
-- Description: Create document_requests table for resident document request management
-- Date: 2025-02-02

CREATE TABLE IF NOT EXISTS document_requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    document_type VARCHAR(100) NOT NULL, -- surat_pengantar_ktp, surat_pengantar_kk, surat_keterangan_domisili, surat_izin_keramaian
    purpose TEXT NOT NULL,
    status VARCHAR(20) DEFAULT 'pending', -- pending, approved, rejected, completed
    approved_by UUID REFERENCES users(id) ON DELETE SET NULL,
    approved_at TIMESTAMP,
    rejected_reason TEXT,
    attachment_ids UUID[], -- array of file_attachments.id (future)
    notes TEXT,
    metadata JSONB DEFAULT '{}', -- untuk extensibility
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_document_requests_tenant_id ON document_requests(tenant_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_document_requests_user_id ON document_requests(user_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_document_requests_status ON document_requests(status) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_document_requests_document_type ON document_requests(document_type) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_document_requests_created_at ON document_requests(created_at DESC) WHERE deleted_at IS NULL;







