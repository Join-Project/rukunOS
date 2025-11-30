package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"rukunos-backend/db"
	"rukunos-backend/middleware"
	"rukunos-backend/models"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// CreateUnit creates a new unit
func CreateUnit(c echo.Context) error {
	tenantID := c.Get(string(middleware.CtxTenantID)).(string)

	req := new(models.CreateUnitRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// Validate request
	if req.Code == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Code is required"})
	}
	if req.Type == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Type is required"})
	}

	// Check if unit code already exists in tenant
	var exists bool
	err := db.DB.Get(&exists, `
		SELECT EXISTS(
			SELECT 1 FROM units 
			WHERE tenant_id = $1 AND code = $2 AND deleted_at IS NULL
		)
	`, tenantID, req.Code)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}
	if exists {
		return c.JSON(http.StatusConflict, map[string]string{"error": "Unit code already exists"})
	}

	// Create unit
	unitID := uuid.New().String()
	query := `INSERT INTO units (id, tenant_id, code, type, owner_name, owner_phone, owner_email, address, status)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, 'active')
	          RETURNING id, tenant_id, code, type, owner_name, owner_phone, owner_email, address, status, created_at, updated_at`
	
	var unit models.Unit
	err = db.DB.QueryRow(query, unitID, tenantID, req.Code, req.Type, req.OwnerName, req.OwnerPhone, req.OwnerEmail, req.Address).Scan(
		&unit.ID, &unit.TenantID, &unit.Code, &unit.Type, &unit.OwnerName, &unit.OwnerPhone, 
		&unit.OwnerEmail, &unit.Address, &unit.Status, &unit.CreatedAt, &unit.UpdatedAt)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create unit: " + err.Error()})
	}

	return c.JSON(http.StatusCreated, unit)
}

// ListUnits lists all units with pagination and filters
func ListUnits(c echo.Context) error {
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

	unitType := c.QueryParam("type")
	search := c.QueryParam("search")

	// Build query
	query := `SELECT id, tenant_id, code, type, owner_name, owner_phone, owner_email, address, status, created_at, updated_at
	          FROM units
	          WHERE tenant_id = $1 AND deleted_at IS NULL`
	args := []interface{}{tenantID}
	argIndex := 2

	if unitType != "" {
		query += ` AND type = $` + strconv.Itoa(argIndex)
		args = append(args, unitType)
		argIndex++
	}

	if search != "" {
		query += ` AND (code ILIKE $` + strconv.Itoa(argIndex) + ` OR owner_name ILIKE $` + strconv.Itoa(argIndex) + `)`
		searchPattern := "%" + search + "%"
		args = append(args, searchPattern)
		argIndex++
	}

	query += ` ORDER BY code ASC LIMIT $` + strconv.Itoa(argIndex) + ` OFFSET $` + strconv.Itoa(argIndex+1)
	args = append(args, limit, offset)

	var units []models.Unit
	err := db.DB.Select(&units, query, args...)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch units: " + err.Error()})
	}

	// Get total count
	countQuery := `SELECT COUNT(*) FROM units WHERE tenant_id = $1 AND deleted_at IS NULL`
	countArgs := []interface{}{tenantID}
	if unitType != "" {
		countQuery += ` AND type = $2`
		countArgs = append(countArgs, unitType)
	}
	if search != "" {
		searchPattern := "%" + search + "%"
		if unitType != "" {
			countQuery += ` AND (code ILIKE $3 OR owner_name ILIKE $3)`
			countArgs = append(countArgs, searchPattern)
		} else {
			countQuery += ` AND (code ILIKE $2 OR owner_name ILIKE $2)`
			countArgs = append(countArgs, searchPattern)
		}
	}

	var total int
	err = db.DB.Get(&total, countQuery, countArgs...)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to count units"})
	}

	totalPages := (total + limit - 1) / limit

	return c.JSON(http.StatusOK, map[string]interface{}{
		"units": units,
		"pagination": map[string]interface{}{
			"page":       page,
			"limit":      limit,
			"total":      total,
			"total_pages": totalPages,
		},
	})
}

// GetUnit gets a unit by ID with assigned users
func GetUnit(c echo.Context) error {
	tenantID := c.Get(string(middleware.CtxTenantID)).(string)
	unitID := c.Param("unit_id")

	var unit models.Unit
	err := db.DB.Get(&unit, `
		SELECT id, tenant_id, code, type, owner_name, owner_phone, owner_email, address, status, created_at, updated_at
		FROM units
		WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL
	`, unitID, tenantID)
	if err == sql.ErrNoRows {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Unit not found"})
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}

	// Get users assigned to this unit
	var users []map[string]interface{}
	rows, err := db.DB.Query(`
		SELECT 
			u.id, u.email, u.full_name, u.phone, u.avatar_url, u.status,
			tu.role_id, tu.status as tenant_user_status
		FROM users u
		INNER JOIN tenant_users tu ON u.id = tu.user_id
		WHERE tu.unit_id = $1 AND tu.tenant_id = $2
		AND u.deleted_at IS NULL AND tu.deleted_at IS NULL
		ORDER BY u.full_name ASC
	`, unitID, tenantID)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var userID, email, fullName, status, tenantUserStatus string
			var phone, avatarURL, roleID sql.NullString

			err := rows.Scan(
				&userID, &email, &fullName, &phone, &avatarURL, &status,
				&roleID, &tenantUserStatus)
			if err != nil {
				continue
			}

			userData := map[string]interface{}{
				"id":        userID,
				"email":     email,
				"full_name": fullName,
				"status":    status,
			}

			if phone.Valid {
				userData["phone"] = phone.String
			}
			if avatarURL.Valid {
				userData["avatar_url"] = avatarURL.String
			}

			// Get role info
			if roleID.Valid {
				var roleName string
				err = db.DB.QueryRow(`
					SELECT name FROM roles 
					WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL
				`, roleID.String, tenantID).Scan(&roleName)
				if err == nil {
					userData["role"] = map[string]interface{}{
						"id":   roleID.String,
						"name": roleName,
					}
				}
			}

			users = append(users, userData)
		}
	}

	// Return unit with users
	return c.JSON(http.StatusOK, map[string]interface{}{
		"id":          unit.ID,
		"tenant_id":   unit.TenantID,
		"code":        unit.Code,
		"type":        unit.Type,
		"owner_name":  unit.OwnerName,
		"owner_phone": unit.OwnerPhone,
		"owner_email": unit.OwnerEmail,
		"address":     unit.Address,
		"status":      unit.Status,
		"created_at":  unit.CreatedAt,
		"updated_at":  unit.UpdatedAt,
		"users":       users,
	})
}

// UpdateUnit updates a unit
func UpdateUnit(c echo.Context) error {
	tenantID := c.Get(string(middleware.CtxTenantID)).(string)
	unitID := c.Param("unit_id")

	req := new(models.UpdateUnitRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// Check if unit exists
	var exists bool
	err := db.DB.Get(&exists, `
		SELECT EXISTS(
			SELECT 1 FROM units 
			WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL
		)
	`, unitID, tenantID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}
	if !exists {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Unit not found"})
	}

	// Check if code already exists (if updating code)
	if req.Code != nil {
		var codeExists bool
		err = db.DB.Get(&codeExists, `
			SELECT EXISTS(
				SELECT 1 FROM units 
				WHERE tenant_id = $1 AND code = $2 AND id != $3 AND deleted_at IS NULL
			)
		`, tenantID, *req.Code, unitID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
		}
		if codeExists {
			return c.JSON(http.StatusConflict, map[string]string{"error": "Unit code already exists"})
		}
	}

	// Build update query dynamically
	updates := []string{}
	args := []interface{}{}
	argIndex := 1

	if req.Code != nil {
		updates = append(updates, "code = $"+strconv.Itoa(argIndex))
		args = append(args, *req.Code)
		argIndex++
	}
	if req.Type != nil {
		updates = append(updates, "type = $"+strconv.Itoa(argIndex))
		args = append(args, *req.Type)
		argIndex++
	}
	if req.OwnerName != nil {
		updates = append(updates, "owner_name = $"+strconv.Itoa(argIndex))
		args = append(args, *req.OwnerName)
		argIndex++
	}
	if req.OwnerPhone != nil {
		updates = append(updates, "owner_phone = $"+strconv.Itoa(argIndex))
		args = append(args, *req.OwnerPhone)
		argIndex++
	}
	if req.OwnerEmail != nil {
		updates = append(updates, "owner_email = $"+strconv.Itoa(argIndex))
		args = append(args, *req.OwnerEmail)
		argIndex++
	}
	if req.Address != nil {
		updates = append(updates, "address = $"+strconv.Itoa(argIndex))
		args = append(args, *req.Address)
		argIndex++
	}
	if req.Status != nil {
		updates = append(updates, "status = $"+strconv.Itoa(argIndex))
		args = append(args, *req.Status)
		argIndex++
	}

	if len(updates) == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "No fields to update"})
	}

	updates = append(updates, "updated_at = NOW()")
	args = append(args, unitID, tenantID)

	// Join updates
	updateStr := ""
	for i, update := range updates {
		if i > 0 {
			updateStr += ", "
		}
		updateStr += update
	}

	query := `UPDATE units SET ` + updateStr + ` 
	          WHERE id = $` + strconv.Itoa(argIndex) + ` AND tenant_id = $` + strconv.Itoa(argIndex+1) + `
	          RETURNING id, tenant_id, code, type, owner_name, owner_phone, owner_email, address, status, created_at, updated_at`

	var unit models.Unit
	err = db.DB.QueryRow(query, args...).Scan(
		&unit.ID, &unit.TenantID, &unit.Code, &unit.Type, &unit.OwnerName, &unit.OwnerPhone,
		&unit.OwnerEmail, &unit.Address, &unit.Status, &unit.CreatedAt, &unit.UpdatedAt)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update unit: " + err.Error()})
	}

	return c.JSON(http.StatusOK, unit)
}

// DeleteUnit soft deletes a unit
func DeleteUnit(c echo.Context) error {
	tenantID := c.Get(string(middleware.CtxTenantID)).(string)
	unitID := c.Param("unit_id")

	// Check if unit exists
	var exists bool
	err := db.DB.Get(&exists, `
		SELECT EXISTS(
			SELECT 1 FROM units 
			WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL
		)
	`, unitID, tenantID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}
	if !exists {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Unit not found"})
	}

	// Soft delete
	_, err = db.DB.Exec(`
		UPDATE units 
		SET deleted_at = NOW(), updated_at = NOW()
		WHERE id = $1 AND tenant_id = $2
	`, unitID, tenantID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete unit"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Unit deleted successfully"})
}

// AssignUserToUnit assigns a user to a unit
func AssignUserToUnit(c echo.Context) error {
	tenantID := c.Get(string(middleware.CtxTenantID)).(string)
	unitID := c.Param("unit_id")

	var req struct {
		UserID string `json:"user_id" validate:"required"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// Verify unit exists and belongs to tenant
	var unitExists bool
	err := db.DB.Get(&unitExists, `
		SELECT EXISTS(
			SELECT 1 FROM units 
			WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL
		)
	`, unitID, tenantID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}
	if !unitExists {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Unit not found"})
	}

	// Verify user exists and is member of tenant
	var userExists bool
	err = db.DB.Get(&userExists, `
		SELECT EXISTS(
			SELECT 1 FROM tenant_users 
			WHERE user_id = $1 AND tenant_id = $2 AND deleted_at IS NULL AND status = 'active'
		)
	`, req.UserID, tenantID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}
	if !userExists {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found or not a member of this tenant"})
	}

	// Update tenant_users with unit_id
	_, err = db.DB.Exec(`
		UPDATE tenant_users 
		SET unit_id = $1, updated_at = NOW()
		WHERE user_id = $2 AND tenant_id = $3
	`, unitID, req.UserID, tenantID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to assign user to unit"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "User assigned to unit successfully"})
}

