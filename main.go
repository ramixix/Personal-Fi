package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func get_valid_transaction_type(reader *bufio.Reader) string {
	for {
		fmt.Print("Type (Income/Expense): ")
		input, _ := reader.ReadString('\n')
		transaction_type := strings.ToLower(strings.TrimSpace(input))

		if transaction_type == "income" || transaction_type == "expense" {
			return transaction_type
		}
		fmt.Println("[Warning] Invalid type! Pleae Enter 'income' or 'expense'")
	}
}

func get_valid_amount(reader *bufio.Reader) float64 {
	for {
		fmt.Print("Amount: $")
		amount_input, _ := reader.ReadString('\n')
		transcation_amount, err := strconv.ParseFloat(strings.TrimSpace(amount_input), 64)

		if err != nil {
			fmt.Println("Invalid amount, please enter a valid number (float/integer)")
			continue
		}

		if transcation_amount <= 0 {
			fmt.Println("Amount must be greater than 0")
			continue
		}

		return transcation_amount
	}
}

func get_non_empty_string(reader *bufio.Reader, prompt string) string {
	for {
		fmt.Print(prompt)
		input, _ := reader.ReadString('\n')
		text := strings.TrimSpace(input)

		if text != "" {
			return text
		}
		fmt.Println("[Warning] This Field Can Not Be Empty!")
	}
}

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

// Handle adding a new transaction
func handleAddTransaction() {
	fmt.Println("Add Transaction")
	fmt.Println("================")

	reader := bufio.NewReader(os.Stdin)

	// Get validated inputs from users
	transactionType := get_valid_transaction_type(reader)
	amount := get_valid_amount(reader)
	category := get_non_empty_string(reader, "Category: ")
	description := get_non_empty_string(reader, "Description: ")

	newTransaction := Transaction{
		ID:          nextTransactionID,
		Date:        time.Now(),
		Amount:      amount,
		Category:    category,
		Description: description,
		Type:        transactionType,
	}

	addTransaction(newTransaction)
	fmt.Printf("Added transaction: %+v\n", newTransaction)
	fmt.Printf("Total transactions: %d\n", len(transactions))
}

// Handle listing all transactions
func handleListTransactions() {
	fmt.Println("List Transactions")
	fmt.Println("=================")

	if len(os.Args) > 2 {
		filter := strings.ToLower(os.Args[2])
		switch filter {
		case "week":
			listFilteredTransactions(getTransactionsByDateRange(7), "Last 7 Days (Last Week)")
		case "month":
			listFilteredTransactions(getTransactionsByDateRange(30), "Last 30 Days (Last Month)")
		case "year":
			listFilteredTransactions(getTransactionsByDateRange(365), "Last 365 Days (Last Year)")
		case "income":
			listFilteredTransactions(getTransactionsByType("income"), "All Incomes")
		case "expense":
			listFilteredTransactions(getTransactionsByType("expense"), "All Expenses")
		case "category":
			handleCategoryFilter()
		case "categories":
			showCategories()
		case "custom":
			handleCustomRange()
		default:
			fmt.Printf("[Warning] Unknow Filter: %s\n", filter)
			fmt.Println("Available filters: week, month, year, income, expense, category, categories, custom(define date range to show)")
			return
		}
	} else {
		listFilteredTransactions(transactions, "All Transactions")
	}
}

func handleCategoryFilter() {
	if len(transactions) == 0 {
		fmt.Println("No transactions available.")
		return
	}

	fmt.Println("Available Categories:")
	categories := getCategories()

	for i, cat := range categories {
		fmt.Printf("  %d. %s\n", i+1, cat)
	}

	reader := bufio.NewReader(os.Stdin)
	category := get_non_empty_string(reader, "Please Enter The Category Name You Want to Filtere:")
	filtered := getTransactionsByCategory(category)
	listFilteredTransactions(filtered, fmt.Sprintf("Transactions in category: %s", category))

}

func showCategories() {
	fmt.Println("All Categroies And Amount In Total")
	fmt.Println("==================================")

	categoryCount := make(map[string]int)
	categoryTotal := make(map[string]float64)
	categoryIncome := make(map[string]float64)
	categoryExpense := make(map[string]float64)

	for _, transac := range transactions {
		categoryCount[transac.Category]++
		if transac.Type == "income" {
			categoryIncome[transac.Category] += transac.Amount
			categoryTotal[transac.Category] += transac.Amount
		} else {
			categoryExpense[transac.Category] -= transac.Amount
			categoryTotal[transac.Category] -= transac.Amount
		}
	}

	for category, amount := range categoryTotal {
		fmt.Printf("%s: %d Transactions | Total Income: %.2f | Total Expenses: %2.f | Total Amount: $%.2f\n", category, categoryCount[category], categoryIncome[category], categoryExpense[category], amount)
	}
}

func parseDate(date_input string) (time.Time, error) {
	layout := "2006-01-02"
	date, err := time.Parse(layout, date_input)
	return date, err
}

func handleCustomRange() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("\nCustom Range Filter")
	fmt.Println("=====================")
	fmt.Println("Date Format : YYYY-MM-DD Year-Month-Day (e.g 2024-01-23)")

	fmt.Print("Start Date: ")
	start_input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("[Error] Could not read the entered date (handleCustomeRange)")
		return
	}
	start_input = strings.TrimSpace(start_input)
	start_date, err := parseDate(start_input)
	if err != nil {
		fmt.Println("[Error] Not a date in specified format (Formant : YYYY-MM-DD)")
		return
	}

	fmt.Print("End Date: ")
	end_input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("[Error] Could not read the entered date (handleCustomeRange)")
		return
	}
	end_input = strings.TrimSpace(end_input)
	end_date, err := parseDate(end_input)
	if err != nil {
		fmt.Println("[Error] Not a date in specified format (Formant : YYYY-MM-DD)")
		return
	}

	if end_date.Before(start_date) {
		fmt.Println("[Error] End date must be before start date, there is a logical error")
		return
	}

	filtered := getTransactionsByCustomRange(start_date, end_date)
	title := fmt.Sprintf("Transaction From %s to %s", start_date.Format("2006-01-02"), end_date.Format("2006-01-02"))
	listFilteredTransactions(filtered, title)
}

func handleAccounts() {
	if len(os.Args) < 3 {
		fmt.Println("Account Management")
		fmt.Println("==================")
		fmt.Printf("Usage: %s accounts [command]\n", os.Args[0])
		fmt.Println("Commands:")
		fmt.Println("  create    Create a new account")
		fmt.Println("  list      List all accounts")
		fmt.Println("  add       Add money to an account")
		fmt.Println("  history   Show account transaction history")
		fmt.Println("  delete    Delete account along all transaction belong to it")
		return
	}

	subCommand := os.Args[2]
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

func handleCreateAccount() {
	fmt.Println("Create New Account")
	fmt.Println("==================")

	reader := bufio.NewReader(os.Stdin)
	accountName := get_non_empty_string(reader, "Account Name: ")

	newAccount := Account{
		ID:      nextAccountID,
		Name:    accountName,
		Balance: 0.0,
		Created: time.Now(),
	}

	accounts = append(accounts, newAccount)
	nextAccountID += 1

	fmt.Printf("✓ Account '%s' created successfully! ID: %d\n", accountName, newAccount.ID)
}

func handleListAccounts() {
	fmt.Println("List Accounts")
	fmt.Println("==============")

	if len(accounts) == 0 {
		fmt.Println("\n[Info] No Accounts Found.")
		return
	}

	var totalBalance float64
	for _, account := range accounts {
		fmt.Printf("ID: %d | Account Name: %s | Balance: $%.2f | Created: %s\n", account.ID, account.Name, account.Balance, account.Created.Format("2006-01-02"))
		totalBalance += account.Balance
	}

	fmt.Printf("\n Total Across All Acounts: $%.2f\n", totalBalance)
}

func handleAddToAccount() {
	if len(accounts) == 0 {
		fmt.Println("No accounts available. Create an account first.")
		return
	}

	fmt.Println("Add Money To Account")
	fmt.Println("====================")

	handleListAccounts()

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("\nEnter Account ID: ")
	id_input, _ := reader.ReadString('\n')
	id, err := strconv.Atoi(strings.TrimSpace(id_input))

	if err != nil {
		fmt.Println("[Warning] Invalid Account ID!")
		return
	}

	if id > len(accounts) || id < 0 {
		fmt.Printf("[Warning] Please Pay Attention to Range of Accounts Available %d-%d\n", 0, len(accounts)-1)
		return
	}

	account := findAccount(id)
	if account == nil {
		fmt.Println("Account Not found")
		return
	}

	amountToAdd := get_valid_amount(reader)
	fmt.Print("Note (Optional): ")
	noteInput, _ := reader.ReadString('\n')
	note := strings.TrimSpace(noteInput)

	addMoneyToAccount(account, amountToAdd, note)
	fmt.Printf("✓ Added $%.2f to '%s'. New balance: $%.2f\n", amountToAdd, account.Name, account.Balance)
}

func handleAccountHistory() {
	if len(accounts) == 0 {
		fmt.Println("[Warning] No Account Available.")
		return
	}

	fmt.Println("Acouunt History")
	fmt.Println("===============")

	fmt.Println("Available Account:")

	for _, account := range accounts {
		fmt.Printf("Account ID: %d | Account Name: %s\n", account.ID, account.Name)
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter Account ID: ")
	id_input, _ := reader.ReadString('\n')
	id, err := strconv.Atoi(strings.TrimSpace(id_input))

	if err != nil || id < 0 || id > len(accounts) {
		fmt.Println("[Error] Invalid ID. Please Enter An Integer In The Available Range!!!")
		return
	}

	var accountName string
	account := findAccount(id)
	if account == nil {
		fmt.Println("No Account Found!")
		return
	}
	accountName = account.Name

	fmt.Printf("\nHistory for Account %s:\n", accountName)
	fmt.Println(strings.Repeat("=", len(accountName)+20))

	found := false
	for _, transac := range accountTransactions {
		if transac.AccountID == id {
			fmt.Printf("%s | +$%.2f | %s, %d\n", transac.Date.Format("2006-01-02 15:04"), transac.Amount, transac.Note, transac.AccountID)
			found = true
		}
	}

	if !found {
		fmt.Println("[Warning] No Transaction Found For This Account.")
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

func getTransactionNumberToShow(reader *bufio.Reader, defaultValue int) int {
	transactionsToShow := defaultValue
InputLoop:
	for {
		input, _ := reader.ReadString('\n')
		input = strings.ToLower(strings.TrimSpace(input))

		switch input {
		case "":
			fmt.Printf("\nDisplaying %d Recent transactions:\n", transactionsToShow)
			break InputLoop
		case "all":
			fmt.Println("\nDisplaying All transactions:")
			transactionsToShow = len(transactions)
			break InputLoop
		default:
			number, err := strconv.Atoi(input)
			if err != nil || number <= 0 {
				fmt.Println("[Warning] Not a Valid Number, Try Again.")
				continue
			}
			fmt.Printf("\nDisplaying Last %d Transactions:\n", number)
			transactionsToShow = number
			break InputLoop
		}
	}
	return transactionsToShow
}

func handleDeleteTransaction() {
	if len(transactions) == 0 {
		fmt.Println("[Info] No Transaction to Delete")
		return
	}

	fmt.Println("Delete Transaction")
	fmt.Println("===================")

	fmt.Printf("\n[Info] Total Transaction Number: %d.\n", len(transactions))
	transactionsToShow := 10
	fmt.Println("\nPress Enter to show the latest 10 transactions (default).\nTo see a different number, enter it.\nType 'all' to display all transactions. Enter: ")

	reader := bufio.NewReader(os.Stdin)
	transactionsToShow = getTransactionNumberToShow(reader, transactionsToShow)

	startingIndex := len(transactions) - transactionsToShow
	if startingIndex < 0 {
		startingIndex = 0
	}

	for i := startingIndex; i < len(transactions); i++ {
		t := transactions[i]
		fmt.Printf("Transaction ID: %d | %s |  %s | %.2f | %s \n", t.ID, t.Date.Format("2006-01-02"), t.Type, t.Amount, t.Category)
	}

	fmt.Print("\nPlease Enter The Transaction ID to Delete (or 'cancel'): ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "cancel" {
		fmt.Println("Deletion cancelled.")
		return
	}
	transacId, err := strconv.Atoi(input)
	if err != nil {
		fmt.Println("[Warning] Not a Valid Transaction ID.")
		return
	}

	isTransactionDeleted := deleteTransaction(startingIndex, transacId)
	if isTransactionDeleted {
		fmt.Printf("✓ Transaction ID %d deleted successfully!\n", transacId)
		return
	}
	fmt.Printf("Transaction ID %d is not in the list please choice an ID that presented in the list!\n", transacId)

}

func handleDeleteAccount() {
	if len(accounts) == 0 {
		fmt.Println("[Info] No Accounts to Delete")
		return
	}

	fmt.Println("Delete Account")
	fmt.Println("==============")

	fmt.Printf("\n[Info] Total Account Number: %d.\n", len(accounts))

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

	fmt.Print("[Info] Are You SURE? This will delete all account hitory and all account transactions. (yes/anythin else considered as no): ")
	confirm, _ := reader.ReadString('\n')
	confirm = strings.ToLower(strings.TrimSpace(confirm))

	if confirm != "yes" {
		fmt.Println("Deletion cancelled.")
		return
	}

	isDeleted := deleteAccount(accountId)
	if isDeleted {
		fmt.Printf("✓ Account ID %d deleted successfully!\n", accountId)
		return
	}
	fmt.Printf("Account ID %d not found!\n", accountId)
}

func handleEditTransaction() {
	if len(transactions) == 0 {
		fmt.Println("[Info] No Transaction to Delete")
		return
	}

	fmt.Println("\nEdit Transaction")
	fmt.Println("=================")

	fmt.Printf("\n[Info] Total Transaction Number: %d.\n", len(transactions))
	transactionsToShow := 10
	fmt.Println("\nPress Enter to show the latest 10 transactions (default).\nTo see a different number, enter it.\nType 'all' to display all transactions.")

	reader := bufio.NewReader(os.Stdin)
	getTransactionNumberToShow(reader, transactionsToShow)

	startingIndex := len(transactions) - transactionsToShow
	if startingIndex < 0 {
		startingIndex = 0
	}

	for i := startingIndex; i < len(transactions); i++ {
		t := transactions[i]
		fmt.Printf("Transaction ID: %d | %s |  %s | %.2f | %s \n", t.ID, t.Date.Format("2006-01-02"), t.Type, t.Amount, t.Category)
	}

	fmt.Print("\nPlease Enter The Transaction ID to Edit (or 'cancel'): ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "cancel" {
		fmt.Println("Edit cancelled.")
		return
	}
	transacId, err := strconv.Atoi(input)
	if err != nil {
		fmt.Println("[Warning] Not a Valid Transaction ID.")
		return
	}

	transac := findTransaction(transacId)
	fmt.Printf("transcation: %v\n", transac)
	// .... continue
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
