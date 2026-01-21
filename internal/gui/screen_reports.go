package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type ReportsScreen struct {
	guiApp *GuiApp
}

func NewReportsScreen(app *GuiApp) *ReportsScreen {
	return &ReportsScreen{guiApp: app}
}

func (r *ReportsScreen) Render() fyne.CanvasObject {
	title := widget.NewLabelWithStyle("Reports & Analytics", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	content := container.NewVBox(
		title,
		widget.NewSeparator(),
		widget.NewLabel("Reports and analytics page - Coming soon!"),
	)

	return container.NewScroll(content)
}
