package models

import (
	"database/sql"
	"time"
)

type Bill struct {
	ID              string         `json:"id" db:"id"`
	TenantID        string         `json:"tenant_id" db:"tenant_id"`
	UnitID          string         `json:"unit_id" db:"unit_id"`
	Category        string         `json:"category" db:"category"`
	Period          string         `json:"period" db:"period"`
	Amount          float64        `json:"amount" db:"amount"`
	LateFee         float64        `json:"late_fee" db:"late_fee"`
	DueDate         sql.NullTime   `json:"due_date,omitempty" db:"due_date"`
	Status          string         `json:"status" db:"status"`
	PaidAt          sql.NullTime   `json:"paid_at,omitempty" db:"paid_at"`
	PaymentMethod   sql.NullString `json:"payment_method,omitempty" db:"payment_method"`
	PaymentReference sql.NullString `json:"payment_reference,omitempty" db:"payment_reference"`
	Notes           sql.NullString `json:"notes,omitempty" db:"notes"`
	BillNumber      sql.NullString `json:"bill_number,omitempty" db:"bill_number"`
	CreatedBy       sql.NullString `json:"created_by,omitempty" db:"created_by"`
	CreatedAt       time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at" db:"updated_at"`
	DeletedAt       sql.NullTime   `json:"-" db:"deleted_at"`
	// Joined fields
	UnitCode        sql.NullString `json:"unit_code,omitempty" db:"unit_code"`
	UnitType        sql.NullString `json:"unit_type,omitempty" db:"unit_type"`
	UserName        sql.NullString `json:"user_name,omitempty" db:"user_name"`
	UserEmail       sql.NullString `json:"user_email,omitempty" db:"user_email"`
}

type CreateBillRequest struct {
	UnitID    string  `json:"unit_id" validate:"required"`
	Category  string  `json:"category" validate:"required"`
	Period    string  `json:"period" validate:"required"`
	Amount    float64 `json:"amount" validate:"required,min=0"`
	LateFee   *float64 `json:"late_fee,omitempty"`
	DueDate   *string  `json:"due_date,omitempty"` // Optional
	Notes     *string `json:"notes,omitempty"`
}

type UpdateBillRequest struct {
	Category  *string  `json:"category,omitempty"`
	Period    *string  `json:"period,omitempty"`
	Amount    *float64 `json:"amount,omitempty" validate:"omitempty,min=0"`
	LateFee   *float64 `json:"late_fee,omitempty"`
	DueDate   *string  `json:"due_date,omitempty"`
	Status    *string  `json:"status,omitempty" validate:"omitempty,oneof=pending paid overdue cancelled"`
	Notes     *string  `json:"notes,omitempty"`
}

type ProcessPaymentRequest struct {
	PaymentMethod   string  `json:"payment_method" validate:"required"`
	PaymentReference *string `json:"payment_reference,omitempty"`
	Amount          *float64 `json:"amount,omitempty"`
}

