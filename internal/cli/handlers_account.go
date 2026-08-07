package cli

import (
	"bufio"
	"financial_tracker/internal/core"
	"financial_tracker/internal/models"
	"financial_tracker/internal/utils"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// handleAccounts handles accounts command with subcommands
func handleAccounts() {
	if len(os.Args) < 3 {
		fmt.Println("Account Management")
		fmt.Println("==================")
		fmt.Printf("Usage: %s accounts [command]\n", filepath.Base(os.Args[0]))
		fmt.Println("Commands:")
		fmt.Println("  create    Create a new account")
		fmt.Println("  list      List last 100 accounts")
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
		handleListAccounts(recent100)
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

// handleCreateAccount creates a new account by receiving necessary information from user
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

// handleListAccounts lists recent accounts
func handleListAccounts(numberofAccountsToShow int) {
	fmt.Println("List Accounts")
	fmt.Println("==============")

	accountsToShow := core.GetRecentAccounts(numberofAccountsToShow)

	for _, account := range accountsToShow {
		fmt.Printf("ID: %-8s | Account Name: %-20s | Balance: %-15s | Created: %s\n",
			account.ID,
			account.Name,
			utils.FormatCurrency(account.Balance, account.CurrencyCode),
			account.Created.Format("2006-01-02"))
	}

	balances := core.GetTotalAccountsBalanceByCurrency()

	fmt.Println("\nTotal Across All Accounts:")
	fmt.Println("--------------------------")
	for currency, total := range balances {
		fmt.Printf("%s: %.2f\n", currency, total)
	}
}

// GetAccountsNumberToShow asks users for a specific number N to list last N Accounts.
func GetAccountsNumberToShow(reader *bufio.Reader, defaultValue int) int {
	accountsToShow := defaultValue
	totalAccounts := core.GetAccountsLength()
	if defaultValue > totalAccounts {
		accountsToShow = totalAccounts
	}
InputLoop:
	for {
		input, _ := reader.ReadString('\n')
		input = strings.ToLower(strings.TrimSpace(input))

		switch input {
		case "":
			fmt.Printf("\nDisplaying %d recent accounts:\n", accountsToShow)
			break InputLoop
		case "all":
			fmt.Println("\nDisplaying all accounts:")
			accountsToShow = totalAccounts
			break InputLoop
		default:
			number, err := strconv.Atoi(input)
			if err != nil || number <= 0 || number >= totalAccounts {
				fmt.Println("[Warning] Not a valid number, try again.")
				continue
			}
			fmt.Printf("\nDisplaying last %d accounts:\n", number)
			accountsToShow = number
			break InputLoop
		}
	}
	return accountsToShow
}

// handleAddToAccount adds money to an account
func handleAddToAccount() {
	accountsCount := core.GetAccountsLength()
	if accountsCount <= 0 {
		fmt.Println("No accounts available. Create an account first.")
		return
	}

	fmt.Println("Add Money To Account")
	fmt.Println("====================")
	fmt.Printf("\n[Info] Total Account Number: %d.\n", accountsCount)
	fmt.Printf("\nPress Enter to show the latest %d accounts (if possible by default).\nTo see a different number, enter it.\nType 'all' to display all accounts.", recent100)
	reader := bufio.NewReader(os.Stdin)

	accountsToShow := GetAccountsNumberToShow(reader, recent100)
	handleListAccounts(accountsToShow)

	id := utils.GetNonEmptyString(reader, "\nPlease Enter Account ID: ")

	account := core.FindAccount(id)
	if account == nil {
		fmt.Printf("Account with ID %s not found", id)
		return
	}

	amountToAdd := utils.GetValidAmount(reader)
	fmt.Print("Note (Optional): ")
	noteInput, _ := reader.ReadString('\n')
	note := strings.TrimSpace(noteInput)

	core.AddMoneyToAccount(id, amountToAdd, note)
	fmt.Printf("✓ Added %.2f to '%s'", amountToAdd, account.Name)
}

// handleAccountHistory shows specific account history
func handleAccountHistory() {
	accountsCount := core.GetAccountsLength()
	if accountsCount <= 0 {
		fmt.Println("[Warning] No Account Available.")
		return
	}

	fmt.Println("Acouunt History")
	fmt.Println("===============")
	fmt.Printf("\n[Info] Total Account Number: %d.\n", accountsCount)
	fmt.Printf("\nPress Enter to show the latest %d accounts (if possible by default).\nTo see a different number, enter it.\nType 'all' to display all accounts.", recent100)
	reader := bufio.NewReader(os.Stdin)

	accountsToShow := GetAccountsNumberToShow(reader, recent100)
	handleListAccounts(accountsToShow)

	id := utils.GetNonEmptyString(reader, "\nPlease Enter Account ID: ")

	var accountName string
	account := core.FindAccount(id)
	if account == nil {
		fmt.Printf("Account with ID %s not found", id)
		return
	}
	accountName = account.Name

	fmt.Printf("\nHistory for Account %s:\n", accountName)
	fmt.Println(strings.Repeat("=", len(accountName)+20))

	accountTransactions := core.GetOneAccountTransactions(id)
	if len(accountTransactions) > 0 {
		for _, transac := range accountTransactions {
			fmt.Printf("%-12s | +%-14s | %s\n",
				transac.Date.Format("2006-01-02 15:04"),
				utils.FormatCurrency(transac.Amount, account.CurrencyCode),
				transac.Note)
		}
	} else {
		fmt.Println("[Warning] No Transaction Found For This Account.")
	}
}

// handleDeleteAccount deletes users selected account
func handleDeleteAccount() {
	accountsCount := core.GetAccountsLength()
	if accountsCount <= 0 {
		fmt.Println("[Info] No Accounts to Delete")
		return
	}

	fmt.Println("Delete Account")
	fmt.Println("==============")
	fmt.Printf("\n[Info] Total Account Number: %d.\n", accountsCount)
	fmt.Printf("\nPress Enter to show the latest %d accounts (if possible by default).\nTo see a different number, enter it.\nType 'all' to display all accounts.", recent100)
	reader := bufio.NewReader(os.Stdin)

	accountsToShow := GetAccountsNumberToShow(reader, recent100)
	handleListAccounts(accountsToShow)

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
