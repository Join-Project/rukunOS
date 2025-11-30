package models

import (
	"database/sql"
	"time"
)

type Unit struct {
	ID         string         `json:"id" db:"id"`
	TenantID   string         `json:"tenant_id" db:"tenant_id"`
	Code       string         `json:"code" db:"code"`
	Type       string         `json:"type" db:"type"`
	OwnerName  sql.NullString `json:"owner_name,omitempty" db:"owner_name"`
	OwnerPhone sql.NullString `json:"owner_phone,omitempty" db:"owner_phone"`
	OwnerEmail sql.NullString `json:"owner_email,omitempty" db:"owner_email"`
	Address    sql.NullString `json:"address,omitempty" db:"address"`
	Status     string         `json:"status" db:"status"`
	CreatedAt  time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at" db:"updated_at"`
	DeletedAt  sql.NullTime   `json:"-" db:"deleted_at"`
}

type CreateUnitRequest struct {
	Code       string  `json:"code" validate:"required"`
	Type       string  `json:"type" validate:"required,oneof=rumah ruko kios"`
	OwnerName  *string `json:"owner_name,omitempty"`
	OwnerPhone *string `json:"owner_phone,omitempty"`
	OwnerEmail *string `json:"owner_email,omitempty" validate:"omitempty,email"`
	Address    *string `json:"address,omitempty"`
}

type UpdateUnitRequest struct {
	Code       *string `json:"code,omitempty"`
	Type       *string `json:"type,omitempty" validate:"omitempty,oneof=rumah ruko kios"`
	OwnerName  *string `json:"owner_name,omitempty"`
	OwnerPhone *string `json:"owner_phone,omitempty"`
	OwnerEmail *string `json:"owner_email,omitempty" validate:"omitempty,email"`
	Address    *string `json:"address,omitempty"`
	Status     *string `json:"status,omitempty" validate:"omitempty,oneof=active inactive vacant"`
}

