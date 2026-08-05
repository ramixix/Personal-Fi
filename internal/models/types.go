package models

import (
	"time"
)

type TransactionType string
type GoalStatus string
type GoalPriority string

const (
	Income  TransactionType = "income"
	Expense TransactionType = "expense"
)

const (
	StatusActive    GoalStatus = "active"
	StatusCompleted GoalStatus = "completed"
	StatusPaused    GoalStatus = "paused"
)

const (
	HighPriority   GoalPriority = "high"
	MediumPriority GoalPriority = "medium"
	LowPriority    GoalPriority = "low"
)

// Transaction represents a financial transaction (income or expense)
type Transaction struct {
	ID           string    `json:"id"`
	Date         time.Time `json:"date"`
	Amount       float64   `json:"amount"`
	Category     string    `json:"category"`
	Description  string    `json:"description"`
	Type         string    `json:"type"` // "income" or "expense"
	CurrencyCode string    `json:"currency_code"`
}

// Account represents a savings/purpose account
type Account struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Balance      float64   `json:"balance"`
	CurrencyCode string    `json:"currency_code"`
	Created      time.Time `json:"created"`
}

// AccountTransaction represents money added/removed from an account
type AccountTransaction struct {
	ID        string    `json:"id"`
	AccountID string    `json:"account_id"`
	Amount    float64   `json:"amount"`
	Date      time.Time `json:"date"`
	Note      string    `json:"note"`
	Automatic bool      `json:"automatic"`
}

// Goal represents a financial goal
type Goal struct {
	ID              string    `json:"id"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	TargetAmount    float64   `json:"target_amount"`
	CurrentAmount   float64   `json:"current_amount"`
	CurrencyCode    string    `json:"currency_code"`
	Deadline        time.Time `json:"deadline"`
	HasDeadline     bool      `json:"has_deadline"`
	Category        string    `json:"category"`
	Priority        string    `json:"priority"`
	Status          string    `json:"status"`
	LinkedAccountID string    `json:"linked_account_id"`
	Created         time.Time `json:"created"`
	CompletedDate   time.Time `json:"completed_date"`
}

// GoalContribution represents a contribution toward a goal
type GoalContribution struct {
	ID        string    `json:"id"`
	GoalID    string    `json:"goal_id"`
	Amount    float64   `json:"amount"`
	Date      time.Time `json:"date"`
	Note      string    `json:"note"`
	Automatic bool      `json:"automatic"`
}

// CurrencyTotal holds income and expenses for a specific currency
type CurrencyTotal struct {
	Income   float64
	Expenses float64
}

// CategorySummary holds the aggregated data for a single category
type CategorySummary struct {
	Category     string
	CurrencyCode string
	Count        int
	TotalIncome  float64
	TotalExpense float64
}

// =================================================
// SearchCriteria holds multiple search parameters. Used in search.go files mostly
// =================================================
type SearchCriteria struct {
	Keyword         string
	Categories      []string
	TransactionType string // "income", "expense", or "" for both
	MinAmount       float64
	MaxAmount       float64
	StartDate       time.Time
	EndDate         time.Time
	HasDateRange    bool
}

// ====================================================================
//
//	Report and Analytics Related Structs (Mostly used in analytics.go)
//
// ===================================================================
// MonthlyReport represents financial data for a specific month
type MonthlyReport struct {
	Year         int // because time.Date.Year() function returns int
	Month        time.Month
	Income       float64
	Expenses     float64
	Net          float64
	TxCount      int
	CurrencyCode string
}

// CategoryReport represents spending/income for a category
// type CategoryReport struct {
// 	Category     string
// 	Amount       float64
// 	Count        int
// 	Percent      float64
// 	CurrencyCode string
// }

// ComparisonReport represents comparison between two periods
type ComparisonReport struct {
	Period1Income   float64
	Period2Income   float64
	Period1Expenses float64
	Period2Expenses float64
	IncomeChange    float64
	ExpenseChange   float64
	IncomePercent   float64
	ExpensePercent  float64
	CurrencyCode    string
}
