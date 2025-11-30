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

// ListDocumentRequests lists all document requests for the tenant (admin) or current user (warga)
func ListDocumentRequests(c echo.Context) error {
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
	documentType := c.QueryParam("document_type")

	offset := (page - 1) * limit

	// Check if user is admin (can see all requests)
	var isAdmin bool
	err := db.DB.Get(&isAdmin, `
		SELECT EXISTS(
			SELECT 1 FROM tenant_users tu
			JOIN roles r ON tu.role_id = r.id
			WHERE tu.user_id = $1 AND tu.tenant_id = $2
			AND (r.name = 'Admin' OR r.name = 'Bendahara')
			AND tu.deleted_at IS NULL AND r.deleted_at IS NULL
		)
	`, userID, tenantID)

	// Build query
	query := `
		SELECT 
			d.id, d.tenant_id, d.user_id, d.document_type, d.purpose, d.status,
			d.approved_by, d.approved_at, d.rejected_reason, 
			CASE WHEN d.attachment_ids IS NULL THEN ARRAY[]::TEXT[] ELSE ARRAY(SELECT unnest(d.attachment_ids)::TEXT) END,
			d.notes, d.created_at, d.updated_at,
			u.full_name as user_name, u.email as user_email,
			au.full_name as approved_by_name
		FROM document_requests d
		LEFT JOIN users u ON d.user_id = u.id
		LEFT JOIN users au ON d.approved_by = au.id
		WHERE d.tenant_id = $1 AND d.deleted_at IS NULL
	`
	args := []interface{}{tenantID}
	argIndex := 2

	// Filter by user if not admin
	if err == nil && !isAdmin {
		query += ` AND d.user_id = $` + strconv.Itoa(argIndex)
		args = append(args, userID)
		argIndex++
	}

	if status != "" {
		query += ` AND d.status = $` + strconv.Itoa(argIndex)
		args = append(args, status)
		argIndex++
	}

	if documentType != "" {
		query += ` AND d.document_type = $` + strconv.Itoa(argIndex)
		args = append(args, documentType)
		argIndex++
	}

	query += ` ORDER BY d.created_at DESC LIMIT $` + strconv.Itoa(argIndex) + ` OFFSET $` + strconv.Itoa(argIndex+1)
	args = append(args, limit, offset)

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		c.Logger().Errorf("Error executing query: %v, query: %s", err, query)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error: " + err.Error()})
	}
	defer rows.Close()

	requests := []map[string]interface{}{}
	for rows.Next() {
		var d models.DocumentRequest
		var approvedBy, userName, userEmail, approvedByName, rejectedReason, notes sql.NullString
		var approvedAt sql.NullTime
		var attachmentIDs pq.StringArray

		err := rows.Scan(
			&d.ID, &d.TenantID, &d.UserID, &d.DocumentType, &d.Purpose, &d.Status,
			&approvedBy, &approvedAt, &rejectedReason, &attachmentIDs,
			&notes, &d.CreatedAt, &d.UpdatedAt,
			&userName, &userEmail, &approvedByName,
		)
		if err != nil {
			// Log error but continue processing other rows
			c.Logger().Errorf("Error scanning document request row: %v", err)
			continue
		}

		requestData := map[string]interface{}{
			"id":            d.ID,
			"user_id":       d.UserID,
			"document_type": d.DocumentType,
			"purpose":       d.Purpose,
			"status":        d.Status,
			"created_at":    d.CreatedAt.Format(time.RFC3339),
			"updated_at":    d.UpdatedAt.Format(time.RFC3339),
		}

		if userName.Valid {
			requestData["user_name"] = userName.String
		}
		if userEmail.Valid {
			requestData["user_email"] = userEmail.String
		}
		if approvedBy.Valid {
			requestData["approved_by"] = approvedBy.String
		}
		if approvedByName.Valid {
			requestData["approved_by_name"] = approvedByName.String
		}
		if approvedAt.Valid {
			requestData["approved_at"] = approvedAt.Time.Format(time.RFC3339)
		}
		if rejectedReason.Valid {
			requestData["rejected_reason"] = rejectedReason.String
		}
		if notes.Valid {
			requestData["notes"] = notes.String
		}
		if len(attachmentIDs) > 0 {
			requestData["attachment_ids"] = []string(attachmentIDs)
		}

		requests = append(requests, requestData)
	}

	// Get total count
	countQuery := `
		SELECT COUNT(*) FROM document_requests d
		WHERE d.tenant_id = $1 AND d.deleted_at IS NULL
	`
	countArgs := []interface{}{tenantID}
	countArgIndex := 2

	if !isAdmin {
		countQuery += ` AND d.user_id = $` + strconv.Itoa(countArgIndex)
		countArgs = append(countArgs, userID)
		countArgIndex++
	}

	if status != "" {
		countQuery += ` AND d.status = $` + strconv.Itoa(countArgIndex)
		countArgs = append(countArgs, status)
		countArgIndex++
	}

	if documentType != "" {
		countQuery += ` AND d.document_type = $` + strconv.Itoa(countArgIndex)
		countArgs = append(countArgs, documentType)
		countArgIndex++
	}

	var total int
	err = db.DB.Get(&total, countQuery, countArgs...)
	if err != nil {
		total = len(requests)
	}

	totalPages := (total + limit - 1) / limit

	return c.JSON(http.StatusOK, map[string]interface{}{
		"requests": requests,
		"pagination": map[string]interface{}{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": totalPages,
		},
	})
}

// GetDocumentRequest gets a single document request by ID
func GetDocumentRequest(c echo.Context) error {
	tenantID := c.Get(string(middleware.CtxTenantID)).(string)
	userID := c.Get(string(middleware.CtxUserID)).(string)
	requestID := c.Param("request_id")

	// Check if user is admin or owns the request
	var isAdmin bool
	err := db.DB.Get(&isAdmin, `
		SELECT EXISTS(
			SELECT 1 FROM tenant_users tu
			JOIN roles r ON tu.role_id = r.id
			WHERE tu.user_id = $1 AND tu.tenant_id = $2
			AND (r.name = 'Admin' OR r.name = 'Bendahara')
			AND tu.deleted_at IS NULL AND r.deleted_at IS NULL
		)
	`, userID, tenantID)

	query := `
		SELECT 
			d.id, d.tenant_id, d.user_id, d.document_type, d.purpose, d.status,
			d.approved_by, d.approved_at, d.rejected_reason, 
			CASE WHEN d.attachment_ids IS NULL THEN ARRAY[]::TEXT[] ELSE ARRAY(SELECT unnest(d.attachment_ids)::TEXT) END,
			d.notes, d.created_at, d.updated_at,
			u.full_name as user_name, u.email as user_email,
			au.full_name as approved_by_name
		FROM document_requests d
		LEFT JOIN users u ON d.user_id = u.id
		LEFT JOIN users au ON d.approved_by = au.id
		WHERE d.id = $1 AND d.tenant_id = $2 AND d.deleted_at IS NULL
	`
	args := []interface{}{requestID, tenantID}

	if !isAdmin {
		query += ` AND d.user_id = $3`
		args = append(args, userID)
	}

	var d models.DocumentRequest
	var approvedBy, userName, userEmail, approvedByName, rejectedReason, notes sql.NullString
	var approvedAt sql.NullTime
	var attachmentIDs pq.StringArray

	err = db.DB.QueryRow(query, args...).Scan(
		&d.ID, &d.TenantID, &d.UserID, &d.DocumentType, &d.Purpose, &d.Status,
		&approvedBy, &approvedAt, &rejectedReason, &attachmentIDs,
		&notes, &d.CreatedAt, &d.UpdatedAt,
		&userName, &userEmail, &approvedByName,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Document request not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}

	requestData := map[string]interface{}{
		"id":            d.ID,
		"user_id":       d.UserID,
		"document_type": d.DocumentType,
		"purpose":       d.Purpose,
		"status":        d.Status,
		"created_at":    d.CreatedAt.Format(time.RFC3339),
		"updated_at":    d.UpdatedAt.Format(time.RFC3339),
	}

	if userName.Valid {
		requestData["user_name"] = userName.String
	}
	if userEmail.Valid {
		requestData["user_email"] = userEmail.String
	}
	if approvedBy.Valid {
		requestData["approved_by"] = approvedBy.String
	}
	if approvedByName.Valid {
		requestData["approved_by_name"] = approvedByName.String
	}
	if approvedAt.Valid {
		requestData["approved_at"] = approvedAt.Time.Format(time.RFC3339)
	}
	if rejectedReason.Valid {
		requestData["rejected_reason"] = rejectedReason.String
	}
	if notes.Valid {
		requestData["notes"] = notes.String
	}
	if len(attachmentIDs) > 0 {
		requestData["attachment_ids"] = attachmentIDs
	}

	return c.JSON(http.StatusOK, requestData)
}

// CreateDocumentRequest creates a new document request
func CreateDocumentRequest(c echo.Context) error {
	tenantID := c.Get(string(middleware.CtxTenantID)).(string)
	userID := c.Get(string(middleware.CtxUserID)).(string)

	var req models.CreateDocumentRequestRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// Validate required fields
	if req.DocumentType == "" || req.Purpose == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Document type and purpose are required"})
	}

	// Create document request
	requestID := uuid.New().String()
	query := `
		INSERT INTO document_requests (id, tenant_id, user_id, document_type, purpose, attachment_ids)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at
	`

	var createdAt time.Time
	var returnedID string
	err := db.DB.QueryRow(query,
		requestID, tenantID, userID, req.DocumentType, req.Purpose, pq.Array(req.AttachmentIDs),
	).Scan(&returnedID, &createdAt)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create document request: " + err.Error()})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"id":         returnedID,
		"created_at": createdAt.Format(time.RFC3339),
	})
}

// UpdateDocumentRequest updates a document request (admin only for approval)
func UpdateDocumentRequest(c echo.Context) error {
	tenantID := c.Get(string(middleware.CtxTenantID)).(string)
	userID := c.Get(string(middleware.CtxUserID)).(string)
	requestID := c.Param("request_id")

	var req models.UpdateDocumentRequestRequest
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
			AND (r.name = 'Admin' OR r.name = 'Bendahara')
			AND tu.deleted_at IS NULL AND r.deleted_at IS NULL
		)
	`, userID, tenantID)

	if !isAdmin {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "Only admin can update document requests"})
	}

	// Build update query dynamically
	updates := []string{"updated_at = CURRENT_TIMESTAMP"}
	args := []interface{}{requestID, tenantID}
	argIndex := 3

	if req.Status != nil {
		updates = append(updates, "status = $"+strconv.Itoa(argIndex))
		args = append(args, *req.Status)
		argIndex++

		// If status is approved, set approved_by and approved_at
		if *req.Status == "approved" {
			updates = append(updates, "approved_by = $"+strconv.Itoa(argIndex))
			args = append(args, userID)
			argIndex++
			updates = append(updates, "approved_at = CURRENT_TIMESTAMP")
		}
	}

	if req.RejectedReason != nil {
		updates = append(updates, "rejected_reason = $"+strconv.Itoa(argIndex))
		args = append(args, *req.RejectedReason)
		argIndex++
	}

	if req.Notes != nil {
		updates = append(updates, "notes = $"+strconv.Itoa(argIndex))
		args = append(args, *req.Notes)
		argIndex++
	}

	if len(updates) == 1 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "No fields to update"})
	}

	query := `UPDATE document_requests SET ` + updates[0]
	for i := 1; i < len(updates); i++ {
		query += `, ` + updates[i]
	}
	query += ` WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL RETURNING id`

	var returnedID string
	err = db.DB.QueryRow(query, args...).Scan(&returnedID)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Document request not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update document request: " + err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Document request updated successfully"})
}

// DeleteDocumentRequest soft deletes a document request
func DeleteDocumentRequest(c echo.Context) error {
	tenantID := c.Get(string(middleware.CtxTenantID)).(string)
	userID := c.Get(string(middleware.CtxUserID)).(string)
	requestID := c.Param("request_id")

	// Check if user is admin or owns the request
	var isAdmin bool
	err := db.DB.Get(&isAdmin, `
		SELECT EXISTS(
			SELECT 1 FROM tenant_users tu
			JOIN roles r ON tu.role_id = r.id
			WHERE tu.user_id = $1 AND tu.tenant_id = $2
			AND (r.name = 'Admin' OR r.name = 'Bendahara')
			AND tu.deleted_at IS NULL AND r.deleted_at IS NULL
		)
	`, userID, tenantID)

	query := `UPDATE document_requests SET deleted_at = CURRENT_TIMESTAMP WHERE id = $1 AND tenant_id = $2`
	args := []interface{}{requestID, tenantID}

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
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Document request not found"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Document request deleted successfully"})
}





