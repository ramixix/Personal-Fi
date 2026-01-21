package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type TransactionsScreen struct {
	guiApp *GuiApp
}

func NewTransactionsScreen(app *GuiApp) *TransactionsScreen {
	return &TransactionsScreen{guiApp: app}
}

func (t *TransactionsScreen) Render() fyne.CanvasObject {
	title := widget.NewLabelWithStyle("Transactions", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	content := container.NewVBox(
		title,
		widget.NewSeparator(),
		widget.NewLabel("Transaction management page - Coming soon!"),
	)

	return container.NewScroll(content)
}
