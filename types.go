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
	Goals                    []Goal               `json:"goals"`
	GoalContributions        []GoalContribution   `json:"goal_contributions"`
	NextTransactionID        int                  `json:"next_transaction_id"`
	NextAccountID            int                  `json:"next_account_id"`
	NextAccountTransactionID int                  `json:"next_account_transaction_id"`
	NextGoalID               int                  `json:"next_goal_id"`
	NextGoalContributionID   int                  `json:"next_goal_contribution_id"`
}

// Goal represents a financial goal
type Goal struct {
	ID              int       `json:"id"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	TargetAmount    float64   `json:"target_amount"`
	CurrentAmount   float64   `json:"current_amount"`
	Deadline        time.Time `json:"deadline"`          // Optional
	HasDeadline     bool      `json:"has_deadline"`      // Whether deadline is set
	Category        string    `json:"category"`          // "savings", "debt", "investment", etc.
	Priority        string    `json:"priority"`          // "high", "medium", "low"
	Status          string    `json:"status"`            // "active", "completed", "paused", "cancelled"
	LinkedAccountID int       `json:"linked_account_id"` // Optional - 0 if not linked
	Created         time.Time `json:"created"`
	CompletedDate   time.Time `json:"completed_date"` // When goal was completed
}

// GoalContribution represents a contribution toward a goal
type GoalContribution struct {
	ID        int       `json:"id"`
	GoalID    int       `json:"goal_id"`
	Amount    float64   `json:"amount"`
	Date      time.Time `json:"date"`
	Note      string    `json:"note"`
	Automatic bool      `json:"automatic"` // Was it automatically tracked or manual?
}
