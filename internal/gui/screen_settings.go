package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type SettingsScreen struct {
	guiApp *GuiApp
}

func NewSettingsScreen(app *GuiApp) *SettingsScreen {
	return &SettingsScreen{guiApp: app}
}

func (s *SettingsScreen) Render() fyne.CanvasObject {
	header := s.createHeader()

	content := container.NewVBox(header, widget.NewSeparator())
	return container.NewScroll(content)
}

func (s *SettingsScreen) createHeader() fyne.CanvasObject {
	title := widget.NewLabelWithStyle("⚙️ Settings", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	subtitle := widget.NewLabel("Manage your application preferences and data")

	return container.NewVBox(title, subtitle)
}
