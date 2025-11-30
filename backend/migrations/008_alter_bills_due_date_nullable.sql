-- Migration: Make due_date nullable in bills table
-- Description: Change due_date column to allow NULL values (optional due date)
-- Date: 2025-01

ALTER TABLE bills ALTER COLUMN due_date DROP NOT NULL;










