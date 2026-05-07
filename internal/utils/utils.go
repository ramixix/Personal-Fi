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
func GetValidCurrency(reader *bufio.Reader, defaultCurrency string) string {
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

	valid := IsValidCurrency(currency)
	if valid != true {
		fmt.Println("Invalid currency code. Using default.")
		return defaultCurrency
	}

	return currency
}

// All currency codes
var currencyCodes = []string{
	"AED", "AFN", "ALL", "AMD", "ANG", "AOA", "ARS", "AUD", "AWG", "AZN",
	"BAM", "BBD", "BDT", "BGN", "BHD", "BMD", "BND", "BOB", "BRL", "BSD",
	"BTN", "BWP", "BYN", "BZD", "CAD", "CDF", "CHF", "CLP", "CNY", "COP",
	"CRC", "CUC", "CUP", "CVE", "CZK", "DJF", "DKK", "DOP", "DZD", "EGP",
	"ERN", "ETB", "EUR", "FJD", "FKP", "GBP", "GEL", "GHS", "GIP", "GMD",
	"GNF", "GTQ", "GYD", "HKD", "HNL", "HRK", "HTG", "HUF", "IDR", "ILS",
	"INR", "IQD", "IRR", "ISK", "JMD", "JOD", "JPY", "KES", "KGS", "KHR",
	"KMF", "KRW", "KWD", "KYD", "KZT", "LAK", "LBP", "LKR", "LRD", "LSL",
	"LYD", "MAD", "MDL", "MGA", "MKD", "MMK", "MNT", "MOP", "MRU", "MUR",
	"MVR", "MWK", "MXN", "MYR", "MZN", "NAD", "NGN", "NIO", "NOK", "NPR",
	"NZD", "OMR", "PAB", "PEN", "PGK", "PHP", "PKR", "PLN", "PYG", "QAR",
	"RON", "RSD", "RUB", "RWF", "SAR", "SBD", "SCR", "SDG", "SEK", "SGD",
	"SHP", "SLL", "SOS", "SRD", "STN", "SVC", "SYP", "SZL", "THB", "TJS",
	"TMT", "TND", "TOP", "TRY", "TTD", "TWD", "TZS", "UAH", "UGX", "USD",
	"UYU", "UZS", "VES", "VND", "VUV", "WST", "XAF", "XCD", "XOF", "XPF",
	"YER", "ZAR", "ZMW", "ZWL",
}

// checks if currency code is in common currencies list (CommonCurrencies)
func IsValidCurrency(currencyCode string) bool {
	for _, code := range currencyCodes {
		if code == strings.ToUpper(currencyCode) {
			return true
		}
	}
	return false
}

func GetCurrencySymbol(currencyCode string) string {
	symbols := map[string]string{
		"AED": "د.إ",
		"AFN": "؋",
		"ALL": "L",
		"AMD": "֏",
		"ANG": "ƒ",
		"AOA": "Kz",
		"ARS": "$",
		"AUD": "A$",
		"AWG": "ƒ",
		"AZN": "₼",
		"BAM": "KM",
		"BBD": "B$",
		"BDT": "৳",
		"BGN": "лв",
		"BHD": ".د.ب",
		"BMD": "$",
		"BND": "B$",
		"BOB": "Bs.",
		"BRL": "R$",
		"BSD": "$",
		"BTN": "Nu.",
		"BWP": "P",
		"BYN": "Br",
		"BZD": "BZ$",
		"CAD": "C$",
		"CDF": "FC",
		"CHF": "Fr",
		"CLP": "$",
		"CNY": "¥",
		"COP": "$",
		"CRC": "₡",
		"CUC": "$",
		"CUP": "$",
		"CVE": "$",
		"CZK": "Kč",
		"DJF": "Fdj",
		"DKK": "kr",
		"DOP": "RD$",
		"DZD": "د.ج",
		"EGP": "£",
		"ERN": "Nfk",
		"ETB": "Br",
		"EUR": "€",
		"FJD": "FJ$",
		"FKP": "£",
		"GBP": "£",
		"GEL": "₾",
		"GHS": "₵",
		"GIP": "£",
		"GMD": "D",
		"GNF": "FG",
		"GTQ": "Q",
		"GYD": "$",
		"HKD": "HK$",
		"HNL": "L",
		"HRK": "kn",
		"HTG": "G",
		"HUF": "Ft",
		"IDR": "Rp",
		"ILS": "₪",
		"INR": "₹",
		"IQD": "ع.د",
		"IRR": "﷼",
		"ISK": "kr",
		"JMD": "J$",
		"JOD": "د.ا",
		"JPY": "¥",
		"KES": "KSh",
		"KGS": "лв",
		"KHR": "៛",
		"KMF": "CF",
		"KRW": "₩",
		"KWD": "د.ك",
		"KYD": "$",
		"KZT": "₸",
		"LAK": "₭",
		"LBP": "ل.ل",
		"LKR": "₨",
		"LRD": "$",
		"LSL": "L",
		"LYD": "ل.د",
		"MAD": "د.م.",
		"MDL": "L",
		"MGA": "Ar",
		"MKD": "ден",
		"MMK": "Ks",
		"MNT": "₮",
		"MOP": "MOP$",
		"MRU": "UM",
		"MUR": "₨",
		"MVR": "Rf",
		"MWK": "MK",
		"MXN": "$",
		"MYR": "RM",
		"MZN": "MT",
		"NAD": "$",
		"NGN": "₦",
		"NIO": "C$",
		"NOK": "kr",
		"NPR": "₨",
		"NZD": "NZ$",
		"OMR": "﷼",
		"PAB": "B/.",
		"PEN": "S/.",
		"PGK": "K",
		"PHP": "₱",
		"PKR": "₨",
		"PLN": "zł",
		"PYG": "₲",
		"QAR": "﷼",
		"RON": "lei",
		"RSD": "дин.",
		"RUB": "₽",
		"RWF": "FRw",
		"SAR": "﷼",
		"SBD": "$",
		"SCR": "₨",
		"SDG": "£",
		"SEK": "kr",
		"SGD": "S$",
		"SHP": "£",
		"SLL": "Le",
		"SOS": "S",
		"SRD": "$",
		"STN": "Db",
		"SVC": "$",
		"SYP": "£",
		"SZL": "L",
		"THB": "฿",
		"TJS": "ЅМ",
		"TMT": "m",
		"TND": "د.ت",
		"TOP": "T$",
		"TRY": "₺",
		"TTD": "TT$",
		"TWD": "NT$",
		"TZS": "TSh",
		"UAH": "₴",
		"UGX": "USh",
		"USD": "$",
		"UYU": "$U",
		"UZS": "so'm",
		"VES": "Bs.",
		"VND": "₫",
		"VUV": "VT",
		"WST": "T",
		"XAF": "FCFA",
		"XCD": "EC$",
		"XOF": "CFA",
		"XPF": "CFPF",
		"YER": "﷼",
		"ZAR": "R",
		"ZMW": "ZK",
		"ZWL": "$",
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
