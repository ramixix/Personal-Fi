package main

import "time"

func deleteAccount(accountId int) bool {
	for index, ac := range accounts {
		if ac.ID == accountId {
			var filteredAccountTransactions []AccountTransaction
			for _, account_transac := range accountTransactions {
				if account_transac.AccountID != accountId {
					// get only transactions that are not part of account id we want to remove
					filteredAccountTransactions = append(filteredAccountTransactions, account_transac)
				}
			}
			// set new account transaction
			accountTransactions = filteredAccountTransactions
			// remove account
			accounts = append(accounts[:index], accounts[index+1:]...)
			return true
		}
	}
	return false
}

// Find account by ID
func findAccount(id int) *Account {
	for i := range accounts {
		if accounts[i].ID == id {
			return &accounts[i]
		}
	}
	return nil
}

// Find account index by ID
func findAccountIndex(id int) int {
	for i, account := range accounts {
		if account.ID == id {
			return i
		}
	}
	return -1
}

// Add money to account and record transaction
func addMoneyToAccount(account *Account, amount float64, note string) {
	account.Balance += amount

	newAccountTransaction := AccountTransaction{
		ID:        nextAccountTransactionID,
		AccountID: account.ID,
		Amount:    amount,
		Date:      time.Now(),
		Note:      note,
	}

	accountTransactions = append(accountTransactions, newAccountTransaction)
	nextAccountTransactionID++
}

// Calculate total balance across all accounts
func getTotalAccountBalance() float64 {
	var total float64
	for _, account := range accounts {
		total += account.Balance
	}
	return total
}
