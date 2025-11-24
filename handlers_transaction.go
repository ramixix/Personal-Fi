package main

import (
	"bufio"
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
	transactionType := getValidTransactionType(reader)
	amount := getValidAmount(reader)
	category := getNonEmptyString(reader, "Category: ")
	description := getNonEmptyString(reader, "Description: ")

	newTransaction := Transaction{
		ID:          nextTransactionID,
		Date:        time.Now(),
		Amount:      amount,
		Category:    category,
		Description: description,
		Type:        transactionType,
	}

	addTransaction(newTransaction)
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

// Handle category-specific filtering
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
	category := getNonEmptyString(reader, "Please Enter The Category Name You Want to Filtere:")
	filtered := getTransactionsByCategory(category)
	listFilteredTransactions(filtered, fmt.Sprintf("Transactions in category: %s", category))
}

// Show all categories with transaction counts, total income/expense
func showCategories() {
	fmt.Println("All Categroies And Amount In Total")
	fmt.Println("==================================")

	if len(transactions) == 0 {
		fmt.Println("No transactions available.")
		return
	}

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
		fmt.Printf("%s: %d Transactions | Total Income: %.2f | Total Expenses: %.2f | Total Amount: $%.2f\n", category, categoryCount[category], categoryIncome[category], categoryExpense[category], amount)
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

// Handle delete transaction command
func handleDeleteTransaction() {
	if len(transactions) == 0 {
		fmt.Println("[Info] No Transaction to Delete")
		return
	}

	fmt.Println("Delete Transaction")
	fmt.Println("===================")

	fmt.Printf("\n[Info] Total Transaction Number: %d.\n", len(transactions))
	transactionsToShow := 10
	fmt.Println("\nPress Enter to show the latest 10 transactions (default).\nEnter different number if you want.\nType 'all' to display all transactions. Your Input: ")

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
	fmt.Printf("Transaction ID %d is not in the list please enter an ID that is in the list!\n", transacId)

}

// Handle edit transaction command
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
	transactionsToShow = getTransactionNumberToShow(reader, transactionsToShow)

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
		newDate, err := parseDate(dateInput)
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
