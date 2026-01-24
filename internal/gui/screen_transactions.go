package gui

import (
	"financial_tracker/internal/core"
	"financial_tracker/internal/storage"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type TransactionsScreen struct {
	guiApp *GuiApp

	// values needed for filtering
	searchText     string
	filterType     string
	filterCategory string
	filterPeriod   string
}

func NewTransactionsScreen(app *GuiApp) *TransactionsScreen {
	return &TransactionsScreen{guiApp: app}
}

func (t *TransactionsScreen) Render() fyne.CanvasObject {
	// Header with title and stats and simple add button
	header := t.createHeader()

	// filter section
	filterBar := t.createFilterBar()

	content := container.NewVBox(
		header,
		widget.NewSeparator(),
		filterBar,
	)

	return container.NewScroll(content)
}

func (t *TransactionsScreen) createHeader() fyne.CanvasObject {
	title := widget.NewLabelWithStyle("💰 Transactions", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	// Quick stats
	totalIncome, totalExpenses := core.CalculateTotals()
	totalCount := len(storage.Transactions)
	stats := fmt.Sprintf("Total: %d transactions  |  Income: $%.2f  |  Expenses: $%.2f", totalCount, totalIncome, totalExpenses)

	statLabel := widget.NewLabel(stats)

	// Add transaction button
	transacAddBtn := widget.NewButton("Add Transaction", func() {})
	transacAddBtn.Importance = widget.HighImportance

	// titleRow := container.NewGridWithColumns(2, title, statLabel)
	titleRow := container.NewBorder(nil, nil, title, statLabel)
	content := container.NewVBox(titleRow, transacAddBtn)
	return content
}

func (t *TransactionsScreen) createFilterBar() fyne.CanvasObject {
	// Search Entery
	searchEntry := widget.NewEntry()
	searchEntry.SetPlaceHolder("🔍 Search transactions...")
	searchEntry.OnChanged = func(text string) { t.searchText = text }

	// Filtering Transactions by Type of transaction
	typeSelect := widget.NewSelect([]string{"All Types", "Income", "Expenses"},
		func(value string) {
			switch value {
			case "Income":
				t.filterType = "income"
			case "Expneses":
				t.filterType = "expense"
			default:
				t.filterType = "all"
			}
		})
	typeSelect.SetSelected("All Types")

	// Filtering transactions by period => week, month, year, all time
	periodSelect := widget.NewSelect([]string{"All Time", "This Week", "This Month", "This Year"},
		func(value string) {
			switch value {
			case "This Week":
				t.filterPeriod = "week"
			case "This Month":
				t.filterPeriod = "month"
			case "This Year":
				t.filterPeriod = "year"
			default:
				t.filterPeriod = "all"
			}
		})
	periodSelect.SetSelected("All Time")

	// filtering by categories
	categories := []string{"All Categories"}
	categories = append(categories, core.GetCategories()...)
	categorySelect := widget.NewSelect(categories,
		func(value string) {
			if value == "All Categories" {
				t.filterCategory = "all"
			} else {
				t.filterCategory = value
			}
		})
	categorySelect.SetSelected("All Categories")

	// button to clear all filters and set default selects for selections
	clearBtn := widget.NewButton("Clear Filters", func() {
		t.searchText = ""
		t.filterType = "all"
		t.filterCategory = "all"
		t.filterPeriod = "all"
		searchEntry.SetText("")
		typeSelect.SetSelected("All Types")
		periodSelect.SetSelected("All Time")
		categorySelect.SetSelected("All Categories")
	})

	filterRow1 := container.NewGridWithColumns(2, searchEntry, clearBtn)
	filterRow2 := container.NewGridWithColumns(3, typeSelect, periodSelect, categorySelect)

	return container.NewVBox(filterRow1, filterRow2)
}
