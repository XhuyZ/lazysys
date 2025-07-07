package main

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type servicesLoadedMsg struct {
	allServices     []list.Item
	runningServices []list.Item
}

type messageMsg struct {
	text string
}

func loadServices() tea.Cmd {
	return func() tea.Msg {
		allServices, err := getAllServices()
		if err != nil {
			return messageMsg{text: fmt.Sprintf("Error loading all services: %v", err)}
		}

		runningServices, err := getRunningServices()
		if err != nil {
			return messageMsg{text: fmt.Sprintf("Error loading running services: %v", err)}
		}

		return servicesLoadedMsg{
			allServices:     allServices,
			runningServices: runningServices,
		}
	}
}

func getAllServices() ([]list.Item, error) {
	// Get all loaded services
	cmd := exec.Command("systemctl", "list-units", "--type=service", "--all")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	servicesMap := make(map[string]service)
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		name := fields[0]
		if !strings.HasSuffix(name, ".service") {
			continue
		}
		loaded := ""
		active := ""
		description := ""
		if len(fields) > 1 {
			loaded = fields[1]
		}
		if len(fields) > 2 {
			active = fields[2]
		}
		if len(fields) > 3 {
			description = strings.Join(fields[3:], " ")
		}
		servicesMap[name] = service{
			name:        name,
			description: description,
			loaded:      loaded,
			active:      active,
			enabled:     "enabled", // default, will update below if found disabled
		}
	}

	// Get all disabled services
	disabledCmd := exec.Command("systemctl", "list-unit-files", "--type=service", "--state=disabled", "--no-legend")
	disabledOutput, err := disabledCmd.Output()
	if err == nil {
		disabledLines := strings.Split(string(disabledOutput), "\n")
		for _, line := range disabledLines {
			fields := strings.Fields(line)
			if len(fields) == 0 {
				continue
			}
			name := fields[0]
			if !strings.HasSuffix(name, ".service") {
				continue
			}
			if s, ok := servicesMap[name]; ok {
				s.enabled = "disabled"
				servicesMap[name] = s
			} else {
				servicesMap[name] = service{
					name:    name,
					loaded:  "",
					active:  "",
					description: "",
					enabled: "disabled",
				}
			}
		}
	}

	// Convert map to slice
	var services []list.Item
	for _, s := range servicesMap {
		services = append(services, s)
	}

	return services, nil
}

func getRunningServices() ([]list.Item, error) {
	cmd := exec.Command("systemctl", "list-units", "--type=service", "--state=running")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var services []list.Item
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) == 0 || fields[0] == "UNIT" || fields[0] == "" {
			continue
		}
		// Skip legend and summary lines
		if strings.HasPrefix(line, "Legend:") ||
			strings.Contains(line, "loaded units listed") ||
			fields[0] == "LOAD" || fields[0] == "ACTIVE" || fields[0] == "SUB" {
			continue
		}
			name := fields[0]
		loaded := ""
		active := ""
		description := ""
		if len(fields) > 1 {
			loaded = fields[1]
		}
		if len(fields) > 2 {
			active = fields[2]
		}
		if len(fields) > 3 {
			description = strings.Join(fields[3:], " ")
		}
			services = append(services, service{
				name:        name,
				description: description,
				loaded:      loaded,
				active:      active,
			})
	}

	return services, nil
}

func showAllServicesMenu(item list.Item) tea.Cmd {
	return func() tea.Msg {
		if item == nil {
			return messageMsg{text: "No service selected"}
		}
		
		s, ok := item.(service)
		if !ok {
			return messageMsg{text: "Invalid service item"}
		}

		// For now, we'll just execute a default action
		// In a full implementation, you'd want to show a proper TUI menu
		return executeServiceCommand(s.name, "status")
	}
}

func showRunningServicesMenu(item list.Item) tea.Cmd {
	return func() tea.Msg {
		if item == nil {
			return messageMsg{text: "No service selected"}
		}
		
		s, ok := item.(service)
		if !ok {
			return messageMsg{text: "Invalid service item"}
		}

		// For now, we'll just execute a default action
		// In a full implementation, you'd want to show a proper TUI menu
		return executeServiceCommand(s.name, "status")
	}
}

func executeServiceCommand(serviceName, action string) tea.Cmd {
	return func() tea.Msg {
		cmd := exec.Command("systemctl", action, serviceName)
		err := cmd.Run()
		
		if err != nil {
			return messageMsg{text: fmt.Sprintf("❌ Failed to %s %s: %v", action, serviceName, err)}
		}
		
		return messageMsg{text: fmt.Sprintf("✅ Successfully %sed %s", action, serviceName)}
	}
}

func refreshServices() tea.Cmd {
	return tea.Tick(2*time.Second, func(t time.Time) tea.Msg {
		return loadServices()
	})
}

func performSearch(searchTerm string, focused int) tea.Cmd {
	return func() tea.Msg {
		var services []list.Item
		var err error

		if focused == 0 {
			services, err = getAllServices()
		} else {
			services, err = getRunningServices()
		}

		if err != nil {
			return messageMsg{text: fmt.Sprintf("Error searching services: %v", err)}
		}

		// Filter services based on search term
		var filteredServices []list.Item
		searchLower := strings.ToLower(searchTerm)
		for _, s := range services {
			if service, ok := s.(service); ok {
				if strings.Contains(strings.ToLower(service.name), searchLower) ||
					strings.Contains(strings.ToLower(service.description), searchLower) {
					filteredServices = append(filteredServices, s)
				}
			}
		}

		if len(filteredServices) == 0 {
			return messageMsg{text: fmt.Sprintf("No services found matching '%s'", searchTerm)}
		}

		return messageMsg{text: fmt.Sprintf("Found %d services matching '%s'", len(filteredServices), searchTerm)}
	}
} 