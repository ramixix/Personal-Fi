package core

import (
	"encoding/csv"
	"financial_tracker/internal/models"
	"financial_tracker/internal/utils"
	"fmt"
	"os"
	"strconv"
)

// ExportTransactionsToCSV exports all transactions to CSV
func ExportTransactionsToCSV(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{"ID", "Date", "Type", "Amount", "Currency", "Category", "Description"}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("failed to write CSV header: %w", err)
	}

	const batchSize = 500
	offset := 0

	for {
		transactions := GetTransactionBatch(batchSize, offset)
		if len(transactions) == 0 {
			break
		}

		for _, transaction := range transactions {
			record := []string{
				transaction.ID,
				transaction.Date.Format("2006-01-02"),
				transaction.Type,
				fmt.Sprintf("%.2f", transaction.Amount),
				transaction.CurrencyCode,
				transaction.Category,
				transaction.Description,
			}
			if err := writer.Write(record); err != nil {
				return fmt.Errorf("failed to write record %s: %w", transaction.ID, err)
			}
		}

		writer.Flush()
		if err := writer.Error(); err != nil {
			return fmt.Errorf("error flushing to file: %w", err)
		}

		offset += batchSize
	}
	return nil
}

// ExportAccountsToCSV exports all accounts to CSV
func ExportAccountsToCSV(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{"ID", "Name", "Balance", "Currency", "Created"}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("failed to write CSV header: %w", err)
	}

	const batchSize = 500
	offset := 0

	for {
		accounts := GetAccountsBatch(batchSize, offset)
		if len(accounts) == 0 {
			break
		}

		for _, account := range accounts {
			record := []string{
				account.ID,
				account.Name,
				fmt.Sprintf("%.2f", account.Balance),
				account.CurrencyCode,
				account.Created.Format("2006-01-02"),
			}
			if err := writer.Write(record); err != nil {
				return fmt.Errorf("failed to write record %s: %w", account.ID, err)
			}
		}

		writer.Flush()
		if err := writer.Error(); err != nil {
			return fmt.Errorf("error flushing to file: %w", err)
		}

		offset += batchSize
	}
	return nil
}

// ExportGoalsToCSV exports all goals to CSV
func ExportGoalsToCSV(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{"ID", "Name", "Description", "TargetAmount", "CurrentAmount", "Currency", "Deadline", "HasDeadline", "Category", "Priority", "Status", "LinkedAccounts", "Created", "CompletedDate"}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("failed to write CSV header: %w", err)
	}

	const batchSize = 500
	offset := 0

	for {
		goals := GetGoalsBatch(batchSize, offset)
		if len(goals) == 0 {
			break
		}

		for _, goal := range goals {
			record := []string{
				goal.ID,
				goal.Name,
				goal.Description,
				fmt.Sprintf("%.2f", goal.TargetAmount),
				fmt.Sprintf("%.2f", goal.CurrentAmount),
				goal.CurrencyCode,
				goal.Deadline.Format("2006-01-02"),
				fmt.Sprintf("%v", goal.HasDeadline),
				goal.Category,
				goal.Priority,
				goal.Status,
				goal.LinkedAccountID,
				goal.Created.Format("2006-01-02"),
				goal.CompletedDate.Format("2006-01-02"),
			}
			if err := writer.Write(record); err != nil {
				return fmt.Errorf("failed to write record %s: %w", goal.ID, err)
			}
		}

		writer.Flush()
		if err := writer.Error(); err != nil {
			return fmt.Errorf("error flushing to file: %w", err)
		}

		offset += batchSize
	}
	return nil
}

// Export account transactions to CSV
func ExportAccountTransactionsToCSV(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{"ID", "Account ID", "Amount", "Date", "Note"}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("failed to write CSV header: %w", err)
	}

	const batchSize = 500
	offset := 0

	for {
		accountTransactions := GetAccountTransactionsBatch(batchSize, offset)
		if len(accountTransactions) == 0 {
			break
		}

		for _, accTx := range accountTransactions {

			record := []string{
				accTx.ID,
				accTx.AccountID,
				fmt.Sprintf("%.2f", accTx.Amount),
				accTx.Date.Format("2006-01-02"),
				accTx.Note,
			}
			if err := writer.Write(record); err != nil {
				return fmt.Errorf("failed to write record %s: %w", accTx.ID, err)
			}
		}

		writer.Flush()
		if err := writer.Error(); err != nil {
			return fmt.Errorf("error flushing to file: %w", err)
		}

		offset += batchSize
	}
	return nil
}

// Import transactions from CSV
func ImportTransactionsFromCSV(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// Read header
	header, err := reader.Read()
	if err != nil {
		return err
	}

	headerLenght := len(header)

	// Read data
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	for _, record := range records {
		if len(record) != headerLenght {
			continue // Skip invalid records
		}

		// Parse record
		amount, err := strconv.ParseFloat(record[3], 64)
		if err != nil {
			fmt.Printf("Warning: Skipping invalid amount in row: %v\n", record)
			continue
		}

		date, err := utils.ParseDate(record[1])
		if err != nil {
			fmt.Printf("Warning: Skipping invalid date in row: %v\n", record)
			continue
		}

		// Create transaction
		transaction := models.Transaction{
			ID:           record[0],
			Date:         date,
			Type:         record[2],
			Amount:       amount,
			CurrencyCode: record[4],
			Category:     record[5],
			Description:  record[6],
		}

		err = AddTransaction(transaction)
		if err != nil {
			fmt.Printf("Warning: Skipping, transaction is read but could not added to db: %v\n", record)
			continue
		}
	}

	return nil
}

// Import accounts from CSV
func ImportAccountsFromCSV(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// Read header
	header, err := reader.Read()
	if err != nil {
		return err
	}
	headerLenght := len(header)

	// Read data
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	for _, record := range records {
		if len(record) != headerLenght {
			continue
		}

		balance, err := strconv.ParseFloat(record[2], 64)
		if err != nil {
			fmt.Printf("Warning: Skipping invalid balance in row: %v\n", record)
			continue
		}

		created, err := utils.ParseDate(record[4])
		if err != nil {
			fmt.Printf("Warning: Skipping invalid date in row: %v\n", record)
			continue
		}

		account := models.Account{
			ID:           record[0],
			Name:         record[1],
			Balance:      balance,
			CurrencyCode: record[3],
			Created:      created,
		}

		err = AddAccount(account)
		if err != nil {
			fmt.Printf("Warning: Skipping, account is read but could not added to db: %v\n", record)
			continue
		}
	}

	return nil
}
