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

	filteredTransactions *fyne.Container
	// values needed for filtering
	searchText     string
	filterType     string
	filterCategory string
	filterPeriod   string
}

func NewTransactionsScreen(app *GuiApp) *TransactionsScreen {
	return &TransactionsScreen{guiApp: app, filterType: "", filterCategory: "", filterPeriod: "", filteredTransactions: container.NewVBox()}
}

func (t *TransactionsScreen) Render() fyne.CanvasObject {
	// Header with title and stats and simple add button
	header := t.createHeader()

	// filter section
	filterBar := t.createFilterBar()

	// section to list transaction based on filters
	// t.filteredTransactions = container.NewVBox()
	// t.refereshTransactionList()
	// transactionList := t.createTransactionList()

	content := container.NewVBox(
		header,
		widget.NewSeparator(),
		filterBar,
		widget.NewSeparator(),
		t.filteredTransactions,
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
	searchEntry.OnChanged = func(text string) {
		t.searchText = text
		t.refereshTransactionList()
	}

	// Filtering Transactions by Type of transaction
	typeSelect := widget.NewSelect([]string{"All Types", "Income", "Expenses"},
		func(value string) {
			switch value {
			case "Income":
				t.filterType = "income"
			case "Expenses":
				t.filterType = "expense"
			default:
				t.filterType = ""
			}
			t.refereshTransactionList()
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
				t.filterPeriod = ""
			}
			t.refereshTransactionList()
		})
	periodSelect.SetSelected("All Time")

	// filtering by categories
	categories := []string{"All Categories"}
	categories = append(categories, core.GetCategories()...)
	categorySelect := widget.NewSelect(categories,
		func(value string) {
			if value == "All Categories" {
				t.filterCategory = ""
			} else {
				t.filterCategory = value
			}
			t.refereshTransactionList()
		})
	categorySelect.SetSelected("All Categories")

	// button to clear all filters and set default selects for selections
	clearBtn := widget.NewButton("Clear Filters", func() {
		t.searchText = ""
		t.filterType = ""
		t.filterCategory = ""
		t.filterPeriod = ""
		searchEntry.SetText("")
		typeSelect.SetSelected("All Types")
		periodSelect.SetSelected("All Time")
		categorySelect.SetSelected("All Categories")
		t.refereshTransactionList()
	})

	filterRow1 := container.NewGridWithColumns(2, searchEntry, clearBtn)
	filterRow2 := container.NewGridWithColumns(3, typeSelect, periodSelect, categorySelect)

	return container.NewVBox(filterRow1, filterRow2)
}

func (t *TransactionsScreen) refereshTransactionList() {
	var transactions []models.Transaction = t.getFilteredTransactions()
	t.filteredTransactions.Objects = nil

	if len(transactions) == 0 {
		t.filteredTransactions.Add(t.createEmptyState())
	}

	groupedTransacs := t.groupTransactionsByMonth(transactions)

	var cards []fyne.CanvasObject
	for monthYear, transactionList := range groupedTransacs {
		monthYearLable := widget.NewLabelWithStyle(monthYear, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
		cards = append(cards, monthYearLable)

		for _, transac := range transactionList {
			transactText := fmt.Sprintf("%d, %s, %.2f, %s, %s", transac.ID, transac.Category, transac.Amount, transac.Type, transac.Date)
			cards = append(cards, widget.NewLabel(transactText))
		}
	}

	t.filteredTransactions.Add(container.NewVBox(cards...))

	t.filteredTransactions.Refresh()
}

func (t *TransactionsScreen) getFilteredTransactions() []models.Transaction {
	startDate := time.Now()
	endDate := time.Now()
	hasRange := true
	if t.filterPeriod == "week" {
		startDate = startDate.AddDate(0, 0, -7)
	} else if t.filterPeriod == "month" {
		startDate = startDate.AddDate(0, -1, 0)
	} else if t.filterPeriod == "year" {
		startDate = startDate.AddDate(-1, 0, 0)
	} else {
		hasRange = false
	}

	category := []string{t.filterCategory}
	if strings.TrimSpace(t.filterCategory) == "" {
		category = []string{}
	}
	itemsToSearch := models.SearchCriteria{
		Keyword:         strings.TrimSpace(t.searchText),
		Categories:      category,
		TransactionType: t.filterType,
		HasDateRange:    hasRange,
		StartDate:       startDate,
		EndDate:         endDate,
	}

	var transactions []models.Transaction = core.AdvancedSearchTransactions(itemsToSearch)
	fmt.Printf("%d : %s %s %s\n", len(transactions), t.searchText, startDate, endDate)
	return transactions
}

func (t *TransactionsScreen) createEmptyState() fyne.CanvasObject {
	emptyIcon := widget.NewLabelWithStyle("📭", fyne.TextAlignCenter, fyne.TextStyle{})
	emptyText := widget.NewLabelWithStyle("No Transaction Found!", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	emptyHint := widget.NewLabel("Try adjusting your filters or add a new transaction")
	addBtn := widget.NewButton("➕ Add Your First Transaction", func() { t.showAddTransactionDialog() })
	addBtn.Importance = widget.HighImportance

	return container.NewCenter(container.NewVBox(emptyIcon, emptyText, emptyHint, addBtn))
}

func (t *TransactionsScreen) groupTransactionsByMonth(transactions []models.Transaction) map[string][]models.Transaction {
	groups := make(map[string][]models.Transaction)

	for _, transac := range transactions {
		monthYear := transac.Date.Format("January 2006")
		groups[monthYear] = append(groups[monthYear], transac)
	}
	return groups
}
