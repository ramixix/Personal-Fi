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
)

func handleSearch() {
	if len(os.Args) < 3 {
		fmt.Println("Search")
		fmt.Println("======")
		fmt.Println("Usage: financial-tracker search [type]")
		fmt.Println("Types:")
		fmt.Println("  keyword      Search transactions by keyword")
		fmt.Println("  amount       Search transactions by amount range")
		fmt.Println("  advanced     Advanced search with multiple criteria")
		fmt.Println("  similar      Find similar transactions")
		fmt.Println("  recent       Show recent transactions")
		return
	}

	searchType := os.Args[2]
	switch searchType {
	case "keyword":
		handleKeywordSearch()
	case "amount":
		handleAmountSearch()
	case "advanced":
		handleAdvancedSearch()
	case "recent":
		handleRecentTransactions()
	default:
		fmt.Printf("Unknown search type: %s\n", searchType)
		fmt.Println("Available types: keyword, amount, advanced, similar, top, recent")
	}
}

// Handle keyword search
func handleKeywordSearch() {
	fmt.Println("Keyword Search")
	fmt.Println("==============")

	reader := bufio.NewReader(os.Stdin)
	keyword := utils.GetNonEmptyString(reader, "Enter keyword to search: ")

	results := core.SearchTransactionsByKeyword(keyword)

	if len(results) == 0 {
		fmt.Printf("No transactions found containing '%s'\n", keyword)
		return
	}

	title := fmt.Sprintf("Search results for '%s'", keyword)
	ListFilteredTransactions(results, title)
}

// Handle amount range search
func handleAmountSearch() {
	fmt.Println("Amount Range Search")
	fmt.Println("===================")

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Minimum  ")
	minAmount := utils.GetValidAmount(reader)
	fmt.Print("Maximum ")
	maxAmount := utils.GetValidAmount(reader)

	if minAmount > maxAmount {
		fmt.Println("Minimum amount cannot be greater than maximum amount!")
		return
	}

	results := core.SearchTransactionsByAmountRange(minAmount, maxAmount)

	if len(results) == 0 {
		fmt.Printf("No transactions found between $%.2f and $%.2f\n", minAmount, maxAmount)
		return
	}

	title := fmt.Sprintf("Transactions between $%.2f and $%.2f", minAmount, maxAmount)
	ListFilteredTransactions(results, title)
}

// Handle advanced search
func handleAdvancedSearch() {
	fmt.Println("Advanced Search")
	fmt.Println("===============")
	fmt.Println("Enter search criteria (press Enter to skip any field)")

	var criteria models.SearchCriteria

	reader := bufio.NewReader(os.Stdin)

	// Keyword
	fmt.Print("\nKeyword (searches in description and category): ")
	keyword, _ := reader.ReadString('\n')
	criteria.Keyword = strings.TrimSpace(keyword)

	// Transaction Type
	fmt.Print("Transaction type (income/expense) (Leave this field empty if the type can be either income or expense.): ")
	transactionType, _ := reader.ReadString('\n')
	transactionType = strings.TrimSpace(transactionType)
	if transactionType == "income" || transactionType == "expense" {
		criteria.TransactionType = transactionType
	}

	// Categories
	categoryInput, _ := reader.ReadString('\n')
	categoryInput = strings.TrimSpace(categoryInput)
	if categoryInput != "" {
		categories := strings.Split(categoryInput, ",")
		for i := range categories {
			categories[i] = strings.TrimSpace(categories[i])
		}
		criteria.Categories = categories
	}

	// Amount Range
	fmt.Print("Minimum amount (or 0 to skip): $")
	minInput, _ := reader.ReadString('\n')
	minAmount, err := strconv.ParseFloat(strings.TrimSpace(minInput), 64)
	if err == nil && minAmount > 0 {
		criteria.MinAmount = minAmount
	}

	fmt.Print("Maximum amount (or 0 to skip): $")
	maxInput, _ := reader.ReadString('\n')
	maxAmount, err := strconv.ParseFloat(strings.TrimSpace(maxInput), 64)
	if err == nil && maxAmount > 0 {
		criteria.MaxAmount = maxAmount
	}

	// Date range
	fmt.Print("Filter by date range? (yes/no): ")
	dateChoice, _ := reader.ReadString('\n')
	if strings.ToLower(strings.TrimSpace(dateChoice)) == "yes" {
		fmt.Print("Start date (YYYY-MM-DD): ")
		startInput, _ := reader.ReadString('\n')
		startDate, err := utils.ParseDate(strings.TrimSpace(startInput))
		if err != nil {
			fmt.Println("Invalid start date, skipping date filter")
		} else {
			fmt.Print("End date (YYYY-MM-DD): ")
			endInput, _ := reader.ReadString('\n')
			endDate, err := utils.ParseDate(strings.TrimSpace(endInput))
			if err != nil {
				fmt.Println("Invalid end date, skipping date filter")
			} else {
				criteria.StartDate = startDate
				criteria.EndDate = endDate
				criteria.HasDateRange = true
			}
		}
	}

	// Perform search
	results := core.GetTransactionsAdvanceSearch(criteria)

	if len(results) == 0 {
		fmt.Println("\nNo transactions found matching your criteria.")
		return
	}

	fmt.Println() // Empty line
	ListFilteredTransactions(results, "Advanced Search Results")
}

// handleRecentTransactions returns recent transactions
func handleRecentTransactions() {
	fmt.Println("Recent Transactions")
	fmt.Println("===================")

	reader := bufio.NewReader(os.Stdin)
	prompt := "How many recent transactions to show? (default: 10): "
	limitInput, err := utils.GetIntInput(reader, prompt)

	limit := 10
	if err == nil && limitInput > 0 {
		limit = limitInput
	}

	results := core.GetRecentTransactions(limit)

	if len(results) == 0 {
		fmt.Println("No transactions found.")
		return
	}

	title := fmt.Sprintf("Last %d Transactions", len(results))
	ListFilteredTransactions(results, title)
}
