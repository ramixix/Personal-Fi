package gui

import (
	"financial_tracker/internal/core"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
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

func (g *GoalsScreen) createHeader() fyne.CanvasObject {
	title := widget.NewLabelWithStyle("🎯 Goals", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	activeGoals := len(core.GetActiveGoals())
	completedGoals := len(core.GetCompletedGoals())
	totalSaved := core.GetTotalGoalsSaved()

	statsLabel := widget.NewLabel(fmt.Sprintf("Active: %d  |  Completed: %d  |  Total Saved: %.2f", activeGoals, completedGoals, totalSaved))

	createBtn := widget.NewButton("➕ New Goal", func() {})

	head := container.NewBorder(nil, nil, title, statsLabel)
	return container.NewVBox(head, createBtn)
}
