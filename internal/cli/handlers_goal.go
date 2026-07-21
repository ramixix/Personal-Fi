package cli

import (
	"bufio"
	"financial_tracker/internal/core"
	"financial_tracker/internal/models"
	"financial_tracker/internal/utils"
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
		fmt.Println("  list        List last 100 goals")
		fmt.Println("  active      List active goals")
		fmt.Println("  completed   List completed goals")
		fmt.Println("  fullDetail  View detailed goal information")
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
		handleListGoals(recent100)
	case "active":
		handleListActiveGoals()
	case "completed":
		handleListCompletedGoals()
	case "fullDetail":
		handleViewSpecificGoalDetails()
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

// handleCreateGoal handles creating goals through CLI
func handleCreateGoal() {
	fmt.Println("Create New Goal")
	fmt.Println("===============")

	reader := bufio.NewReader(os.Stdin)

	// Get goal name
	name := utils.GetNonEmptyString(reader, "Goal name: ")

	// Get description
	fmt.Print("Description (optional): ")
	descInput, _ := reader.ReadString('\n')
	description := strings.TrimSpace(descInput)

	// Get target amount
	fmt.Print("Target ")
	targetAmount := utils.GetValidAmount(reader)

	currencyCode := utils.GetValidCurrency(reader, "USD")

	// Get category
	fmt.Println("\nCategory options: savings, debt, investment, purchase, other ...(if empty  by default will set to 'saving')")
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
		parsedDeadline, err := utils.ParseDate(strings.TrimSpace(deadlineInput))
		if err != nil {
			fmt.Println("Invalid date format, proceeding without deadline")
			hasDeadline = false
		} else {
			deadline = parsedDeadline
		}
	}

	// Ask about linking to account
	var linkedAccountID string
	accounts := core.GetRecentAccounts(recent100)
	if len(accounts) > 0 {
		fmt.Print("\nDo you want to link this goal to an account? (yes/no): ")
		linkChoice, _ := reader.ReadString('\n')
		if strings.ToLower(strings.TrimSpace(linkChoice)) == "yes" {
			fmt.Printf("\n[Info] Total Account Number: %d.\n", core.GetAccountsLength())
			fmt.Println("\nPress Enter to show the latest %d accounts (if possible by default).\nTo see a different number, enter it.\nType 'all' to display all accounts.", recent100)
			accountsToShow := GetAccountsNumberToShow(reader, recent100)
			handleListAccounts(accountsToShow)

			accountID := utils.GetNonEmptyString(reader, "Account ID: ")
			if core.FindAccount(accountID) != nil {
				linkedAccountID = accountID
			}
		}
	}

	// Create goal
	newGoal := models.Goal{
		ID:              utils.MustGenerateUUID(),
		Name:            name,
		Description:     description,
		TargetAmount:    targetAmount,
		CurrentAmount:   0.0,
		CurrencyCode:    currencyCode,
		Deadline:        deadline,
		HasDeadline:     hasDeadline,
		Category:        category,
		Priority:        priority,
		Status:          "active",
		LinkedAccountID: linkedAccountID,
		Created:         time.Now(),
	}

	core.AddGoal(newGoal)
	fmt.Printf("\n✓ Goal '%s' created successfully! ID: %s\n", name, newGoal.ID)

	// Show required monthly contribution if deadline is set
	if hasDeadline {
		required := core.GetRequiredMonthlyContribution(newGoal)
		fmt.Printf("💡 To reach your goal by the deadline, you need to save $%.2f per month\n", required)
	}
}

// handleListGoals(N) lists recent goals
func handleListGoals(numberofGoalsToShow int) {
	fmt.Println("List Goals")
	fmt.Println("==========")
	goalsCount := core.GetGoalsLength("")

	if goalsCount <= 0 {
		fmt.Println("No goals found. Create one with: goals create")
		return
	}

	goalsToShow := core.GetRecentGoals(numberofGoalsToShow)
	for _, goal := range goalsToShow {
		displayGoalSummary(goal)
	}

	// Show totals
	fmt.Printf("\n--- Summary ---\n")
	fmt.Printf("Total Goals: %d\n", goalsCount)
	fmt.Println("\nTotal Across All Goals:")
	fmt.Println("--------------------------")
	totals := core.GetTotalAccountsBalanceByCurrency()
	for currency, total := range totals {
		fmt.Printf("%s: %.2f\n", currency, total)
	}
}

// handleListActiveGoals lists all active goals
func handleListActiveGoals() {
	fmt.Println("Active Goals")
	fmt.Println("============")

	activeGoals := core.GetActiveGoals()

	if len(activeGoals) == 0 {
		fmt.Println("No active goals.")
		return
	}

	for _, goal := range activeGoals {
		displayGoalSummary(goal)
	}
	fmt.Printf("\nTotal Active Goals: %d\n", len(activeGoals))
}

// handleListCompletedGoals lists all completed goals
func handleListCompletedGoals() {
	fmt.Println("Completed Goals")
	fmt.Println("===============")

	completedGoals := core.GetCompletedGoals()

	if len(completedGoals) == 0 {
		fmt.Println("No completed goals yet. Keep going!")
		return
	}

	for _, goal := range completedGoals {
		displayGoalSummary(goal)
	}

	fmt.Printf("\nTotal Completed Goals: %d\n", len(completedGoals))
}

// displayGoalSummary(goal) displays short summary of given gool(one line)
func displayGoalSummary(goal models.Goal) {
	progress := core.GetGoalProgress(goal)
	statusSymbol := "🎯"
	if goal.Status == "completed" {
		statusSymbol = "✅"
	} else if goal.Status == "paused" {
		statusSymbol = "⏸️"
	}

	fmt.Printf("%s ID: %s | %s | $%.2f / $%.2f (%.1f%%) | %s | %s | Priority: %s\n",
		statusSymbol,
		goal.ID,
		goal.Name,
		goal.CurrentAmount,
		goal.TargetAmount,
		progress,
		goal.CurrencyCode,
		goal.Status,
		goal.Priority)
}

// handleViewSpecificGoalDetails views detailed goal information
func handleViewSpecificGoalDetails() {
	goalsCount := core.GetGoalsLength("")
	if goalsCount <= 0 {
		fmt.Println("No goals available.")
		return
	}

	fmt.Println("View Specific Goal With Full Details")
	fmt.Println("====================================")
	fmt.Printf("\n[Info] Total Goals Number: %d.\n", goalsCount)
	fmt.Println("\nPress Enter to show the latest %d goals (if possible by default).\nTo see a different number, enter it.\nType 'all' to display all goals.", recent100)
	reader := bufio.NewReader(os.Stdin)

	goalsToShow := GetGoalsNumberToShow(reader, recent100)
	handleListGoals(goalsToShow)

	goalID := utils.GetNonEmptyString(reader, "\nEnter ID of the goal you want to view in details: ")

	goal := core.FindGoal(goalID)
	if goal == nil {
		fmt.Println("Goal not found. Please verify the ID you entered.")
		return
	}

	// Display detailed information
	fmt.Printf("\n=== %s ===\n", goal.Name)
	fmt.Printf("ID: %s\n", goal.ID)
	fmt.Printf("Status: %s\n", goal.Status)
	fmt.Printf("Category: %s\n", goal.Category)
	fmt.Printf("Priority: %s\n", goal.Priority)
	fmt.Printf("Currency: %s\n", goal.CurrencyCode)

	if goal.Description != "" {
		fmt.Printf("Description: %s\n", goal.Description)
	}

	fmt.Printf("\n--- Progress ---\n")
	fmt.Printf("Current Amount: %s\n", utils.FormatCurrency(goal.CurrentAmount, goal.CurrencyCode))
	fmt.Printf("Target Amount:  $%.2f\n", goal.TargetAmount)
	fmt.Printf("Remaining Amount:      $%.2f\n", goal.TargetAmount-goal.CurrentAmount)

	var goalProgress float64 = core.GetGoalProgress(*goal)
	fmt.Printf("Progress:       %.1f%%\n", goalProgress)
	// Show progress bar
	progressBar := generateProgressBar(goalProgress, 30)
	fmt.Printf("%s\n", progressBar)

	// Deadline information
	if goal.HasDeadline {
		fmt.Printf("\n--- Deadline ---\n")
		fmt.Printf("Target Date: %s\n", goal.Deadline.Format("2006-01-02"))
		daysRemaining := core.GetRemainingDays(*goal)
		if daysRemaining >= 0 {
			fmt.Printf("Remaining Days: %d\n", daysRemaining)
			required := core.GetRequiredMonthlyContribution(*goal)
			fmt.Printf("Required Monthly: $%.2f\n", required)
		} else {
			fmt.Println("⚠️  Deadline has passed!")
		}
	}

	// Projection
	if projectedDate, canProject := core.GetProjectedCompletionDate(*goal); canProject {
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
	if goal.LinkedAccountID != "" {
		account := core.FindAccount(goal.LinkedAccountID)
		if account != nil {
			fmt.Printf("\n--- Linked Account ---\n")
			fmt.Printf("%s (Balance: $%.2f)\n", account.Name, account.Balance)
		}
	}

	// Contribution history
	contributions := core.GetSpecificGoalContributions(goalID)
	if len(contributions) > 0 {
		contributionsCount := len(contributions)
		fmt.Printf("\n--- Total Contributions For This Goal %d ", contributionsCount)
		count := recent100
		if contributionsCount < recent100 {
			count = contributionsCount
		}
		fmt.Printf("\n--- Showing Recent %d Goal Contributions ---\n", count)

		startIndex := contributionsCount - count
		for i := startIndex; i < contributionsCount; i++ {
			c := contributions[i]
			fmt.Printf("%s | +$%.2f | %s\n", c.Date.Format("2006-01-02"), c.Amount, c.Note)
		}
	}

	// Created date + if finished then the finished date
	fmt.Printf("\nCreated: %s\n", goal.Created.Format("2006-01-02"))
	if goal.Status == "completed" {
		fmt.Printf("Completed: %s\n", goal.CompletedDate.Format("2006-01-02"))
	}
}

// generateProgressBar creates a bar based percentage of goals current amount and target amount
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

// handleContributeToGoal Contribute to a goal
func handleContributeToGoal() {
	goalsCount := core.GetGoalsLength("")
	if goalsCount <= 0 {
		fmt.Println("No goals available. Create one first.")
		return
	}

	fmt.Println("Contribute to Goal")
	fmt.Println("==================")

	// Show active goals
	activeGoals := core.GetActiveGoals()
	if len(activeGoals) == 0 {
		fmt.Println("No active goals available.")
		return
	}

	fmt.Println("Active goals:")
	for _, goal := range activeGoals {
		fmt.Printf("  ID: %s - %s ($%.2f / $%.2f %s)\n",
			goal.ID, goal.Name, goal.CurrentAmount, goal.TargetAmount, goal.CurrencyCode)
	}

	reader := bufio.NewReader(os.Stdin)
	goalID := utils.GetNonEmptyString(reader, "Enter Goal ID: ")

	goal := core.FindGoal(goalID)
	if goal == nil {
		fmt.Println("Goal not found. Please verify the ID you entered.")
		return
	}

	if goal.Status != "active" {
		fmt.Printf("Goal is %s and cannot accept contributions.\n", goal.Status)
		return
	}

	// Get amount
	amount := utils.GetValidAmount(reader)

	// Get note
	fmt.Print("Note (optional): ")
	noteInput, _ := reader.ReadString('\n')
	note := strings.TrimSpace(noteInput)
	automatic := false

	if core.AddGoalContribution(goalID, amount, note, automatic) {
		fmt.Printf("\n✓ Added $%.2f to '%s'\n", amount, goal.Name)
		fmt.Printf("New balance: $%.2f / $%.2f (%.1f%%)\n", goal.CurrentAmount, goal.TargetAmount, core.GetGoalProgress(*goal))

		// Check if goal is completed
		if goal.Status == "completed" {
			fmt.Println("\n🎉 Congratulations! Goal completed! 🎉")
		}
	} else {
		fmt.Println("Failed to add contribution.")
	}
}

// handleUpdateGoalStatus updates goal status
func handleUpdateGoalStatus() {
	goalsCount := core.GetGoalsLength("")
	if goalsCount == 0 {
		fmt.Println("No goals available.")
		return
	}

	fmt.Println("Update Goal Status")
	fmt.Println("==================")
	fmt.Printf("\n[Info] Total Goals Number: %d.\n", goalsCount)
	fmt.Println("\nPress Enter to show the latest %d goals (if possible by default).\nTo see a different number, enter it.\nType 'all' to display all goals.", recent100)
	reader := bufio.NewReader(os.Stdin)

	goalsToShow := GetGoalsNumberToShow(reader, recent100)
	handleListGoals(goalsToShow)

	goalID := utils.GetNonEmptyString(reader, "\nEnter goal ID: ")

	goal := core.FindGoal(goalID)
	if goal == nil {
		fmt.Println("Goal not found. Please verify the ID you entered.")
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

	if core.UpdateGoalStatus(goalID, newStatus) {
		fmt.Printf("✓ Goal status updated to: %s\n", newStatus)
		if newStatus == "completed" {
			fmt.Println("🎉 Congratulations on completing your goal! 🎉")
		}
	} else {
		fmt.Println("Failed to update goal status.")
	}
}

// handleDeleteGoal deletes a goal
func handleDeleteGoal() {
	goalsCount := core.GetGoalsLength("")
	if goalsCount <= 0 {
		fmt.Println("No goals to delete.")
		return
	}

	fmt.Println("Delete Goal")
	fmt.Println("===========")
	fmt.Printf("\n[Info] Total Goals Number: %d.\n", goalsCount)
	fmt.Println("\nPress Enter to show the latest %d goals (if possible by default).\nTo see a different number, enter it.\nType 'all' to display all goals.", recent100)
	reader := bufio.NewReader(os.Stdin)

	goalsToShow := GetGoalsNumberToShow(reader, recent100)
	handleListGoals(goalsToShow)

	fmt.Print("\nEnter goal ID to delete (or 'cancel'): ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "cancel" {
		fmt.Println("Deletion cancelled.")
		return
	}

	if !utils.GetConfirmation(reader, "Are you sure? This will delete the goal and all contribution history. (yes/no): ") {
		fmt.Println("Deletion cancelled.")
		return
	}

	if core.DeleteGoal(input) == nil {
		fmt.Printf("✓ Goal ID %s deleted successfully!\n", input)
	} else {
		fmt.Printf("Goal ID %d not found!\n", input)
	}
}

// GetAccountsNumberToShow asks users for a specific number N to list last N Accounts.
func GetGoalsNumberToShow(reader *bufio.Reader, defaultValue int) int {
	goalsToShow := defaultValue
	totalGoals := core.GetGoalsLength("")
	if defaultValue > totalGoals {
		goalsToShow = totalGoals
	}
InputLoop:
	for {
		input, _ := reader.ReadString('\n')
		input = strings.ToLower(strings.TrimSpace(input))

		switch input {
		case "":
			fmt.Printf("\nDisplaying %d recent goals:\n", goalsToShow)
			break InputLoop
		case "all":
			fmt.Println("\nDisplaying all goals:")
			goalsToShow = totalGoals
			break InputLoop
		default:
			number, err := strconv.Atoi(input)
			if err != nil || number <= 0 || number >= totalGoals {
				fmt.Println("[Warning] Not a valid number, try again.")
				continue
			}
			fmt.Printf("\nDisplaying last %d goals:\n", number)
			goalsToShow = number
			break InputLoop
		}
	}
	return goalsToShow
}
