package gui

type AccountsScreen struct {
	guiApp *GuiApp
}

func NewAccountsScreen(app *GuiApp) *AccountsScreen {
	return &AccountsScreen{guiApp: app}
}
