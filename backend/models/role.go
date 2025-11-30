package models

import (
	"database/sql"
	"time"
)

type Role struct {
	ID          string         `json:"id" db:"id"`
	TenantID    string         `json:"tenant_id" db:"tenant_id"`
	Name        string         `json:"name" db:"name"`
	Description sql.NullString `json:"description,omitempty" db:"description"`
	IsSystem    bool           `json:"is_system" db:"is_system"`
	CreatedAt   time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at" db:"updated_at"`
	DeletedAt   sql.NullTime   `json:"-" db:"deleted_at"`
	Permissions []Permission   `json:"permissions,omitempty"`
}

type Permission struct {
	ID          string         `json:"id" db:"id"`
	Key         string         `json:"key" db:"key"`
	Name        string         `json:"name" db:"name"`
	Description sql.NullString `json:"description,omitempty" db:"description"`
	Module      sql.NullString `json:"module,omitempty" db:"module"`
	CreatedAt   time.Time      `json:"created_at" db:"created_at"`
}

type CreateRoleRequest struct {
	Name        string   `json:"name" validate:"required"`
	Description *string  `json:"description,omitempty"`
	Permissions []string `json:"permissions" validate:"required,min=1"`
}

type UpdateRoleRequest struct {
	Name        *string  `json:"name,omitempty"`
	Description *string  `json:"description,omitempty"`
	Permissions *[]string `json:"permissions,omitempty"`
}

type RolePermission struct {
	RoleID       string `json:"role_id" db:"role_id"`
	PermissionID string `json:"permission_id" db:"permission_id"`
}

