package models

import (
	"database/sql"
	"time"
)

type VisitorLog struct {
	ID              string         `json:"id" db:"id"`
	TenantID        string         `json:"tenant_id" db:"tenant_id"`
	UnitID          sql.NullString `json:"unit_id,omitempty" db:"unit_id"`
	VisitorName     string         `json:"visitor_name" db:"visitor_name"`
	VisitorPhone    sql.NullString `json:"visitor_phone,omitempty" db:"visitor_phone"`
	VisitorIDNumber sql.NullString `json:"visitor_id_number,omitempty" db:"visitor_id_number"`
	VisitorVehicle  sql.NullString `json:"visitor_vehicle,omitempty" db:"visitor_vehicle"`
	Purpose         sql.NullString `json:"purpose,omitempty" db:"purpose"`
	HostName        sql.NullString `json:"host_name,omitempty" db:"host_name"`
	CheckedInAt     time.Time      `json:"checked_in_at" db:"checked_in_at"`
	CheckedOutAt    sql.NullTime   `json:"checked_out_at,omitempty" db:"checked_out_at"`
	Notes           sql.NullString `json:"notes,omitempty" db:"notes"`
	CheckedInBy     sql.NullString `json:"checked_in_by,omitempty" db:"checked_in_by"`
	CheckedOutBy    sql.NullString `json:"checked_out_by,omitempty" db:"checked_out_by"`
	CreatedAt       time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at" db:"updated_at"`
	DeletedAt       sql.NullTime   `json:"-" db:"deleted_at"`
	// Joined fields
	UnitCode        sql.NullString `json:"unit_code,omitempty" db:"unit_code"`
}

type CreateVisitorLogRequest struct {
	UnitID          *string `json:"unit_id,omitempty"`
	VisitorName     string  `json:"visitor_name" validate:"required"`
	VisitorPhone    *string `json:"visitor_phone,omitempty"`
	VisitorIDNumber *string `json:"visitor_id_number,omitempty"`
	VisitorVehicle  *string `json:"visitor_vehicle,omitempty"`
	Purpose         *string `json:"purpose,omitempty"`
	HostName        *string `json:"host_name,omitempty"`
	Notes           *string `json:"notes,omitempty"`
}

type UpdateVisitorLogRequest struct {
	VisitorName     *string `json:"visitor_name,omitempty"`
	VisitorPhone    *string `json:"visitor_phone,omitempty"`
	VisitorIDNumber *string `json:"visitor_id_number,omitempty"`
	VisitorVehicle  *string `json:"visitor_vehicle,omitempty"`
	Purpose         *string `json:"purpose,omitempty"`
	HostName        *string `json:"host_name,omitempty"`
	Notes           *string `json:"notes,omitempty"`
}










