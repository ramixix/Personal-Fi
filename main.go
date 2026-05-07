package main

import (
	"financial_tracker/internal/cli"
	"financial_tracker/internal/gui"
	"financial_tracker/internal/storage"
	"fmt"
	"os"
)

const dbFile = "financial_tracker.db"

func main() {
	// Initialize storage (SQLite database)
	err := storage.InitStorage(dbFile)
	if err != nil {
		fmt.Printf("Fatal: Failed to initialize database: %v\n", err)
		os.Exit(1)
	}

	// Ensure database is closed on exit
	defer func() {
		if err := storage.CloseStorage(); err != nil {
			fmt.Printf("Warning: Failed to close database: %v\n", err)
		}
	}()

	// Check if GUI flag is provided
	if len(os.Args) > 1 && os.Args[1] == "--gui" {
		// Launch GUI
		gui.Run()
		return
	}

	// CLI Mode
	cli.Run()
}
