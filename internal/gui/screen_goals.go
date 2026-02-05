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

	content := container.NewVBox(
		header,
		widget.NewSeparator(),
		widget.NewLabel("Goals tracking page - Coming soon!"),
	)

	return container.NewScroll(content)
}

// createHeader creates the header with stats and create button
func (g *GoalsScreen) createHeader() fyne.CanvasObject {
	title := widget.NewLabelWithStyle("🎯 Goals", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	activeGoals := len(core.GetActiveGoals())
	completedGoals := len(core.GetCompletedGoals())
	totalSaved := core.GetTotalGoalsSaved()

	statsLabel := widget.NewLabel(fmt.Sprintf("Active: %d  |  Completed: %d  |  Total Saved: %.2f", activeGoals, completedGoals, totalSaved))

	createBtn := widget.NewButton("➕ New Goal", func() { g.showCreateGoalDialog() })

	head := container.NewBorder(nil, nil, title, statsLabel)
	return container.NewVBox(head, createBtn)
}

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
