package gui

import (
	"financial_tracker/internal/core"
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
		s.restoreFromBackup()
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
		s.reloadData()
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

	transactionsNum := core.GetTransactionsLength("")
	accountsNum := core.GetAccountsLength()
	accountsTransactionNum := core.GetAccountsTransactionsLength()
	goalsNum := core.GetGoalsLength("")
	goalContributionNum := core.GetGoalContributionsLength()

	totalRecords := transactionsNum + accountsNum + accountsTransactionNum + goalsNum + goalContributionNum

	statsGrid := container.NewGridWithColumns(6,
		widget.NewLabel("Total Records"),
		widget.NewLabel("Transactions"),
		widget.NewLabel("Accounts"),
		widget.NewLabel("Goals"),
		widget.NewLabel("Account Transactions"),
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
		s.confirmClearTransactions()
	})
	clearTransactionsBtn.Importance = widget.DangerImportance

	// Clear all data
	clearAllBtn := widget.NewButton("💣 Clear ALL Data", func() {
		s.confirmClearAllData()
	})
	clearAllBtn.Importance = widget.DangerImportance

	// Delete data file
	deleteFileBtn := widget.NewButton("🔥 Delete Data File", func() {
		s.confirmDeleteDataFile()
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
	// timeStamp := time.Now().Format("2006-01-02_15-04-05")
	// backupFile := fmt.Sprintf("financial_data_backup_%s.json", timeStamp)

	// data, err := os.ReadFile(storage.DataFile)
	// if err != nil {
	// 	dialog.ShowError(fmt.Errorf("failed to read data file: %v", err), s.guiApp.GuiWindow)
	// 	return
	// }

	// err = os.WriteFile(backupFile, data, 0644)
	// if err != nil {
	// 	dialog.ShowError(fmt.Errorf("failed to create backup: %v", err), s.guiApp.GuiWindow)
	// 	return
	// }

	// dialog.ShowInformation("Backup Created",
	// 	fmt.Sprintf("Backup created successfully!\n\nFile: %s\nSize: %.2f KB",
	// 		backupFile, float64(len(data))/1024),
	// 	s.guiApp.GuiWindow)
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
	err := core.ExportTransactionsToCSV(transactionsCSVFile)
	if err != nil {
		dialog.ShowError(fmt.Errorf("failed to export transactions: %v", err), s.guiApp.GuiWindow)
		return
	}

	// Export accounts
	err = core.ExportAccountsToCSV(accountsCSVFile)
	if err != nil {
		dialog.ShowError(fmt.Errorf("failed to export accounts: %v", err), s.guiApp.GuiWindow)
		return
	}

	// Export goals
	err = core.ExportGoalsToCSV(goalsCSVFile)
	if err != nil {
		dialog.ShowError(fmt.Errorf("failed to export goals: %v", err), s.guiApp.GuiWindow)
		return
	}

	successMessage := fmt.Sprintf("All data exported successfully!\n\nFiles created:\n- %s\n- %s\n- %s", transactionsCSVFile, accountsCSVFile, goalsCSVFile)
	dialog.ShowInformation("Export Complete", successMessage, s.guiApp.GuiWindow)
}

// -----------------------------
// Reloading data from file
// -----------------------------
func (s *SettingsScreen) reloadData() {
	// err := storage.LoadData()
	// if err != nil {
	// 	dialog.ShowError(fmt.Errorf("failed to reload data: %v", err), s.guiApp.GuiWindow)
	// 	return
	// }

	// dialog.ShowInformation("Data Reloaded", "Data reloaded successfully from disk!", s.guiApp.GuiWindow)
	// s.guiApp.ShowDashboardScreen()
}

// ----------------------------------------------------------
// Ask for confirmation before clearing transactions
// ----------------------------------------------------------
func (s *SettingsScreen) confirmClearTransactions() {
	// dialog.ShowConfirm(
	// 	"Clear All Transactions",
	// 	fmt.Sprintf("Are you sure you want to delete ALL %d transactions?\n\nThis action cannot be undone!", len(storage.Transactions)),
	// 	func(confirmed bool) {
	// 		if confirmed {
	// 			storage.Transactions = []models.Transaction{}
	// 			dialog.ShowInformation("Cleared", "All transactions have been deleted", s.guiApp.GuiWindow)
	// 			s.guiApp.ShowSettingsScreen()
	// 		}
	// 	},
	// 	s.guiApp.GuiWindow)
}

// ---------------------------------------------------
// Ask for confirmation before clearing all data
// ---------------------------------------------------
func (s *SettingsScreen) confirmClearAllData() {
	// confirmationMessage := fmt.Sprintf(
	// 	"WARNING: You are about to delete ALL data!\n\n"+
	// 		"This includes:\n"+
	// 		"- %d Transactions\n"+
	// 		"- %d Accounts\n"+
	// 		"- %d Goals\n"+
	// 		"- All contribution history\n\n"+
	// 		"This action CANNOT be undone!\n\n"+
	// 		"Are you absolutely sure?",
	// 	len(storage.Transactions),
	// 	len(storage.Accounts),
	// 	len(storage.Goals),
	// )

	// dialog.ShowConfirm(
	// 	"⚠️ Clear ALL Data",
	// 	confirmationMessage,
	// 	func(confirmed bool) {
	// 		if confirmed {
	// 			// Double confirmation
	// 			dialog.ShowConfirm(
	// 				"Final Confirmation",
	// 				"This is your last chance!\n\nDelete all data permanently?",
	// 				func(reallyConfirmed bool) {
	// 					if reallyConfirmed {
	// 						s.clearAllData()
	// 					}
	// 				},
	// 				s.guiApp.GuiWindow,
	// 			)
	// 		}
	// 	},
	// 	s.guiApp.GuiWindow,
	// )
}

// --------------------------------
// clear all application data
// --------------------------------
func (s *SettingsScreen) clearAllData() {
	// storage.Transactions = []models.Transaction{}
	// storage.Accounts = []models.Account{}
	// storage.Goals = []models.Goal{}
	// storage.AccountTransactions = []models.AccountTransaction{}
	// storage.GoalContributions = []models.GoalContribution{}

	// storage.NextTransactionID = 1
	// storage.NextAccountID = 1
	// storage.NextGoalID = 1
	// storage.NextAccountTransactionID = 1
	// storage.NextGoalContributionID = 1

	// storage.SaveData()

	// dialog.ShowInformation("All Data Cleared", "All application data has been permanently deleted", s.guiApp.GuiWindow)
	// s.guiApp.ShowDashboardScreen()
}

// -------------------------------------------------------
// Ask for confirmation before deleting data file
// -------------------------------------------------------
func (s *SettingsScreen) confirmDeleteDataFile() {
	// dialog.ShowConfirm(
	// 	"Delete Data File",
	// 	"This will delete the data file from disk.\nThe application will start fresh next time.\n\nContinue?",
	// 	func(confirmed bool) {
	// 		if confirmed {
	// 			err := os.Remove(storage.DataFile)
	// 			if err != nil {
	// 				dialog.ShowError(fmt.Errorf("failed to delete file: %v", err), s.guiApp.GuiWindow)
	// 				return
	// 			}

	// 			dialog.ShowInformation("File Deleted", "Data file deleted successfully.\nThe application will restart with empty data.", s.guiApp.GuiWindow)
	// 			// Clear in-memory data too
	// 			s.clearAllData()
	// 		}
	// 	},
	// 	s.guiApp.GuiWindow,
	// )
}
