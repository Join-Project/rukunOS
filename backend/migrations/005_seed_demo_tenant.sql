-- Migration: Seed Demo Tenant
-- Description: Create default DEMO001 tenant for development/testing
-- Date: 2025-01

-- Create DEMO001 tenant if it doesn't exist
INSERT INTO tenants (id, name, code, address, phone, email, status)
VALUES (
    gen_random_uuid(),
    'RukunOS Demo',
    'DEMO001',
    'Jl. Demo No. 1, Jakarta',
    '081234567890',
    'demo@rukunos.id',
    'active'
)
ON CONFLICT (code) DO NOTHING;

-- Create default roles for DEMO001 tenant if they don't exist
DO $$
DECLARE
    v_tenant_id UUID;
BEGIN
    SELECT id INTO v_tenant_id FROM tenants WHERE code = 'DEMO001' AND deleted_at IS NULL;
    
    IF v_tenant_id IS NOT NULL THEN
        -- Check if roles exist, if not create them
        IF NOT EXISTS (SELECT 1 FROM roles WHERE tenant_id = v_tenant_id AND deleted_at IS NULL) THEN
            PERFORM create_default_roles_for_tenant(v_tenant_id);
        END IF;
    END IF;
END $$;










