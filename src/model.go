package main

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
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
	enabled     string
}

func (s service) Title() string {
	statusIcon := "🔘"
	if s.enabled == "disabled" {
		statusIcon = "🔒"
	} else if strings.Contains(s.active, "running") {
		statusIcon = "🟢"
	} else if strings.Contains(s.active, "exited") {
		statusIcon = "🟡"
	} else if strings.Contains(s.active, "inactive") {
		statusIcon = "◯"
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
	db                 *sql.DB
	allServices        list.Model
	runningServices    list.Model
	focused            int // 0 = all services, 1 = running services
	loading            bool
	spinner            spinner.Model
	searchMode         bool
	searchInput        textinput.Model
	showHelp           bool
	showAbout          bool
	showMenu           bool
	showDescription    bool
	editingDescription bool
	descriptionInput   textarea.Model
	selectedService    service
	menuChoice         int
	message            string
	messageTimer       *time.Timer
}

type descriptionLoadedMsg struct {
	description string
}

func initialModel(db *sql.DB) model {
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

	ta := textarea.New()
	ta.Placeholder = "Enter a description for the service..."
	ta.SetWidth(50)
	ta.SetHeight(5)

	return model{
		db:                 db,
		allServices:        allList,
		runningServices:    runningList,
		focused:            0,
		loading:            true,
		spinner:            s,
		searchMode:         false,
		searchInput:        ti,
		showHelp:           false,
		showAbout:          false,
		showMenu:           false,
		showDescription:    false,
		editingDescription: false,
		descriptionInput:   ta,
		selectedService:    service{},
		menuChoice:         0,
		message:            "",
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		loadServices(m.db),
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.showDescription {
			if m.editingDescription {
				switch msg.String() {
				case "enter":
					m.editingDescription = false
					return m, updateServiceDescriptionCommand(m.db, m.selectedService.name, m.descriptionInput.Value())
				case "esc":
					m.editingDescription = false
					m.descriptionInput.Blur()
				}
				var cmd tea.Cmd
				m.descriptionInput, cmd = m.descriptionInput.Update(msg)
				cmds = append(cmds, cmd)
				return m, tea.Batch(cmds...)
			}

			switch msg.String() {
			case "e":
				m.editingDescription = true
				m.descriptionInput.Focus()
				return m, nil
			case "q", "ctrl+c", "esc", "U":
				m.showDescription = false
				m.editingDescription = false
				m.descriptionInput.Blur()
				m.descriptionInput.SetValue("")
			}
			return m, nil
		}

		if m.searchMode {
			switch msg.String() {
			case "enter":
				searchTerm := m.searchInput.Value()
				m.searchMode = false
				m.searchInput.Blur()
				m.searchInput.SetValue("")
				if searchTerm != "" {
					// Perform search and focus
					if m.focused == 0 {
						found := -1
						for i, item := range m.allServices.Items() {
							if s, ok := item.(service); ok {
								if strings.Contains(strings.ToLower(s.name), strings.ToLower(searchTerm)) || strings.Contains(strings.ToLower(s.description), strings.ToLower(searchTerm)) {
									found = i
									break
								}
							}
						}
						if found >= 0 {
							m.allServices.Select(found)
							m.message = fmt.Sprintf("Found and focused '%s'", searchTerm)
						} else {
							m.message = fmt.Sprintf("No services found matching '%s'", searchTerm)
						}
					} else {
						found := -1
						for i, item := range m.runningServices.Items() {
							if s, ok := item.(service); ok {
								if strings.Contains(strings.ToLower(s.name), strings.ToLower(searchTerm)) || strings.Contains(strings.ToLower(s.description), strings.ToLower(searchTerm)) {
									found = i
									break
								}
							}
						}
						if found >= 0 {
							m.runningServices.Select(found)
							m.message = fmt.Sprintf("Found and focused '%s'", searchTerm)
						} else {
							m.message = fmt.Sprintf("No services found matching '%s'", searchTerm)
						}
					}
					return m, nil
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
			if m.showMenu || m.showHelp || m.showAbout || m.searchMode || m.showDescription {
				m.showMenu = false
				m.showHelp = false
				m.showAbout = false
				m.searchMode = false
				m.showDescription = false
			} else {
				return m, tea.Quit
			}
		case "esc":
			if m.showMenu || m.showHelp || m.showAbout || m.searchMode || m.showDescription {
				m.showMenu = false
				m.showHelp = false
				m.showAbout = false
				m.searchMode = false
				m.showDescription = false
			}
		case "U":
			var item list.Item
			if m.focused == 0 {
				item = m.allServices.SelectedItem()
			} else {
				item = m.runningServices.SelectedItem()
			}
			if item != nil {
				if s, ok := item.(service); ok {
					m.selectedService = s
					m.showDescription = true
					return m, loadDescriptionCommand(m.db, s.name)
				}
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
						m.showMenu = false
						return m, tea.Batch(executeServiceCommand(m.selectedService.name, "start"), loadServices(m.db))
					case 2:
						m.showMenu = false
						return m, tea.Batch(executeServiceCommand(m.selectedService.name, "restart"), loadServices(m.db))
					case 3:
						m.showMenu = false
						return m, tea.Batch(executeServiceCommand(m.selectedService.name, "stop"), loadServices(m.db))
					case 4:
						m.showMenu = false
						return m, tea.Batch(executeServiceCommand(m.selectedService.name, "disable"), loadServices(m.db))
					case 5:
						m.showMenu = false
						return m, tea.Batch(executeServiceCommand(m.selectedService.name, "enable"), loadServices(m.db))
					}
				} else {
					// Running services menu
					switch m.menuChoice {
					case 1:
						m.showMenu = false
						return m, tea.Batch(executeServiceCommand(m.selectedService.name, "stop"), loadServices(m.db))
					case 2:
						m.showMenu = false
						return m, tea.Batch(executeServiceCommand(m.selectedService.name, "restart"), loadServices(m.db))
					case 3:
						m.showMenu = false
						return m, tea.Batch(executeServiceCommand(m.selectedService.name, "disable"), loadServices(m.db))
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
						m.showMenu = false
						return m, tea.Batch(executeServiceCommand(m.selectedService.name, "start"), loadServices(m.db))
					case 2:
						m.showMenu = false
						return m, tea.Batch(executeServiceCommand(m.selectedService.name, "restart"), loadServices(m.db))
					case 3:
						m.showMenu = false
						return m, tea.Batch(executeServiceCommand(m.selectedService.name, "stop"), loadServices(m.db))
					case 4:
						m.showMenu = false
						return m, tea.Batch(executeServiceCommand(m.selectedService.name, "disable"), loadServices(m.db))
					case 5:
						m.showMenu = false
						return m, tea.Batch(executeServiceCommand(m.selectedService.name, "enable"), loadServices(m.db))
					}
				} else {
					// Running services menu
					switch m.menuChoice {
					case 1:
						m.showMenu = false
						return m, tea.Batch(executeServiceCommand(m.selectedService.name, "stop"), loadServices(m.db))
					case 2:
						m.showMenu = false
						return m, tea.Batch(executeServiceCommand(m.selectedService.name, "restart"), loadServices(m.db))
					case 3:
						m.showMenu = false
						return m, tea.Batch(executeServiceCommand(m.selectedService.name, "disable"), loadServices(m.db))
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
		case "r":
			return m, loadServices(m.db)
		}

	case tea.WindowSizeMsg:
		h, v := lipgloss.NewStyle().Margin(1, 2).GetFrameSize()
		m.allServices.SetSize(msg.Width/2-h, msg.Height-v-10)
		m.runningServices.SetSize(msg.Width/2-h, msg.Height-v-10)

	case servicesLoadedMsg:
		m.loading = false
		m.allServices.SetItems(msg.allServices)
		m.runningServices.SetItems(msg.runningServices)

	case descriptionLoadedMsg:
		m.descriptionInput.SetValue(msg.description)
		var cmd tea.Cmd
		m.descriptionInput, cmd = m.descriptionInput.Update(msg)
		return m, cmd

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

func loadDescriptionCommand(db *sql.DB, serviceName string) tea.Cmd {
	return func() tea.Msg {
		description, err := getServiceDescription(db, serviceName)
		if err != nil {
			return messageMsg{text: fmt.Sprintf("Error loading description: %v", err)}
		}
		return descriptionLoadedMsg{description: description}
	}
}

func updateServiceDescriptionCommand(db *sql.DB, serviceName, description string) tea.Cmd {
	return func() tea.Msg {
		err := updateServiceDescription(db, serviceName, description)
		if err != nil {
			return messageMsg{text: fmt.Sprintf("Error updating description: %v", err)}
		}
		return messageMsg{text: "✅ Successfully updated description"}
	}
} 
