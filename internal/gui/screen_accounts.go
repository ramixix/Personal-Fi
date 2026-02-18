package gui

import (
	"financial_tracker/internal/core"
	"financial_tracker/internal/storage"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type AccountsScreen struct {
	guiApp *GuiApp
}

func NewAccountsScreen(app *GuiApp) *AccountsScreen {
	return &AccountsScreen{guiApp: app}
}

func (a *AccountsScreen) Render() fyne.CanvasObject {
	header := a.createHeader()

	content := container.NewVBox(
		header,
		widget.NewSeparator(),
		widget.NewLabel("Accounts management page - Coming soon!"),
	)

	return container.NewScroll(content)
}

func (a *AccountsScreen) createHeader() fyne.CanvasObject {
	title := widget.NewLabelWithStyle("🏦 Accounts", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	accountCount := len(storage.Accounts)
	accountTotals := core.GetTotalAccountBalance()

	statsLabel := widget.NewLabel(fmt.Sprintf("Total Accounts: %d | Total Balance: %.2f", accountCount, accountTotals))

	addNewAccountBtn := widget.NewButton("Add new account", func() {})
	addNewAccountBtn.Importance = widget.HighImportance

	header := container.NewBorder(nil, nil, title, statsLabel)
	return container.NewVBox(header, addNewAccountBtn)
}
