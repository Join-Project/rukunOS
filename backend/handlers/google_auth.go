package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
	"rukunos-backend/db"
	"rukunos-backend/models"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
)

var googleOAuthConfig *oauth2.Config

func init() {
	// Initialize Google OAuth config
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	redirectURL := os.Getenv("GOOGLE_REDIRECT_URL")
	
	if clientID == "" || clientSecret == "" {
		return
	}
	
	if redirectURL == "" {
		redirectURL = "http://localhost:3000/auth/google/callback"
	}

	googleOAuthConfig = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.google.com/o/oauth2/auth",
			TokenURL: "https://oauth2.googleapis.com/token",
		},
	}
}

func GetGoogleAuthURL(c echo.Context) error {
	if googleOAuthConfig == nil {
		return c.JSON(http.StatusServiceUnavailable, map[string]string{
			"error": "Google OAuth is not configured",
		})
	}

	state := generateStateToken()
	url := googleOAuthConfig.AuthCodeURL(state, oauth2.AccessTypeOnline)
	
	return c.JSON(http.StatusOK, map[string]string{
		"auth_url": url,
		"state":    state,
	})
}

func GoogleAuthCallback(c echo.Context) error {
	if googleOAuthConfig == nil {
		return c.JSON(http.StatusServiceUnavailable, map[string]string{
			"error": "Google OAuth is not configured",
		})
	}

	code := c.QueryParam("code")
	if code == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Authorization code not provided",
		})
	}

	// Exchange code for token
	token, err := googleOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to exchange token: " + err.Error(),
		})
	}

	// Get user info from Google
	userInfo, err := getGoogleUserInfo(token.AccessToken)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to get user info: " + err.Error(),
		})
	}

	// Find or create user
	user, err := findOrCreateGoogleUser(userInfo)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create/find user: " + err.Error(),
		})
	}

	// Get user's tenants (for now, we'll use first tenant or create one)
	// TODO: Implement proper tenant selection for OAuth users
	var tenantIDs []string
	err = db.DB.Select(&tenantIDs, `
		SELECT tenant_id 
		FROM tenant_users 
		WHERE user_id = $1 
		AND deleted_at IS NULL 
		AND status = 'active'
	`, user.ID)
	
	var selectedTenantID string
	if len(tenantIDs) > 0 {
		selectedTenantID = tenantIDs[0]
	} else {
		// If user has no tenant, we can't generate token
		// This should be handled by creating a tenant or joining one
		return c.JSON(http.StatusForbidden, map[string]string{
			"error": "User is not a member of any tenant",
		})
	}

	// Generate JWT
	jwtToken := jwt.New(jwt.SigningMethodHS256)
	claims := jwtToken.Claims.(jwt.MapClaims)
	claims["user_id"] = user.ID
	claims["tenant_id"] = selectedTenantID
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "secret"
	}

	t, err := jwtToken.SignedString([]byte(jwtSecret))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to generate token",
		})
	}

	// Redirect to frontend
	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:3000"
	}
	
	// Get user tenants for response (reuse tenantIDs from above)
	responseData := map[string]interface{}{
		"token":     t,
		"user":      user,
		"tenant_id": selectedTenantID,
	}
	if len(tenantIDs) > 1 {
		responseData["tenants"] = tenantIDs
	}
	
	userJSON, _ := json.Marshal(responseData)
	redirectHTML := fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head><title>Redirecting...</title></head>
		<body>
			<script>
				const data = %s;
				window.location.href = '%s/auth/google/callback?token=' + encodeURIComponent(data.token) + '&data=' + encodeURIComponent(JSON.stringify(data));
			</script>
			<p>Redirecting...</p>
		</body>
		</html>
	`, userJSON, frontendURL)

	return c.HTML(http.StatusOK, redirectHTML)
}

type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
}

func getGoogleUserInfo(accessToken string) (*GoogleUserInfo, error) {
	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + accessToken)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var userInfo GoogleUserInfo
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, err
	}

	return &userInfo, nil
}

func findOrCreateGoogleUser(googleUser *GoogleUserInfo) (*models.User, error) {
	var user models.User
	var passwordHash sql.NullString
	var googleID sql.NullString
	var authProvider sql.NullString
	
	// Try to find by Google ID
	query := `SELECT id, email, password_hash, full_name, google_id, auth_provider, status, created_at, updated_at
	          FROM users WHERE google_id = $1 AND deleted_at IS NULL`
	var status sql.NullString
	err := db.DB.QueryRow(query, googleUser.ID).Scan(
		&user.ID, &user.Email, &passwordHash, &user.FullName, 
		&googleID, &authProvider, &status, &user.CreatedAt, &user.UpdatedAt)
	
	if err == nil {
		if googleID.Valid {
			user.GoogleID = &googleID.String
		}
		if authProvider.Valid {
			user.AuthProvider = authProvider.String
		} else {
			user.AuthProvider = "google"
		}
		if status.Valid {
			user.Status = status.String
		} else {
			user.Status = "active"
		}
		return &user, nil
	}
	
	if err != sql.ErrNoRows {
		return nil, err
	}
	
	// Try to find by email
	query = `SELECT id, email, password_hash, full_name, google_id, auth_provider, status, created_at, updated_at
	         FROM users WHERE email = $1 AND deleted_at IS NULL`
	err = db.DB.QueryRow(query, googleUser.Email).Scan(
		&user.ID, &user.Email, &passwordHash, &user.FullName, 
		&googleID, &authProvider, &status, &user.CreatedAt, &user.UpdatedAt)
	
	if err == nil {
		if authProvider.Valid && authProvider.String != "google" {
			return nil, fmt.Errorf("Email sudah terdaftar dengan akun email/password")
		}
		if googleID.Valid {
			user.GoogleID = &googleID.String
		}
		if status.Valid {
			user.Status = status.String
		} else {
			user.Status = "active"
		}
		return &user, nil
	}
	
	if err != sql.ErrNoRows {
		return nil, err
	}
	
	// Create new user
	fullName := googleUser.Name
	if fullName == "" {
		fullName = googleUser.Email
	}
	
	insertQuery := `INSERT INTO users (email, full_name, google_id, auth_provider, status) 
	                VALUES ($1, $2, $3, 'google', 'active') 
	                RETURNING id, email, full_name, auth_provider, status, created_at, updated_at`
	var insertedStatus sql.NullString
	err = db.DB.QueryRow(insertQuery, googleUser.Email, fullName, googleUser.ID).Scan(
		&user.ID, &user.Email, &user.FullName, &authProvider, &insertedStatus, &user.CreatedAt, &user.UpdatedAt)
	
	if err != nil {
		return nil, err
	}
	
	user.GoogleID = &googleUser.ID
	user.AuthProvider = "google"
	if insertedStatus.Valid {
		user.Status = insertedStatus.String
	} else {
		user.Status = "active"
	}
	
	return &user, nil
}

func generateStateToken() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

