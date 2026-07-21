package cli

import (
	"financial_tracker/internal/core"
	"financial_tracker/internal/models"
	"financial_tracker/internal/utils"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func handleSummary() {
	fmt.Println("Financial Summary")
	fmt.Println("=================")
	transactionsCount := core.GetTransactionsLength("")
	accountsCount := core.GetAccountsLength()

	if transactionsCount <= 0 && accountsCount <= 0 {
		fmt.Println("No data available yet. Start by adding transactions or creating accounts!")
		return
	}

	fmt.Println("\n--- Transaction Stats ---")
	fmt.Printf("Total number of Income transactions:\t%d\n", core.GetTransactionsLength(models.Income))
	fmt.Printf("Total number of expense transactions:\t%d\n", core.GetTransactionsLength(models.Expense))
	fmt.Printf("Total number of transactions:\t%d\n", transactionsCount)

	currency_totals, err := core.CalculateTotalsByCurrency()
	if err == nil {
		fmt.Println("\n--- Overall Totals ---")
		for currency, totals := range currency_totals {
			netAmount := totals.Income - totals.Expenses
			fmt.Printf("Currency:\t%s\n", currency)
			fmt.Printf("\tTotal Income:\t%s\n", utils.FormatCurrency(totals.Income, currency))
			fmt.Printf("\tTotal Expenses:\t%s\n", utils.FormatCurrency(totals.Expenses, currency))
			fmt.Printf("\tNet Amount:\t%s", utils.FormatCurrency(netAmount, currency))
			if netAmount >= 0 {
				fmt.Println(" ✓\n\n")
			} else {
				fmt.Println(" ⚠\n\n")
			}
		}
	}

	// Account summary
	if accountsCount > 0 {
		fmt.Println("\n--- Accounts Summary ---")
		totalAccountBalanceByCurrency := core.GetTotalAccountsBalanceByCurrency()
		for currency, total := range totalAccountBalanceByCurrency {
			fmt.Printf("%-20s %s\n", utils.FormatCurrency(total, currency), currency)
		}
	}

	currency_totals, err = core.GetMonthlyAverage()
	if err == nil {
		fmt.Println("\n--- All Months Average ---")
		for currency, totals := range currency_totals {
			netAmount := totals.Income - totals.Expenses
			fmt.Printf("Currency:\t%s\n", currency)
			fmt.Printf("\tAvg Income:\t%s\n", utils.FormatCurrency(totals.Income, currency))
			fmt.Printf("\tAvg Expenses:\t%s\n", utils.FormatCurrency(totals.Expenses, currency))
			fmt.Printf("\tAvg Amount:\t%s", utils.FormatCurrency(netAmount, currency))
			if netAmount >= 0 {
				fmt.Println(" ✓\n\n")
			} else {
				fmt.Println(" ⚠\n\n")
			}
		}
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
