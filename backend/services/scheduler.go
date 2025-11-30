package services

import (
	"database/sql"
	"log"
	"time"
	"rukunos-backend/db"
)

// StartScheduler starts all scheduled background jobs
func StartScheduler() {
	log.Println("Starting scheduler...")
	
	// Start late fee calculation job (runs daily at 00:00)
	go runDailyJob(calculateLateFees, time.Hour*24, "00:00")
	
	// Start bill status update job (runs every hour)
	go runHourlyJob(updateBillStatus)
	
	log.Println("Scheduler started")
}

// runDailyJob runs a job daily at a specific time
func runDailyJob(job func(), interval time.Duration, timeStr string) {
	// Parse time string (HH:MM)
	now := time.Now()
	jobTime, err := time.Parse("15:04", timeStr)
	if err != nil {
		log.Printf("Error parsing time %s: %v", timeStr, err)
		return
	}
	
	// Calculate next run time
	nextRun := time.Date(now.Year(), now.Month(), now.Day(), jobTime.Hour(), jobTime.Minute(), 0, 0, now.Location())
	if nextRun.Before(now) {
		nextRun = nextRun.Add(24 * time.Hour)
	}
	
	// Wait until next run time
	time.Sleep(time.Until(nextRun))
	
	// Run job immediately
	job()
	
	// Then run every interval
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	
	for range ticker.C {
		job()
	}
}

// runHourlyJob runs a job every hour
func runHourlyJob(job func()) {
	// Run immediately
	job()
	
	// Then run every hour
	ticker := time.NewTicker(time.Hour)
	defer ticker.Stop()
	
	for range ticker.C {
		job()
	}
}

// calculateLateFees calculates and updates late fees for overdue bills
func calculateLateFees() {
	log.Println("Running late fee calculation job...")
	
	// Get all pending/overdue bills that have due_date and are not paid
	query := `
		SELECT 
			b.id, b.tenant_id, b.amount, b.late_fee, b.due_date, b.status,
			bt.late_fee_type, bt.late_fee as template_late_fee, 
			bt.late_fee_percentage, bt.late_fee_max
		FROM bills b
		LEFT JOIN billing_templates bt ON b.category = bt.name AND b.tenant_id = bt.tenant_id AND bt.deleted_at IS NULL
		WHERE b.status IN ('pending', 'overdue')
		AND b.due_date IS NOT NULL
		AND b.due_date < CURRENT_DATE
		AND b.deleted_at IS NULL
	`
	
	rows, err := db.DB.Query(query)
	if err != nil {
		log.Printf("Error querying bills for late fee calculation: %v", err)
		return
	}
	defer rows.Close()
	
	updatedCount := 0
	for rows.Next() {
		var billID, tenantID, status, lateFeeType string
		var amount, currentLateFee, templateLateFee sql.NullFloat64
		var lateFeePercentage, lateFeeMax sql.NullFloat64
		var dueDate time.Time
		
		err := rows.Scan(&billID, &tenantID, &amount, &currentLateFee, &dueDate, &status,
			&lateFeeType, &templateLateFee, &lateFeePercentage, &lateFeeMax)
		if err != nil {
			log.Printf("Error scanning bill: %v", err)
			continue
		}
		
		// Calculate days overdue
		daysOverdue := int(time.Since(dueDate).Hours() / 24)
		if daysOverdue <= 0 {
			continue
		}
		
		// Calculate new late fee
		var newLateFee float64
		
		if lateFeeType == "percentage" && lateFeePercentage.Valid && amount.Valid {
			// Percentage-based: percentage of amount per day
			dailyLateFee := (amount.Float64 * lateFeePercentage.Float64) / 100.0
			newLateFee = dailyLateFee * float64(daysOverdue)
		} else {
			// Fixed amount per day
			dailyLateFee := 0.0
			if templateLateFee.Valid {
				dailyLateFee = templateLateFee.Float64
			} else if currentLateFee.Valid {
				// Use existing late fee as daily rate (if it was set manually)
				dailyLateFee = currentLateFee.Float64 / float64(daysOverdue)
			} else {
				// Default: 5000 per day
				dailyLateFee = 5000.0
			}
			newLateFee = dailyLateFee * float64(daysOverdue)
		}
		
		// Apply max late fee if configured
		if lateFeeMax.Valid && newLateFee > lateFeeMax.Float64 {
			newLateFee = lateFeeMax.Float64
		}
		
		// Update bill
		updateQuery := `
			UPDATE bills 
			SET late_fee = $1, 
			    status = CASE WHEN status = 'pending' THEN 'overdue' ELSE status END,
			    updated_at = NOW()
			WHERE id = $2 AND tenant_id = $3
		`
		
		_, err = db.DB.Exec(updateQuery, newLateFee, billID, tenantID)
		if err != nil {
			log.Printf("Error updating late fee for bill %s: %v", billID, err)
			continue
		}
		
		updatedCount++
	}
	
	log.Printf("Late fee calculation completed. Updated %d bills.", updatedCount)
}

// updateBillStatus updates bill status from pending to overdue
func updateBillStatus() {
	log.Println("Running bill status update job...")
	
	// Update bills that are pending and past due date to overdue
	query := `
		UPDATE bills 
		SET status = 'overdue', updated_at = NOW()
		WHERE status = 'pending'
		AND due_date IS NOT NULL
		AND due_date < CURRENT_DATE
		AND deleted_at IS NULL
	`
	
	result, err := db.DB.Exec(query)
	if err != nil {
		log.Printf("Error updating bill status: %v", err)
		return
	}
	
	rowsAffected, _ := result.RowsAffected()
	log.Printf("Bill status update completed. Updated %d bills to overdue.", rowsAffected)
}

