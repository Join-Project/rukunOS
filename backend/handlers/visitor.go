package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"
	"rukunos-backend/db"
	"rukunos-backend/middleware"
	"rukunos-backend/models"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// ListVisitorLogs lists all visitor logs for the tenant
func ListVisitorLogs(c echo.Context) error {
	tenantID := c.Get(string(middleware.CtxTenantID)).(string)

	// Get query parameters
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit < 1 {
		limit = 20
	}
	unitID := c.QueryParam("unit_id")
	status := c.QueryParam("status") // "checked_in" or "checked_out"
	search := c.QueryParam("search")

	offset := (page - 1) * limit

	// Build query
	query := `
		SELECT 
			v.id, v.tenant_id, v.unit_id, v.visitor_name, v.visitor_phone,
			v.visitor_id_number, v.visitor_vehicle, v.purpose, v.host_name,
			v.checked_in_at, v.checked_out_at, v.notes, v.checked_in_by, v.checked_out_by,
			v.created_at, v.updated_at,
			u.code as unit_code
		FROM visitor_logs v
		LEFT JOIN units u ON v.unit_id = u.id
		WHERE v.tenant_id = $1 AND v.deleted_at IS NULL
	`
	args := []interface{}{tenantID}
	argIndex := 2

	if unitID != "" {
		query += ` AND v.unit_id = $` + strconv.Itoa(argIndex)
		args = append(args, unitID)
		argIndex++
	}

	if status == "checked_in" {
		query += ` AND v.checked_out_at IS NULL`
	} else if status == "checked_out" {
		query += ` AND v.checked_out_at IS NOT NULL`
	}

	if search != "" {
		searchPattern := "%" + search + "%"
		query += ` AND (v.visitor_name ILIKE $` + strconv.Itoa(argIndex) + ` OR v.host_name ILIKE $` + strconv.Itoa(argIndex) + `)`
		args = append(args, searchPattern)
		argIndex++
	}

	query += ` ORDER BY v.checked_in_at DESC LIMIT $` + strconv.Itoa(argIndex) + ` OFFSET $` + strconv.Itoa(argIndex+1)
	args = append(args, limit, offset)

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}
	defer rows.Close()

	visitors := []map[string]interface{}{}
	for rows.Next() {
		var v models.VisitorLog
		var phone, idNumber, vehicle, purpose, hostName, notes, checkedInBy, checkedOutBy, unitCode sql.NullString
		var checkedOutAt sql.NullTime

		err := rows.Scan(
			&v.ID, &v.TenantID, &v.UnitID, &v.VisitorName, &phone,
			&idNumber, &vehicle, &purpose, &hostName,
			&v.CheckedInAt, &checkedOutAt, &notes, &checkedInBy, &checkedOutBy,
			&v.CreatedAt, &v.UpdatedAt, &unitCode,
		)
		if err != nil {
			continue
		}

		visitorData := map[string]interface{}{
			"id":            v.ID,
			"visitor_name":  v.VisitorName,
			"checked_in_at": v.CheckedInAt.Format(time.RFC3339),
			"created_at":    v.CreatedAt.Format(time.RFC3339),
		}

		if v.UnitID.Valid {
			visitorData["unit_id"] = v.UnitID.String
		}
		if phone.Valid {
			visitorData["visitor_phone"] = phone.String
		}
		if idNumber.Valid {
			visitorData["visitor_id_number"] = idNumber.String
		}
		if vehicle.Valid {
			visitorData["visitor_vehicle"] = vehicle.String
		}
		if purpose.Valid {
			visitorData["purpose"] = purpose.String
		}
		if hostName.Valid {
			visitorData["host_name"] = hostName.String
		}
		if notes.Valid {
			visitorData["notes"] = notes.String
		}
		if checkedOutAt.Valid {
			visitorData["checked_out_at"] = checkedOutAt.Time.Format(time.RFC3339)
		}
		if unitCode.Valid {
			visitorData["unit_code"] = unitCode.String
		}

		visitors = append(visitors, visitorData)
	}

	// Get total count
	var total int
	countQuery := `SELECT COUNT(*) FROM visitor_logs WHERE tenant_id = $1 AND deleted_at IS NULL`
	countArgs := []interface{}{tenantID}
	countArgIndex := 2

	if unitID != "" {
		countQuery += ` AND unit_id = $` + strconv.Itoa(countArgIndex)
		countArgs = append(countArgs, unitID)
		countArgIndex++
	}
	if status == "checked_in" {
		countQuery += ` AND checked_out_at IS NULL`
	} else if status == "checked_out" {
		countQuery += ` AND checked_out_at IS NOT NULL`
	}

	err = db.DB.Get(&total, countQuery, countArgs...)
	if err != nil {
		total = len(visitors)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"visitors": visitors,
		"pagination": map[string]interface{}{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": (total + limit - 1) / limit,
		},
	})
}

// CreateVisitorLog creates a new visitor log entry
func CreateVisitorLog(c echo.Context) error {
	tenantID := c.Get(string(middleware.CtxTenantID)).(string)
	userID := c.Get(string(middleware.CtxUserID)).(string)

	req := new(models.CreateVisitorLogRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// Validate required fields
	if req.VisitorName == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "visitor_name is required"})
	}

	// Validate unit if provided
	if req.UnitID != nil && *req.UnitID != "" {
		var unitExists bool
		err := db.DB.Get(&unitExists, `
			SELECT EXISTS(
				SELECT 1 FROM units 
				WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL
			)
		`, *req.UnitID, tenantID)
		if err != nil || !unitExists {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Unit not found"})
		}
	}

	// Create visitor log
	visitorID := uuid.New().String()
	query := `
		INSERT INTO visitor_logs (id, tenant_id, unit_id, visitor_name, visitor_phone, visitor_id_number, visitor_vehicle, purpose, host_name, notes, checked_in_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, checked_in_at
	`

	var checkedInAt time.Time
	var returnedID string
	err := db.DB.QueryRow(query,
		visitorID, tenantID, req.UnitID, req.VisitorName, req.VisitorPhone,
		req.VisitorIDNumber, req.VisitorVehicle, req.Purpose, req.HostName, req.Notes, userID,
	).Scan(&returnedID, &checkedInAt)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create visitor log: " + err.Error()})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"id":            returnedID,
		"message":       "Visitor logged successfully",
		"checked_in_at": checkedInAt.Format(time.RFC3339),
	})
}

// CheckOutVisitor checks out a visitor
func CheckOutVisitor(c echo.Context) error {
	tenantID := c.Get(string(middleware.CtxTenantID)).(string)
	userID := c.Get(string(middleware.CtxUserID)).(string)
	visitorID := c.Param("visitor_id")

	// Check if visitor log exists and is not already checked out
	var exists bool
	var alreadyCheckedOut bool
	err := db.DB.QueryRow(`
		SELECT 
			EXISTS(SELECT 1 FROM visitor_logs WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL),
			EXISTS(SELECT 1 FROM visitor_logs WHERE id = $1 AND tenant_id = $2 AND checked_out_at IS NOT NULL)
	`, visitorID, tenantID).Scan(&exists, &alreadyCheckedOut)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}
	if !exists {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Visitor log not found"})
	}
	if alreadyCheckedOut {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Visitor already checked out"})
	}

	// Update checked_out_at
	_, err = db.DB.Exec(`
		UPDATE visitor_logs 
		SET checked_out_at = CURRENT_TIMESTAMP, checked_out_by = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2 AND tenant_id = $3
	`, userID, visitorID, tenantID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to check out visitor: " + err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Visitor checked out successfully",
	})
}

// DeleteVisitorLog deletes a visitor log (soft delete)
func DeleteVisitorLog(c echo.Context) error {
	tenantID := c.Get(string(middleware.CtxTenantID)).(string)
	visitorID := c.Param("visitor_id")

	// Check if visitor log exists
	var exists bool
	err := db.DB.Get(&exists, `
		SELECT EXISTS(
			SELECT 1 FROM visitor_logs 
			WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL
		)
	`, visitorID, tenantID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}
	if !exists {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Visitor log not found"})
	}

	// Soft delete
	_, err = db.DB.Exec(`
		UPDATE visitor_logs 
		SET deleted_at = CURRENT_TIMESTAMP 
		WHERE id = $1 AND tenant_id = $2
	`, visitorID, tenantID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete visitor log: " + err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Visitor log deleted successfully",
	})
}






