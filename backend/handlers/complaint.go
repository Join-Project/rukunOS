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
	"github.com/lib/pq"
)

// ListComplaints lists all complaints for the tenant (admin) or current user (warga)
func ListComplaints(c echo.Context) error {
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
	category := c.QueryParam("category")
	priority := c.QueryParam("priority")

	offset := (page - 1) * limit

	// Check if user is admin (can see all complaints)
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
			c.id, c.tenant_id, c.user_id, c.unit_id, c.category, c.priority,
			c.title, c.description, c.status, c.assigned_to, c.resolved_at,
			c.resolution_notes, COALESCE(c.attachment_urls, ARRAY[]::TEXT[]), c.created_at, c.updated_at,
			u.full_name as user_name, u.email as user_email,
			un.code as unit_code,
			au.full_name as assigned_to_name
		FROM complaints c
		LEFT JOIN users u ON c.user_id = u.id
		LEFT JOIN units un ON c.unit_id = un.id
		LEFT JOIN users au ON c.assigned_to = au.id
		WHERE c.tenant_id = $1 AND c.deleted_at IS NULL
	`
	args := []interface{}{tenantID}
	argIndex := 2

	// Filter by user if not admin
	if err == nil && !isAdmin {
		query += ` AND c.user_id = $` + strconv.Itoa(argIndex)
		args = append(args, userID)
		argIndex++
	}

	if status != "" {
		query += ` AND c.status = $` + strconv.Itoa(argIndex)
		args = append(args, status)
		argIndex++
	}

	if category != "" {
		query += ` AND c.category = $` + strconv.Itoa(argIndex)
		args = append(args, category)
		argIndex++
	}

	if priority != "" {
		query += ` AND c.priority = $` + strconv.Itoa(argIndex)
		args = append(args, priority)
		argIndex++
	}

	query += ` ORDER BY c.created_at DESC LIMIT $` + strconv.Itoa(argIndex) + ` OFFSET $` + strconv.Itoa(argIndex+1)
	args = append(args, limit, offset)

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		c.Logger().Errorf("Error executing query: %v, query: %s", err, query)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error: " + err.Error()})
	}
	defer rows.Close()

	complaints := []map[string]interface{}{}
	for rows.Next() {
		var complaint models.Complaint
		var unitID, assignedTo, userName, userEmail, unitCode, assignedToName sql.NullString
		var resolvedAt sql.NullTime
		var resolutionNotes sql.NullString
		var attachmentURLs pq.StringArray

		err := rows.Scan(
			&complaint.ID, &complaint.TenantID, &complaint.UserID, &unitID, &complaint.Category, &complaint.Priority,
			&complaint.Title, &complaint.Description, &complaint.Status, &assignedTo, &resolvedAt,
			&resolutionNotes, &attachmentURLs, &complaint.CreatedAt, &complaint.UpdatedAt,
			&userName, &userEmail, &unitCode, &assignedToName,
		)
		if err != nil {
			// Log error but continue processing other rows
			c.Logger().Errorf("Error scanning complaint row: %v", err)
			continue
		}

		complaintData := map[string]interface{}{
			"id":          complaint.ID,
			"user_id":     complaint.UserID,
			"category":    complaint.Category,
			"priority":    complaint.Priority,
			"title":       complaint.Title,
			"description": complaint.Description,
			"status":      complaint.Status,
			"created_at":  complaint.CreatedAt.Format(time.RFC3339),
			"updated_at":  complaint.UpdatedAt.Format(time.RFC3339),
		}

		if unitID.Valid {
			complaintData["unit_id"] = unitID.String
		}
		if unitCode.Valid {
			complaintData["unit_code"] = unitCode.String
		}
		if userName.Valid {
			complaintData["user_name"] = userName.String
		}
		if userEmail.Valid {
			complaintData["user_email"] = userEmail.String
		}
		if assignedTo.Valid {
			complaintData["assigned_to"] = assignedTo.String
		}
		if assignedToName.Valid {
			complaintData["assigned_to_name"] = assignedToName.String
		}
		if resolvedAt.Valid {
			complaintData["resolved_at"] = resolvedAt.Time.Format(time.RFC3339)
		}
		if resolutionNotes.Valid {
			complaintData["resolution_notes"] = resolutionNotes.String
		}
		if len(attachmentURLs) > 0 {
			complaintData["attachment_urls"] = []string(attachmentURLs)
		}

		complaints = append(complaints, complaintData)
	}

	// Get total count
	countQuery := `
		SELECT COUNT(*) FROM complaints c
		WHERE c.tenant_id = $1 AND c.deleted_at IS NULL
	`
	countArgs := []interface{}{tenantID}
	countArgIndex := 2

	if !isAdmin {
		countQuery += ` AND c.user_id = $` + strconv.Itoa(countArgIndex)
		countArgs = append(countArgs, userID)
		countArgIndex++
	}

	if status != "" {
		countQuery += ` AND c.status = $` + strconv.Itoa(countArgIndex)
		countArgs = append(countArgs, status)
		countArgIndex++
	}

	if category != "" {
		countQuery += ` AND c.category = $` + strconv.Itoa(countArgIndex)
		countArgs = append(countArgs, category)
		countArgIndex++
	}

	if priority != "" {
		countQuery += ` AND c.priority = $` + strconv.Itoa(countArgIndex)
		countArgs = append(countArgs, priority)
		countArgIndex++
	}

	var total int
	err = db.DB.Get(&total, countQuery, countArgs...)
	if err != nil {
		total = len(complaints)
	}

	totalPages := (total + limit - 1) / limit

	return c.JSON(http.StatusOK, map[string]interface{}{
		"complaints": complaints,
		"pagination": map[string]interface{}{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": totalPages,
		},
	})
}

// GetComplaint gets a single complaint by ID
func GetComplaint(c echo.Context) error {
	tenantID := c.Get(string(middleware.CtxTenantID)).(string)
	userID := c.Get(string(middleware.CtxUserID)).(string)
	complaintID := c.Param("complaint_id")

	// Check if user is admin or owns the complaint
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

	query := `
		SELECT 
			c.id, c.tenant_id, c.user_id, c.unit_id, c.category, c.priority,
			c.title, c.description, c.status, c.assigned_to, c.resolved_at,
			c.resolution_notes, COALESCE(c.attachment_urls, ARRAY[]::TEXT[]), c.created_at, c.updated_at,
			u.full_name as user_name, u.email as user_email,
			un.code as unit_code,
			au.full_name as assigned_to_name
		FROM complaints c
		LEFT JOIN users u ON c.user_id = u.id
		LEFT JOIN units un ON c.unit_id = un.id
		LEFT JOIN users au ON c.assigned_to = au.id
		WHERE c.id = $1 AND c.tenant_id = $2 AND c.deleted_at IS NULL
	`
	args := []interface{}{complaintID, tenantID}

	if !isAdmin {
		query += ` AND c.user_id = $3`
		args = append(args, userID)
	}

	var complaint models.Complaint
	var unitID, assignedTo, userName, userEmail, unitCode, assignedToName sql.NullString
	var resolvedAt sql.NullTime
	var resolutionNotes sql.NullString
	var attachmentURLs pq.StringArray

	err = db.DB.QueryRow(query, args...).Scan(
		&complaint.ID, &complaint.TenantID, &complaint.UserID, &unitID, &complaint.Category, &complaint.Priority,
		&complaint.Title, &complaint.Description, &complaint.Status, &assignedTo, &resolvedAt,
		&resolutionNotes, &attachmentURLs, &complaint.CreatedAt, &complaint.UpdatedAt,
		&userName, &userEmail, &unitCode, &assignedToName,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Complaint not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}

	complaintData := map[string]interface{}{
		"id":          complaint.ID,
		"user_id":     complaint.UserID,
		"category":    complaint.Category,
		"priority":    complaint.Priority,
		"title":       complaint.Title,
		"description": complaint.Description,
		"status":      complaint.Status,
		"created_at":  complaint.CreatedAt.Format(time.RFC3339),
		"updated_at":  complaint.UpdatedAt.Format(time.RFC3339),
	}

	if unitID.Valid {
		complaintData["unit_id"] = unitID.String
	}
	if unitCode.Valid {
		complaintData["unit_code"] = unitCode.String
	}
	if userName.Valid {
		complaintData["user_name"] = userName.String
	}
	if userEmail.Valid {
		complaintData["user_email"] = userEmail.String
	}
	if assignedTo.Valid {
		complaintData["assigned_to"] = assignedTo.String
	}
	if assignedToName.Valid {
		complaintData["assigned_to_name"] = assignedToName.String
	}
	if resolvedAt.Valid {
		complaintData["resolved_at"] = resolvedAt.Time.Format(time.RFC3339)
	}
	if resolutionNotes.Valid {
		complaintData["resolution_notes"] = resolutionNotes.String
	}
	if len(attachmentURLs) > 0 {
		complaintData["attachment_urls"] = attachmentURLs
	}

	return c.JSON(http.StatusOK, complaintData)
}

// CreateComplaint creates a new complaint
func CreateComplaint(c echo.Context) error {
	tenantID := c.Get(string(middleware.CtxTenantID)).(string)
	userID := c.Get(string(middleware.CtxUserID)).(string)

	var req models.CreateComplaintRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// Validate required fields
	if req.Title == "" || req.Description == "" || req.Category == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Title, description, and category are required"})
	}

	// Set default priority
	priority := req.Priority
	if priority == "" {
		priority = "normal"
	}

	// Get user's unit_id
	var unitID sql.NullString
	err := db.DB.Get(&unitID, `
		SELECT unit_id FROM tenant_users
		WHERE user_id = $1 AND tenant_id = $2 AND deleted_at IS NULL
	`, userID, tenantID)

	// Create complaint
	complaintID := uuid.New().String()
	query := `
		INSERT INTO complaints (id, tenant_id, user_id, unit_id, category, priority, title, description, attachment_urls)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at
	`

	var createdAt time.Time
	var returnedID string
	err = db.DB.QueryRow(query,
		complaintID, tenantID, userID, unitID, req.Category, priority, req.Title, req.Description, pq.Array(req.AttachmentURLs),
	).Scan(&returnedID, &createdAt)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create complaint: " + err.Error()})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"id":         returnedID,
		"created_at": createdAt.Format(time.RFC3339),
	})
}

// UpdateComplaint updates a complaint (admin only for status/assignment)
func UpdateComplaint(c echo.Context) error {
	tenantID := c.Get(string(middleware.CtxTenantID)).(string)
	userID := c.Get(string(middleware.CtxUserID)).(string)
	complaintID := c.Param("complaint_id")

	var req models.UpdateComplaintRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// Check if user is admin
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

	if !isAdmin {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "Only admin can update complaints"})
	}

	// Build update query dynamically
	updates := []string{"updated_at = CURRENT_TIMESTAMP"}
	args := []interface{}{complaintID, tenantID}
	argIndex := 3

	if req.Status != nil {
		updates = append(updates, "status = $"+strconv.Itoa(argIndex))
		args = append(args, *req.Status)
		argIndex++

		// If status is resolved, set resolved_at
		if *req.Status == "resolved" {
			updates = append(updates, "resolved_at = CURRENT_TIMESTAMP")
		}
	}

	if req.AssignedTo != nil {
		updates = append(updates, "assigned_to = $"+strconv.Itoa(argIndex))
		args = append(args, *req.AssignedTo)
		argIndex++
	}

	if req.ResolutionNotes != nil {
		updates = append(updates, "resolution_notes = $"+strconv.Itoa(argIndex))
		args = append(args, *req.ResolutionNotes)
		argIndex++
	}

	if len(updates) == 1 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "No fields to update"})
	}

	query := `UPDATE complaints SET ` + updates[0]
	for i := 1; i < len(updates); i++ {
		query += `, ` + updates[i]
	}
	query += ` WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL RETURNING id`

	var returnedID string
	err = db.DB.QueryRow(query, args...).Scan(&returnedID)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Complaint not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update complaint: " + err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Complaint updated successfully"})
}

// DeleteComplaint soft deletes a complaint
func DeleteComplaint(c echo.Context) error {
	tenantID := c.Get(string(middleware.CtxTenantID)).(string)
	userID := c.Get(string(middleware.CtxUserID)).(string)
	complaintID := c.Param("complaint_id")

	// Check if user is admin or owns the complaint
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

	query := `UPDATE complaints SET deleted_at = CURRENT_TIMESTAMP WHERE id = $1 AND tenant_id = $2`
	args := []interface{}{complaintID, tenantID}

	if !isAdmin {
		query += ` AND user_id = $3`
		args = append(args, userID)
	}

	result, err := db.DB.Exec(query, args...)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Complaint not found"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Complaint deleted successfully"})
}





