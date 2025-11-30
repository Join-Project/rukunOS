package models

import (
	"database/sql"
	"time"
)

type BillingTemplate struct {
	ID                 string         `json:"id" db:"id"`
	TenantID           string         `json:"tenant_id" db:"tenant_id"`
	Name               string         `json:"name" db:"name"`
	Category           string         `json:"category" db:"category"`
	Type               string         `json:"type" db:"type"`
	Description        sql.NullString `json:"description,omitempty" db:"description"`
	Amount             float64        `json:"amount" db:"amount"`
	LateFee            float64        `json:"late_fee" db:"late_fee"`
	DueDay             sql.NullInt64  `json:"due_day,omitempty" db:"due_day"`
	RecurringType      string         `json:"recurring_type" db:"recurring_type"`
	LateFeeType        string         `json:"late_fee_type" db:"late_fee_type"`
	LateFeePercentage  sql.NullFloat64 `json:"late_fee_percentage,omitempty" db:"late_fee_percentage"`
	LateFeeMax         sql.NullFloat64 `json:"late_fee_max,omitempty" db:"late_fee_max"`
	IsActive           bool           `json:"is_active" db:"is_active"`
	IsSystem           bool           `json:"is_system" db:"is_system"`
	CreatedBy          sql.NullString `json:"created_by,omitempty" db:"created_by"`
	CreatedAt          time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at" db:"updated_at"`
	DeletedAt          sql.NullTime   `json:"-" db:"deleted_at"`
}

type BillingTemplateAmountRule struct {
	ID         string    `json:"id" db:"id"`
	TemplateID string    `json:"template_id" db:"template_id"`
	UnitType   string    `json:"unit_type" db:"unit_type"`
	Amount     float64   `json:"amount" db:"amount"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

type CreateBillingTemplateRequest struct {
	Name               string   `json:"name" validate:"required"`
	Category           string   `json:"category" validate:"required"`
	Type               string   `json:"type" validate:"required,oneof=Bulanan Tahunan One-time"`
	Description        *string  `json:"description,omitempty"`
	Amount             float64  `json:"amount" validate:"required,min=0"`
	LateFee            *float64 `json:"late_fee,omitempty"`
	DueDay             *int     `json:"due_day,omitempty" validate:"omitempty,min=1,max=31"`
	RecurringType      *string  `json:"recurring_type,omitempty" validate:"omitempty,oneof=monthly yearly one-time"`
	LateFeeType        *string  `json:"late_fee_type,omitempty" validate:"omitempty,oneof=fixed percentage"`
	LateFeePercentage  *float64 `json:"late_fee_percentage,omitempty" validate:"omitempty,min=0,max=100"`
	LateFeeMax         *float64 `json:"late_fee_max,omitempty" validate:"omitempty,min=0"`
	IsActive           *bool    `json:"is_active,omitempty"`
	AmountRules        []AmountRuleRequest `json:"amount_rules,omitempty"`
}

type AmountRuleRequest struct {
	UnitType string  `json:"unit_type" validate:"required,oneof=rumah ruko kios"`
	Amount   float64 `json:"amount" validate:"required,min=0"`
}

type UpdateBillingTemplateRequest struct {
	Name               *string  `json:"name,omitempty"`
	Category           *string  `json:"category,omitempty"`
	Type               *string  `json:"type,omitempty" validate:"omitempty,oneof=Bulanan Tahunan One-time"`
	Description        *string  `json:"description,omitempty"`
	Amount             *float64 `json:"amount,omitempty" validate:"omitempty,min=0"`
	LateFee            *float64 `json:"late_fee,omitempty"`
	DueDay             *int     `json:"due_day,omitempty" validate:"omitempty,min=1,max=31"`
	RecurringType      *string  `json:"recurring_type,omitempty" validate:"omitempty,oneof=monthly yearly one-time"`
	LateFeeType        *string  `json:"late_fee_type,omitempty" validate:"omitempty,oneof=fixed percentage"`
	LateFeePercentage  *float64 `json:"late_fee_percentage,omitempty" validate:"omitempty,min=0,max=100"`
	LateFeeMax         *float64 `json:"late_fee_max,omitempty" validate:"omitempty,min=0"`
	IsActive           *bool    `json:"is_active,omitempty"`
	AmountRules        *[]AmountRuleRequest `json:"amount_rules,omitempty"`
}

type GenerateBillsFromTemplateRequest struct {
	TemplateID string   `json:"template_id" validate:"required"`
	Period     string   `json:"period" validate:"required"` // Format: YYYY-MM
	UnitIDs    []string `json:"unit_ids,omitempty"`         // Empty = all units
}






