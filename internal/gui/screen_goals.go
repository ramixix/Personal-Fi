package gui

import (
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
	title := widget.NewLabelWithStyle("Goals", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	content := container.NewVBox(
		title,
		widget.NewSeparator(),
		widget.NewLabel("Goals tracking page - Coming soon!"),
	)

	return container.NewScroll(content)
}
