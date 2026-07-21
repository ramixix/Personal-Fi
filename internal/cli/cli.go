package cli

import (
	"fmt"
	"os"
)

const recent100 = 100

// Run starts the CLI application
func Run() {
	fmt.Println("=== Personal Financial Tracker ===")

	// Check if any command line arguments were provided
	if len(os.Args) < 2 {
		showHelp()
		return
	}

	command := os.Args[1]

	// Basic command routing
	switch command {
	case "help":
		showHelp()
	case "version":
		showVersion()
	case "add":
		handleAddTransaction()
	case "list":
		handleListTransactions()
	case "edit":
		handleEditTransaction()
	case "delete":
		handleDelete()
	case "accounts":
		handleAccounts()
	case "goals":
		handleGoals()
	case "search":
		handleSearch()
	case "reports":
		handleReports()
	case "summary":
		handleSummary()
	case "export":
		handleExport()
	case "import":
		handleImport()
	default:
		fmt.Printf("Unknown command: %s\n", command)
		showHelp()
	}
}
