-- Migration: Add Bill Number to Bills Table
-- Description: Add bill_number field for auto-generated bill numbers
-- Date: 2025-11-30

-- Add bill_number column
ALTER TABLE bills 
ADD COLUMN IF NOT EXISTS bill_number VARCHAR(100);

-- Create index for bill_number
CREATE INDEX IF NOT EXISTS idx_bills_bill_number ON bills(tenant_id, bill_number) WHERE deleted_at IS NULL;

-- Function to generate bill number
CREATE OR REPLACE FUNCTION generate_bill_number()
RETURNS TRIGGER AS $$
DECLARE
    bill_count INT;
    bill_num VARCHAR(100);
    tenant_code VARCHAR(50);
BEGIN
    -- Get tenant code (first 3 chars of tenant_id or use default)
    SELECT SUBSTRING(t.id::TEXT, 1, 8) INTO tenant_code
    FROM tenants t
    WHERE t.id = NEW.tenant_id;
    
    IF tenant_code IS NULL THEN
        tenant_code := 'TENANT';
    END IF;
    
    -- Count bills for this tenant in current month
    SELECT COUNT(*) + 1 INTO bill_count
    FROM bills
    WHERE tenant_id = NEW.tenant_id
    AND DATE_TRUNC('month', created_at) = DATE_TRUNC('month', CURRENT_TIMESTAMP)
    AND deleted_at IS NULL;
    
    -- Generate bill number: BILL-YYYYMM-TENANT-0001
    bill_num := 'BILL-' || TO_CHAR(CURRENT_TIMESTAMP, 'YYYYMM') || '-' || tenant_code || '-' || LPAD(bill_count::TEXT, 4, '0');
    NEW.bill_number := bill_num;
    
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Trigger to auto-generate bill number on insert
DROP TRIGGER IF EXISTS generate_bill_number_trigger ON bills;
CREATE TRIGGER generate_bill_number_trigger 
    BEFORE INSERT ON bills
    FOR EACH ROW 
    WHEN (NEW.bill_number IS NULL)
    EXECUTE FUNCTION generate_bill_number();




