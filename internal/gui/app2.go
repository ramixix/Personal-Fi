package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

type GuiApp struct {
	application fyne.App
	GuiWindow   fyne.Window

	dashboardScreen *DashboardScreen
	// transactionsScreen *TransactionScreen
	// goalsScreen *GoalsScreen
	// accountsScreen *AccountsScreen
	// reportsScreen *ReportsScreen
}

func Run() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Financial Tracker")
	// set window size and center it on the screen
	myWindow.Resize(fyne.NewSize(1200, 800))
	myWindow.CenterOnScreen()

	gui := GuiApp{application: myApp, GuiWindow: myWindow}

	// Initialize screens
	gui.dashboardScreen = NewDashboardScreen(&gui)

	gui.ShowDashboardScreen()
	myWindow.ShowAndRun()
}

// Render dashboard screen
func (a *GuiApp) ShowDashboardScreen() {
	a.GuiWindow.SetContent(a.dashboardScreen.Render())
}

// Render Transactions screen
func (a *GuiApp) ShowTransactionsScreen() {
	// TODO: Will implement next
	a.ShowDashboardScreen()
}

// Render Goals screen
func (a *GuiApp) ShowGoalsScreen() {
	// TODO: Will implement next
	a.ShowDashboardScreen()
}
