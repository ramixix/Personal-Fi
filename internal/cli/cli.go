package cli

import (
	"financial_tracker/internal/storage"
	"fmt"
	"os"
)

// Run starts the CLI application
func Run() {
	fmt.Println("=== Personal Financial Tracker ===")

	// Load existing data
	err := storage.LoadData()
	if err != nil {
		fmt.Printf("Warning: Could not load data: %v\n", err)
	}

	// Check if any command line arguments were provided
	if len(os.Args) < 2 {
		showHelp()
		storage.SaveData()
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

	// Save data after any command
	err = storage.SaveData()
	if err != nil {
		fmt.Printf("Error saving data: %v\n", err)
	}
}
