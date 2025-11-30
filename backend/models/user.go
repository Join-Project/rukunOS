package models

import (
	"time"
)

type User struct {
	ID           string     `json:"id" db:"id"`
	Email        string     `json:"email" db:"email"`
	PasswordHash string     `json:"-" db:"password_hash"`
	FullName     string     `json:"full_name" db:"full_name"`
	Phone        *string    `json:"phone,omitempty" db:"phone"`
	GoogleID     *string    `json:"-" db:"google_id"`
	AuthProvider string     `json:"auth_provider" db:"auth_provider"`
	AvatarURL    *string    `json:"avatar_url,omitempty" db:"avatar_url"`
	Status       string     `json:"status" db:"status"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
	// Tenant-specific fields (from tenant_users)
	TenantID *string `json:"tenant_id,omitempty"`
	RoleID   *string `json:"role_id,omitempty"`
	UnitID   *string `json:"unit_id,omitempty"`
}

type RegisterRequest struct {
	Email           string  `json:"email" validate:"required,email"`
	Password        string  `json:"password" validate:"required,min=6"`
	FullName        string  `json:"full_name" validate:"required"`
	RegistrationType string `json:"registration_type" validate:"required,oneof=tenant warga"` // "tenant" or "warga"
	// For tenant registration
	TenantName      *string `json:"tenant_name,omitempty"`
	TenantCode      *string `json:"tenant_code,omitempty"`
	TenantAddress   *string `json:"tenant_address,omitempty"`
	// For warga registration (join existing tenant)
	TenantCodeJoin  *string `json:"tenant_code_join,omitempty"` // Alias for tenant_code when joining
}

type LoginRequest struct {
	Email     string  `json:"email" validate:"required,email"`
	Password  string  `json:"password" validate:"required"`
	TenantID  *string `json:"tenant_id,omitempty"` // Optional: jika user punya multiple tenants
}

type AuthResponse struct {
	Token    string   `json:"token"`
	User     User     `json:"user"`
	TenantID string   `json:"tenant_id"`
	Tenants  []string `json:"tenants,omitempty"` // List of tenant IDs if user has multiple
}

