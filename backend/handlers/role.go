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

// CreateRole creates a new custom role
func CreateRole(c echo.Context) error {
	tenantID := c.Get(string(middleware.CtxTenantID)).(string)

	req := new(models.CreateRoleRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// Validate request
	if req.Name == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Role name is required"})
	}

	// Check if role name already exists in tenant
	var exists bool
	err := db.DB.Get(&exists, `
		SELECT EXISTS(
			SELECT 1 FROM roles 
			WHERE tenant_id = $1 AND name = $2 AND deleted_at IS NULL
		)
	`, tenantID, req.Name)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}
	if exists {
		return c.JSON(http.StatusConflict, map[string]string{"error": "Role name already exists"})
	}

	// Start transaction
	tx, err := db.DB.Beginx()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to start transaction"})
	}
	defer tx.Rollback()

	// Create role
	roleID := uuid.New().String()
	var description sql.NullString
	if req.Description != nil && *req.Description != "" {
		description = sql.NullString{String: *req.Description, Valid: true}
	}
	
	query := `INSERT INTO roles (id, tenant_id, name, description, is_system)
	          VALUES ($1, $2, $3, $4, false)
	          RETURNING id, tenant_id, name, description, is_system, created_at, updated_at`
	
	var role models.Role
	err = tx.QueryRow(query, roleID, tenantID, req.Name, description).Scan(
		&role.ID, &role.TenantID, &role.Name, &role.Description, &role.IsSystem, 
		&role.CreatedAt, &role.UpdatedAt)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create role: " + err.Error()})
	}

	// Assign permissions if provided
	if len(req.Permissions) > 0 {
		// Get permission IDs from keys
		permissionMap := make(map[string]string)
		permissionKeys := make([]interface{}, len(req.Permissions))
		for i, key := range req.Permissions {
			permissionKeys[i] = key
		}

		// Build query to get permission IDs
		placeholders := ""
		for i := range req.Permissions {
			if i > 0 {
				placeholders += ","
			}
			placeholders += "$" + strconv.Itoa(i+1)
		}

		var permissions []struct {
			ID  string `db:"id"`
			Key string `db:"key"`
		}
		err = tx.Select(&permissions, `SELECT id, key FROM permissions WHERE key IN (`+placeholders+`)`, permissionKeys...)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch permissions"})
		}

		for _, p := range permissions {
			permissionMap[p.Key] = p.ID
		}

		// Insert role_permissions
		for _, permKey := range req.Permissions {
			permID, ok := permissionMap[permKey]
			if !ok {
				continue // Skip invalid permissions
			}

			_, err = tx.Exec(`
				INSERT INTO role_permissions (id, role_id, permission_id)
				VALUES ($1, $2, $3)
			`, uuid.New().String(), roleID, permID)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to assign permissions"})
			}
		}
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to commit transaction"})
	}

	// Get role with permissions for response
	roleWithPerms, err := getRoleWithPermissions(roleID, tenantID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch role details"})
	}

	return c.JSON(http.StatusCreated, roleWithPerms)
}

// ListRoles lists all roles for a tenant
func ListRoles(c echo.Context) error {
	tenantID := c.Get(string(middleware.CtxTenantID)).(string)

	var roles []models.Role
	err := db.DB.Select(&roles, `
		SELECT id, tenant_id, name, description, is_system, created_at, updated_at
		FROM roles
		WHERE tenant_id = $1 AND deleted_at IS NULL
		ORDER BY is_system DESC, name ASC
	`, tenantID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch roles"})
	}

	// Get permissions for each role
	rolesWithPerms := make([]map[string]interface{}, len(roles))
	for i, role := range roles {
		permissions, _ := getRolePermissionKeys(role.ID)
		roleData := map[string]interface{}{
			"id":          role.ID,
			"name":        role.Name,
			"is_system":   role.IsSystem,
			"permissions": permissions,
			"created_at":  role.CreatedAt,
			"updated_at":  role.UpdatedAt,
		}
		if role.Description.Valid {
			roleData["description"] = role.Description.String
		}
		rolesWithPerms[i] = roleData
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"roles": rolesWithPerms,
	})
}

// GetRole gets a role by ID with permissions
func GetRole(c echo.Context) error {
	tenantID := c.Get(string(middleware.CtxTenantID)).(string)
	roleID := c.Param("role_id")

	roleWithPerms, err := getRoleWithPermissions(roleID, tenantID)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Role not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}

	return c.JSON(http.StatusOK, roleWithPerms)
}

// UpdateRole updates a role and its permissions
func UpdateRole(c echo.Context) error {
	tenantID := c.Get(string(middleware.CtxTenantID)).(string)
	roleID := c.Param("role_id")

	req := new(models.UpdateRoleRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// Check if role exists and belongs to tenant
	var role models.Role
	err := db.DB.Get(&role, `
		SELECT id, tenant_id, name, description, is_system
		FROM roles
		WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL
	`, roleID, tenantID)
	if err == sql.ErrNoRows {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Role not found"})
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}

	// System roles cannot be modified
	if role.IsSystem {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "System roles cannot be modified"})
	}

	// Start transaction
	tx, err := db.DB.Beginx()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to start transaction"})
	}
	defer tx.Rollback()

	// Update role name/description if provided
	if req.Name != nil || req.Description != nil {
		updates := []string{}
		args := []interface{}{}
		argIndex := 1

		if req.Name != nil {
			// Check if new name already exists
			var nameExists bool
			err = tx.Get(&nameExists, `
				SELECT EXISTS(
					SELECT 1 FROM roles 
					WHERE tenant_id = $1 AND name = $2 AND id != $3 AND deleted_at IS NULL
				)
			`, tenantID, *req.Name, roleID)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
			}
			if nameExists {
				return c.JSON(http.StatusConflict, map[string]string{"error": "Role name already exists"})
			}

			updates = append(updates, "name = $"+strconv.Itoa(argIndex))
			args = append(args, *req.Name)
			argIndex++
		}

		if req.Description != nil {
			updates = append(updates, "description = $"+strconv.Itoa(argIndex))
			args = append(args, *req.Description)
			argIndex++
		}

		if len(updates) > 0 {
			updates = append(updates, "updated_at = NOW()")
			args = append(args, roleID, tenantID)

			updateStr := ""
			for i, update := range updates {
				if i > 0 {
					updateStr += ", "
				}
				updateStr += update
			}

			_, err = tx.Exec(`UPDATE roles SET `+updateStr+` WHERE id = $`+strconv.Itoa(argIndex)+` AND tenant_id = $`+strconv.Itoa(argIndex+1), args...)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update role"})
			}
		}
	}

	// Update permissions if provided
	if req.Permissions != nil {
		// Delete existing permissions
		_, err = tx.Exec(`DELETE FROM role_permissions WHERE role_id = $1`, roleID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to remove existing permissions"})
		}

		// Add new permissions
		if len(*req.Permissions) > 0 {
			permissionKeys := make([]interface{}, len(*req.Permissions))
			for i, key := range *req.Permissions {
				permissionKeys[i] = key
			}

			placeholders := ""
			for i := range *req.Permissions {
				if i > 0 {
					placeholders += ","
				}
				placeholders += "$" + strconv.Itoa(i+1)
			}

			var permissions []struct {
				ID  string `db:"id"`
				Key string `db:"key"`
			}
			err = tx.Select(&permissions, `SELECT id, key FROM permissions WHERE key IN (`+placeholders+`)`, permissionKeys...)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch permissions"})
			}

			permissionMap := make(map[string]string)
			for _, p := range permissions {
				permissionMap[p.Key] = p.ID
			}

			for _, permKey := range *req.Permissions {
				permID, ok := permissionMap[permKey]
				if !ok {
					continue
				}

				_, err = tx.Exec(`
					INSERT INTO role_permissions (id, role_id, permission_id)
					VALUES ($1, $2, $3)
				`, uuid.New().String(), roleID, permID)
				if err != nil {
					return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to assign permissions"})
				}
			}
		}
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to commit transaction"})
	}

	// Get updated role with permissions
	roleWithPerms, err := getRoleWithPermissions(roleID, tenantID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch updated role"})
	}

	return c.JSON(http.StatusOK, roleWithPerms)
}

// DeleteRole soft deletes a role
func DeleteRole(c echo.Context) error {
	tenantID := c.Get(string(middleware.CtxTenantID)).(string)
	roleID := c.Param("role_id")

	// Check if role exists and is not system role
	var role models.Role
	err := db.DB.Get(&role, `
		SELECT id, tenant_id, is_system
		FROM roles
		WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL
	`, roleID, tenantID)
	if err == sql.ErrNoRows {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Role not found"})
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}

	// System roles cannot be deleted
	if role.IsSystem {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "System roles cannot be deleted"})
	}

	// Check if role is assigned to any users
	var assignedCount int
	err = db.DB.Get(&assignedCount, `
		SELECT COUNT(*) FROM tenant_users 
		WHERE role_id = $1 AND tenant_id = $2 AND deleted_at IS NULL
	`, roleID, tenantID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}
	if assignedCount > 0 {
		return c.JSON(http.StatusConflict, map[string]string{"error": "Cannot delete role that is assigned to users"})
	}

	// Soft delete
	_, err = db.DB.Exec(`
		UPDATE roles 
		SET deleted_at = NOW(), updated_at = NOW()
		WHERE id = $1 AND tenant_id = $2
	`, roleID, tenantID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete role"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Role deleted successfully"})
}

// ListPermissions lists all available permissions
func ListPermissions(c echo.Context) error {
	var permissions []models.Permission
	err := db.DB.Select(&permissions, `
		SELECT id, key, name, description, module, created_at
		FROM permissions
		ORDER BY module, key
	`)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch permissions"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"permissions": permissions,
	})
}

// AssignRoleToUser assigns a role to a user
func AssignRoleToUser(c echo.Context) error {
	tenantID := c.Get(string(middleware.CtxTenantID)).(string)
	userID := c.Param("user_id")

	var req struct {
		RoleID string `json:"role_id" validate:"required"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// Verify user exists and is member of tenant
	var userExists bool
	err := db.DB.Get(&userExists, `
		SELECT EXISTS(
			SELECT 1 FROM tenant_users 
			WHERE user_id = $1 AND tenant_id = $2 AND deleted_at IS NULL
		)
	`, userID, tenantID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}
	if !userExists {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found or not a member of this tenant"})
	}

	// Verify role exists and belongs to tenant
	var roleExists bool
	err = db.DB.Get(&roleExists, `
		SELECT EXISTS(
			SELECT 1 FROM roles 
			WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL
		)
	`, req.RoleID, tenantID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}
	if !roleExists {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Role not found"})
	}

	// Update tenant_users with new role_id
	_, err = db.DB.Exec(`
		UPDATE tenant_users 
		SET role_id = $1, updated_at = NOW()
		WHERE user_id = $2 AND tenant_id = $3
	`, req.RoleID, userID, tenantID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to assign role to user"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Role assigned to user successfully"})
}

// Helper functions

func getRoleWithPermissions(roleID, tenantID string) (map[string]interface{}, error) {
	var role models.Role
	err := db.DB.Get(&role, `
		SELECT id, tenant_id, name, description, is_system, created_at, updated_at
		FROM roles
		WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL
	`, roleID, tenantID)
	if err != nil {
		return nil, err
	}

	permissions, _ := getRolePermissionKeys(roleID)

	return map[string]interface{}{
		"id":          role.ID,
		"name":        role.Name,
		"description": role.Description,
		"is_system":   role.IsSystem,
		"permissions": permissions,
		"created_at":  role.CreatedAt,
		"updated_at":  role.UpdatedAt,
	}, nil
}

func getRolePermissionKeys(roleID string) ([]string, error) {
	var permissions []string
	err := db.DB.Select(&permissions, `
		SELECT p.key
		FROM role_permissions rp
		JOIN permissions p ON rp.permission_id = p.id
		WHERE rp.role_id = $1
		ORDER BY p.key
	`, roleID)
	return permissions, err
}

