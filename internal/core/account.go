package core

import (
	"financial_tracker/internal/models"
	"financial_tracker/internal/storage"
	"financial_tracker/internal/utils"
	"time"
)

// DeleteAccount deletes an account
func DeleteAccount(id string) error {
	return storage.Store.DeleteAccount(id)
}

// AddAccount adds account
func AddAccount(account models.Account) (id string) {
	id, err := storage.Store.InsertAccount(account)
	if err != nil {
		return ""
	}
	return id
}

// FindAccount finds an account by ID
func FindAccount(id string) *models.Account {
	account, err := storage.Store.GetAccountByID(id)
	if err != nil {
		return nil
	}
	return account
}

func GetAllAccounts() []models.Account {
	accounts, err := storage.Store.GetAllAccounts()
	if err != nil {
		return []models.Account{}
	}
	return accounts
}

// AddMoneyToAccount adds money to an account as transaction
func AddMoneyToAccount(accountID string, amount float64, note string) bool {
	tx := models.AccountTransaction{
		ID:        utils.MustGenerateUUID(),
		AccountID: accountID,
		Amount:    amount,
		Date:      time.Now(),
		Note:      note,
		Automatic: false,
	}
	_, err := storage.Store.InsertAccountTransaction(tx)
	return err == nil
}

// GetTotalAccountBalance returns total balance across all accounts
func GetTotalAccountBalance(accountID string) float64 {
	balance, err := storage.Store.GetTotalAccountsBalance()
	if err != nil {
		return 0
	}
	return balance
}

// GetAccountTransactions returns transactions for an account
func GetAccountTransactions(accoundID string) []models.AccountTransaction {
	transactions, err := storage.Store.GetAccountTransactions(accoundID)
	if err != nil {
		return []models.AccountTransaction{}
	}
	return transactions
}
