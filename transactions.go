package main

import (
	"fmt"
	"strings"
	"time"
)

// Add a transaction to storage
func addTransaction(t Transaction) {
	transactions = append(transactions, t)
	nextTransactionID++
}

// Delete a transaction by ID
// inside handleDeleteTransaction i can call this function at the end of the function but
// in orginal code i start from starting index so that is faster i think.
func deleteTransaction(startingIndex int, id int) bool {
	for i := startingIndex; i < len(transactions); i++ {
		if transactions[i].ID == id {
			// Remove transaction from slice
			transactions = append(transactions[:i], transactions[i+1:]...)
			return true
		}
	}
	return false
}

// Find transaction by ID
func findTransaction(transacId int) *Transaction {
	for i := range transactions {
		if transactions[i].ID == transacId {
			return &transactions[i]
		}
	}
	return nil
}

// Get transactions within date range (last N days)
func getTransactionsByDateRange(days int) []Transaction {
	var filtered []Transaction
	cutoff_data := time.Now().AddDate(0, 0, -days)

	for _, transaction := range transactions {
		if transaction.Date.After(cutoff_data) {
			filtered = append(filtered, transaction)
		}
	}
	return filtered
}

// Get transactions by category
func getTransactionsByCategory(category string) []Transaction {
	category = strings.ToLower(category)
	var filtered []Transaction

	for _, transac := range transactions {
		if strings.ToLower(transac.Category) == category {
			filtered = append(filtered, transac)
		}
	}
	return filtered
}

// Get transactions by type (income or expense)
func getTransactionsByType(transactionType string) []Transaction {
	var filtered []Transaction

	for _, transac := range transactions {
		if transac.Type == transactionType {
			filtered = append(filtered, transac)
		}
	}
	return filtered
}

// Get transactions between two dates
func getTransactionsByCustomRange(start, end time.Time) []Transaction {
	var filtered []Transaction

	for _, transac := range transactions {
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
func getCategories() []string {
	categories_map := make(map[string]bool)
	var categories []string

	for _, transac := range transactions {
		if !categories_map[transac.Category] {
			categories_map[transac.Category] = true
			categories = append(categories, transac.Category)
		}
	}
	return categories
}

// Calculate total income and expenses
func calculateTotals() (totalIncome, totalExpenses float64) {
	for _, transaction := range transactions {
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
func getMonthlyAverage() (avgIncome float64, avgExpenses float64) {

	var oldestDate, newestDate time.Time
	for i, transac := range transactions {
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

	totalIncome, totalExpenses := calculateTotals()
	avgIncome = totalIncome / float64(totalMonths)
	avgExpenses = totalExpenses / float64(totalMonths)

	return avgIncome, avgExpenses
}

// Display filtered transactions with summary
func listFilteredTransactions(transaction_list []Transaction, title string) {
	fmt.Printf("%s\n", title)
	fmt.Printf("%s\n", strings.Repeat("=", len(title)))

	if len(transaction_list) == 0 {
		fmt.Println("No transactions found.")
		return
	}

	var total_income, total_expense float64

	for _, transaction := range transaction_list {
		fmt.Printf("ID: %d | %s | %s | $%.2f | %s | %s\n",
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
	fmt.Printf("Total Income:  $%.2f\n", total_income)
	fmt.Printf("Total Expenses: $%.2f\n", total_expense)
	fmt.Printf("Net Amount:    $%.2f\n", total_income-total_expense)
}
