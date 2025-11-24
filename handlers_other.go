package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func hadleSummary() {
	fmt.Println("Financial Summary")
	fmt.Println("=================")

	if len(transactions) == 0 && len(accounts) == 0 {
		fmt.Println("No data available yet. Start by adding transactions or creating accounts!")
		return
	}

	totalIncome, totalExpenses := calculateTotals()
	netAmount := totalIncome - totalExpenses
	fmt.Println("\n--- Overall Totals ---")
	fmt.Printf("Total Income:    $%.2f\n", totalIncome)
	fmt.Printf("Total Expenses:  $%.2f\n", totalExpenses)
	fmt.Printf("Net Amount:      $%.2f", netAmount)
	if netAmount >= 0 {
		fmt.Println(" ✓")
	} else {
		fmt.Println(" ⚠")
	}

	if len(transactions) > 0 {
		avgIncome, avgExpenses := getMonthlyAverage()
		fmt.Println("\n--- Monthly Averages ---")
		fmt.Printf("Avg Income:      $%.2f\n", avgIncome)
		fmt.Printf("Avg Expenses:    $%.2f\n", avgExpenses)
		fmt.Printf("Avg Net:         $%.2f\n", avgIncome-avgExpenses)
	}

	// Account summary
	if len(accounts) > 0 {
		fmt.Println("\n--- Accounts Summary ---")
		var totalAccountBalance float64
		for _, account := range accounts {
			totalAccountBalance += account.Balance
			fmt.Printf("%-20s $%.2f\n", account.Name+":", account.Balance)
		}
		fmt.Printf("%-20s $%.2f\n", "Total in Accounts:", totalAccountBalance)
	}

	// Transaction count
	fmt.Println("\n--- Transaction Statistics ---")
	fmt.Printf("Total Transactions: %d\n", len(transactions))

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
	if len(transactions) > 0 {
		var oldestDate, newestDate = transactions[0].Date, transactions[0].Date
		for _, transaction := range transactions {
			if transaction.Date.Before(oldestDate) {
				oldestDate = transaction.Date
			}
			if transaction.Date.After(newestDate) {
				newestDate = transaction.Date
			}
		}
		fmt.Printf("Date Range:         %s to %s\n", oldestDate.Format("2006-01-02"), newestDate.Format("2006-01-02"))
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
	fmt.Println("  export        	Export data to CSV files")
	fmt.Println("  import        	Import data from CSV files")
}

// Show version information
func showVersion() {
	fmt.Println("Financial Tracker v0.1.0")
	fmt.Println("A personal finance management tool")
}

func handleExport() {
	if len(os.Args) < 3 {
		fmt.Println("Export Data")
		fmt.Println("===========")
		fmt.Println("Usage: financial-tracker export [type]")
		fmt.Println("Types:")
		fmt.Println("  transactions       Export all transactions to CSV")
		fmt.Println("  accounts           Export all accounts to CSV")
		fmt.Println("  account-history    Export account transaction history to CSV")
		fmt.Println("  all                Export everything to separate CSV files")
		return
	}

	exportType := strings.ToLower(os.Args[2])

	switch exportType {
	case "transactions":
		filename := "transactions_export.csv"
		if err := exportTransactionsToCSV(filename); err != nil {
			fmt.Printf("Error exporting transactions: %v\n", err)
			return
		}
		fmt.Printf("✓ Transactions exported to %s (%d records)\n", filename, len(transactions))

	case "accounts":
		filename := "accounts_export.csv"
		if err := exportAccountsToCSV(filename); err != nil {
			fmt.Printf("Error exporting accounts: %v\n", err)
			return
		}
		fmt.Printf("✓ Accounts exported to %s (%d records)\n", filename, len(accounts))

	case "account-history":
		filename := "account_history_export.csv"
		if err := exportAccountTransactionsToCSV(filename); err != nil {
			fmt.Printf("Error exporting account history: %v\n", err)
			return
		}
		fmt.Printf("✓ Account history exported to %s (%d records)\n", filename, len(accountTransactions))

	case "all":
		fmt.Println("Exporting all data...")

		// Export transactions
		if err := exportTransactionsToCSV("transactions_export.csv"); err != nil {
			fmt.Printf("Error exporting transactions: %v\n", err)
		} else {
			fmt.Printf("✓ Transactions exported (%d records)\n", len(transactions))
		}

		// Export accounts
		if err := exportAccountsToCSV("accounts_export.csv"); err != nil {
			fmt.Printf("Error exporting accounts: %v\n", err)
		} else {
			fmt.Printf("✓ Accounts exported (%d records)\n", len(accounts))
		}

		// Export account history
		if err := exportAccountTransactionsToCSV("account_history_export.csv"); err != nil {
			fmt.Printf("Error exporting account history: %v\n", err)
		} else {
			fmt.Printf("✓ Account history exported (%d records)\n", len(accountTransactions))
		}

		fmt.Println("\n✓ All data exported successfully!")

	default:
		fmt.Printf("Unknown export type: %s\n", exportType)
		fmt.Println("Available types: transactions, accounts, account-history, all")
	}
}

// Handle import command
func handleImport() {
	if len(os.Args) < 3 {
		fmt.Println("Import Data")
		fmt.Println("===========")
		fmt.Println("Usage: financial-tracker import [type] [filename]")
		fmt.Println("Types:")
		fmt.Println("  transactions    Import transactions from CSV")
		fmt.Println("  accounts        Import accounts from CSV")
		fmt.Println("\nExample:")
		fmt.Println("  financial-tracker import transactions my_transactions.csv")
		return
	}

	importType := strings.ToLower(os.Args[2])

	// Check for filename
	if len(os.Args) < 4 {
		fmt.Println("Error: Please specify a filename")
		fmt.Println("Example: financial-tracker import transactions my_transactions.csv")
		return
	}

	filename := os.Args[3]

	// Check if file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		fmt.Printf("Error: File '%s' does not exist\n", filename)
		return
	}

	switch importType {
	case "transactions":
		oldCount := len(transactions)

		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("Import transactions from '%s'?\n", filename)
		if !getConfirmation(reader, "This will add to existing transactions. Continue? (yes/no): ") {
			fmt.Println("Import cancelled.")
			return
		}

		if err := importTransactionsFromCSV(filename); err != nil {
			fmt.Printf("Error importing transactions: %v\n", err)
			return
		}

		newCount := len(transactions) - oldCount
		fmt.Printf("✓ Imported %d transactions from %s\n", newCount, filename)
		fmt.Printf("Total transactions: %d\n", len(transactions))

	case "accounts":
		oldCount := len(accounts)

		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("Import accounts from '%s'?\n", filename)
		if !getConfirmation(reader, "This will add to existing accounts. Continue? (yes/no): ") {
			fmt.Println("Import cancelled.")
			return
		}

		if err := importAccountsFromCSV(filename); err != nil {
			fmt.Printf("Error importing accounts: %v\n", err)
			return
		}

		newCount := len(accounts) - oldCount
		fmt.Printf("✓ Imported %d accounts from %s\n", newCount, filename)
		fmt.Printf("Total accounts: %d\n", len(accounts))

	default:
		fmt.Printf("Unknown import type: %s\n", importType)
		fmt.Println("Available types: transactions, accounts")
	}
}
