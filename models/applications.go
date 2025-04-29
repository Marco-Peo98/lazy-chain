package models

import (
	tea "github.com/charmbracelet/bubbletea"
)

// ApplicationsModel isn't ready yet, but it's a placeholder for now.

type ApplicationsModel struct {
	CurrentState SessionState
}

func NewApplicationsModel() *ApplicationsModel {
	return &ApplicationsModel{
		CurrentState: ApplicationsView,
	}
}

func (m *ApplicationsModel) Init() tea.Cmd {
	return nil
}

func (m *ApplicationsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m *ApplicationsModel) View() string {
	return "Applications View\n\nPress ESC to go back to the main view.\n"
}
