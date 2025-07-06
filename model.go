package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type service struct {
	name        string
	description string
	status      string
	loaded      string
	active      string
}

func (s service) Title() string {
	statusIcon := "🔴"
	if strings.Contains(s.active, "active") {
		statusIcon = "🟢"
	} else if strings.Contains(s.active, "inactive") {
		statusIcon = "⚪"
	}
	return fmt.Sprintf("%s %s", statusIcon, s.name)
}

func (s service) Description() string {
	return s.description
}

func (s service) FilterValue() string {
	return s.name
}

type model struct {
	allServices     list.Model
	runningServices list.Model
	focused         int // 0 = all services, 1 = running services
	loading         bool
	spinner         spinner.Model
	searchMode      bool
	searchInput     textinput.Model
	showHelp        bool
	showAbout       bool
	showMenu        bool
	selectedService service
	menuChoice      int
	message         string
	messageTimer    *time.Timer
}

func initialModel() model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	ti := textinput.New()
	ti.Placeholder = "Search services..."
	ti.CharLimit = 50
	ti.Width = 40

	allList := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	allList.Title = "📋 All Services"
	allList.SetShowHelp(false)

	runningList := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	runningList.Title = "🟢 Running Services"
	runningList.SetShowHelp(false)

	return model{
		allServices:     allList,
		runningServices: runningList,
		focused:         0,
		loading:         true,
		spinner:         s,
		searchMode:      false,
		searchInput:     ti,
		showHelp:        false,
		showAbout:       false,
		showMenu:        false,
		selectedService: service{},
		menuChoice:      0,
		message:         "",
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		loadServices(),
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.searchMode {
			switch msg.String() {
			case "enter":
				searchTerm := m.searchInput.Value()
				m.searchMode = false
				m.searchInput.Blur()
				m.searchInput.SetValue("")
				if searchTerm != "" {
					return m, performSearch(searchTerm, m.focused)
				}
				return m, tea.Batch(cmds...)
			case "esc":
				m.searchMode = false
				m.searchInput.Blur()
				m.searchInput.SetValue("")
				return m, tea.Batch(cmds...)
			}
			var cmd tea.Cmd
			m.searchInput, cmd = m.searchInput.Update(msg)
			cmds = append(cmds, cmd)
			return m, tea.Batch(cmds...)
		}

		switch msg.String() {
		case "q", "ctrl+c":
			if m.showMenu {
				m.showMenu = false
			} else {
				return m, tea.Quit
			}
		case "esc":
			if m.showMenu {
				m.showMenu = false
			}
		case "H":
			m.focused = 0
		case "L":
			m.focused = 1
		case "j":
			if m.showMenu {
				if m.focused == 0 && m.menuChoice < 5 {
					m.menuChoice++
				} else if m.focused == 1 && m.menuChoice < 3 {
					m.menuChoice++
				}
			} else {
				if m.focused == 0 {
					var cmd tea.Cmd
					m.allServices, cmd = m.allServices.Update(msg)
					cmds = append(cmds, cmd)
				} else {
					var cmd tea.Cmd
					m.runningServices, cmd = m.runningServices.Update(msg)
					cmds = append(cmds, cmd)
				}
			}
		case "k":
			if m.showMenu {
				if m.menuChoice > 0 {
					m.menuChoice--
				}
			} else {
				if m.focused == 0 {
					var cmd tea.Cmd
					m.allServices, cmd = m.allServices.Update(msg)
					cmds = append(cmds, cmd)
				} else {
					var cmd tea.Cmd
					m.runningServices, cmd = m.runningServices.Update(msg)
					cmds = append(cmds, cmd)
				}
			}
		case "s":
			m.searchMode = true
			m.searchInput.Focus()
			m.searchInput.SetValue("")
		case "?":
			m.showHelp = !m.showHelp
		case "P":
			m.showAbout = !m.showAbout
		case "1", "2", "3", "4", "5":
			if m.showMenu {
				m.menuChoice = int(msg.String()[0] - '0')
				// Execute action based on focused window and choice
				if m.focused == 0 {
					// All services menu
					switch m.menuChoice {
					case 1:
						return m, executeServiceCommand(m.selectedService.name, "start")
					case 2:
						return m, executeServiceCommand(m.selectedService.name, "restart")
					case 3:
						return m, executeServiceCommand(m.selectedService.name, "stop")
					case 4:
						return m, executeServiceCommand(m.selectedService.name, "disable")
					case 5:
						return m, executeServiceCommand(m.selectedService.name, "enable")
					}
				} else {
					// Running services menu
					switch m.menuChoice {
					case 1:
						return m, executeServiceCommand(m.selectedService.name, "stop")
					case 2:
						return m, executeServiceCommand(m.selectedService.name, "restart")
					case 3:
						return m, executeServiceCommand(m.selectedService.name, "disable")
					}
				}
				m.showMenu = false
			}
		case "enter":
			if m.showMenu {
				// Execute action based on focused window and choice
				if m.focused == 0 {
					// All services menu
					switch m.menuChoice {
					case 1:
						return m, executeServiceCommand(m.selectedService.name, "start")
					case 2:
						return m, executeServiceCommand(m.selectedService.name, "restart")
					case 3:
						return m, executeServiceCommand(m.selectedService.name, "stop")
					case 4:
						return m, executeServiceCommand(m.selectedService.name, "disable")
					case 5:
						return m, executeServiceCommand(m.selectedService.name, "enable")
					}
				} else {
					// Running services menu
					switch m.menuChoice {
					case 1:
						return m, executeServiceCommand(m.selectedService.name, "stop")
					case 2:
						return m, executeServiceCommand(m.selectedService.name, "restart")
					case 3:
						return m, executeServiceCommand(m.selectedService.name, "disable")
					}
				}
				m.showMenu = false
			} else {
				if m.focused == 0 {
					if item := m.allServices.SelectedItem(); item != nil {
						if s, ok := item.(service); ok {
							m.selectedService = s
							m.showMenu = true
							m.menuChoice = 0
						}
					}
				} else {
					if item := m.runningServices.SelectedItem(); item != nil {
						if s, ok := item.(service); ok {
							m.selectedService = s
							m.showMenu = true
							m.menuChoice = 0
						}
					}
				}
			}
		}

	case tea.WindowSizeMsg:
		h, v := lipgloss.NewStyle().Margin(1, 2).GetFrameSize()
		m.allServices.SetSize(msg.Width/2-h, msg.Height-v-10)
		m.runningServices.SetSize(msg.Width/2-h, msg.Height-v-10)

	case servicesLoadedMsg:
		m.loading = false
		m.allServices.SetItems(msg.allServices)
		m.runningServices.SetItems(msg.runningServices)

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)

	case messageMsg:
		m.message = msg.text
		if m.messageTimer != nil {
			m.messageTimer.Stop()
		}
		m.messageTimer = time.AfterFunc(3*time.Second, func() {
			// This will be handled in the view
		})

	default:
		if m.focused == 0 {
			var cmd tea.Cmd
			m.allServices, cmd = m.allServices.Update(msg)
			cmds = append(cmds, cmd)
		} else {
			var cmd tea.Cmd
			m.runningServices, cmd = m.runningServices.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
} 