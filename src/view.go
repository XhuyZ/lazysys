package main

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
)

var ( titleStyle = lipgloss.NewStyle(). Foreground(lipgloss.Color("#FAFAFA")). Background(lipgloss.Color("#7D56F4")). Padding(0, 1). Bold(true)

	focusedStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7D56F4")).
			Padding(1, 0)

	unfocusedStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#626262")).
			Padding(1, 0)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262")).
			Italic(true)

	messageStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#04B575")).
			Bold(true)

	searchStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#04B575")).
			Padding(0, 1)

	modalStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7D56F4")).
			Padding(1, 2).
			Background(lipgloss.Color("#1A1A1A"))

	aboutStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFD700")).
			Bold(true)

	dimStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#444444"))
)

func (m model) View() string {
	if m.loading {
		return m.loadingView()
	}

	main := m.mainView()

	// Get terminal size for centering
	w, h := 80, 25
	if m.allServices.Width() > 0 && m.allServices.Height() > 0 {
		w = m.allServices.Width()*2 + 8 // 2 panes + padding
		h = m.allServices.Height() + 10 // add for title/help
	}

	if m.showHelp {
		return dimStyle.Render(main) + "\n" + m.floatingModal(m.helpView(), w, h)
	}
	if m.showAbout {
		return dimStyle.Render(main) + "\n" + m.floatingModal(m.aboutView(), w, h)
	}
	if m.searchMode {
		return dimStyle.Render(main) + "\n" + m.floatingModal(m.searchView(), w, h)
	}
	if m.showMenu {
		return dimStyle.Render(main) + "\n" + m.floatingModal(m.menuView(), w, h)
	}
	if m.showDescription {
		return dimStyle.Render(main) + "\n" + m.floatingModal(m.descriptionView(), w, h)
	}

	return main
}

func (m model) floatingModal(content string, w, h int) string {
	return lipgloss.Place(
		w, h,
		lipgloss.Center, lipgloss.Center,
		content,
	)
}

func (m model) loadingView() string {
	return lipgloss.Place(
		80, 25,
		lipgloss.Center, lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Center,
			titleStyle.Render("🔧 LazySys Service Manager"),
			"",
			m.spinner.View()+" Loading services...",
		),
	)
}

func (m model) searchView() string {
	windowName := "All Services"
	if m.focused == 1 {
		windowName = "Running Services"
	}

	searchBox := searchStyle.Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			fmt.Sprintf("🔍 Search %s", windowName),
			m.searchInput.View(),
			"\nEnter: Search | Esc: Cancel",
		),
	)

	return lipgloss.Place(
		80, 25,
		lipgloss.Center, lipgloss.Center,
		searchBox,
	)
}

func (m model) helpView() string {
	help := `
🔧 LazySys Service Manager - Help

Navigation:
  Shift+H / Shift+L  Navigate between windows
  j / k              Navigate up/down in lists
  Enter              Select service for action
  s                  Search services
  r                  Reload UI/services
  U                  View/Edit service description
  ?                  Toggle this help
  P                  Show about/coffee info
  q / Esc / Ctrl+C   Quit/close window

Service Actions:
  All Services:      1=Start, 2=Restart, 3=Stop, 4=Disable, 5=Enable
  Running Services:  1=Stop, 2=Restart, 3=Disable
`
	return modalStyle.Render(help)
}

func (m model) aboutView() string {
	about := `
╔══════════════════════════════════════════════════════════════╗
║                    ☕ Buy Me a Coffee ☕                     ║
╠══════════════════════════════════════════════════════════════╣
║                                                              ║
║  🎉 Thanks for using LazySys!                                
║                                                              ║
║  If you find this tool helpful, consider buying me a coffee  ║
║  to support further development!                             ║
║                                                              ║
║                                                              ║
║  💳 Nah its free :))                                         ║
║  ☕ How about giving this repo a star ?                      ║
║                                                              ║
║  Made with ❤️  using BubbleTea                               ║
║                                                              ║
╚══════════════════════════════════════════════════════════════╝
`
	return modalStyle.Render(about)
}

func (m model) mainView() string {
	var s string

	// Title
	s += titleStyle.Render("🔧 LazySys Service Manager") + "\n\n"

	// Service counts
	allCount := len(m.allServices.Items())
	runningCount := len(m.runningServices.Items())
	s += fmt.Sprintf("📊 Total Services: %d | 🟢 Running: %d\n\n", allCount, runningCount)

	// Lists
	allServicesView := m.allServices.View()
	runningServicesView := m.runningServices.View()

	// Apply focus styling
	if m.focused == 0 {
		allServicesView = focusedStyle.Render(allServicesView)
		runningServicesView = unfocusedStyle.Render(runningServicesView)
	} else {
		allServicesView = unfocusedStyle.Render(allServicesView)
		runningServicesView = focusedStyle.Render(runningServicesView)
	}

	// Split layout
	lists := lipgloss.JoinHorizontal(
		lipgloss.Top,
		allServicesView,
		"  ",
		runningServicesView,
	)

	s += lists + "\n\n"

	// Help bar
	helpText := "H/L: Navigate | j/k: Scroll | Enter: Action | s: Search | r: Reload || U: Show services info | ?: Help | P: About | q: Quit"
	s += helpStyle.Render(helpText)

	// Message
	if m.message != "" {
		s += "\n" + messageStyle.Render(m.message)
	}

	return s
}

func (m model) menuView() string {
	var menuItems []string
	var title string

	if m.focused == 0 {
		title = fmt.Sprintf("🔧 Service: %s", m.selectedService.name)
		menuItems = []string{
			"1. Start",
			"2. Restart",
			"3. Stop",
			"4. Disable",
			"5. Enable",
		}
	} else {
		title = fmt.Sprintf("🔧 Service: %s", m.selectedService.name)
		menuItems = []string{
			"1. Stop",
			"2. Restart",
			"3. Disable",
		}
	}

	var menuContent string
	menuContent += title + "\n\n"
	for i, item := range menuItems {
		if i == m.menuChoice {
			menuContent += "▶ " + item + "\n"
		} else {
			menuContent += "  " + item + "\n"
		}
	}
	menuContent += "\nEnter: Execute | Esc/q: Cancel"
	return modalStyle.Render(menuContent)
}

func (m model) descriptionView() string {
	var content string
	title := fmt.Sprintf("📖 Description: %s", m.selectedService.name)
	content += title + "\n\n"

	if m.editingDescription {
		content += m.descriptionInput.View()
		content += "\n\nEnter: Save | Esc: Cancel"
	} else {
		content += m.descriptionInput.Value()
		content += "\n\ne: Edit | q/Esc/U: Close"
	}

	return modalStyle.Render(content)
}
