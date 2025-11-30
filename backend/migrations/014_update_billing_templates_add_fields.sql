-- Migration: Update Billing Templates - Add Due Day, Recurring Type, Late Fee Type, Is Active
-- Description: Add fields for recurring bill generation and late fee configuration
-- Date: 2025-11-30

-- Add new columns to billing_templates
ALTER TABLE billing_templates 
ADD COLUMN IF NOT EXISTS due_day INTEGER CHECK (due_day >= 1 AND due_day <= 31),
ADD COLUMN IF NOT EXISTS recurring_type VARCHAR(20) DEFAULT 'one-time' CHECK (recurring_type IN ('monthly', 'yearly', 'one-time')),
ADD COLUMN IF NOT EXISTS late_fee_type VARCHAR(20) DEFAULT 'fixed' CHECK (late_fee_type IN ('fixed', 'percentage')),
ADD COLUMN IF NOT EXISTS late_fee_percentage DECIMAL(5,2) DEFAULT 0 CHECK (late_fee_percentage >= 0 AND late_fee_percentage <= 100),
ADD COLUMN IF NOT EXISTS late_fee_max DECIMAL(15,2) DEFAULT NULL,
ADD COLUMN IF NOT EXISTS is_active BOOLEAN DEFAULT true;

-- Update existing templates to have default values
UPDATE billing_templates 
SET recurring_type = 'one-time' 
WHERE recurring_type IS NULL;

UPDATE billing_templates 
SET late_fee_type = 'fixed' 
WHERE late_fee_type IS NULL;

UPDATE billing_templates 
SET is_active = true 
WHERE is_active IS NULL;

-- Add index for active templates
CREATE INDEX IF NOT EXISTS idx_billing_templates_is_active ON billing_templates(tenant_id, is_active) WHERE deleted_at IS NULL AND is_active = true;




