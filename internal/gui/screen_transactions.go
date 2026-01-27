package gui

import (
	"errors"
	"financial_tracker/internal/core"
	"financial_tracker/internal/models"
	"financial_tracker/internal/storage"
	"fmt"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
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
	return &TransactionsScreen{guiApp: app, filterType: "all", filterCategory: "all", filterPeriod: "all"}
}

func (t *TransactionsScreen) Render() fyne.CanvasObject {
	// Header with title and stats and simple add button
	header := t.createHeader()

	// filter section
	filterBar := t.createFilterBar()

	// section to list transaction based on filters
	// transactionList := t.createTransactionList()

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
	transacAddBtn := widget.NewButton("Add Transaction", func() { t.showAddTransactionDialog() })
	transacAddBtn.Importance = widget.HighImportance

	// titleRow := container.NewGridWithColumns(2, title, statLabel)
	titleRow := container.NewBorder(nil, nil, title, statLabel)
	content := container.NewVBox(titleRow, transacAddBtn)
	return content
}

// dialog to add transactions
func (t *TransactionsScreen) showAddTransactionDialog() {
	// create form fields
	typeSelect := widget.NewSelect([]string{"income", "expense"}, nil)
	typeSelect.SetSelected("income")

	amountEntery := widget.NewEntry()
	amountEntery.SetPlaceHolder("0.00")

	amountEntery.Validator = func(value string) error {
		val, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return errors.New("Not a valid number!")
		}
		if val <= 0 {
			return errors.New("Amount must be positive (greater than 0)")
		}
		return nil
	}

	categorySelectEntry := widget.NewSelectEntry(core.GetCategories())
	categorySelectEntry.SetPlaceHolder("Category")

	categorySelectEntry.Validator = func(value string) error {
		trimCat := strings.TrimSpace(value)
		if len(trimCat) <= 2 {
			return errors.New("Category name too short!")
		}
		return nil
	}

	descriptionEntery := widget.NewEntry()
	descriptionEntery.SetPlaceHolder("Description")

	items := []*widget.FormItem{
		widget.NewFormItem("Type*", typeSelect),
		widget.NewFormItem("Amount*", amountEntery),
		widget.NewFormItem("Category*", categorySelectEntry),
		widget.NewFormItem("Description", descriptionEntery),
	}

	formDialog := dialog.NewForm("Add New Transaction", "Add", "Cancel", items, func(confirmed bool) {
		if !confirmed {
			return
		}

		amount, _ := strconv.ParseFloat(amountEntery.Text, 64)

		newTransaction := models.Transaction{
			ID:          storage.NextTransactionID,
			Date:        time.Now(),
			Amount:      amount,
			Category:    strings.TrimSpace(categorySelectEntry.Text),
			Description: strings.TrimSpace(descriptionEntery.Text),
			Type:        typeSelect.Selected,
		}
		core.AddTransaction(newTransaction)
		storage.SaveData()

		dialog.ShowInformation("Success", "Transaction added successfully!", t.guiApp.GuiWindow)
		t.guiApp.ShowTransactionsScreen()

	}, t.guiApp.GuiWindow)

	formDialog.Resize(fyne.NewSize(450, 300))
	formDialog.Show()
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

// func (t *TransactionsScreen) createTransactionList() []fyne.CanvasObject{

// }
