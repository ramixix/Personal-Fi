package gui

import (
	"financial_tracker/internal/core"
	"financial_tracker/internal/models"
	"financial_tracker/internal/storage"
	"fmt"
	"strconv"
	"strings"
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

	reportContent := r.createReportContent()

	content := container.NewVBox(
		header,
		widget.NewSeparator(),
		reportSelector,
	)

	return container.NewScroll(container.NewBorder(content, nil, nil, nil, reportContent))
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

// --------------------------------------------------------------------
//
//	Create the main report content based on report type selection
//
// -------------------------------------------------------------------
func (r *ReportsScreen) createReportContent() fyne.CanvasObject {
	switch r.reportType {
	case "monthly":
		return r.createMonthlyReport()
	case "category":
		return r.createCategoryReport()
	case "comparison":
		// return r.createComparisonReport()
		return nil
	case "trends":
		// return r.createTrendsReport()
		return nil
	default:
		return r.createOverviewReport()
	}
}

// --------------------------------------------
//
//	Create the overview report dashboard
//
// --------------------------------------------
func (r *ReportsScreen) createOverviewReport() fyne.CanvasObject {
	title := widget.NewLabelWithStyle("Financial Overview", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	totalIncome, totalExpenses := core.CalculateTotals()
	totalNet := totalIncome - totalExpenses
	totalAccounts := core.GetTotalAccountBalance()

	// total income, outcome and accounts summary cards
	totalsCard1 := widget.NewCard("💰 Total Income", "", widget.NewLabel(fmt.Sprintf("%.2f", totalIncome)))
	totalsCard2 := widget.NewCard("💸 Total Expenses", "", widget.NewLabel(fmt.Sprintf("%.2f", totalExpenses)))
	totalsCard3 := widget.NewCard("📊 Net Worth", "", widget.NewLabel(fmt.Sprintf("%.2f", totalNet)))
	totalsCard4 := widget.NewCard("🏦 In Accounts", "", widget.NewLabel(fmt.Sprintf("%.2f", totalAccounts)))

	totalsCardGrid := container.NewGridWithColumns(2, totalsCard1, totalsCard2, totalsCard3, totalsCard4)

	// monthly average summary section
	avgIncome, avgExpenses := core.GetMonthlyAverage()
	averageSection := widget.NewLabelWithStyle("Monthly Averages", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	avgCard1 := widget.NewCard("", "Average Income", widget.NewLabel(fmt.Sprintf("%.2f", avgIncome)))
	avgCard2 := widget.NewCard("", "Average Expenses", widget.NewLabel(fmt.Sprintf("%.2f", avgExpenses)))
	avgCard3 := widget.NewCard("", "Average Net", widget.NewLabel(fmt.Sprintf("%.2f", avgIncome-avgExpenses)))

	avgCardGrid := container.NewGridWithColumns(3, avgCard1, avgCard2, avgCard3)

	// active and complete Gaols summary
	activeGoals := len(core.GetActiveGoals())
	completedGoals := len(core.GetCompletedGoals())
	goalsSection := widget.NewLabelWithStyle("Goals Summary", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	goalsCard1 := widget.NewCard("", "Total Goals", widget.NewLabel(fmt.Sprintf("%d", len(storage.Goals))))
	goalsCard2 := widget.NewCard("", "Active Goals", widget.NewLabel(fmt.Sprintf("%d", activeGoals)))
	goalsCard3 := widget.NewCard("", "Completed Goals", widget.NewLabel(fmt.Sprintf("%d", completedGoals)))

	goalsCardGrid := container.NewGridWithColumns(3, goalsCard1, goalsCard2, goalsCard3)

	// total income and outcome transactions count summary
	incomeCount := 0
	expenseCount := 0
	for _, transac := range storage.Transactions {
		if transac.Type == "income" {
			incomeCount++
		} else {
			expenseCount++
		}
	}

	transacSection := widget.NewLabelWithStyle("Transaction Statistics", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	transacCard1 := widget.NewCard("", "Total Transactions", widget.NewLabel(fmt.Sprintf("%d", len(storage.Transactions))))
	transacCard2 := widget.NewCard("", "Income Entries", widget.NewLabel(fmt.Sprintf("%d", incomeCount)))
	transacCard3 := widget.NewCard("", "Expense Entries", widget.NewLabel(fmt.Sprintf("%d", expenseCount)))

	transacCardGrid := container.NewGridWithColumns(3, transacCard1, transacCard2, transacCard3)

	// Top categories
	var incomeCategory models.CategoryReport = core.GetCategoryBreakdown("income")[0]
	var expenseCategory models.CategoryReport = core.GetCategoryBreakdown("expense")[0]
	topCategoriesSection := widget.NewLabelWithStyle("Top Income & Expenses Categories", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	incomeCardMessage := fmt.Sprintf("Category Name: %s", incomeCategory.Category)
	expenseCardMessage := fmt.Sprintf("Category Name: %s", expenseCategory.Category)
	incomeCard := widget.NewCard("", incomeCardMessage, widget.NewLabel(fmt.Sprintf("%.2f (Count: %d | %.2f Total Percent of Incomes)", incomeCategory.Amount, incomeCategory.Count, incomeCategory.Percent)))
	expenseCard := widget.NewCard("", expenseCardMessage, widget.NewLabel(fmt.Sprintf("%2.f (Count: %d | %.2f Total Percent of Expense)", expenseCategory.Amount, expenseCategory.Count, expenseCategory.Percent)))

	categoryCardGrid := container.NewGridWithColumns(2, incomeCard, expenseCard)

	content := container.NewVBox(
		title,
		totalsCardGrid,
		widget.NewSeparator(),
		averageSection,
		avgCardGrid,
		widget.NewSeparator(),
		goalsSection,
		goalsCardGrid,
		widget.NewSeparator(),
		transacSection,
		transacCardGrid,
		widget.NewSeparator(),
		topCategoriesSection,
		categoryCardGrid,
	)
	return content
}

// ---------------------------------------------
//
//	creates monthly financial report screen
//
// ---------------------------------------------
func (r *ReportsScreen) createMonthlyReport() fyne.CanvasObject {
	monthlyReportTitle := widget.NewLabelWithStyle("Monthly Financial Report", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	currentYear := time.Now().Year()
	var decadeList []string
	for year := currentYear; year >= currentYear-10; year-- {
		decadeList = append(decadeList, fmt.Sprintf("%d", year))
	}

	yearSelectEntry := widget.NewSelect(decadeList, nil)
	yearSelectEntry.SetSelected(fmt.Sprintf("%d", r.selectedYear))
	yearSelectEntry.OnChanged = func(year string) {
		yearInt, err := strconv.Atoi(year)
		if err == nil {
			r.selectedYear = yearInt
			r.guiApp.ShowReportsScreen()
		}
	}

	months := []string{"January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"}
	monthSelect := widget.NewSelect(months, nil)
	monthSelect.SetSelected(r.selectedMonth.String())
	monthSelect.OnChanged = func(monthSelected string) {
		for index, m := range months {
			if m == monthSelected {
				r.selectedMonth = time.Month(index + 1)
				break
			}
		}
		r.guiApp.ShowReportsScreen()
	}

	selectorGrid := container.NewGridWithColumns(2,
		container.NewBorder(nil, nil, widget.NewLabel("Year:"), nil, yearSelectEntry),
		container.NewBorder(nil, nil, widget.NewLabel("Month:"), nil, monthSelect),
	)

	report := core.GetSpecificMonthYearReport(r.selectedYear, r.selectedMonth)

	reportTitle := fmt.Sprintf("%s %d Report", r.selectedMonth.String(), r.selectedYear)
	reportContent := widget.NewLabelWithStyle(
		fmt.Sprintf("Income:\t%-10.2f\nExpenses:\t%-10.2f\nNet:\t%-10.2f", report.Income, report.Expenses, report.Net),
		fyne.TextAlignLeading,
		fyne.TextStyle{Bold: true})
	reportCard := widget.NewCard(reportTitle, fmt.Sprintf("Transactions Number:\t%d", report.TxCount), reportContent)

	// Category Section
	var specificYearMonthTransactions []models.Transaction
	for _, transac := range storage.Transactions {
		if transac.Date.Year() == r.selectedYear && transac.Date.Month() == r.selectedMonth {
			specificYearMonthTransactions = append(specificYearMonthTransactions, transac)
		}
	}

	categoriesAndAmounts := make(map[string]float64)
	for _, transac := range specificYearMonthTransactions {
		if transac.Type == "income" {
			categoriesAndAmounts[transac.Category] += transac.Amount
		} else {
			categoriesAndAmounts[transac.Category] -= transac.Amount
		}
	}

	var categoryWidgetList []fyne.CanvasObject
	for cat, amount := range categoriesAndAmounts {
		cateType := "Income"
		if amount < 0 {
			cateType = "Expense"
		}
		lable := widget.NewLabelWithStyle(fmt.Sprintf("(%8s) %-10s: %-10.2f", cateType, cat, amount), fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
		categoryWidgetList = append(categoryWidgetList, lable)
	}
	categoryCard := widget.NewCard(fmt.Sprintf("%s %d Category Breakdown", r.selectedMonth.String(), r.selectedYear), "", container.NewVBox(categoryWidgetList...))

	content := container.NewVBox(
		monthlyReportTitle,
		selectorGrid,
		widget.NewSeparator(),
		reportCard,
		widget.NewSeparator(),
		categoryCard,
	)

	return content
}

// ---------------------------------------
//
//	category breakdown report screen
//
// ---------------------------------------
func (r *ReportsScreen) createCategoryReport() fyne.CanvasObject {
	sectionTitle := widget.NewLabelWithStyle("Category Breakdown", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	var table *widget.Table
	categoryReports := core.GetCategoryBreakdown("expense")

	categoryTypeSelection := widget.NewSelect([]string{"Income", "Expenses", "Both"}, nil)
	categoryTypeSelection.SetSelected("Expenses")
	categoryTypeSelection.OnChanged = func(selected string) {
		var transactionType string
		switch selected {
		case "Income":
			transactionType = "income"
		case "Expenses":
			transactionType = "expense"
		default:
			transactionType = ""
		}

		categoryReports = core.GetCategoryBreakdown(transactionType)
		table.Refresh()

	}

	if len(categoryReports) == 0 {
		return container.NewVBox(
			sectionTitle,
			categoryTypeSelection,
			widget.NewLabel("No data available"),
		)
	}

	table = widget.NewTable(
		func() (int, int) {
			return len(categoryReports) + 1, 5 // rows + 1(one extra for header), columns
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(id widget.TableCellID, obj fyne.CanvasObject) {
			label := obj.(*widget.Label)

			// Header row
			if id.Row == 0 {
				headers := []string{"No.", "Category", "Amount", "Count", "Percentage + Visual Bar"}
				label.SetText(headers[id.Col])
				label.TextStyle = fyne.TextStyle{Bold: true}
				return
			}

			category := categoryReports[id.Row-1]

			switch id.Col {
			case 0:
				label.SetText(fmt.Sprintf("%d", id.Row))
			case 1:
				label.SetText(category.Category)
			case 2:
				label.SetText(fmt.Sprintf("%.2f", category.Amount))
			case 3:
				label.SetText(fmt.Sprintf("%d", category.Count))
			case 4:
				barLength := int(category.Percent / 4) // Scale down
				if barLength > 25 {
					barLength = 25
				}
				bar := strings.Repeat("█", barLength)
				label.SetText(fmt.Sprintf("%-5.1f%%   %s", category.Percent, bar))

			}
		},
	)

	table.SetColumnWidth(0, 50)
	table.SetColumnWidth(1, 200)
	table.SetColumnWidth(2, 120)
	table.SetColumnWidth(3, 100)
	table.SetColumnWidth(4, 120)

	header := container.NewVBox(
		sectionTitle,
		categoryTypeSelection,
		widget.NewSeparator(),
	)

	return container.NewBorder(
		header, // top
		nil,    // bottom
		nil,    // left
		nil,    // right
		table,  // center -> fills all remaining space
	)
}
