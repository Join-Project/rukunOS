package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"
	"rukunos-backend/db"
	"rukunos-backend/middleware"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

// ListUsers lists all users in a tenant with pagination and filters
func ListUsers(c echo.Context) error {
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

	roleID := c.QueryParam("role_id")
	unitID := c.QueryParam("unit_id")
	search := c.QueryParam("search")

	// Build query
	query := `
		SELECT 
			u.id, u.email, u.full_name, u.phone, u.avatar_url, u.status, 
			u.created_at, u.updated_at,
			tu.role_id, tu.unit_id, tu.status as tenant_user_status
		FROM users u
		INNER JOIN tenant_users tu ON u.id = tu.user_id
		WHERE tu.tenant_id = $1 
		AND u.deleted_at IS NULL 
		AND tu.deleted_at IS NULL`
	
	args := []interface{}{tenantID}
	argIndex := 2

	if roleID != "" {
		query += ` AND tu.role_id = $` + strconv.Itoa(argIndex)
		args = append(args, roleID)
		argIndex++
	}

	if unitID != "" {
		query += ` AND tu.unit_id = $` + strconv.Itoa(argIndex)
		args = append(args, unitID)
		argIndex++
	}

	if search != "" {
		query += ` AND (u.email ILIKE $` + strconv.Itoa(argIndex) + ` OR u.full_name ILIKE $` + strconv.Itoa(argIndex) + `)`
		searchPattern := "%" + search + "%"
		args = append(args, searchPattern)
		argIndex++
	}

	query += ` ORDER BY u.created_at DESC LIMIT $` + strconv.Itoa(argIndex) + ` OFFSET $` + strconv.Itoa(argIndex+1)
	args = append(args, limit, offset)

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch users: " + err.Error()})
	}
	defer rows.Close()

	var users []map[string]interface{}
	for rows.Next() {
		var userID, email, fullName, status, tenantUserStatus string
		var phone, avatarURL, roleIDFromDB, unitID sql.NullString
		var createdAt, updatedAt sql.NullTime

		err := rows.Scan(
			&userID, &email, &fullName, &phone, &avatarURL, &status,
			&createdAt, &updatedAt, &roleIDFromDB, &unitID, &tenantUserStatus)
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
		if createdAt.Valid {
			userData["created_at"] = createdAt.Time
		}
		if updatedAt.Valid {
			userData["updated_at"] = updatedAt.Time
		}

		// Get role info
		if roleIDFromDB.Valid {
			var roleName string
			var roleDescription sql.NullString
			err = db.DB.QueryRow(`
				SELECT name, description FROM roles 
				WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL
			`, roleIDFromDB.String, tenantID).Scan(&roleName, &roleDescription)
			if err == nil {
				roleData := map[string]interface{}{
					"id":   roleIDFromDB.String,
					"name": roleName,
				}
				if roleDescription.Valid {
					roleData["description"] = roleDescription.String
				}
				userData["role"] = roleData
			}
		}

		// Get unit info
		if unitID.Valid {
			var unitCode, unitType string
			err = db.DB.QueryRow(`
				SELECT code, type FROM units 
				WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL
			`, unitID.String, tenantID).Scan(&unitCode, &unitType)
			if err == nil {
				userData["unit"] = map[string]interface{}{
					"id":   unitID.String,
					"code": unitCode,
					"type": unitType,
				}
			}
		}

		users = append(users, userData)
	}

	// Get total count
	countQuery := `
		SELECT COUNT(*) 
		FROM users u
		INNER JOIN tenant_users tu ON u.id = tu.user_id
		WHERE tu.tenant_id = $1 
		AND u.deleted_at IS NULL 
		AND tu.deleted_at IS NULL`
	countArgs := []interface{}{tenantID}
	
	countArgIndex := 2
	if roleID != "" {
		countQuery += ` AND tu.role_id = $` + strconv.Itoa(countArgIndex)
		countArgs = append(countArgs, roleID)
		countArgIndex++
	}
	if unitID != "" {
		countQuery += ` AND tu.unit_id = $` + strconv.Itoa(countArgIndex)
		countArgs = append(countArgs, unitID)
		countArgIndex++
	}
	if search != "" {
		searchPattern := "%" + search + "%"
		countQuery += ` AND (u.email ILIKE $` + strconv.Itoa(countArgIndex) + ` OR u.full_name ILIKE $` + strconv.Itoa(countArgIndex) + `)`
		countArgs = append(countArgs, searchPattern)
	}

	var total int
	err = db.DB.Get(&total, countQuery, countArgs...)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to count users"})
	}

	totalPages := (total + limit - 1) / limit

	return c.JSON(http.StatusOK, map[string]interface{}{
		"users": users,
		"pagination": map[string]interface{}{
			"page":       page,
			"limit":      limit,
			"total":      total,
			"total_pages": totalPages,
		},
	})
}

// GetUser gets a user by ID with role and unit info
func GetUser(c echo.Context) error {
	tenantID := c.Get(string(middleware.CtxTenantID)).(string)
	userID := c.Param("user_id")

	// Get user with tenant info
	var userIDFromDB, email, fullName, status string
	var phone, avatarURL, roleID, unitID sql.NullString
	var createdAt, updatedAt sql.NullTime
	
	err := db.DB.QueryRow(`
		SELECT 
			u.id, u.email, u.full_name, u.phone, u.avatar_url, u.status, 
			u.created_at, u.updated_at,
			tu.role_id, tu.unit_id
		FROM users u
		INNER JOIN tenant_users tu ON u.id = tu.user_id
		WHERE u.id = $1 AND tu.tenant_id = $2 
		AND u.deleted_at IS NULL AND tu.deleted_at IS NULL
	`, userID, tenantID).Scan(
		&userIDFromDB, &email, &fullName, &phone, &avatarURL, &status,
		&createdAt, &updatedAt, &roleID, &unitID)
	
	if err == sql.ErrNoRows {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error: " + err.Error()})
	}

	userData := map[string]interface{}{
		"id":        userIDFromDB,
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
	if createdAt.Valid {
		userData["created_at"] = createdAt.Time
	}
	if updatedAt.Valid {
		userData["updated_at"] = updatedAt.Time
	}

	// Get role info
	if roleID.Valid {
		var roleName string
		var roleDescription sql.NullString
		err = db.DB.QueryRow(`
			SELECT name, description FROM roles 
			WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL
		`, roleID.String, tenantID).Scan(&roleName, &roleDescription)
		if err == nil {
			roleData := map[string]interface{}{
				"id":   roleID.String,
				"name": roleName,
			}
			if roleDescription.Valid {
				roleData["description"] = roleDescription.String
			}
			userData["role"] = roleData
		}
	}

	// Get unit info
	if unitID.Valid {
		var unitCode, unitType string
		var ownerName sql.NullString
		err = db.DB.QueryRow(`
			SELECT code, type, owner_name FROM units 
			WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL
		`, unitID.String, tenantID).Scan(&unitCode, &unitType, &ownerName)
		if err == nil {
			unitData := map[string]interface{}{
				"id":   unitID.String,
				"code": unitCode,
				"type": unitType,
			}
			if ownerName.Valid {
				unitData["owner_name"] = ownerName.String
			}
			userData["unit"] = unitData
		}
	}

	// Get user permissions
	if roleID.Valid {
		permissions, _ := GetUserPermissions(userID, tenantID)
		userData["permissions"] = permissions
	}

	return c.JSON(http.StatusOK, userData)
}

// UpdateUser updates user role or status in tenant
func UpdateUser(c echo.Context) error {
	tenantID := c.Get(string(middleware.CtxTenantID)).(string)
	userID := c.Param("user_id")

	var req struct {
		RoleID *string `json:"role_id,omitempty"`
		UnitID *string `json:"unit_id,omitempty"` // Support unit assignment
		Status *string `json:"status,omitempty" validate:"omitempty,oneof=active inactive"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// Check if user exists and is member of tenant
	var exists bool
	err := db.DB.Get(&exists, `
		SELECT EXISTS(
			SELECT 1 FROM tenant_users 
			WHERE user_id = $1 AND tenant_id = $2 AND deleted_at IS NULL
		)
	`, userID, tenantID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}
	if !exists {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found or not a member of this tenant"})
	}

	// Update role_id if provided
	if req.RoleID != nil {
		// Verify role exists and belongs to tenant
		var roleExists bool
		err = db.DB.Get(&roleExists, `
			SELECT EXISTS(
				SELECT 1 FROM roles 
				WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL
			)
		`, *req.RoleID, tenantID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
		}
		if !roleExists {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Role not found"})
		}

		_, err = db.DB.Exec(`
			UPDATE tenant_users 
			SET role_id = $1, updated_at = NOW()
			WHERE user_id = $2 AND tenant_id = $3
		`, *req.RoleID, userID, tenantID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update user role"})
		}
	}

	// Update unit_id if provided
	if req.UnitID != nil {
		// If unit_id is empty string, set to NULL (unassign from unit)
		if *req.UnitID == "" {
			_, err = db.DB.Exec(`
				UPDATE tenant_users 
				SET unit_id = NULL, updated_at = NOW()
				WHERE user_id = $1 AND tenant_id = $2
			`, userID, tenantID)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to unassign user from unit"})
			}
		} else {
			// Verify unit exists and belongs to tenant
			var unitExists bool
			err = db.DB.Get(&unitExists, `
				SELECT EXISTS(
					SELECT 1 FROM units 
					WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL
				)
			`, *req.UnitID, tenantID)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
			}
			if !unitExists {
				return c.JSON(http.StatusNotFound, map[string]string{"error": "Unit not found"})
			}

			_, err = db.DB.Exec(`
				UPDATE tenant_users 
				SET unit_id = $1, updated_at = NOW()
				WHERE user_id = $2 AND tenant_id = $3
			`, *req.UnitID, userID, tenantID)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update user unit"})
			}
		}
	}

	// Update status if provided
	if req.Status != nil {
		_, err = db.DB.Exec(`
			UPDATE tenant_users 
			SET status = $1, updated_at = NOW()
			WHERE user_id = $2 AND tenant_id = $3
		`, *req.Status, userID, tenantID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update user status"})
		}
	}

	// Return updated user
	return GetUser(c)
}

// CreateUserByAdmin creates a new user and adds to current tenant (admin only)
// This is different from Register endpoint - it uses tenant_id from context
func CreateUserByAdmin(c echo.Context) error {
	tenantID := c.Get(string(middleware.CtxTenantID)).(string)
	
	var req struct {
		Email    string  `json:"email" validate:"required,email"`
		Password string  `json:"password" validate:"required,min=6"`
		FullName string  `json:"full_name" validate:"required"`
		RoleID   *string `json:"role_id,omitempty"`
		UnitID   *string `json:"unit_id,omitempty"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to hash password"})
	}

	// Start transaction
	tx, err := db.DB.Beginx()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to start transaction"})
	}
	defer tx.Rollback()

	// Check if email already exists
	var existingUserID string
	err = tx.Get(&existingUserID, `
		SELECT id FROM users 
		WHERE email = $1 AND deleted_at IS NULL
	`, req.Email)
	if err == nil {
		tx.Rollback()
		return c.JSON(http.StatusConflict, map[string]string{"error": "Email already exists"})
	} else if err != sql.ErrNoRows {
		tx.Rollback()
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to check email"})
	}

	// Insert user
	userID := uuid.New().String()
	query := `INSERT INTO users (id, email, password_hash, full_name, auth_provider, status) 
	          VALUES ($1, $2, $3, $4, 'email', 'active') 
	          RETURNING id, created_at, updated_at`
	
	var userIDFromInsert string
	var userCreatedAt, userUpdatedAt time.Time
	err = tx.QueryRow(query, userID, req.Email, string(hashedPassword), req.FullName).Scan(&userIDFromInsert, &userCreatedAt, &userUpdatedAt)
	if err != nil {
		tx.Rollback()
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create user: " + err.Error()})
	}

	// Determine role_id - use provided role_id or default to Warga role
	var roleID string
	if req.RoleID != nil && *req.RoleID != "" {
		// Verify role exists and belongs to tenant
		var roleExists bool
		err = tx.Get(&roleExists, `
			SELECT EXISTS(
				SELECT 1 FROM roles 
				WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL
			)
		`, *req.RoleID, tenantID)
		if err != nil {
			tx.Rollback()
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
		}
		if !roleExists {
			tx.Rollback()
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Role not found"})
		}
		roleID = *req.RoleID
	} else {
		// Get default Warga role
		err = tx.Get(&roleID, `
			SELECT id FROM roles 
			WHERE tenant_id = $1 
			AND (name = 'Warga' OR name = 'Member')
			AND is_system = true 
			AND deleted_at IS NULL
			LIMIT 1
		`, tenantID)
		if err != nil {
			tx.Rollback()
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to find default Warga role"})
		}
	}

	// Add user to tenant with assigned role and unit
	_, err = tx.Exec(`
		INSERT INTO tenant_users (tenant_id, user_id, role_id, unit_id, status)
		VALUES ($1, $2, $3, $4, 'active')
	`, tenantID, userID, roleID, req.UnitID)
	if err != nil {
		tx.Rollback()
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to add user to tenant: " + err.Error()})
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to commit transaction"})
	}

	// Return created user info - get full details
	var userIDFromDB, email, fullName, status string
	var phone, avatarURL, roleIDFromDB, unitID sql.NullString
	var userCreatedAtDB, userUpdatedAtDB time.Time
	
	err = db.DB.QueryRow(`
		SELECT 
			u.id, u.email, u.full_name, u.phone, u.avatar_url, u.status, 
			u.created_at, u.updated_at,
			tu.role_id, tu.unit_id
		FROM users u
		INNER JOIN tenant_users tu ON u.id = tu.user_id
		WHERE u.id = $1 AND tu.tenant_id = $2 
		AND u.deleted_at IS NULL AND tu.deleted_at IS NULL
	`, userIDFromInsert, tenantID).Scan(
		&userIDFromDB, &email, &fullName, &phone, &avatarURL, &status,
		&userCreatedAtDB, &userUpdatedAtDB, &roleIDFromDB, &unitID)
	
	if err != nil {
		// If we can't get full details, return basic info
		return c.JSON(http.StatusCreated, map[string]interface{}{
			"id":        userIDFromInsert,
			"email":     req.Email,
			"full_name": req.FullName,
			"message":   "User created successfully",
		})
	}

	userData := map[string]interface{}{
		"id":         userIDFromDB,
		"email":      email,
		"full_name":  fullName,
		"status":     status,
		"created_at": userCreatedAtDB,
		"updated_at": userUpdatedAtDB,
	}

	if phone.Valid {
		userData["phone"] = phone.String
	}
	if avatarURL.Valid {
		userData["avatar_url"] = avatarURL.String
	}

	// Get role info
	if roleIDFromDB.Valid {
		var roleName string
		var roleDescription sql.NullString
		err = db.DB.QueryRow(`
			SELECT name, description FROM roles 
			WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL
		`, roleIDFromDB.String, tenantID).Scan(&roleName, &roleDescription)
		if err == nil {
			roleData := map[string]interface{}{
				"id":   roleIDFromDB.String,
				"name": roleName,
			}
			if roleDescription.Valid {
				roleData["description"] = roleDescription.String
			}
			userData["role"] = roleData
		}
	}

	// Get unit info
	if unitID.Valid {
		var unitCode, unitType string
		err = db.DB.QueryRow(`
			SELECT code, type FROM units 
			WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL
		`, unitID.String, tenantID).Scan(&unitCode, &unitType)
		if err == nil {
			userData["unit"] = map[string]interface{}{
				"id":   unitID.String,
				"code": unitCode,
				"type": unitType,
			}
		}
	}

	return c.JSON(http.StatusCreated, userData)
}

// RemoveUserFromTenant removes a user from tenant (soft delete tenant_users)
func RemoveUserFromTenant(c echo.Context) error {
	tenantID := c.Get(string(middleware.CtxTenantID)).(string)
	userID := c.Param("user_id")

	// Check if user exists and is member of tenant
	var exists bool
	err := db.DB.Get(&exists, `
		SELECT EXISTS(
			SELECT 1 FROM tenant_users 
			WHERE user_id = $1 AND tenant_id = $2 AND deleted_at IS NULL
		)
	`, userID, tenantID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}
	if !exists {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found or not a member of this tenant"})
	}

	// Soft delete tenant_users entry
	_, err = db.DB.Exec(`
		UPDATE tenant_users 
		SET deleted_at = NOW(), updated_at = NOW()
		WHERE user_id = $1 AND tenant_id = $2
	`, userID, tenantID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to remove user from tenant"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "User removed from tenant successfully"})
}





