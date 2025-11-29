package main

import (
	"fmt"
	"time"
)

func getMonthlyReports() []MonthlyReport {
	monthlyData := make(map[string]*MonthlyReport)

	for _, transac := range transactions {
		var key string = fmt.Sprintf("%d-%02d", transac.Date.Year(), transac.Date.Month())

		// Initialize if not exists
		_, exists := monthlyData[key]
		if !exists {
			monthlyData[key] = &MonthlyReport{Year: transac.Date.Year(), Month: transac.Date.Month()}
		}
		// Update totals
		if transac.Type == "income" {
			monthlyData[key].Income += transac.Amount
		} else {
			monthlyData[key].Expenses += transac.Amount
		}
		monthlyData[key].TxCount++
	}

	// Convert map to slice and calculate net
	var reports []MonthlyReport
	for _, report := range monthlyData {
		report.Net = report.Income - report.Expenses
		reports = append(reports, *report)
	}

	for i := 0; i < len(reports); i++ {
		for j := i + 1; j < len(reports); i++ {
			time1 := time.Date(reports[i].Year, reports[i].Month, 1, 0, 0, 0, 0, time.UTC)
			time2 := time.Date(reports[j].Year, reports[j].Month, 1, 0, 0, 0, 0, time.UTC)
			if time2.Before(time1) {
				reports[i], reports[j] = reports[j], reports[i]
			}
		}
	}
	return reports
}

// Get report for a specific month
func getMonthReport(year int, month time.Month) MonthlyReport {
	report := MonthlyReport{
		Year:  year,
		Month: month,
	}

	for _, transaction := range transactions {
		if transaction.Date.Year() == year && transaction.Date.Month() == month {
			if transaction.Type == "income" {
				report.Income += transaction.Amount
			} else {
				report.Expenses += transaction.Amount
			}
			report.TxCount++
		}
	}

	report.Net = report.Income - report.Expenses
	return report
}

// Get category breakdown report
func getCategoryBreakdown(transactionType string) []CategoryReport {
	categoryData := make(map[string]*CategoryReport)

	var total float64
	for _, transac := range transactions {
		if transactionType != "" && transac.Type != transactionType {
			continue
		}

		var transactionCategory string = transac.Category
		_, doesExists := categoryData[transactionCategory]
		if !doesExists {
			categoryData[transactionCategory] = &CategoryReport{Category: transactionCategory}
		}
		categoryData[transactionCategory].Amount = transac.Amount
		categoryData[transactionCategory].Count++
		total += transac.Amount
	}

	// Convert to slice and calculate percentages
	var reports []CategoryReport
	for _, report := range categoryData {
		if total > 0 {
			report.Percent = (report.Amount / total) * 100
		}
		reports = append(reports, *report)
	}

	// Sort by amount (descending)
	for i := 0; i < len(reports); i++ {
		for j := i + 1; j < len(reports); j++ {
			if reports[j].Amount > reports[i].Amount {
				reports[i], reports[j] = reports[j], reports[i]
			}
		}
	}

	return reports
}

// Compare two time periods
func comparePeriods(start1, end1, start2, end2 time.Time) ComparisonReport {
	var report ComparisonReport

	for _, transaction := range transactions {
		// Period 1
		if (transaction.Date.After(start1) || transaction.Date.Equal(start1)) &&
			(transaction.Date.Before(end1) || transaction.Date.Equal(end1)) {
			if transaction.Type == "income" {
				report.Period1Income += transaction.Amount
			} else {
				report.Period1Expenses += transaction.Amount
			}
		}

		// Period 2
		if (transaction.Date.After(start2) || transaction.Date.Equal(start2)) &&
			(transaction.Date.Before(end2) || transaction.Date.Equal(end2)) {
			if transaction.Type == "income" {
				report.Period2Income += transaction.Amount
			} else {
				report.Period2Expenses += transaction.Amount
			}
		}
	}

	// Calculate changes
	report.IncomeChange = report.Period2Income - report.Period1Income
	report.ExpenseChange = report.Period2Expenses - report.Period1Expenses

	// Calculate percentage changes
	if report.Period1Income > 0 {
		report.IncomePercent = (report.IncomeChange / report.Period1Income) * 100
	}
	if report.Period1Expenses > 0 {
		report.ExpensePercent = (report.ExpenseChange / report.Period1Expenses) * 100
	}

	return report
}

// Get year-over-year comparison
func getYearOverYearComparison(year1, year2 int) ComparisonReport {
	start1 := time.Date(year1, 1, 1, 0, 0, 0, 0, time.UTC)
	end1 := time.Date(year1, 12, 31, 23, 59, 59, 0, time.UTC)
	start2 := time.Date(year2, 1, 1, 0, 0, 0, 0, time.UTC)
	end2 := time.Date(year2, 12, 31, 23, 59, 59, 0, time.UTC)

	return comparePeriods(start1, end1, start2, end2)
}

// Get month-over-month comparison
func getMonthOverMonthComparison(year1 int, month1 time.Month, year2 int, month2 time.Month) ComparisonReport {
	start1 := time.Date(year1, month1, 1, 0, 0, 0, 0, time.UTC)
	end1 := time.Date(year1, month1+1, 0, 23, 59, 59, 0, time.UTC) // Last day of month
	start2 := time.Date(year2, month2, 1, 0, 0, 0, 0, time.UTC)
	end2 := time.Date(year2, month2+1, 0, 23, 59, 59, 0, time.UTC)

	return comparePeriods(start1, end1, start2, end2)
}

// Get quarterly report
func getQuarterlyReport(year int, quarter int) MonthlyReport {
	if quarter < 1 || quarter > 4 {
		return MonthlyReport{}
	}

	startMonth := time.Month((quarter-1)*3 + 1)
	endMonth := time.Month(quarter * 3)

	report := MonthlyReport{
		Year:  year,
		Month: startMonth, // Represents start of quarter
	}

	for _, transaction := range transactions {
		if transaction.Date.Year() != year {
			continue
		}
		if transaction.Date.Month() < startMonth || transaction.Date.Month() > endMonth {
			continue
		}

		if transaction.Type == "income" {
			report.Income += transaction.Amount
		} else {
			report.Expenses += transaction.Amount
		}
		report.TxCount++
	}

	report.Net = report.Income - report.Expenses
	return report
}

// Get average transaction amount by category
func getAverageAmountByCategory(category string, transactionType string) float64 {
	var total float64
	var count int

	for _, transac := range transactions {
		if transac.Type == transactionType || transactionType == "" {
			if transac.Category == category {
				total += transac.Amount
				count++
			}
		}
	}

	if count == 0 {
		return 0.0
	}
	return total / float64(count)
}

// Detect anomalies (unusually high spending)
func detectHighSpending(threshold float64) []Transaction {
	var total float64
	expenseCount := 0

	for _, transaction := range transactions {
		if transaction.Type == "expense" {
			total += transaction.Amount
			expenseCount++
		}
	}

	if expenseCount == 0 {
		return []Transaction{}
	}

	average := total / float64(expenseCount)
	anomalyThreshold := average * threshold

	// Find transactions above threshold
	var anomalies []Transaction
	for _, transaction := range transactions {
		if transaction.Type == "expense" && transaction.Amount > anomalyThreshold {
			anomalies = append(anomalies, transaction)
		}
	}
	return anomalies
}

// Get spending trend (increasing, decreasing, stable)
func getCategoryTrend(category string, months int) string {
	if len(transactions) == 0 {
		return "insufficient data"
	}

	// Get monthly spending for this category over last N months
	type MonthData struct {
		Date   time.Time
		Amount float64
	}
	monthlySpending := make(map[string]*MonthData)
	cutoffDate := time.Now().AddDate(0, -months, 0)

	for _, transaction := range transactions {
		if transaction.Type != "expense" {
			continue
		}
		if transaction.Category != category {
			continue
		}
		if transaction.Date.Before(cutoffDate) {
			continue
		}

		key := fmt.Sprintf("%d-%02d", transaction.Date.Year(), transaction.Date.Month())
		_, exists := monthlySpending[key]
		if !exists {
			monthlySpending[key] = &MonthData{Date: time.Date(transaction.Date.Year(), transaction.Date.Month(), 1, 0, 0, 0, 0, time.UTC)}
		}
		monthlySpending[key].Amount += transaction.Amount
	}

	if len(monthlySpending) < 2 {
		return "insufficient data"
	}

	// Convert to sorted slice
	var monthData []MonthData
	for _, data := range monthlySpending {
		monthData = append(monthData, *data)
	}

	// Sort by date
	for i := 0; i < len(monthData); i++ {
		for j := i + 1; j < len(monthData); j++ {
			if monthData[j].Date.Before(monthData[i].Date) {
				monthData[i], monthData[j] = monthData[j], monthData[i]
			}
		}
	}

	// Calculate trend (simple: compare first half vs second half)
	halfPoint := len(monthData) / 2
	var firstHalfAvg, secondHalfAvg float64

	for i := 0; i < halfPoint; i++ {
		firstHalfAvg += monthData[i].Amount
	}
	firstHalfAvg /= float64(halfPoint)

	for i := halfPoint; i < len(monthData); i++ {
		secondHalfAvg += monthData[i].Amount
	}
	secondHalfAvg /= float64(len(monthData) - halfPoint)

	// Determine trend
	difference := secondHalfAvg - firstHalfAvg
	percentChange := (difference / firstHalfAvg) * 100

	if percentChange > 10 {
		return "increasing"
	} else if percentChange < -10 {
		return "decreasing"
	} else {
		return "stable"
	}
}
