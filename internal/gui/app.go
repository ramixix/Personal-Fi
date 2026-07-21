package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
)

const recent100 = 100

type GuiApp struct {
	application fyne.App
	GuiWindow   fyne.Window

	mainContent *fyne.Container

	dashboardScreen    *DashboardScreen
	transactionsScreen *TransactionsScreen
	goalsScreen        *GoalsScreen
	accountsScreen     *AccountsScreen
	reportsScreen      *ReportsScreen
	settingsScreen     *SettingsScreen
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
	gui.transactionsScreen = NewTransactionsScreen(&gui)
	gui.accountsScreen = NewAccountsScreen(&gui)
	gui.goalsScreen = NewGoalsScreen(&gui)
	gui.reportsScreen = NewReportsScreen(&gui)
	gui.settingsScreen = NewSettingsScreen(&gui)

	gui.mainContent = container.NewStack()

	sideBard := gui.createSidebar()

	split := container.NewHSplit(sideBard, gui.mainContent)
	split.SetOffset(0)

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
	a.mainContent.Objects = []fyne.CanvasObject{a.transactionsScreen.Render()}
	a.mainContent.Refresh()
}

// Render Accoutns screen
func (a *GuiApp) ShowAccountsScreen() {
	a.mainContent.Objects = []fyne.CanvasObject{a.accountsScreen.Render()}
	a.mainContent.Refresh()
}

// Render Goals screen
func (a *GuiApp) ShowGoalsScreen() {
	a.mainContent.Objects = []fyne.CanvasObject{a.goalsScreen.Render()}
	a.mainContent.Refresh()
}

func (a *GuiApp) ShowReportsScreen() {
	a.mainContent.Objects = []fyne.CanvasObject{a.reportsScreen.Render()}
	a.mainContent.Refresh()
}

func (a *GuiApp) ShowSettingsScreen() {
	a.mainContent.Objects = []fyne.CanvasObject{a.settingsScreen.Render()}
	a.mainContent.Refresh()
}
