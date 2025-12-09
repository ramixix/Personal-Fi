package core

import (
	"financial_tracker/internal/models"
	"financial_tracker/internal/storage"
	"fmt"
	"strings"
	"time"
)

// Add a transaction to storage
func AddTransaction(t models.Transaction) {
	storage.Transactions = append(storage.Transactions, t)
	storage.NextTransactionID++
}

// Delete a transaction by ID
// inside handleDeleteTransaction i can call this function at the end of the function but
// in orginal code i start from starting index so that is faster i think.
func DeleteTransaction(startingIndex int, id int) bool {
	for i := startingIndex; i < len(storage.Transactions); i++ {
		if storage.Transactions[i].ID == id {
			// Remove transaction from slice
			storage.Transactions = append(storage.Transactions[:i], storage.Transactions[i+1:]...)
			return true
		}
	}
	return false
}

// Find transaction by ID
func FindTransaction(transacId int) *models.Transaction {
	for i := range storage.Transactions {
		if storage.Transactions[i].ID == transacId {
			return &storage.Transactions[i]
		}
	}
	return nil
}

// Get transactions within date range (last N days)
func GetTransactionsByDateRange(days int) []models.Transaction {
	var filtered []models.Transaction
	cutoff_data := time.Now().AddDate(0, 0, -days)

	for _, transaction := range storage.Transactions {
		if transaction.Date.After(cutoff_data) {
			filtered = append(filtered, transaction)
		}
	}
	return filtered
}

// Get transactions by category
func GetTransactionsByCategory(category string) []models.Transaction {
	category = strings.ToLower(category)
	var filtered []models.Transaction

	for _, transac := range storage.Transactions {
		if strings.ToLower(transac.Category) == category {
			filtered = append(filtered, transac)
		}
	}
	return filtered
}

// Get transactions by type (income or expense)
func GetTransactionsByType(transactionType string) []models.Transaction {
	var filtered []models.Transaction

	for _, transac := range storage.Transactions {
		if transac.Type == transactionType {
			filtered = append(filtered, transac)
		}
	}
	return filtered
}

// Get transactions between two dates
func GetTransactionsByCustomRange(start, end time.Time) []models.Transaction {
	var filtered []models.Transaction

	for _, transac := range storage.Transactions {
		// when doing comparision it also compare the time (hour, minute, second) so in that case even the start or end date is equal the transaction date the time will be
		// most probable different. we know that start and end times are set to 0 for all hour, minute, and second so we do the same for transaction.
		transacDate := time.Date(transac.Date.Year(), transac.Date.Month(), transac.Date.Day(), 0, 0, 0, 0, transac.Date.Location())

		if (transacDate.Equal(start) || transacDate.After(start)) && (transacDate.Equal(end) || transacDate.Before(end)) {
			filtered = append(filtered, transac)
		}
	}
	return filtered
}

// Get all unique categories
func GetCategories() []string {
	categories_map := make(map[string]bool)
	var categories []string

	for _, transac := range storage.Transactions {
		if !categories_map[transac.Category] {
			categories_map[transac.Category] = true
			categories = append(categories, transac.Category)
		}
	}
	return categories
}

// Calculate total income and expenses
func CalculateTotals() (totalIncome, totalExpenses float64) {
	for _, transaction := range storage.Transactions {
		if transaction.Type == "income" {
			totalIncome += transaction.Amount
		} else {
			totalExpenses += transaction.Amount
		}
	}
	return //equivalent to : return totalIncome, totalExpenses

	// The function is equivalent to this, which uses explicit return values:
	// func calculateTotals() (float64, float64) {
	// 	var totalIncome float64 = 0.0
	// 	var totalExpenses float64 = 0.0

	// 	// ... (calculation logic is the same)

	// 	// Return the final accumulated values
	// 	return totalIncome, totalExpenses
	// }
}

// Get monthly average income and expneses
func GetMonthlyAverage() (float64, float64) {
	avgIncome := 0.0
	avgExpenses := 0.0
	var oldestDate, newestDate time.Time
	for i, transac := range storage.Transactions {
		if i == 0 {
			oldestDate = transac.Date
			newestDate = transac.Date
		} else {
			if transac.Date.Before(oldestDate) {
				oldestDate = transac.Date
			} else if transac.Date.After(newestDate) {
				newestDate = transac.Date
			}
		}
	}

	years := newestDate.Year() - oldestDate.Year()    // returns int type
	months := newestDate.Month() - oldestDate.Month() // return time.Month type an alias for base type int so we need to convert it to int if we want perform arithmetic operations.
	totalMonths := years*12 + int(months) + 1         // +1 to inlcude the current month

	// redundant security check. so in any way if dates get swap we would get negative number. but since above the dates are correctly set this would be redundant. but keep it if later we change there would be security check.
	if totalMonths < 1 {
		totalMonths = 1
	}

	totalIncome, totalExpenses := CalculateTotals()
	avgIncome = totalIncome / float64(totalMonths)
	avgExpenses = totalExpenses / float64(totalMonths)

	return avgIncome, avgExpenses
}

// Display filtered transactions with summary
func ListFilteredTransactions(transaction_list []models.Transaction, title string) {
	fmt.Printf("%s\n", title)
	fmt.Printf("%s\n", strings.Repeat("=", len(title)))

	if len(transaction_list) == 0 {
		fmt.Println("No transactions found.")
		return
	}

	var total_income, total_expense float64

	for _, transaction := range transaction_list {
		fmt.Printf("ID: %-6d | %15s | %-8s | $%-10.2f | %-20s | %s\n",
			transaction.ID,
			transaction.Date.Format("2006-01-02 15:04"),
			transaction.Type,
			transaction.Amount,
			transaction.Category,
			transaction.Description)

		if transaction.Type == "income" {
			total_income += transaction.Amount
		} else {
			total_expense += transaction.Amount
		}
	}

	fmt.Printf("\n--- Summary ---\n")
	fmt.Printf("Total Income:   $%-.2f\n", total_income)
	fmt.Printf("Total Expenses: $%-.2f\n", total_expense)
	fmt.Printf("Net Amount:     $%-.2f\n", total_income-total_expense)
}
