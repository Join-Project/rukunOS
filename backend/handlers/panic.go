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

// ListPanicAlerts lists all panic alerts for the tenant
func ListPanicAlerts(c echo.Context) error {
	tenantID := c.Get(string(middleware.CtxTenantID)).(string)
	userID := c.Get(string(middleware.CtxUserID)).(string)

	// Get query parameters
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit < 1 {
		limit = 20
	}
	status := c.QueryParam("status")

	offset := (page - 1) * limit

	// Check if user is admin/security (can see all alerts) or just their own
	var isAdmin bool
	err := db.DB.Get(&isAdmin, `
		SELECT EXISTS(
			SELECT 1 FROM tenant_users tu
			JOIN roles r ON tu.role_id = r.id
			WHERE tu.user_id = $1 AND tu.tenant_id = $2
			AND (r.name = 'Admin' OR r.name = 'Security')
			AND tu.deleted_at IS NULL AND r.deleted_at IS NULL
		)
	`, userID, tenantID)

	// Build query
	query := `
		SELECT 
			p.id, p.tenant_id, p.user_id, p.unit_id, p.location, p.status,
			p.responded_by, p.responded_at, p.resolved_at, p.notes,
			p.created_at, p.updated_at,
			u.full_name as user_name, u.email as user_email, u.phone as user_phone,
			un.code as unit_code,
			ru.full_name as responder_name
		FROM panic_alerts p
		LEFT JOIN users u ON p.user_id = u.id
		LEFT JOIN units un ON p.unit_id = un.id
		LEFT JOIN users ru ON p.responded_by = ru.id
		WHERE p.tenant_id = $1 AND p.deleted_at IS NULL
	`
	args := []interface{}{tenantID}
	argIndex := 2

	// Filter by user if not admin
	if err == nil && !isAdmin {
		query += ` AND p.user_id = $` + strconv.Itoa(argIndex)
		args = append(args, userID)
		argIndex++
	}

	if status != "" {
		query += ` AND p.status = $` + strconv.Itoa(argIndex)
		args = append(args, status)
		argIndex++
	}

	query += ` ORDER BY p.created_at DESC LIMIT $` + strconv.Itoa(argIndex) + ` OFFSET $` + strconv.Itoa(argIndex+1)
	args = append(args, limit, offset)

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}
	defer rows.Close()

	alerts := []map[string]interface{}{}
	for rows.Next() {
		var p models.PanicAlert
		var location, notes, unitID, respondedBy, userName, userEmail, userPhone, unitCode, responderName sql.NullString
		var respondedAt, resolvedAt sql.NullTime

		err := rows.Scan(
			&p.ID, &p.TenantID, &p.UserID, &unitID, &location, &p.Status,
			&respondedBy, &respondedAt, &resolvedAt, &notes,
			&p.CreatedAt, &p.UpdatedAt,
			&userName, &userEmail, &userPhone, &unitCode, &responderName,
		)
		if err != nil {
			continue
		}

		alertData := map[string]interface{}{
			"id":         p.ID,
			"user_id":    p.UserID,
			"status":     p.Status,
			"created_at": p.CreatedAt.Format(time.RFC3339),
		}

		if unitID.Valid {
			alertData["unit_id"] = unitID.String
		}
		if location.Valid {
			alertData["location"] = location.String
		}
		if notes.Valid {
			alertData["notes"] = notes.String
		}
		if respondedBy.Valid {
			alertData["responded_by"] = respondedBy.String
		}
		if respondedAt.Valid {
			alertData["responded_at"] = respondedAt.Time.Format(time.RFC3339)
		}
		if resolvedAt.Valid {
			alertData["resolved_at"] = resolvedAt.Time.Format(time.RFC3339)
		}
		if userName.Valid {
			alertData["user_name"] = userName.String
		}
		if userEmail.Valid {
			alertData["user_email"] = userEmail.String
		}
		if userPhone.Valid {
			alertData["user_phone"] = userPhone.String
		}
		if unitCode.Valid {
			alertData["unit_code"] = unitCode.String
		}
		if responderName.Valid {
			alertData["responder_name"] = responderName.String
		}

		alerts = append(alerts, alertData)
	}

	// Get total count
	var total int
	countQuery := `SELECT COUNT(*) FROM panic_alerts WHERE tenant_id = $1 AND deleted_at IS NULL`
	countArgs := []interface{}{tenantID}
	countArgIndex := 2

	if err == nil && !isAdmin {
		countQuery += ` AND user_id = $` + strconv.Itoa(countArgIndex)
		countArgs = append(countArgs, userID)
		countArgIndex++
	}
	if status != "" {
		countQuery += ` AND status = $` + strconv.Itoa(countArgIndex)
		countArgs = append(countArgs, status)
	}

	err = db.DB.Get(&total, countQuery, countArgs...)
	if err != nil {
		total = len(alerts)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"alerts": alerts,
		"pagination": map[string]interface{}{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": (total + limit - 1) / limit,
		},
	})
}

// CreatePanicAlert creates a new panic alert
func CreatePanicAlert(c echo.Context) error {
	tenantID := c.Get(string(middleware.CtxTenantID)).(string)
	userID := c.Get(string(middleware.CtxUserID)).(string)

	req := new(models.CreatePanicAlertRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// Get user's unit_id
	var unitID sql.NullString
	_ = db.DB.Get(&unitID, `
		SELECT unit_id FROM tenant_users
		WHERE user_id = $1 AND tenant_id = $2 AND deleted_at IS NULL
	`, userID, tenantID)

	// Create panic alert
	alertID := uuid.New().String()
	query := `
		INSERT INTO panic_alerts (id, tenant_id, user_id, unit_id, location, status)
		VALUES ($1, $2, $3, $4, $5, 'active')
		RETURNING id, created_at
	`

	var location sql.NullString
	if req.Location != nil {
		location = sql.NullString{String: *req.Location, Valid: true}
	}

	var createdAt time.Time
	var returnedID string
	err := db.DB.QueryRow(query, alertID, tenantID, userID, unitID, location).Scan(&returnedID, &createdAt)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create panic alert: " + err.Error()})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"id":         returnedID,
		"message":    "Panic alert created successfully",
		"created_at": createdAt.Format(time.RFC3339),
	})
}

// UpdatePanicAlert updates a panic alert (for responding/resolving)
func UpdatePanicAlert(c echo.Context) error {
	tenantID := c.Get(string(middleware.CtxTenantID)).(string)
	userID := c.Get(string(middleware.CtxUserID)).(string)
	alertID := c.Param("alert_id")

	req := new(models.UpdatePanicAlertRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// Check if alert exists
	var exists bool
	err := db.DB.Get(&exists, `
		SELECT EXISTS(
			SELECT 1 FROM panic_alerts 
			WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL
		)
	`, alertID, tenantID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}
	if !exists {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Panic alert not found"})
	}

	// Build update query
	updates := []string{}
	args := []interface{}{}
	argIndex := 1

	if req.Status != nil {
		updates = append(updates, "status = $"+strconv.Itoa(argIndex))
		args = append(args, *req.Status)
		argIndex++

		// Set responded_at if status changes to "responded"
		if *req.Status == "responded" {
			updates = append(updates, "responded_by = $"+strconv.Itoa(argIndex))
			updates = append(updates, "responded_at = CURRENT_TIMESTAMP")
			args = append(args, userID)
			argIndex++
		}

		// Set resolved_at if status changes to "resolved"
		if *req.Status == "resolved" {
			updates = append(updates, "resolved_at = CURRENT_TIMESTAMP")
		}
	}

	if req.Notes != nil {
		updates = append(updates, "notes = $"+strconv.Itoa(argIndex))
		args = append(args, *req.Notes)
		argIndex++
	}

	if len(updates) == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "No fields to update"})
	}

	updates = append(updates, "updated_at = CURRENT_TIMESTAMP")
	query := `UPDATE panic_alerts SET ` + updates[0]
	for i := 1; i < len(updates); i++ {
		query += `, ` + updates[i]
	}
	query += ` WHERE id = $` + strconv.Itoa(argIndex) + ` AND tenant_id = $` + strconv.Itoa(argIndex+1)
	args = append(args, alertID, tenantID)

	_, err = db.DB.Exec(query, args...)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update panic alert: " + err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Panic alert updated successfully",
		"id":      alertID,
	})
}






