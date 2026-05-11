package database

import (
	"database/sql"
	"financial_tracker/internal/models"
	"financial_tracker/internal/utils"
	"fmt"
	"time"
)

//=============================================== Create(C) =================================================

// InsertAccount creates a new account
func (s *Store) InsertAccount(account models.Account) (string, error) {
	// Generate UUID if not provided
	if account.ID == "" {
		account.ID = utils.MustGenerateUUID()
	}

	// Set default currency if not provided
	if account.CurrencyCode == "" {
		account.CurrencyCode = "USD"
	}

	// Set created time if not provided
	if account.Created.IsZero() {
		account.Created = time.Now()
	}

	query := `
		INSERT INTO accounts (id, name, balance, currency_code, created)
		VALUES (?, ?, ?, ?, ?)
	`

	_, err := s.db.Exec(query, account.ID, account.Name, account.Balance, account.CurrencyCode, account.Created)
	if err != nil {
		s.logger.Error("Failed to insert account", "error", err, "id", account.ID)
		return "", fmt.Errorf("failed to insert account: %w", err)
	}

	s.logger.Info("Account created", "id", account.ID, "name", account.Name)
	return account.ID, nil
}

//=============================================== Read(R) =================================================

// GetAllAccounts retrieves all accounts
func (s *Store) GetAllAccounts() ([]models.Account, error) {
	query := `
		SELECT id, name, balance, currency_code, created
		FROM accounts
		ORDERED BY created DESC
	`

	rows, err := s.db.Query(query)
	if err != nil {
		s.logger.Error("Failed to query accounts(GetAllAccounts)", "error", err)
		return nil, fmt.Errorf("failed to read all account: %w", err)
	}
	defer rows.Close()

	return s.scanAccounts(rows)
}

func (s *Store) scanAccounts(rows *sql.Rows) ([]models.Account, error) {
	var accounts []models.Account

	for rows.Next() {
		var account models.Account
		err := rows.Scan(&account.ID, &account.Name, &account.Balance, &account.CurrencyCode, &account.Created)
		if err != nil {
			s.logger.Error("Failed to scan account", "error", err)
			continue
		}

		accounts = append(accounts, account)
	}
	return accounts, nil
}

// GetAccountByID retrieves a single account by ID
func (s *Store) GetAccountByID(id string) (*models.Account, error) {
	query := `
		SELECT id, name, balance, currency_code, created 
		FROM accounts
		WHERE id = ?
	`
	var account models.Account
	err := s.db.QueryRow(query, id).Scan(&account.ID, account.Name, account.Balance, account.CurrencyCode, account.Created)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("account not found: %s", id)
	}

	if err != nil {
		s.logger.Error("Failed to get account", "error", err, "id", id)
		return nil, fmt.Errorf("failed to get account: %w", err)
	}

	return &account, nil
}

// GetAccountByName retrieves account by name
func (s *Store) GetAccountByName(name string) (*models.Account, error) {
	query := `
		SELECT id, name, balance, currency_code, created 
		FROM accounts
		WHERE name = ?
	`

	var account models.Account
	err := s.db.QueryRow(query, name).Scan(&account.ID, account.Name, account.Balance, account.CurrencyCode, account.Created)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("account not found: %s", name)
	}

	if err != nil {
		s.logger.Error("Failed to get account", "error", err, "name", name)
		return nil, fmt.Errorf("failed to get account: %w", err)
	}

	return &account, nil
}

// SearchAccounts searches accounts by keyword in name
func (s *Store) SearchAccounts(keyword string) ([]models.Account, error) {
	query := `
		SELECT id, name, balance, currency_code, created
		FROM accounts
		WHERE name LIKE ?
		ORDER BY name
	`
	searchPattern := "%" + keyword + "%"

	rows, err := s.db.Query(query, searchPattern)
	if err != nil {
		return nil, fmt.Errorf("failed to get accounts with given keyword(%s): %w", keyword, err)
	}
	rows.Close()

	return s.scanAccounts(rows)
}

//=============================================== Update(U) =================================================

// UpdateAccount updates an existing account
func (s *Store) UpdateAccount(newAccountValue models.Account) error {
	query := `
		UPDATE accounts
		SET name = ?, balance = ?, currency_code = ?, created = ?
		WHERE id = ?
	`

	result, err := s.db.Exec(query, newAccountValue.Name, newAccountValue.Balance, newAccountValue.CurrencyCode, newAccountValue.Created, newAccountValue.ID)
	if err != nil {
		s.logger.Error("Failed to update account", "error", err, "id", newAccountValue.ID)
		return fmt.Errorf("failed to update account: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("account not found: %s", newAccountValue.ID)
	}

	s.logger.Info("Account updated", "id", newAccountValue.ID, "name", newAccountValue.Name)
	return nil
}

// UpdateAccountBalance updates only the balance (for deposits/withdrawals)
func (s *Store) UpdateAccountBalance(id string, newBalance float64) error {
	query := `
		UPDATE accounts
		SET balance = ?
		WHERE id = ?
	`

	result, err := s.db.Exec(query, newBalance, id)
	if err != nil {
		s.logger.Error("Failed to Update balance", "error", err, "ID", id)
		return fmt.Errorf("failed to update account balance: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("account not found: %s", id)
	}

	s.logger.Info("Account updated", "id", id, "new_balance", newBalance)
	return nil
}

// AddToAccountBalance adds an amount to account balance (atomic operation)
func (s *Store) AddToAccountBalance(id string, amount float64) error {
	query := `
		UPDATE accounts
		SET balance = balance + ?
		WHERE id = ?
	`

	result, err := s.db.Exec(query, amount, id)
	if err != nil {
		s.logger.Error("Failed to add amount to account balance", "error", err, "ID", id)
		return fmt.Errorf("failed to add amount to %s account balance: %w", id, err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("account not found : %s", id)
	}

	s.logger.Info("Amount added to account", "id", id, "amount", amount)
	return nil
}

//=============================================== DELETE(D) =================================================

// DeleteAccount deletes an account (and cascades to account_transactions)
func (s *Store) DeleteAccount(id string) error {
	query := `DELETE FROM accounts WHERE id = ?`

	result, err := s.db.Exec(query, id)
	if err != nil {
		s.logger.Error("Failed to delete account", "error", err, "ID", id)
		return fmt.Errorf("failed to delete account %s: %w", id, err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("account not found: %s", id)
	}

	s.logger.Info("Account deleted", "id", id)
	return nil
}

//=============================================== ANALYTICS =================================================

// GetTotalAccountsBalance calculates total balance across all accounts
func (s *Store) GetTotalAccountsBalance() (float64, error) {
	query := `Select COALESCE(sum(balance), 0) FROM accounts`

	var totalBalance float64
	err := s.db.QueryRow(query).Scan(&totalBalance)
	if err != nil {
		return 0, fmt.Errorf("failed to get total balance: %w", err)
	}

	return totalBalance, nil
}

// GetTotalAccountsBalanceByCurrency calculates total balance per currency
func (s *Store) GetTotalAccountsBalanceByCurrency() (map[string]float64, error) {
	query := `
		SELECT currency_code, SUM(balance)
		FROM accounts
		GROUP BY currency_code
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get balance by currency: %w", err)
	}
	defer rows.Close()

	var balances = make(map[string]float64)
	for rows.Next() {
		var currency_code string
		var amount float64
		err := rows.Scan(&currency_code, &amount)
		if err != nil {
			continue
		}
		balances[currency_code] = amount
	}

	return balances, nil
}

// GetAccountCount returns the number of accounts
func (s *Store) GetAccountCount() (int, error) {
	query := `SELECT COUNT(*) FROM accounts`

	var count int
	err := s.db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get account count: %w", err)
	}

	return count, nil
}

//=============================================== ACCOUNT TRANSACTIONS =================================================
//======================================================================================================================
//======================================================================================================================

//=============================================== CREATE(C) =================================================

// InsertAccountTransaction creates a new account transaction and updates balance atomically
func (s *Store) InsertAccountTransaction(tx models.AccountTransaction) (string, error) {
	// Generate UUID if not provided
	if tx.ID == "" {
		tx.ID = utils.MustGenerateUUID()
	}

	// Set date if not provided
	if tx.Date.IsZero() {
		tx.Date = time.Now()
	}

	// Use database transaction for atomicity
	dbTx, err := s.db.Begin()
	if err != nil {
		return "", fmt.Errorf("failed to begin transaction: %w", err)

	}
	defer dbTx.Rollback()

	query := `
		INSERT INTO account_transactions (id, account_id, amount, date, note, automatic)
		VALUES(?, ?, ?, ?, ?, ?)`

	_, err = s.db.Exec(query, tx.ID, tx.AccountID, tx.Amount, tx.Date, tx.Note, tx.Automatic)
	if err != nil {
		s.logger.Error("Failed to insert account transaction", "error", err, "id", tx.ID)
		return "", fmt.Errorf("failed to insert account transaction: %w", err)
	}

	// Update account balance
	updateBalanceQuery := `UPDATE accounts SET balance = balance + ? WHERE id = ?`

	result, err := s.db.Exec(updateBalanceQuery, tx.Amount, tx.AccountID)
	if err != nil {
		s.logger.Error("Failed to update account balance", "error", err, "account_id", tx.AccountID)
		return "", fmt.Errorf("failed to update account balance: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return "", fmt.Errorf("account not found: %s", tx.AccountID)
	}

	// Commit transaction
	err = dbTx.Commit()
	if err != nil {
		return "", fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("Account transaction created", "id", tx.ID, "account_id", tx.AccountID, "amount", tx.Amount)
	return tx.ID, nil
}

//=============================================== Read(R) =================================================

// GetAccountTransactions retrieves all transactions for an account
func (s *Store) GetAccountTransactions(accountID string) ([]models.AccountTransaction, error) {
	query := `
		SELECT id, account_id, amount, date, note, automatic
		FROM account_transactions
		WHERE account_id = ?
		ORDER BY date DESC`

	rows, err := s.db.Query(query, accountID)
	if err != nil {
		// s.logger.Error("Failed to find transactions for given account", "error", err, "Account ID", accountID)
		return nil, fmt.Errorf("failed to find transactions for given account(%s) : %w", accountID, err)
	}
	defer rows.Close()

	return s.scanAccountTransactions(rows)
}

// scanAccountTransactions is a helper to scan account transaction rows
func (s *Store) scanAccountTransactions(rows *sql.Rows) ([]models.AccountTransaction, error) {
	var accountTransactions []models.AccountTransaction
	for rows.Next() {
		var tx models.AccountTransaction
		err := rows.Scan(&tx.ID, &tx.AccountID, &tx.Amount, &tx.Date, &tx.Note, &tx.Automatic)
		if err != nil {
			s.logger.Error("Failed to scan account transaction", "error", err)
			continue
		}
		accountTransactions = append(accountTransactions, tx)
	}

	return accountTransactions, nil
}

// GetAllAccountTransactions retrieves all transactions for all accounts
func (s *Store) GetAllAccountTransactions() ([]models.AccountTransaction, error) {
	query := `
		SELECT id, account_id, amount, date, note, automatic
		FROM account_transactions
		ORDER BY date DESC`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to find transactions for any account : %w", err)
	}
	defer rows.Close()

	return s.scanAccountTransactions(rows)
}

// GetAccountTransactionByID retrieves a single account transaction
func (s *Store) GetAccountTransactionByID(txID string) (*models.AccountTransaction, error) {
	query := `
		SELECT id, account_id, amount, date, note, automatic
		FROM account_transactions
		WHERE id = ?`

	var tx models.AccountTransaction
	err := s.db.QueryRow(query, txID).Scan(&tx.ID, &tx.AccountID, &tx.Amount, &tx.Date, &tx.Note, &tx.Automatic)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("account transaction not found: %s", txID)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get account transaction: %w", err)
	}

	return &tx, nil

}

//=============================================== Update(U) =================================================

// UpdateAccountTransaction updates an account transaction and adjusts balance
func (s *Store) UpdateAccountTransaction(newTxValue models.AccountTransaction) error {
	// Get transaction first to know the amount
	transac, err := s.GetAccountTransactionByID(newTxValue.ID)
	if err != nil {
		return err
	}

	dbTx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction(DeleteAccountTransaction): %w", err)
	}
	defer dbTx.Rollback()

	query := `
		UPDATE account_transactions
		SET account_id = ?, amount = ?, date = ?, note = ?, automatic = ?
		WHERE id = ?
	`

	_, err = s.db.Exec(query, newTxValue.AccountID, newTxValue.Amount, newTxValue.Date, newTxValue.Note, newTxValue.Automatic, newTxValue.ID)
	if err != nil {
		return fmt.Errorf("failed to update account transaction: %w", err)
	}

	balanceDifference := newTxValue.Amount - transac.Amount
	updateQuery := `
		UPDATE accounts
		SET balance = balance + ?
		WHERE id = ?`
	_, err = s.db.Exec(updateQuery, balanceDifference, transac.AccountID)
	if err != nil {
		return fmt.Errorf("failed to update account balance: %w", err)
	}

	// Commit transaction
	if err := dbTx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("Account transaction updated", "id", newTxValue.ID, "account_id", transac.AccountID)
	return nil
}

//=============================================== DELETE(D) =================================================

// DeleteAccountTransaction deletes an account transaction and adjusts balance
func (s *Store) DeleteAccountTransaction(txID string) error {
	// Get transaction first to know the amount
	transac, err := s.GetAccountTransactionByID(txID)
	if err != nil {
		return err
	}

	// Use database transaction for atomicity
	dbTx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction(DeleteAccountTransaction): %w", err)
	}
	defer dbTx.Rollback()

	query := `DELETE FROM account_transactions WHERE id = ?`
	_, err = s.db.Exec(query, txID)
	if err != nil {
		return fmt.Errorf("failed to delete account transaction: %w", err)
	}

	updateQuery := `
		UPDATE accounts
		SET balance = balance - ?
		WHERE id = ?`
	_, err = s.db.Exec(updateQuery, transac.Amount, transac.AccountID)
	if err != nil {
		return fmt.Errorf("failed to update account balance: %w", err)
	}

	// Commit transaction
	if err := dbTx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("Account transaction deleted", "id", txID, "account_id", transac.AccountID)
	return nil
}
