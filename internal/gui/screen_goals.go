package gui

type GoalsScreen struct {
	guiApp *GuiApp
}

func NewGoalsScreen(app *GuiApp) *GoalsScreen {
	return &GoalsScreen{guiApp: app}
}
