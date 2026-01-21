package gui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func (a *GuiApp) createSidebar() fyne.CanvasObject {
	sideBarTitle := widget.NewLabelWithStyle("Finance Tracker", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	// Create navigation buttons
	dashboardBtn := createNavButton("🏠 Dashboard", func() { a.ShowDashboardScreen() })
	transactionsBtn := createNavButton("💰 Transactions", func() { a.ShowTransactionsScreen() })
	accountsBtn := createNavButton("🏦 Accounts", func() { a.ShowAccountsScreen() })
	goalsBtn := createNavButton("🎯 Goals", func() { a.ShowGoalsScreen() })
	reportsBtn := createNavButton("📊 Reports", func() { a.ShowReportsScreen() })

	// Settings button (at the bottom)
	settingsBtn := createNavButton("⚙️ Settings", func() { a.ShowSettingsScreen() })

	// Spacer to push settings to bottom
	spacer := layout.NewSpacer()

	navItems := container.NewVBox(
		sideBarTitle,
		widget.NewSeparator(),
		dashboardBtn,
		transactionsBtn,
		accountsBtn,
		goalsBtn,
		reportsBtn,
		spacer,
		widget.NewSeparator(),
		settingsBtn,
	)

	sidebar := container.NewPadded(navItems)
	bg := canvas.NewRectangle(color.RGBA{R: 50, G: 50, B: 60, A: 255})
	return container.NewStack(bg, sidebar)
}

func createNavButton(label string, onClickFunc func()) *widget.Button {
	btn := widget.NewButton(label, onClickFunc)
	btn.Alignment = widget.ButtonAlignLeading
	btn.Importance = widget.LowImportance
	return btn
}
