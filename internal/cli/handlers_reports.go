package cli

import (
	"bufio"
	"financial_tracker/internal/core"
	"financial_tracker/internal/models"
	"financial_tracker/internal/utils"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// Handle reports command with subcommands
func handleReports() {
	if len(os.Args) < 3 {
		fmt.Println("Reports & Analytics")
		fmt.Println("===================")
		fmt.Printf("Usage: %s reports [type]\n", filepath.Base(os.Args[0]))
		fmt.Println("Types:")
		fmt.Println("  monthly       Monthly financial reports")
		fmt.Println("  quarterly     Quarterly financial reports")
		fmt.Println("  yearly        Yearly financial reports")
		fmt.Println("  category      Category breakdown analysis")
		fmt.Println("  compare       Compare two time periods")
		fmt.Println("  anomalies     Detect unusual spending")
		return
	}

	reportType := strings.ToLower(os.Args[2])
	switch reportType {
	case "monthly":
		handleMonthlyReport()
	case "quarterly":
		handleQuarterlyReport()
	case "yearly":
		handleYearlyReport()
	case "category":
		handleCategoryReport()
	case "compare":
		handleComparisonReport()
	case "anomalies":
		handleAnomaliesReport()
	default:
		fmt.Printf("Unknown report type: %s\n", reportType)
		fmt.Println("Available types: monthly, quarterly, yearly, category, compare, trends, anomalies")
	}
}

// handleMonthlyReport prints monthly reports
func handleMonthlyReport() {
	fmt.Println("Monthly Financial Reports")
	fmt.Println("=========================")

	reportsMap := core.GetTransactionsMonthlyReports()

	if len(reportsMap) == 0 {
		fmt.Println("No monthly data available.")
		return
	}

	for currency, reports := range reportsMap {
		if len(reports) == 0 {
			continue
		}

		fmt.Printf("\n=== %s REPORT ===\n", currency)
		fmt.Printf("%-15s | %-12s | %-12s | %-12s | %s\n", "Month", "Income", "Expenses", "Net", "Transactions")
		fmt.Println(strings.Repeat("-", 80))

		bestMonth := reports[0]
		worstMonth := reports[0]

		var totalIncome, totalExpenses, totalNet float64
		var totalTx int

		for _, report := range reports {
			netSymbol := ""
			if report.Net >= 0 {
				netSymbol = "+"
			}
			fmt.Printf("%-15s | %-11s | %-11s | %s%-11s | %d\n",
				fmt.Sprintf("%s %d", report.Month.String()[:3], report.Year),
				utils.FormatCurrency(report.Income, currency),
				utils.FormatCurrency(report.Expenses, currency),
				netSymbol,
				utils.FormatCurrency(report.Net, currency),
				report.TxCount)

			totalIncome += report.Income
			totalExpenses += report.Expenses
			totalNet += report.Net
			totalTx += report.TxCount

			if report.Net > bestMonth.Net {
				bestMonth = report
			} else if report.Net < worstMonth.Net {
				worstMonth = report
			}
		}

		fmt.Println(strings.Repeat("-", 80))
		fmt.Printf("%-15s | %-11s | %-11s | %-11s | %d\n", "TOTAL",
			utils.FormatCurrency(totalIncome, currency),
			utils.FormatCurrency(totalExpenses, currency),
			utils.FormatCurrency(totalNet, currency),
			totalTx)

		avgIncome := totalIncome / float64(len(reports))
		avgExpenses := totalExpenses / float64(len(reports))
		avgNet := totalNet / float64(len(reports))

		fmt.Printf("\n%-15s | %-11s | %-11s | %-11s\n", "AVERAGE",
			utils.FormatCurrency(avgIncome, currency),
			utils.FormatCurrency(avgExpenses, currency),
			utils.FormatCurrency(avgNet, currency))

		fmt.Printf("\n--- %s Insights ---\n", currency)
		fmt.Printf("Best Month:  %s %d (%.2f net)\n", bestMonth.Month.String(), bestMonth.Year, bestMonth.Net)
		fmt.Printf("Worst Month: %s %d (%.2f net)\n", worstMonth.Month.String(), worstMonth.Year, worstMonth.Net)
		fmt.Println()
	}
}

// handleYearlyReport prints yearly reports
func handleYearlyReport() {
	fmt.Println("Yearly Financial Reports")
	fmt.Println("========================")

	reportsMap := core.GetTransactionsYearlyReports()
	if len(reportsMap) == 0 {
		fmt.Println("No monthly data available.")
		return
	}

	for currency, reports := range reportsMap {
		if len(reports) == 0 {
			continue
		}

		fmt.Printf("\n=== %s REPORT ===\n", currency)
		fmt.Printf("\n%-10s | %-12s | %-12s | %-12s | %s\n", "Year", "Income", "Expenses", "Net", "Transactions")
		fmt.Println(strings.Repeat("-", 70))

		var totalIncome, totalExpenses, totalNet float64
		var totalTx int

		for _, yearReport := range reports {
			netSymbol := ""
			if yearReport.Net >= 0 {
				netSymbol = "+"
			}

			fmt.Printf("%-10d | %-11s | %-11s | %s%-11s | %d\n",
				yearReport.Year,
				utils.FormatCurrency(yearReport.Income, currency),
				utils.FormatCurrency(yearReport.Expenses, currency),
				netSymbol,
				utils.FormatCurrency(yearReport.Net, currency),
				yearReport.TxCount)

			totalIncome += yearReport.Income
			totalExpenses += yearReport.Expenses
			totalNet += yearReport.Net
			totalTx += yearReport.TxCount
		}

		fmt.Println(strings.Repeat("-", 70))
		fmt.Printf("%-10s | %-11.2f | %-11.2f | %-11.2f | %d\n", "TOTAL",
			utils.FormatCurrency(totalIncome, currency),
			utils.FormatCurrency(totalExpenses, currency),
			utils.FormatCurrency(totalNet, currency),
			totalTx)

		// Year-over-year growth
		if len(reports) >= 2 {
			fmt.Println("\n--- Year-over-Year Growth ---")
			for y := 1; y < len(reports); y++ {
				currentYear := reports[y]
				previousYear := reports[y-1]
				comparison := core.GetYearOverYearComparison(previousYear.Year, currentYear.Year)
				for currency, report := range comparison {
					fmt.Printf("  Income:   %.2f → %.2f (%+.1f%%)\n",
						utils.FormatCurrency(report.Period1Income, currency),
						utils.FormatCurrency(report.Period2Income, currency),
						report.IncomePercent)
					fmt.Printf("  Expenses: %.2f → %.2f (%+.1f%%)\n",
						utils.FormatCurrency(report.Period1Expenses, currency),
						utils.FormatCurrency(report.Period2Expenses, currency),
						report.ExpensePercent)
				}
			}
		}
	}
}

// handleQuarterlyReport prints each quarter report for users given year
func handleQuarterlyReport() {
	fmt.Println("Quarterly Financial Reports")
	fmt.Println("===========================")

	reader := bufio.NewReader(os.Stdin)
	year, err := utils.GetIntInput(reader, "Enter Year: ")
	if err != nil || len(strconv.Itoa(year)) != 4 {
		fmt.Println("Invalid year!")
		return
	}
	fmt.Printf("\nQuarterly Report for %d\n", year)
	fmt.Println(strings.Repeat("=", 60))

	totalIncome := make(map[string]float64)
	totalExpenses := make(map[string]float64)
	totalTx := 0

	for quarter := 1; quarter <= 4; quarter++ {
		fmt.Printf("\n%-10s | %s | %-12s | %-12s | %-12s | %s\n", "Quarter", "Currency", "Income", "Expenses", "Net", "Transactions")
		fmt.Println(strings.Repeat("-", 70))
		quarterReport := core.GetTransactionsQuarterlyReport(year, quarter)

		for currency, quarterReport := range quarterReport {
			netSymbol := ""
			if quarterReport.Net >= 0 {
				netSymbol = "+"
			}
			fmt.Printf("Q%-9d | %-5s | %-11s | %-11s | %s%-11s | %d\n",
				quarter,
				quarterReport.CurrencyCode,
				utils.FormatCurrency(quarterReport.Income, currency),
				utils.FormatCurrency(quarterReport.Expenses, currency),
				netSymbol,
				utils.FormatCurrency(quarterReport.Net, currency),
				quarterReport.TxCount)

			totalIncome[currency] += quarterReport.Income
			totalExpenses[currency] += quarterReport.Expenses
			totalTx += quarterReport.TxCount
		}
	}

	fmt.Printf("Totals (%d transactions of year %d)", totalTx, year)
	for currency, total := range totalIncome {
		fmt.Printf("Currency: %-5s | Total Income: %-12s\n", currency, utils.FormatCurrency(total, currency))
	}

	fmt.Print("\n")
	for currency, total := range totalIncome {
		fmt.Printf("Currency: %-5s | Total Expense: %-12s\n", currency, utils.FormatCurrency(total, currency))
	}
}

// handleCategoryReport prints category report
func handleCategoryReport() {
	showCategories()
	// fmt.Println("Category Breakdown Analysis")
	// fmt.Println("===========================")

	// reader := bufio.NewReader(os.Stdin)
	// fmt.Println("Select what you want to analyze:")
	// fmt.Println("\t1) Expenses")
	// fmt.Println("\t2) Income")
	// fmt.Println("\t3) Both")
	// desireType, err := utils.GetIntInput(reader, "Enter the number corresponding to your choice: ")
	// if err != nil {
	// 	fmt.Println("Invalid Int")
	// 	return
	// }

	// var transactionType string
	// switch desireType {
	// case 1:
	// 	transactionType = "expense"
	// 	fmt.Println("\n--- Expense Categories ---")
	// case 2:
	// 	transactionType = "income"
	// 	fmt.Println("\n--- Income Categories ---")
	// case 3:
	// 	transactionType = ""
	// 	fmt.Println("\n--- All Categories ---")
	// default:
	// 	transactionType = "expense"
	// 	fmt.Println("\n--- Expense Categories ---")
	// }

	// categoryReport := core.GetCategoryBreakdown(transactionType)

	// if len(categoryReport) == 0 {
	// 	fmt.Println("No data available.")
	// 	return
	// }
	// fmt.Printf("\n%-20s | %-12s | %-8s | %s\n", "Category", "Amount", "Count", "Percentage")
	// fmt.Println(strings.Repeat("-", 65))

	// var total float64
	// for _, report := range categoryReport {
	// 	fmt.Printf("%-20s | %-11.2f | %-8d | %.1f%%\n", report.Category, report.Amount, report.Count, report.Percent)
	// 	total += report.Amount
	// }
	// fmt.Println(strings.Repeat("-", 65))
	// fmt.Printf("%-20s | %-11.2f\n", "TOTAL", total)

	// // Show top 3
	// if len(categoryReport) >= 3 {
	// 	fmt.Println("\n--- Top 3 Categories ---")
	// 	for i := 0; i < 3; i++ {
	// 		fmt.Printf("%d. %s: %.2f (%.1f%%)\n", i+1, categoryReport[i].Category, categoryReport[i].Amount, categoryReport[i].Percent)
	// 	}
	// }
}

// handleComparisonReport allows users to execute and print different comparison report type
func handleComparisonReport() {
	fmt.Println("Period Comparison")
	fmt.Println("=================")

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Select a comparison type:")
	fmt.Println("\t1) Month-over-Month")
	fmt.Println("\t2) Year-over-Year")
	fmt.Println("\t3) Custom Periods")
	desireType, err := utils.GetIntInput(reader, "Enter the number corresponding to your choice: ")
	if err != nil {
		fmt.Println("Invalid Int")
		return
	}

	switch desireType {
	case 1:
		handleMonthOverMonthComparison(reader)
	case 2:
		handleYearOverYearComparisonInteractive(reader)
	case 3:
		handleCustomPeriodComparison(reader)
	default:
		fmt.Println("Invalid choice!")
	}
}

// handleMonthOverMonthComparison prints month-over-month comparison
func handleMonthOverMonthComparison(reader *bufio.Reader) {
	fmt.Println("\nMonth-over-Month Comparison")
	fmt.Println("---------------------------")

	fmt.Print("First month (YYYY-MM): ")
	input, _ := reader.ReadString('\n')
	var year1, month1 int16
	_, err := fmt.Sscanf(strings.TrimSpace(input), "%d-%d", &year1, &month1)
	if err != nil {
		fmt.Println("Invalid Format. Format should be YYYY-MM for instance 2022-02")
		return
	}

	fmt.Print("Second month (YYYY-MM): ")
	input2, _ := reader.ReadString('\n')
	var year2, month2 int16
	_, err = fmt.Sscanf(strings.TrimSpace(input2), "%d-%d", &year2, &month2)
	if err != nil {
		fmt.Println("Invalid Format. Format should be YYYY-MM for instance 2022-02")
		return
	}

	comparison := core.GetMonthOverMonthComparison(int(year1), time.Month(month1), int(year2), time.Month(month2))

	fmt.Printf("\nComparison: %s %d vs %s %d\n", time.Month(month1).String(), year1, time.Month(month2).String(), year2)
	fmt.Println(strings.Repeat("=", 60))

	displayComparisonReport(comparison)
}

// handleYearOverYearComparisonInteractive prints year-over-year comparison
func handleYearOverYearComparisonInteractive(reader *bufio.Reader) {
	fmt.Println("\nYear-over-Year Comparison")
	fmt.Println("-------------------------")

	year1, err := utils.GetIntInput(reader, "First Year: ")
	if err != nil || len(strconv.Itoa(year1)) != 4 {
		fmt.Println("Invalid year!")
		return
	}

	year2, err := utils.GetIntInput(reader, "Second Year: ")
	if err != nil || len(strconv.Itoa(year2)) != 4 {
		fmt.Println("Invalid year!")
		return
	}

	comparison := core.GetYearOverYearComparison(year1, year2)

	fmt.Printf("\nComparison: %d vs %d\n", year1, year2)
	fmt.Println(strings.Repeat("=", 60))

	displayComparisonReport(comparison)
}

// handleCustomPeriodComparison allows users to enter custom period and prints comparison
func handleCustomPeriodComparison(reader *bufio.Reader) {
	fmt.Println("\nCustom Period Comparison")
	fmt.Println("------------------------")

	// Period 1
	fmt.Println("\nPeriod 1:")
	fmt.Print("  Start date (YYYY-MM-DD): ")
	start1Input, _ := reader.ReadString('\n')
	start1, err := utils.ParseDate(start1Input)
	if err != nil {
		fmt.Println("Invalid date!")
		return
	}

	fmt.Print("  End date (YYYY-MM-DD): ")
	end1Input, _ := reader.ReadString('\n')
	end1, err := utils.ParseDate(strings.TrimSpace(end1Input))
	if err != nil {
		fmt.Println("Invalid date!")
		return
	}

	// Period 2
	fmt.Println("\nPeriod 2:")
	fmt.Print("  Start date (YYYY-MM-DD): ")
	start2Input, _ := reader.ReadString('\n')
	start2, err := utils.ParseDate(strings.TrimSpace(start2Input))
	if err != nil {
		fmt.Println("Invalid date!")
		return
	}

	fmt.Print("  End date (YYYY-MM-DD): ")
	end2Input, _ := reader.ReadString('\n')
	end2, err := utils.ParseDate(strings.TrimSpace(end2Input))
	if err != nil {
		fmt.Println("Invalid date!")
		return
	}

	comparison := core.ComparePeriods(start1, end1, start2, end2)

	fmt.Printf("\nComparison:\n")
	fmt.Printf("Period 1: %s to %s\n", start1.Format("2006-01-02"), end1.Format("2006-01-02"))
	fmt.Printf("Period 2: %s to %s\n", start2.Format("2006-01-02"), end2.Format("2006-01-02"))
	fmt.Println(strings.Repeat("=", 60))

	displayComparisonReport(comparison)
}

// displayComparisonReport prints comparison report result to stdout (helper function)
func displayComparisonReport(comparisonReport map[string]models.ComparisonReport) {
	fmt.Printf("\n%-15s | %-15s | %-15s | %-15s | %s\n", "", "Currency", "Period 1", "Period 2", "Change")
	fmt.Println(strings.Repeat("-", 70))
	for currency, report := range comparisonReport {
		// Income
		incomeSymbol := ""
		if report.IncomeChange >= 0 {
			incomeSymbol = "+"
		}
		fmt.Printf("%-15s | %-14.2f | %-14.2f | %s%.2f (%+.1f%%)\n",
			"Income",
			utils.FormatCurrency(report.Period1Income, currency),
			utils.FormatCurrency(report.Period2Income, currency),
			incomeSymbol,
			utils.FormatCurrency(report.IncomeChange, currency),
			report.IncomePercent)

		// Expenses
		expenseSymbol := ""
		if report.ExpenseChange >= 0 {
			expenseSymbol = "+"
		}
		fmt.Printf("%-15s | %-14.2f | %-14.2f | %s%.2f (%+.1f%%)\n",
			"Expenses",
			utils.FormatCurrency(report.Period1Expenses, currency),
			utils.FormatCurrency(report.Period2Expenses, currency),
			expenseSymbol,
			utils.FormatCurrency(report.ExpenseChange, currency),
			report.ExpensePercent)

		// Net
		net1 := report.Period1Income - report.Period1Expenses
		net2 := report.Period2Income - report.Period2Expenses
		netChange := net2 - net1

		netSymbol := ""
		if netChange >= 0 {
			netSymbol = "+"

		}
		var netPercent float64
		if net1 != 0 {
			netPercent = (netChange / net1) * 100
		}
		fmt.Printf("%-15s | %-14.2f | %-14.2f | %s%.2f (%+.1f%%)\n",
			"Net",
			utils.FormatCurrency(net1, currency),
			utils.FormatCurrency(net2, currency),
			netSymbol,
			utils.FormatCurrency(netChange, currency),
			netPercent)

		// Insights
		fmt.Println("\n--- Insights ---")
		if report.IncomeChange > 0 {
			fmt.Printf("✓ Income increased by %.2f\n", report.IncomeChange)
		} else if report.IncomeChange < 0 {
			fmt.Printf("⚠ Income decreased by %.2f\n", -report.IncomeChange)
		}

		if report.ExpenseChange > 0 {
			fmt.Printf("⚠ Expenses increased by %.2f\n", report.ExpenseChange)
		} else if report.ExpenseChange < 0 {
			fmt.Printf("✓ Expenses decreased by %.2f\n", -report.ExpenseChange)
		}

		if netChange > 0 {
			fmt.Printf("✓ Overall improvement: +%.2f\n", netChange)
		} else if netChange < 0 {
			fmt.Printf("⚠ Overall decline: %.2f\n", netChange)
		}
	}
}

// handleAnomaliesReport prints anomaly transactiosn (transactions) that are above user specified threshold.
func handleAnomaliesReport() {
	fmt.Println("Unusual Spending Detection (find expense transactions that are x time of average)")
	fmt.Println("=================================================================================")
	fmt.Println("Threshold value will multiply by average expense of all time and then transactions above that value will be returned. For example if threshold is set to 2 then any expense transaction that is above 2xAverage will be considered as Unusual Spending. (Defualt 3.0) ")
	fmt.Print("Your Input: ")

	reader := bufio.NewReader(os.Stdin)
	thresholdInput, _ := reader.ReadString('\n')
	thresholdInput = strings.TrimSpace(thresholdInput)

	threshold := 3.0
	if thresholdInput != "" {
		parsed, err := strconv.ParseFloat(thresholdInput, 64)
		if err == nil && parsed > 0 {
			threshold = parsed
		}
	}

	fmt.Printf("Finding Expenses Above %.1f X Average.", threshold)
	var anomalyTransactions []models.Transaction = core.DetectHighSpendingTransactions(threshold)

	if len(anomalyTransactions) == 0 {
		fmt.Printf("\n✓ No unusual spending detected (threshold: %.1fx average)\n", threshold)
		return
	}

	fmt.Printf("\n⚠ Found %d unusual transactions (threshold: %.1fx average):\n\n", len(anomalyTransactions), threshold)
	ListFilteredTransactions(anomalyTransactions, "Unusual Spending")
}
