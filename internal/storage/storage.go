package storage

import (
	"encoding/json"
	"financial_tracker/internal/models"
	"fmt"
	"os"
)

// file to save necessary infromation
const dataFile = "financial_data.json"

// Global storage - in memory
var Transactions []models.Transaction
var Accounts []models.Account
var AccountTransactions []models.AccountTransaction

var NextTransactionID = 1
var NextAccountID = 1
var NextAccountTransactionID = 1

// Goals variables
var Goals []models.Goal
var GoalContributions []models.GoalContribution

var NextGoalID = 1
var NextGoalContributionID = 1

// Save all data to JSON file
func SaveData() error {
	data := models.AppData{
		Transactions:             Transactions,
		Accounts:                 Accounts,
		AccountTransactions:      AccountTransactions,
		Goals:                    Goals,
		GoalContributions:        GoalContributions,
		NextTransactionID:        NextTransactionID,
		NextAccountID:            NextAccountID,
		NextAccountTransactionID: NextAccountTransactionID,
		NextGoalID:               NextGoalID,
		NextGoalContributionID:   NextGoalContributionID,
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
func LoadData() error {
	if _, err := os.Stat(dataFile); os.IsNotExist(err) {
		fmt.Println("[Info] DataFile to Load Information From Does Not Exists!")
		return nil
	}

	jsonData, err := os.ReadFile(dataFile)
	if err != nil {
		return err
	}

	var data models.AppData
	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		return err
	}

	// Load data into global variables
	Transactions = data.Transactions
	Accounts = data.Accounts
	AccountTransactions = data.AccountTransactions
	Goals = data.Goals
	GoalContributions = data.GoalContributions
	NextTransactionID = data.NextTransactionID
	NextAccountID = data.NextAccountID
	NextAccountTransactionID = data.NextAccountTransactionID
	NextGoalID = data.NextGoalID
	NextGoalContributionID = data.NextGoalContributionID

	return nil
}
