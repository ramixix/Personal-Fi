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
}

func NewTransactionsScreen(app *GuiApp) *TransactionsScreen {
	return &TransactionsScreen{guiApp: app}
}

func (t *TransactionsScreen) Render() fyne.CanvasObject {
	header := t.createHeader()

	content := container.NewVBox(
		header,
		widget.NewSeparator(),
		widget.NewLabel("Transaction management page - Coming soon!"),
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
