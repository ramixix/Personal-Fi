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

	// Buttons
	addMoneyBtn := widget.NewButton("➕ Add Money", func() { a.addMoneyToAccountDialog(account) })
	addMoneyBtn.Importance = widget.HighImportance

	viewHistoryBtn := widget.NewButton("📊 History", func() { a.showAccountHistoryDialog(account) })

	editBtn := widget.NewButton("Edit", func() { a.editAccountDialog(account) })
	editBtn.Importance = widget.WarningImportance

	deleteBtn := widget.NewButton("Delete", func() { a.deleteAccountDialog(account) })
	deleteBtn.Importance = widget.DangerImportance

	actions := container.NewVBox(
		container.NewGridWithColumns(2, addMoneyBtn, viewHistoryBtn),
		container.NewGridWithColumns(2, editBtn, deleteBtn),
	)

	cardContent := container.NewVBox(
		nameLabel,
		balanceLabel,
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

// ----------------------------------------------------------------------
//
//	Dialog to show account history (history == account transactions)
//
// ----------------------------------------------------------------------
func (a *AccountsScreen) showAccountHistoryDialog(account models.Account) {
	accountTransactions := core.GetAccountTransactions(account.ID)

	if len(accountTransactions) == 0 {
		dialog.ShowInformation(fmt.Sprintf("%s Account History", account.Name), "No transactions for this account yet!", a.guiApp.GuiWindow)
		return
	}

	historyText := fmt.Sprintf("Account: %s\nCurrent Balance: $%.2f\n\n", account.Name, account.Balance)
	historyText += "--- Transaction History ---\n\n"

	for i := len(accountTransactions) - 1; i >= 0; i-- {
		accTransac := accountTransactions[i]
		historyText += fmt.Sprintf("Date: %s\nAmount: %.2f\nNote: %s\n\n", accTransac.Date.Format("02/01/2006 15:04:05"), accTransac.Amount, accTransac.Note)
	}

	historyLabel := widget.NewLabel(historyText)
	scrollContainer := container.NewScroll(historyLabel)
	scrollContainer.SetMinSize(fyne.NewSize(400, 300))

	statsLabel := widget.NewLabel(fmt.Sprintf("\n--- Statistics ---\nTotal Contributions: %d\nTotal Added: $%.2f\nAverage: $%.2f",
		len(accountTransactions),
		account.Balance,
		account.Balance/float64(len(accountTransactions))))

	bottomSection := container.NewVBox(
		widget.NewSeparator(),
		statsLabel,
	)

	// Unlike VBox, which treats everyone "equally squished," the Border layout gives the bottomSection exactly the height it needs and then forces the scrollContainer in the center to stretch and fill the rest of the window.
	content := container.NewBorder(nil, bottomSection, nil, nil, scrollContainer)
	dialog.ShowCustom(fmt.Sprintf("%s Account History", account.Name), "Close", content, a.guiApp.GuiWindow)
}

// ----------------------------------------
//
//	Dialog to edit account information
//
// ----------------------------------------
func (a *AccountsScreen) editAccountDialog(account models.Account) {
	nameEntry := widget.NewEntry()
	nameEntry.SetText(account.Name)
	nameEntry.Validator = func(value string) error {
		if len(strings.TrimSpace(value)) < 2 {
			return errors.New("Name too short. Must be at least 2 characters.")
		}
		return nil
	}

	formItems := []*widget.FormItem{widget.NewFormItem("Account Name", nameEntry)}

	editDialog := dialog.NewForm("Edit Account", "Save", "Cancel", formItems, func(confirmed bool) {
		if !confirmed {
			return
		}
		accountPointer := core.FindAccount(account.ID)
		if accountPointer == nil {
			dialog.ShowError(fmt.Errorf("account not found"), a.guiApp.GuiWindow)
			return
		}

		if nameEntry.Text == "" {
			dialog.ShowError(fmt.Errorf("please enter a valid account name"), a.guiApp.GuiWindow)
			return
		}
		accountPointer.Name = nameEntry.Text
		storage.SaveData()
		dialog.ShowInformation("Success", "Account updated successfully!", a.guiApp.GuiWindow)
		a.guiApp.ShowAccountsScreen()
	},
		a.guiApp.GuiWindow)

	editDialog.Resize(fyne.NewSize(400, 200))
	editDialog.Show()
}

// --------------------------------
//
//	Dialog to delete accounts
//
// --------------------------------
func (a *AccountsScreen) deleteAccountDialog(account models.Account) {
	deleteDialog := dialog.NewConfirm(fmt.Sprintf("%s Account Deletion", account.Name),
		fmt.Sprintf("Are you sure you want to remove %s from account list. All transaction belonging to this account will be also deleted", account.Name),
		func(confirmed bool) {
			if !confirmed {
				return
			}
			deleted := core.DeleteAccount(account.ID)
			if !deleted {
				dialog.ShowError(fmt.Errorf("Could not delete account and related transactions"), a.guiApp.GuiWindow)
				return
			}
			dialog.ShowInformation("Successful Deletion", "Account and related transaction are all deleted.", a.guiApp.GuiWindow)
			storage.SaveData()
			a.guiApp.ShowAccountsScreen()
		},
		a.guiApp.GuiWindow,
	)

	deleteDialog.Resize(fyne.NewSize(400, 200))
	deleteDialog.Show()
}
