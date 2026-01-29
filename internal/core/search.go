package core

import (
	"financial_tracker/internal/models"
	"financial_tracker/internal/storage"
	"strings"
)

// Search transactions by keyword in description and category
func SearchTransactionsByKeyword(keyword string) []models.Transaction {
	if keyword == "" {
		return nil
	}
	var results []models.Transaction
	keyword = strings.ToLower(keyword)

	for _, transac := range storage.Transactions {
		description := strings.ToLower(transac.Description)
		category := strings.ToLower(transac.Category)

		if strings.Contains(description, keyword) || strings.Contains(category, keyword) {
			results = append(results, transac)
		}
	}
	return results
}

// Search transactions by amount range
func SearchTransactionsByAmountRange(min, max float64) []models.Transaction {
	var results []models.Transaction

	for _, transaction := range storage.Transactions {
		if transaction.Amount >= min && transaction.Amount <= max {
			results = append(results, transaction)
		}
	}
	return results
}

// Advanced search with multiple criteria
func AdvancedSearchTransactions(criteria models.SearchCriteria) []models.Transaction {
	var results []models.Transaction

	for _, transac := range storage.Transactions {
		// If keyword is set then check if the transaction have the keyword if not go to next transaction
		if criteria.Keyword != "" {
			keyword := strings.ToLower(criteria.Keyword)
			description := strings.ToLower(transac.Description)
			category := strings.ToLower(transac.Category)

			if !strings.Contains(description, keyword) && !strings.Contains(category, keyword) {
				continue
			}
		}

		// After keyword if type is set then check if the types are equal if not go to next transaction
		if criteria.TransactionType != "" {
			if transac.Type != criteria.TransactionType {
				continue
			}
		}

		// Now check if category of transaction is mentioned if not go to next transaction
		if len(criteria.Categories) > 0 {
			found := false
			for _, cat := range criteria.Categories {
				if strings.EqualFold(transac.Category, cat) {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		// If max/min or any one of them are set then check transaction amount to be in that range if not go to the next transaction
		if criteria.MinAmount > 0 || criteria.MaxAmount > 0 {
			if criteria.MaxAmount > 0 {
				if transac.Amount > criteria.MaxAmount || transac.Amount < criteria.MinAmount {
					continue
				}
			} else {
				if transac.Amount < criteria.MinAmount {
					continue
				}
			}
		}

		// check if user had set date range and check if transaction date is in that range if not skip
		if criteria.HasDateRange {
			if transac.Date.Before(criteria.StartDate) || transac.Date.After(criteria.EndDate) {
				continue
			}
		}

		// if a transaction make to this point it means it passed all the test and completely meets all ceriteria
		results = append(results, transac)
	}
	return results
}

// Find similar transactions to specific transaction by controling that the categories are same and have similar amount (by similar mean be in range of tolerance)
func FindSimilarTransactions(referenceTransaction models.Transaction, amountTolerance float64) []models.Transaction {
	var results []models.Transaction

	for _, transac := range storage.Transactions {
		// Skip the reference transaction itself
		if transac.ID == referenceTransaction.ID {
			continue
		}

		// Check same category
		if !strings.EqualFold(transac.Category, referenceTransaction.Category) {
			continue
		}

		// Check similar amount (within tolerance percentage)
		tolerance := referenceTransaction.Amount * (amountTolerance / 100)
		minAmount := referenceTransaction.Amount - tolerance
		maxAmount := referenceTransaction.Amount + tolerance

		if transac.Amount >= minAmount && transac.Amount <= maxAmount {
			results = append(results, transac)
		}
	}
	return results
}

// Search accounts by name
func SearchAccountsByName(keyword string) []models.Account {
	var results []models.Account
	keywordLower := strings.ToLower(keyword)

	for _, account := range storage.Accounts {
		nameLower := strings.ToLower(account.Name)
		if strings.Contains(nameLower, keywordLower) {
			results = append(results, account)
		}
	}

	return results
}

// Search goals by name or description
func SearchGoalsByKeyword(keyword string) []models.Goal {
	var results []models.Goal
	keywordLower := strings.ToLower(keyword)

	for _, goal := range storage.Goals {
		nameLower := strings.ToLower(goal.Name)
		descLower := strings.ToLower(goal.Description)

		if strings.Contains(nameLower, keywordLower) || strings.Contains(descLower, keywordLower) {
			results = append(results, goal)
		}
	}

	return results
}

// Get transactions by multiple categories
func GetTransactionsByMultipleCategories(categories []string) []models.Transaction {
	var results []models.Transaction

	for _, transaction := range storage.Transactions {
		for _, cat := range categories {
			if strings.EqualFold(transaction.Category, cat) {
				results = append(results, transaction)
				break
			}
		}
	}

	return results
}

// Get highest spending transactions (top N)
func GetTopSpendingTransactions(limit int) []models.Transaction {

	if limit <= 0 {
		return []models.Transaction{}
	}

	// Get only expenses
	expenses := GetTransactionsByType("expense")

	// Sort by amount (bubble sort for simplicity)
	for i := 0; i < len(expenses); i++ {
		for j := i + 1; j < len(expenses); j++ {
			if expenses[j].Amount > expenses[i].Amount {
				expenses[i], expenses[j] = expenses[j], expenses[i]
			}
		}
	}

	// Return top N
	if limit > len(expenses) {
		limit = len(expenses)
	}

	return expenses[:limit]
}

// Get recent transactions (last N)
func GetRecentTransactions(limit int) []models.Transaction {
	if limit > len(storage.Transactions) {
		limit = len(storage.Transactions)
	}

	if limit <= 0 {
		return []models.Transaction{}
	}

	// Return last N transactions
	return storage.Transactions[len(storage.Transactions)-limit:]
}

// Get recent goals (last N)
func GetRecentGoalContributions(limit int) []models.GoalContribution {
	if limit > len(storage.GoalContributions) {
		limit = len(storage.GoalContributions)
	}

	if limit <= 0 {
		return []models.GoalContribution{}
	}

	// Return last N transactions
	return storage.GoalContributions[len(storage.GoalContributions)-limit:]
}
