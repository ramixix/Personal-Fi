package main

import (
	"financial_tracker/internal/cli"
	"financial_tracker/internal/gui"
	"financial_tracker/internal/storage"
	"fmt"
	"os"
)

func main() {
	// Load data first
	err := storage.LoadData()
	if err != nil {
		fmt.Printf("Warning: Could not load data: %v\n", err)
	}

	// Check if GUI flag is provided
	if len(os.Args) > 1 && os.Args[1] == "--gui" {
		// Launch GUI
		gui.Run()
		return
	}

	// CLI Mode
	cli.Run()
}
