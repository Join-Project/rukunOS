-- Migration: Create Core Tables (Phase 1)
-- Description: Create tenants, users, and tenant_users tables for multi-tenant architecture
-- Date: 2025-01

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- 1. Tenants table
CREATE TABLE IF NOT EXISTS tenants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    code VARCHAR(50) UNIQUE NOT NULL,
    address TEXT,
    phone VARCHAR(50),
    email VARCHAR(255),
    settings JSONB DEFAULT '{}',
    modules JSONB DEFAULT '{
        "billing": true,
        "communication": true,
        "administration": false,
        "security": false
    }',
    status VARCHAR(20) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_tenants_code ON tenants(code) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_tenants_status ON tenants(status) WHERE deleted_at IS NULL;

-- 2. Users table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255),
    full_name VARCHAR(255) NOT NULL,
    phone VARCHAR(50),
    google_id VARCHAR(255),
    auth_provider VARCHAR(50) DEFAULT 'email',
    avatar_url TEXT,
    status VARCHAR(20) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email_auth_provider 
    ON users (email, auth_provider) 
    WHERE deleted_at IS NULL AND auth_provider IN ('email', 'both');

CREATE UNIQUE INDEX IF NOT EXISTS idx_users_google_id 
    ON users (google_id) 
    WHERE deleted_at IS NULL AND google_id IS NOT NULL;

-- 3. Tenant Users junction table
CREATE TABLE IF NOT EXISTS tenant_users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role_id UUID, -- Will reference roles table (created in next migration)
    unit_id UUID, -- Will reference units table (created later)
    status VARCHAR(20) DEFAULT 'active',
    joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    UNIQUE(tenant_id, user_id)
);

CREATE INDEX IF NOT EXISTS idx_tenant_users_tenant_id ON tenant_users(tenant_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_tenant_users_user_id ON tenant_users(user_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_tenant_users_role_id ON tenant_users(role_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_tenant_users_unit_id ON tenant_users(unit_id) WHERE deleted_at IS NULL;

-- 4. Schema migrations tracking table
CREATE TABLE IF NOT EXISTS schema_migrations (
    version VARCHAR(255) PRIMARY KEY,
    description TEXT,
    applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Trigger untuk auto-update updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

DROP TRIGGER IF EXISTS update_tenants_updated_at ON tenants;
CREATE TRIGGER update_tenants_updated_at 
    BEFORE UPDATE ON tenants
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS update_users_updated_at ON users;
CREATE TRIGGER update_users_updated_at 
    BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS update_tenant_users_updated_at ON tenant_users;
CREATE TRIGGER update_tenant_users_updated_at 
    BEFORE UPDATE ON tenant_users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

