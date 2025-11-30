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

// ListAnnouncements lists all announcements for the tenant
func ListAnnouncements(c echo.Context) error {
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
	priority := c.QueryParam("priority")
	category := c.QueryParam("category")
	search := c.QueryParam("search")

	offset := (page - 1) * limit

	// Build query
	query := `
		SELECT 
			a.id, a.tenant_id, a.author_id, a.title, a.content, a.priority, 
			a.category, a.is_pinned, a.sent_notification, a.sent_whatsapp,
			a.sent_at, a.expires_at, a.metadata, a.created_at, a.updated_at,
			u.full_name as author_name, u.email as author_email,
			CASE WHEN ar.id IS NOT NULL THEN true ELSE false END as is_read
		FROM announcements a
		LEFT JOIN users u ON a.author_id = u.id
		LEFT JOIN announcement_reads ar ON a.id = ar.announcement_id AND ar.user_id = $1
		WHERE a.tenant_id = $2 AND a.deleted_at IS NULL
	`
	args := []interface{}{userID, tenantID}
	argIndex := 3

	if priority != "" {
		query += ` AND a.priority = $` + strconv.Itoa(argIndex)
		args = append(args, priority)
		argIndex++
	}

	if category != "" {
		query += ` AND a.category = $` + strconv.Itoa(argIndex)
		args = append(args, category)
		argIndex++
	}

	if search != "" {
		searchPattern := "%" + search + "%"
		query += ` AND (a.title ILIKE $` + strconv.Itoa(argIndex) + ` OR a.content ILIKE $` + strconv.Itoa(argIndex) + `)`
		args = append(args, searchPattern)
		argIndex++
	}

	// Order by pinned first, then by created_at DESC
	query += ` ORDER BY a.is_pinned DESC, a.created_at DESC LIMIT $` + strconv.Itoa(argIndex) + ` OFFSET $` + strconv.Itoa(argIndex+1)
	args = append(args, limit, offset)

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}
	defer rows.Close()

	announcements := []map[string]interface{}{}
	for rows.Next() {
		var ann models.Announcement
		var category, metadata sql.NullString
		var sentAt, expiresAt sql.NullTime
		var isRead bool

		err := rows.Scan(
			&ann.ID, &ann.TenantID, &ann.AuthorID, &ann.Title, &ann.Content, &ann.Priority,
			&category, &ann.IsPinned, &ann.SentNotification, &ann.SentWhatsApp,
			&sentAt, &expiresAt, &metadata, &ann.CreatedAt, &ann.UpdatedAt,
			&ann.AuthorName, &ann.AuthorEmail, &isRead,
		)
		if err != nil {
			continue
		}

		annData := map[string]interface{}{
			"id":                ann.ID,
			"title":             ann.Title,
			"content":           ann.Content,
			"priority":          ann.Priority,
			"is_pinned":         ann.IsPinned,
			"sent_notification": ann.SentNotification,
			"sent_whatsapp":     ann.SentWhatsApp,
			"created_at":        ann.CreatedAt.Format(time.RFC3339),
			"updated_at":        ann.UpdatedAt.Format(time.RFC3339),
			"is_read":           isRead,
		}

		if category.Valid {
			annData["category"] = category.String
		}
		if sentAt.Valid {
			annData["sent_at"] = sentAt.Time.Format(time.RFC3339)
		}
		if expiresAt.Valid {
			annData["expires_at"] = expiresAt.Time.Format(time.RFC3339)
		}
		if ann.AuthorName.Valid {
			annData["author_name"] = ann.AuthorName.String
		}
		if ann.AuthorEmail.Valid {
			annData["author_email"] = ann.AuthorEmail.String
		}

		announcements = append(announcements, annData)
	}

	// Get total count
	var total int
	countQuery := `SELECT COUNT(*) FROM announcements WHERE tenant_id = $1 AND deleted_at IS NULL`
	if priority != "" {
		countQuery += ` AND priority = $2`
		err = db.DB.Get(&total, countQuery, tenantID, priority)
	} else {
		err = db.DB.Get(&total, countQuery, tenantID)
	}
	if err != nil {
		total = len(announcements)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"announcements": announcements,
		"pagination": map[string]interface{}{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": (total + limit - 1) / limit,
		},
	})
}

// GetAnnouncement gets an announcement by ID
func GetAnnouncement(c echo.Context) error {
	tenantID := c.Get(string(middleware.CtxTenantID)).(string)
	userID := c.Get(string(middleware.CtxUserID)).(string)
	announcementID := c.Param("announcement_id")

	var ann models.Announcement
	var category, metadata sql.NullString
	var sentAt, expiresAt sql.NullTime
	var isRead bool

	err := db.DB.QueryRow(`
		SELECT 
			a.id, a.tenant_id, a.author_id, a.title, a.content, a.priority,
			a.category, a.is_pinned, a.sent_notification, a.sent_whatsapp,
			a.sent_at, a.expires_at, a.metadata, a.created_at, a.updated_at,
			u.full_name as author_name, u.email as author_email,
			CASE WHEN ar.id IS NOT NULL THEN true ELSE false END as is_read
		FROM announcements a
		LEFT JOIN users u ON a.author_id = u.id
		LEFT JOIN announcement_reads ar ON a.id = ar.announcement_id AND ar.user_id = $1
		WHERE a.id = $2 AND a.tenant_id = $3 AND a.deleted_at IS NULL
	`, userID, announcementID, tenantID).Scan(
		&ann.ID, &ann.TenantID, &ann.AuthorID, &ann.Title, &ann.Content, &ann.Priority,
		&category, &ann.IsPinned, &ann.SentNotification, &ann.SentWhatsApp,
		&sentAt, &expiresAt, &metadata, &ann.CreatedAt, &ann.UpdatedAt,
		&ann.AuthorName, &ann.AuthorEmail, &isRead,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Announcement not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}

	// Mark as read if not already read
	if !isRead {
		_, _ = db.DB.Exec(`
			INSERT INTO announcement_reads (announcement_id, user_id)
			VALUES ($1, $2)
			ON CONFLICT (announcement_id, user_id) DO NOTHING
		`, announcementID, userID)
	}

	annData := map[string]interface{}{
		"id":                ann.ID,
		"title":             ann.Title,
		"content":           ann.Content,
		"priority":          ann.Priority,
		"is_pinned":         ann.IsPinned,
		"sent_notification": ann.SentNotification,
		"sent_whatsapp":     ann.SentWhatsApp,
		"created_at":        ann.CreatedAt.Format(time.RFC3339),
		"updated_at":        ann.UpdatedAt.Format(time.RFC3339),
		"is_read":           true, // Marked as read
	}

	if category.Valid {
		annData["category"] = category.String
	}
	if sentAt.Valid {
		annData["sent_at"] = sentAt.Time.Format(time.RFC3339)
	}
	if expiresAt.Valid {
		annData["expires_at"] = expiresAt.Time.Format(time.RFC3339)
	}
	if ann.AuthorName.Valid {
		annData["author_name"] = ann.AuthorName.String
	}
	if ann.AuthorEmail.Valid {
		annData["author_email"] = ann.AuthorEmail.String
	}

	return c.JSON(http.StatusOK, annData)
}

// CreateAnnouncement creates a new announcement (admin/sekretariat only)
func CreateAnnouncement(c echo.Context) error {
	tenantID := c.Get(string(middleware.CtxTenantID)).(string)
	userID := c.Get(string(middleware.CtxUserID)).(string)

	req := new(models.CreateAnnouncementRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// Validate required fields
	if req.Title == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "title is required"})
	}
	if req.Content == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "content is required"})
	}

	priority := "medium"
	if req.Priority != "" {
		priority = req.Priority
	}

	isPinned := false
	if req.IsPinned != nil {
		isPinned = *req.IsPinned
	}

	// Parse expires_at if provided
	var expiresAt sql.NullTime
	if req.ExpiresAt != nil && *req.ExpiresAt != "" {
		parsedTime, err := time.Parse(time.RFC3339, *req.ExpiresAt)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid expires_at format. Use ISO 8601"})
		}
		expiresAt = sql.NullTime{Time: parsedTime, Valid: true}
	}

	// Create announcement
	announcementID := uuid.New().String()
	query := `
		INSERT INTO announcements (id, tenant_id, author_id, title, content, priority, category, is_pinned, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at
	`

	var createdAt time.Time
	var returnedID string
	
	// Handle category - can be nil
	var categoryValue interface{}
	if req.Category != nil && *req.Category != "" {
		categoryValue = *req.Category
	} else {
		categoryValue = nil
	}
	
	err := db.DB.QueryRow(query,
		announcementID, tenantID, userID, req.Title, req.Content, priority,
		categoryValue, isPinned, expiresAt,
	).Scan(&returnedID, &createdAt)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create announcement: " + err.Error()})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"id":         returnedID,
		"message":    "Announcement created successfully",
		"created_at": createdAt.Format(time.RFC3339),
	})
}

// UpdateAnnouncement updates an announcement
func UpdateAnnouncement(c echo.Context) error {
	tenantID := c.Get(string(middleware.CtxTenantID)).(string)
	announcementID := c.Param("announcement_id")

	req := new(models.UpdateAnnouncementRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// Check if announcement exists and belongs to tenant
	var exists bool
	err := db.DB.Get(&exists, `
		SELECT EXISTS(
			SELECT 1 FROM announcements 
			WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL
		)
	`, announcementID, tenantID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}
	if !exists {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Announcement not found"})
	}

	// Build update query dynamically
	updates := []string{}
	args := []interface{}{}
	argIndex := 1

	if req.Title != nil {
		updates = append(updates, "title = $"+strconv.Itoa(argIndex))
		args = append(args, *req.Title)
		argIndex++
	}
	if req.Content != nil {
		updates = append(updates, "content = $"+strconv.Itoa(argIndex))
		args = append(args, *req.Content)
		argIndex++
	}
	if req.Priority != nil {
		updates = append(updates, "priority = $"+strconv.Itoa(argIndex))
		args = append(args, *req.Priority)
		argIndex++
	}
	if req.Category != nil {
		if *req.Category == "" {
			updates = append(updates, "category = NULL")
		} else {
			updates = append(updates, "category = $"+strconv.Itoa(argIndex))
			args = append(args, *req.Category)
			argIndex++
		}
	}
	if req.IsPinned != nil {
		updates = append(updates, "is_pinned = $"+strconv.Itoa(argIndex))
		args = append(args, *req.IsPinned)
		argIndex++
	}
	if req.ExpiresAt != nil {
		if *req.ExpiresAt == "" {
			updates = append(updates, "expires_at = NULL")
		} else {
			parsedTime, err := time.Parse(time.RFC3339, *req.ExpiresAt)
			if err != nil {
				return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid expires_at format. Use ISO 8601"})
			}
			updates = append(updates, "expires_at = $"+strconv.Itoa(argIndex))
			args = append(args, parsedTime)
			argIndex++
		}
	}

	if len(updates) == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "No fields to update"})
	}

	// Add WHERE clause
	updates = append(updates, "updated_at = CURRENT_TIMESTAMP")
	query := `UPDATE announcements SET ` + updates[0]
	for i := 1; i < len(updates); i++ {
		query += `, ` + updates[i]
	}
	query += ` WHERE id = $` + strconv.Itoa(argIndex) + ` AND tenant_id = $` + strconv.Itoa(argIndex+1)
	args = append(args, announcementID, tenantID)

	_, err = db.DB.Exec(query, args...)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update announcement: " + err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Announcement updated successfully",
		"id":      announcementID,
	})
}

// DeleteAnnouncement deletes an announcement (soft delete)
func DeleteAnnouncement(c echo.Context) error {
	tenantID := c.Get(string(middleware.CtxTenantID)).(string)
	announcementID := c.Param("announcement_id")

	// Check if announcement exists and belongs to tenant
	var exists bool
	err := db.DB.Get(&exists, `
		SELECT EXISTS(
			SELECT 1 FROM announcements 
			WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL
		)
	`, announcementID, tenantID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}
	if !exists {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Announcement not found"})
	}

	// Soft delete
	_, err = db.DB.Exec(`
		UPDATE announcements 
		SET deleted_at = CURRENT_TIMESTAMP 
		WHERE id = $1 AND tenant_id = $2
	`, announcementID, tenantID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete announcement: " + err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Announcement deleted successfully",
	})
}
