package models

import (
	"database/sql"
	"time"
)

type PanicAlert struct {
	ID          string         `json:"id" db:"id"`
	TenantID    string         `json:"tenant_id" db:"tenant_id"`
	UserID      string         `json:"user_id" db:"user_id"`
	UnitID      sql.NullString `json:"unit_id,omitempty" db:"unit_id"`
	Location    sql.NullString `json:"location,omitempty" db:"location"` // JSONB stored as string
	Status      string         `json:"status" db:"status"`
	RespondedBy sql.NullString `json:"responded_by,omitempty" db:"responded_by"`
	RespondedAt sql.NullTime   `json:"responded_at,omitempty" db:"responded_at"`
	ResolvedAt  sql.NullTime   `json:"resolved_at,omitempty" db:"resolved_at"`
	Notes       sql.NullString `json:"notes,omitempty" db:"notes"`
	CreatedAt   time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at" db:"updated_at"`
	DeletedAt   sql.NullTime   `json:"-" db:"deleted_at"`
	// Joined fields
	UserName    sql.NullString `json:"user_name,omitempty" db:"user_name"`
	UserEmail   sql.NullString `json:"user_email,omitempty" db:"user_email"`
	UserPhone   sql.NullString `json:"user_phone,omitempty" db:"user_phone"`
	UnitCode    sql.NullString `json:"unit_code,omitempty" db:"unit_code"`
	ResponderName sql.NullString `json:"responder_name,omitempty" db:"responder_name"`
}

type CreatePanicAlertRequest struct {
	Location *string `json:"location,omitempty"` // JSON string: {"lat": 0, "lng": 0}
}

type UpdatePanicAlertRequest struct {
	Status *string `json:"status,omitempty" validate:"omitempty,oneof=active responded resolved"`
	Notes  *string `json:"notes,omitempty"`
}










