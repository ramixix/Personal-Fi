package main

import "time"

// Add a new goal
func addGoal(goal Goal) {
	goals = append(goals, goal)
	nextGoalID++
}

// Find goal by ID
func findGoal(id int) *Goal {
	for i := range goals {
		if goals[i].ID == id {
			return &goals[i]
		}
	}
	return nil
}

// Find goal index by ID
func findGoalIndex(id int) int {
	for i, goal := range goals {
		if goal.ID == id {
			return i
		}
	}
	return -1
}

// Delete goal by ID, Also delete all goal contributions for the specified goal
func deleteGoal(id int) bool {
	for i := range goals {
		if goals[i].ID == id {
			var filteredGoalContributions []GoalContribution
			for _, contribution := range goalContributions {
				if contribution.GoalID != id {
					filteredGoalContributions = append(filteredGoalContributions, contribution)
				}
			}
			goalContributions = filteredGoalContributions
			goals = append(goals[:i], goals[i+1:]...)
			return true
		}
	}
	return false
}

// Add contribution to a goal
func addGoalContribution(goalID int, amount float64, note string, automatic bool) bool {
	goal := findGoal(goalID)
	if goal == nil {
		return false
	}

	goal.CurrentAmount += amount
	if goal.CurrentAmount >= goal.TargetAmount && goal.Status == "active" {
		goal.CompletedDate = time.Now()
		goal.Status = "complete"
	}

	contribution := GoalContribution{ID: nextGoalContributionID, GoalID: goalID, Amount: amount, Date: time.Now(), Note: note, Automatic: automatic}
	goalContributions = append(goalContributions, contribution)
	nextGoalContributionID++
	return true
}

// Get contributions for a specific goal
func getGoalContributions(goalID int) []GoalContribution {
	var contributions []GoalContribution
	for _, contribute := range goalContributions {
		if contribute.GoalID == goalID {
			contributions = append(contributions, contribute)
		}
	}
	return contributions
}

// Get active goals
func getActiveGoals() []Goal {
	var activeGoals []Goal
	for _, goal := range goals {
		if goal.Status == "active" {
			activeGoals = append(activeGoals, goal)
		}
	}
	return activeGoals
}

// Get completed goals
func getCompletedGoals() []Goal {
	var completed []Goal
	for _, goal := range goals {
		if goal.Status == "completed" {
			completed = append(completed, goal)
		}
	}
	return completed
}

// Calculate goal progress percentage
func getGoalProgress(goal Goal) float64 {
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
func getRemainingDays(goal Goal) int {
	if !goal.HasDeadline {
		return -1
	}
	duration := time.Until(goal.Deadline)
	remainingDays := int(duration.Hours() / 24)
	return remainingDays
}

// Calculate required monthly contribution to meet deadline
func getRequiredMonthlyContribution(goal Goal) float64 {
	if !goal.HasDeadline {
		return 0
	}
	remainingAmount := goal.TargetAmount - goal.CurrentAmount
	if remainingAmount <= 0 {
		return 0
	}

	remainingDays := getRemainingDays(goal)
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
func getProjectedCompletionDate(goal Goal) (time.Time, bool) {
	contributions := getGoalContributions(goal.ID)
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
func getTotalGoalsSaved() float64 {
	var total float64
	for _, goal := range goals {
		total += goal.CurrentAmount
	}
	return total
}

// Update goal status
func updateGoalStatus(goalID int, newStatus string) bool {
	goalIndex := findGoalIndex(goalID)
	if goalIndex == -1 {
		return false
	}

	goals[goalIndex].Status = newStatus

	if newStatus == "completed" && goals[goalIndex].CompletedDate.IsZero() {
		goals[goalIndex].CompletedDate = time.Now()
	}

	return true
}
