package models

import (
	tea "github.com/charmbracelet/bubbletea"
)

// SettingsModel isn't ready yet, but it's a placeholder for now.

type SettingsModel struct {
	CurrentState SessionState
}

func NewSettingsModel() *SettingsModel {
	return &SettingsModel{
		CurrentState: SettingsView,
	}
}

func (m *SettingsModel) Init() tea.Cmd {
	return nil
}

func (m *SettingsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m *SettingsModel) View() string {
	return "Settings View\n\nPress ESC to go back to the main view.\n"
}
