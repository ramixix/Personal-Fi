package cli

import (
	"bufio"
	"financial_tracker/internal/core"
	"financial_tracker/internal/models"
	"financial_tracker/internal/utils"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// ================================ Handels Adding & Editing & Deleting Transactions ================================
// handleAddTransaction handlse adding a new transaction
func handleAddTransaction() {
	fmt.Println("Add Transaction")
	fmt.Println("================")

	reader := bufio.NewReader(os.Stdin)

	// Get validated inputs from users
	transactionType := GetValidTransactionType(reader)
	amount := utils.GetValidAmount(reader)
	category := utils.GetNonEmptyString(reader, "Category: ")
	description := utils.GetNonEmptyString(reader, "Description: ")
	currencyCode := utils.GetValidCurrency(reader, "USD")

	if len(currencyCode) != 3 {
		currencyCode = ""
	}

	newTransaction := models.Transaction{
		ID:           utils.MustGenerateUUID(),
		Date:         time.Now(),
		Amount:       amount,
		Category:     category,
		Description:  description,
		Type:         transactionType,
		CurrencyCode: currencyCode,
	}

	err := core.AddTransaction(newTransaction)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("\n✓ Transaction added successfully! ID: %s\n", newTransaction.ID)
}

// handleEditTransaction handles editing a transaction
func handleEditTransaction() {
	transactionsLength := core.GetTransactionsLength("")
	if transactionsLength == 0 {
		fmt.Println("[Info] No Transaction to Edit")
		return
	}

	fmt.Println("\nEdit Transaction")
	fmt.Println("=================")
	fmt.Printf("\n[Info] Total Transaction Number: %d.\n", transactionsLength)

	transactionsToShow := 10
	if transactionsLength < transactionsToShow {
		transactionsToShow = transactionsLength
	}

	fmt.Printf("\nPress Enter to show the latest %d transactions (by default).\nTo see a different number, enter it.\nType 'all' to display all transactions.", transactionsToShow)
	reader := bufio.NewReader(os.Stdin)
	transactionsToShow = GetTransactionNumberToShow(reader, transactionsToShow)

	ListFilteredTransactions(core.GetRecentTransactions(transactionsToShow), fmt.Sprintf("listing last %d transactions", transactionsToShow))

	input := utils.GetNonEmptyString(reader, "\nPlease Enter The Transaction ID to edit (or 'cancel'): ")
	if input == "cancel" {
		fmt.Println("Edit cancelled.")
		return
	}
	txID := input

	transac := core.FindTransaction(txID)
	if transac == nil {
		fmt.Printf("Transaction ID %s not found!\n", txID)
		return
	}

	fmt.Printf("\n--- Editing Transaction ID %s ---\n", txID)
	fmt.Println("Press Enter to keep current value")

	// Edit type
	fmt.Printf("Type income/expnese (current %s): ", transac.Type)
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
	fmt.Printf("Amount (current: %s): ", utils.FormatCurrency(transac.Amount, transac.CurrencyCode))
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

	// Edit Currency code
	fmt.Printf("Currency Code(current %s): ", transac.CurrencyCode)
	transac.CurrencyCode = utils.GetValidCurrency(reader, transac.CurrencyCode)

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

	fmt.Printf("\n✓ Transaction ID %s updated successfully!\n", txID)
	fmt.Printf("Updated: %s | %s | %s | %s | %s | %s\n",
		transac.Date.Format("2006-01-02"),
		transac.Type,
		utils.FormatCurrency(transac.Amount, transac.CurrencyCode),
		transac.CurrencyCode,
		transac.Category,
		transac.Description)
}

// handleDeleteTransaction handles deleting a transaction
func handleDeleteTransaction() {
	transactionsLength := core.GetTransactionsLength("")
	if transactionsLength == 0 {
		fmt.Println("[Info] No Transaction to Delete")
		return
	}

	fmt.Println("Delete Transaction")
	fmt.Println("===================")
	fmt.Printf("\n[Info] Total Transaction Number: %d.\n", transactionsLength)

	transactionsToShow := 10
	if transactionsLength < transactionsToShow {
		transactionsToShow = transactionsLength
	}

	fmt.Printf("\nPress Enter to show the latest %d transactions (default).\nEnter different number if you want.\nType 'all' to display all transactions. Your Input: ", transactionsToShow)

	reader := bufio.NewReader(os.Stdin)
	transactionsToShow = GetTransactionNumberToShow(reader, transactionsToShow)

	ListFilteredTransactions(core.GetRecentTransactions(transactionsToShow), fmt.Sprintf("listing recent %d transactions", transactionsToShow))

	input := utils.GetNonEmptyString(reader, "\nPlease Enter The Transaction ID to Delete (or 'cancel'): ")
	if input == "cancel" {
		fmt.Println("Deletion cancelled.")
		return
	}
	txID := input
	if transac := core.FindTransaction(txID); transac == nil {
		fmt.Printf("Transaction ID %s not found!\n", txID)
		return
	}

	err := core.DeleteTransaction(txID)
	if err != nil {
		fmt.Printf("error : %v", err)
		fmt.Printf("Transaction ID %s is not in the list please enter an ID that is in the list!\n", txID)
		return
	}
	fmt.Printf("✓ Transaction ID %s deleted successfully!\n", txID)

}

// ================================ Handels Listing & Filtering Transactions ================================

// handleListTransactions handles listing transactions
func handleListTransactions() {
	fmt.Println("List Transactions")
	fmt.Println("=================")

	if len(os.Args) > 2 {
		filter := strings.ToLower(os.Args[2])
		switch filter {
		case "week":
			ListFilteredTransactions(core.GetTransactionsForLastDays(7), "Last 7 Days (Last Week)")
		case "month":
			ListFilteredTransactions(core.GetTransactionsForLastDays(30), "Last 30 Days (Last Month)")
		case "year":
			ListFilteredTransactions(core.GetTransactionsForLastDays(365), "Last 365 Days (Last Year)")
		case "income":
			ListFilteredTransactions(core.GetTransactionsByType("income"), "All Incomes")
		case "expense":
			ListFilteredTransactions(core.GetTransactionsByType("expense"), "All Expenses")
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
		ListFilteredTransactions(core.GetRecentTransactions(recent100), fmt.Sprintf("Last %d Transactions if available", recent100))
	}
}

// handleCategoryFilter handles category-specific filtering
func handleCategoryFilter() {
	if core.GetTransactionsLength("") == 0 {
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
	ListFilteredTransactions(filtered, fmt.Sprintf("Transactions in category: %s", category))
}

// showCategories shows all categories with transaction counts, total income/expense
func showCategories() {
	fmt.Println("All Categroies And Amount In Total")
	fmt.Println("===================================")

	summaries := core.GetCategorySummary()
	if len(summaries) == 0 {
		fmt.Println("No transactions available.")
		return
	}

	for _, summary := range summaries {
		netAmount := summary.TotalIncome - summary.TotalExpense

		fmt.Printf("%-15s (%s): %d Transactions | Total Income: %-10s | Total Expenses: %-10s | Net Amount: %s\n",
			summary.Category,
			summary.CurrencyCode,
			summary.Count,
			utils.FormatCurrency(summary.TotalIncome, summary.CurrencyCode),
			utils.FormatCurrency(summary.TotalExpense, summary.CurrencyCode),
			utils.FormatCurrency(netAmount, summary.CurrencyCode))
	}
}

// handleCustomRange handles custom date range filtering
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
		fmt.Println("[Error] End date must be After start date, there is a logical error")
		return
	}

	filtered := core.GetTransactionsForDateRange(start_date, end_date)
	title := fmt.Sprintf("Transaction From %s to %s", start_date.Format("2006-01-02"), end_date.Format("2006-01-02"))
	ListFilteredTransactions(filtered, title)
}

// GetTransactionNumberToShow asks users for a specific number N to list last N transactions.
func GetTransactionNumberToShow(reader *bufio.Reader, defaultValue int) int {
	transactionsToShow := defaultValue
	totalTransactions := core.GetTransactionsLength("")
	if defaultValue > totalTransactions {
		transactionsToShow = totalTransactions
	}
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
			transactionsToShow = totalTransactions
			break InputLoop
		default:
			number, err := strconv.Atoi(input)
			if err != nil || number <= 0 || number > totalTransactions {
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

// GetValidTransactionType ensures that the user selects one of the keywords “income” or “expense” as the transaction type.
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

// ListFilteredTransactions displays filtered transactions with summary
func ListFilteredTransactions(transaction_list []models.Transaction, title string) {
	fmt.Printf("%s\n", title)
	fmt.Printf("%s\n", strings.Repeat("=", len(title)))

	if len(transaction_list) == 0 {
		fmt.Println("No transactions found.")
		return
	}

	var totals map[string]models.CurrencyTotal = make(map[string]models.CurrencyTotal)

	for _, transaction := range transaction_list {
		fmt.Printf("ID: %-6s | %15s | %-8s | %-10s | %-3s | %-20s | %s\n",
			transaction.ID,
			transaction.Date.Format("2006-01-02 15:04"),
			transaction.Type,
			utils.FormatCurrency(transaction.Amount, transaction.CurrencyCode),
			transaction.CurrencyCode,
			transaction.Category,
			transaction.Description)

		total := totals[transaction.CurrencyCode]
		if transaction.Type == "income" {
			total.Income += transaction.Amount
		} else {
			total.Expenses += transaction.Amount
		}
		totals[transaction.CurrencyCode] = total
	}

	fmt.Printf("\n--- Summary ---\n")
	for cur, total := range totals {
		fmt.Printf("Currency:   %2s\n", cur)
		fmt.Printf("Total Income:   %.2f\n", total.Income)
		fmt.Printf("Total Expenses: %.2f\n", total.Expenses)
		fmt.Printf("Net Amount:     %.2f\n", total.Income-total.Expenses)

	}
}
