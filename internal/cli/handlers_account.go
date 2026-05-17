package cli

import (
	"bufio"
	"financial_tracker/internal/core"
	"financial_tracker/internal/models"
	"financial_tracker/internal/utils"
	"fmt"
	"os"
	"path/filepath"
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
	currencyCode := utils.GetValidCurrency(reader, "USD")

	newAccount := models.Account{
		ID:           utils.MustGenerateUUID(),
		Name:         accountName,
		Balance:      0.0,
		CurrencyCode: currencyCode,
		Created:      time.Now(),
	}

	core.AddAccount(newAccount)

	fmt.Printf("✓ Account '%s' created successfully! ID: %s\n", accountName, newAccount.ID)
}

// List all accounts
func handleListAccounts() {
	fmt.Println("List Accounts")
	fmt.Println("==============")

	accounts := core.GetAllAccounts()
	if len(accounts) == 0 {
		fmt.Println("\n[Info] No Accounts Found.")
		return
	}

	var totalBalance float64
	for _, account := range accounts {
		fmt.Printf("ID: %-8s | Account Name: %-20s | Balance: %-15s | Created: %s\n", account.ID, account.Name, utils.FormatCurrency(account.Balance, account.CurrencyCode), account.Created.Format("2006-01-02"))
		totalBalance += account.Balance
	}

	fmt.Printf("\n Total Across All Acounts: $%.2f\n", totalBalance)
}

// Add money to an account
func handleAddToAccount() {
	accounts := core.GetAllAccounts()
	if len(accounts) == 0 {
		fmt.Println("No accounts available. Create an account first.")
		return
	}

	fmt.Println("Add Money To Account")
	fmt.Println("====================")

	handleListAccounts()

	reader := bufio.NewReader(os.Stdin)
	id := utils.GetNonEmptyString(reader, "\nPlease Enter Account ID: ")

	account := core.FindAccount(id)
	if account == nil {
		fmt.Println("Account with ID %s not found", id)
		return
	}

	amountToAdd := utils.GetValidAmount(reader)
	fmt.Print("Note (Optional): ")
	noteInput, _ := reader.ReadString('\n')
	note := strings.TrimSpace(noteInput)

	core.AddMoneyToAccount(id, amountToAdd, note)
	fmt.Printf("✓ Added $%.2f to '%s'", amountToAdd, account.Name)
}

// Show account transaction history
func handleAccountHistory() {
	accounts := core.GetAllAccounts()
	if len(accounts) == 0 {
		fmt.Println("[Warning] No Account Available.")
		return
	}

	fmt.Println("Acouunt History")
	fmt.Println("===============")

	fmt.Println("Available Accounts:")

	for _, account := range accounts {
		fmt.Printf("Account ID: %s | Account Name: %s\n", account.ID, account.Name)
	}

	reader := bufio.NewReader(os.Stdin)
	id := utils.GetNonEmptyString(reader, "\nPlease Enter Account ID: ")

	var accountName string
	account := core.FindAccount(id)
	if account == nil {
		fmt.Println("Account with ID %s not found", id)
		return
	}
	accountName = account.Name

	fmt.Printf("\nHistory for Account %s:\n", accountName)
	fmt.Println(strings.Repeat("=", len(accountName)+20))

	accountTransactions := core.GetAccountTransactions(id)
	if len(accountTransactions) > 0 {
		for _, transac := range accountTransactions {
			fmt.Printf("%-12s | +$%-14.2f | %s\n", transac.Date.Format("2006-01-02 15:04"), transac.Amount, transac.Note)
		}
	} else {
		fmt.Println("[Warning] No Transaction Found For This Account.")
	}
}

// Handle delete account command
func handleDeleteAccount() {
	accounts := core.GetAllAccounts()
	if len(accounts) == 0 {
		fmt.Println("[Info] No Accounts to Delete")
		return
	}

	fmt.Println("Delete Account")
	fmt.Println("==============")

	fmt.Printf("\n[Info] Total Account Number: %d.\n", len(accounts))

	handleListAccounts()

	reader := bufio.NewReader(os.Stdin)
	input := utils.GetNonEmptyString(reader, "\nPlease Enter Account ID to Delete (or 'cancel'): ")

	if input == "cancel" {
		fmt.Println("Deletion cancelled.")
		return
	}

	doYouConfirm := utils.GetConfirmation(reader, "[Info] ARE YOU SURE? This will delete all account history and all account transactions. (yes/y or anything else considered as no): ")

	if !doYouConfirm {
		fmt.Println("Deletion cancelled.")
		return
	}

	isDeleted := core.DeleteAccount(input)
	if isDeleted == nil {
		fmt.Printf("✓ Account with ID %s has deleted successfully!\n", input)
		return
	}
	fmt.Printf("Account with ID %s not found!\n", input)
}
