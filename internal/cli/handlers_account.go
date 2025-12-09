package cli

import (
	"bufio"
	"financial_tracker/internal/core"
	"financial_tracker/internal/models"
	"financial_tracker/internal/storage"
	"financial_tracker/internal/utils"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// Handle accounts command with subcommands
func handleAccounts() {
	if len(os.Args) < 3 {
		fmt.Println("Account Management")
		fmt.Println("==================")
		fmt.Printf("Usage: %s accounts [command]\n", filepath.Base(os.Args[0]))
		fmt.Println("Commands:")
		fmt.Println("  create    Create a new account")
		fmt.Println("  list      List all accounts")
		fmt.Println("  add       Add money to an account")
		fmt.Println("  history   Show account transaction history")
		fmt.Println("  delete    Delete account along all transaction belong to it")
		return
	}

	subCommand := strings.ToLower(os.Args[2])
	switch subCommand {
	case "create":
		handleCreateAccount()
	case "list":
		handleListAccounts()
	case "add":
		handleAddToAccount()
	case "history":
		handleAccountHistory()
	case "delete":
		handleDeleteAccount()
	default:
		fmt.Println("[Warning] Unknown Accounts Command, Please Pay Attention To The Usage Guide.")
	}
}

// Create a new account
func handleCreateAccount() {
	fmt.Println("Create New Account")
	fmt.Println("==================")

	reader := bufio.NewReader(os.Stdin)
	accountName := utils.GetNonEmptyString(reader, "Account Name: ")

	newAccount := models.Account{
		ID:      storage.NextAccountID,
		Name:    accountName,
		Balance: 0.0,
		Created: time.Now(),
	}

	storage.Accounts = append(storage.Accounts, newAccount)
	storage.NextAccountID += 1

	fmt.Printf("✓ Account '%s' created successfully! ID: %d\n", accountName, newAccount.ID)
}

// List all accounts
func handleListAccounts() {
	fmt.Println("List Accounts")
	fmt.Println("==============")

	if len(storage.Accounts) == 0 {
		fmt.Println("\n[Info] No Accounts Found.")
		return
	}

	var totalBalance float64
	for _, account := range storage.Accounts {
		fmt.Printf("ID: %-8d | Account Name: %-20s | Balance: $%-15.2f | Created: %s\n", account.ID, account.Name, account.Balance, account.Created.Format("2006-01-02"))
		totalBalance += account.Balance
	}

	fmt.Printf("\n Total Across All Acounts: $%.2f\n", totalBalance)
}

// Add money to an account
func handleAddToAccount() {
	if len(storage.Accounts) == 0 {
		fmt.Println("No accounts available. Create an account first.")
		return
	}

	fmt.Println("Add Money To Account")
	fmt.Println("====================")

	handleListAccounts()

	reader := bufio.NewReader(os.Stdin)
	id, err := utils.GetIntInput(reader, "\nEnter Account ID: ")

	if err != nil {
		fmt.Println("[Warning] Invalid Account ID!")
		return
	}

	if id > len(storage.Accounts) || id < 0 {
		fmt.Printf("[Warning] Please Pay Attention to Range of Accounts Available %d-%d\n", 0, len(storage.Accounts)-1)
		return
	}

	account := core.FindAccount(id)
	if account == nil {
		fmt.Println("Account Not found")
		return
	}

	amountToAdd := utils.GetValidAmount(reader)
	fmt.Print("Note (Optional): ")
	noteInput, _ := reader.ReadString('\n')
	note := strings.TrimSpace(noteInput)

	core.AddMoneyToAccount(account, amountToAdd, note)
	fmt.Printf("✓ Added $%.2f to '%s'. New balance: $%.2f\n", amountToAdd, account.Name, account.Balance)
}

// Show account transaction history
func handleAccountHistory() {
	if len(storage.Accounts) == 0 {
		fmt.Println("[Warning] No Account Available.")
		return
	}

	fmt.Println("Acouunt History")
	fmt.Println("===============")

	fmt.Println("Available Accounts:")

	for _, account := range storage.Accounts {
		fmt.Printf("Account ID: %d | Account Name: %s\n", account.ID, account.Name)
	}

	reader := bufio.NewReader(os.Stdin)
	id, err := utils.GetIntInput(reader, "Enter Account ID: ")

	if err != nil || id < 0 || id > len(storage.Accounts) {
		fmt.Println("[Error] Invalid ID. Please Enter An Integer In The Available Range!!!")
		return
	}

	var accountName string
	account := core.FindAccount(id)
	if account == nil {
		fmt.Println("No Account Found!")
		return
	}
	accountName = account.Name

	fmt.Printf("\nHistory for Account %s:\n", accountName)
	fmt.Println(strings.Repeat("=", len(accountName)+20))

	found := false
	for _, transac := range storage.AccountTransactions {
		if transac.AccountID == id {
			fmt.Printf("%-12s | +$%-14.2f | %s\n", transac.Date.Format("2006-01-02 15:04"), transac.Amount, transac.Note)
			found = true
		}
	}

	if !found {
		fmt.Println("[Warning] No Transaction Found For This Account.")
	}
}

// Handle delete account command
func handleDeleteAccount() {
	if len(storage.Accounts) == 0 {
		fmt.Println("[Info] No Accounts to Delete")
		return
	}

	fmt.Println("Delete Account")
	fmt.Println("==============")

	fmt.Printf("\n[Info] Total Account Number: %d.\n", len(storage.Accounts))

	handleListAccounts()

	fmt.Print("\nPlease Enter The Account ID to Delete (or 'cancel'): ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "cancel" {
		fmt.Println("Deletion cancelled.")
		return
	}
	accountId, err := strconv.Atoi(input)
	if err != nil {
		fmt.Println("[Warning] Not a Valid Transaction ID.")
		return
	}

	doYouConfirm := utils.GetConfirmation(reader, "[Info] Are You SURE? This will delete all account hitory and all account transactions. (yes/y or anything else considered as no): ")

	if !doYouConfirm {
		fmt.Println("Deletion cancelled.")
		return
	}

	isDeleted := core.DeleteAccount(accountId)
	if isDeleted {
		fmt.Printf("✓ Account ID %d deleted successfully!\n", accountId)
		return
	}
	fmt.Printf("Account ID %d not found!\n", accountId)
}
