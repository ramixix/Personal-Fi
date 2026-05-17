package cli

import (
	"financial_tracker/internal/core"
	"financial_tracker/internal/utils"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func handleSummary() {
	fmt.Println("Financial Summary")
	fmt.Println("=================")
	transactions := core.GetAllTransactions()
	accounts := core.GetAllAccounts()
	transactionsCount := len(transactions)
	accountsCount := len(accounts)

	if transactionsCount == 0 && accountsCount == 0 {
		fmt.Println("No data available yet. Start by adding transactions or creating accounts!")
		return
	}

	totalIncome, totalExpenses := core.CalculateTotals()
	netAmount := totalIncome - totalExpenses
	fmt.Println("\n--- Overall Totals ---")
	fmt.Printf("Total Income:    %.2f\n", totalIncome)
	fmt.Printf("Total Expenses:  %.2f\n", totalExpenses)
	fmt.Printf("Net Amount:      %.2f", netAmount)
	if netAmount >= 0 {
		fmt.Println(" ✓")
	} else {
		fmt.Println(" ⚠")
	}

	if transactionsCount > 0 {
		avgIncome, avgExpenses := core.GetMonthlyAverage()
		fmt.Println("\n--- Monthly Averages ---")
		fmt.Printf("Avg Income:      %.2f\n", avgIncome)
		fmt.Printf("Avg Expenses:    %.2f\n", avgExpenses)
		fmt.Printf("Avg Net:         %.2f\n", avgIncome-avgExpenses)
	}

	// Account summary
	if accountsCount > 0 {
		fmt.Println("\n--- Accounts Summary ---")
		var totalAccountBalance float64
		for _, account := range accounts {
			totalAccountBalance += account.Balance
			fmt.Printf("%-20s %s\n", account.Name+":", utils.FormatCurrency(account.Balance, account.CurrencyCode))
		}
		fmt.Printf("%-20s %.2f\n", "Total in Accounts:", totalAccountBalance)
	}

	// Transaction count
	fmt.Println("\n--- Transaction Statistics ---")
	fmt.Printf("Total Transactions: %d\n", transactionsCount)

	incomeCount := 0
	expenseCount := 0
	for _, t := range transactions {
		if t.Type == "income" {
			incomeCount++
		} else {
			expenseCount++
		}
	}
	fmt.Printf("Income Entries:     %d\n", incomeCount)
	fmt.Printf("Expense Entries:    %d\n", expenseCount)

	// Date range
	if transactionsCount > 0 {
		var oldestDate, newestDate = transactions[0].Date, transactions[0].Date
		for _, transaction := range transactions {
			if transaction.Date.Before(oldestDate) {
				oldestDate = transaction.Date
			}
			if transaction.Date.After(newestDate) {
				newestDate = transaction.Date
			}
		}
		fmt.Printf("Date Range: %s to %s\n", oldestDate.Format("2006-01-02"), newestDate.Format("2006-01-02"))
	}
}

// Handle delete command with subcommands
func handleDelete() {
	if len(os.Args) < 3 || os.Args[2] == "help" {
		fmt.Println("Delete")
		fmt.Println("======")
		fmt.Printf("Usage: %s delete [type]\n", filepath.Base(os.Args[0]))
		fmt.Println("Types:")
		fmt.Println("  transaction   Delete a transaction")
		fmt.Println("  account       Delete an account")
		return
	}

	deleteType := strings.ToLower(os.Args[2])
	switch deleteType {
	case "transaction":
		handleDeleteTransaction()
	case "account":
		handleDeleteAccount()
	default:
		fmt.Printf("[Error] Unknown delete type: %s\n", deleteType)
		fmt.Println("Available types: transaction, account")
	}
}

// Show help message
func showHelp() {
	fmt.Printf("\nUsage: %s [command]\n", filepath.Base(os.Args[0]))
	fmt.Println("\nCommands:")
	fmt.Println("  help      		Show this help message")
	fmt.Println("  version   		Show version information")
	fmt.Println("  add       		Add a new transaction")
	fmt.Println("  list [?filter]	List transactions(filters: week, month, year, income, expense, category, custom-range)")
	fmt.Println("  accounts  		Manage accounts (accounts help => to see what commands you can run with accounts)")
	fmt.Println("  delete			Delete transaction or accounts. ''/'delete help' => to show the options.")
	fmt.Println("  edit				Edit an existing transaction")
	fmt.Println("  goals         	Manage goals (create, list, contribute, view)")
	fmt.Println("  search        	Search and filter transactions")
	fmt.Println("  summary       	Show financial summary and statistics")
	fmt.Println("  reports       	Generate reports and analytics")
	fmt.Println("  export        	Export data to CSV files")
	fmt.Println("  import        	Import data from CSV files")
}

// Show version information
func showVersion() {
	fmt.Println("Financial Tracker v0.1.0")
	fmt.Println("A personal finance management tool")
}
