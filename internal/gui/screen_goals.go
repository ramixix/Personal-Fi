package gui

import (
	"errors"
	"financial_tracker/internal/core"
	"financial_tracker/internal/models"
	"financial_tracker/internal/storage"
	"financial_tracker/internal/utils"
	"fmt"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type GoalsScreen struct {
	guiApp *GuiApp
}

func NewGoalsScreen(app *GuiApp) *GoalsScreen {
	return &GoalsScreen{guiApp: app}
}

func (g *GoalsScreen) Render() fyne.CanvasObject {

	header := g.createHeader()

	activeGoals := g.createActiveGoalSection()

	content := container.NewVBox(
		header,
		widget.NewSeparator(),
		activeGoals,
	)

	return container.NewScroll(content)
}

// --------------------------
//
//	Header
//
// --------------------------
// createHeader creates the header with stats and create button
func (g *GoalsScreen) createHeader() fyne.CanvasObject {
	title := widget.NewLabelWithStyle("🎯 Goals", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	activeGoals := len(core.GetActiveGoals())
	completedGoals := len(core.GetCompletedGoals())
	totalSaved := core.GetTotalGoalsSaved()

	statsLabel := widget.NewLabel(fmt.Sprintf("Active: %d  |  Completed: %d  |  Total Saved: %.2f", activeGoals, completedGoals, totalSaved))

	createBtn := widget.NewButton("➕ New Goal", func() { g.showCreateGoalDialog() })
	createBtn.Importance = widget.HighImportance

	head := container.NewBorder(nil, nil, title, statsLabel)
	return container.NewVBox(head, createBtn)
}

// ----------------------------------
//
//	Add Goal Dialog
//
// ----------------------------------
func (g *GoalsScreen) showCreateGoalDialog() {

	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("Goal name (e.g., Emergency Fund)")
	nameEntry.Validator = func(name string) error {
		if len(strings.TrimSpace(name)) < 2 {
			return errors.New("Name too short. Must at least 2 characters.")
		}
		return nil
	}

	descriptionEntry := widget.NewEntry()
	descriptionEntry.SetPlaceHolder("Description")

	targetAmountEntry := widget.NewEntry()
	targetAmountEntry.SetPlaceHolder("0.0")
	targetAmountEntry.Validator = func(value string) error {
		trimmed := strings.TrimSpace(value)
		if trimmed == "" {
			return errors.New("amount is required")
		}
		amount, err := strconv.ParseFloat(targetAmountEntry.Text, 64)
		if err != nil {
			return errors.New("Amount must be a number.")
		}
		if amount <= 0 {
			return errors.New("Amount can not be negative and must be greater than 0.")
		}
		return nil
	}

	categorySelect := widget.NewSelect([]string{"savings", "debt", "investment", "purchase", "other"}, nil)
	categorySelect.SetSelected("savings")

	prioritySelect := widget.NewSelect([]string{"high", "medium", "low"}, nil)
	prioritySelect.SetSelected("medium")

	hasDeadlineCheck := widget.NewCheck("Set a deadline", nil)

	dealineEntry := widget.NewEntry()
	dealineEntry.SetPlaceHolder("YYYY-MM-DD")
	dealineEntry.Disable()
	dealineEntry.Validator = func(s string) error {
		_, err := utils.ParseDate(dealineEntry.Text)
		if err != nil {
			return errors.New("Dealine must be in YYYY-MM-DD format.")
		}
		return nil
	}

	hasDeadlineCheck.OnChanged = func(checked bool) {
		if checked {
			dealineEntry.Enable()
		} else {
			dealineEntry.Disable()
		}
	}

	formItems := []*widget.FormItem{
		widget.NewFormItem("Name", nameEntry),
		widget.NewFormItem("Description", descriptionEntry),
		widget.NewFormItem("Target Amount", targetAmountEntry),
		widget.NewFormItem("Category", categorySelect),
		widget.NewFormItem("Priority", prioritySelect),
		widget.NewFormItem("", hasDeadlineCheck),
		widget.NewFormItem("Dealine", dealineEntry),
	}

	goalCreationDialog := dialog.NewForm("Create New Goals", "Create", "Cancel", formItems, func(confirmed bool) {
		if !confirmed {
			return
		}

		amount, err := strconv.ParseFloat(strings.TrimSpace(targetAmountEntry.Text), 64)
		name := strings.TrimSpace(nameEntry.Text)
		if err != nil || amount <= 0 || name == "" {
			dialog.ShowError(fmt.Errorf("Invalid data submitted"), g.guiApp.GuiWindow)
			return
		}

		var deadline time.Time
		if hasDeadlineCheck.Checked {
			deadline, _ = utils.ParseDate(dealineEntry.Text)
		}

		newGoal := models.Goal{
			ID:            storage.NextGoalID,
			Name:          nameEntry.Text,
			Description:   descriptionEntry.Text,
			TargetAmount:  amount,
			CurrentAmount: 0,
			HasDeadline:   hasDeadlineCheck.Checked,
			Deadline:      deadline,
			Category:      categorySelect.Selected,
			Priority:      prioritySelect.Selected,
			Status:        "active",
			Created:       time.Now(),
		}
		core.AddGoal(newGoal)
		storage.SaveData()

		dialog.ShowInformation("Success", "Goal created successfully! 🎯", g.guiApp.GuiWindow)
		g.guiApp.ShowGoalsScreen()
	}, g.guiApp.GuiWindow)

	goalCreationDialog.Resize(fyne.NewSize(450, 300))
	goalCreationDialog.Show()

}

// --------------------------------------
//
//	Active Goal Section
//
// --------------------------------------
func (g *GoalsScreen) createActiveGoalSection() fyne.CanvasObject {
	activeGoals := core.GetActiveGoals()

	sectionTitle := widget.NewLabelWithStyle("🔥 Active Goals", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	if len(activeGoals) == 0 {
		emptyState := container.NewVBox(
			widget.NewLabel("🎯"),
			widget.NewLabel("No active goals yet"),
			widget.NewLabel("Create your first goal to start tracking your progress!"),
		)
		return container.NewVBox(sectionTitle, emptyState)
	}

	var goalsCard []fyne.CanvasObject
	for _, goal := range activeGoals {
		card := g.createGoalCard(goal)
		goalsCard = append(goalsCard, card)
	}

	return container.NewVBox(
		sectionTitle,
		container.NewVBox(goalsCard...),
	)

}

// --------------------------------------
//
//	Create beautiful card for Goal
//
// --------------------------------------
func (g *GoalsScreen) createGoalCard(goal models.Goal) fyne.CanvasObject {
	priorityIcon := "🔵"
	switch goal.Priority {
	case "high":
		priorityIcon = "🔴 (High Priority)"
	case "medium":
		priorityIcon = "🟡 (Medium Priority)"
	case "low":
		priorityIcon = "🟢 (Low Priority)"
	}

	nameLabel := widget.NewLabelWithStyle(fmt.Sprintf("%s %s", priorityIcon, goal.Name), fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	var descriptionLabel *widget.Label
	if goal.Description != "" {
		descriptionLabel = widget.NewLabel(goal.Description)
	}

	progress := core.GetGoalProgress(goal)
	progressPercent := widget.NewLabelWithStyle(fmt.Sprintf("%.1f%% Completed", progress), fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	// progressBar := g.createProgressBar(progress)

	amountLabel := widget.NewLabel(fmt.Sprintf("$%.2f / $%.2f", goal.CurrentAmount, goal.TargetAmount))

	var deadlineLabel *widget.Label
	var monthlyRequiredContribution *widget.Label
	if goal.HasDeadline {
		remainingDays := core.GetRemainingDays(goal)
		if remainingDays >= 0 {
			deadlineLabel = widget.NewLabel(fmt.Sprintf("⏰ %d days remaining (Due: %s)", remainingDays, goal.Deadline.Format("Jan 02, 2006")))
		} else {
			deadlineLabel = widget.NewLabel(fmt.Sprintf("⚠️ Deadline passed (%s)", goal.Deadline.Format("Jan 02, 2006")))
		}

		requiredContribution := core.GetRequiredMonthlyContribution(goal)
		if requiredContribution > 0 {
			monthlyRequiredContribution = widget.NewLabel(fmt.Sprintf("💡 Need $%.2f/month to meet deadline", requiredContribution))
		}
	}

	// Build card content
	cardContent := container.NewVBox(nameLabel)
	if descriptionLabel != nil {
		cardContent.Add(descriptionLabel)
	}
	cardContent.Add(progressPercent)
	// cardContent.Add(progressBar)
	cardContent.Add(amountLabel)
	if deadlineLabel != nil {
		cardContent.Add(deadlineLabel)
	}
	if monthlyRequiredContribution != nil {
		cardContent.Add(monthlyRequiredContribution)
	}

	// Action buttons
	contributeBtn := widget.NewButton("➕ Contribute", func() {})
	contributeBtn.Importance = widget.HighImportance

	viewDetailsBtn := widget.NewButton("📊 Details", func() {})

	editBtn := widget.NewButton("Edit", func() {})
	editBtn.Importance = widget.LowImportance

	deleteBtn := widget.NewButton("Delete", func() {})
	deleteBtn.Importance = widget.DangerImportance

	actions := container.NewHBox(
		contributeBtn,
		viewDetailsBtn,
		editBtn,
		deleteBtn,
	)

	cardContent.Add(actions)

	card := widget.NewCard("", "", cardContent)
	return card
}
