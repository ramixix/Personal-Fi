package main

import (
	"fmt"
	"os"
	"path/filepath"
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
	case "accounts":
		handleAccounts()
	case "delete":
		handleDelete()
	case "edit":
		handleEditTransaction()
	default:
		fmt.Printf("Unknown command: %s\n", command)
		showHelp()
	}

	err = saveData()
	if err != nil {
		fmt.Printf("[Warning] saveData function could not Save the data, error given: %v\n", err)
	}
}

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

	deleteType := os.Args[2]
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

func showHelp() {
	fmt.Printf("\nUsage: %s [command]\n", filepath.Base(os.Args[0]))
	fmt.Println("\nCommands:")
	fmt.Println("  help      		Show this help message")
	fmt.Println("  version   		Show version information")
	fmt.Println("  add       		Add a new transaction")
	fmt.Println("  list [?filter]	List transactions(filters: week, month, year, income, expense, category, custom-range)")
	fmt.Println("  accounts  		Manage accounts (accounts help => to see what commands you can run with accounts)")
	fmt.Println("  delete			Delete transaction or accounts. Enter help to show the options.")
	fmt.Println("  edit             Edit an existing transaction")
}

func showVersion() {
	fmt.Println("Financial Tracker v0.1.0")
}
