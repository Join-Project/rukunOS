package models

import (
	"time"
)

type Tenant struct {
	ID        string    `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Code      string    `json:"code" db:"code"`
	Address   *string   `json:"address,omitempty" db:"address"`
	Phone     *string   `json:"phone,omitempty" db:"phone"`
	Email     *string   `json:"email,omitempty" db:"email"`
	Settings  string    `json:"settings" db:"settings"` // JSONB stored as string
	Modules   string    `json:"modules" db:"modules"`   // JSONB stored as string
	Status    string    `json:"status" db:"status"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

type CreateTenantRequest struct {
	Name    string  `json:"name" validate:"required"`
	Code    string  `json:"code" validate:"required,alphanum"`
	Address *string `json:"address,omitempty"`
	Phone   *string `json:"phone,omitempty"`
	Email   *string `json:"email,omitempty" validate:"omitempty,email"`
}

type UpdateTenantRequest struct {
	Name    *string `json:"name,omitempty"`
	Address *string `json:"address,omitempty"`
	Phone   *string `json:"phone,omitempty"`
	Email   *string `json:"email,omitempty" validate:"omitempty,email"`
	Status  *string `json:"status,omitempty" validate:"omitempty,oneof=active suspended inactive"`
}











