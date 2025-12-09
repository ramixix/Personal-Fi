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
