-- Migration: Add Super Admin Support and Fix Role System
-- Description: 
-- 1. Add is_super_admin column to users table
-- 2. Update create_default_roles_for_tenant to create correct roles (Warga, Admin, Bendahara, Sekretariat, Satpam)
-- 3. Add security permissions
-- Date: 2025-01

-- 1. Add is_super_admin column to users table
ALTER TABLE users 
ADD COLUMN IF NOT EXISTS is_super_admin BOOLEAN DEFAULT false;

CREATE INDEX IF NOT EXISTS idx_users_is_super_admin ON users(is_super_admin) WHERE is_super_admin = true;

-- 2. Add security permissions
INSERT INTO permissions (key, name, description, module) VALUES
('security.visitor.create', 'Create Visitor', 'Mencatat tamu masuk', 'security'),
('security.visitor.view', 'View Visitor', 'Melihat data tamu', 'security'),
('security.visitor.update', 'Update Visitor', 'Mengubah data tamu', 'security'),
('security.visitor.delete', 'Delete Visitor', 'Menghapus data tamu', 'security'),
('security.alert.view', 'View Alert', 'Melihat panic alerts', 'security'),
('security.alert.respond', 'Respond Alert', 'Merespons panic alerts', 'security')
ON CONFLICT (key) DO NOTHING;

-- 3. Update create_default_roles_for_tenant function
CREATE OR REPLACE FUNCTION create_default_roles_for_tenant(p_tenant_id UUID)
RETURNS VOID AS $$
DECLARE
    v_admin_role_id UUID;
    v_warga_role_id UUID;
    v_bendahara_role_id UUID;
    v_sekretariat_role_id UUID;
    v_satpam_role_id UUID;
    v_permission_id UUID;
BEGIN
    -- Create Admin role (full access)
    INSERT INTO roles (tenant_id, name, description, is_system)
    VALUES (p_tenant_id, 'Admin', 'Administrator tenant dengan akses penuh', true)
    ON CONFLICT (tenant_id, name) DO NOTHING
    RETURNING id INTO v_admin_role_id;
    
    -- Get admin role ID if it already exists
    IF v_admin_role_id IS NULL THEN
        SELECT id INTO v_admin_role_id FROM roles WHERE tenant_id = p_tenant_id AND name = 'Admin' AND deleted_at IS NULL;
    END IF;

    -- Create Warga role (default role untuk registrasi)
    INSERT INTO roles (tenant_id, name, description, is_system)
    VALUES (p_tenant_id, 'Warga', 'Warga dengan akses dasar', true)
    ON CONFLICT (tenant_id, name) DO NOTHING
    RETURNING id INTO v_warga_role_id;
    
    -- Get warga role ID if it already exists
    IF v_warga_role_id IS NULL THEN
        SELECT id INTO v_warga_role_id FROM roles WHERE tenant_id = p_tenant_id AND name = 'Warga' AND deleted_at IS NULL;
    END IF;

    -- Create Bendahara role (optional, untuk keuangan)
    INSERT INTO roles (tenant_id, name, description, is_system)
    VALUES (p_tenant_id, 'Bendahara', 'Pengelola keuangan dan tagihan', true)
    ON CONFLICT (tenant_id, name) DO NOTHING
    RETURNING id INTO v_bendahara_role_id;
    
    IF v_bendahara_role_id IS NULL THEN
        SELECT id INTO v_bendahara_role_id FROM roles WHERE tenant_id = p_tenant_id AND name = 'Bendahara' AND deleted_at IS NULL;
    END IF;

    -- Create Sekretariat role (optional, untuk komunikasi)
    INSERT INTO roles (tenant_id, name, description, is_system)
    VALUES (p_tenant_id, 'Sekretariat', 'Pengelola komunikasi dan pengumuman', true)
    ON CONFLICT (tenant_id, name) DO NOTHING
    RETURNING id INTO v_sekretariat_role_id;
    
    IF v_sekretariat_role_id IS NULL THEN
        SELECT id INTO v_sekretariat_role_id FROM roles WHERE tenant_id = p_tenant_id AND name = 'Sekretariat' AND deleted_at IS NULL;
    END IF;

    -- Create Satpam role (optional, untuk keamanan)
    INSERT INTO roles (tenant_id, name, description, is_system)
    VALUES (p_tenant_id, 'Satpam', 'Petugas keamanan', true)
    ON CONFLICT (tenant_id, name) DO NOTHING
    RETURNING id INTO v_satpam_role_id;
    
    IF v_satpam_role_id IS NULL THEN
        SELECT id INTO v_satpam_role_id FROM roles WHERE tenant_id = p_tenant_id AND name = 'Satpam' AND deleted_at IS NULL;
    END IF;

    -- Assign all permissions to Admin role
    FOR v_permission_id IN SELECT id FROM permissions
    LOOP
        INSERT INTO role_permissions (role_id, permission_id)
        VALUES (v_admin_role_id, v_permission_id)
        ON CONFLICT DO NOTHING;
    END LOOP;

    -- Assign basic permissions to Warga role
    FOR v_permission_id IN 
        SELECT id FROM permissions 
        WHERE key IN (
            'billing.view',
            'communication.announcement.view',
            'user.view'
        )
    LOOP
        INSERT INTO role_permissions (role_id, permission_id)
        VALUES (v_warga_role_id, v_permission_id)
        ON CONFLICT DO NOTHING;
    END LOOP;

    -- Assign billing permissions to Bendahara role
    FOR v_permission_id IN 
        SELECT id FROM permissions 
        WHERE key LIKE 'billing.%' OR key IN ('user.view', 'unit.view')
    LOOP
        INSERT INTO role_permissions (role_id, permission_id)
        VALUES (v_bendahara_role_id, v_permission_id)
        ON CONFLICT DO NOTHING;
    END LOOP;

    -- Assign communication permissions to Sekretariat role
    FOR v_permission_id IN 
        SELECT id FROM permissions 
        WHERE key LIKE 'communication.%' OR key = 'user.view'
    LOOP
        INSERT INTO role_permissions (role_id, permission_id)
        VALUES (v_sekretariat_role_id, v_permission_id)
        ON CONFLICT DO NOTHING;
    END LOOP;

    -- Assign security permissions to Satpam role
    FOR v_permission_id IN 
        SELECT id FROM permissions 
        WHERE key LIKE 'security.%' OR key = 'user.view'
    LOOP
        INSERT INTO role_permissions (role_id, permission_id)
        VALUES (v_satpam_role_id, v_permission_id)
        ON CONFLICT DO NOTHING;
    END LOOP;
END;
$$ LANGUAGE plpgsql;

-- 4. Update existing roles: Rename "Member" to "Warga" if exists
UPDATE roles 
SET name = 'Warga', description = 'Warga dengan akses dasar'
WHERE name = 'Member' AND is_system = true;










