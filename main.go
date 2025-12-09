package main

import (
	"financial_tracker/internal/cli"
	"fmt"
	"os"
)

func main() {
	// Check if GUI flag is provided
	if len(os.Args) > 1 && os.Args[1] == "--gui" {
		// TODO: Launch GUI
		fmt.Println("GUI mode will be implemented next!")
		fmt.Println("For now, use CLI mode (run without --gui flag)")
		return
	}

	// CLI Mode
	cli.Run()
}
