package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("=== Personal Financial Tracker ===")

	err := loadData()
	if err != nil {
		fmt.Printf("[Warning] loadData function could not load the data, error given: %v\n", err)
	}

	// Check if any command line arguments were provided
	if len(os.Args) < 2 {
		showHelp()
		return
	}

	command := os.Args[1]
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
	case "summary":
		hadleSummary()
	case "export":
		handleExport()
	case "import":
		handleImport()
	default:
		fmt.Printf("Unknown command: %s\n", command)
		showHelp()
	}

	err = saveData()
	if err != nil {
		fmt.Printf("[Warning] saveData function could not Save the data, error given: %v\n", err)
	}
}
