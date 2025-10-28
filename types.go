package main

import (
	"time"
)

// Transaction represents a financial transaction (income or expense)
type Transaction struct {
	ID          int       `json:"id"`
	Date        time.Time `json:"date"`
	Amount      float64   `json:"amount"`
	Category    string    `json:"category"`
	Description string    `json:"description"`
	Type        string    `json:"type"` // "income" or "expense"
}

// Account represents a savings/purpose account
type Account struct {
	ID      int       `json:"id"`
	Name    string    `json:"name"`
	Balance float64   `json:"balance"`
	Created time.Time `json:"created"`
}

// AccountTransaction represents money added/removed from an account
type AccountTransaction struct {
	ID        int       `json:"id"`
	AccountID int       `json:"account_id"`
	Amount    float64   `json:"amount"`
	Date      time.Time `json:"date"`
	Note      string    `json:"note"`
}

// AppData structure to hold all app data for persistence
type AppData struct {
	Transactions             []Transaction        `json:"transactions"`
	Accounts                 []Account            `json:"accounts"`
	AccountTransactions      []AccountTransaction `json:"account_transactions"`
	NextTransactionID        int                  `json:"next_transaction_id"`
	NextAccountID            int                  `json:"next_account_id"`
	NextAccountTransactionID int                  `json:"next_account_transaction_id"`
}
