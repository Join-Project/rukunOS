package middleware

import (
	"net/http"
	"rukunos-backend/db"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// TenantContextKey is used to store tenant_id in echo.Context
type TenantContextKey string
const CtxTenantID TenantContextKey = "tenantID"
const CtxUserID TenantContextKey = "userID"
const CtxUserRole TenantContextKey = "userRole"
const CtxUserPermissions TenantContextKey = "userPermissions"

// TenantMiddleware extracts tenant_id from JWT and sets it in context
func TenantMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get user from JWT (set by JWTMiddleware)
			user := c.Get("user")
			if user == nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Unauthorized",
				})
			}

			token := user.(*jwt.Token)
			
			// Handle both *jwt.MapClaims and jwt.MapClaims
			var claims jwt.MapClaims
			switch v := token.Claims.(type) {
			case *jwt.MapClaims:
				claims = *v
			case jwt.MapClaims:
				claims = v
			default:
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Invalid token claims",
				})
			}

			// Extract tenant_id from claims
			tenantID, ok := claims["tenant_id"].(string)
			if !ok || tenantID == "" {
				return c.JSON(http.StatusBadRequest, map[string]string{
					"error": "Tenant ID not found in token",
				})
			}

			// Verify tenant exists and is active
			var exists bool
			err := db.DB.Get(&exists, `
				SELECT EXISTS(
					SELECT 1 FROM tenants 
					WHERE id = $1 AND status = 'active' AND deleted_at IS NULL
				)
			`, tenantID)
			if err != nil || !exists {
				return c.JSON(http.StatusForbidden, map[string]string{
					"error": "Invalid or inactive tenant",
				})
			}

			// Set tenant_id in context
			c.Set(string(CtxTenantID), tenantID)
			
			// Set user_id if available
			if userID, ok := claims["user_id"].(string); ok {
				c.Set(string(CtxUserID), userID)
			}

			// Set user role if available
			if role, ok := claims["role"].(string); ok {
				c.Set(string(CtxUserRole), role)
			}

			// Fetch and set user permissions for the tenant
			if userID, ok := claims["user_id"].(string); ok {
				permissions, err := GetUserPermissions(userID, tenantID)
				if err == nil {
					c.Set(string(CtxUserPermissions), permissions)
				} else {
					c.Logger().Warnf("Failed to get permissions for user %s in tenant %s: %v", userID, tenantID, err)
				}
			}

			return next(c)
		}
	}
}

// RequireTenantMembership ensures the user is part of a tenant
func RequireTenantMembership() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			tenantID := c.Get(string(CtxTenantID))
			if tenantID == nil || tenantID.(string) == "" {
				return c.JSON(http.StatusForbidden, map[string]string{"error": "Access denied: User not associated with a tenant"})
			}
			return next(c)
		}
	}
}

// RequirePermission checks if the user has the required permission
func RequirePermission(permissionKey string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			permissions, ok := c.Get(string(CtxUserPermissions)).([]string)
			if !ok {
				return c.JSON(http.StatusForbidden, map[string]string{"error": "Access denied: Permissions not loaded"})
			}

			hasPermission := false
			for _, p := range permissions {
				if p == permissionKey {
					hasPermission = true
					break
				}
			}

			if !hasPermission {
				return c.JSON(http.StatusForbidden, map[string]string{"error": "Access denied: Insufficient permissions"})
			}
			return next(c)
		}
	}
}

// GetUserPermissions fetches all permission keys for a given user and tenant
func GetUserPermissions(userID, tenantID string) ([]string, error) {
	var permissions []string
	query := `
		SELECT DISTINCT p.key
		FROM tenant_users tu
		JOIN roles r ON tu.role_id = r.id
		JOIN role_permissions rp ON r.id = rp.role_id
		JOIN permissions p ON rp.permission_id = p.id
		WHERE tu.user_id = $1 AND tu.tenant_id = $2 AND tu.status = 'active'
		AND tu.deleted_at IS NULL AND r.deleted_at IS NULL
	`
	err := db.DB.Select(&permissions, query, userID, tenantID)
	if err != nil {
		return nil, err
	}
	return permissions, nil
}
