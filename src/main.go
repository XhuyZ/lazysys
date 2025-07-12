package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Check if running with sudo
	if os.Geteuid() != 0 {
		fmt.Println("❌ This application requires sudo privileges to manage systemd services.")
		fmt.Println("Please run: sudo lazysys")
		os.Exit(1)
	}

	db, err := initDB()
	if err != nil {
		fmt.Printf("Error initializing database: %v", err)
		os.Exit(1)
	}
	defer db.Close()

	p := tea.NewProgram(initialModel(db), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v", err)
		os.Exit(1)
	}
} 