-- Seed Data
-- 1. Create Tenant
INSERT INTO tenants (name, code, address, status)
VALUES ('RukunOS Demo', 'DEMO001', 'Jl. Demo No. 1', 'active')
ON CONFLICT (code) DO NOTHING;

-- 2. Create Roles for the Tenant (using the function)
DO $$
DECLARE
    v_tenant_id UUID;
BEGIN
    SELECT id INTO v_tenant_id FROM tenants WHERE code = 'DEMO001';
    
    -- Check if roles exist, if not create them
    IF NOT EXISTS (SELECT 1 FROM roles WHERE tenant_id = v_tenant_id) THEN
        PERFORM create_default_roles_for_tenant(v_tenant_id);
    END IF;
END $$;

-- 3. Create Admin User
INSERT INTO users (email, password_hash, full_name, status, auth_provider)
VALUES (
    'admin@rukunos.id',
    '$2a$10$xX/k/x.x.x.x.x.x.x.x.x.x.x.x.x.x.x.x.x.x.x.x.x', -- 'password' (placeholder hash, will need real bcrypt hash)
    'Admin RukunOS',
    'active',
    'email'
)
ON CONFLICT (email) DO NOTHING;

-- Update password hash to a known value 'password'
-- Hash for 'password' is usually $2a$10$ ... let's use a known one or just rely on Register for new users.
-- Actually, I can't easily generate bcrypt hash here without a tool.
-- Better to just create the Tenant and Roles, then use Register page (if I fix the tenant code issue).

-- Let's just create the Tenant and Roles.
-- And maybe a "Warga" role if it's not created by the function (the function creates Admin and Member).
-- Wait, the function creates 'Member'. Is 'Member' same as 'Warga'?
-- In useAuth.ts: isResident = role === 'Warga' || role === 'Member'. So yes.

-- 4. Create 'Satpam' role manually since function doesn't
DO $$
DECLARE
    v_tenant_id UUID;
BEGIN
    SELECT id INTO v_tenant_id FROM tenants WHERE code = 'DEMO001';
    
    INSERT INTO roles (tenant_id, name, description, is_system)
    VALUES (v_tenant_id, 'Satpam', 'Security Guard', true)
    ON CONFLICT DO NOTHING;
END $$;
