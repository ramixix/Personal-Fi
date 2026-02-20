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

type AccountsScreen struct {
	guiApp *GuiApp
}

func NewAccountsScreen(app *GuiApp) *AccountsScreen {
	return &AccountsScreen{guiApp: app}
}

func (a *AccountsScreen) Render() fyne.CanvasObject {
	header := a.createHeader()

	accountGrid := a.createAccountsGrid()

	content := container.NewVBox(
		header,
		widget.NewSeparator(),
		accountGrid,
	)

	return container.NewScroll(content)
}

// -------------------------------------------------------------------------------------------------------------------
//
//	Account header which shows account number, total balance across accounts and a button to add new accouts
//
// -------------------------------------------------------------------------------------------------------------------
func (a *AccountsScreen) createHeader() fyne.CanvasObject {
	title := widget.NewLabelWithStyle("🏦 Accounts", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	accountCount := len(storage.Accounts)
	accountTotals := core.GetTotalAccountBalance()

	statsLabel := widget.NewLabel(fmt.Sprintf("Total Accounts: %d | Total Balance: %.2f", accountCount, accountTotals))

	addNewAccountBtn := widget.NewButton("Add new account", func() { a.showCreateAccountDialog() })
	addNewAccountBtn.Importance = widget.HighImportance

	header := container.NewBorder(nil, nil, title, statsLabel)
	return container.NewVBox(header, addNewAccountBtn)
}

// --------------------------------------------------
//
//	create grid of account cards with 2 columns
//
// --------------------------------------------------
func (a *AccountsScreen) createAccountsGrid() fyne.CanvasObject {
	if len(storage.Accounts) == 0 {
		return a.createEmptyState()
	}

	var cards []fyne.CanvasObject
	for _, acc := range storage.Accounts {
		card := a.createAccountCard(acc)
		cards = append(cards, card)
	}

	return container.NewGridWithColumns(2, cards...)
}

// ---------------------------
//
//	Create Account Card
//
// ---------------------------
func (a *AccountsScreen) createAccountCard(account models.Account) fyne.CanvasObject {
	nameLabel := widget.NewLabelWithStyle(account.Name, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	balanceLabel := widget.NewLabelWithStyle(fmt.Sprintf("%.2f", account.Balance), fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	createdTimeLabel := widget.NewLabel(fmt.Sprintf("Created: %s", account.Created.Format("Jan 02, 2006")))

	AccountTransactionCount := len(core.GetAccountTransactions(account.ID))
	AccountTransactionCountLabel := widget.NewLabel(fmt.Sprintf("%d Number of Account Transactions", AccountTransactionCount))

	// balanceBar := a.createBalanceIndicator(account.Balance)

	// Buttons
	addMoneyBtn := widget.NewButton("➕ Add Money", func() { a.addMoneyToAccountDialog(account) })
	addMoneyBtn.Importance = widget.HighImportance

	viewHistoryBtn := widget.NewButton("📊 History", func() {})

	editBtn := widget.NewButton("Edit", func() {})
	editBtn.Importance = widget.WarningImportance

	deleteBtn := widget.NewButton("Delete", func() {})
	deleteBtn.Importance = widget.DangerImportance

	actions := container.NewVBox(
		container.NewGridWithColumns(2, addMoneyBtn, viewHistoryBtn),
		container.NewGridWithColumns(2, editBtn, deleteBtn),
	)

	cardContent := container.NewVBox(
		nameLabel,
		balanceLabel,
		// balanceBar,
		createdTimeLabel,
		AccountTransactionCountLabel,
		widget.NewSeparator(),
		actions,
	)

	card := widget.NewCard("", "", cardContent)
	return card

}

// ---------------------------------------------------------------
//
//	message to show if there is not account availble to show
//
// ---------------------------------------------------------------
func (a *AccountsScreen) createEmptyState() fyne.CanvasObject {
	emptyIcon := widget.NewLabelWithStyle("🏦", fyne.TextAlignCenter, fyne.TextStyle{})
	message := widget.NewLabelWithStyle("No accounts yet!", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	createBtn := widget.NewButton("➕ Create Your First Account", func() { a.showCreateAccountDialog() })
	createBtn.Importance = widget.HighImportance

	return container.NewVBox(emptyIcon, message, createBtn)
}

// ----------------------------------
//
//	Dialog to create new accounts
//
// ----------------------------------
func (a *AccountsScreen) showCreateAccountDialog() {
	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("Account name (e.g., Emergency Fund, Vacation)")
	nameEntry.Validator = func(value string) error {
		if len(strings.TrimSpace(value)) < 2 {
			return errors.New("Name too short. Must be at least 2 characters.")
		}
		return nil
	}

	formItems := []*widget.FormItem{
		widget.NewFormItem("Account Name", nameEntry),
	}

	createAccountDialog := dialog.NewForm("Create Account", "Create", "Cancel", formItems, func(confirmed bool) {
		if !confirmed {
			return
		}

		if len(strings.TrimSpace(nameEntry.Text)) < 2 {
			dialog.ShowError(fmt.Errorf("Not a valid account name. Must at least contains 2 characters."), a.guiApp.GuiWindow)
			return
		}

		newAccount := models.Account{ID: storage.NextAccountID, Name: nameEntry.Text, Balance: 0, Created: time.Now()}
		core.AddAccount(newAccount)
		storage.SaveData()

		dialog.ShowInformation("Success", fmt.Sprintf("Account '%s' created! 🎉", newAccount.Name), a.guiApp.GuiWindow)
		a.guiApp.ShowAccountsScreen()
	},
		a.guiApp.GuiWindow)

	createAccountDialog.Resize(fyne.NewSize(400, 200))
	createAccountDialog.Show()
}

// -----------------------------------
//
//	Dialog to add money to account
//
// -----------------------------------
func (a *AccountsScreen) addMoneyToAccountDialog(account models.Account) {
	amountEntry := widget.NewEntry()
	amountEntry.SetPlaceHolder("0.0")
	amountEntry.Validator = func(value string) error {
		amount, err := strconv.ParseFloat(strings.TrimSpace(value), 64)
		if err != nil {
			return errors.New("Please enter a valid number.")
		}
		if amount <= 0 {
			return errors.New("Amount to add must be positive and greater than zero")
		}
		return nil
	}

	noteEntry := widget.NewEntry()
	noteEntry.SetPlaceHolder("Note (optional)")

	addMoneyFormitems := []*widget.FormItem{
		widget.NewFormItem("Amount", amountEntry),
		widget.NewFormItem("Note", noteEntry),
	}

	addMoneyDialog := dialog.NewForm(fmt.Sprintf("Add Money to: %s", account.Name), "Add", "Cancel", addMoneyFormitems, func(confirmed bool) {
		if !confirmed {
			return
		}

		amount, err := strconv.ParseFloat(strings.TrimSpace(amountEntry.Text), 64)
		if err != nil || amount <= 0 {
			dialog.ShowError(fmt.Errorf("please enter a valid amount"), a.guiApp.GuiWindow)
			return
		}

		accountPointer := core.FindAccount(account.ID)
		core.AddMoneyToAccount(accountPointer, amount, noteEntry.Text)
		storage.SaveData()

		dialog.ShowInformation("Success", fmt.Sprintf("Added $%.2f to %s! 💰", amount, account.Name), a.guiApp.GuiWindow)
		a.guiApp.ShowAccountsScreen()
	}, a.guiApp.GuiWindow)

	addMoneyDialog.Resize(fyne.NewSize(400, 250))
	addMoneyDialog.Show()
}
