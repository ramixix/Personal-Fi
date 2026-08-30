package database

import (
	"database/sql"
	"financial_tracker/internal/models"
	"financial_tracker/internal/utils"
	"fmt"
	"strings"
	"time"
)

//=============================================== Create(C) =================================================

// InsertTransaction creates a new transaction
func (s *Store) InsertTransaction(transaction models.Transaction) (string, error) {
	// Generate UUID if not provided
	if transaction.ID == "" {
		transaction.ID = utils.MustGenerateUUID()
	}

	// Set default currency if not provided
	if transaction.CurrencyCode == "" {
		transaction.CurrencyCode = "USD"
	}

	// Set date if not provided
	if transaction.Date.IsZero() {
		transaction.Date = time.Now()
	}

	query := `
	INSERT INTO transactions(id, date, amount, category, description, type, currency_code)
	VALUES(?, ?, ?, ?, ?, ?, ?)`

	_, err := s.db.Exec(query, transaction.ID, transaction.Date, transaction.Amount, transaction.Category, transaction.Description, transaction.Type, transaction.CurrencyCode)

	if err != nil {
		s.logger.Error("Failed to insert transaction", "error", err, "id", transaction.ID)
		return "", fmt.Errorf("failed to insert transaction: %w", err)
	}

	s.logger.Info("Transaction created", "id", transaction.ID, "amount", transaction.Amount, "type", transaction.Type)
	return transaction.ID, nil
}

//=============================================== READ(R) =================================================

// GetTransactionsLength returns the total number of transactions efficiently
func (s *Store) GetTransactionsLength(transactionType models.TransactionType) (int, error) {
	query := `SELECT COUNT(*) FROM transactions`
	args := []any{}

	if transactionType != "" && (transactionType == models.Expense || transactionType == models.Income) {
		query += ` WHERE type = ?`
		args = append(args, transactionType)
	}

	var count int
	err := s.db.QueryRow(query, args...).Scan(&count)
	if err != nil {
		s.logger.Error("Failed to return transactions Length", "error", err)
		return 0, fmt.Errorf("failed to get transaction Length: %w", err)
	}
	return count, nil
}

// GetTransactionsPaginated retrieves a specific "page" of transactions
func (s *Store) GetTransactionsPaginated(limit, offset int) ([]models.Transaction, error) {
	query := `
    SELECT id, date, amount, category, description, type, currency_code
    FROM transactions
    ORDER BY date DESC
    LIMIT ? OFFSET ?`

	rows, err := s.db.Query(query, limit, offset)
	if err != nil {
		s.logger.Error("Failed to query paginated transactions", "error", err)
		return nil, fmt.Errorf("failed to get paginated transactions: %w", err)
	}
	defer rows.Close()

	return s.scanTransactions(rows)
}

// scanTransactions is a helper to scan multiple transaction rows
func (s *Store) scanTransactions(rows *sql.Rows) ([]models.Transaction, error) {
	var transactions []models.Transaction
	for rows.Next() {
		var tx models.Transaction
		err := rows.Scan(&tx.ID, &tx.Date, &tx.Amount, &tx.Category, &tx.Description, &tx.Type, &tx.CurrencyCode)

		if err != nil {
			s.logger.Error("Failed to scan transaction", "error", err)
			continue
		}
		transactions = append(transactions, tx)
	}
	return transactions, nil
}

// GetRecentTransactions fetches the N most recent transactions
func (s *Store) GetRecentTransactions(limit int) ([]models.Transaction, error) {
	query := `
    SELECT id, date, amount, category, description, type, currency_code
    FROM transactions
    ORDER BY date DESC, created_at DESC
    LIMIT ?`

	rows, err := s.db.Query(query, limit)
	if err != nil {
		s.logger.Error("Failed to query recent N transactions", "error", err)
		return nil, fmt.Errorf("failed to query recent N transactions: %w", err)
	}
	defer rows.Close()

	return s.scanTransactions(rows)
}

// GetTransactionByID retrieves a single transaction by ID
func (s *Store) GetTransactionByID(id string) (*models.Transaction, error) {
	query := `
	SELECT id, date, amount, category, description, type, currency_code
	FROM transactions
	WHERE id = ?`

	var tx models.Transaction
	err := s.db.QueryRow(query, id).Scan(&tx.ID, &tx.Date, &tx.Amount, &tx.Category, &tx.Description, &tx.Type, &tx.CurrencyCode)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("transaction not found: %s", id)
	}
	if err != nil {
		s.logger.Error("Failed to get transaction", "error", err, "id", id)
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}

	return &tx, nil
}

// // GetTransactionsByDateRange retrieves transactions within date range
// func (s *Store) GetTransactionsByDateRange(startDate, endDate time.Time) ([]models.Transaction, error) {
// 	query := `
// 	SELECT id, date, amount, category, description, type, currency_code
// 	FROM transactions
// 	WHERE date BETWEEN ? AND ?
// 	ORDER BY date DESC`

// 	rows, err := s.db.Query(query, startDate, endDate)
// 	if err != nil {
// 		s.logger.Error("Failed to query transactions by date range", "error", err)
// 		return nil, fmt.Errorf("failed to query transactions by date range: %w", err)
// 	}
// 	defer rows.Close()

// 	return s.scanTransactions(rows)
// }

// // GetTransactionsByType retrieves transactions by type (income/expense)
// func (s *Store) GetTransactionsByType(txType string) ([]models.Transaction, error) {
// 	query := `
// 	SELECT id, date, amount, category, description, type, currency_code
// 	FROM transactions
// 	WHERE type = ?
// 	ORDER BY date DESC`

// 	rows, err := s.db.Query(query, txType)
// 	if err != nil {
// 		s.logger.Error("Failed to query transactions by type", "error", err)
// 		return nil, fmt.Errorf("failed to query transactions by type: %w", err)
// 	}
// 	defer rows.Close()

// 	return s.scanTransactions(rows)
// }

// // GetTransactionsByCategory retrieves transactions by category
// func (s *Store) GetTransactionsByCategory(category string) ([]models.Transaction, error) {
// 	query := `
// 	SELECT id, date, amount, category, description, type, currency_code
// 	FROM transactions
// 	WHERE category = ?
// 	ORDER BY date DESC`

// 	rows, err := s.db.Query(query, category)
// 	if err != nil {
// 		s.logger.Error("Failed to query transactions by category", "error", err)
// 		return nil, fmt.Errorf("failed to query transactions by category: %w", err)
// 	}
// 	defer rows.Close()

// 	return s.scanTransactions(rows)
// }

// GetTransactionsAdvanceSearch returns transactions by filtering based on criteria given(such as keyword, category, type, min/max amount, start/end date)
func (s *Store) GetTransactionsAdvanceSearch(search models.SearchCriteria) ([]models.Transaction, error) {
	query := "SELECT id, date, amount, category, description, type, currency_code FROM transactions WHERE 1=1"
	var args []interface{}

	if search.Keyword != "" {
		query += " AND (description LIKE ? OR category LIKE ?)"
		searchPattern := "%" + search.Keyword + "%"
		args = append(args, searchPattern, searchPattern)
	}

	if len(search.Categories) > 0 {
		placeholders := make([]string, len(search.Categories))
		for i, cat := range search.Categories {
			placeholders[i] = "?"
			args = append(args, cat)
		}
		query += fmt.Sprintf(" AND category IN (%s)", strings.Join(placeholders, ","))
	}

	if search.TransactionType != "" {
		query += " AND type = ?"
		args = append(args, search.TransactionType)
	}

	if search.MinAmount > 0 {
		query += " AND amount >= ?"
		args = append(args, search.MinAmount)
	}
	if search.MaxAmount > 0 && search.MaxAmount >= search.MinAmount {
		query += " AND amount <= ?"
		args = append(args, search.MaxAmount)
	}

	if search.HasDateRange {
		query += " AND date >= ? AND date <= ?"
		args = append(args, search.StartDate, search.EndDate)
	}

	query += " ORDER BY date DESC"

	rows, err := s.db.Query(query, args...)
	if err != nil {
		s.logger.Error("Failed to query transactions by given criteria", "error", err)
		return nil, fmt.Errorf("failed to query transactions by given criteria: %w", err)
	}
	defer rows.Close()

	return s.scanTransactions(rows)
}

// // SearchTransactions searches by keyword in description or category
// func (s *Store) SearchTransactions(keyword string) ([]models.Transaction, error) {
// 	query := `
// 	SELECT id, date, amount, category, description, type, currency_code
// 	FROM transactions
// 	WHERE description LIKE ? OR category LIKE ?
// 	ORDER BY date DESC`

// 	searchPattern := "%" + keyword + "%"
// 	rows, err := s.db.Query(query, searchPattern, searchPattern)
// 	if err != nil {
// 		s.logger.Error("Failed to query transactions by keyword in description and category", "error", err)
// 		return nil, fmt.Errorf("failed to search transactions: %w", err)
// 	}
// 	defer rows.Close()

// 	return s.scanTransactions(rows)
// }

// GetCategories retrieves all unique categories
func (s *Store) GetCategories() ([]string, error) {
	query := `SELECT DISTINCT category FROM transactions ORDER BY category`

	rows, err := s.db.Query(query)
	if err != nil {
		s.logger.Error("Failed to query all unique categories", "error", err)
		return nil, fmt.Errorf("failed to query all distinct categories: %w", err)
	}
	defer rows.Close()

	var categories []string
	for rows.Next() {
		var category string
		err := rows.Scan(&category)
		if err != nil {
			continue
		}
		categories = append(categories, category)
	}

	return categories, nil
}

// GetCategorySummaries calculates stats for all categories based on currecy and category name
func (s *Store) GetTransactionsCategorySummary() ([]models.CategorySummary, error) {
	query := `
    SELECT 
        category,
        currency_code,
        COUNT(*) as tx_count,
        COALESCE(SUM(CASE WHEN type = 'income' THEN amount ELSE 0 END), 0) as total_income,
        COALESCE(SUM(CASE WHEN type = 'expense' THEN amount ELSE 0 END), 0) as total_expense
    FROM transactions
    GROUP BY category, currency_code
    ORDER BY total_income DESC, total_expense DESC
    `

	rows, err := s.db.Query(query)
	if err != nil {
		s.logger.Error("Failed to query category summaries", "error", err)
		return nil, fmt.Errorf("failed to get category summaries: %w", err)
	}
	defer rows.Close()

	var summaries []models.CategorySummary
	for rows.Next() {
		var cs models.CategorySummary
		err := rows.Scan(&cs.Category, &cs.CurrencyCode, &cs.Count, &cs.TotalIncome, &cs.TotalExpense)
		if err != nil {
			continue
		}
		summaries = append(summaries, cs)
	}

	return summaries, nil
}

//=============================================== Update(U) =================================================

// UpdateTransaction updates an existing transaction
func (s *Store) UpdateTransaction(newTx models.Transaction) error {
	query := `
	UPDATE transactions
	SET date = ?, amount = ?, category = ?, description = ?, type = ?, currency_code = ?
	WHERE id = ?`

	result, err := s.db.Exec(query, newTx.Date, newTx.Amount, newTx.Category, newTx.Description, newTx.Type, newTx.CurrencyCode, newTx.ID)
	if err != nil {
		s.logger.Error("Failed to update transaction", "error", err, "id", newTx.ID)
		return fmt.Errorf("failed to update transaction: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("transaction not found: %s", newTx.ID)
	}

	s.logger.Info("Transaction updated", "id", newTx.ID)
	return nil
}

//=============================================== DELET(D) =================================================

// DeleteTransaction deletes a transaction
func (s *Store) DeleteTransaction(id string) error {
	query := `
	DELETE FROM transactions WHERE id = ?`

	result, err := s.db.Exec(query, id)
	if err != nil {
		s.logger.Error("Failed to delete transaction", "error", err, "id", id)
		return fmt.Errorf("failed to delete transaction: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("transaction not found: %s", id)
	}

	s.logger.Info("Transaction deleted", "id", id)
	return nil
}

// =============================================== ANALYTICS =================================================

// GetCurrencyTotals calculates total income and expenses grouped by currency
func (s *Store) GetCurrencyTotals() (map[string]models.CurrencyTotal, error) {
	query := `
    SELECT 
        currency_code,
        COALESCE(SUM(CASE WHEN type = 'income' THEN amount ELSE 0 END), 0) as income,
        COALESCE(SUM(CASE WHEN type = 'expense' THEN amount ELSE 0 END), 0) as expenses
    FROM transactions
    GROUP BY currency_code
    `

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get totals: %w", err)
	}
	defer rows.Close()

	totals := make(map[string]models.CurrencyTotal)

	for rows.Next() {
		var currency string
		var ct models.CurrencyTotal

		if err := rows.Scan(&currency, &ct.Income, &ct.Expenses); err != nil {
			s.logger.Error("Failed to scan row", "error", err)
			continue
		}
		totals[currency] = ct
	}

	return totals, nil
}

// GetAllMonthsAverageByCurrency calculates average monthly income and expenses per currency
func (s *Store) GetAllMonthsAverageByCurrency() (map[string]models.CurrencyTotal, error) {
	query := `
    SELECT
        currency_code,
        COALESCE(AVG(monthly_income), 0) as avg_income,
        COALESCE(AVG(monthly_expenses), 0) as avg_expenses
    FROM (
        SELECT
            currency_code,
            strftime('%Y-%m', date) as month,
            SUM(CASE WHEN type = 'income' THEN amount ELSE 0 END) as monthly_income,
            SUM(CASE WHEN type = 'expense' THEN amount ELSE 0 END) as monthly_expenses
        FROM transactions
        GROUP BY currency_code, month
    )
    GROUP BY currency_code
    `

	rows, err := s.db.Query(query)
	if err != nil {
		s.logger.Error("Failed to query monthly averages", "error", err)
		return nil, fmt.Errorf("failed to get monthly averages: %w", err)
	}
	defer rows.Close()

	averages := make(map[string]models.CurrencyTotal)

	for rows.Next() {
		var currency string
		var ct models.CurrencyTotal

		if err := rows.Scan(&currency, &ct.Income, &ct.Expenses); err != nil {
			s.logger.Error("Failed to scan monthly average row", "error", err)
			continue
		}
		averages[currency] = ct
	}

	return averages, nil
}

// GetSpecificMonthYearReport gets aggregated data for a specific month and year per currency
func (s *Store) GetSpecificMonthYearReport(year int, month time.Month) (map[string]models.MonthlyReport, error) {
	query := `
		SELECT 
			currency_code,
			COUNT(*) as tx_count,
			COALESCE(SUM(CASE WHEN type = 'income' THEN amount ELSE 0 END), 0) as income,
			COALESCE(SUM(CASE WHEN type = 'expense' THEN amount ELSE 0 END), 0) as expenses
		FROM transactions
		WHERE CAST(strftime('%Y', date) AS INTEGER) = ? 
		  AND CAST(strftime('%m', date) AS INTEGER) = ?
		GROUP BY currency_code
	`

	// We cast time.Month to an int (e.g., time.June becomes 6)
	rows, err := s.db.Query(query, year, int(month))
	if err != nil {
		return nil, fmt.Errorf("failed to query specific month report: %w", err)
	}
	defer rows.Close()

	reportsMap := make(map[string]models.MonthlyReport)

	for rows.Next() {
		// Initialize the report with the known year and month
		report := models.MonthlyReport{
			Year:  year,
			Month: month,
		}

		if err := rows.Scan(&report.CurrencyCode, &report.TxCount, &report.Income, &report.Expenses); err != nil {
			s.logger.Error("Failed to scan specific month report", "error", err)
			continue
		}

		report.Net = report.Income - report.Expenses
		reportsMap[report.CurrencyCode] = report
	}

	return reportsMap, nil
}

// GetTransactionsMonthlyReports calculates aggregated monthly data
func (s *Store) GetTransactionsMonthlyReports() (map[string][]models.MonthlyReport, error) {
	query := `
		SELECT 
			currency_code,
			CAST(strftime('%Y', date) AS INTEGER) as year,
			CAST(strftime('%m', date) AS INTEGER) as month,
			COUNT(*) as tx_count,
			COALESCE(SUM(CASE WHEN type = 'income' THEN amount ELSE 0 END), 0) as income,
			COALESCE(SUM(CASE WHEN type = 'expense' THEN amount ELSE 0 END), 0) as expenses
		FROM transactions
		GROUP BY currency_code, year, month
		ORDER BY year DESC, month DESC
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to generate monthly reports: %w", err)
	}
	defer rows.Close()

	reportsMap := make(map[string][]models.MonthlyReport)
	for rows.Next() {
		var r models.MonthlyReport
		if err := rows.Scan(&r.CurrencyCode, &r.Year, &r.Month, &r.TxCount, &r.Income, &r.Expenses); err != nil {
			continue
		}
		r.Net = r.Income - r.Expenses
		reportsMap[r.CurrencyCode] = append(reportsMap[r.CurrencyCode], r)
	}

	return reportsMap, nil
}

// GetTransactionsYearlyReports calculates aggregated yearly data
func (s *Store) GetTransactionsYearlyReports() (map[string][]models.MonthlyReport, error) {
	query := `
		SELECT 
			currency_code,
			CAST(strftime('%Y', date) AS INTEGER) as year,
			COUNT(*) as tx_count,
			COALESCE(SUM(CASE WHEN type = 'income' THEN amount ELSE 0 END), 0) as income,
			COALESCE(SUM(CASE WHEN type = 'expense' THEN amount ELSE 0 END), 0) as expenses
		FROM transactions
		GROUP BY currency_code, year
		ORDER BY year DESC
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to generate yearly reports: %w", err)
	}
	defer rows.Close()

	reportsMap := make(map[string][]models.MonthlyReport)
	for rows.Next() {
		var r models.MonthlyReport
		if err := rows.Scan(&r.CurrencyCode, &r.Year, &r.TxCount, &r.Income, &r.Expenses); err != nil {
			continue
		}
		r.Net = r.Income - r.Expenses
		reportsMap[r.CurrencyCode] = append(reportsMap[r.CurrencyCode], r)
	}

	return reportsMap, nil
}

// ComparePeriods calculates income/expenses for two specific date ranges per currency
func (s *Store) ComparePeriods(start1, end1, start2, end2 time.Time) (map[string]models.ComparisonReport, error) {
	query := `
		SELECT 
			currency_code,
			-- Period 1 Sums
			COALESCE(SUM(CASE WHEN type = 'income' AND date >= ? AND date <= ? THEN amount ELSE 0 END), 0) as p1_income,
			COALESCE(SUM(CASE WHEN type = 'expense' AND date >= ? AND date <= ? THEN amount ELSE 0 END), 0) as p1_expense,
			-- Period 2 Sums
			COALESCE(SUM(CASE WHEN type = 'income' AND date >= ? AND date <= ? THEN amount ELSE 0 END), 0) as p2_income,
			COALESCE(SUM(CASE WHEN type = 'expense' AND date >= ? AND date <= ? THEN amount ELSE 0 END), 0) as p2_expense
		FROM transactions
		WHERE (date >= ? AND date <= ?) OR (date >= ? AND date <= ?)
		GROUP BY currency_code
	`

	// We pass the dates in multiple times to satisfy the WHERE and CASE clauses
	rows, err := s.db.Query(query,
		start1, end1, start1, end1, // For P1 CASE statements
		start2, end2, start2, end2, // For P2 CASE statements
		start1, end1, start2, end2, // For WHERE clause
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get periods data reports: %w", err)
	}
	defer rows.Close()

	reportsMap := make(map[string]models.ComparisonReport)
	for rows.Next() {
		var r models.ComparisonReport
		if err := rows.Scan(&r.CurrencyCode, &r.Period1Income, &r.Period1Expenses, &r.Period2Income, &r.Period2Expenses); err != nil {
			continue
		}
		r.IncomeChange = r.Period2Income - r.Period1Income
		r.ExpenseChange = r.Period2Expenses - r.Period1Expenses

		// Calculate percentage changes
		if r.Period1Income > 0 {
			r.IncomePercent = (r.IncomeChange / r.Period1Income) * 100
		}
		if r.Period1Expenses > 0 {
			r.ExpensePercent = (r.ExpenseChange / r.Period1Expenses) * 100
		}
		reportsMap[r.CurrencyCode] = r
	}
	return reportsMap, nil
}

// DetectHighSpendingTransactions finds transactions exceeding the average by a given threshold
func (s *Store) DetectHighSpendingTransactions(thresholdMultiplier float64) ([]models.Transaction, error) {
	query := `
		WITH CurrencyAverages AS (
			SELECT currency_code, AVG(amount) as avg_amount
			FROM transactions
			WHERE type = 'expense'
			GROUP BY currency_code
		)
		SELECT t.id, t.date, t.amount, t.category, t.description, t.type, t.currency_code
		FROM transactions t
		JOIN CurrencyAverages a ON t.currency_code = a.currency_code
		WHERE t.type = 'expense' AND t.amount > (a.avg_amount * ?)
		ORDER BY t.amount DESC
	`

	rows, err := s.db.Query(query, thresholdMultiplier)
	if err != nil {
		return nil, fmt.Errorf("failed to detect anomalies: %w", err)
	}
	defer rows.Close()

	return s.scanTransactions(rows) // Reuse your existing scanner!
}

// GetTransactionsQuarterlyReport gets aggregated data for a specific quarter per currency
func (s *Store) GetTransactionsQuarterlyReport(year int, quarter int) (map[string]models.MonthlyReport, error) {
	if quarter < 1 || quarter > 4 {
		return nil, fmt.Errorf("invalid quarter: %d", quarter)
	}

	// Keep your excellent Go logic for finding the boundaries
	startMonth := (quarter-1)*3 + 1
	endMonth := quarter * 3

	query := `
		SELECT 
			currency_code,
			COUNT(*) as tx_count,
			COALESCE(SUM(CASE WHEN type = 'income' THEN amount ELSE 0 END), 0) as income,
			COALESCE(SUM(CASE WHEN type = 'expense' THEN amount ELSE 0 END), 0) as expenses
		FROM transactions
		WHERE CAST(strftime('%Y', date) AS INTEGER) = ? 
		  AND CAST(strftime('%m', date) AS INTEGER) >= ?
		  AND CAST(strftime('%m', date) AS INTEGER) <= ?
		GROUP BY currency_code
	`

	rows, err := s.db.Query(query, year, startMonth, endMonth)
	if err != nil {
		return nil, fmt.Errorf("failed to query quarterly report: %w", err)
	}
	defer rows.Close()

	reportsMap := make(map[string]models.MonthlyReport)

	for rows.Next() {
		report := models.MonthlyReport{
			Year:  year,
			Month: time.Month(startMonth), // Represents the start of the quarter
		}

		if err := rows.Scan(&report.CurrencyCode, &report.TxCount, &report.Income, &report.Expenses); err != nil {
			s.logger.Error("Failed to scan quarterly report", "error", err)
			continue
		}

		report.Net = report.Income - report.Expenses
		reportsMap[report.CurrencyCode] = report
	}

	return reportsMap, nil
}
