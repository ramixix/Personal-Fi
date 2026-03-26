package gui

import (
	"financial_tracker/internal/storage"
	"fmt"
	"os"
	"time"

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

	//display danger zone (a zone where users can delete some parts or information / whole file)
	dangerZoneSection := s.createDangerZoneSection()

	content := container.NewVBox(
		headerPlusAppInfo,
		widget.NewSeparator(),
		dataManagementSection,
		widget.NewSeparator(),
		statisticsSection,
		widget.NewSeparator(),
		displaySettingsSection,
		widget.NewSeparator(),
		dangerZoneSection,
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
		s.createBackup()
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
		s.exportAllData()
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

// -------------------------------------------------------------------------------------------
//
// creates display preferences, where users can change theme, currancy and date format
//
// -------------------------------------------------------------------------------------------
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

// -----------------------------------------------------------------------------------------------
//
// creates danger zone, where users can remove a part/whole of information or delete whole file
//
// -----------------------------------------------------------------------------------------------
func (s *SettingsScreen) createDangerZoneSection() fyne.CanvasObject {
	sectionTitle := widget.NewLabelWithStyle("⚠️ Danger Zone", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	warningLabel := widget.NewLabel("⚠️ Warning: These actions cannot be undone!")
	warningLabel.TextStyle = fyne.TextStyle{Bold: true}

	// Clear all transactions
	clearTransactionsBtn := widget.NewButton("🗑️ Clear All Transactions", func() {
		// s.confirmClearTransactions()
	})
	clearTransactionsBtn.Importance = widget.DangerImportance

	// Clear all data
	clearAllBtn := widget.NewButton("💣 Clear ALL Data", func() {
		// s.confirmClearAllData()
	})
	clearAllBtn.Importance = widget.DangerImportance

	// Delete data file
	deleteFileBtn := widget.NewButton("🔥 Delete Data File", func() {
		// s.confirmDeleteDataFile()
	})
	deleteFileBtn.Importance = widget.DangerImportance

	content := container.NewVBox(
		sectionTitle,
		warningLabel,
		widget.NewLabel("Clear all transactions (keeps accounts and goals)"),
		clearTransactionsBtn,
		widget.NewLabel("Clear all data (transactions, accounts, goals, everything)"),
		clearAllBtn,
		widget.NewLabel("Delete the data file from disk"),
		deleteFileBtn,
	)

	return widget.NewCard("", "", content)
}

// ------------------------------------
// create a timestamped backup
// ------------------------------------
func (s *SettingsScreen) createBackup() {
	timeStamp := time.Now().Format("2006-01-02_15-04-05")
	backupFile := fmt.Sprintf("financial_data_backup_%s.json", timeStamp)

	data, err := os.ReadFile(storage.DataFile)
	if err != nil {
		dialog.ShowError(fmt.Errorf("failed to read data file: %v", err), s.guiApp.GuiWindow)
		return
	}

	err = os.WriteFile(backupFile, data, 0644)
	if err != nil {
		dialog.ShowError(fmt.Errorf("failed to create backup: %v", err), s.guiApp.GuiWindow)
		return
	}

	dialog.ShowInformation("Backup Created",
		fmt.Sprintf("Backup created successfully!\n\nFile: %s\nSize: %.2f KB",
			backupFile, float64(len(data))/1024),
		s.guiApp.GuiWindow)
}

// -----------------------------------------
// restore data from a backup file
// -----------------------------------------
func (s *SettingsScreen) restoreFromBackup() {
	// For now, show instructions
	// In a full implementation, you'd use a file picker
	dialog.ShowInformation("Restore from Backup",
		"To restore from a backup:\n\n"+
			"1. Close the application\n"+
			"2. Rename your backup file to 'financial_data.json'\n"+
			"3. Replace the current data file\n"+
			"4. Restart the application\n\n"+
			"A file picker will be added in a future version.",
		s.guiApp.GuiWindow)
}

// --------------------------------------------------------------------------------------------------------
// exporting all data to CSV files (Transactions, Acounts, Goals each with their dedicated csv files)
// --------------------------------------------------------------------------------------------------------
func (s *SettingsScreen) exportAllData() {
	transactionsCSVFile := "transactions_export.csv"
	accountsCSVFile := "accounts_export.csv"
	goalsCSVFile := "goals_export.csv"

	// Export transactions
	err := s.exportTransactionsCSV(transactionsCSVFile)
	if err != nil {
		dialog.ShowError(fmt.Errorf("failed to export transactions: %v", err), s.guiApp.GuiWindow)
		return
	}

	// Export accounts
	err = s.exportAccountsCSV(accountsCSVFile)
	if err != nil {
		dialog.ShowError(fmt.Errorf("failed to export accounts: %v", err), s.guiApp.GuiWindow)
		return
	}

	// Export goals
	err = s.exportGoalsCSV(goalsCSVFile)
	if err != nil {
		dialog.ShowError(fmt.Errorf("failed to export goals: %v", err), s.guiApp.GuiWindow)
		return
	}

	successMessage := fmt.Sprintf("All data exported successfully!\n\nFiles created:\n- %s\n- %s\n- %s", transactionsCSVFile, accountsCSVFile, goalsCSVFile)
	dialog.ShowInformation("Export Complete", successMessage, s.guiApp.GuiWindow)
}

// ---------------------------------------------------------------
// Helper funcitons to export Transactions, Accounts, and Goals
// ---------------------------------------------------------------
func (s *SettingsScreen) exportTransactionsCSV(exportedFileName string) error {
	file, err := os.Create(exportedFileName)
	if err != nil {
		return err
	}
	defer file.Close()

	// Header
	file.WriteString("ID,Date,Type,Amount,Category,Description")
	// Write data
	for _, transac := range storage.Transactions {
		file.WriteString(fmt.Sprintf("%d,%s,%s,%.2f,%s,%s\n", transac.ID, transac.Date.Format("2006-01-02"), transac.Type, transac.Amount, transac.Category, transac.Description))
	}
	return nil
}

func (s *SettingsScreen) exportAccountsCSV(exportedFileName string) error {
	file, err := os.Create(exportedFileName)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write header
	file.WriteString("ID,Name,Balance,Created\n")
	// Write data
	for _, acc := range storage.Accounts {
		file.WriteString(fmt.Sprintf("%d,%s,%.2f,%s\n", acc.ID, acc.Name, acc.Balance, acc.Created.Format("2006-01-02")))
	}
	return nil
}

func (s *SettingsScreen) exportGoalsCSV(exportedFileName string) error {
	file, err := os.Create(exportedFileName)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write header
	file.WriteString("ID,Name,TargetAmount,CurrentAmount,Status,Created\n")
	// Write data
	for _, goal := range storage.Goals {
		file.WriteString(fmt.Sprintf("%d,%s,%.2f,%.2f,%s,%s\n",
			goal.ID,
			goal.Name,
			goal.TargetAmount,
			goal.CurrentAmount,
			goal.Status,
			goal.Created.Format("2006-01-02"),
		))
	}
	return nil
}
