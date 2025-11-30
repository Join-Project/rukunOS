package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"
	"rukunos-backend/db"
	"rukunos-backend/middleware"
	"rukunos-backend/models"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// GenerateBillsFromTemplate generates bills from a template for specified units
func GenerateBillsFromTemplate(c echo.Context) error {
	tenantID := c.Get(string(middleware.CtxTenantID)).(string)
	userID := c.Get(string(middleware.CtxUserID)).(string)

	req := new(models.GenerateBillsFromTemplateRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// Validate template exists and belongs to tenant
	var template models.BillingTemplate
	err := db.DB.QueryRow(`
		SELECT id, tenant_id, name, category, type, amount, due_day, recurring_type, is_active
		FROM billing_templates
		WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL
	`, req.TemplateID, tenantID).Scan(
		&template.ID, &template.TenantID, &template.Name, &template.Category,
		&template.Type, &template.Amount, &template.DueDay, &template.RecurringType, &template.IsActive,
	)

	if err == sql.ErrNoRows {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Template not found"})
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}

	if !template.IsActive {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Template is not active"})
	}

	// Get units to generate bills for
	var unitRows *sql.Rows
	if len(req.UnitIDs) > 0 {
		// Generate for specific units - build IN clause manually
		placeholders := ""
		args := []interface{}{tenantID}
		for i, unitID := range req.UnitIDs {
			if i > 0 {
				placeholders += ","
			}
			placeholders += fmt.Sprintf("$%d", len(args)+1)
			args = append(args, unitID)
		}
		query := fmt.Sprintf(`
			SELECT id, code, type
			FROM units
			WHERE tenant_id = $1 AND id IN (%s) AND deleted_at IS NULL
		`, placeholders)
		unitRows, err = db.DB.Query(query, args...)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
		}
	} else {
		// Generate for all units
		unitRows, err = db.DB.Query(`
			SELECT id, code, type
			FROM units
			WHERE tenant_id = $1 AND deleted_at IS NULL
		`, tenantID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
		}
	}
	defer unitRows.Close()

	// Calculate due date based on period and due_day
	dueDate := calculateDueDate(req.Period, template.DueDay)

	// Start transaction
	tx, err := db.DB.Beginx()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to start transaction"})
	}
	defer tx.Rollback()

	generatedCount := 0
	errors := []string{}

	for unitRows.Next() {
		var unitID, unitCode, unitType string
		if err := unitRows.Scan(&unitID, &unitCode, &unitType); err != nil {
			errors = append(errors, fmt.Sprintf("Error scanning unit: %v", err))
			continue
		}

		// Get amount for this unit type from amount rules
		amount := template.Amount
		// Check if there's a specific amount rule for this unit type
		var ruleAmount sql.NullFloat64
		err := tx.Get(&ruleAmount, `
			SELECT amount
			FROM billing_template_amount_rules
			WHERE template_id = $1 AND unit_type = $2
		`, req.TemplateID, unitType)
		if err == nil && ruleAmount.Valid {
			amount = ruleAmount.Float64
		}

		// Check if bill already exists for this period
		var exists bool
		err = tx.Get(&exists, `
			SELECT EXISTS(
				SELECT 1 FROM bills
				WHERE tenant_id = $1 AND unit_id = $2 AND category = $3 AND period = $4 AND deleted_at IS NULL
			)
		`, tenantID, unitID, template.Category, req.Period)
		if err != nil {
			errors = append(errors, fmt.Sprintf("Error checking existing bill for unit %s: %v", unitCode, err))
			continue
		}
		if exists {
			errors = append(errors, fmt.Sprintf("Bill already exists for unit %s for period %s", unitCode, req.Period))
			continue
		}

		// Create bill
		billID := uuid.New().String()
		_, err = tx.Exec(`
			INSERT INTO bills (id, tenant_id, unit_id, category, period, amount, late_fee, due_date, status, created_by)
			VALUES ($1, $2, $3, $4, $5, $6, 0, $7, 'pending', $8)
		`, billID, tenantID, unitID, template.Category, req.Period, amount, dueDate, userID)
		if err != nil {
			errors = append(errors, fmt.Sprintf("Error creating bill for unit %s: %v", unitCode, err))
			continue
		}

		generatedCount++
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to commit transaction"})
	}

	response := map[string]interface{}{
		"message":        fmt.Sprintf("Generated %d bills successfully", generatedCount),
		"generated_count": generatedCount,
	}
	if len(errors) > 0 {
		response["errors"] = errors
	}

	return c.JSON(http.StatusOK, response)
}

// calculateDueDate calculates the due date based on period and due_day
func calculateDueDate(period string, dueDay sql.NullInt64) time.Time {
	// Parse period (YYYY-MM)
	periodTime, err := time.Parse("2006-01", period)
	if err != nil {
		// Default to current month if parsing fails
		periodTime = time.Now()
	}

	day := 1
	if dueDay.Valid && dueDay.Int64 >= 1 && dueDay.Int64 <= 31 {
		day = int(dueDay.Int64)
	}

	// Get last day of month
	lastDay := time.Date(periodTime.Year(), periodTime.Month()+1, 0, 0, 0, 0, 0, periodTime.Location()).Day()
	if day > lastDay {
		day = lastDay
	}

	return time.Date(periodTime.Year(), periodTime.Month(), day, 0, 0, 0, 0, periodTime.Location())
}

