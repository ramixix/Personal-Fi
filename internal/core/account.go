package core

import (
	"financial_tracker/internal/models"
	"financial_tracker/internal/storage"
	"time"
)

func DeleteAccount(accountId int) bool {
	for index, ac := range storage.Accounts {
		if ac.ID == accountId {
			var filteredAccountTransactions []models.AccountTransaction
			for _, account_transac := range storage.AccountTransactions {
				if account_transac.AccountID != accountId {
					// get only transactions that are not part of account id we want to remove
					filteredAccountTransactions = append(filteredAccountTransactions, account_transac)
				}
			}
			// set new account transaction
			storage.AccountTransactions = filteredAccountTransactions
			// remove account
			storage.Accounts = append(storage.Accounts[:index], storage.Accounts[index+1:]...)
			return true
		}
	}
	return false
}

// Find account by ID
func FindAccount(id int) *models.Account {
	for i := range storage.Accounts {
		if storage.Accounts[i].ID == id {
			return &storage.Accounts[i]
		}
	}
	return nil
}

// Find account index by ID
func FindAccountIndex(id int) int {
	for i, account := range storage.Accounts {
		if account.ID == id {
			return i
		}
	}
	return -1
}

// Add money to account and record transaction
func AddMoneyToAccount(account *models.Account, amount float64, note string) {
	account.Balance += amount

	newAccountTransaction := models.AccountTransaction{
		ID:        storage.NextAccountTransactionID,
		AccountID: account.ID,
		Amount:    amount,
		Date:      time.Now(),
		Note:      note,
	}

	storage.AccountTransactions = append(storage.AccountTransactions, newAccountTransaction)
	storage.NextAccountTransactionID++
}

// Calculate total balance across all accounts
func GetTotalAccountBalance() float64 {
	var total float64
	for _, account := range storage.Accounts {
		total += account.Balance
	}
	return total
}
