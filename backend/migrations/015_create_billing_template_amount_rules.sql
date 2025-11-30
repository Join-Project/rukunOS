-- Migration: Create Billing Template Amount Rules Table
-- Description: Store amount rules per unit type for billing templates
-- Date: 2025-11-30

CREATE TABLE IF NOT EXISTS billing_template_amount_rules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    template_id UUID NOT NULL REFERENCES billing_templates(id) ON DELETE CASCADE,
    unit_type VARCHAR(50) NOT NULL CHECK (unit_type IN ('rumah', 'ruko', 'kios')),
    amount DECIMAL(15, 2) NOT NULL CHECK (amount >= 0),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(template_id, unit_type)
);

CREATE INDEX IF NOT EXISTS idx_billing_template_amount_rules_template_id ON billing_template_amount_rules(template_id);
CREATE INDEX IF NOT EXISTS idx_billing_template_amount_rules_unit_type ON billing_template_amount_rules(unit_type);

-- Trigger untuk auto-update updated_at
DROP TRIGGER IF EXISTS update_billing_template_amount_rules_updated_at ON billing_template_amount_rules;
CREATE TRIGGER update_billing_template_amount_rules_updated_at 
    BEFORE UPDATE ON billing_template_amount_rules
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();


