package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// file to save necessary infromation
const dataFile = "financial_data.json"

// Global storage - in memory
var transactions []Transaction
var accounts []Account
var accountTransactions []AccountTransaction

var nextTransactionID = 1
var nextAccountID = 1
var nextAccountTransactionID = 1

// Goals variables
var goals []Goal
var goalContributions []GoalContribution

var nextGoalID = 1
var nextGoalContributionID = 1

// Save all data to JSON file
func saveData() error {
	data := AppData{
		Transactions:             transactions,
		Accounts:                 accounts,
		AccountTransactions:      accountTransactions,
		Goals:                    goals,
		GoalContributions:        goalContributions,
		NextTransactionID:        nextTransactionID,
		NextAccountID:            nextAccountID,
		NextAccountTransactionID: nextAccountTransactionID,
		NextGoalID:               nextGoalID,
		NextGoalContributionID:   nextGoalContributionID,
	}

	dataJsonFormat, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		fmt.Println("[Warning] Could not convert struct to json format in saveData function!")
		return err
	}

	err = os.WriteFile(dataFile, dataJsonFormat, 0664)
	if err != nil {
		fmt.Println("[Warning] Could not Write to file to save data in saveData function!")
		return err
	}
	return nil
}

// Load all data from JSON file
func loadData() error {
	if _, err := os.Stat(dataFile); os.IsNotExist(err) {
		fmt.Println("[Info] DataFile to Load Information From Does Not Exists!")
		return nil
	}

	jsonData, err := os.ReadFile(dataFile)
	if err != nil {
		return err
	}

	var data AppData
	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		return err
	}

	// Load data into global variables
	transactions = data.Transactions
	accounts = data.Accounts
	accountTransactions = data.AccountTransactions
	nextTransactionID = data.NextTransactionID
	nextAccountID = data.NextAccountID
	nextAccountTransactionID = data.NextAccountTransactionID
	goals = data.Goals
	goalContributions = data.GoalContributions
	nextGoalID = data.NextGoalID
	nextGoalContributionID = data.NextGoalContributionID

	return nil
}
