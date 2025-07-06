package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1).
			Bold(true)

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
)

func (m model) View() string {
	if m.loading {
		return m.loadingView()
	}

	if m.searchMode {
		return m.searchView()
	}

	if m.showHelp {
		return m.helpView()
	}

	if m.showAbout {
		return m.aboutView()
	}

	if m.showMenu {
		return m.menuView()
	}

	return m.mainView()
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
  ?                  Toggle this help
  P                  Show about/coffee info
  q / Ctrl+C         Quit

Service Actions:
  All Services:      1=Start, 2=Restart, 3=Stop, 4=Disable, 5=Enable
  Running Services:  1=Stop, 2=Restart, 3=Disable

Press any key to return...
`
	return modalStyle.Render(help)
}

func (m model) aboutView() string {
	about := `
╔══════════════════════════════════════════════════════════════╗
║                    ☕ Buy Me a Coffee ☕                      ║
╠══════════════════════════════════════════════════════════════╣
║                                                              ║
║  🎉 Thanks for using LazySys!                               ║
║                                                              ║
║  If you find this tool helpful, consider buying me a coffee ║
║  to support further development!                             ║
║                                                              ║
║  💳 Bitcoin: bc1qxy2kgdygjrsqtzq2n0yrf2493p83kkfjhx0wlh   ║
║  💳 Ethereum: 0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6  ║
║  ☕ Ko-fi: https://ko-fi.com/lazysys                        ║
║                                                              ║
║  Made with ❤️  using BubbleTea                              ║
║                                                              ║
╚══════════════════════════════════════════════════════════════╝

Press any key to return...
`
	return modalStyle.Render(about)
}

func (m model) mainView() string {
	var s strings.Builder

	// Title
	s.WriteString(titleStyle.Render("🔧 LazySys Service Manager"))
	s.WriteString("\n\n")

	// Service counts
	allCount := len(m.allServices.Items())
	runningCount := len(m.runningServices.Items())
	
	s.WriteString(fmt.Sprintf("📊 Total Services: %d | 🟢 Running: %d\n\n", allCount, runningCount))

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

	s.WriteString(lists)
	s.WriteString("\n\n")

	// Help bar
	helpText := "H/L: Navigate | j/k: Scroll | Enter: Action | s: Search | ?: Help | P: About | q: Quit"
	s.WriteString(helpStyle.Render(helpText))

	// Message
	if m.message != "" {
		s.WriteString("\n")
		s.WriteString(messageStyle.Render(m.message))
	}

	return s.String()
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

	var menuContent strings.Builder
	menuContent.WriteString(title)
	menuContent.WriteString("\n\n")

	for i, item := range menuItems {
		if i == m.menuChoice {
			menuContent.WriteString("▶ " + item + "\n")
		} else {
			menuContent.WriteString("  " + item + "\n")
		}
	}

	menuContent.WriteString("\nEnter: Execute | Esc: Cancel")

	return modalStyle.Render(menuContent.String())
} 