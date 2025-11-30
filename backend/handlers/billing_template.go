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

// ListBillingTemplates lists all billing templates for the tenant with pagination
func ListBillingTemplates(c echo.Context) error {
	tenantID := c.Get(string(middleware.CtxTenantID)).(string)

	// Parse query parameters
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	search := c.QueryParam("search")
	category := c.QueryParam("category")
	isActive := c.QueryParam("is_active")

	// Build query
	query := `
		SELECT id, tenant_id, name, category, type, description, amount, late_fee, 
		       due_day, recurring_type, late_fee_type, late_fee_percentage, late_fee_max, 
		       is_active, is_system, created_by, created_at, updated_at
		FROM billing_templates
		WHERE tenant_id = $1 AND deleted_at IS NULL
	`
	args := []interface{}{tenantID}
	argIndex := 2

	if search != "" {
		query += ` AND (name ILIKE $` + strconv.Itoa(argIndex) + ` OR category ILIKE $` + strconv.Itoa(argIndex) + `)`
		searchPattern := "%" + search + "%"
		args = append(args, searchPattern)
		argIndex++
	}

	if category != "" {
		query += ` AND category = $` + strconv.Itoa(argIndex)
		args = append(args, category)
		argIndex++
	}

	if isActive != "" {
		activeBool := isActive == "true"
		query += ` AND is_active = $` + strconv.Itoa(argIndex)
		args = append(args, activeBool)
		argIndex++
	}

	query += ` ORDER BY is_system DESC, name ASC LIMIT $` + strconv.Itoa(argIndex) + ` OFFSET $` + strconv.Itoa(argIndex+1)
	args = append(args, limit, offset)

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		c.Logger().Errorf("Error querying billing templates: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error: " + err.Error()})
	}
	defer rows.Close()

	var templates []map[string]interface{}
	for rows.Next() {
		var template models.BillingTemplate
		err := rows.Scan(
			&template.ID, &template.TenantID, &template.Name, &template.Category,
			&template.Type, &template.Description, &template.Amount, &template.LateFee,
			&template.DueDay, &template.RecurringType, &template.LateFeeType,
			&template.LateFeePercentage, &template.LateFeeMax,
			&template.IsActive, &template.IsSystem, &template.CreatedBy, 
			&template.CreatedAt, &template.UpdatedAt,
		)
		if err != nil {
			c.Logger().Errorf("Error scanning template row: %v", err)
			continue
		}

		templateData := map[string]interface{}{
			"id":             template.ID,
			"name":           template.Name,
			"category":       template.Category,
			"type":           template.Type,
			"amount":         template.Amount,
			"late_fee":       template.LateFee,
			"recurring_type": template.RecurringType,
			"late_fee_type":  template.LateFeeType,
			"is_active":      template.IsActive,
			"is_system":      template.IsSystem,
			"created_at":     template.CreatedAt,
		}

		if template.Description.Valid {
			templateData["description"] = template.Description.String
		}
		if template.DueDay.Valid {
			templateData["due_day"] = template.DueDay.Int64
		}
		if template.LateFeePercentage.Valid {
			templateData["late_fee_percentage"] = template.LateFeePercentage.Float64
		}
		if template.LateFeeMax.Valid {
			templateData["late_fee_max"] = template.LateFeeMax.Float64
		}

		// Get amount rules for this template
		amountRules, err := getAmountRules(template.ID)
		if err == nil && len(amountRules) > 0 {
			templateData["amount_rules"] = amountRules
		}

		templates = append(templates, templateData)
	}

	// Get total count
	countQuery := `SELECT COUNT(*) FROM billing_templates WHERE tenant_id = $1 AND deleted_at IS NULL`
	countArgs := []interface{}{tenantID}
	countArgIndex := 2

	if search != "" {
		searchPattern := "%" + search + "%"
		countQuery += ` AND (name ILIKE $` + strconv.Itoa(countArgIndex) + ` OR category ILIKE $` + strconv.Itoa(countArgIndex) + `)`
		countArgs = append(countArgs, searchPattern)
		countArgIndex++
	}

	if category != "" {
		countQuery += ` AND category = $` + strconv.Itoa(countArgIndex)
		countArgs = append(countArgs, category)
		countArgIndex++
	}

	if isActive != "" {
		activeBool := isActive == "true"
		countQuery += ` AND is_active = $` + strconv.Itoa(countArgIndex)
		countArgs = append(countArgs, activeBool)
		countArgIndex++
	}

	var total int
	err = db.DB.Get(&total, countQuery, countArgs...)
	if err != nil {
		total = len(templates)
	}

	totalPages := (total + limit - 1) / limit

	return c.JSON(http.StatusOK, map[string]interface{}{
		"templates": templates,
		"pagination": map[string]interface{}{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": totalPages,
		},
	})
}

// GetBillingTemplate gets a billing template by ID
func GetBillingTemplate(c echo.Context) error {
	tenantID := c.Get(string(middleware.CtxTenantID)).(string)
	templateID := c.Param("template_id")

	var template models.BillingTemplate
	err := db.DB.QueryRow(`
		SELECT id, tenant_id, name, category, type, description, amount, late_fee, 
		       due_day, recurring_type, late_fee_type, late_fee_percentage, late_fee_max, 
		       is_active, is_system, created_by, created_at, updated_at
		FROM billing_templates
		WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL
	`, templateID, tenantID).Scan(
		&template.ID, &template.TenantID, &template.Name, &template.Category,
		&template.Type, &template.Description, &template.Amount, &template.LateFee,
		&template.DueDay, &template.RecurringType, &template.LateFeeType,
		&template.LateFeePercentage, &template.LateFeeMax,
		&template.IsActive, &template.IsSystem, &template.CreatedBy, 
		&template.CreatedAt, &template.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Template not found"})
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error: " + err.Error()})
	}

	templateData := map[string]interface{}{
		"id":             template.ID,
		"name":           template.Name,
		"category":       template.Category,
		"type":           template.Type,
		"amount":         template.Amount,
		"late_fee":       template.LateFee,
		"recurring_type": template.RecurringType,
		"late_fee_type":  template.LateFeeType,
		"is_active":      template.IsActive,
		"is_system":      template.IsSystem,
		"created_at":     template.CreatedAt,
		"updated_at":     template.UpdatedAt,
	}

	if template.Description.Valid {
		templateData["description"] = template.Description.String
	}
	if template.DueDay.Valid {
		templateData["due_day"] = template.DueDay.Int64
	}
	if template.LateFeePercentage.Valid {
		templateData["late_fee_percentage"] = template.LateFeePercentage.Float64
	}
	if template.LateFeeMax.Valid {
		templateData["late_fee_max"] = template.LateFeeMax.Float64
	}

	// Get amount rules
	amountRules, err := getAmountRules(templateID)
	if err == nil {
		templateData["amount_rules"] = amountRules
	}

	return c.JSON(http.StatusOK, templateData)
}

// getAmountRules retrieves amount rules for a template
func getAmountRules(templateID string) ([]map[string]interface{}, error) {
	rows, err := db.DB.Query(`
		SELECT unit_type, amount
		FROM billing_template_amount_rules
		WHERE template_id = $1
		ORDER BY unit_type
	`, templateID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rules []map[string]interface{}
	for rows.Next() {
		var unitType string
		var amount float64
		if err := rows.Scan(&unitType, &amount); err != nil {
			continue
		}
		rules = append(rules, map[string]interface{}{
			"unit_type": unitType,
			"amount":    amount,
		})
	}
	return rules, nil
}

// CreateBillingTemplate creates a new billing template
func CreateBillingTemplate(c echo.Context) error {
	tenantID := c.Get(string(middleware.CtxTenantID)).(string)
	userID := c.Get(string(middleware.CtxUserID)).(string)

	req := new(models.CreateBillingTemplateRequest)
	if err := c.Bind(req); err != nil {
		c.Logger().Errorf("Error binding request: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request: " + err.Error()})
	}
	
	// Validate required fields
	if req.Name == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "name is required"})
	}
	if req.Category == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "category is required"})
	}
	if req.Type == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "type is required"})
	}

	// Check if template name already exists in tenant
	var exists bool
	err := db.DB.Get(&exists, `
		SELECT EXISTS(
			SELECT 1 FROM billing_templates 
			WHERE tenant_id = $1 AND name = $2 AND deleted_at IS NULL
		)
	`, tenantID, req.Name)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}
	if exists {
		return c.JSON(http.StatusConflict, map[string]string{"error": "Template name already exists"})
	}

	lateFee := 0.0
	if req.LateFee != nil {
		lateFee = *req.LateFee
	}

	dueDay := sql.NullInt64{Valid: false}
	if req.DueDay != nil {
		dueDay = sql.NullInt64{Int64: int64(*req.DueDay), Valid: true}
	}

	recurringType := "one-time"
	if req.RecurringType != nil {
		recurringType = *req.RecurringType
	}

	lateFeeType := "fixed"
	if req.LateFeeType != nil {
		lateFeeType = *req.LateFeeType
	}

	lateFeePercentage := sql.NullFloat64{Valid: false}
	if req.LateFeePercentage != nil {
		lateFeePercentage = sql.NullFloat64{Float64: *req.LateFeePercentage, Valid: true}
	}

	lateFeeMax := sql.NullFloat64{Valid: false}
	if req.LateFeeMax != nil {
		lateFeeMax = sql.NullFloat64{Float64: *req.LateFeeMax, Valid: true}
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	// Start transaction
	tx, err := db.DB.Beginx()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to start transaction"})
	}
	defer tx.Rollback()

	// Handle description (nullable)
	description := sql.NullString{Valid: false}
	if req.Description != nil && *req.Description != "" {
		description = sql.NullString{String: *req.Description, Valid: true}
	}

	// Create template
	templateID := uuid.New().String()
	query := `INSERT INTO billing_templates 
	          (id, tenant_id, name, category, type, description, amount, late_fee, 
	           due_day, recurring_type, late_fee_type, late_fee_percentage, late_fee_max, 
	           is_active, is_system, created_by)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, false, $15)
	          RETURNING id, created_at`
	
	var createdAt time.Time
	err = tx.QueryRow(query, templateID, tenantID, req.Name, req.Category, req.Type, description, 
		req.Amount, lateFee, dueDay, recurringType, lateFeeType, lateFeePercentage, lateFeeMax, 
		isActive, userID).Scan(&templateID, &createdAt)
	if err != nil {
		c.Logger().Errorf("Error creating billing template: %v, query: %s", err, query)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create template: " + err.Error()})
	}

	// Create amount rules if provided
	if len(req.AmountRules) > 0 {
		for _, rule := range req.AmountRules {
			// Skip rules with invalid amount
			if rule.Amount <= 0 {
				continue
			}
			ruleID := uuid.New().String()
			_, err = tx.Exec(`
				INSERT INTO billing_template_amount_rules (id, template_id, unit_type, amount)
				VALUES ($1, $2, $3, $4)
			`, ruleID, templateID, rule.UnitType, rule.Amount)
			if err != nil {
				c.Logger().Errorf("Error creating amount rule: %v", err)
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create amount rule: " + err.Error()})
			}
		}
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to commit transaction"})
	}

	// Return created template
	c.SetParamNames("template_id")
	c.SetParamValues(templateID)
	return GetBillingTemplate(c)
}

// UpdateBillingTemplate updates a billing template
func UpdateBillingTemplate(c echo.Context) error {
	tenantID := c.Get(string(middleware.CtxTenantID)).(string)
	templateID := c.Param("template_id")

	req := new(models.UpdateBillingTemplateRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// Check if template exists and belongs to tenant
	var exists bool
	var isSystem bool
	err := db.DB.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM billing_templates 
			WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL
		), is_system
		FROM billing_templates
		WHERE id = $1 AND tenant_id = $2
	`, templateID, tenantID).Scan(&exists, &isSystem)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}
	if !exists {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Template not found"})
	}
	if isSystem {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "Cannot update system template"})
	}

	// Build update query dynamically
	updates := []string{}
	args := []interface{}{}
	argIndex := 1

	if req.Name != nil {
		// Check if new name conflicts with existing template
		var nameExists bool
		err = db.DB.Get(&nameExists, `
			SELECT EXISTS(
				SELECT 1 FROM billing_templates 
				WHERE tenant_id = $1 AND name = $2 AND id != $3 AND deleted_at IS NULL
			)
		`, tenantID, *req.Name, templateID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
		}
		if nameExists {
			return c.JSON(http.StatusConflict, map[string]string{"error": "Template name already exists"})
		}

		updates = append(updates, "name = $"+strconv.Itoa(argIndex))
		args = append(args, *req.Name)
		argIndex++
	}
	if req.Category != nil {
		updates = append(updates, "category = $"+strconv.Itoa(argIndex))
		args = append(args, *req.Category)
		argIndex++
	}
	if req.Type != nil {
		updates = append(updates, "type = $"+strconv.Itoa(argIndex))
		args = append(args, *req.Type)
		argIndex++
	}
	if req.Description != nil {
		updates = append(updates, "description = $"+strconv.Itoa(argIndex))
		args = append(args, *req.Description)
		argIndex++
	}
	if req.Amount != nil {
		updates = append(updates, "amount = $"+strconv.Itoa(argIndex))
		args = append(args, *req.Amount)
		argIndex++
	}
	if req.LateFee != nil {
		updates = append(updates, "late_fee = $"+strconv.Itoa(argIndex))
		args = append(args, *req.LateFee)
		argIndex++
	}
	if req.DueDay != nil {
		updates = append(updates, "due_day = $"+strconv.Itoa(argIndex))
		args = append(args, *req.DueDay)
		argIndex++
	}
	if req.RecurringType != nil {
		updates = append(updates, "recurring_type = $"+strconv.Itoa(argIndex))
		args = append(args, *req.RecurringType)
		argIndex++
	}
	if req.LateFeeType != nil {
		updates = append(updates, "late_fee_type = $"+strconv.Itoa(argIndex))
		args = append(args, *req.LateFeeType)
		argIndex++
	}
	if req.LateFeePercentage != nil {
		updates = append(updates, "late_fee_percentage = $"+strconv.Itoa(argIndex))
		args = append(args, *req.LateFeePercentage)
		argIndex++
	}
	if req.LateFeeMax != nil {
		updates = append(updates, "late_fee_max = $"+strconv.Itoa(argIndex))
		args = append(args, *req.LateFeeMax)
		argIndex++
	}
	if req.IsActive != nil {
		updates = append(updates, "is_active = $"+strconv.Itoa(argIndex))
		args = append(args, *req.IsActive)
		argIndex++
	}

	if len(updates) == 0 && req.AmountRules == nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "No fields to update"})
	}

	// Start transaction
	tx, err := db.DB.Beginx()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to start transaction"})
	}
	defer tx.Rollback()

	// Update template fields if any
	if len(updates) > 0 {
		updates = append(updates, "updated_at = NOW()")
		args = append(args, templateID, tenantID)

		query := "UPDATE billing_templates SET " + updates[0]
		for i := 1; i < len(updates); i++ {
			query += ", " + updates[i]
		}
		query += " WHERE id = $" + strconv.Itoa(argIndex) + " AND tenant_id = $" + strconv.Itoa(argIndex+1)

		_, err = tx.Exec(query, args...)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update template: " + err.Error()})
		}
	}

	// Update amount rules if provided
	if req.AmountRules != nil {
		// Delete existing rules
		_, err = tx.Exec(`
			DELETE FROM billing_template_amount_rules
			WHERE template_id = $1
		`, templateID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete existing amount rules: " + err.Error()})
		}

		// Insert new rules
		for _, rule := range *req.AmountRules {
			ruleID := uuid.New().String()
			_, err = tx.Exec(`
				INSERT INTO billing_template_amount_rules (id, template_id, unit_type, amount)
				VALUES ($1, $2, $3, $4)
			`, ruleID, templateID, rule.UnitType, rule.Amount)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create amount rule: " + err.Error()})
			}
		}
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to commit transaction"})
	}

	return GetBillingTemplate(c)
}

// DeleteBillingTemplate deletes a billing template (soft delete)
func DeleteBillingTemplate(c echo.Context) error {
	tenantID := c.Get(string(middleware.CtxTenantID)).(string)
	templateID := c.Param("template_id")

	// Check if template exists and is not system template
	var exists bool
	var isSystem bool
	err := db.DB.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM billing_templates 
			WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL
		), is_system
		FROM billing_templates
		WHERE id = $1 AND tenant_id = $2
	`, templateID, tenantID).Scan(&exists, &isSystem)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}
	if !exists {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Template not found"})
	}
	if isSystem {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "Cannot delete system template"})
	}

	// Soft delete
	_, err = db.DB.Exec(`
		UPDATE billing_templates 
		SET deleted_at = NOW(), updated_at = NOW()
		WHERE id = $1 AND tenant_id = $2
	`, templateID, tenantID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete template: " + err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Template deleted successfully",
	})
}

