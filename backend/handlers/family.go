package handlers

import (
	"database/sql"
	"net/http"
	"rukunos-backend/db"
	"rukunos-backend/middleware"
	"time"

	"github.com/labstack/echo/v4"
)

// GetFamilyMembers gets all family members (users) in the same unit
func GetFamilyMembers(c echo.Context) error {
	tenantID := c.Get(string(middleware.CtxTenantID)).(string)
	userID := c.Get(string(middleware.CtxUserID)).(string)

	// Get user's unit_id
	var unitID sql.NullString
	err := db.DB.Get(&unitID, `
		SELECT unit_id FROM tenant_users
		WHERE user_id = $1 AND tenant_id = $2 AND deleted_at IS NULL
	`, userID, tenantID)

	if err != nil || !unitID.Valid {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"family_members": []interface{}{},
			"unit_code":      nil,
		})
	}

	// Get unit code
	var unitCode sql.NullString
	db.DB.Get(&unitCode, `SELECT code FROM units WHERE id = $1`, unitID.String)

	// Get all users in the same unit
	query := `
		SELECT 
			u.id, u.full_name, u.email, u.phone, u.avatar_url,
			tu.role_id, r.name as role_name,
			tu.created_at as joined_at
		FROM tenant_users tu
		INNER JOIN users u ON tu.user_id = u.id
		LEFT JOIN roles r ON tu.role_id = r.id
		WHERE tu.tenant_id = $1 AND tu.unit_id = $2 
		AND tu.deleted_at IS NULL AND u.deleted_at IS NULL
		ORDER BY tu.created_at ASC
	`

	rows, err := db.DB.Query(query, tenantID, unitID.String)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}
	defer rows.Close()

	familyMembers := []map[string]interface{}{}
	for rows.Next() {
		var id, fullName, email, phone, avatarURL sql.NullString
		var roleID, roleName sql.NullString
		var joinedAt time.Time

		err := rows.Scan(&id, &fullName, &email, &phone, &avatarURL, &roleID, &roleName, &joinedAt)
		if err != nil {
			continue
		}

		memberData := map[string]interface{}{
			"id":        id.String,
			"joined_at": joinedAt.Format(time.RFC3339),
		}

		if fullName.Valid {
			memberData["full_name"] = fullName.String
		}
		if email.Valid {
			memberData["email"] = email.String
		}
		if phone.Valid {
			memberData["phone"] = phone.String
		}
		if avatarURL.Valid {
			memberData["avatar_url"] = avatarURL.String
		}
		if roleID.Valid {
			memberData["role_id"] = roleID.String
		}
		if roleName.Valid {
			memberData["role_name"] = roleName.String
			// Determine relationship based on role
			if roleName.String == "Warga" {
				memberData["relationship"] = "Kepala Keluarga"
			} else {
				memberData["relationship"] = roleName.String
			}
		}

		familyMembers = append(familyMembers, memberData)
	}

	// Mark the first member (oldest) as head of family
	if len(familyMembers) > 0 {
		familyMembers[0]["relationship"] = "Kepala Keluarga"
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"family_members": familyMembers,
		"unit_code":      unitCode.String,
	})
}

