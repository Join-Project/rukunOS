package handlers

import (
	"database/sql"
	"net/http"
	"rukunos-backend/db"
	"rukunos-backend/middleware"
	"rukunos-backend/models"

	"github.com/labstack/echo/v4"
)

// CreateTenant creates a new tenant (super admin only for now)
func CreateTenant(c echo.Context) error {
	req := new(models.CreateTenantRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// Start transaction
	tx, err := db.DB.Beginx()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to start transaction"})
	}
	defer tx.Rollback()

	// Create tenant
	var tenant models.Tenant
	query := `INSERT INTO tenants (name, code, address, phone, email, status)
	          VALUES ($1, $2, $3, $4, $5, 'active')
	          RETURNING id, name, code, address, phone, email, settings, modules, status, created_at, updated_at`
	
	err = tx.QueryRow(query, req.Name, req.Code, req.Address, req.Phone, req.Email).Scan(
		&tenant.ID, &tenant.Name, &tenant.Code, &tenant.Address, &tenant.Phone, 
		&tenant.Email, &tenant.Settings, &tenant.Modules, &tenant.Status, 
		&tenant.CreatedAt, &tenant.UpdatedAt)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create tenant: " + err.Error(),
		})
	}

	// Create default roles for tenant
	_, err = tx.Exec(`SELECT create_default_roles_for_tenant($1)`, tenant.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create default roles: " + err.Error(),
		})
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to commit transaction",
		})
	}

	return c.JSON(http.StatusCreated, tenant)
}

// GetTenant gets tenant information
// If tenant_id is provided in URL param, use it (for admin viewing other tenants)
// Otherwise, use tenant_id from context (for current user's tenant)
func GetTenant(c echo.Context) error {
	var tenantID string
	var ok bool
	
	// Try to get tenant_id from URL parameter first
	tenantIDParam := c.Param("tenant_id")
	if tenantIDParam != "" {
		// Verify that the requested tenant_id matches the context tenant_id
		// (for security, users can only view their own tenant)
		contextTenantIDInterface := c.Get("tenant_id")
		if contextTenantIDInterface != nil {
			contextTenantID, ok := contextTenantIDInterface.(string)
			if ok && contextTenantID == tenantIDParam {
				tenantID = tenantIDParam
			} else {
				// If param doesn't match context, use context tenant_id
				tenantID = contextTenantID
			}
		} else {
			tenantID = tenantIDParam
		}
	} else {
		// Get tenant_id from context
		tenantIDInterface := c.Get("tenant_id")
		if tenantIDInterface == nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Tenant ID not found in context"})
		}
		tenantID, ok = tenantIDInterface.(string)
		if !ok || tenantID == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid tenant ID in context"})
		}
	}

	var tenant models.Tenant
	err := db.DB.Get(&tenant, `
		SELECT id, name, code, address, phone, email, settings, modules, status, created_at, updated_at
		FROM tenants
		WHERE id = $1 AND deleted_at IS NULL
	`, tenantID)
	if err == sql.ErrNoRows {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Tenant not found"})
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}

	return c.JSON(http.StatusOK, tenant)
}

// GetCurrentUser gets current user with tenant info
func GetCurrentUser(c echo.Context) error {
	// Get user_id from context with proper nil checking
	// Try both "user_id" and CtxUserID key
	userIDInterface := c.Get("user_id")
	if userIDInterface == nil {
		userIDInterface = c.Get(string(middleware.CtxUserID))
	}
	if userIDInterface == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "User ID not found in context"})
	}
	userID, ok := userIDInterface.(string)
	if !ok || userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid user ID in context"})
	}

	// Get tenant_id from context with proper nil checking
	// Try both "tenant_id" and CtxTenantID key
	tenantIDInterface := c.Get("tenant_id")
	if tenantIDInterface == nil {
		tenantIDInterface = c.Get(string(middleware.CtxTenantID))
	}
	if tenantIDInterface == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Tenant ID not found in context"})
	}
	tenantID, ok := tenantIDInterface.(string)
	if !ok || tenantID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid tenant ID in context"})
	}

	var userIDFromDB, email, fullName, status string
	var phone, avatarURL, roleID, unitID sql.NullString
	var createdAt, updatedAt sql.NullTime
	
	err := db.DB.QueryRow(`
		SELECT u.id, u.email, u.full_name, u.phone, u.avatar_url, u.status, 
		       u.created_at, u.updated_at, tu.role_id, tu.unit_id
		FROM users u
		JOIN tenant_users tu ON u.id = tu.user_id
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

	// Get tenant info
	var tenantName, tenantCode string
	err = db.DB.QueryRow(`SELECT name, code FROM tenants WHERE id = $1`, tenantID).Scan(&tenantName, &tenantCode)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get tenant info: " + err.Error()})
	}

	// Get user permissions
	permissions, _ := GetUserPermissions(userID, tenantID)

	// Get role info
	var roleName string
	if roleID.Valid {
		err = db.DB.QueryRow(`SELECT name FROM roles WHERE id = $1`, roleID.String).Scan(&roleName)
		if err != nil {
			// Log error but don't fail request
			roleName = ""
		}
	}

	response := map[string]interface{}{
		"id":          userIDFromDB,
		"email":       email,
		"full_name":   fullName,
		"phone":       nil,
		"tenant_id":   tenantID,
		"tenant_name": tenantName,
		"tenant_code": tenantCode, // Add tenant_code to response
		"role_id":     nil,
		"role_name":   roleName,
		"unit_id":     nil,
		"permissions": permissions,
	}

	if phone.Valid {
		response["phone"] = phone.String
	}
	if roleID.Valid {
		response["role_id"] = roleID.String
	}
	if unitID.Valid {
		response["unit_id"] = unitID.String
	}

	return c.JSON(http.StatusOK, response)
}

// Helper function to get user permissions
func GetUserPermissions(userID, tenantID string) ([]string, error) {
	var permissions []string
	err := db.DB.Select(&permissions, `
		SELECT DISTINCT p.key
		FROM tenant_users tu
		JOIN roles r ON tu.role_id = r.id
		JOIN role_permissions rp ON r.id = rp.role_id
		JOIN permissions p ON rp.permission_id = p.id
		WHERE tu.user_id = $1
		AND tu.tenant_id = $2
		AND tu.deleted_at IS NULL
		AND r.deleted_at IS NULL
		AND tu.status = 'active'
	`, userID, tenantID)
	return permissions, err
}
