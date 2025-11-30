-- Migration: Create Announcements Table
-- Description: Create announcements table for communication management
-- Date: 2025-01

CREATE TABLE IF NOT EXISTS announcements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    author_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    priority VARCHAR(20) DEFAULT 'medium', -- low, medium, high
    category VARCHAR(100), -- umum, penting, darurat, dll
    is_pinned BOOLEAN DEFAULT false,
    sent_notification BOOLEAN DEFAULT false,
    sent_whatsapp BOOLEAN DEFAULT false,
    sent_at TIMESTAMP,
    expires_at TIMESTAMP, -- untuk pengumuman dengan expiry
    metadata JSONB DEFAULT '{}', -- untuk extensibility
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_announcements_tenant_id ON announcements(tenant_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_announcements_author_id ON announcements(author_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_announcements_priority ON announcements(priority) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_announcements_created_at ON announcements(created_at DESC) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_announcements_is_pinned ON announcements(is_pinned) WHERE deleted_at IS NULL AND is_pinned = true;

-- Trigger untuk auto-update updated_at
DROP TRIGGER IF EXISTS update_announcements_updated_at ON announcements;
CREATE TRIGGER update_announcements_updated_at 
    BEFORE UPDATE ON announcements
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Table untuk tracking pembacaan pengumuman
CREATE TABLE IF NOT EXISTS announcement_reads (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    announcement_id UUID NOT NULL REFERENCES announcements(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    read_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(announcement_id, user_id)
);

CREATE INDEX IF NOT EXISTS idx_announcement_reads_announcement_id ON announcement_reads(announcement_id);
CREATE INDEX IF NOT EXISTS idx_announcement_reads_user_id ON announcement_reads(user_id);










