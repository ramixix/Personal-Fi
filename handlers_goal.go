package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func handleGoals() {
	if len(os.Args) < 3 {
		fmt.Println("Goal Management")
		fmt.Println("===============")
		fmt.Println("Usage: financial-tracker goals [command]")
		fmt.Println("Commands:")
		fmt.Println("  create      Create a new goal")
		fmt.Println("  list        List all goals")
		fmt.Println("  active      List active goals")
		fmt.Println("  completed   List completed goals")
		fmt.Println("  view        View detailed goal information")
		fmt.Println("  contribute  Add contribution to a goal")
		fmt.Println("  update      Update goal status")
		fmt.Println("  delete      Delete a goal")
		return
	}

	subCommand := strings.ToLower(os.Args[2])
	switch subCommand {
	case "create":
		handleCreateGoal()
	case "list":
		handleListGoals()
	case "active":
		handleListActiveGoals()
	case "completed":
		handleListCompletedGoals()
	case "view":
		handleViewGoal()
	case "contribute":
		handleContributeToGoal()
	case "update":
		handleUpdateGoalStatus()
	case "delete":
		handleDeleteGoal()
	default:
		fmt.Printf("Unknown goals command: %s\n", subCommand)
	}
}

func handleCreateGoal() {
	fmt.Println("Create New Goal")
	fmt.Println("===============")

	reader := bufio.NewReader(os.Stdin)

	// Get goal name
	name := getNonEmptyString(reader, "Goal name: ")

	// Get description
	fmt.Print("Description (optional): ")
	descInput, _ := reader.ReadString('\n')
	description := strings.TrimSpace(descInput)

	// Get target amount
	targetAmount := getValidAmount(reader)

	// Get category
	fmt.Println("\nCategory options: savings, debt, investment, purchase, other ...(if empty set to saving by default)")
	fmt.Print("Category: ")
	categoryInput, _ := reader.ReadString('\n')
	category := strings.TrimSpace(categoryInput)
	if category == "" {
		category = "savings"
	}

	// Get priority
	fmt.Println("\nPriority options: high, medium, low (anything else would cause the priority to get set as medium by default)")
	fmt.Print("Priority (default: medium): ")
	priorityInput, _ := reader.ReadString('\n')
	priority := strings.ToLower(strings.TrimSpace(priorityInput))
	if priority != "high" && priority != "medium" && priority != "low" {
		priority = "medium"
	}

	// Ask about deadline
	fmt.Print("\nDo you want to set a deadline? (yes/no): ")
	deadlineChoice, _ := reader.ReadString('\n')
	hasDeadline := strings.ToLower(strings.TrimSpace(deadlineChoice)) == "yes"

	var deadline time.Time
	if hasDeadline {
		fmt.Print("Deadline (YYYY-MM-DD): ")
		deadlineInput, _ := reader.ReadString('\n')
		parsedDeadline, err := parseDate(strings.TrimSpace(deadlineInput))
		if err != nil {
			fmt.Println("Invalid date format, proceeding without deadline")
			hasDeadline = false
		} else {
			deadline = parsedDeadline
		}
	}

	// Ask about linking to account
	var linkedAccountID int
	if len(accounts) > 0 {
		fmt.Print("\nDo you want to link this goal to an account? (yes/no): ")
		linkChoice, _ := reader.ReadString('\n')
		if strings.ToLower(strings.TrimSpace(linkChoice)) == "yes" {
			fmt.Println("\nAvailable accounts:")
			for _, account := range accounts {
				fmt.Printf("  ID: %d - %s\n", account.ID, account.Name)
			}
			accountID, err := getIntInput(reader, "Account ID (or 0 for none): ")
			if err == nil && findAccountIndex(accountID) != -1 {
				linkedAccountID = accountID
			}
		}
	}

	// Create goal
	newGoal := Goal{
		ID:              nextGoalID,
		Name:            name,
		Description:     description,
		TargetAmount:    targetAmount,
		CurrentAmount:   0.0,
		Deadline:        deadline,
		HasDeadline:     hasDeadline,
		Category:        category,
		Priority:        priority,
		Status:          "active",
		LinkedAccountID: linkedAccountID,
		Created:         time.Now(),
	}

	addGoal(newGoal)
	fmt.Printf("\n✓ Goal '%s' created successfully! ID: %d\n", name, newGoal.ID)

	// Show required monthly contribution if deadline is set
	if hasDeadline {
		required := getRequiredMonthlyContribution(newGoal)
		fmt.Printf("💡 To reach your goal by the deadline, you need to save $%.2f per month\n", required)
	}
}

// List all goals
func handleListGoals() {
	fmt.Println("All Goals")
	fmt.Println("=========")

	if len(goals) == 0 {
		fmt.Println("No goals found. Create one with: goals create")
		return
	}

	for _, goal := range goals {
		displayGoalSummary(goal)
	}

	// Show totals
	fmt.Printf("\n--- Summary ---\n")
	fmt.Printf("Total Goals: %d\n", len(goals))
	fmt.Printf("Total Saved: $%.2f\n", getTotalGoalsSaved())
}

// List active goals
func handleListActiveGoals() {
	fmt.Println("Active Goals")
	fmt.Println("============")

	activeGoals := getActiveGoals()

	if len(activeGoals) == 0 {
		fmt.Println("No active goals.")
		return
	}

	for _, goal := range activeGoals {
		displayGoalSummary(goal)
	}
	fmt.Printf("\nTotal Active Goals: %d\n", len(activeGoals))
}

// List completed goals
func handleListCompletedGoals() {
	fmt.Println("Completed Goals")
	fmt.Println("===============")

	completedGoals := getCompletedGoals()

	if len(completedGoals) == 0 {
		fmt.Println("No completed goals yet. Keep going!")
		return
	}

	for _, goal := range completedGoals {
		displayGoalSummary(goal)
	}

	fmt.Printf("\nTotal Completed Goals: %d\n", len(completedGoals))
}

// Display goal summary (one line)
func displayGoalSummary(goal Goal) {
	progress := getGoalProgress(goal)
	statusSymbol := "🎯"
	if goal.Status == "completed" {
		statusSymbol = "✅"
	} else if goal.Status == "paused" {
		statusSymbol = "⏸️"
	}

	fmt.Printf("%s ID: %d | %s | $%.2f / $%.2f (%.1f%%) | %s | Priority: %s\n",
		statusSymbol,
		goal.ID,
		goal.Name,
		goal.CurrentAmount,
		goal.TargetAmount,
		progress,
		goal.Status,
		goal.Priority)
}

// View detailed goal information
func handleViewGoal() {
	if len(goals) == 0 {
		fmt.Println("No goals available.")
		return
	}

	fmt.Println("View Specific Goal With Full Details")
	fmt.Println("====================================")

	fmt.Println("Available goals:")
	for _, goal := range goals {
		fmt.Printf("  ID: %d - %s\n", goal.ID, goal.Name)
	}

	reader := bufio.NewReader(os.Stdin)
	goalID, err := getIntInput(reader, "\nEnter ID of goal you want to view in details: ")
	if err != nil {
		fmt.Println("Invalid goal ID! Select What is presented in given list.")
		return
	}

	goal := findGoal(goalID)
	if goal == nil {
		fmt.Println("Goal not found!")
		return
	}

	// Display detailed information
	fmt.Printf("\n=== %s ===\n", goal.Name)
	fmt.Printf("ID: %d\n", goal.ID)
	fmt.Printf("Status: %s\n", goal.Status)
	fmt.Printf("Category: %s\n", goal.Category)
	fmt.Printf("Priority: %s\n", goal.Priority)

	if goal.Description != "" {
		fmt.Printf("Description: %s\n", goal.Description)
	}

	fmt.Printf("\n--- Progress ---\n")
	fmt.Printf("Current Amount: $%.2f\n", goal.CurrentAmount)
	fmt.Printf("Target Amount:  $%.2f\n", goal.TargetAmount)
	fmt.Printf("Remaining Amount:      $%.2f\n", goal.TargetAmount-goal.CurrentAmount)
	fmt.Printf("Progress:       %.1f%%\n", getGoalProgress(*goal))
	// Show progress bar
	progressBar := generateProgressBar(getGoalProgress(*goal), 30)
	fmt.Printf("%s\n", progressBar)

	// Deadline information
	if goal.HasDeadline {
		fmt.Printf("\n--- Deadline ---\n")
		fmt.Printf("Target Date: %s\n", goal.Deadline.Format("2006-01-02"))
		daysRemaining := getRemainingDays(*goal)
		if daysRemaining >= 0 {
			fmt.Printf("Remaining Days: %d\n", daysRemaining)
			required := getRequiredMonthlyContribution(*goal)
			fmt.Printf("Required Monthly: $%.2f\n", required)
		} else {
			fmt.Println("⚠️  Deadline has passed!")
		}
	}

	// Projection
	if projectedDate, canProject := getProjectedCompletionDate(*goal); canProject {
		fmt.Printf("\n--- Projection ---\n")
		fmt.Printf("Projected Completion: %s\n", projectedDate.Format("2006-01-02"))

		if goal.HasDeadline {
			if projectedDate.Before(goal.Deadline) {
				fmt.Println("✅ On track to meet deadline!")
			} else {
				daysLate := int(projectedDate.Sub(goal.Deadline).Hours() / 24)
				fmt.Printf("⚠️  Currently %d days behind schedule\n", daysLate)
			}
		}
	}

	// Linked account
	if goal.LinkedAccountID != 0 {
		account := findAccount(goal.LinkedAccountID)
		if account != nil {
			fmt.Printf("\n--- Linked Account ---\n")
			fmt.Printf("%s (Balance: $%.2f)\n", account.Name, account.Balance)
		}
	}

	// Contribution history
	contributions := getGoalContributions(goalID)
	if len(contributions) > 0 {
		fmt.Printf("\n--- Recent 10 Goal Contributions ---\n")
		count := 10
		if len(contributions) < 10 {
			count = len(contributions)
		}
		startIndex := len(contributions) - count
		for i := startIndex; i < len(contributions); i++ {
			c := contributions[i]
			fmt.Printf("%s | +$%.2f | %s\n", c.Date.Format("2006-01-02"), c.Amount, c.Note)
		}
		if len(contributions) > 10 {
			fmt.Printf("(... and %d more)\n", len(contributions)-10)
		}
	}

	// Created date + if finished then the finished date
	fmt.Printf("\nCreated: %s\n", goal.Created.Format("2006-01-02"))
	if goal.Status == "completed" {
		fmt.Printf("Completed: %s\n", goal.CompletedDate.Format("2006-01-02"))
	}
}

func generateProgressBar(percentage float64, width int) string {
	filled := int((percentage / 100) * float64(width))
	if filled > width {
		filled = width
	}

	bar := "["
	for i := 0; i < width; i++ {
		if i < filled {
			bar += "="
		} else {
			bar += " "
		}
	}
	bar += fmt.Sprintf("] %.1f Percentage", percentage)
	return bar
}

// Contribute to a goal
func handleContributeToGoal() {
	if len(goals) == 0 {
		fmt.Println("No goals available. Create one first.")
		return
	}

	fmt.Println("Contribute to Goal")
	fmt.Println("==================")

	// Show active goals
	activeGoals := getActiveGoals()
	if len(activeGoals) == 0 {
		fmt.Println("No active goals available.")
		return
	}

	fmt.Println("Active goals:")
	for _, goal := range activeGoals {
		fmt.Printf("  ID: %d - %s ($%.2f / $%.2f)\n",
			goal.ID, goal.Name, goal.CurrentAmount, goal.TargetAmount)
	}

	reader := bufio.NewReader(os.Stdin)
	goalID, err := getIntInput(reader, "Enter Goal ID: ")
	if err != nil {
		fmt.Println("Invalid goal ID!")
		return
	}

	goal := findGoal(goalID)
	if goal == nil {
		fmt.Println("Goal not found!")
		return
	}

	if goal.Status != "active" {
		fmt.Printf("Goal is %s and cannot accept contributions.\n", goal.Status)
		return
	}

	// Get amount
	amount := getValidAmount(reader)

	// Get note
	fmt.Print("Note (optional): ")
	noteInput, _ := reader.ReadString('\n')
	note := strings.TrimSpace(noteInput)
	automatic := false

	if addGoalContribution(goalID, amount, note, automatic) {
		fmt.Printf("\n✓ Added $%.2f to '%s'\n", amount, goal.Name)
		fmt.Printf("New balance: $%.2f / $%.2f (%.1f%%)\n", goal.CurrentAmount, goal.TargetAmount, getGoalProgress(*goal))

		// Check if goal is completed
		if goal.Status == "completed" {
			fmt.Println("\n🎉 Congratulations! Goal completed! 🎉")
		}
	} else {
		fmt.Println("Failed to add contribution.")
	}
}

func handleUpdateGoalStatus() {
	if len(goals) == 0 {
		fmt.Println("No goals available.")
		return
	}

	fmt.Println("Update Goal Status")
	fmt.Println("==================")

	// Show all goals
	fmt.Println("All goals:")
	for _, goal := range goals {
		fmt.Printf("  ID: %d - %s (Current status: %s)\n", goal.ID, goal.Name, goal.Status)
	}

	reader := bufio.NewReader(os.Stdin)
	goalID, err := getIntInput(reader, "\nEnter goal ID: ")
	if err != nil {
		fmt.Println("Invalid goal ID!")
		return
	}

	goal := findGoal(goalID)
	if goal == nil {
		fmt.Println("Goal not found!")
		return
	}

	validStatuses := []string{"active", "paused", "cancelled", "completed"}
	var newStatus string
	fmt.Printf("\nStatus options: %s\n", strings.Join(validStatuses, ", "))
	fmt.Printf("Current status: %s\n", goal.Status)
	for {
		fmt.Print("New Status: ")
		newStatus, _ = reader.ReadString('\n')
		newStatus = strings.ToLower(strings.TrimSpace(newStatus))

		valid := false
		for _, s := range validStatuses {
			if newStatus == s {
				valid = true
				break
			}
		}
		if valid {
			break
		}
		fmt.Println("Invalid status! Please select from the list provide to you.")
	}

	if updateGoalStatus(goalID, newStatus) {
		fmt.Printf("✓ Goal status updated to: %s\n", newStatus)
		if newStatus == "completed" {
			fmt.Println("🎉 Congratulations on completing your goal! 🎉")
		}
	} else {
		fmt.Println("Failed to update goal status.")
	}
}

// Delete a goal
func handleDeleteGoal() {
	if len(goals) == 0 {
		fmt.Println("No goals to delete.")
		return
	}

	fmt.Println("Delete Goal")
	fmt.Println("===========")

	// Show all goals
	handleListGoals()

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("\nEnter goal ID to delete (or 'cancel'): ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "cancel" {
		fmt.Println("Deletion cancelled.")
		return
	}

	id, err := strconv.Atoi(input)
	if err != nil {
		fmt.Println("Invalid ID!")
		return
	}

	if !getConfirmation(reader, "Are you sure? This will delete the goal and all contribution history. (yes/no): ") {
		fmt.Println("Deletion cancelled.")
		return
	}

	if deleteGoal(id) {
		fmt.Printf("✓ Goal ID %d deleted successfully!\n", id)
	} else {
		fmt.Printf("Goal ID %d not found!\n", id)
	}
}
