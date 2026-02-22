package gui

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type ReportsScreen struct {
	guiApp *GuiApp

	reportType    string // "monthly", "category", "comparison", "trends"
	selectedYear  int
	selectedMonth time.Month
}

func NewReportsScreen(app *GuiApp) *ReportsScreen {
	return &ReportsScreen{
		guiApp:        app,
		reportType:    "overview",
		selectedYear:  time.Now().Year(),
		selectedMonth: time.Now().Month(),
	}
}

func (r *ReportsScreen) Render() fyne.CanvasObject {
	header := r.createHeader()

	content := container.NewVBox(
		header,
		widget.NewSeparator(),
		widget.NewLabel("Reports and analytics page - Coming soon!"),
	)

	return container.NewScroll(content)
}

// ----------------------------------------------
//
//	Create a simple header for report screen
//
// ----------------------------------------------
func (r *ReportsScreen) createHeader() fyne.CanvasObject {
	title := widget.NewLabelWithStyle("📊 Reports & Analytics", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	subtitle := widget.NewLabel("Analyze your financial data and discover insights")
	return container.NewVBox(title, subtitle)
}
