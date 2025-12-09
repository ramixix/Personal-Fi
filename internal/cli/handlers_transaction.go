package cli

import (
	"bufio"
	"financial_tracker/internal/core"
	"financial_tracker/internal/models"
	"financial_tracker/internal/storage"
	"financial_tracker/internal/utils"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Handle adding a new transaction
func handleAddTransaction() {
	fmt.Println("Add Transaction")
	fmt.Println("================")

	reader := bufio.NewReader(os.Stdin)

	// Get validated inputs from users
	transactionType := GetValidTransactionType(reader)
	amount := utils.GetValidAmount(reader)
	category := utils.GetNonEmptyString(reader, "Category: ")
	description := utils.GetNonEmptyString(reader, "Description: ")

	newTransaction := models.Transaction{
		ID:          storage.NextTransactionID,
		Date:        time.Now(),
		Amount:      amount,
		Category:    category,
		Description: description,
		Type:        transactionType,
	}

	core.AddTransaction(newTransaction)
	fmt.Printf("\n✓ Transaction added successfully! ID: %d\n", newTransaction.ID)
}

// Handle listing all transactions
func handleListTransactions() {
	fmt.Println("List Transactions")
	fmt.Println("=================")

	if len(os.Args) > 2 {
		filter := strings.ToLower(os.Args[2])
		switch filter {
		case "week":
			core.ListFilteredTransactions(core.GetTransactionsByDateRange(7), "Last 7 Days (Last Week)")
		case "month":
			core.ListFilteredTransactions(core.GetTransactionsByDateRange(30), "Last 30 Days (Last Month)")
		case "year":
			core.ListFilteredTransactions(core.GetTransactionsByDateRange(365), "Last 365 Days (Last Year)")
		case "income":
			core.ListFilteredTransactions(core.GetTransactionsByType("income"), "All Incomes")
		case "expense":
			core.ListFilteredTransactions(core.GetTransactionsByType("expense"), "All Expenses")
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
		core.ListFilteredTransactions(storage.Transactions, "All Transactions")
	}
}

// Handle category-specific filtering
func handleCategoryFilter() {
	if len(storage.Transactions) == 0 {
		fmt.Println("No transactions available.")
		return
	}

	fmt.Println("Available Categories:")
	categories := core.GetCategories()

	for i, cat := range categories {
		fmt.Printf("  %d. %s\n", i+1, cat)
	}

	reader := bufio.NewReader(os.Stdin)
	category := utils.GetNonEmptyString(reader, "Please Enter The Category Name You Want to Filtere:")
	filtered := core.GetTransactionsByCategory(category)
	core.ListFilteredTransactions(filtered, fmt.Sprintf("Transactions in category: %s", category))
}

// Show all categories with transaction counts, total income/expense
func showCategories() {
	fmt.Println("All Categroies And Amount In Total")
	fmt.Println("==================================================================================================...")

	if len(storage.Transactions) == 0 {
		fmt.Println("No transactions available.")
		return
	}

	categoryCount := make(map[string]int)
	categoryTotal := make(map[string]float64)
	categoryIncome := make(map[string]float64)
	categoryExpense := make(map[string]float64)

	for _, transac := range storage.Transactions {
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
		fmt.Printf("%-15s: %d Transactions | Total Income: %-10.2f | Total Expenses: %-10.2f | Total Amount: $%.2f\n", category, categoryCount[category], categoryIncome[category], categoryExpense[category], amount)
	}
}

// Handle custom date range filtering
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
	start_date, err := utils.ParseDate(start_input)
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
	end_date, err := utils.ParseDate(end_input)
	if err != nil {
		fmt.Println("[Error] Not a date in specified format (Formant : YYYY-MM-DD)")
		return
	}

	if end_date.Before(start_date) {
		fmt.Println("[Error] End date must be before start date, there is a logical error")
		return
	}

	filtered := core.GetTransactionsByCustomRange(start_date, end_date)
	title := fmt.Sprintf("Transaction From %s to %s", start_date.Format("2006-01-02"), end_date.Format("2006-01-02"))
	core.ListFilteredTransactions(filtered, title)
}

// Handle delete transaction command
func handleDeleteTransaction() {
	if len(storage.Transactions) == 0 {
		fmt.Println("[Info] No Transaction to Delete")
		return
	}

	fmt.Println("Delete Transaction")
	fmt.Println("===================")

	fmt.Printf("\n[Info] Total Transaction Number: %d.\n", len(storage.Transactions))
	transactionsToShow := 10
	fmt.Println("\nPress Enter to show the latest 10 transactions (default).\nEnter different number if you want.\nType 'all' to display all transactions. Your Input: ")

	reader := bufio.NewReader(os.Stdin)
	transactionsToShow = GetTransactionNumberToShow(reader, transactionsToShow)

	startingIndex := len(storage.Transactions) - transactionsToShow
	if startingIndex < 0 {
		startingIndex = 0
	}

	for i := startingIndex; i < len(storage.Transactions); i++ {
		t := storage.Transactions[i]
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

	isTransactionDeleted := core.DeleteTransaction(startingIndex, transacId)
	if isTransactionDeleted {
		fmt.Printf("✓ Transaction ID %d deleted successfully!\n", transacId)
		return
	}
	fmt.Printf("Transaction ID %d is not in the list please enter an ID that is in the list!\n", transacId)

}

// Handle edit transaction command
func handleEditTransaction() {
	if len(storage.Transactions) == 0 {
		fmt.Println("[Info] No Transaction to Delete")
		return
	}

	fmt.Println("\nEdit Transaction")
	fmt.Println("=================")

	fmt.Printf("\n[Info] Total Transaction Number: %d.\n", len(storage.Transactions))
	transactionsToShow := 10
	fmt.Println("\nPress Enter to show the latest 10 transactions (default).\nTo see a different number, enter it.\nType 'all' to display all transactions.")

	reader := bufio.NewReader(os.Stdin)
	transactionsToShow = GetTransactionNumberToShow(reader, transactionsToShow)

	startingIndex := len(storage.Transactions) - transactionsToShow
	if startingIndex < 0 {
		startingIndex = 0
	}

	for i := startingIndex; i < len(storage.Transactions); i++ {
		t := storage.Transactions[i]
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

	transac := core.FindTransaction(transacId)
	if transac == nil {
		fmt.Printf("Transaction ID %d not found!\n", transacId)
		return
	}

	fmt.Printf("\n--- Editing Transaction ID %d ---\n", transacId)
	fmt.Println("Press Enter to keep current value")

	// Edit type
	fmt.Printf("Type income/expnese (current %s)", transac.Type)
	typeInput, _ := reader.ReadString('\n')
	typeInput = strings.TrimSpace(typeInput)
	if typeInput != "" {
		typeInput = strings.ToLower(typeInput)
		if typeInput == "income" || typeInput == "expense" {
			transac.Type = typeInput
		} else {
			fmt.Println("Invalid type, keeping current value")
		}
	}

	// Edit Amount
	fmt.Printf("Amount (current: $%.2f): ", transac.Amount)
	amountInput, _ := reader.ReadString('\n')
	amountInput = strings.TrimSpace(amountInput)
	if amountInput != "" {
		amount, err := strconv.ParseFloat(amountInput, 64)
		if err == nil && amount > 0 {
			transac.Amount = amount
		} else {
			fmt.Println("Invalid amount, keeping current value")
		}
	}

	// Edit category
	fmt.Printf("Category (current: %s): ", transac.Category)
	categoryInput, _ := reader.ReadString('\n')
	categoryInput = strings.TrimSpace(categoryInput)
	if categoryInput != "" {
		transac.Category = categoryInput
	}

	// Edit description
	fmt.Printf("Description (current: %s): ", transac.Description)
	descInput, _ := reader.ReadString('\n')
	descInput = strings.TrimSpace(descInput)
	if descInput != "" {
		transac.Description = descInput
	}

	// Edit date
	fmt.Printf("Date (current: %s, format YYYY-MM-DD): ", transac.Date.Format("2006-01-02"))
	dateInput, _ := reader.ReadString('\n')
	dateInput = strings.TrimSpace(dateInput)
	if dateInput != "" {
		newDate, err := utils.ParseDate(dateInput)
		if err == nil {
			transac.Date = newDate
		} else {
			fmt.Println("Invalid date format, keeping current value")
		}
	}

	fmt.Printf("\n✓ Transaction ID %d updated successfully!\n", transacId)
	fmt.Printf("Updated: %s | %s | $%.2f | %s | %s\n",
		transac.Date.Format("2006-01-02"),
		transac.Type,
		transac.Amount,
		transac.Category,
		transac.Description)

}

func GetTransactionNumberToShow(reader *bufio.Reader, defaultValue int) int {
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
			transactionsToShow = len(storage.Transactions)
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

func GetValidTransactionType(reader *bufio.Reader) string {
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
