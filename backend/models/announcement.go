package models

import (
	"database/sql"
	"time"
)

type Announcement struct {
	ID              string         `json:"id" db:"id"`
	TenantID        string         `json:"tenant_id" db:"tenant_id"`
	AuthorID        string         `json:"author_id" db:"author_id"`
	Title           string         `json:"title" db:"title"`
	Content         string         `json:"content" db:"content"`
	Priority        string         `json:"priority" db:"priority"`
	Category        sql.NullString `json:"category,omitempty" db:"category"`
	IsPinned        bool           `json:"is_pinned" db:"is_pinned"`
	SentNotification bool          `json:"sent_notification" db:"sent_notification"`
	SentWhatsApp    bool           `json:"sent_whatsapp" db:"sent_whatsapp"`
	SentAt          sql.NullTime   `json:"sent_at,omitempty" db:"sent_at"`
	ExpiresAt       sql.NullTime   `json:"expires_at,omitempty" db:"expires_at"`
	Metadata        sql.NullString `json:"metadata,omitempty" db:"metadata"`
	CreatedAt       time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at" db:"updated_at"`
	DeletedAt       sql.NullTime   `json:"-" db:"deleted_at"`
	// Joined fields
	AuthorName      sql.NullString `json:"author_name,omitempty" db:"author_name"`
	AuthorEmail     sql.NullString `json:"author_email,omitempty" db:"author_email"`
	IsRead          bool           `json:"is_read,omitempty" db:"is_read"`
}

type CreateAnnouncementRequest struct {
	Title     string  `json:"title" validate:"required"`
	Content   string  `json:"content" validate:"required"`
	Priority  string  `json:"priority" validate:"omitempty,oneof=low medium high"`
	Category  *string `json:"category,omitempty"`
	IsPinned  *bool   `json:"is_pinned,omitempty"`
	ExpiresAt *string `json:"expires_at,omitempty"` // ISO 8601 format
}

type UpdateAnnouncementRequest struct {
	Title     *string `json:"title,omitempty"`
	Content   *string `json:"content,omitempty"`
	Priority  *string `json:"priority,omitempty" validate:"omitempty,oneof=low medium high"`
	Category  *string `json:"category,omitempty"`
	IsPinned  *bool   `json:"is_pinned,omitempty"`
	ExpiresAt *string `json:"expires_at,omitempty"`
}









