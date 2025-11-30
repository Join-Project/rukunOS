package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"
	"rukunos-backend/db"
	"rukunos-backend/middleware"
	"rukunos-backend/models"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// ListBills lists all bills for the tenant (admin) or current user (warga)
func ListBills(c echo.Context) error {
	tenantID := c.Get(string(middleware.CtxTenantID)).(string)
	userID := c.Get(string(middleware.CtxUserID)).(string)

	// Get query parameters
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit < 1 {
		limit = 20
	}
	status := c.QueryParam("status")
	unitID := c.QueryParam("unit_id")
	search := c.QueryParam("search")

	offset := (page - 1) * limit

	// Build query
	query := `
		SELECT 
			b.id, b.tenant_id, b.unit_id, b.category, b.period, b.amount, b.late_fee,
			b.due_date, b.status, b.paid_at, b.payment_method, b.payment_reference,
			b.notes, b.created_by, b.created_at, b.updated_at,
			b.bill_number, (b.amount + COALESCE(b.late_fee, 0)) as total_amount,
			u.code as unit_code, u.type as unit_type
		FROM bills b
		INNER JOIN units u ON b.unit_id = u.id
		WHERE b.tenant_id = $1 AND b.deleted_at IS NULL
	`
	args := []interface{}{tenantID}
	argIndex := 2

	// Filter by user (for warga - only show their bills)
	// Check if user has permission to view all bills (admin/bendahara)
	// For now, if status is not provided and user is not admin, filter by user's unit
	var isAdmin bool
	err := db.DB.Get(&isAdmin, `
		SELECT EXISTS(
			SELECT 1 FROM tenant_users tu
			JOIN roles r ON tu.role_id = r.id
			WHERE tu.user_id = $1 AND tu.tenant_id = $2
			AND (r.name = 'Admin' OR r.name = 'Bendahara')
			AND tu.deleted_at IS NULL AND r.deleted_at IS NULL
		)
	`, userID, tenantID)
	if err == nil && !isAdmin {
		// Filter by user's unit
		var userUnitID sql.NullString
		err = db.DB.Get(&userUnitID, `
			SELECT unit_id FROM tenant_users
			WHERE user_id = $1 AND tenant_id = $2 AND deleted_at IS NULL
		`, userID, tenantID)
		if err == nil && userUnitID.Valid {
			query += ` AND b.unit_id = $` + strconv.Itoa(argIndex)
			args = append(args, userUnitID.String)
			argIndex++
		}
	}

	if status != "" {
		query += ` AND b.status = $` + strconv.Itoa(argIndex)
		args = append(args, status)
		argIndex++
	}

	if unitID != "" {
		query += ` AND b.unit_id = $` + strconv.Itoa(argIndex)
		args = append(args, unitID)
		argIndex++
	}

	if search != "" {
		searchPattern := "%" + search + "%"
		query += ` AND (u.code ILIKE $` + strconv.Itoa(argIndex) + ` OR b.category ILIKE $` + strconv.Itoa(argIndex) + `)`
		args = append(args, searchPattern)
		argIndex++
	}

	query += ` ORDER BY b.due_date DESC NULLS LAST, b.created_at DESC LIMIT $` + strconv.Itoa(argIndex) + ` OFFSET $` + strconv.Itoa(argIndex+1)
	args = append(args, limit, offset)

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		c.Logger().Errorf("Error executing bills query: %v, query: %s, args: %v", err, query, args)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error: " + err.Error()})
	}
	defer rows.Close()

	var bills []map[string]interface{}
	for rows.Next() {
		var bill models.Bill
		var unitCode, unitType sql.NullString
		var paidAt sql.NullTime
		var paymentMethod, paymentReference, notes, createdBy sql.NullString
		var billNumber sql.NullString
		var totalAmount float64
		
		err := rows.Scan(
			&bill.ID, &bill.TenantID, &bill.UnitID, &bill.Category, &bill.Period,
			&bill.Amount, &bill.LateFee, &bill.DueDate, &bill.Status,
			&paidAt, &paymentMethod, &paymentReference,
			&notes, &createdBy, &bill.CreatedAt, &bill.UpdatedAt,
			&billNumber, &totalAmount,
			&unitCode, &unitType,
		)
		if err != nil {
			c.Logger().Errorf("Error scanning bill row: %v", err)
			continue
		}

		billData := map[string]interface{}{
			"id":          bill.ID,
			"unit_id":     bill.UnitID,
			"category":    bill.Category,
			"period":      bill.Period,
			"amount":      bill.Amount,
			"late_fee":    bill.LateFee,
			"total_amount": totalAmount,
			"status":      bill.Status,
			"created_at":  bill.CreatedAt.Format(time.RFC3339),
		}

		if billNumber.Valid {
			billData["bill_number"] = billNumber.String
		}
		if bill.DueDate.Valid {
			billData["due_date"] = bill.DueDate.Time.Format("2006-01-02")
		}
		if unitCode.Valid {
			billData["unit_code"] = unitCode.String
		}
		if unitType.Valid {
			billData["unit_type"] = unitType.String
		}
		if paidAt.Valid {
			billData["paid_at"] = paidAt.Time.Format(time.RFC3339)
		}
		if paymentMethod.Valid {
			billData["payment_method"] = paymentMethod.String
		}
		if paymentReference.Valid {
			billData["payment_reference"] = paymentReference.String
		}
		if notes.Valid {
			billData["notes"] = notes.String
		}

		bills = append(bills, billData)
	}

	// Get total count
	countQuery := `SELECT COUNT(*) FROM bills b
		INNER JOIN units u ON b.unit_id = u.id
		WHERE b.tenant_id = $1 AND b.deleted_at IS NULL`
	countArgs := []interface{}{tenantID}
	countArgIndex := 2

	if !isAdmin {
		var userUnitID sql.NullString
		err = db.DB.Get(&userUnitID, `
			SELECT unit_id FROM tenant_users
			WHERE user_id = $1 AND tenant_id = $2 AND deleted_at IS NULL
		`, userID, tenantID)
		if err == nil && userUnitID.Valid {
			countQuery += ` AND b.unit_id = $` + strconv.Itoa(countArgIndex)
			countArgs = append(countArgs, userUnitID.String)
			countArgIndex++
		}
	}

	if status != "" {
		countQuery += ` AND b.status = $` + strconv.Itoa(countArgIndex)
		countArgs = append(countArgs, status)
		countArgIndex++
	}
	if unitID != "" {
		countQuery += ` AND b.unit_id = $` + strconv.Itoa(countArgIndex)
		countArgs = append(countArgs, unitID)
		countArgIndex++
	}
	if search != "" {
		searchPattern := "%" + search + "%"
		countQuery += ` AND (u.code ILIKE $` + strconv.Itoa(countArgIndex) + ` OR b.category ILIKE $` + strconv.Itoa(countArgIndex) + `)`
		countArgs = append(countArgs, searchPattern)
		countArgIndex++
	}

	var total int
	err = db.DB.Get(&total, countQuery, countArgs...)
	if err != nil {
		total = len(bills)
	}

	totalPages := (total + limit - 1) / limit

	return c.JSON(http.StatusOK, map[string]interface{}{
		"bills": bills,
		"pagination": map[string]interface{}{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": totalPages,
		},
	})
}

// GetBill gets a bill by ID
func GetBill(c echo.Context) error {
	tenantID := c.Get(string(middleware.CtxTenantID)).(string)
	billID := c.Param("bill_id")

	var bill models.Bill
	err := db.DB.QueryRow(`
		SELECT 
			b.id, b.tenant_id, b.unit_id, b.category, b.period, b.amount, b.late_fee,
			b.due_date, b.status, b.paid_at, b.payment_method, b.payment_reference,
			b.notes, b.created_by, b.created_at, b.updated_at,
			u.code as unit_code, u.type as unit_type
		FROM bills b
		INNER JOIN units u ON b.unit_id = u.id
		WHERE b.id = $2 AND b.tenant_id = $1 AND b.deleted_at IS NULL
	`, tenantID, billID).Scan(
		&bill.ID, &bill.TenantID, &bill.UnitID, &bill.Category, &bill.Period,
		&bill.Amount, &bill.LateFee, &bill.DueDate, &bill.Status,
		&bill.PaidAt, &bill.PaymentMethod, &bill.PaymentReference,
		&bill.Notes, &bill.CreatedBy, &bill.CreatedAt, &bill.UpdatedAt,
		&bill.UnitCode, &bill.UnitType,
	)

	if err == sql.ErrNoRows {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Bill not found"})
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error: " + err.Error()})
	}

	billData := map[string]interface{}{
		"id":        bill.ID,
		"unit_id":   bill.UnitID,
		"category":  bill.Category,
		"period":    bill.Period,
		"amount":    bill.Amount,
		"late_fee":  bill.LateFee,
		"status":    bill.Status,
		"created_at": bill.CreatedAt,
		"updated_at": bill.UpdatedAt,
	}
	
	if bill.DueDate.Valid {
		billData["due_date"] = bill.DueDate.Time.Format("2006-01-02")
	}

	if bill.UnitCode.Valid {
		billData["unit_code"] = bill.UnitCode.String
	}
	if bill.UnitType.Valid {
		billData["unit_type"] = bill.UnitType.String
	}
	if bill.PaidAt.Valid {
		billData["paid_at"] = bill.PaidAt.Time
	}
	if bill.PaymentMethod.Valid {
		billData["payment_method"] = bill.PaymentMethod.String
	}
	if bill.PaymentReference.Valid {
		billData["payment_reference"] = bill.PaymentReference.String
	}
	if bill.Notes.Valid {
		billData["notes"] = bill.Notes.String
	}

	return c.JSON(http.StatusOK, billData)
}

// CreateBill creates a new bill (admin/bendahara only)
func CreateBill(c echo.Context) error {
	tenantID := c.Get(string(middleware.CtxTenantID)).(string)
	userID := c.Get(string(middleware.CtxUserID)).(string)

	req := new(models.CreateBillRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// Validate unit belongs to tenant
	var unitExists bool
	err := db.DB.Get(&unitExists, `
		SELECT EXISTS(
			SELECT 1 FROM units 
			WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL
		)
	`, req.UnitID, tenantID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}
	if !unitExists {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Unit not found"})
	}

	// Parse due date (optional)
	var dueDate sql.NullTime
	if req.DueDate != nil && *req.DueDate != "" {
		parsedDate, err := time.Parse("2006-01-02", *req.DueDate)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid due_date format. Use YYYY-MM-DD"})
		}
		dueDate = sql.NullTime{Time: parsedDate, Valid: true}
	}

	lateFee := 0.0
	if req.LateFee != nil {
		lateFee = *req.LateFee
	}

	// Create bill
	billID := uuid.New().String()
	
	// Build INSERT query based on whether due_date is provided
	var query string
	var args []interface{}
	
	if dueDate.Valid {
		query = `INSERT INTO bills (id, tenant_id, unit_id, category, period, amount, late_fee, due_date, status, notes, created_by)
		          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, 'pending', $9, $10)
		          RETURNING id, created_at`
		args = []interface{}{billID, tenantID, req.UnitID, req.Category, req.Period, req.Amount, lateFee, dueDate.Time, req.Notes, userID}
	} else {
		query = `INSERT INTO bills (id, tenant_id, unit_id, category, period, amount, late_fee, status, notes, created_by)
		          VALUES ($1, $2, $3, $4, $5, $6, $7, 'pending', $8, $9)
		          RETURNING id, created_at`
		args = []interface{}{billID, tenantID, req.UnitID, req.Category, req.Period, req.Amount, lateFee, req.Notes, userID}
	}
	
	var createdAt time.Time
	var returnedBillID string
	err = db.DB.QueryRow(query, args...).Scan(&returnedBillID, &createdAt)
	if returnedBillID != "" {
		billID = returnedBillID
	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create bill: " + err.Error()})
	}

	// Return created bill - get full details
	c.SetParamNames("bill_id")
	c.SetParamValues(billID)
	return GetBill(c)
}

// BulkCreateBills creates multiple bills at once (for mass generation)
func BulkCreateBills(c echo.Context) error {
	tenantID := c.Get(string(middleware.CtxTenantID)).(string)
	userID := c.Get(string(middleware.CtxUserID)).(string)

	var req struct {
		Category  string   `json:"category" validate:"required"`
		Period    string   `json:"period" validate:"required"`
		Amount    float64  `json:"amount" validate:"required,min=0"`
		LateFee   *float64 `json:"late_fee,omitempty"`
		DueDate   *string  `json:"due_date,omitempty"` // Optional
		Notes     *string  `json:"notes,omitempty"`
		UnitIDs   []string `json:"unit_ids" validate:"required,min=1"` // Array of unit IDs
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request: " + err.Error()})
	}

	// Validate required fields
	if req.Category == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "category is required"})
	}
	if req.Period == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "period is required"})
	}
	if req.Amount <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "amount must be greater than 0"})
	}
	if len(req.UnitIDs) == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "unit_ids is required and must contain at least one unit"})
	}

	// Parse due date (optional)
	var dueDate sql.NullTime
	if req.DueDate != nil && *req.DueDate != "" {
		parsedDate, err := time.Parse("2006-01-02", *req.DueDate)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid due_date format. Use YYYY-MM-DD"})
		}
		dueDate = sql.NullTime{Time: parsedDate, Valid: true}
	}

	lateFee := 0.0
	if req.LateFee != nil {
		lateFee = *req.LateFee
	}

	// Start transaction
	tx, err := db.DB.Beginx()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to start transaction"})
	}
	defer tx.Rollback()

	// Validate all units belong to tenant
	for _, unitID := range req.UnitIDs {
		var unitExists bool
		err = tx.Get(&unitExists, `
			SELECT EXISTS(
				SELECT 1 FROM units 
				WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL
			)
		`, unitID, tenantID)
		if err != nil {
			tx.Rollback()
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
		}
		if !unitExists {
			tx.Rollback()
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Unit not found: " + unitID})
		}
	}

	// Create bills for each unit
	createdBills := []string{}
	for _, unitID := range req.UnitIDs {
		billID := uuid.New().String()
		
		// Build INSERT query based on whether due_date is provided
		var insertQuery string
		var args []interface{}
		
		if dueDate.Valid {
			// Include due_date in INSERT
			insertQuery = `
				INSERT INTO bills (id, tenant_id, unit_id, category, period, amount, late_fee, due_date, status, notes, created_by)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, 'pending', $9, $10)
			`
			args = []interface{}{billID, tenantID, unitID, req.Category, req.Period, req.Amount, lateFee, dueDate.Time, req.Notes, userID}
		} else {
			// Exclude due_date from INSERT (will be NULL)
			insertQuery = `
				INSERT INTO bills (id, tenant_id, unit_id, category, period, amount, late_fee, status, notes, created_by)
				VALUES ($1, $2, $3, $4, $5, $6, $7, 'pending', $8, $9)
			`
			args = []interface{}{billID, tenantID, unitID, req.Category, req.Period, req.Amount, lateFee, req.Notes, userID}
		}
		
		_, err = tx.Exec(insertQuery, args...)
		if err != nil {
			tx.Rollback()
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create bill for unit " + unitID + ": " + err.Error()})
		}
		createdBills = append(createdBills, billID)
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to commit transaction"})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message":      "Bills created successfully",
		"created_count": len(createdBills),
		"bill_ids":     createdBills,
	})
}

// UpdateBill updates a bill
func UpdateBill(c echo.Context) error {
	tenantID := c.Get(string(middleware.CtxTenantID)).(string)
	billID := c.Param("bill_id")

	req := new(models.UpdateBillRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// Check if bill exists and belongs to tenant
	var exists bool
	err := db.DB.Get(&exists, `
		SELECT EXISTS(
			SELECT 1 FROM bills 
			WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL
		)
	`, billID, tenantID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}
	if !exists {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Bill not found"})
	}

	// Build update query dynamically
	updates := []string{}
	args := []interface{}{}
	argIndex := 1

	if req.Category != nil {
		updates = append(updates, "category = $"+strconv.Itoa(argIndex))
		args = append(args, *req.Category)
		argIndex++
	}
	if req.Period != nil {
		updates = append(updates, "period = $"+strconv.Itoa(argIndex))
		args = append(args, *req.Period)
		argIndex++
	}
	if req.Amount != nil {
		updates = append(updates, "amount = $"+strconv.Itoa(argIndex))
		args = append(args, *req.Amount)
		argIndex++
	}
	if req.LateFee != nil {
		updates = append(updates, "late_fee = $"+strconv.Itoa(argIndex))
		args = append(args, *req.LateFee)
		argIndex++
	}
	if req.DueDate != nil {
		if *req.DueDate == "" {
			// Allow clearing due_date by setting to NULL
			updates = append(updates, "due_date = NULL")
		} else {
			dueDate, err := time.Parse("2006-01-02", *req.DueDate)
			if err != nil {
				return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid due_date format. Use YYYY-MM-DD"})
			}
			updates = append(updates, "due_date = $"+strconv.Itoa(argIndex))
			args = append(args, dueDate)
			argIndex++
		}
	}
	if req.Status != nil {
		updates = append(updates, "status = $"+strconv.Itoa(argIndex))
		args = append(args, *req.Status)
		argIndex++
	}
	if req.Notes != nil {
		updates = append(updates, "notes = $"+strconv.Itoa(argIndex))
		args = append(args, *req.Notes)
		argIndex++
	}

	if len(updates) == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "No fields to update"})
	}

	updates = append(updates, "updated_at = NOW()")
	args = append(args, billID, tenantID)

	query := "UPDATE bills SET " + updates[0]
	for i := 1; i < len(updates); i++ {
		query += ", " + updates[i]
	}
	query += " WHERE id = $" + strconv.Itoa(argIndex) + " AND tenant_id = $" + strconv.Itoa(argIndex+1)

	_, err = db.DB.Exec(query, args...)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update bill: " + err.Error()})
	}

	return GetBill(c)
}

// DeleteBill deletes a bill (soft delete)
func DeleteBill(c echo.Context) error {
	tenantID := c.Get(string(middleware.CtxTenantID)).(string)
	billID := c.Param("bill_id")

	// Check if bill exists and belongs to tenant
	var exists bool
	err := db.DB.Get(&exists, `
		SELECT EXISTS(
			SELECT 1 FROM bills 
			WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL
		)
	`, billID, tenantID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}
	if !exists {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Bill not found"})
	}

	// Soft delete
	_, err = db.DB.Exec(`
		UPDATE bills 
		SET deleted_at = NOW(), updated_at = NOW()
		WHERE id = $1 AND tenant_id = $2
	`, billID, tenantID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete bill: " + err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Bill deleted successfully",
	})
}

// ProcessPayment processes a payment for a bill
func ProcessPayment(c echo.Context) error {
	tenantID := c.Get(string(middleware.CtxTenantID)).(string)
	billID := c.Param("bill_id")

	req := new(models.ProcessPaymentRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// Check if bill exists and belongs to tenant
	var billStatus string
	var billAmount float64
	err := db.DB.QueryRow(`
		SELECT status, amount FROM bills 
		WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL
	`, billID, tenantID).Scan(&billStatus, &billAmount)
	if err == sql.ErrNoRows {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Bill not found"})
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}

	if billStatus == "paid" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Bill is already paid"})
	}

	// Process payment
	paymentRef := ""
	if req.PaymentReference != nil {
		paymentRef = *req.PaymentReference
	}

	_, err = db.DB.Exec(`
		UPDATE bills 
		SET status = 'paid', 
		    paid_at = NOW(),
		    payment_method = $1,
		    payment_reference = $2,
		    updated_at = NOW()
		WHERE id = $3 AND tenant_id = $4
	`, req.PaymentMethod, paymentRef, billID, tenantID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to process payment: " + err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Payment processed successfully",
		"id":      billID,
	})
}


