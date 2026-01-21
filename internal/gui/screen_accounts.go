package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type AccountsScreen struct {
	guiApp *GuiApp
}

func NewAccountsScreen(app *GuiApp) *AccountsScreen {
	return &AccountsScreen{guiApp: app}
}

func (a *AccountsScreen) Render() fyne.CanvasObject {
	title := widget.NewLabelWithStyle("Accounts", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	content := container.NewVBox(
		title,
		widget.NewSeparator(),
		widget.NewLabel("Accounts management page - Coming soon!"),
	)

	return container.NewScroll(content)
}
