package database

import (
	"database/sql"
	"financial_tracker/internal/models"
	"financial_tracker/internal/utils"
	"fmt"
	"time"
)

type DBTX interface {
	Exec(query string, args ...any) (sql.Result, error)
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
}

//=============================================== Create(C) =================================================

// InsertGoal creates a new goal
func (s *Store) InsertGoal(goal models.Goal) (string, error) {
	// Generate UUID if not provided
	if goal.ID == "" {
		goal.ID = utils.MustGenerateUUID()
	}

	// Set default currency if not provided
	if goal.CurrencyCode == "" {
		goal.CurrencyCode = "USD"
	}

	// Set created time if not provided
	if goal.Created.IsZero() {
		goal.Created = time.Now()
	}

	// Set default priority if not provided
	if goal.Priority == "" {
		goal.Priority = string(models.MediumPriority)
	}

	// Set default status if not provided
	if goal.Status == "" {
		goal.Status = string(models.StatusActive)
	}

	query := `
        INSERT INTO goals (
            id, name, description, target_amount, current_amount, currency_code,
            deadline, has_deadline, category, priority, status, 
            linked_account_id, created, completed_date
        )
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
    `

	_, err := s.db.Exec(query,
		goal.ID,
		goal.Name,
		goal.Description,
		goal.TargetAmount,
		goal.CurrentAmount,
		goal.CurrencyCode,
		goal.Deadline,
		goal.HasDeadline,
		goal.Category,
		goal.Priority,
		goal.Status,
		goal.LinkedAccountID,
		goal.Created,
		goal.CompletedDate,
	)

	if err != nil {
		s.logger.Error("Failed to insert goal", "error", err, "id", goal.ID)
		return "", fmt.Errorf("failed to insert goal: %w", err)
	}

	s.logger.Info("Goal created", "id", goal.ID, "name", goal.Name, "target", goal.TargetAmount)
	return goal.ID, nil
}

//=============================================== Create(C) =================================================

// GetAllGoals retrieves all goals

func (s *Store) GetGoalsLength(status models.GoalStatus) (int, error) {
	query := `SELECT COUNT(*) FROM goals`
	args := []any{}

	if status != "" && status == models.StatusActive || status == models.StatusCompleted {
		query += ` WHERE status = ?`
		args = append(args, status)
	}

	var count int
	err := s.db.QueryRow(query, args...).Scan(&count)
	if err != nil {
		s.logger.Error("Failed to return Goals Length", "error", err)
		return 0, fmt.Errorf("failed to get Goals Length: %w", err)
	}
	return count, nil
}

// GetGoalsPaginated return a specific "page" of goals
func (s *Store) GetGoalsPaginated(limit, offset int) ([]models.Goal, error) {
	query := `
		SELECT id, name, description, target_amount, current_amount, currency_code,
			   deadline, has_deadline, category, priority, status,
			   linked_account_id, created, completed_date
		FROM goals
		ORDER BY created DESC
		LIMIT ? OFFSET ?
	`

	rows, err := s.db.Query(query, limit, offset)
	if err != nil {
		s.logger.Error("Failed to query goal batch", "error", err)
		return nil, fmt.Errorf("failed to read batch of goals: %w", err)
	}
	defer rows.Close()

	return s.scanGoals(rows)
}

// GetRecentGoals returns N recent goals
func (s *Store) GetRecentGoals(limit int) ([]models.Goal, error) {
	query := `
		SELECT id, name, description, target_amount, current_amount, currency_code,
			   deadline, has_deadline, category, priority, status,
			   linked_account_id, created, completed_date
		FROM goals
		ORDER BY created DESC
		LIMIT ?
	`

	rows, err := s.db.Query(query, limit)
	if err != nil {
		s.logger.Error("Failed to query N recent goals", "error", err)
		return nil, fmt.Errorf("failed to query N recent goals: %w", err)
	}
	defer rows.Close()

	return s.scanGoals(rows)
}

// scanGoals is a helper to scan multiple goal rows
func (s *Store) scanGoals(rows *sql.Rows) ([]models.Goal, error) {
	var goals []models.Goal

	for rows.Next() {
		var goal models.Goal
		err := rows.Scan(
			&goal.ID,
			&goal.Name,
			&goal.Description,
			&goal.TargetAmount,
			&goal.CurrentAmount,
			&goal.CurrencyCode,
			&goal.Deadline,
			&goal.HasDeadline,
			&goal.Category,
			&goal.Priority,
			&goal.Status,
			&goal.LinkedAccountID,
			&goal.Created,
			&goal.CompletedDate,
		)
		if err != nil {
			s.logger.Error("Failed to scan goal", "error", err)
			continue
		}
		goals = append(goals, goal)
	}

	if err := rows.Err(); err != nil {
		s.logger.Error("Error iterating over rows", "error", err)
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return goals, nil
}

// GetGoalByID retrieves a single goal by ID
func (s *Store) GetGoalByID(goalID string) (*models.Goal, error) {
	return s.GetGoalByIDTx(s.db, goalID)
}

// GetGoalByIDTx retrieves a single goal by ID (for a specified db connection)
func (s *Store) GetGoalByIDTx(db DBTX, goalID string) (*models.Goal, error) {
	query := `
		SELECT id, name, description, target_amount, current_amount, currency_code,
			   deadline, has_deadline, category, priority, status,
			   linked_account_id, created, completed_date
		FROM goals
        WHERE id = ?
	`

	var goal models.Goal
	err := db.QueryRow(query, goalID).Scan(
		&goal.ID,
		&goal.Name,
		&goal.Description,
		&goal.TargetAmount,
		&goal.CurrentAmount,
		&goal.CurrencyCode,
		&goal.Deadline,
		&goal.HasDeadline,
		&goal.Category,
		&goal.Priority,
		&goal.Status,
		&goal.LinkedAccountID,
		&goal.Created,
		&goal.CompletedDate,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("goal not found: %s", goalID)
	}

	if err != nil {
		s.logger.Error("Failed to get goal", "error", err, "id", goalID)
		return nil, fmt.Errorf("failed to get goal: %w", err)
	}

	return &goal, nil
}

// GetGoalsByStatus retrieves goals by status (active, completed, paused, cancelled)
func (s *Store) GetGoalsByStatus(status models.GoalStatus) ([]models.Goal, error) {
	query := `
		SELECT id, name, description, target_amount, current_amount, currency_code,
			   deadline, has_deadline, category, priority, status,
			   linked_account_id, created, completed_date
		FROM goals
		WHERE status = ?
		ORDER BY created DESC
	`

	rows, err := s.db.Query(query, status)
	if err != nil {
		return nil, fmt.Errorf("failed to query goals by status: %w", err)
	}
	defer rows.Close()

	return s.scanGoals(rows)
}

// GetActiveGoals retrieves all active goals
func (s *Store) GetActiveGoals() ([]models.Goal, error) {
	return s.GetGoalsByStatus(models.StatusActive)
}

// GetCompletedGoals retrieves all completed goals
func (s *Store) GetCompletedGoals() ([]models.Goal, error) {
	return s.GetGoalsByStatus(models.StatusCompleted)
}

// GetGoalsByPriority retrieves goals by priority
func (s *Store) GetGoalsByPriority(priority models.GoalPriority) ([]models.Goal, error) {
	query := `
		SELECT id, name, description, target_amount, current_amount, currency_code,
			   deadline, has_deadline, category, priority, status,
			   linked_account_id, created, completed_date
		FROM goals
		WHERE priority = ?
		ORDER BY created DESC
	`

	rows, err := s.db.Query(query, priority)
	if err != nil {
		return nil, fmt.Errorf("failed to query goals by priority: %w", err)
	}
	defer rows.Close()

	return s.scanGoals(rows)
}

// GetGoalsByCategory retrieves goals by category
func (s *Store) GetGoalsByCategory(category string) ([]models.Goal, error) {
	query := `
		SELECT id, name, description, target_amount, current_amount, currency_code,
			   deadline, has_deadline, category, priority, status,
			   linked_account_id, created, completed_date
		FROM goals
		WHERE category = ?
		ORDER BY created DESC
	`

	rows, err := s.db.Query(query, category)
	if err != nil {
		return nil, fmt.Errorf("failed to query goals by category: %w", err)
	}
	defer rows.Close()

	return s.scanGoals(rows)
}

// SearchGoals searches goals by keyword in name or description or category
func (s *Store) SearchGoals(keyword string) ([]models.Goal, error) {
	query := `
		SELECT id, name, description, target_amount, current_amount, currency_code,
			   deadline, has_deadline, category, priority, status,
			   linked_account_id, created, completed_date
		FROM goals
		WHERE name LIKE ? OR description LIKE ? OR category LIKE ?
		ORDER BY created DESC
	`

	searchPattern := "%" + keyword + "%"
	rows, err := s.db.Query(query, searchPattern, searchPattern, searchPattern)
	if err != nil {
		return nil, fmt.Errorf("failed to search goals: %w", err)
	}
	defer rows.Close()

	return s.scanGoals(rows)
}

//=============================================== UPDATE(U) =================================================

// UpdateGoal updates an existing goal
func (s *Store) UpdateGoal(newGoalValue models.Goal) error {
	query := `
		UPDATE goals
		SET name = ?, description = ?, target_amount = ?, current_amount = ?,
			currency_code = ?, deadline = ?, has_deadline = ?, category = ?,
			priority = ?, status = ?, linked_account_id = ?, completed_date = ?
		WHERE id = ?
	`

	result, err := s.db.Exec(query,
		newGoalValue.Name,
		newGoalValue.Description,
		newGoalValue.TargetAmount,
		newGoalValue.CurrentAmount,
		newGoalValue.CurrencyCode,
		newGoalValue.Deadline,
		newGoalValue.HasDeadline,
		newGoalValue.Category,
		newGoalValue.Priority,
		newGoalValue.Status,
		newGoalValue.LinkedAccountID,
		newGoalValue.CompletedDate,
		newGoalValue.ID,
	)

	if err != nil {
		s.logger.Error("Failed to update goal", "error", err, "id", newGoalValue.ID)
		return fmt.Errorf("failed to update goal: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("goal not found: %s", newGoalValue.ID)
	}

	s.logger.Info("Goal updated", "id", newGoalValue.ID, "name", newGoalValue.Name)
	return nil
}

// UpdateGoalStatus updates only the status of a goal
func (s *Store) UpdateGoalStatus(goalID string, status string) error {
	return s.UpdateGoalStatusTx(s.db, goalID, status)
}

// UpdateGoalStatusTx updates only the status of a goal (For specified db connection)
func (s *Store) UpdateGoalStatusTx(db DBTX, goalID string, status string) error {
	var query string
	var args []interface{}

	if status == string(models.StatusCompleted) {
		query = `UPDATE goals SET status = ?, completed_date = ? WHERE id = ?`
		args = []interface{}{status, time.Now(), goalID}
	} else {
		query = `UPDATE goals SET status = ? WHERE id = ?`
		args = []interface{}{status, goalID}
	}

	result, err := db.Exec(query, args...)
	if err != nil {
		s.logger.Error("Failed to update status of given goal", "error", err, "id", goalID)
		return fmt.Errorf("failed to chage status of given goal(%s): %w", goalID, err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("foal not found: %s", goalID)
	}

	s.logger.Info("Goal status updated", "id", goalID, "status", status)
	return nil
}

// UpdateGoalCurrentAmount updates the current amount (for contributions)
func (s *Store) UpdateGoalCurrentAmount(goalID string, newAmount float64) error {
	// Get current goal to check if it should be marked as completed
	goal, err := s.GetGoalByID(goalID)
	if err != nil {
		return err
	}

	if newAmount >= goal.TargetAmount && goal.Status == string(models.StatusActive) {
		// Use transaction to update both amount and status
		return s.completeGoalWithAmount(goalID, newAmount)
	}

	query := `UPDATE goals SET current_amount = ? WHERE id = ?`

	result, err := s.db.Exec(query, newAmount, goalID)
	if err != nil {
		s.logger.Error("Failed to update goal amount", "error", err, "id", goalID)
		return fmt.Errorf("failed to update goal amount: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("goal not found: %s", goalID)
	}

	s.logger.Info("Goal amount updated", "id", goalID, "new_amount", newAmount)
	return nil
}

// AddToGoalAmount adds an amount to goal's current amount (atomic)
func (s *Store) AddToGoalAmountTx(tx *sql.Tx, goalID string, amount float64) error {
	// Get current goal to check if it should be marked as completed
	goal, err := s.GetGoalByIDTx(tx, goalID)
	if err != nil {
		return err
	}

	newAmount := goal.CurrentAmount + amount

	// Check if goal is now completed
	if newAmount >= goal.TargetAmount && goal.Status == string(models.StatusActive) {
		// Use transaction to update both amount and status
		return s.completeGoalWithAmountTx(tx, goalID, newAmount)
	}

	// Just update the amount
	query := `UPDATE goals SET current_amount = current_amount + ? WHERE id = ?`
	result, err := tx.Exec(query, amount, goalID)
	if err != nil {
		s.logger.Error("Failed to add to goal amount", "error", err, "id", goalID)
		return fmt.Errorf("failed to add to goal amount: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("goal not found: %s", goalID)
	}

	s.logger.Info("Amount added to goal", "id", goalID, "amount", amount)
	return nil
}

// completeGoalWithAmount marks goal as completed and updates amount atomically
func (s *Store) completeGoalWithAmount(goalID string, newAmount float64) error {
	return s.completeGoalWithAmountTx(s.db, goalID, newAmount)
}

// completeGoalWithAmountTx marks goal as completed and updates amount atomically (For a specified db connection)
func (s *Store) completeGoalWithAmountTx(db DBTX, goalID string, newAmount float64) error {
	query := `
		UPDATE goals 
		SET current_amount = ?, status = ?, completed_date = ?
		WHERE id = ?
	`

	result, err := db.Exec(query, newAmount, models.StatusCompleted, time.Now(), goalID)
	if err != nil {
		return fmt.Errorf("failed to complete goal: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("goal not found: %s", goalID)
	}

	s.logger.Info("Goal completed!", "id", goalID, "final_amount", newAmount)
	return nil
}

//=============================================== DELETE(D) =================================================

// DeleteGoal deletes a goal (and cascades to goal_contributions)
func (s *Store) DeleteGoal(goalID string) error {
	query := `DELETE FROM goals WHERE id = ?`

	result, err := s.db.Exec(query, goalID)
	if err != nil {
		s.logger.Error("Failed to delete goal", "error", err, "id", goalID)
		return fmt.Errorf("failed to delete goal: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("goal not found: %s", goalID)
	}

	s.logger.Info("Goal deleted", "id", goalID)
	return nil
}

//=============================================== ANALYTICS =================================================

// GetTotalGoalsAmountByCurrency calculates total saved across all goals
func (s *Store) GetTotalGoalsAmountByCurrency() (map[string]float64, error) {
	query := `
		SELECT currency_code, COALESCE(SUM(current_amount), 0)
		FROM goals
		GROUP BY currency_code
	`

	rows, err := s.db.Query(query)
	if err != nil {
		s.logger.Error("Failed to query goals balances by currency", "error", err)
		return nil, fmt.Errorf("failed to get goals balances by currency: %w", err)
	}
	defer rows.Close()

	balances := make(map[string]float64)

	for rows.Next() {
		var currencyCode string
		var amount float64

		if err := rows.Scan(&currencyCode, &amount); err != nil {
			s.logger.Error("Failed to scan account balance row", "error", err)
			continue
		}
		balances[currencyCode] = amount
	}

	if err := rows.Err(); err != nil {
		s.logger.Error("Error iterating over balance rows", "error", err)
		return nil, fmt.Errorf("error iterating over balances: %w", err)
	}

	return balances, nil
}

// GetGoalProgress calculates progress percentage for a goal
func (s *Store) GetGoalProgress(id string) (float64, error) {
	goal, err := s.GetGoalByID(id)
	if err != nil {
		return 0, err
	}

	if goal.TargetAmount == 0 {
		return 0, nil
	}

	progress := (goal.CurrentAmount / goal.TargetAmount) * 100
	if progress > 100 {
		progress = 100
	}

	return progress, nil
}

// GetGoalCountByStatus returns count of goals by status
func (s *Store) GetGoalCountByStatus() (map[string]int, error) {
	query := `
		SELECT status, COUNT(*)
		FROM goals
		GROUP BY status
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get goal count by status: %w", err)
	}
	defer rows.Close()

	counts := make(map[string]int)
	for rows.Next() {
		var status string
		var count int
		if err := rows.Scan(&status, &count); err != nil {
			continue
		}
		counts[status] = count
	}

	return counts, nil
}

//=============================================== Goal Contributions =================================================
//====================================================================================================================
//====================================================================================================================

//=============================================== CREATE(C) =================================================

// InsertGoalContribution creates a new goal contribution
func (s *Store) InsertGoalContribution(contribution models.GoalContribution) (string, error) {
	// Generate UUID if not provided
	if contribution.ID == "" {
		contribution.ID = utils.MustGenerateUUID()
	}

	// Set date if not provided
	if contribution.Date.IsZero() {
		contribution.Date = time.Now()
	}

	// Use database transaction for atomicity
	dbTx, err := s.db.Begin()
	if err != nil {
		return "", fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer dbTx.Rollback()

	// Insert contribution
	query := `
		INSERT INTO goal_contributions (id, goal_id, amount, date, note, automatic)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	_, err = dbTx.Exec(query,
		contribution.ID,
		contribution.GoalID,
		contribution.Amount,
		contribution.Date,
		contribution.Note,
		contribution.Automatic,
	)

	if err != nil {
		s.logger.Error("Failed to insert goal contribution", "error", err, "id", contribution.ID)
		return "", fmt.Errorf("failed to insert goal contribution: %w", err)
	}

	// Get goal to check completion
	err = s.AddToGoalAmountTx(dbTx, contribution.GoalID, contribution.Amount)
	if err != nil {
		return "", fmt.Errorf("failed to update goal amount after creating new goal contribution: %w", err)
	}

	// Commit transaction
	if err := dbTx.Commit(); err != nil {
		return "", fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("Goal contribution created", "id", contribution.ID, "goal_id", contribution.GoalID, "amount", contribution.Amount)
	return contribution.ID, nil
}

//=============================================== READ(R) =================================================

// GetContributionsLength returns total number goal contributions
func (s *Store) GetGoalContributionsLength() (int, error) {
	query := `SELECT COUNT(*) FROM goal_contributions`

	var count int
	err := s.db.QueryRow(query).Scan(&count)
	if err != nil {
		s.logger.Error("Failed to return contributions Length", "error", err)
		return 0, fmt.Errorf("failed to get contirbution Length: %w", err)
	}
	return count, nil
}

// GetSpecificGoalContributions returns All contributions for a specific goal
func (s *Store) GetSpecificGoalContributions(goalID string) ([]models.GoalContribution, error) {
	query := `
		SELECT id, goal_id, amount, date, note, automatic
		FROM goal_contributions
		WHERE goal_id = ?
		ORDER BY date DESC
	`

	rows, err := s.db.Query(query, goalID)
	if err != nil {
		return nil, fmt.Errorf("failed to get goal contributions: %w", err)
	}
	defer rows.Close()

	return s.scanGoalContributions(rows)
}

// scanGoalContributions is a helper to scan contribution rows
func (s *Store) scanGoalContributions(rows *sql.Rows) ([]models.GoalContribution, error) {
	var contributions []models.GoalContribution

	for rows.Next() {
		var contribution models.GoalContribution
		err := rows.Scan(
			&contribution.ID,
			&contribution.GoalID,
			&contribution.Amount,
			&contribution.Date,
			&contribution.Note,
			&contribution.Automatic,
		)
		if err != nil {
			s.logger.Error("Failed to scan goal contribution", "error", err)
			continue
		}
		contributions = append(contributions, contribution)
	}

	if err := rows.Err(); err != nil {
		s.logger.Error("Error iterating over balance rows", "error", err)
		return nil, fmt.Errorf("error iterating over balances: %w", err)
	}

	return contributions, nil
}

// GetGoalContributionByID retrieves a single contribution
func (s *Store) GetGoalContributionByID(id string) (*models.GoalContribution, error) {
	query := `
		SELECT id, goal_id, amount, date, note, automatic
		FROM goal_contributions
		WHERE id = ?
	`

	var contribution models.GoalContribution
	err := s.db.QueryRow(query, id).Scan(
		&contribution.ID,
		&contribution.GoalID,
		&contribution.Amount,
		&contribution.Date,
		&contribution.Note,
		&contribution.Automatic,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("goal contribution not found: %s", id)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get goal contribution: %w", err)
	}

	return &contribution, nil
}

// GetRecentContributions retrieves N contributions related to all goals
func (s *Store) GetRecentGoalContributions(limit int) ([]models.GoalContribution, error) {
	query := `
		SELECT id, goal_id, amount, date, note, automatic
		FROM goal_contributions
		ORDER BY date DESC
		LIMIT ?
	`

	rows, err := s.db.Query(query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get goal contributions: %w", err)
	}
	defer rows.Close()

	return s.scanGoalContributions(rows)
}

//=============================================== UPDATE(U) =================================================

// UpdateGoalContribution updates a contribution and adjusts goal amount
func (s *Store) UpdateGoalContribution(newContributionValue models.GoalContribution) error {
	oldGoalContribution, err := s.GetGoalContributionByID(newContributionValue.ID)
	if err != nil {
		return err
	}

	dbTx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("Failed to start transaction (UpdateGoalContribution)")
	}
	defer dbTx.Rollback()

	query := `
		UPDATE goal_contributions
		SET goal_id = ?, amount = ?, date = ?, note = ?, automatic = ?
		WHERE id = ?
	`

	result, err := dbTx.Exec(
		query,
		newContributionValue.GoalID,
		newContributionValue.Amount,
		newContributionValue.Date,
		newContributionValue.Note,
		newContributionValue.Automatic,
		newContributionValue.ID,
	)

	if err != nil {
		s.logger.Error("Faield to update goal contribution", "error", err, "id", newContributionValue.ID)
		return fmt.Errorf("failed to update goal contribution(ID: %s): %w", newContributionValue.ID, err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("goal contribution not found: %s", newContributionValue.ID)
	}

	balanceDifference := newContributionValue.Amount - oldGoalContribution.Amount
	err = s.AddToGoalAmountTx(dbTx, newContributionValue.GoalID, balanceDifference)
	if err != nil {
		return fmt.Errorf("Failed to update goal amount after updating goal contribution: %w", err)
	}

	if err := dbTx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction(UpdateGoalContribution): %w", err)
	}

	s.logger.Info("Goal contribution updated", "id", newContributionValue.ID, "goal_id", newContributionValue.GoalID, "amount", newContributionValue.Amount)
	return nil
}

//=============================================== DELETE(D) =================================================

// DeleteGoalContribution deletes a contribution and adjusts goal amount
func (s *Store) DeleteGoalContribution(id string) error {
	// Get contribution first
	contribution, err := s.GetGoalContributionByID(id)
	if err != nil {
		return err
	}

	// Use transaction for atomicity
	dbTx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer dbTx.Rollback()

	// Delete contribution
	deleteQuery := `DELETE FROM goal_contributions WHERE id = ?`
	_, err = dbTx.Exec(deleteQuery, id)
	if err != nil {
		return fmt.Errorf("failed to delete contribution: %w", err)
	}

	// Reverse the amount (and possibly status change)
	balanceDifference := -contribution.Amount
	err = s.AddToGoalAmountTx(dbTx, contribution.GoalID, balanceDifference)
	if err != nil {
		return fmt.Errorf("Failed to update goal amount after updating goal contribution: %w", err)
	}

	goal, err := s.GetGoalByIDTx(dbTx, contribution.GoalID)
	if err != nil {
		return err
	}

	if goal.Status == string(models.StatusCompleted) && goal.CurrentAmount < goal.TargetAmount {
		err := s.UpdateGoalStatus(goal.ID, string(models.StatusActive))
		if err != nil {
			return err
		}
		// TO DO : we also need to set completed_date = value to nil
	}

	if err := dbTx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("Goal contribution deleted", "id", id)
	return nil
}
