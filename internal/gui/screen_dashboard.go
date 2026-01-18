package gui

import (
	"financial_tracker/internal/core"
	"fmt"
	"image/color"

	// "financial_tracker/internal/models"
	"financial_tracker/internal/storage"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type DashboardScreen struct {
	guiApp *GuiApp
}

// NewDashboardScreen creates a new dashboard screen
func NewDashboardScreen(app *GuiApp) *DashboardScreen {
	return &DashboardScreen{guiApp: app}
}

func (d *DashboardScreen) Render() fyne.CanvasObject {
	// page title
	title := widget.NewLabelWithStyle("Dashboard", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	// create summary cards
	summaryCards := d.createSummaryCards()

	content := container.NewVBox(
		title,
		widget.NewSeparator(),
		summaryCards,
	)
	return content
}

func (d *DashboardScreen) createSummaryCards() fyne.CanvasObject {
	// total number of transactions so far + number of accounts + balance across all acounts
	transactionsCount := fmt.Sprintf("%d total transaction", len(storage.Transactions))
	accountsCount := fmt.Sprintf("%d accounts", len(storage.Accounts))
	totalAccounts := fmt.Sprintf("$%.2f total balance across acounts", core.GetTotalAccountBalance())

	// total number of goals + only active ones
	goalsCount := fmt.Sprintf("%d goals", len(storage.Goals))
	activeGoals := fmt.Sprintf("%d active goasl", len(core.GetActiveGoals()))

	// total income + total expense + net value
	totalIncome, totalExpenses := core.CalculateTotals()
	netWorth := totalIncome - totalExpenses

	totalIncomeText := fmt.Sprintf("$%.2f", totalIncome)
	totalExpensesText := fmt.Sprintf("$%.2f", totalExpenses)
	// if netWorth >0 then + otherwise -
	netWorthText := ""
	if netWorth > 0 {
		netWorthText = fmt.Sprintf("+ $%.2f", netWorth)
	}

	// income card
	incomeCard := widget.NewCard("💰 Income", "", widget.NewLabel(totalIncomeText))

	// expenses card
	expensesCard := widget.NewCard("💸 Expenses", "", widget.NewLabel(totalExpensesText))

	// Net Worth Card
	netWorthCard := widget.NewCard("📊 Net Worth", "", widget.NewLabel(netWorthText))

	// Total Accounts balance card
	accountsBalanceCard := widget.NewCard("🏦 Accounts Balance", "", widget.NewLabel(totalAccounts))

	// total Number of accounts
	accoutNumCard := widget.NewCard("Accounts", "", widget.NewLabel(accountsCount))

	// total goals card
	GoalsNumCard := widget.NewCard("🎯 Goals", "", widget.NewLabel(goalsCount))

	// Active goals card
	activeGoalsCard := widget.NewCard("🎯 Active Goals", "", widget.NewLabel(activeGoals))

	// transaction count card
	transactionsNumCard := widget.NewCard("Transactions", "", widget.NewLabel(transactionsCount))

	cardsGrid_1 := container.NewGridWithColumns(4, transactionsNumCard, accoutNumCard, GoalsNumCard, activeGoalsCard)
	cardsGrid_2 := container.NewGridWithColumns(4, accountsBalanceCard, incomeCard, expensesCard, netWorthCard)

	content := container.NewVBox(
		widget.NewLabelWithStyle("Financial Review", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		cardsGrid_1,
		cardsGrid_2,
	)
	return content
}

func NewSummaryCard(title, value string) fyne.CanvasObject {
	bg := canvas.NewRectangle(color.NRGBA{R: 30, G: 30, B: 30, A: 255})
	bg.StrokeColor = theme.PrimaryColor()
	bg.StrokeWidth = 2

	titleLabel := widget.NewLabelWithStyle(
		title,
		fyne.TextAlignLeading,
		fyne.TextStyle{Bold: true},
	)

	valueLabel := widget.NewLabelWithStyle(
		value,
		fyne.TextAlignLeading,
		fyne.TextStyle{Bold: true},
	)

	content := container.NewVBox(
		titleLabel,
		valueLabel,
	)

	card := container.NewPadded(content)

	return container.NewStack(bg, card)
}

// incomeCard := NewSummaryCard("Income", totalIncomeText)
// expensesCard := NewSummaryCard("Expenses", totalExpensesText)
// netWorthCard := NewSummaryCard("Net Worth", netWorthText)

// accountsBalanceCard := NewSummaryCard("Total Balance", totalAccounts)
// accoutNumCard := NewSummaryCard("Accounts", accountsCount)
// GoalsNumCard := NewSummaryCard("Goals", goalsCount)
// activeGoalsCard := NewSummaryCard("Active Goals", activeGoals)
// transactionsNumCard := NewSummaryCard("Transactions", transactionsCount)
