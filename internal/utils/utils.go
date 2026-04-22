package utils

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Parse date string in format YYYY-MM-DD
func ParseDate(date_input string) (time.Time, error) {
	layout := "2006-01-02"
	date, err := time.Parse(layout, date_input)
	return date, err
}

func GetValidAmount(reader *bufio.Reader) float64 {
	for {
		fmt.Print("Amount: $")
		amount_input, _ := reader.ReadString('\n')
		transcation_amount, err := strconv.ParseFloat(strings.TrimSpace(amount_input), 64)

		if err != nil {
			fmt.Println("Invalid amount, please enter a valid number (float/integer)")
			continue
		}

		if transcation_amount <= 0 {
			fmt.Println("Amount must be greater than 0")
			continue
		}

		return transcation_amount
	}
}

func GetNonEmptyString(reader *bufio.Reader, prompt string) string {
	for {
		fmt.Print(prompt)
		input, _ := reader.ReadString('\n')
		text := strings.TrimSpace(input)

		if text != "" {
			return text
		}
		fmt.Println("[Warning] This Field Can Not Be Empty!")
	}
}

// Get yes/no confirmation
func GetConfirmation(reader *bufio.Reader, prompt string) bool {
	fmt.Print(prompt)
	input, _ := reader.ReadString('\n')
	response := strings.ToLower(strings.TrimSpace(input))
	return response == "yes" || response == "y"
}

// Get integer input
func GetIntInput(reader *bufio.Reader, prompt string) (int, error) {
	fmt.Print(prompt)
	input, _ := reader.ReadString('\n')
	return strconv.Atoi(strings.TrimSpace(input))
}

// prompts for and validates currency code
func GetValidCurrency(reader bufio.Reader, defaultCurrency string) string {
	fmt.Printf("Default Currency: %s", defaultCurrency)
	input, _ := reader.ReadString('\n')
	currency := strings.ToUpper(strings.TrimSpace(input))
	if currency == "" {
		return defaultCurrency
	}
	if len(currency) != 3 {
		fmt.Println("Invalid currency code. Using default.")
		return defaultCurrency
	}
	return currency
}

// Common currency codes
var CommonCurrencies = []string{
	"USD", "EUR", "GBP", "JPY", "CNY", "AUD", "CAD", "CHF", "HKD", "SGD",
	"SEK", "KRW", "NOK", "NZD", "INR", "MXN", "ZAR", "BRL", "TRY", "RUB",
}

// checks if currency code is in common currencies list (CommonCurrencies)
func IsValidCurrency(currencyCode string) bool {
	for _, code := range CommonCurrencies {
		if code == strings.ToUpper(currencyCode) {
			return true
		}
	}
	return false
}

func GetCurrencySymbol(currencyCode string) string {
	symbols := map[string]string{
		"USD": "$",
		"EUR": "€",
		"GBP": "£",
		"JPY": "¥",
		"CNY": "¥",
		"INR": "₹",
		"RUB": "₽",
		"KRW": "₩",
		"TRY": "₺",
		"BRL": "R$",
		"AUD": "A$",
		"CAD": "C$",
		"CHF": "Fr",
		"SEK": "kr",
		"NOK": "kr",
	}

	if symbol, ok := symbols[strings.ToUpper(currencyCode)]; ok {
		return symbol
	}

	return currencyCode + " "
}

// FormatCurrency formats amount with currency symbol
func FormatCurrency(amount float64, currencyCode string) string {
	symbol := GetCurrencySymbol(currencyCode)
	return fmt.Sprintf("%s %.2f", symbol, amount)
}
