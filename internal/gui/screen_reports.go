package gui

type ReportsScreen struct {
	guiApp *GuiApp
}

func NewReportsScreen(app *GuiApp) *ReportsScreen {
	return &ReportsScreen{guiApp: app}
}
