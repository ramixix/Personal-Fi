package gui

import (
	"financial_tracker/internal/core"
	"fmt"

	// "financial_tracker/internal/models"
	"financial_tracker/internal/storage"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
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

	// create quick actions
	// quickActions := d.createQuickActions()

	// Create recent transaction and goal contribution section
	recentTransacGoalContributions := d.createRecentEvents()

	content := container.NewVBox(
		title,
		widget.NewSeparator(),
		summaryCards,
		widget.NewSeparator(),
		recentTransacGoalContributions,
	)
	return container.NewScroll(content)
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
		netWorthText = "+ "
	}
	netWorthText += fmt.Sprintf("%.2f", netWorth)

	// Average Montly Income/Expenses and Net of All time
	// avgIncome, avgExpenses := core.GetMonthlyAverage()

	// income card
	incomeCard := widget.NewCard("💰 Total Income", "", widget.NewLabel(totalIncomeText))

	// expenses card
	expensesCard := widget.NewCard("💸 Total Expenses", "", widget.NewLabel(totalExpensesText))

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

func (d *DashboardScreen) createRecentEvents() fyne.CanvasObject {
	title := widget.NewLabelWithStyle("Transaction & Goals", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	recentTransactions := core.GetRecentTransactions(5)
	recentGoalContributions := core.GetRecentGoalContributions(5)

	var transacCards []fyne.CanvasObject
	var goalContributionCards []fyne.CanvasObject

	if len(recentTransactions) == 0 {
		transacCards = append(transacCards, widget.NewLabel("No Transaction Found. Add Transaction First!"))

	} else {
		for i := len(recentTransactions) - 1; i >= 0; i-- {
			transac := recentTransactions[i]
			typeIcon := "💰 ✅"
			if transac.Type == "expense" {
				typeIcon = "💸 🔴"
			}

			cardTitle := fmt.Sprintf("%s Transaction ID: %d ", typeIcon, transac.ID)
			cardSubtitle := fmt.Sprintf("Category: %s\nType: %s\nNote: %s", transac.Category, transac.Type, transac.Description)
			cardAmount := fmt.Sprintf("%.2f", transac.Amount)
			date := transac.Date.Format("Jan 02, 2006")

			cardContent := container.NewVBox(
				widget.NewLabelWithStyle(cardTitle, fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
				widget.NewLabel(cardSubtitle),
				container.NewHBox(
					widget.NewLabel(date),
					widget.NewLabel("   |   "),
					widget.NewLabelWithStyle(cardAmount, fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
				),
			)
			card := widget.NewCard("", "", cardContent)
			transacCards = append(transacCards, card)
		}
	}

	if len(recentGoalContributions) == 0 {
		goalContributionCards = append(goalContributionCards, widget.NewLabel("No Goal Contribution Yet. Add Contirbution First!"))
	} else {

		for i := len(recentGoalContributions) - 1; i >= 0; i-- {
			contribution := recentGoalContributions[i]
			goal := core.FindGoal(contribution.GoalID)

			cardTitle := fmt.Sprintf("Contribution ID %d\nContributed To Goal: %s", contribution.ID, goal.Name)
			cardSubtitle := fmt.Sprintf("Note: %s", contribution.Note)
			cardAmount := fmt.Sprintf("Amount: %.2f", contribution.Amount)
			cardDate := contribution.Date.Format("Jan 02, 2006")

			cardContent := container.NewVBox(
				widget.NewLabelWithStyle(cardTitle, fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
				widget.NewLabel(cardSubtitle),
				container.NewHBox(
					widget.NewLabel(cardDate),
					widget.NewLabel("   |   "),
					widget.NewLabelWithStyle(cardAmount, fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
				),
			)

			card := widget.NewCard("", "", cardContent)
			goalContributionCards = append(goalContributionCards, card)
		}
	}

	// put all cards on top of each other
	transactions := container.NewVBox(transacCards...)
	contributions := container.NewVBox(goalContributionCards...)

	// create two columns one for transction one for goal contributions
	cardsGrid := container.NewGridWithColumns(2,
		transactions,
		contributions,
	)

	// craete button to view all transction in transction screen and goals in goals screen
	viewAllTransaction := widget.NewButton("View All Transactions", func() { d.guiApp.ShowTransactionsScreen() })
	viewAllGoals := widget.NewButton("View All Goals", func() { d.guiApp.ShowGoalsScreen() })

	return container.NewVBox(title, cardsGrid, viewAllTransaction, viewAllGoals)

}

// func NewSummaryCard(title, value string) fyne.CanvasObject {
// 	bg := canvas.NewRectangle(color.NRGBA{R: 30, G: 30, B: 30, A: 255})
// 	bg.StrokeColor = theme.PrimaryColor()
// 	bg.StrokeWidth = 2

// 	titleLabel := widget.NewLabelWithStyle(
// 		title,
// 		fyne.TextAlignLeading,
// 		fyne.TextStyle{Bold: true},
// 	)

// 	valueLabel := widget.NewLabelWithStyle(
// 		value,
// 		fyne.TextAlignLeading,
// 		fyne.TextStyle{Bold: true},
// 	)

// 	content := container.NewVBox(
// 		titleLabel,
// 		valueLabel,
// 	)

// 	card := container.NewPadded(content)

// 	return container.NewStack(bg, card)
// }

// incomeCard := NewSummaryCard("Income", totalIncomeText)
// expensesCard := NewSummaryCard("Expenses", totalExpensesText)
// netWorthCard := NewSummaryCard("Net Worth", netWorthText)

// accountsBalanceCard := NewSummaryCard("Total Balance", totalAccounts)
// accoutNumCard := NewSummaryCard("Accounts", accountsCount)
// GoalsNumCard := NewSummaryCard("Goals", goalsCount)
// activeGoalsCard := NewSummaryCard("Active Goals", activeGoals)
// transactionsNumCard := NewSummaryCard("Transactions", transactionsCount)
