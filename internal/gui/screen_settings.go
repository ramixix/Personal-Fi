package gui

import (
	"financial_tracker/internal/storage"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type SettingsScreen struct {
	guiApp *GuiApp
}

func NewSettingsScreen(app *GuiApp) *SettingsScreen {
	return &SettingsScreen{guiApp: app}
}

func (s *SettingsScreen) Render() fyne.CanvasObject {
	headerPlusAppInfo := s.HeaderAndAppInfoSection()

	// Data Management Section
	dataManagementSection := s.createDataManagementSection()

	// Statistics Section
	statisticsSection := s.createStatisticsSection()

	// diplay settings (theme, currancy and date format)
	displaySettingsSection := s.createDisplaySettingsSection()

	content := container.NewVBox(
		headerPlusAppInfo,
		widget.NewSeparator(),
		dataManagementSection,
		widget.NewSeparator(),
		statisticsSection,
		widget.NewSeparator(),
		displaySettingsSection,
	)
	return container.NewScroll(content)
}

// -------------------------------------
//
//	Simple header for setting screen
//
// -------------------------------------
func (s *SettingsScreen) HeaderAndAppInfoSection() fyne.CanvasObject {
	title := widget.NewLabelWithStyle("⚙️ Settings", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	subtitle := widget.NewLabel("Manage your application preferences and data")

	// app version info
	version := widget.NewLabelWithStyle(fmt.Sprintf("Version: %s", storage.AppVersion), fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	appInfo := widget.NewCard("Application Information", "Financial Tracker: A personal finance management application", version)

	return container.NewVBox(title, subtitle, widget.NewSeparator(), appInfo)
}

// ----------------------------------------------------------------------------------------------------------------
//
//	Section for data management functions such as creating/restoring backup, exporting to csv and reloading data
//
// ----------------------------------------------------------------------------------------------------------------
func (s *SettingsScreen) createDataManagementSection() fyne.CanvasObject {
	sectionTitle := widget.NewLabelWithStyle("💾 Data Management", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	// Backup button
	backupBtn := widget.NewButton("📦 Create Backup", func() {
		// s.createBackup()
	})
	backupBtn.Importance = widget.HighImportance
	backupInfo := widget.NewLabel("Create a timestamped backup of your data: ")
	backupGrid := container.NewGridWithColumns(2, backupInfo, backupBtn)

	// Restore button
	restoreBtn := widget.NewButton("📥 Restore from Backup", func() {
		// s.restoreFromBackup()
	})
	restoreInfo := widget.NewLabel("Restore data from a previous backup file: ")
	restoreGrid := container.NewGridWithColumns(2, restoreInfo, restoreBtn)

	// Export all data
	exportBtn := widget.NewButton("📤 Export All Data (CSV)", func() {
		// s.exportAllData()
	})
	exportInfo := widget.NewLabel("Export transactions, accounts, and goals to CSV files: ")
	exportGrid := container.NewGridWithColumns(2, exportInfo, exportBtn)

	// Refresh data
	refreshBtn := widget.NewButton("🔄 Reload Data", func() {
		// s.reloadData()
	})
	refreshInfo := widget.NewLabel("Reload data from disk (useful if file was modified externally): ")
	refreshGrid := container.NewGridWithColumns(2, refreshInfo, refreshBtn)

	content := container.NewVBox(
		sectionTitle,
		backupGrid,
		restoreGrid,
		exportGrid,
		refreshGrid,
	)
	return content
}

// -----------------------------------------------------------------------------------------------------------
//
//	Section for showing breif statistic about the data stored so user can have an idea of what they have
//
// -----------------------------------------------------------------------------------------------------------
func (s *SettingsScreen) createStatisticsSection() fyne.CanvasObject {
	sectionTitle := widget.NewLabelWithStyle("📊 Data Statistics", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	transactionsNum := len(storage.Transactions)
	accountsNum := len(storage.Accounts)
	accountsTransactionNum := len(storage.AccountTransactions)
	goalsNum := len(storage.Goals)
	goalContributionNum := len(storage.GoalContributions)

	totalRecords := transactionsNum + accountsNum + accountsTransactionNum + goalsNum + goalContributionNum

	statsGrid := container.NewGridWithColumns(6,
		widget.NewLabel("Total Records"),
		widget.NewLabel("Transactions"),
		widget.NewLabel("Accounts"),
		widget.NewLabel("Goals"),
		widget.NewLabel("Account  Transactions"),
		widget.NewLabel("Goal Contributions"),

		widget.NewLabel(fmt.Sprintf("%d", totalRecords)),
		widget.NewLabel(fmt.Sprintf("%d", transactionsNum)),
		widget.NewLabel(fmt.Sprintf("%d", accountsNum)),
		widget.NewLabel(fmt.Sprintf("%d", goalsNum)),
		widget.NewLabel(fmt.Sprintf("%d", accountsTransactionNum)),
		widget.NewLabel(fmt.Sprintf("%d", goalContributionNum)),
	)

	return container.NewVBox(sectionTitle, statsGrid)
}

// creates display preferences
func (s *SettingsScreen) createDisplaySettingsSection() fyne.CanvasObject {
	sectionTitle := widget.NewLabelWithStyle("🎨 Display Settings", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	// Theme selector (for future implementation)
	themeLabel := widget.NewLabel("Theme:")
	themeSelect := widget.NewSelect([]string{"Default", "Light", "Dark"}, nil)
	themeSelect.SetSelected("Default")
	themeSelect.OnChanged = func(value string) {
		dialog.ShowInformation("Theme", "Theme switching will be implemented in a future version", s.guiApp.GuiWindow)
	}

	// Currency (for future implementation)
	currencyLabel := widget.NewLabel("Currency:")
	currencySelect := widget.NewSelect([]string{"USD ($)", "EUR (€)", "GBP (£)", "JPY (¥)"}, nil)
	currencySelect.SetSelected("USD ($)")
	currencySelect.OnChanged = func(value string) {
		dialog.ShowInformation("Currency", "Currency switching will be implemented in a future version", s.guiApp.GuiWindow)
	}

	// Date format (for future implementation)
	dateFormatLabel := widget.NewLabel("Date Format:")
	dateFormatSelect := widget.NewSelect([]string{"MM/DD/YYYY", "DD/MM/YYYY", "YYYY-MM-DD"}, nil)
	dateFormatSelect.SetSelected("YYYY-MM-DD")
	dateFormatSelect.OnChanged = func(value string) {
		dialog.ShowInformation("Date Format", "Date format switching will be implemented in a future version", s.guiApp.GuiWindow)
	}

	content := container.NewVBox(
		sectionTitle,
		themeLabel,
		themeSelect,
		currencyLabel,
		currencySelect,
		dateFormatLabel,
		dateFormatSelect,
		widget.NewLabel("💡 Note: These settings will be fully functional in a future update"),
	)
	return widget.NewCard("", "", content)
}
