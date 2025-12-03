package main

import (
	"bufio"
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
		fmt.Println("  trends        Analyze spending trends by category")
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
	case "trends":
		handleTrendsReport()
	case "anomalies":
		handleAnomaliesReport()
	default:
		fmt.Printf("Unknown report type: %s\n", reportType)
		fmt.Println("Available types: monthly, quarterly, yearly, category, compare, trends, anomalies")
	}
}

// Handle monthly report
func handleMonthlyReport() {
	fmt.Println("Monthly Financial Reports")
	fmt.Println("=========================")

	if len(transactions) == 0 {
		fmt.Println("No transactions available.")
		return
	}

	var reports []MonthlyReport = getMonthlyReports()

	if len(reports) == 0 {
		fmt.Println("No monthly data available.")
		return
	}

	fmt.Printf("\n%-15s | %-12s | %-12s | %-12s | %s\n", "Month", "Income", "Expenses", "Net", "Transactions")
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
		fmt.Printf("%-15s | $%-11.2f | $%-11.2f | %s$%-11.2f | %d\n",
			fmt.Sprintf("%s %d", report.Month.String()[:3], report.Year),
			report.Income,
			report.Expenses,
			netSymbol,
			report.Net,
			report.TxCount)

		totalIncome += report.Income
		totalExpenses += report.Expenses
		totalNet += report.Net
		totalTx += report.TxCount

		// Find best and worst months. Maximum and Minimum Net among months
		if report.Net > bestMonth.Net {
			bestMonth = report
		} else if report.Net < worstMonth.Net {
			worstMonth = report
		}
	}

	fmt.Println(strings.Repeat("-", 80))
	fmt.Printf("%-15s | $%-11.2f | $%-11.2f | $%-11.2f | %d\n", "TOTAL", totalIncome, totalExpenses, totalNet, totalTx)

	// Calculate averages
	avgIncome := totalIncome / float64(len(reports))
	avgExpenses := totalExpenses / float64(len(reports))
	avgNet := totalNet / float64(len(reports))

	fmt.Printf("\n%-15s | $%-11.2f | $%-11.2f | $%-11.2f\n", "AVERAGE", avgIncome, avgExpenses, avgNet)

	fmt.Printf("\n--- Insights ---\n")
	fmt.Printf("Best Month:  %s %d ($%.2f net)\n", bestMonth.Month.String(), bestMonth.Year, bestMonth.Net)
	fmt.Printf("Worst Month: %s %d ($%.2f net)\n", worstMonth.Month.String(), worstMonth.Year, worstMonth.Net)

}

// Handle quarterly report
func handleQuarterlyReport() {
	fmt.Println("Quarterly Financial Reports")
	fmt.Println("===========================")

	reader := bufio.NewReader(os.Stdin)
	year, err := getIntInput(reader, "Enter Year: ")
	if err != nil {
		fmt.Println("Invalid year!")
		return
	}

	fmt.Printf("\nQuarterly Report for %d\n", year)
	fmt.Println(strings.Repeat("=", 60))
	fmt.Printf("\n%-10s | %-12s | %-12s | %-12s | %s\n", "Quarter", "Income", "Expenses", "Net", "Transactions")
	fmt.Println(strings.Repeat("-", 70))

	var totalIncome, totalExpenses, totalNet float64
	var totalTx int

	for quarter := 1; quarter <= 4; quarter++ {
		quarterReport := getQuarterlyReport(year, quarter)

		netSymbol := ""
		if quarterReport.Net >= 0 {
			netSymbol = "+"
		}

		fmt.Printf("Q%-9d | $%-11.2f | $%-11.2f | %s$%-11.2f | %d\n",
			quarter,
			quarterReport.Income,
			quarterReport.Expenses,
			netSymbol,
			quarterReport.Net,
			quarterReport.TxCount)

		totalIncome += quarterReport.Income
		totalExpenses += quarterReport.Expenses
		totalNet += quarterReport.Net
		totalTx += quarterReport.TxCount
	}

	fmt.Println(strings.Repeat("-", 70))
	fmt.Printf("%-10s | $%-11.2f | $%-11.2f | $%-11.2f | %d\n", "TOTAL", totalIncome, totalExpenses, totalNet, totalTx)
}

// Handle yearly report
func handleYearlyReport() {
	fmt.Println("Yearly Financial Reports")
	fmt.Println("========================")

	if len(transactions) == 0 {
		fmt.Println("No transactions available.")
		return
	}

	var yearData []MonthlyReport = getYearlyReports()

	fmt.Printf("\n%-10s | %-12s | %-12s | %-12s | %s\n", "Year", "Income", "Expenses", "Net", "Transactions")
	fmt.Println(strings.Repeat("-", 70))

	var totalIncome, totalExpenses, totalNet float64
	var totalTx int

	for _, yearReport := range yearData {
		netSymbol := ""
		if yearReport.Net >= 0 {
			netSymbol = "+"
		}

		fmt.Printf("%-10d | $%-11.2f | $%-11.2f | %s$%-11.2f | %d\n",
			yearReport.Year,
			yearReport.Income,
			yearReport.Expenses,
			netSymbol,
			yearReport.Net,
			yearReport.TxCount)

		totalIncome += yearReport.Income
		totalExpenses += yearReport.Expenses
		totalNet += yearReport.Net
		totalTx += yearReport.TxCount
	}

	fmt.Println(strings.Repeat("-", 70))
	fmt.Printf("%-10s | $%-11.2f | $%-11.2f | $%-11.2f | %d\n", "TOTAL", totalIncome, totalExpenses, totalNet, totalTx)

	// Year-over-year growth
	if len(yearData) >= 2 {
		fmt.Println("\n--- Year-over-Year Growth ---")
		for y := 1; y < len(yearData); y++ {
			currentYear := yearData[y]
			previousYear := yearData[y-1]
			comparison := getYearOverYearComparison(previousYear.Year, currentYear.Year)

			fmt.Printf("  Income:   $%.2f → $%.2f (%+.1f%%)\n", comparison.Period1Income, comparison.Period2Income, comparison.IncomePercent)
			fmt.Printf("  Expenses: $%.2f → $%.2f (%+.1f%%)\n", comparison.Period1Expenses, comparison.Period2Expenses, comparison.ExpensePercent)
		}
	}
}

// Handle category report
func handleCategoryReport() {
	fmt.Println("Category Breakdown Analysis")
	fmt.Println("===========================")

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Analyze: (1) Expenses (2) Income (3) Both")
	desireType, err := getIntInput(reader, "Enter Number Of Type You want: ")
	if err != nil {
		fmt.Println("Invalid Int")
		return
	}

	var transactionType string
	switch desireType {
	case 1:
		transactionType = "expense"
		fmt.Println("\n--- Expense Categories ---")
	case 2:
		transactionType = "income"
		fmt.Println("\n--- Income Categories ---")
	case 3:
		transactionType = ""
		fmt.Println("\n--- All Categories ---")
	default:
		transactionType = "expense"
		fmt.Println("\n--- Expense Categories ---")
	}

	categoryReport := getCategoryBreakdown(transactionType)

	if len(categoryReport) == 0 {
		fmt.Println("No data available.")
		return
	}
	fmt.Printf("\n%-20s | %-12s | %-8s | %s\n", "Category", "Amount", "Count", "Percentage")
	fmt.Println(strings.Repeat("-", 65))

	var total float64
	for _, report := range categoryReport {
		fmt.Printf("%-20s | $%-11.2f | %-8d | %.1f%%\n", report.Category, report.Amount, report.Count, report.Percent)
		total += report.Amount
	}
	fmt.Println(strings.Repeat("-", 65))
	fmt.Printf("%-20s | $%-11.2f\n", "TOTAL", total)

	// Show top 3
	if len(categoryReport) >= 3 {
		fmt.Println("\n--- Top 3 Categories ---")
		for i := 0; i < 3; i++ {
			fmt.Printf("%d. %s: $%.2f (%.1f%%)\n", i+1, categoryReport[i].Category, categoryReport[i].Amount, categoryReport[i].Percent)
		}
	}
}

// Handle comparison report
func handleComparisonReport() {
	fmt.Println("Period Comparison")
	fmt.Println("=================")

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Compare: (1) Month-over-Month (2) Year-over-Year (3) Custom Periods")
	fmt.Print("Choice: ")
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)

	switch choice {
	case "1":
		handleMonthOverMonthComparison(reader)
	case "2":
		handleYearOverYearComparisonInteractive(reader)
	case "3":
		handleCustomPeriodComparison(reader)
	default:
		fmt.Println("Invalid choice!")
	}
}

// Handle month-over-month comparison
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

	var comparison ComparisonReport = getMonthOverMonthComparison(int(year1), time.Month(month1), int(year2), time.Month(month2))

	fmt.Printf("\nComparison: %s %d vs %s %d\n", time.Month(month1).String(), year1, time.Month(month2).String(), year2)
	fmt.Println(strings.Repeat("=", 60))

	displayComparisonReport(comparison)
}

// Handle year-over-year comparison (interactive)
func handleYearOverYearComparisonInteractive(reader *bufio.Reader) {
	fmt.Println("\nYear-over-Year Comparison")
	fmt.Println("-------------------------")

	year1, err := getIntInput(reader, "First Year: ")
	if err != nil {
		fmt.Println("Invalid year!")
		return
	}

	year2, err := getIntInput(reader, "Second Year: ")
	if err != nil {
		fmt.Println("Invalid year!")
		return
	}

	comparison := getYearOverYearComparison(year1, year2)

	fmt.Printf("\nComparison: %d vs %d\n", year1, year2)
	fmt.Println(strings.Repeat("=", 60))

	displayComparisonReport(comparison)
}

// Handle custom period comparison
func handleCustomPeriodComparison(reader *bufio.Reader) {
	fmt.Println("\nCustom Period Comparison")
	fmt.Println("------------------------")

	// Period 1
	fmt.Println("\nPeriod 1:")
	fmt.Print("  Start date (YYYY-MM-DD): ")
	start1Input, _ := reader.ReadString('\n')
	start1, err := parseDate(start1Input)
	if err != nil {
		fmt.Println("Invalid date!")
		return
	}

	fmt.Print("  End date (YYYY-MM-DD): ")
	end1Input, _ := reader.ReadString('\n')
	end1, err := parseDate(strings.TrimSpace(end1Input))
	if err != nil {
		fmt.Println("Invalid date!")
		return
	}

	// Period 2
	fmt.Println("\nPeriod 2:")
	fmt.Print("  Start date (YYYY-MM-DD): ")
	start2Input, _ := reader.ReadString('\n')
	start2, err := parseDate(strings.TrimSpace(start2Input))
	if err != nil {
		fmt.Println("Invalid date!")
		return
	}

	fmt.Print("  End date (YYYY-MM-DD): ")
	end2Input, _ := reader.ReadString('\n')
	end2, err := parseDate(strings.TrimSpace(end2Input))
	if err != nil {
		fmt.Println("Invalid date!")
		return
	}

	comparison := comparePeriods(start1, end1, start2, end2)

	fmt.Printf("\nComparison:\n")
	fmt.Printf("Period 1: %s to %s\n", start1.Format("2006-01-02"), end1.Format("2006-01-02"))
	fmt.Printf("Period 2: %s to %s\n", start2.Format("2006-01-02"), end2.Format("2006-01-02"))
	fmt.Println(strings.Repeat("=", 60))

	displayComparisonReport(comparison)
}

// Display comparison report (helper function)
func displayComparisonReport(comparison ComparisonReport) {
	fmt.Printf("\n%-15s | %-15s | %-15s | %s\n", "", "Period 1", "Period 2", "Change")
	fmt.Println(strings.Repeat("-", 70))

	// Income
	incomeSymbol := ""
	if comparison.IncomeChange >= 0 {
		incomeSymbol = "+"
	}
	fmt.Printf("%-15s | $%-14.2f | $%-14.2f | %s$%.2f (%+.1f%%)\n",
		"Income",
		comparison.Period1Income,
		comparison.Period2Income,
		incomeSymbol,
		comparison.IncomeChange,
		comparison.IncomePercent)

	// Expenses
	expenseSymbol := ""
	if comparison.ExpenseChange >= 0 {
		expenseSymbol = "+"
	}
	fmt.Printf("%-15s | $%-14.2f | $%-14.2f | %s$%.2f (%+.1f%%)\n",
		"Expenses",
		comparison.Period1Expenses,
		comparison.Period2Expenses,
		expenseSymbol,
		comparison.ExpenseChange,
		comparison.ExpensePercent)

	// Net
	net1 := comparison.Period1Income - comparison.Period1Expenses
	net2 := comparison.Period2Income - comparison.Period2Expenses
	netChange := net2 - net1

	netSymbol := ""
	if netChange >= 0 {
		netSymbol = "+"
	}

	var netPercent float64
	if net1 != 0 {
		netPercent = (netChange / net1) * 100
	}

	fmt.Printf("%-15s | $%-14.2f | $%-14.2f | %s$%.2f (%+.1f%%)\n", "Net", net1, net2, netSymbol, netChange, netPercent)

	// Insights
	fmt.Println("\n--- Insights ---")
	if comparison.IncomeChange > 0 {
		fmt.Printf("✓ Income increased by $%.2f\n", comparison.IncomeChange)
	} else if comparison.IncomeChange < 0 {
		fmt.Printf("⚠ Income decreased by $%.2f\n", -comparison.IncomeChange)
	}

	if comparison.ExpenseChange > 0 {
		fmt.Printf("⚠ Expenses increased by $%.2f\n", comparison.ExpenseChange)
	} else if comparison.ExpenseChange < 0 {
		fmt.Printf("✓ Expenses decreased by $%.2f\n", -comparison.ExpenseChange)
	}

	if netChange > 0 {
		fmt.Printf("✓ Overall improvement: +$%.2f\n", netChange)
	} else if netChange < 0 {
		fmt.Printf("⚠ Overall decline: $%.2f\n", netChange)
	}
}

// Handle anomalies report
func handleAnomaliesReport() {
	fmt.Println("Unusual Spending Detection (find expense transactions that are x time of average)")
	fmt.Println("=================================================================================")

	fmt.Println("Threshold value will multiply by average expense of all time and then transactions above that value will be returned. For example if 2 then any expense transaction that is above 2xAverage will be considered as Unusual Spending. (Defualt 3.0) ")
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
	var anomalyTransactions []Transaction = detectHighSpending(threshold)

	if len(anomalyTransactions) == 0 {
		fmt.Printf("\n✓ No unusual spending detected (threshold: %.1fx average)\n", threshold)
		return
	}

	fmt.Printf("\n⚠ Found %d unusual transactions (threshold: %.1fx average):\n\n", len(anomalyTransactions), threshold)

	listFilteredTransactions(anomalyTransactions, "Unusual Spending")

}

// Handle trends report
func handleTrendsReport() {
	fmt.Println("Spending Trends Analysis")
	fmt.Println("========================")

	categories := getCategories()

	if len(categories) == 0 {
		fmt.Println("No categories available.")
		return
	}

	fmt.Println("\nAnalyzing trends over the last 6 months...\n")

	fmt.Printf("%-20s | %-15s | %s\n", "Category", "Trend", "Status")
	fmt.Println(strings.Repeat("-", 60))

	for _, category := range categories {
		trend := getCategoryTrend(category, 6)

		trendSymbol := "→"
		if trend == "increasing" {
			trendSymbol = "↑"
		} else if trend == "decreasing" {
			trendSymbol = "↓"
		}

		fmt.Printf("%-20s | %-15s | %s\n", category, trend, trendSymbol)
	}

	fmt.Println("\n💡 Legend: ↑ Increasing | ↓ Decreasing | → Stable")
}
