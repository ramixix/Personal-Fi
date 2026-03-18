package gui

import (
	"financial_tracker/internal/storage"
	"fmt"

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
	// App Info Section
	appInfoSection := s.createAppInfoSection()

	content := container.NewVBox(header, widget.NewSeparator(), appInfoSection)
	return container.NewScroll(content)
}

// -------------------------------------
//
//	Simple header for setting screen
//
// -------------------------------------
func (s *SettingsScreen) createHeader() fyne.CanvasObject {
	title := widget.NewLabelWithStyle("⚙️ Settings", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	subtitle := widget.NewLabel("Manage your application preferences and data")

	return container.NewVBox(title, subtitle)
}

// ------------------------------------------------------------------------
//
//	Section to show information about current version of application
//
// ------------------------------------------------------------------------
// createAppInfoSection creates app information section
func (s *SettingsScreen) createAppInfoSection() fyne.CanvasObject {
	version := widget.NewLabelWithStyle(fmt.Sprintf("Version: %s", storage.AppVersion), fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	appInfo := widget.NewCard("Application Information", "Financial Tracker: A personal finance management application", version)

	return container.NewVBox(appInfo)
}
