-- Migration: Seed Default Roles (Phase 1)
-- Description: Create default system roles for each tenant
-- Note: This will be run after tenant creation via application logic
-- Date: 2025-01

-- This migration is a template for seeding default roles
-- Actual seeding will be done via application code when tenant is created
-- Example roles to create:
-- 1. Admin - Full access (all permissions)
-- 2. Member - Basic access (view own bills, view announcements)

-- Function to create default roles for a tenant
CREATE OR REPLACE FUNCTION create_default_roles_for_tenant(p_tenant_id UUID)
RETURNS VOID AS $$
DECLARE
    v_admin_role_id UUID;
    v_member_role_id UUID;
    v_permission_id UUID;
BEGIN
    -- Create Admin role
    INSERT INTO roles (tenant_id, name, description, is_system)
    VALUES (p_tenant_id, 'Admin', 'Full access to all features', true)
    RETURNING id INTO v_admin_role_id;

    -- Create Member role
    INSERT INTO roles (tenant_id, name, description, is_system)
    VALUES (p_tenant_id, 'Member', 'Basic member access', true)
    RETURNING id INTO v_member_role_id;

    -- Assign all permissions to Admin role
    FOR v_permission_id IN SELECT id FROM permissions
    LOOP
        INSERT INTO role_permissions (role_id, permission_id)
        VALUES (v_admin_role_id, v_permission_id)
        ON CONFLICT DO NOTHING;
    END LOOP;

    -- Assign basic permissions to Member role
    -- View bills, view announcements, view own profile
    FOR v_permission_id IN 
        SELECT id FROM permissions 
        WHERE key IN (
            'billing.view',
            'communication.announcement.view',
            'user.view'
        )
    LOOP
        INSERT INTO role_permissions (role_id, permission_id)
        VALUES (v_member_role_id, v_permission_id)
        ON CONFLICT DO NOTHING;
    END LOOP;
END;
$$ LANGUAGE plpgsql;













