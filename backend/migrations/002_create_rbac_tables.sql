-- Migration: Create RBAC Tables (Phase 1)
-- Description: Create roles, permissions, and role_permissions tables for dynamic RBAC
-- Date: 2025-01

-- 1. Permissions table (global, not tenant-specific)
CREATE TABLE IF NOT EXISTS permissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    key VARCHAR(100) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    module VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 2. Roles table (tenant-specific)
CREATE TABLE IF NOT EXISTS roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    is_system BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    UNIQUE(tenant_id, name)
);

CREATE INDEX IF NOT EXISTS idx_roles_tenant_id ON roles(tenant_id) WHERE deleted_at IS NULL;

-- 3. Role Permissions junction table
CREATE TABLE IF NOT EXISTS role_permissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    role_id UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    permission_id UUID NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(role_id, permission_id)
);

CREATE INDEX IF NOT EXISTS idx_role_permissions_role_id ON role_permissions(role_id);
CREATE INDEX IF NOT EXISTS idx_role_permissions_permission_id ON role_permissions(permission_id);

-- Add foreign key constraint for role_id in tenant_users
DO $$ 
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint 
        WHERE conname = 'fk_tenant_users_role_id'
    ) THEN
        ALTER TABLE tenant_users 
            ADD CONSTRAINT fk_tenant_users_role_id 
            FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE SET NULL;
    END IF;
END $$;

-- Trigger untuk auto-update updated_at
DROP TRIGGER IF EXISTS update_roles_updated_at ON roles;
CREATE TRIGGER update_roles_updated_at 
    BEFORE UPDATE ON roles
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Seed default permissions
INSERT INTO permissions (key, name, description, module) VALUES
('billing.create', 'Create Bill', 'Membuat tagihan baru', 'billing'),
('billing.view', 'View Bill', 'Melihat tagihan', 'billing'),
('billing.update', 'Update Bill', 'Mengubah tagihan', 'billing'),
('billing.delete', 'Delete Bill', 'Menghapus tagihan', 'billing'),
('billing.payment', 'Process Payment', 'Memproses pembayaran', 'billing'),
('user.create', 'Create User', 'Membuat user baru', 'user'),
('user.view', 'View User', 'Melihat data user', 'user'),
('user.update', 'Update User', 'Mengubah data user', 'user'),
('user.delete', 'Delete User', 'Menghapus user', 'user'),
('user.manage', 'Manage User', 'Mengelola user (full access)', 'user'),
('communication.announcement.create', 'Create Announcement', 'Membuat pengumuman', 'communication'),
('communication.announcement.view', 'View Announcement', 'Melihat pengumuman', 'communication'),
('communication.announcement.update', 'Update Announcement', 'Mengubah pengumuman', 'communication'),
('communication.announcement.delete', 'Delete Announcement', 'Menghapus pengumuman', 'communication'),
('tenant.settings', 'Manage Tenant Settings', 'Mengelola pengaturan tenant', 'tenant'),
('role.create', 'Create Role', 'Membuat role baru', 'role'),
('role.view', 'View Role', 'Melihat role', 'role'),
('role.update', 'Update Role', 'Mengubah role', 'role'),
('role.delete', 'Delete Role', 'Menghapus role', 'role'),
('unit.create', 'Create Unit', 'Membuat unit baru', 'unit'),
('unit.view', 'View Unit', 'Melihat data unit', 'unit'),
('unit.update', 'Update Unit', 'Mengubah data unit', 'unit'),
('unit.delete', 'Delete Unit', 'Menghapus unit', 'unit')
ON CONFLICT (key) DO NOTHING;

