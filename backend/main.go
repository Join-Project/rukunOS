package main

import (
	"net/http"
	"os"
	"strings"
	"rukunos-backend/db"
	"rukunos-backend/handlers"
	customMiddleware "rukunos-backend/middleware"
	"rukunos-backend/services"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Initialize Database
	db.Init()

	// Start background scheduler
	go func() {
		services.StartScheduler()
	}()

	e := EchoServer()
	
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	e.Logger.Fatal(e.Start(":" + port))
}

func EchoServer() *echo.Echo {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	
	// CORS Configuration
	corsOrigins := []string{"http://localhost:3000"}
	if corsEnv := os.Getenv("CORS_ALLOWED_ORIGINS"); corsEnv != "" {
		corsOrigins = append(corsOrigins, strings.Split(corsEnv, ",")...)
	}
	
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: corsOrigins,
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"message": "Welcome to API",
			"status":  "healthy",
		})
	})

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "up",
		})
	})

	// Auth Routes
	auth := e.Group("/api/auth")
	auth.POST("/register", handlers.Register)
	auth.POST("/login", handlers.Login)
	auth.GET("/google", handlers.GetGoogleAuthURL)
	auth.GET("/google/callback", handlers.GoogleAuthCallback)

	// Public: Create Tenant (for initial setup)
	e.POST("/api/tenants", handlers.CreateTenant)

	// Protected Routes
	api := e.Group("/api")
	api.Use(customMiddleware.JWTMiddleware())
	api.Use(customMiddleware.TenantMiddleware())
	api.Use(customMiddleware.RequireTenantMembership())
	
	// User routes
	api.GET("/me", handlers.GetCurrentUser)
	
	// Tenant routes
	api.GET("/tenants/:tenant_id", handlers.GetTenant)

	// Unit routes
	units := api.Group("/units")
	units.POST("", handlers.CreateUnit)
	units.GET("", handlers.ListUnits)
	units.GET("/:unit_id", handlers.GetUnit)
	units.PUT("/:unit_id", handlers.UpdateUnit)
	units.DELETE("/:unit_id", handlers.DeleteUnit)
	units.POST("/:unit_id/assign", handlers.AssignUserToUnit)

	// Role routes
	roles := api.Group("/roles")
	roles.POST("", handlers.CreateRole)
	roles.GET("", handlers.ListRoles)
	roles.GET("/:role_id", handlers.GetRole)
	roles.PUT("/:role_id", handlers.UpdateRole)
	roles.DELETE("/:role_id", handlers.DeleteRole)
	roles.GET("/permissions", handlers.ListPermissions)

	// User routes
	users := api.Group("/users")
	users.POST("", handlers.CreateUserByAdmin) // Admin create user
	users.GET("", handlers.ListUsers)
	users.GET("/:user_id", handlers.GetUser)
	users.PUT("/:user_id", handlers.UpdateUser)
	users.DELETE("/:user_id", handlers.RemoveUserFromTenant)

	// User role assignment (alternative endpoint)
	api.POST("/users/:user_id/roles", handlers.AssignRoleToUser)

	// Billing routes
	billing := api.Group("/billing")
	billing.GET("", handlers.ListBills)
	billing.POST("", handlers.CreateBill)
	billing.POST("/bulk", handlers.BulkCreateBills) // Bulk create bills
	billing.GET("/:bill_id", handlers.GetBill)
	billing.PUT("/:bill_id", handlers.UpdateBill)
	billing.DELETE("/:bill_id", handlers.DeleteBill)
	billing.POST("/:bill_id/payment", handlers.ProcessPayment)

	// Billing template routes
	billingTemplates := api.Group("/billing/templates")
	billingTemplates.GET("", handlers.ListBillingTemplates)
	billingTemplates.POST("", handlers.CreateBillingTemplate)
	billingTemplates.GET("/:template_id", handlers.GetBillingTemplate)
	billingTemplates.PUT("/:template_id", handlers.UpdateBillingTemplate)
	billingTemplates.DELETE("/:template_id", handlers.DeleteBillingTemplate)
	billingTemplates.POST("/:template_id/generate", handlers.GenerateBillsFromTemplate)

	// Announcement routes
	announcements := api.Group("/announcements")
	announcements.GET("", handlers.ListAnnouncements)
	announcements.POST("", handlers.CreateAnnouncement)
	announcements.GET("/:announcement_id", handlers.GetAnnouncement)
	announcements.PUT("/:announcement_id", handlers.UpdateAnnouncement)
	announcements.DELETE("/:announcement_id", handlers.DeleteAnnouncement)

	// Visitor routes
	visitors := api.Group("/visitors")
	visitors.GET("", handlers.ListVisitorLogs)
	visitors.POST("", handlers.CreateVisitorLog)
	visitors.POST("/:visitor_id/checkout", handlers.CheckOutVisitor)
	visitors.DELETE("/:visitor_id", handlers.DeleteVisitorLog)

	// Panic alert routes
	panicAlerts := api.Group("/panic-alerts")
	panicAlerts.GET("", handlers.ListPanicAlerts)
	panicAlerts.POST("", handlers.CreatePanicAlert)
	panicAlerts.PUT("/:alert_id", handlers.UpdatePanicAlert)

	// Complaint routes
	complaints := api.Group("/complaints")
	complaints.GET("", handlers.ListComplaints)
	complaints.POST("", handlers.CreateComplaint)
	complaints.GET("/:complaint_id", handlers.GetComplaint)
	complaints.PUT("/:complaint_id", handlers.UpdateComplaint)
	complaints.DELETE("/:complaint_id", handlers.DeleteComplaint)

	// Document request routes
	documents := api.Group("/document-requests")
	documents.GET("", handlers.ListDocumentRequests)
	documents.POST("", handlers.CreateDocumentRequest)
	documents.GET("/:request_id", handlers.GetDocumentRequest)
	documents.PUT("/:request_id", handlers.UpdateDocumentRequest)
	documents.DELETE("/:request_id", handlers.DeleteDocumentRequest)

	// Family routes (get family members in same unit)
	api.GET("/family", handlers.GetFamilyMembers)

	// Dashboard routes
	api.GET("/dashboard/warga/summary", handlers.GetWargaDashboardSummary)
	api.GET("/billing/dashboard", handlers.GetBillingDashboard)

	return e
}

