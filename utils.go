package main

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Parse date string in format YYYY-MM-DD
func parseDate(date_input string) (time.Time, error) {
	layout := "2006-01-02"
	date, err := time.Parse(layout, date_input)
	return date, err
}

func getValidTransactionType(reader *bufio.Reader) string {
	for {
		fmt.Print("Type (Income/Expense): ")
		input, _ := reader.ReadString('\n')
		transaction_type := strings.ToLower(strings.TrimSpace(input))

		if transaction_type == "income" || transaction_type == "expense" {
			return transaction_type
		}
		fmt.Println("[Warning] Invalid type! Pleae Enter 'income' or 'expense'")
	}
}

func getValidAmount(reader *bufio.Reader) float64 {
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

func getNonEmptyString(reader *bufio.Reader, prompt string) string {
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
func getConfirmation(reader *bufio.Reader, prompt string) bool {
	fmt.Print(prompt)
	input, _ := reader.ReadString('\n')
	response := strings.ToLower(strings.TrimSpace(input))
	return response == "yes" || response == "y"
}

// Get integer input
func getIntInput(reader *bufio.Reader, prompt string) (int, error) {
	fmt.Print(prompt)
	input, _ := reader.ReadString('\n')
	return strconv.Atoi(strings.TrimSpace(input))
}

func getTransactionNumberToShow(reader *bufio.Reader, defaultValue int) int {
	transactionsToShow := defaultValue
InputLoop:
	for {
		input, _ := reader.ReadString('\n')
		input = strings.ToLower(strings.TrimSpace(input))

		switch input {
		case "":
			fmt.Printf("\nDisplaying %d Recent transactions:\n", transactionsToShow)
			break InputLoop
		case "all":
			fmt.Println("\nDisplaying All transactions:")
			transactionsToShow = len(transactions)
			break InputLoop
		default:
			number, err := strconv.Atoi(input)
			if err != nil || number <= 0 {
				fmt.Println("[Warning] Not a Valid Number, Try Again.")
				continue
			}
			fmt.Printf("\nDisplaying Last %d Transactions:\n", number)
			transactionsToShow = number
			break InputLoop
		}
	}
	return transactionsToShow
}
