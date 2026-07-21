package gui

import (
	"errors"
	"financial_tracker/internal/core"
	"financial_tracker/internal/models"
	"financial_tracker/internal/utils"
	"fmt"
	"image/color"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
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

	completedGoals := g.createCompletedGoals()

	content := container.NewVBox(
		header,
		widget.NewSeparator(),
		activeGoals,
		widget.NewSeparator(),
		completedGoals,
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

	activeGoals := core.GetGoalsLength(models.StatusActive)
	completedGoals := core.GetGoalsLength(models.StatusCompleted)
	totalSaved := core.GetTotalGoalsAmountByCurrency()

	var statuText strings.Builder
	statuText.WriteString(fmt.Sprintf("Active: %d  |  Completed: %d\n", activeGoals, completedGoals))
	for currency, total := range totalSaved {
		statuText.WriteString(fmt.Sprintf("Total Saved in %s: %s", currency, utils.FormatCurrency(total, currency)))
	}

	statsLabel := widget.NewLabel(strings.TrimSpace(statuText.String()))

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
			return errors.New("Name too short. Must be at least 2 characters.")
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

	formattedCurrencies := utils.GetFormattedCurrencyOptions()
	currencySelect := widget.NewSelect(formattedCurrencies, func(value string) {})
	currencySelect.SetSelected("USD ($)")

	categorySelect := widget.NewSelect([]string{"savings", "debt", "investment", "purchase", "other"}, nil)
	categorySelect.SetSelected("savings")

	prioritySelect := widget.NewSelect([]string{string(models.HighPriority), string(models.MediumPriority), string(models.LowPriority)}, nil)
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
		widget.NewFormItem("Currency", currencySelect),
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

		if currencySelect.Selected == "" {
			dialog.ShowError(fmt.Errorf("Please select a currency."), g.guiApp.GuiWindow)
			return
		}

		var deadline time.Time
		if hasDeadlineCheck.Checked {
			deadline, _ = utils.ParseDate(dealineEntry.Text)
		}

		selectedCurrencyCode := currencySelect.Selected[:3]
		newGoal := models.Goal{
			ID:            utils.MustGenerateUUID(),
			Name:          nameEntry.Text,
			Description:   descriptionEntry.Text,
			TargetAmount:  amount,
			CurrentAmount: 0,
			CurrencyCode:  selectedCurrencyCode,
			HasDeadline:   hasDeadlineCheck.Checked,
			Deadline:      deadline,
			Category:      categorySelect.Selected,
			Priority:      prioritySelect.Selected,
			Status:        "active",
			Created:       time.Now(),
		}
		core.AddGoal(newGoal)

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
	case string(models.HighPriority):
		priorityIcon = "🔴 (High Priority)"
	case string(models.MediumPriority):
		priorityIcon = "🟡 (Medium Priority)"
	case string(models.LowPriority):
		priorityIcon = "🟢 (Low Priority)"
	}

	nameLabel := widget.NewLabelWithStyle(fmt.Sprintf("%s %s", priorityIcon, goal.Name), fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	var descriptionLabel *widget.Label
	if goal.Description != "" {
		descriptionLabel = widget.NewLabel(goal.Description)
	}

	progress := core.GetGoalProgress(goal)
	progressPercent := widget.NewLabelWithStyle(fmt.Sprintf("%.1f%% Completed", progress), fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	progressBar := g.createProgressBar(progress)

	amountLabel := widget.NewLabel(fmt.Sprintf("%s / %s",
		utils.FormatCurrency(goal.CurrentAmount, goal.CurrencyCode),
		utils.FormatCurrency(goal.TargetAmount, goal.CurrencyCode)))

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
			monthlyRequiredContribution = widget.NewLabel(fmt.Sprintf("💡 Need %.2f/month to meet deadline", requiredContribution))
		}
	}

	// Build card content
	cardContent := container.NewVBox(nameLabel)
	if descriptionLabel != nil {
		cardContent.Add(descriptionLabel)
	}
	cardContent.Add(progressPercent)
	cardContent.Add(progressBar)
	cardContent.Add(amountLabel)
	if deadlineLabel != nil {
		cardContent.Add(deadlineLabel)
	}
	if monthlyRequiredContribution != nil {
		cardContent.Add(monthlyRequiredContribution)
	}

	// Action buttons
	contributeBtn := widget.NewButton("➕ Contribute", func() { g.showContributeDialog(goal) })
	contributeBtn.Importance = widget.HighImportance

	viewDetailsBtn := widget.NewButton("📊 Details", func() { g.showGoalDetails(goal) })

	editBtn := widget.NewButton("Edit", func() { g.showEditGoalDialog(goal) })
	editBtn.Importance = widget.WarningImportance

	deleteBtn := widget.NewButton("Delete", func() { g.showDeleteGoalDialog(goal) })
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

// -------------------------------------------
//
//	Show progress bar for each active goal
//
// -------------------------------------------
type ratioLayout struct {
	ratio float64
}

func (r *ratioLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	// The first object is the green bar
	objects[0].Resize(fyne.NewSize(size.Width*float32(r.ratio), size.Height))
	objects[0].Move(fyne.NewPos(0, 0))
}

func (r *ratioLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	return fyne.NewSize(0, 15) // Minimum height of 15px
}

func (g *GoalsScreen) createProgressBar(progressPercentage float64) fyne.CanvasObject {
	if progressPercentage < 0 {
		progressPercentage = 0
	}
	if progressPercentage > 100 {
		progressPercentage = 100
	}
	ratio := progressPercentage / 100.0

	background := canvas.NewRectangle(color.RGBA{R: 180, G: 180, B: 180, A: 160})
	background.SetMinSize(fyne.NewSize(0, 15))

	var progressColor color.Color
	if progressPercentage < 25 {
		progressColor = color.RGBA{R: 239, G: 68, B: 68, A: 255} // Red
	} else if progressPercentage < 50 {
		progressColor = color.RGBA{R: 251, G: 191, B: 36, A: 255} // Orange
	} else if progressPercentage < 75 {
		progressColor = color.RGBA{R: 234, G: 179, B: 8, A: 255} // Yellow
	} else if progressPercentage < 100 {
		progressColor = color.RGBA{R: 34, G: 197, B: 94, A: 255} // Green
	} else {
		progressColor = color.RGBA{R: 16, G: 185, B: 129, A: 255} // Bright Green
	}

	foreground := canvas.NewRectangle(progressColor)
	progressLayer := container.New(&ratioLayout{ratio: ratio}, foreground)

	progressText := widget.NewLabelWithStyle(fmt.Sprintf("%.2f%%", progressPercentage), fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	return container.NewStack(background, progressLayer, container.NewCenter(progressText))
}

// -----------------------
//
//	Show Completed Goals
//
// -----------------------
func (g *GoalsScreen) createCompletedGoals() fyne.CanvasObject {
	completedGoals := core.GetCompletedGoals()

	title := widget.NewLabelWithStyle("✅ Completed Goals", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	if len(completedGoals) == 0 {
		return container.NewVBox(title, widget.NewLabel("No completed goals yet. Keep working on your active goals!"))
	}

	var cards []fyne.CanvasObject
	for _, goal := range completedGoals {
		goalCard := g.createCompletedGoalCard(goal)
		cards = append(cards, goalCard)
	}

	return container.NewVBox(title, container.NewVBox(cards...))
}

// -------------------------------------
//
//	Create Card for Completed Goals
//
// -------------------------------------
func (g *GoalsScreen) createCompletedGoalCard(goal models.Goal) fyne.CanvasObject {
	nameLabel := widget.NewLabelWithStyle(fmt.Sprintf("🏆 %s", goal.Name), fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	infoLabel := widget.NewLabel(
		fmt.Sprintf("Completed: %s | Amount: %s", goal.CompletedDate.Format("Jan 02, 2006"), utils.FormatCurrency(goal.CurrentAmount, goal.CurrencyCode)))

	viewBtn := widget.NewButton("View Details", func() { g.showGoalDetails(goal) })
	deleteBtn := widget.NewButton("Delete", func() { g.showDeleteGoalDialog(goal) })
	deleteBtn.Importance = widget.DangerImportance

	buttons := container.NewHBox(viewBtn, deleteBtn)

	content := container.NewVBox(nameLabel, infoLabel, buttons)
	return widget.NewCard("", "", content)
}

// -------------------------------------------
//
//	Dialog to Add Contribution to a Goal
//
// -------------------------------------------
func (g *GoalsScreen) showContributeDialog(goal models.Goal) {
	amountEntry := widget.NewEntry()
	amountEntry.SetPlaceHolder("0.0")
	amountEntry.Validator = func(value string) error {
		amount, err := strconv.ParseFloat(strings.TrimSpace(amountEntry.Text), 64)
		if err != nil {
			return errors.New("Not valid amount")
		}
		if amount <= 0 {
			return errors.New("Amount must be positive (greater than 0)")
		}
		return nil
	}

	noteEntry := widget.NewEntry()
	noteEntry.SetPlaceHolder("Add Note (Optional)")

	formItems := []*widget.FormItem{
		widget.NewFormItem("Amount", amountEntry),
		widget.NewFormItem("Note", noteEntry),
	}

	contributeDilog := dialog.NewForm(fmt.Sprintf("Contribute to: %s", goal.Name), "Contribute", "Cancel", formItems,
		func(confirmed bool) {
			if !confirmed {
				return
			}

			amount, err := strconv.ParseFloat(strings.TrimSpace(amountEntry.Text), 64)
			if err != nil || amount <= 0 {
				dialog.ShowError(errors.New("Invalid data submitted"), g.guiApp.GuiWindow)
				return
			}

			success := core.AddGoalContribution(goal.ID, amount, strings.TrimSpace(noteEntry.Text), false)
			if !success {
				dialog.ShowError(errors.New("failed to add contribution"), g.guiApp.GuiWindow)
				return
			}

			updatedGoal := core.FindGoal(goal.ID)
			if updatedGoal != nil && updatedGoal.Status == "complete" {
				g.showGoalCompletedCelebration(updatedGoal)
			} else {
				dialog.ShowInformation("Success", fmt.Sprintf("Added %.2f(%s) to %s! 🎉", amount, goal.CurrencyCode, goal.Name), g.guiApp.GuiWindow)
			}
			g.guiApp.ShowGoalsScreen()
		}, g.guiApp.GuiWindow)

	contributeDilog.Resize(fyne.NewSize(450, 300))
	contributeDilog.Show()
}

// -----------------------------------------------------------------------------
//
//	show congragratulation message if contribution end up finishing a goal
//
// -----------------------------------------------------------------------------
func (g *GoalsScreen) showGoalCompletedCelebration(goal *models.Goal) {
	message := fmt.Sprintf("🎉🎊 CONGRATULATIONS! 🎊🎉\nYou've completed your goal:\n%s", goal.Name)
	celebrationText := widget.NewLabelWithStyle(message, fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	amountText := widget.NewLabel(fmt.Sprintf("Amount achieved: $%.2f", goal.TargetAmount))

	content := container.NewVBox(celebrationText, amountText, widget.NewLabel("Keep up the great work!"))

	dialog.ShowCustom("Goal Completed!", "Close", content, g.guiApp.GuiWindow)
}

// ---------------------------------
//
//	Showing Goal Details Dialog
//
// ---------------------------------
func (g *GoalsScreen) showGoalDetails(goal models.Goal) {
	details := fmt.Sprintf("Goal Name: %s\nStatus: %s\nCategory: %s\nPriority: %s\n\nProgress: %.1f%%\nCurrent Amount: %.2f\nTarget Amount: %.2f,\n Currency: %s\nReamining: %.2f",
		goal.Name, goal.Status, goal.Category, goal.Priority, core.GetGoalProgress(goal), goal.CurrentAmount, goal.TargetAmount, goal.CurrencyCode, goal.TargetAmount-goal.CurrentAmount)

	if goal.Description != "" {
		details += fmt.Sprintf("\nDescription: %s\n", goal.Description)
	}

	if goal.HasDeadline {
		details += fmt.Sprintf("\nDeadline: %s\n", goal.Deadline.Format("Jan 02, 2006"))
		remainingDays := core.GetRemainingDays(goal)
		if remainingDays > 0 {
			details += fmt.Sprintf("Days remaining: %d\n", remainingDays)
		}
	}

	details += fmt.Sprintf("\nCreated: %s\n", goal.Created.Format("Jan 02, 2006"))

	if goal.Status == string(models.StatusCompleted) {
		details += fmt.Sprintf("Completed: %s\n", goal.CompletedDate.Format("Jan 02, 2006"))
	}

	contributionslist := core.GetSpecificGoalContributions(goal.ID)
	contributionCount := len(contributionslist)
	if contributionCount > 0 {
		details += fmt.Sprintf("\n--- Recent Contributions (%d total) ---\n", contributionCount)
		recent := 10
		for i := 0; i < recent; i++ {
			details += fmt.Sprintf("%s: +%.2f - %s\n", contributionslist[i].Date.Format("Jan 02, 2006"),
				contributionslist[i].Amount,
				contributionslist[i].Note)
		}
	}

	detailsLabel := widget.NewLabel(details)
	scrollableContent := container.NewScroll(detailsLabel)
	scrollableContent.SetMinSize(fyne.NewSize(500, 400))

	dialog.ShowCustom("Goal Details", "Close", scrollableContent, g.guiApp.GuiWindow)
}

// ------------------------
//
//	Edit Goal Dialog
//
// ------------------------
func (g *GoalsScreen) showEditGoalDialog(goal models.Goal) {
	nameEntry := widget.NewEntry()
	nameEntry.SetText(goal.Name)
	nameEntry.Validator = func(name string) error {
		if len(strings.TrimSpace(name)) < 2 {
			return errors.New("Name too short. Must at least 2 characters.")
		}
		return nil
	}

	descriptionEntry := widget.NewEntry()
	descriptionEntry.SetText(goal.Description)

	targetAmountEntry := widget.NewEntry()
	targetAmountEntry.SetText(fmt.Sprintf("%f", goal.TargetAmount))
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

	formattedCurrencies := utils.GetFormattedCurrencyOptions()
	currencySelect := widget.NewSelect(formattedCurrencies, func(value string) {})
	currentSelection := goal.CurrencyCode
	if symbol, exists := utils.CurrencySymbols[goal.CurrencyCode]; exists && symbol != "" {
		currentSelection = fmt.Sprintf("%s (%s)", goal.CurrencyCode, symbol)
	}
	currencySelect.SetSelected(currentSelection)

	categorySelect := widget.NewSelect([]string{"savings", "debt", "investment", "purchase", "other"}, nil)
	categorySelect.SetSelected(goal.Category)

	prioritySelect := widget.NewSelect([]string{string(models.HighPriority), string(models.MediumPriority), string(models.LowPriority)}, nil)
	prioritySelect.SetSelected(goal.Priority)

	hasDeadlineCheck := widget.NewCheck("Set a deadline", nil)

	dealineEntry := widget.NewEntry()
	dealineEntry.SetPlaceHolder("YYYY-MM-DD")
	if goal.HasDeadline {
		dealineEntry.SetText(goal.Deadline.Format("2006-01-02"))
	}
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
		widget.NewFormItem("Currency", currencySelect),
		widget.NewFormItem("Category", categorySelect),
		widget.NewFormItem("Priority", prioritySelect),
		widget.NewFormItem("", hasDeadlineCheck),
		widget.NewFormItem("Dealine", dealineEntry),
	}

	goalCreationDialog := dialog.NewForm("Edit Goals", "Edit", "Cancel", formItems, func(confirmed bool) {
		if !confirmed {
			return
		}

		amount, err := strconv.ParseFloat(strings.TrimSpace(targetAmountEntry.Text), 64)
		name := strings.TrimSpace(nameEntry.Text)
		if err != nil || amount <= 0 || name == "" {
			dialog.ShowError(fmt.Errorf("Invalid data submitted"), g.guiApp.GuiWindow)
			return
		}

		if currencySelect.Selected == "" {
			dialog.ShowError(fmt.Errorf("Please select a currency."), g.guiApp.GuiWindow)
			return
		}

		var deadline time.Time
		if hasDeadlineCheck.Checked {
			deadline, _ = utils.ParseDate(dealineEntry.Text)
		}

		goalPointer := core.FindGoal(goal.ID)
		goalPointer.Name = nameEntry.Text
		goalPointer.Description = descriptionEntry.Text
		goalPointer.TargetAmount = amount
		goalPointer.CurrencyCode = currencySelect.Selected[:3]
		if hasDeadlineCheck.Checked {
			goalPointer.HasDeadline = hasDeadlineCheck.Checked
			goalPointer.Deadline = deadline
		}
		goalPointer.Category = categorySelect.Selected
		goalPointer.Priority = prioritySelect.Selected
		if goalPointer.CurrentAmount > amount {
			goalPointer.Status = "complete"
		}

		err = core.UpdateGoal(*goalPointer)
		if err != nil {
			dialog.ShowError(fmt.Errorf("Could not update goal"), g.guiApp.GuiWindow)
			return
		}

		dialog.ShowInformation("Success", "Goal edited successfully! 🎯", g.guiApp.GuiWindow)
		g.guiApp.ShowGoalsScreen()
	}, g.guiApp.GuiWindow)

	goalCreationDialog.Resize(fyne.NewSize(450, 300))
	goalCreationDialog.Show()
}

// ----------------------------
//
//	Dialog to delete a goal
//
// ----------------------------
func (g *GoalsScreen) showDeleteGoalDialog(goal models.Goal) {
	d := dialog.NewConfirm("Delete Goal",
		fmt.Sprintf("Are you sure you want to delete '%s'?\nThis will also delete all contribution history of this goal.", goal.Name),
		func(confirmed bool) {
			if !confirmed {
				return
			}
			err := core.DeleteGoal(goal.ID)
			if err == nil {
				dialog.ShowInformation("Successfull", "Goal has been deleted successfully", g.guiApp.GuiWindow)
			} else {
				dialog.ShowError(fmt.Errorf("Could not remove selected goal"), g.guiApp.GuiWindow)
			}
			g.guiApp.ShowGoalsScreen()

		},
		g.guiApp.GuiWindow)
	d.Show()
}
