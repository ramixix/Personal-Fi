package core

import (
	"financial_tracker/internal/models"
	"financial_tracker/internal/storage"
	"time"
)

// Add a new goal
func AddGoal(goal models.Goal) {
	storage.Goals = append(storage.Goals, goal)
	storage.NextGoalID++
}

// Find goal by ID
func FindGoal(id int) *models.Goal {
	for i := range storage.Goals {
		if storage.Goals[i].ID == id {
			return &storage.Goals[i]
		}
	}
	return nil
}

// Find goal index by ID
func FindGoalIndex(id int) int {
	for i, goal := range storage.Goals {
		if goal.ID == id {
			return i
		}
	}
	return -1
}

// Delete goal by ID, Also delete all goal contributions for the specified goal
func DeleteGoal(id int) bool {
	for i := range storage.Goals {
		if storage.Goals[i].ID == id {
			var filteredGoalContributions []models.GoalContribution
			for _, contribution := range storage.GoalContributions {
				if contribution.GoalID != id {
					filteredGoalContributions = append(filteredGoalContributions, contribution)
				}
			}
			storage.GoalContributions = filteredGoalContributions
			storage.Goals = append(storage.Goals[:i], storage.Goals[i+1:]...)
			return true
		}
	}
	return false
}

// Add contribution to a goal
func AddGoalContribution(goalID int, amount float64, note string, automatic bool) bool {
	goal := FindGoal(goalID)
	if goal == nil {
		return false
	}

	goal.CurrentAmount += amount
	if goal.CurrentAmount >= goal.TargetAmount && goal.Status == "active" {
		goal.CompletedDate = time.Now()
		goal.Status = "complete"
	}

	contribution := models.GoalContribution{ID: storage.NextGoalContributionID, GoalID: goalID, Amount: amount, Date: time.Now(), Note: note, Automatic: automatic}
	storage.GoalContributions = append(storage.GoalContributions, contribution)
	storage.NextGoalContributionID++
	return true
}

// Get contributions for a specific goal
func GetGoalContributions(goalID int) []models.GoalContribution {
	var contributions []models.GoalContribution
	for _, contribute := range storage.GoalContributions {
		if contribute.GoalID == goalID {
			contributions = append(contributions, contribute)
		}
	}
	return contributions
}

// Get active goals
func GetActiveGoals() []models.Goal {
	var activeGoals []models.Goal
	for _, goal := range storage.Goals {
		if goal.Status == "active" {
			activeGoals = append(activeGoals, goal)
		}
	}
	return activeGoals
}

// Get completed goals
func GetCompletedGoals() []models.Goal {
	var completed []models.Goal
	for _, goal := range storage.Goals {
		if goal.Status == "completed" {
			completed = append(completed, goal)
		}
	}
	return completed
}

// Calculate goal progress percentage
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

// Calculate days remaining to deadline
func GetRemainingDays(goal models.Goal) int {
	if !goal.HasDeadline {
		return -1
	}
	duration := time.Until(goal.Deadline)
	remainingDays := int(duration.Hours() / 24)
	return remainingDays
}

// Calculate required monthly contribution to meet deadline
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

// Calculate projected completion date based on current contribution rate
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

// Get total amount saved across all goals
func GetTotalGoalsSaved() float64 {
	var total float64
	for _, goal := range storage.Goals {
		total += goal.CurrentAmount
	}
	return total
}

// Update goal status
func UpdateGoalStatus(goalID int, newStatus string) bool {
	goalIndex := FindGoalIndex(goalID)
	if goalIndex == -1 {
		return false
	}

	storage.Goals[goalIndex].Status = newStatus

	if newStatus == "completed" && storage.Goals[goalIndex].CompletedDate.IsZero() {
		storage.Goals[goalIndex].CompletedDate = time.Now()
	}

	return true
}
