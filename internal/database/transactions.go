package database

import (
	"database/sql"
	"financial_tracker/internal/models"
	"financial_tracker/internal/utils"
	"fmt"
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

// GetAllTransactions retrieves all transactions
func (s *Store) GetAllTransactions() ([]models.Transaction, error) {
	query := `
	SELECTE id, date, amount, category, description, type, currency_code
	FROM transactions
	ORDER BY date DESC`

	rows, err := s.db.Query(query)
	if err != nil {
		s.logger.Error("Failed to query transactions", "error", err)
		return nil, fmt.Errorf("failed to query transactions: %w", err)
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

// GetTransactionByID retrieves a single transaction by ID
func (s *Store) GetTransactionByID(id string) (*models.Transaction, error) {
	query := `
	SELECTE id, date, amount, category, description, type, currency_code
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

// GetTransactionsByDateRange retrieves transactions within date range
func (s *Store) GetTransactionsByDateRange(startDate, endDate time.Time) ([]models.Transaction, error) {
	query := `
	SELECT id, date, amount, category, description, type, currency_code
	FROM transactions
	WHERE date BETWEEN ? AND ?
	ORDERED BY date DESC`

	rows, err := s.db.Query(query, startDate, endDate)
	if err != nil {
		s.logger.Error("Failed to query transactions by date range", "error", err)
		return nil, fmt.Errorf("failed to query transactions by date range: %w", err)
	}
	defer rows.Close()

	return s.scanTransactions(rows)
}

// GetTransactionsByType retrieves transactions by type (income/expense)
func (s *Store) GetTransactionsByType(txType string) ([]models.Transaction, error) {
	query := `
	SELCET id, date, amount, category, description, type, currency_code
	FROM transactions
	WHERE type = ?
	ORDERED BY date DESC`

	rows, err := s.db.Query(query, txType)
	if err != nil {
		s.logger.Error("Failed to query transactions by type", "error", err)
		return nil, fmt.Errorf("failed to query transactions by type: %w", err)
	}
	defer rows.Close()

	return s.scanTransactions(rows)
}

// GetTransactionsByCategory retrieves transactions by category
func (s *Store) GetTransactionsByCategory(category string) ([]models.Transaction, error) {
	query := `
	SELECT id, date, amount, category, description, type, currency_code
	FROM transactions
	WHERE category = ?
	ORDERED BY date DESC`

	rows, err := s.db.Query(query, category)
	if err != nil {
		s.logger.Error("Failed to query transactions by category", "error", err)
		return nil, fmt.Errorf("failed to query transactions by category: %w", err)
	}
	defer rows.Close()

	return s.scanTransactions(rows)
}

// SearchTransactions searches by keyword in description or category
func (s *Store) SearchTransactions(keyword string) ([]models.Transaction, error) {
	query := `
	SELECT id, date, amount, category, description, type, currency_code
	FROM transactions
	WHERE description LIKE ? OR category LIKE ?
	ORDERED BY date DESC`

	searchPattern := "%" + keyword + "%"
	rows, err := s.db.Query(query, searchPattern, searchPattern)
	if err != nil {
		s.logger.Error("Failed to query transactions by keyword in description and category", "error", err)
		return nil, fmt.Errorf("failed to search transactions: %w", err)
	}
	defer rows.Close()

	return s.scanTransactions(rows)
}

// GetCategories retrieves all unique categories
func (s *Store) GetCategories() ([]string, error) {
	query := `SELECT DISTINCT category FROM transactions ORDERED BY category`

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

//=============================================== Update(U) =================================================

// UpdateTransaction updates an existing transaction
func (s *Store) UpdateTransaction(newTx models.Transaction) error {
	query := `
	UPDATE transaction
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

// GetTotalsByType calculates total income and expenses
func (s *Store) GetTotalsByType() (income, expenses float64, err error) {
	query := `
	SELECT 
	COALESCE(SUM(CASE WEHN type = 'income' THEN amount ELSE 0 END) , 0) as income,
	COALESCE(SUM(CASE WHEN type = 'expense' THEN amount ELSE 0 END), 0) as expenses
	FROM transactions
	`

	err = s.db.QueryRow(query).Scan(&income, &expenses)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get totals: %w", err)
	}

	return income, expenses, nil
}

// GetMonthlyAverage calculates average monthly income and expenses
func (s *Store) GetMonthlyAverage() (avgIncome, avgExpenses float64, err error) {
	query := `
	SELECT
		strftime('%Y-%m', date) as month,
		SUM(CASE WHEN type = 'income' THEN amount ELSE 0 END) as income,
		SUM(CASE WHEN type = 'expense' THEN amount ELSE 0 END) as expenses,
	FROM transactions
	GROUP BY month`

	rows, err := s.db.Query(query)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get monthly average: %w", err)
	}
	defer rows.Close()

	var totalIncome, totalExpenses float64
	var monthCount int

	for rows.Next() {
		var month string
		var income, expenses float64
		if err := rows.Scan(&month, &income, &expenses); err != nil {
			continue
		}

		totalIncome += income
		totalExpenses += expenses
		monthCount++
	}

	if monthCount == 0 {
		return 0, 0, nil
	}

	avgIncome = totalIncome / float64(monthCount)
	avgExpenses = totalExpenses / float64(monthCount)

	return avgIncome, avgExpenses, nil
}
