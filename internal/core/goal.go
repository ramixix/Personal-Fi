package core

import (
	"financial_tracker/internal/models"
	"financial_tracker/internal/storage"
	"financial_tracker/internal/utils"
	"time"
)

// AddGoal adds a new goal
func AddGoal(goal models.Goal) error {
	_, err := storage.Store.InsertGoal(goal)
	return err
}

// FindGoal finds a goal by ID
func FindGoal(id string) *models.Goal {
	goal, err := storage.Store.GetGoalByID(id)
	if err != nil {
		return nil
	}
	return goal
}

// DeleteGoal deletes a goal
func DeleteGoal(id string) error {
	return storage.Store.DeleteGoal(id)
}

// GetActiveGoals returns all active goals
func GetActiveGoals() []models.Goal {
	goals, err := storage.Store.GetActiveGoals()
	if err != nil {
		return []models.Goal{}
	}
	return goals
}

// GetCompletedGoals returns all completed goals
func GetCompletedGoals() []models.Goal {
	goals, err := storage.Store.GetCompletedGoals()
	if err != nil {
		return []models.Goal{}
	}
	return goals
}

// AddGoalContribution adds a contribution to a goal
func AddGoalContribution(goalID string, amount float64, note string, automatic bool) bool {
	contribution := models.GoalContribution{
		ID:        utils.MustGenerateUUID(),
		GoalID:    goalID,
		Amount:    amount,
		Date:      time.Now(),
		Note:      note,
		Automatic: automatic,
	}

	_, err := storage.Store.InsertGoalContribution(contribution)
	return err == nil
}

// GetGoalContributions returns contributions for a goal
func GetGoalContributions(goalID string) []models.GoalContribution {
	contributions, err := storage.Store.GetGoalContributions(goalID)
	if err != nil {
		return []models.GoalContribution{}
	}
	return contributions
}

// GetGoalProgress calculates progress percentage
func GetGoalProgress(goal models.Goal) float64 {
	if goal.TargetAmount == 0 {
		return 0
	}
	progressPercentage := (goal.CurrentAmount / goal.TargetAmount) * 100
	if progressPercentage > 100 {
		progressPercentage = 100
	}
	return progressPercentage
}

// GetDaysRemaining calculates days remaining to deadline
func GetRemainingDays(goal models.Goal) int {
	if !goal.HasDeadline {
		return -1
	}
	duration := time.Until(goal.Deadline)
	remainingDays := int(duration.Hours() / 24)
	return remainingDays
}

// GetRequiredMonthlyContribution calculates required monthly contribution to meet target amount till deadline
func GetRequiredMonthlyContribution(goal models.Goal) float64 {
	if !goal.HasDeadline {
		return 0
	}
	remainingAmount := goal.TargetAmount - goal.CurrentAmount
	if remainingAmount <= 0 {
		return 0
	}

	remainingDays := GetRemainingDays(goal)
	if remainingDays <= 0 {
		return 0
	}

	remainingMonths := float64(remainingDays) / 30.0
	if remainingMonths < 1 {
		remainingMonths = 1.0
	}
	return remainingAmount / remainingMonths
}

// GetProjectedCompletionDate calculates projected completion date based on the current contribution rate
func GetProjectedCompletionDate(goal models.Goal) (time.Time, bool) {
	contributions := GetGoalContributions(goal.ID)
	if len(contributions) == 0 {
		// Can't calculate without history
		return time.Time{}, false
	}

	var oldestDate, newestDate time.Time
	totalContributed := 0.0
	for i, contribute := range contributions {
		if i == 0 {
			oldestDate = contribute.Date
			newestDate = contribute.Date
		} else {
			if contribute.Date.Before(oldestDate) {
				oldestDate = contribute.Date
			}
			if contribute.Date.After(newestDate) {
				newestDate = contribute.Date
			}
		}
		totalContributed += contribute.Amount
	}

	durationDays := newestDate.Sub(oldestDate).Hours() / 24
	if durationDays < 1 {
		durationDays = 1
	}

	// Calculate average daily contribution rate
	dailyRate := totalContributed / durationDays
	if dailyRate <= 0 {
		return time.Time{}, false
	}

	// Calculate days needed to reach target
	remainingContributionMoney := goal.TargetAmount - goal.CurrentAmount
	if remainingContributionMoney <= 0 {
		return time.Now(), true // Already completed
	}
	daysNeeded := remainingContributionMoney / dailyRate
	projectedDate := time.Now().Add(time.Duration(daysNeeded*24) * time.Hour)
	return projectedDate, true
}

// GetTotalGoalsSaved returns total saved across all goals
func GetTotalGoalsSaved() float64 {
	total, err := storage.Store.GetTotalGoalsSaved()
	if err != nil {
		return 0
	}
	return total
}

// UpdateGoalStatus updates goal status
func UpdateGoalStatus(goalID, newStatus string) bool {
	err := storage.Store.UpdateGoalStatus(goalID, newStatus)
	return err == nil
}
