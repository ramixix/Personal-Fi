package core

import (
	"encoding/csv"
	"financial_tracker/internal/models"
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
	header := []string{"ID", "Date", "Type", "Amount", "Currency", "Category", "Description"}
	if err := writer.Write(header); err != nil {
		return err
	}

	// Write data
	transactions := GetAllTransactions()
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
	header := []string{"ID", "Name", "Balance", "Currency", "Created"}
	if err := writer.Write(header); err != nil {
		return err
	}

	// Write data
	accounts := GetAllAccounts()
	for _, account := range accounts {
		record := []string{
			account.ID,
			account.Name,
			fmt.Sprintf("%.2f", account.Balance),
			account.CurrencyCode,
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
	header := []string{"ID", "Account ID", "Amount", "Date", "Note"}
	if err := writer.Write(header); err != nil {
		return err
	}

	// Write data
	accountTransactions := GetAllAccountTransactions()
	for _, accTx := range accountTransactions {

		record := []string{
			accTx.ID,
			accTx.AccountID,
			fmt.Sprintf("%.2f", accTx.Amount),
			accTx.Date.Format("2006-01-02"),
			accTx.Note,
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
		if len(record) < headerLenght {
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
		if len(record) < headerLenght {
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
