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

	reportSelector := r.setReportType()

	content := container.NewVBox(
		header,
		widget.NewSeparator(),
		reportSelector,
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

// ------------------------------------------------
//
//	Create buttons to select and set report type
//
// ------------------------------------------------
func (r *ReportsScreen) setReportType() fyne.CanvasObject {
	sectionTitle := widget.NewLabel("Select Report Type:")

	overviewBtn := widget.NewButton("📈 Overview", func() {
		r.reportType = "overview"
		r.guiApp.ShowReportsScreen()
	})

	monthlyBtn := widget.NewButton("📅 Monthly Report", func() {
		r.reportType = "monthly"
		r.guiApp.ShowReportsScreen()
	})

	categoryBtn := widget.NewButton("🏷️ Category Breakdown", func() {
		r.reportType = "category"
		r.guiApp.ShowReportsScreen()
	})

	comparisonBtn := widget.NewButton("🔄 Comparison", func() {
		r.reportType = "comparison"
		r.guiApp.ShowReportsScreen()
	})

	trendsBtn := widget.NewButton("📉 Trends", func() {
		r.reportType = "trends"
		r.guiApp.ShowReportsScreen()
	})

	// hightlight the selected button
	switch r.reportType {
	case "overview":
		overviewBtn.Importance = widget.HighImportance
	case "monthly":
		monthlyBtn.Importance = widget.HighImportance
	case "category":
		categoryBtn.Importance = widget.HighImportance
	case "comparison":
		comparisonBtn.Importance = widget.HighImportance
	case "trends":
		trendsBtn.Importance = widget.HighImportance
	}

	buttonsGrid := container.NewGridWithColumns(5, overviewBtn, monthlyBtn, categoryBtn, comparisonBtn, trendsBtn)

	return container.NewVBox(sectionTitle, buttonsGrid)

}
