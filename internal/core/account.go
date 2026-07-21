package core

import (
	"financial_tracker/internal/models"
	"financial_tracker/internal/storage"
	"financial_tracker/internal/utils"
	"time"
)

// AddAccount adds account
func AddAccount(account models.Account) error {
	_, err := storage.Store.InsertAccount(account)
	if err != nil {
		return err
	}
	return nil
}

// FindAccount finds an account by ID
func FindAccount(id string) *models.Account {
	account, err := storage.Store.GetAccountByID(id)
	if err != nil {
		return nil
	}
	return account
}

// GetAccountsLength returns total number of accounts
func GetAccountsLength() int {
	length, err := storage.Store.GetAccountsLength()
	if err != nil {
		return 0
	}
	return length
}

// GetAccountsBatch returns specific page/batch of accounts
func GetAccountsBatch(batchsize, offset int) []models.Account {
	accounts, err := storage.Store.GetAccountsPaginated(batchsize, offset)
	if err != nil {
		return []models.Account{}
	}
	return accounts
}

// GetRecentAccounts returns N recent accounts
func GetRecentAccounts(limit int) []models.Account {
	accounts, err := storage.Store.GetRecentAccounts(limit)
	if err != nil {
		return []models.Account{}
	}
	return accounts
}

// GetTotalAccountBalance returns total balance across all accounts
func GetTotalAccountsBalanceByCurrency() map[string]float64 {
	balance, err := storage.Store.GetTotalAccountsBalanceByCurrency()
	if err != nil {
		return nil
	}
	return balance
}

// SearchAccountsByName searchs accounts by name
func GetAccountsByName(name string) *models.Account {
	if name == "" {
		return nil
	}

	results, err := storage.Store.GetAccountByName(name)
	if err != nil {
		return nil
	}
	return results
}

// UpdateAccount updates an account values.
func UpdateAccount(account models.Account) error {
	err := storage.Store.UpdateAccount(account)
	if err != nil {
		return err
	}
	return nil
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

// DeleteAccount deletes an account
func DeleteAccount(id string) error {
	return storage.Store.DeleteAccount(id)
}

// GetOneAccountTransactions returns transactions for an account
func GetOneAccountTransactions(accoundID string) []models.AccountTransaction {
	transactions, err := storage.Store.GetOneAccountTransactions(accoundID)
	if err != nil {
		return []models.AccountTransaction{}
	}
	return transactions
}

// GetAccountsTransactionsLength return total number of account transactions
func GetAccountsTransactionsLength() int {
	length, err := storage.Store.GetAccountsTransactionsLength()
	if err != nil {
		return 0
	}
	return length
}

// GetAccountTransactionsBatch returns a specific page of account transactions
func GetAccountTransactionsBatch(batchsize, offset int) []models.AccountTransaction {
	transactions, err := storage.Store.GetAccountTransactionsPaginated(batchsize, offset)
	if err != nil {
		return []models.AccountTransaction{}
	}
	return transactions
}
