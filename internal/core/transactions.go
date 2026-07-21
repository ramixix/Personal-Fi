package core

import (
	"financial_tracker/internal/models"
	"financial_tracker/internal/storage"
	"fmt"
	"strings"
	"time"
)

// ===================================================
// Inserting & Updating & Deleting Operations
//====================================================

// AddTransaction adds a new transaction
func AddTransaction(tx models.Transaction) error {
	_, err := storage.Store.InsertTransaction(tx)
	return err
}

// UpdateTransaction updates values of a given transactions with new values
func UpdateTransaction(tx *models.Transaction) error {
	err := storage.Store.UpdateTransaction(*tx)
	if err != nil {
		return err
	}
	return nil
}

// DeleteTransaction deletes a transaction
func DeleteTransaction(id string) error {
	return storage.Store.DeleteTransaction(id)
}

// ===================================================
// Retriving Transaction Information
// ====================================================
// GetTransactionsLength returns total number of transactions
func GetTransactionsLength(transactionType models.TransactionType) int {
	length, err := storage.Store.GetTransactionsLength(transactionType)
	if err != nil {
		return 0
	}
	return length
}

// FindTransaction finds a transaction by ID
func FindTransaction(id string) *models.Transaction {
	tx, err := storage.Store.GetTransactionByID(id)
	if err != nil {
		return nil
	}
	return tx
}

// GetRecentTransactions(N) returns the last N transactions
func GetRecentTransactions(limit int) []models.Transaction {
	recentTransactions, err := storage.Store.GetRecentTransactions(limit)
	if err != nil {
		return []models.Transaction{}
	}
	return recentTransactions
}

// GetAllTransactions returns all transactions
func GetTransactionBatch(batchsize int, offset int) []models.Transaction {
	transactions, err := storage.Store.GetTransactionsPaginated(batchsize, offset)
	if err != nil {
		return []models.Transaction{}
	}
	return transactions
}

// GetTransactionsByType returns transactions by type
func GetTransactionsByType(txType string) []models.Transaction {
	transactions, err := storage.Store.GetTransactionsAdvanceSearch(models.SearchCriteria{TransactionType: txType})
	if err != nil {
		return []models.Transaction{}
	}
	return transactions
}

// Get transactions by category
func GetTransactionsByCategory(category string) []models.Transaction {
	transactions, err := storage.Store.GetTransactionsAdvanceSearch(models.SearchCriteria{Categories: []string{category}})
	if err != nil {
		return []models.Transaction{}
	}
	return transactions
}

// Get transactions by multiple categories
func GetTransactionsByMultipleCategories(categories []string) []models.Transaction {
	transactions, err := storage.Store.GetTransactionsAdvanceSearch(models.SearchCriteria{Categories: categories})
	if err != nil {
		return []models.Transaction{}
	}
	return transactions
}

// GetTransactionsForLastDays returns transactions in a date range (transactions in last N days)
func GetTransactionsForLastDays(days int) []models.Transaction {
	cutoffDate := time.Now().AddDate(0, 0, -days)
	endDate := time.Now()

	transactions, err := storage.Store.GetTransactionsAdvanceSearch(
		models.SearchCriteria{HasDateRange: true, StartDate: cutoffDate, EndDate: endDate},
	)
	if err != nil {
		return []models.Transaction{}
	}

	return transactions
}

// GetTransactionsForDateRange Get transactions between two custom start and end dataes
func GetTransactionsForDateRange(start, end time.Time) []models.Transaction {
	transactions, err := storage.Store.GetTransactionsAdvanceSearch(
		models.SearchCriteria{HasDateRange: true, StartDate: start, EndDate: end},
	)
	if err != nil {
		return []models.Transaction{}
	}

	return transactions
}

// Search transactions by amount range
func SearchTransactionsByAmountRange(min, max float64) []models.Transaction {
	results, err := storage.Store.GetTransactionsAdvanceSearch(models.SearchCriteria{MinAmount: min, MaxAmount: max})
	if err != nil {
		return []models.Transaction{}
	}
	return results
}

// SearchTransactionsByKeyword finds transactions by searching given keyword in description and in category
func SearchTransactionsByKeyword(keyword string) []models.Transaction {
	if keyword == "" {
		return []models.Transaction{}
	}

	keyword = strings.ToLower(keyword)

	results, err := storage.Store.GetTransactionsAdvanceSearch(models.SearchCriteria{Keyword: keyword})
	if err != nil {
		return []models.Transaction{}
	}
	return results
}

// SearchTransactionsByKeyword returns transactions by filtering based on criteria given(such as keyword, category, type, min/max amount, start/end date)
func GetTransactionsAdvanceSearch(searchCriteria models.SearchCriteria) []models.Transaction {
	results, err := storage.Store.GetTransactionsAdvanceSearch(searchCriteria)
	if err != nil {
		return []models.Transaction{}
	}
	return results
}

// ============================================
// Retriving Analytics and Category Info
// ============================================

// GetCategories returns all unique categories
func GetCategories() []string {
	categories, err := storage.Store.GetCategories()
	if err != nil {
		return []string{}
	}
	return categories
}

// GetCategorySummary return a summary of each category based on currency code
func GetCategorySummary() []models.CategorySummary {
	summaries, err := storage.Store.GetTransactionsCategorySummary()
	if err != nil {
		fmt.Printf("Error loading categories summary: %v\n", err)
		return nil
	}
	return summaries
}

// CalculateTotalsByCurrency calculates total income and expenses grouped by currency code
func CalculateTotalsByCurrency() (map[string]models.CurrencyTotal, error) {
	totals, err := storage.Store.GetCurrencyTotals()
	if err != nil {
		return nil, fmt.Errorf("calculate totals by currency: %w", err)
	}

	return totals, nil
}

// GetMonthlyAverage calculates average monthly income and expenses per currency
func GetMonthlyAverage() (map[string]models.CurrencyTotal, error) {
	monthlyAverage, err := storage.Store.GetAllMonthsAverageByCurrency()
	if err != nil {
		return nil, fmt.Errorf("calculate monthly average by currency: %w", err)
	}
	return monthlyAverage, nil
}

// GetTransactionsMonthlyReports calculates aggregated monthly data
func GetTransactionsMonthlyReports() map[string][]models.MonthlyReport {
	monthlyReports, err := storage.Store.GetTransactionsMonthlyReports()
	if err != nil {
		return nil
	}
	return monthlyReports
}

// GetTransactionsYearlyReports calculates aggregated yearly data
func GetTransactionsYearlyReports() map[string][]models.MonthlyReport {
	yearlyReports, err := storage.Store.GetTransactionsYearlyReports()
	if err != nil {
		return nil
	}
	return yearlyReports
}

// GetTransactionsQuarterlyReport calculates aggregated quarterly data
func GetTransactionsQuarterlyReport(year, quarter int) map[string]models.MonthlyReport {
	quarterlyReport, err := storage.Store.GetTransactionsQuarterlyReport(year, quarter)
	if err != nil {
		return nil
	}
	return quarterlyReport
}

// GetSpecificMonthYearReport gets aggregated data for a specific month and year per currency
func GetSpecificMonthYearReport(year int, month time.Month) map[string]models.MonthlyReport {
	report, err := storage.Store.GetSpecificMonthYearReport(year, month)
	if err != nil {
		fmt.Printf("failed to get specified year/month reports: %v\n", err)
		return nil
	}
	return report
}

// ComparePeriods calculates income/expenses for two specific date ranges per currency
func ComparePeriods(start1, end1, start2, end2 time.Time) map[string]models.ComparisonReport {
	report, err := storage.Store.ComparePeriods(start1, end1, start2, end2)
	if err != nil {
		fmt.Printf("failed to get periods data reports: %v\n", err)
		return nil
	}
	return report
}

// GetYearOverYearComparison gets year-over-year comparison
func GetYearOverYearComparison(year1, year2 int) map[string]models.ComparisonReport {
	start1 := time.Date(year1, 1, 1, 0, 0, 0, 0, time.UTC)
	end1 := time.Date(year1, 12, 31, 23, 59, 59, 0, time.UTC)
	start2 := time.Date(year2, 1, 1, 0, 0, 0, 0, time.UTC)
	end2 := time.Date(year2, 12, 31, 23, 59, 59, 0, time.UTC)

	return ComparePeriods(start1, end1, start2, end2)
}

// GetMonthOverMonthComparison gets month-over-month comparison
func GetMonthOverMonthComparison(year1 int, month1 time.Month, year2 int, month2 time.Month) map[string]models.ComparisonReport {
	start1 := time.Date(year1, month1, 1, 0, 0, 0, 0, time.UTC)
	end1 := time.Date(year1, month1+1, 0, 23, 59, 59, 0, time.UTC) // Last day of month
	start2 := time.Date(year2, month2, 1, 0, 0, 0, 0, time.UTC)
	end2 := time.Date(year2, month2+1, 0, 23, 59, 59, 0, time.UTC)

	return ComparePeriods(start1, end1, start2, end2)
}

// DetectHighSpendingTransactions finds transactions exceeding the average by a given threshold
func DetectHighSpendingTransactions(thresholdMultiplier float64) []models.Transaction {
	transactions, err := storage.Store.DetectHighSpendingTransactions(thresholdMultiplier)
	if err != nil {
		return nil
	}
	return transactions
}
