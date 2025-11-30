package handlers

import (
	"database/sql"
	"net/http"
	"os"
	"strings"
	"time"
	"rukunos-backend/db"
	"rukunos-backend/models"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func Register(c echo.Context) error {
	req := new(models.RegisterRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// Validate registration type
	if req.RegistrationType != "tenant" && req.RegistrationType != "warga" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "registration_type must be 'tenant' or 'warga'"})
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

	// Insert user
	query := `INSERT INTO users (email, password_hash, full_name, auth_provider, status) 
	          VALUES ($1, $2, $3, 'email', 'active') 
	          RETURNING id, created_at, updated_at`
	var user models.User
	user.Email = req.Email
	user.FullName = req.FullName
	user.AuthProvider = "email"
	user.Status = "active"

	err = tx.QueryRow(query, req.Email, string(hashedPassword), req.FullName).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create user: " + err.Error()})
	}

	var tenantID string
	var roleID string
	var tenantName string

	if req.RegistrationType == "tenant" {
		// Register as Tenant - create new tenant and assign Admin role
		if req.TenantName == nil || *req.TenantName == "" {
			tx.Rollback()
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "tenant_name is required for tenant registration"})
		}
		if req.TenantCode == nil || *req.TenantCode == "" {
			tx.Rollback()
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "tenant_code is required for tenant registration"})
		}

		// Check if tenant code already exists
		var existingTenantID string
		err = tx.Get(&existingTenantID, `
			SELECT id FROM tenants 
			WHERE code = $1 AND deleted_at IS NULL
		`, *req.TenantCode)
		if err == nil {
			tx.Rollback()
			return c.JSON(http.StatusConflict, map[string]string{"error": "Tenant code already exists"})
		} else if err != sql.ErrNoRows {
			tx.Rollback()
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to check tenant code"})
		}

		// Create new tenant
		tenantAddress := ""
		if req.TenantAddress != nil {
			tenantAddress = *req.TenantAddress
		}
		err = tx.QueryRow(`
			INSERT INTO tenants (name, code, address, status)
			VALUES ($1, $2, $3, 'active')
			RETURNING id, name
		`, *req.TenantName, strings.ToUpper(*req.TenantCode), tenantAddress).Scan(&tenantID, &tenantName)
		if err != nil {
			tx.Rollback()
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create tenant: " + err.Error()})
		}

		// Create default roles for the new tenant
		_, err = tx.Exec(`SELECT create_default_roles_for_tenant($1)`, tenantID)
		if err != nil {
			tx.Rollback()
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create default roles: " + err.Error()})
		}

		// Get Admin role for the new tenant
		err = tx.Get(&roleID, `
			SELECT id FROM roles 
			WHERE tenant_id = $1 
			AND name = 'Admin' 
			AND is_system = true 
			AND deleted_at IS NULL
		`, tenantID)
		if err != nil {
			tx.Rollback()
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to find Admin role"})
		}

	} else {
		// Register as Warga - join existing tenant with code
		var tenantCode string
		if req.TenantCodeJoin != nil && *req.TenantCodeJoin != "" {
			tenantCode = strings.ToUpper(*req.TenantCodeJoin)
		} else if req.TenantCode != nil && *req.TenantCode != "" {
			tenantCode = strings.ToUpper(*req.TenantCode)
		} else {
			tx.Rollback()
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "tenant_code is required for warga registration"})
		}

		// Find tenant by code
		err = tx.QueryRow(`
			SELECT id, name FROM tenants 
			WHERE code = $1 
			AND deleted_at IS NULL 
			AND status = 'active'
		`, tenantCode).Scan(&tenantID, &tenantName)
		if err == sql.ErrNoRows {
			tx.Rollback()
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Tenant not found. Please check the tenant code."})
		} else if err != nil {
			tx.Rollback()
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to find tenant"})
		}

		// Get Warga role for the tenant
		err = tx.Get(&roleID, `
			SELECT id FROM roles 
			WHERE tenant_id = $1 
			AND name = 'Warga' 
			AND is_system = true 
			AND deleted_at IS NULL
		`, tenantID)
		if err != nil {
			tx.Rollback()
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to find Warga role"})
		}
	}

	// Add user to tenant with assigned role
	_, err = tx.Exec(`
		INSERT INTO tenant_users (tenant_id, user_id, role_id, status)
		VALUES ($1, $2, $3, 'active')
	`, tenantID, user.ID, roleID)
	if err != nil {
		tx.Rollback()
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to add user to tenant: " + err.Error()})
	}

	user.TenantID = &tenantID
	user.RoleID = &roleID

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to commit transaction"})
	}

	// Generate JWT token
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = user.ID
	claims["tenant_id"] = tenantID
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix() // 3 days

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "secret" // Default for development only
	}

	t, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate token"})
	}

	// Get user permissions
	permissions, _ := GetUserPermissions(user.ID, tenantID)

	// Get role name
	var roleName string
	if roleID != "" {
		err = db.DB.Get(&roleName, `SELECT name FROM roles WHERE id = $1 AND deleted_at IS NULL`, roleID)
		if err != nil {
			roleName = ""
		}
	}

	// Return response with token
	response := models.AuthResponse{
		Token:    t,
		User:     user,
		TenantID: tenantID,
	}
	response.User.TenantID = &tenantID
	response.User.RoleID = &roleID

	// Add tenant name to user object for frontend
	userWithTenant := map[string]interface{}{
		"id":          user.ID,
		"email":       user.Email,
		"full_name":   user.FullName,
		"tenant_id":   tenantID,
		"tenant_name": tenantName,
		"role_id":     roleID,
		"role_name":   roleName,
		"permissions": permissions,
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"token": t,
		"user":  userWithTenant,
		"tenant_id": tenantID,
	})
}

func Login(c echo.Context) error {
	req := new(models.LoginRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// Find user - only look for email/password users
	var user models.User
	var googleID sql.NullString
	var authProvider sql.NullString
	var status sql.NullString
	query := `SELECT id, email, password_hash, full_name, google_id, auth_provider, status
	          FROM users 
	          WHERE email = $1 
	          AND deleted_at IS NULL
	          AND (auth_provider = 'email' OR auth_provider = 'both')`
	err := db.DB.QueryRow(query, req.Email).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.FullName, 
		&googleID, &authProvider, &status)
	if err == sql.ErrNoRows {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "User belum terdaftar"})
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}

	// Set optional fields
	if googleID.Valid {
		user.GoogleID = &googleID.String
	}
	if authProvider.Valid {
		user.AuthProvider = authProvider.String
	} else {
		user.AuthProvider = "email"
	}
	if status.Valid {
		user.Status = status.String
	} else {
		user.Status = "active"
	}

	// Check if user is active
	if user.Status != "active" {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "User account is inactive"})
	}

	// Check password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid credentials"})
	}

	// Get user's tenants
	var tenantIDs []string
	err = db.DB.Select(&tenantIDs, `
		SELECT tenant_id 
		FROM tenant_users 
		WHERE user_id = $1 
		AND deleted_at IS NULL 
		AND status = 'active'
	`, user.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get user tenants"})
	}

	if len(tenantIDs) == 0 {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "User is not a member of any tenant"})
	}

	// Determine which tenant to use
	var selectedTenantID string
	if req.TenantID != nil {
		// Check if user is member of requested tenant
		for _, tid := range tenantIDs {
			if tid == *req.TenantID {
				selectedTenantID = tid
				break
			}
		}
		if selectedTenantID == "" {
			return c.JSON(http.StatusForbidden, map[string]string{"error": "User is not a member of the specified tenant"})
		}
	} else {
		// Use first tenant (or primary tenant in future)
		selectedTenantID = tenantIDs[0]
	}

	// Get user's role in selected tenant
	var roleID sql.NullString
	err = db.DB.Get(&roleID, `
		SELECT role_id 
		FROM tenant_users 
		WHERE user_id = $1 
		AND tenant_id = $2 
		AND deleted_at IS NULL
	`, user.ID, selectedTenantID)
	if err != nil && err != sql.ErrNoRows {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get user role"})
	}
	if roleID.Valid {
		user.RoleID = &roleID.String
	}

	// Generate JWT with tenant_id
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = user.ID
	claims["tenant_id"] = selectedTenantID
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix() // 3 days

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "secret" // Default for development only
	}

	t, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate token"})
	}

	// Get tenant name
	var tenantName string
	err = db.DB.Get(&tenantName, `SELECT name FROM tenants WHERE id = $1`, selectedTenantID)
	if err != nil {
		tenantName = ""
	}

	// Get user permissions
	permissions, _ := GetUserPermissions(user.ID, selectedTenantID)

	// Get role name
	var roleName string
	if user.RoleID != nil {
		err = db.DB.Get(&roleName, `SELECT name FROM roles WHERE id = $1 AND deleted_at IS NULL`, *user.RoleID)
		if err != nil {
			roleName = ""
		}
	}

	// Return response with tenant info
	userWithTenant := map[string]interface{}{
		"id":          user.ID,
		"email":       user.Email,
		"full_name":   user.FullName,
		"tenant_id":   selectedTenantID,
		"tenant_name": tenantName,
		"role_id":     user.RoleID,
		"role_name":   roleName,
		"permissions": permissions,
	}

	response := map[string]interface{}{
		"token":     t,
		"user":      userWithTenant,
		"tenant_id": selectedTenantID,
	}
	if len(tenantIDs) > 1 {
		response["tenants"] = tenantIDs
	}

	return c.JSON(http.StatusOK, response)
}
