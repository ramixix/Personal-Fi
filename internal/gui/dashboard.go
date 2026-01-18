package gui

// import (
// 	"financial_tracker/internal/core"
// 	"financial_tracker/internal/models"
// 	"financial_tracker/internal/storage"
// 	"fmt"
// 	"time"

// 	"fyne.io/fyne/v2"
// 	"fyne.io/fyne/v2/container"
// 	"fyne.io/fyne/v2/widget"
// )

// // // DashboardScreen represents the dashboard page
// // type DashboardScreen struct {
// // 	app *App
// // }

// // // NewDashboardScreen creates a new dashboard screen
// // func NewDashboardScreen(app *App) *DashboardScreen {
// // 	return &DashboardScreen{
// // 		app: app,
// // 	}
// // }

// // Render creates and returns the dashboard content
// func (d *DashboardScreen) Render() fyne.CanvasObject {
// 	// Page title
// 	title := widget.NewLabelWithStyle(
// 		"Dashboard",
// 		fyne.TextAlignCenter,
// 		fyne.TextStyle{Bold: true},
// 	)

// 	// Create summary cards
// 	summaryCards := d.createSummaryCards()

// 	// Create recent transactions section
// 	recentTransactions := d.createRecentTransactions()

// 	// Create quick actions
// 	quickActions := d.createQuickActions()

// 	// Layout everything vertically
// 	content := container.NewVBox(
// 		title,
// 		widget.NewSeparator(),
// 		summaryCards,
// 		widget.NewSeparator(),
// 		quickActions,
// 		widget.NewSeparator(),
// 		recentTransactions,
// 	)

// 	// Wrap in a scrollable container
// 	return container.NewScroll(content)
// }

// // createSummaryCards creates the financial summary cards
// func (d *DashboardScreen) createSummaryCards() fyne.CanvasObject {
// 	// Calculate totals
// 	totalIncome, totalExpenses := core.CalculateTotals()
// 	netWorth := totalIncome - totalExpenses

// 	// Get account totals
// 	totalAccounts := core.GetTotalAccountBalance()

// 	// Get active goals count
// 	activeGoals := core.GetActiveGoals()

// 	// Income Card
// 	incomeCard := widget.NewCard(
// 		"💰 Total Income",
// 		"",
// 		widget.NewLabel(fmt.Sprintf("$%.2f", totalIncome)),
// 	)

// 	// Expenses Card
// 	expensesCard := widget.NewCard(
// 		"💸 Total Expenses",
// 		"",
// 		widget.NewLabel(fmt.Sprintf("$%.2f", totalExpenses)),
// 	)

// 	// Net Worth Card
// 	netWorthLabel := fmt.Sprintf("$%.2f", netWorth)
// 	if netWorth >= 0 {
// 		netWorthLabel = "+" + netWorthLabel
// 	}
// 	netWorthCard := widget.NewCard(
// 		"📊 Net Worth",
// 		"",
// 		widget.NewLabel(netWorthLabel),
// 	)

// 	// Accounts Card
// 	accountsCard := widget.NewCard(
// 		"🏦 Total Accounts Balance",
// 		"",
// 		widget.NewLabel(fmt.Sprintf("$%.2f", totalAccounts)),
// 	)

// 	// Goals Card
// 	goalsCard := widget.NewCard(
// 		"🎯 Active Goals",
// 		"",
// 		widget.NewLabel(fmt.Sprintf("%d goals", len(activeGoals))),
// 	)

// 	// Arrange cards in a grid (3 columns)
// 	cardsGrid := container.NewGridWithColumns(3,
// 		incomeCard,
// 		expensesCard,
// 		netWorthCard,
// 	)

// 	cardsGrid2 := container.NewGridWithColumns(2,
// 		accountsCard,
// 		goalsCard,
// 	)

// 	return container.NewVBox(
// 		widget.NewLabelWithStyle("Financial Overview", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
// 		cardsGrid,
// 		cardsGrid2,
// 	)
// }

// // createRecentTransactions shows the last 5 transactions
// func (d *DashboardScreen) createRecentTransactions() fyne.CanvasObject {
// 	sectionTitle := widget.NewLabelWithStyle(
// 		"Recent Transactions",
// 		fyne.TextAlignLeading,
// 		fyne.TextStyle{Bold: true},
// 	)

// 	// Get recent transactions
// 	recentTx := core.GetRecentTransactions(5)

// 	if len(recentTx) == 0 {
// 		noData := widget.NewLabel("No transactions yet. Add your first transaction!")
// 		return container.NewVBox(sectionTitle, noData)
// 	}

// 	// Create a list of transaction cards
// 	var txCards []fyne.CanvasObject

// 	for _, tx := range recentTx {
// 		// Format the transaction
// 		typeIcon := "💰"
// 		if tx.Type == "expense" {
// 			typeIcon = "💸"
// 		}

// 		title := fmt.Sprintf("%s %s", typeIcon, tx.Category)
// 		subtitle := tx.Description
// 		amount := fmt.Sprintf("$%.2f", tx.Amount)
// 		date := tx.Date.Format("Jan 02, 2006")

// 		// Create card content
// 		cardContent := container.NewVBox(
// 			widget.NewLabelWithStyle(title, fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
// 			widget.NewLabel(subtitle),
// 			container.NewHBox(
// 				widget.NewLabel(date),
// 				widget.NewLabel("  |  "),
// 				widget.NewLabelWithStyle(amount, fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
// 			),
// 		)

// 		card := widget.NewCard("", "", cardContent)
// 		txCards = append(txCards, card)
// 	}

// 	// Add "View All" button
// 	viewAllBtn := widget.NewButton("View All Transactions", func() {
// 		d.app.ShowTransactions()
// 	})

// 	txCards = append(txCards, viewAllBtn)

// 	return container.NewVBox(
// 		sectionTitle,
// 		container.NewVBox(txCards...),
// 	)
// }

// // createQuickActions creates quick action buttons
// func (d *DashboardScreen) createQuickActions() fyne.CanvasObject {
// 	sectionTitle := widget.NewLabelWithStyle(
// 		"Quick Actions",
// 		fyne.TextAlignLeading,
// 		fyne.TextStyle{Bold: true},
// 	)

// 	// Add Transaction button
// 	addTransactionBtn := widget.NewButton("➕ Add Transaction", func() {
// 		d.showAddTransactionDialog()
// 	})

// 	// Add to Goal button
// 	addToGoalBtn := widget.NewButton("🎯 Contribute to Goal", func() {
// 		d.app.ShowGoals()
// 	})

// 	// View Summary button
// 	viewSummaryBtn := widget.NewButton("📊 View Summary", func() {
// 		d.showSummaryDialog()
// 	})

// 	// Arrange buttons
// 	buttonsGrid := container.NewGridWithColumns(3,
// 		addTransactionBtn,
// 		addToGoalBtn,
// 		viewSummaryBtn,
// 	)

// 	return container.NewVBox(
// 		sectionTitle,
// 		buttonsGrid,
// 	)
// }

// // showAddTransactionDialog shows a dialog to add a transaction
// func (d *DashboardScreen) showAddTransactionDialog() {
// 	// Create form fields
// 	typeSelect := widget.NewSelect([]string{"income", "expense"}, func(value string) {})
// 	typeSelect.SetSelected("expense")

// 	amountEntry := widget.NewEntry()
// 	amountEntry.SetPlaceHolder("0.00")

// 	categoryEntry := widget.NewEntry()
// 	categoryEntry.SetPlaceHolder("e.g., Food, Salary")

// 	descriptionEntry := widget.NewEntry()
// 	descriptionEntry.SetPlaceHolder("Description")

// 	// Create form
// 	form := &widget.Form{
// 		Items: []*widget.FormItem{
// 			{Text: "Type", Widget: typeSelect},
// 			{Text: "Amount", Widget: amountEntry},
// 			{Text: "Category", Widget: categoryEntry},
// 			{Text: "Description", Widget: descriptionEntry},
// 		},
// 		OnSubmit: func() {
// 			// Parse and save transaction
// 			amount := 0.0
// 			fmt.Sscanf(amountEntry.Text, "%f", &amount)

// 			if amount <= 0 {
// 				dialog := widget.NewLabel("Please enter a valid amount!")
// 				closebutton := widget.NewButton("Try Again", func() { d.app.dashboardScreen.showAddTransactionDialog() })
// 				d.app.Window.SetContent(container.NewVBox(dialog, closebutton))
// 				return
// 			}

// 			// Create transaction
// 			newTx := models.Transaction{
// 				ID:          storage.NextTransactionID,
// 				Date:        time.Now(),
// 				Amount:      amount,
// 				Category:    categoryEntry.Text,
// 				Description: descriptionEntry.Text,
// 				Type:        typeSelect.Selected,
// 			}

// 			core.AddTransaction(newTx)
// 			storage.SaveData()

// 			// Refresh dashboard
// 			d.app.ShowDashboard()
// 		},
// 		OnCancel: func() {
// 			d.app.ShowDashboard()
// 		},
// 	}

// 	// Show dialog
// 	dialogContent := container.NewVBox(
// 		widget.NewLabelWithStyle("Add Transaction", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
// 		form,
// 	)

// 	d.app.Window.SetContent(dialogContent)
// }

// // showSummaryDialog shows financial summary
// func (d *DashboardScreen) showSummaryDialog() {
// 	totalIncome, totalExpenses := core.CalculateTotals()
// 	netWorth := totalIncome - totalExpenses
// 	avgIncome, avgExpenses := core.GetMonthlyAverage()

// 	summary := fmt.Sprintf(`
// Financial Summary
// =================

// Total Income:    $%.2f
// Total Expenses:  $%.2f
// Net Worth:       $%.2f

// Monthly Average:
//   Income:        $%.2f
//   Expenses:      $%.2f
//   Net:           $%.2f

// Total Transactions: %d
// `,
// 		totalIncome,
// 		totalExpenses,
// 		netWorth,
// 		avgIncome,
// 		avgExpenses,
// 		avgIncome-avgExpenses,
// 		len(storage.Transactions),
// 	)

// 	summaryLabel := widget.NewLabel(summary)

// 	closeBtn := widget.NewButton("Close", func() {
// 		d.app.ShowDashboard()
// 	})

// 	content := container.NewVBox(
// 		summaryLabel,
// 		closeBtn,
// 	)

// 	d.app.Window.SetContent(content)
// }
