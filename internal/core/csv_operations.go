package core

import (
	"encoding/csv"
	"financial_tracker/internal/models"
	"financial_tracker/internal/storage"
	"financial_tracker/internal/utils"
	"fmt"
	"os"
	"strconv"
)

// Export transactions to CSV
func ExportTransactionsToCSV(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{"ID", "Date", "Type", "Amount", "Category", "Description"}
	if err := writer.Write(header); err != nil {
		return err
	}

	// Write data
	for _, transaction := range storage.Transactions {
		record := []string{
			strconv.Itoa(transaction.ID),
			transaction.Date.Format("2006-01-02"),
			transaction.Type,
			fmt.Sprintf("%.2f", transaction.Amount),
			transaction.Category,
			transaction.Description,
		}
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	return nil
}

// Export accounts to CSV
func ExportAccountsToCSV(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{"ID", "Name", "Balance", "Created"}
	if err := writer.Write(header); err != nil {
		return err
	}

	// Write data
	for _, account := range storage.Accounts {
		record := []string{
			strconv.Itoa(account.ID),
			account.Name,
			fmt.Sprintf("%.2f", account.Balance),
			account.Created.Format("2006-01-02"),
		}
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	return nil
}

// Export account transactions to CSV
func ExportAccountTransactionsToCSV(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{"ID", "Account ID", "Account Name", "Amount", "Date", "Note"}
	if err := writer.Write(header); err != nil {
		return err
	}

	// Write data
	for _, at := range storage.AccountTransactions {
		// Find account name
		accountName := "Unknown"
		for _, acc := range storage.Accounts {
			if acc.ID == at.AccountID {
				accountName = acc.Name
				break
			}
		}

		record := []string{
			strconv.Itoa(at.ID),
			strconv.Itoa(at.AccountID),
			accountName,
			fmt.Sprintf("%.2f", at.Amount),
			at.Date.Format("2006-01-02"),
			at.Note,
		}
		if err := writer.Write(record); err != nil {
			return err
		}
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
	_, err = reader.Read()
	if err != nil {
		return err
	}

	// Read data
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	for _, record := range records {
		if len(record) < 6 {
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
			ID:          storage.NextTransactionID,
			Date:        date,
			Type:        record[2],
			Amount:      amount,
			Category:    record[4],
			Description: record[5],
		}

		storage.Transactions = append(storage.Transactions, transaction)
		storage.NextTransactionID++
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
	_, err = reader.Read()
	if err != nil {
		return err
	}

	// Read data
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	for _, record := range records {
		if len(record) < 4 {
			continue
		}

		balance, err := strconv.ParseFloat(record[2], 64)
		if err != nil {
			fmt.Printf("Warning: Skipping invalid balance in row: %v\n", record)
			continue
		}

		created, err := utils.ParseDate(record[3])
		if err != nil {
			fmt.Printf("Warning: Skipping invalid date in row: %v\n", record)
			continue
		}

		account := models.Account{
			ID:      storage.NextAccountID,
			Name:    record[1],
			Balance: balance,
			Created: created,
		}

		storage.Accounts = append(storage.Accounts, account)
		storage.NextAccountID++
	}

	return nil
}
