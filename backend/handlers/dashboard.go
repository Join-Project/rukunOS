package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"
	"rukunos-backend/db"
	"rukunos-backend/middleware"

	"github.com/labstack/echo/v4"
)

// GetWargaDashboardSummary returns summary data for warga dashboard
func GetWargaDashboardSummary(c echo.Context) error {
	tenantID := c.Get(string(middleware.CtxTenantID)).(string)
	userID := c.Get(string(middleware.CtxUserID)).(string)

	// Get user's unit_id
	var unitID sql.NullString
	err := db.DB.Get(&unitID, `
		SELECT unit_id FROM tenant_users
		WHERE user_id = $1 AND tenant_id = $2 AND deleted_at IS NULL
	`, userID, tenantID)

	var unitCode sql.NullString
	if err == nil && unitID.Valid {
		db.DB.Get(&unitCode, `SELECT code FROM units WHERE id = $1`, unitID.String)
	}

	// Get active bills (pending)
	var totalPendingAmount float64
	var pendingBillsCount int
	var activeBill map[string]interface{}

	err = db.DB.QueryRow(`
		SELECT 
			COALESCE(SUM(amount + COALESCE(late_fee, 0)), 0) as total_amount,
			COUNT(*) as count
		FROM bills
		WHERE tenant_id = $1 AND unit_id = $2 AND status = 'pending' AND deleted_at IS NULL
	`, tenantID, unitID.String).Scan(&totalPendingAmount, &pendingBillsCount)

	// Get latest pending bill
	if pendingBillsCount > 0 {
		var billID, category, period string
		var amount, lateFee float64
		var dueDate sql.NullTime

		err = db.DB.QueryRow(`
			SELECT id, category, period, amount, late_fee, due_date
			FROM bills
			WHERE tenant_id = $1 AND unit_id = $2 AND status = 'pending' AND deleted_at IS NULL
			ORDER BY due_date ASC NULLS LAST, created_at DESC
			LIMIT 1
		`, tenantID, unitID.String).Scan(&billID, &category, &period, &amount, &lateFee, &dueDate)

		if err == nil {
			activeBill = map[string]interface{}{
				"id":       billID,
				"category": category,
				"period":   period,
				"amount":   amount,
				"late_fee": lateFee,
			}
			if dueDate.Valid {
				activeBill["due_date"] = dueDate.Time.Format("2006-01-02")
			}
		}
	}

	// Get recent announcements (last 3)
	announcements := []map[string]interface{}{}
	rows, err := db.DB.Query(`
		SELECT 
			a.id, a.title, a.content, a.priority, a.category, a.created_at,
			u.full_name as author_name
		FROM announcements a
		LEFT JOIN users u ON a.author_id = u.id
		WHERE a.tenant_id = $1 AND a.deleted_at IS NULL
		ORDER BY a.is_pinned DESC, a.created_at DESC
		LIMIT 3
	`, tenantID)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var id, title, content, priority, category, authorName sql.NullString
			var createdAt time.Time

			err := rows.Scan(&id, &title, &content, &priority, &category, &createdAt, &authorName)
			if err != nil {
				continue
			}

			announcement := map[string]interface{}{
				"id":         id.String,
				"title":      title.String,
				"content":    content.String,
				"priority":   priority.String,
				"created_at": createdAt.Format(time.RFC3339),
			}
			if category.Valid {
				announcement["category"] = category.String
			}
			if authorName.Valid {
				announcement["author_name"] = authorName.String
			}
			announcements = append(announcements, announcement)
		}
	}

	// Get recent activities (last 5)
	activities := []map[string]interface{}{}

	// Recent paid bills
	billRows, err := db.DB.Query(`
		SELECT category, period, paid_at, amount
		FROM bills
		WHERE tenant_id = $1 AND unit_id = $2 AND status = 'paid' AND deleted_at IS NULL
		ORDER BY paid_at DESC
		LIMIT 3
	`, tenantID, unitID.String)
	if err == nil {
		defer billRows.Close()
		for billRows.Next() {
			var category, period string
			var paidAt sql.NullTime
			var amount float64

			err := billRows.Scan(&category, &period, &paidAt, &amount)
			if err != nil {
				continue
			}

			if paidAt.Valid {
				activities = append(activities, map[string]interface{}{
					"type":      "payment",
					"title":     "Pembayaran " + category + " Diterima",
					"date":      paidAt.Time.Format(time.RFC3339),
					"icon":      "check",
					"iconColor": "green",
				})
			}
		}
	}

	// Recent document requests
	docRows, err := db.DB.Query(`
		SELECT document_type, status, created_at
		FROM document_requests
		WHERE tenant_id = $1 AND user_id = $2 AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT 2
	`, tenantID, userID)
	if err == nil {
		defer docRows.Close()
		for docRows.Next() {
			var docType, status string
			var createdAt time.Time

			err := docRows.Scan(&docType, &status, &createdAt)
			if err != nil {
				continue
			}

			var iconColor string
			if status == "completed" {
				iconColor = "green"
			} else if status == "approved" {
				iconColor = "blue"
			} else {
				iconColor = "yellow"
			}

			activities = append(activities, map[string]interface{}{
				"type":      "document",
				"title":     "Surat " + docType + " " + getDocumentStatusLabel(status),
				"date":      createdAt.Format(time.RFC3339),
				"icon":      "document",
				"iconColor": iconColor,
			})
		}
	}

	// Sort activities by date
	// (In a real implementation, you'd want to merge and sort all activities)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"unit_code":          unitCode.String,
		"total_pending_amount": totalPendingAmount,
		"pending_bills_count":  pendingBillsCount,
		"active_bill":          activeBill,
		"announcements":        announcements,
		"activities":           activities[:min(5, len(activities))], // Limit to 5
	})
}

func getDocumentStatusLabel(status string) string {
	labels := map[string]string{
		"pending":   "Diproses",
		"approved":  "Disetujui",
		"rejected":  "Ditolak",
		"completed": "Selesai",
	}
	if label, ok := labels[status]; ok {
		return label
	}
	return status
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// GetBillingDashboard returns financial dashboard data for admin/bendahara
func GetBillingDashboard(c echo.Context) error {
	tenantID := c.Get(string(middleware.CtxTenantID)).(string)

	// Get period filter (optional)
	// If period is provided, filter by created_at month instead of period field
	// because period field might be in different format (e.g., "Januari 2025")
	periodFilter := c.QueryParam("period") // Format: YYYY-MM or YYYY

	// Build period filter for SQL - use created_at instead of period field
	var periodWhere string
	var periodArgs []interface{}
	if periodFilter != "" {
		if len(periodFilter) == 7 { // YYYY-MM
			periodWhere = "AND TO_CHAR(b.created_at, 'YYYY-MM') = $1"
			periodArgs = []interface{}{periodFilter}
		} else if len(periodFilter) == 4 { // YYYY
			periodWhere = "AND TO_CHAR(b.created_at, 'YYYY') = $1"
			periodArgs = []interface{}{periodFilter}
		}
	} else {
		// Default: no filter (show all data)
		periodFilter = ""
		periodWhere = ""
		periodArgs = []interface{}{}
	}

	// Get summary statistics
	var summary struct {
		TotalPendingCount  int     `db:"total_pending_count"`
		TotalPendingAmount float64 `db:"total_pending_amount"`
		TotalPaidCount     int     `db:"total_paid_count"`
		TotalPaidAmount    float64 `db:"total_paid_amount"`
		TotalOverdueCount  int     `db:"total_overdue_count"`
		TotalOverdueAmount float64 `db:"total_overdue_amount"`
		TotalBillsCount    int     `db:"total_bills_count"`
		TotalBillsAmount   float64 `db:"total_bills_amount"`
	}

	// Build query with proper parameter indexing
	argIndex := len(periodArgs) + 1
	query := fmt.Sprintf(`
		SELECT 
			COUNT(CASE WHEN b.status = 'pending' THEN 1 END) as total_pending_count,
			COALESCE(SUM(CASE WHEN b.status = 'pending' THEN b.amount + COALESCE(b.late_fee, 0) ELSE 0 END), 0) as total_pending_amount,
			COUNT(CASE WHEN b.status = 'paid' THEN 1 END) as total_paid_count,
			COALESCE(SUM(CASE WHEN b.status = 'paid' THEN b.amount + COALESCE(b.late_fee, 0) ELSE 0 END), 0) as total_paid_amount,
			COUNT(CASE WHEN b.status = 'overdue' THEN 1 END) as total_overdue_count,
			COALESCE(SUM(CASE WHEN b.status = 'overdue' THEN b.amount + COALESCE(b.late_fee, 0) ELSE 0 END), 0) as total_overdue_amount,
			COUNT(*) as total_bills_count,
			COALESCE(SUM(b.amount + COALESCE(b.late_fee, 0)), 0) as total_bills_amount
		FROM bills b
		WHERE b.tenant_id = $%d AND b.deleted_at IS NULL %s
	`, argIndex, periodWhere)

	args := append(periodArgs, tenantID)
	err := db.DB.Get(&summary, query, args...)
	if err != nil {
		c.Logger().Errorf("Error executing dashboard summary query: %v, query: %s, args: %v", err, query, args)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error: " + err.Error()})
	}

	// Calculate collection rate
	collectionRate := 0.0
	if summary.TotalBillsCount > 0 {
		collectionRate = (float64(summary.TotalPaidCount) / float64(summary.TotalBillsCount)) * 100
	}

	// Get monthly trend (last 12 months)
	// Generate all 12 months and LEFT JOIN with bills data to ensure all months are shown
	monthlyTrend := []map[string]interface{}{}
	trendRows, err := db.DB.Query(`
		WITH months AS (
			SELECT TO_CHAR(month_series, 'YYYY-MM') as month
			FROM generate_series(
				DATE_TRUNC('month', CURRENT_DATE - INTERVAL '11 months'),
				DATE_TRUNC('month', CURRENT_DATE),
				INTERVAL '1 month'
			) AS month_series
		)
		SELECT 
			m.month,
			COALESCE(COUNT(CASE WHEN b.status = 'paid' THEN 1 END), 0) as paid_count,
			COALESCE(SUM(CASE WHEN b.status = 'paid' THEN b.amount + COALESCE(b.late_fee, 0) ELSE 0 END), 0) as paid_amount,
			COALESCE(COUNT(CASE WHEN b.status = 'pending' THEN 1 END), 0) as pending_count,
			COALESCE(SUM(CASE WHEN b.status = 'pending' THEN b.amount + COALESCE(b.late_fee, 0) ELSE 0 END), 0) as pending_amount,
			COALESCE(COUNT(CASE WHEN b.status = 'overdue' THEN 1 END), 0) as overdue_count,
			COALESCE(SUM(CASE WHEN b.status = 'overdue' THEN b.amount + COALESCE(b.late_fee, 0) ELSE 0 END), 0) as overdue_amount
		FROM months m
		LEFT JOIN bills b ON TO_CHAR(b.created_at, 'YYYY-MM') = m.month
			AND b.tenant_id = $1 
			AND b.deleted_at IS NULL
		GROUP BY m.month
		ORDER BY m.month ASC
	`, tenantID)
	if err != nil {
		c.Logger().Errorf("Error executing monthly trend query: %v", err)
		// Continue with empty trend data
		monthlyTrend = []map[string]interface{}{}
	} else {
		defer trendRows.Close()
		for trendRows.Next() {
			var monthStr string
			var paidCount, pendingCount, overdueCount int
			var paidAmount, pendingAmount, overdueAmount float64

			err := trendRows.Scan(&monthStr, &paidCount, &paidAmount, &pendingCount, &pendingAmount, &overdueCount, &overdueAmount)
			if err != nil {
				c.Logger().Errorf("Error scanning monthly trend row: %v", err)
				continue
			}

			monthlyTrend = append(monthlyTrend, map[string]interface{}{
				"month":          monthStr,
				"paid_count":     paidCount,
				"paid_amount":    paidAmount,
				"pending_count":  pendingCount,
				"pending_amount": pendingAmount,
				"overdue_count":  overdueCount,
				"overdue_amount": overdueAmount,
			})
		}
	}

	// Get top 10 overdue units
	topOverdue := []map[string]interface{}{}
	overdueRows, err := db.DB.Query(`
		SELECT 
			u.code as unit_code,
			u.type as unit_type,
			COUNT(*) as overdue_count,
			COALESCE(SUM(b.amount + COALESCE(b.late_fee, 0)), 0) as total_amount
		FROM bills b
		INNER JOIN units u ON b.unit_id = u.id
		WHERE b.tenant_id = $1 
		AND b.status = 'overdue' 
		AND b.deleted_at IS NULL
		GROUP BY u.id, u.code, u.type
		ORDER BY total_amount DESC
		LIMIT 10
	`, tenantID)
	if err != nil {
		c.Logger().Errorf("Error executing top overdue query: %v", err)
		// Continue with empty overdue data
		topOverdue = []map[string]interface{}{}
	} else {
		defer overdueRows.Close()
		for overdueRows.Next() {
			var unitCode, unitType string
			var overdueCount int
			var totalAmount float64

			err := overdueRows.Scan(&unitCode, &unitType, &overdueCount, &totalAmount)
			if err != nil {
				continue
			}

			topOverdue = append(topOverdue, map[string]interface{}{
				"unit_code":     unitCode,
				"unit_type":     unitType,
				"overdue_count": overdueCount,
				"total_amount":  totalAmount,
			})
		}
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"summary": map[string]interface{}{
			"collection_rate":      collectionRate,
			"total_pending": map[string]interface{}{
				"count":  summary.TotalPendingCount,
				"amount": summary.TotalPendingAmount,
			},
			"total_paid": map[string]interface{}{
				"count":  summary.TotalPaidCount,
				"amount": summary.TotalPaidAmount,
			},
			"total_overdue": map[string]interface{}{
				"count":  summary.TotalOverdueCount,
				"amount": summary.TotalOverdueAmount,
			},
			"total_bills": map[string]interface{}{
				"count":  summary.TotalBillsCount,
				"amount": summary.TotalBillsAmount,
			},
		},
		"monthly_trend": monthlyTrend,
		"top_overdue":   topOverdue,
		"period":        periodFilter,
	})
}





