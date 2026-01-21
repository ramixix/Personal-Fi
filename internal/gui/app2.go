package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
)

type GuiApp struct {
	application fyne.App
	GuiWindow   fyne.Window

	mainContent *fyne.Container

	dashboardScreen    *DashboardScreen
	transactionsScreen *TransactionsScreen
	goalsScreen        *GoalsScreen
	accountsScreen     *AccountsScreen
	reportsScreen      *ReportsScreen
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

	gui.mainContent = container.NewStack()

	sideBard := gui.createSidebar()

	split := container.NewHSplit(sideBard, gui.mainContent)
	split.SetOffset(0.15)

	myWindow.SetContent(split)
	gui.ShowDashboardScreen()
	myWindow.ShowAndRun()
}

// Render dashboard screen
func (a *GuiApp) ShowDashboardScreen() {
	a.mainContent.Objects = []fyne.CanvasObject{a.dashboardScreen.Render()}
	a.mainContent.Refresh()
	// a.GuiWindow.SetContent(a.dashboardScreen.Render())
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

func (a *GuiApp) ShowAccountsScreen() {
	// TODO: Will implement next
	a.ShowDashboardScreen()
}

func (a *GuiApp) ShowReportsScreen() {
	// TODO: Will implement next
	a.ShowDashboardScreen()
}

func (a *GuiApp) ShowSettingsScreen() {
	// TODO: Will implement next
	a.ShowDashboardScreen()
}
