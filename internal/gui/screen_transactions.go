package gui

type TransactionsScreen struct {
	guiApp *GuiApp
}

func NewTransactionsScreen(app *GuiApp) *TransactionsScreen {
	return &TransactionsScreen{guiApp: app}
}
