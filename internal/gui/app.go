package gui

// import (
// 	"fyne.io/fyne/v2"
// 	"fyne.io/fyne/v2/app"
// )

// type App struct {
// 	fyneApp fyne.App
// 	Window  fyne.Window

// 	// screens
// 	dashboardScreen *DashboardScreen
// 	// transactionScreen *TransactionScreen
// 	// goalsScreen       *GoalsScreen
// 	// reportsScreen     *ReportsScreen
// 	// accountsScreen    *AccountsScreen
// }

// // Run starts the GUI application
// func Run() {
// 	// Create a new Fyne application
// 	myApp := app.New()

// 	// Create main window, set size and place it on center of screen
// 	myWindow := myApp.NewWindow("Financial Tracker")
// 	myWindow.Resize(fyne.NewSize(1200, 800))
// 	myWindow.CenterOnScreen()

// 	guiApp := App{
// 		fyneApp: myApp,
// 		Window:  myWindow,
// 	}

// 	// Initialize dashboard screen
// 	guiApp.dashboardScreen = NewDashboardScreen(&guiApp)

// 	guiApp.ShowDashboard()

// 	// Show window and run the app
// 	myWindow.ShowAndRun()
// }

// // ShowDashboard shows the dashboard screen
// func (a *App) ShowDashboard() {
// 	a.Window.SetContent(a.dashboardScreen.Render())
// }

// // ShowTransactions is a placeholder for now
// func (a *App) ShowTransactions() {
// 	// TODO: Will implement next
// 	a.ShowDashboard()
// }

// // ShowGoals is a placeholder for now
// func (a *App) ShowGoals() {
// 	// TODO: Will implement next
// 	a.ShowDashboard()
// }
