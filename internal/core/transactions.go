package core

import (
	"financial_tracker/internal/models"
	"financial_tracker/internal/storage"
	"strings"
	"time"
)

// / AddTransaction adds a new transaction
func AddTransaction(tx models.Transaction) error {
	_, err := storage.Store.InsertTransaction(tx)
	return err
}

// DeleteTransaction deletes a transaction
func DeleteTransaction(id string) error {
	return storage.Store.DeleteTransaction(id)
}

// FindTransaction finds a transaction by ID
func FindTransaction(id string) *models.Transaction {
	tx, err := storage.Store.GetTransactionByID(id)
	if err != nil {
		return nil
	}
	return tx
}

// GetAllTransactions returns all transactions
func GetAllTransactions() []models.Transaction {
	transactions, err := storage.Store.GetAllTransactions()
	if err != nil {
		return []models.Transaction{}
	}
	return transactions
}

// GetTransactionsByDateRange returns transactions in a date range (transactions in last N days)
func GetTransactionsByDateRange(days int) []models.Transaction {
	cutoffDate := time.Now().AddDate(0, 0, -days)
	endDate := time.Now()

	transactions, err := storage.Store.GetTransactionsByDateRange(cutoffDate, endDate)
	if err != nil {
		return []models.Transaction{}
	}

	return transactions
}

// Get transactions by category
func GetTransactionsByCategory(category string) []models.Transaction {
	category = strings.ToLower(category)
	transactions, err := storage.Store.GetTransactionsByCategory(category)
	if err != nil {
		return []models.Transaction{}
	}
	return transactions
}

// GetTransactionsByType returns transactions by type
func GetTransactionsByType(txType string) []models.Transaction {
	transactions, err := storage.Store.GetTransactionsByType(txType)
	if err != nil {
		return []models.Transaction{}
	}
	return transactions
}

// GetTransactionsByCustomRange Get transactions between two custom start and end dataes
func GetTransactionsByCustomRange(start, end time.Time) []models.Transaction {
	transactions, err := storage.Store.GetTransactionsByDateRange(start, end)
	if err != nil {
		return []models.Transaction{}
	}

	return transactions
}

// GetCategories returns all unique categories
func GetCategories() []string {
	categories, err := storage.Store.GetCategories()
	if err != nil {
		return []string{}
	}
	return categories
}

// CalculateTotals calculates total income and expenses
func CalculateTotals() (totalIncome, totalExpenses float64) {
	totalIncome, totalExpenses, _ = storage.Store.GetTotalsByType()
	return totalIncome, totalExpenses
}

// GetMonthlyAverage calculates monthly averages
func GetMonthlyAverage() (avgIncome, avgExpenses float64) {
	avgIncome, avgExpenses, _ = storage.Store.GetMonthlyAverage()
	return avgIncome, avgExpenses
}

// GetRecentTransactions returns the last N transactions
func GetRecentTransactions(limit int) []models.Transaction {
	if limit <= 0 {
		return []models.Transaction{}
	}

	transactions, err := storage.Store.GetAllTransactions()
	if err != nil {
		return []models.Transaction{}
	}

	if limit > len(transactions) {
		limit = len(transactions)
	}

	return transactions[:limit]
}
