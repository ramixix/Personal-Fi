package cli

import (
	"bufio"
	"financial_tracker/internal/core"
	"financial_tracker/internal/utils"
	"fmt"
	"os"
	"strings"
)

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
		if err := core.ExportTransactionsToCSV(filename); err != nil {
			fmt.Printf("Error exporting transactions: %v\n", err)
			return
		}
		fmt.Printf("✓ Transactions are exported to %s.\n", filename)

	case "accounts":
		filename := "accounts_export.csv"
		if err := core.ExportAccountsToCSV(filename); err != nil {
			fmt.Printf("Error exporting accounts: %v\n", err)
			return
		}
		fmt.Printf("✓ Accounts are exported to %s.\n", filename)

	case "account-history":
		filename := "account_history_export.csv"
		if err := core.ExportAccountTransactionsToCSV(filename); err != nil {
			fmt.Printf("Error exporting account history: %v\n", err)
			return
		}
		fmt.Printf("✓ Account history(Account Transactions) are exported to %s\n", filename)

	case "all":
		fmt.Println("Exporting all data...")

		// Export transactions
		if err := core.ExportTransactionsToCSV("transactions_export.csv"); err != nil {
			fmt.Printf("Error exporting transactions: %v\n", err)
		} else {
			fmt.Printf("✓ Transactions are exported.\n")
		}

		// Export accounts
		if err := core.ExportAccountsToCSV("accounts_export.csv"); err != nil {
			fmt.Printf("Error exporting accounts: %v\n", err)
		} else {
			fmt.Printf("✓ Accounts are exported.\n")
		}

		// Export account history
		if err := core.ExportAccountTransactionsToCSV("account_history_export.csv"); err != nil {
			fmt.Printf("Error exporting account history: %v\n", err)
		} else {
			fmt.Printf("✓ Account history(Account Transactions) are exported.\n")
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
		oldCount := core.GetAccountsTransactionsLength()

		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("Import transactions from '%s'?\n", filename)
		if !utils.GetConfirmation(reader, "This will add to existing transactions. Continue? (yes/no): ") {
			fmt.Println("Import cancelled.")
			return
		}

		if err := core.ImportTransactionsFromCSV(filename); err != nil {
			fmt.Printf("Error importing transactions: %v\n", err)
			return
		}

		newCount := core.GetAccountsTransactionsLength()
		fmt.Printf("✓ Imported new %d transactions from %s\n", newCount-oldCount, filename)
		fmt.Printf("Total transactions: %d\n", newCount)

	case "accounts":
		oldCount := core.GetAccountsLength()

		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("Import accounts from '%s'?\n", filename)
		if !utils.GetConfirmation(reader, "This will add to existing accounts. Continue? (yes/no): ") {
			fmt.Println("Import cancelled.")
			return
		}

		if err := core.ImportAccountsFromCSV(filename); err != nil {
			fmt.Printf("Error importing accounts: %v\n", err)
			return
		}

		newCount := core.GetAccountsLength()
		fmt.Printf("✓ Imported new %d accounts from %s\n", newCount-oldCount, filename)
		fmt.Printf("Total accounts: %d\n", newCount)

	default:
		fmt.Printf("Unknown import type: %s\n", importType)
		fmt.Println("Available types: transactions, accounts")
	}
}
