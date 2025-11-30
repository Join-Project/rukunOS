package models

import (
	"database/sql"
	"time"
)

type DocumentRequest struct {
	ID            string         `json:"id" db:"id"`
	TenantID      string         `json:"tenant_id" db:"tenant_id"`
	UserID        string         `json:"user_id" db:"user_id"`
	DocumentType  string         `json:"document_type" db:"document_type"`
	Purpose       string         `json:"purpose" db:"purpose"`
	Status        string         `json:"status" db:"status"`
	ApprovedBy    sql.NullString `json:"approved_by,omitempty" db:"approved_by"`
	ApprovedAt    sql.NullTime   `json:"approved_at,omitempty" db:"approved_at"`
	RejectedReason sql.NullString `json:"rejected_reason,omitempty" db:"rejected_reason"`
	AttachmentIDs []string       `json:"attachment_ids,omitempty" db:"attachment_ids"`
	Notes         sql.NullString `json:"notes,omitempty" db:"notes"`
	Metadata      []byte         `json:"metadata,omitempty" db:"metadata"`
	CreatedAt     time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at" db:"updated_at"`
	DeletedAt     sql.NullTime   `json:"-" db:"deleted_at"`
	// Joined fields
	UserName      sql.NullString `json:"user_name,omitempty" db:"user_name"`
	UserEmail     sql.NullString `json:"user_email,omitempty" db:"user_email"`
	ApprovedByName sql.NullString `json:"approved_by_name,omitempty" db:"approved_by_name"`
}

type CreateDocumentRequestRequest struct {
	DocumentType string   `json:"document_type" validate:"required"`
	Purpose      string   `json:"purpose" validate:"required"`
	AttachmentIDs []string `json:"attachment_ids,omitempty"`
}

type UpdateDocumentRequestRequest struct {
	Status         *string `json:"status,omitempty"` // pending, approved, rejected, completed
	RejectedReason *string `json:"rejected_reason,omitempty"`
	Notes          *string `json:"notes,omitempty"`
}









