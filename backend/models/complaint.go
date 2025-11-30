package models

import (
	"database/sql"
	"time"
)

type Complaint struct {
	ID              string         `json:"id" db:"id"`
	TenantID        string         `json:"tenant_id" db:"tenant_id"`
	UserID          string         `json:"user_id" db:"user_id"`
	UnitID          sql.NullString `json:"unit_id,omitempty" db:"unit_id"`
	Category        string         `json:"category" db:"category"`
	Priority        string         `json:"priority" db:"priority"`
	Title           string         `json:"title" db:"title"`
	Description     string         `json:"description" db:"description"`
	Status          string         `json:"status" db:"status"`
	AssignedTo      sql.NullString `json:"assigned_to,omitempty" db:"assigned_to"`
	ResolvedAt      sql.NullTime   `json:"resolved_at,omitempty" db:"resolved_at"`
	ResolutionNotes sql.NullString `json:"resolution_notes,omitempty" db:"resolution_notes"`
	AttachmentURLs  []string       `json:"attachment_urls,omitempty" db:"attachment_urls"`
	Metadata        []byte         `json:"metadata,omitempty" db:"metadata"`
	CreatedAt       time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at" db:"updated_at"`
	DeletedAt       sql.NullTime   `json:"-" db:"deleted_at"`
	// Joined fields
	UserName        sql.NullString `json:"user_name,omitempty" db:"user_name"`
	UserEmail       sql.NullString `json:"user_email,omitempty" db:"user_email"`
	UnitCode        sql.NullString `json:"unit_code,omitempty" db:"unit_code"`
	AssignedToName  sql.NullString `json:"assigned_to_name,omitempty" db:"assigned_to_name"`
}

type CreateComplaintRequest struct {
	Category    string   `json:"category" validate:"required"`
	Priority    string   `json:"priority,omitempty"` // Default to normal
	Title       string   `json:"title" validate:"required"`
	Description string   `json:"description" validate:"required"`
	AttachmentURLs []string `json:"attachment_urls,omitempty"`
}

type UpdateComplaintRequest struct {
	Status          *string `json:"status,omitempty"`
	AssignedTo      *string `json:"assigned_to,omitempty"`
	ResolutionNotes *string `json:"resolution_notes,omitempty"`
}










