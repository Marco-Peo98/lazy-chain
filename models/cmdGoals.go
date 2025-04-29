package models

import (
	tea "github.com/charmbracelet/bubbletea"
)

// CmdGoalsModel isn't ready yet, but it's a placeholder for now.

type CmdGoalsModel struct {
	CurrentState SessionState
}

func NewCmdGoalsModel() *CmdGoalsModel {
	return &CmdGoalsModel{
		CurrentState: CmdGoalsView,
	}
}

func (m *CmdGoalsModel) Init() tea.Cmd {
	return nil
}

func (m *CmdGoalsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m *CmdGoalsModel) View() string {
	return "Commands Goals View\n\nPress ESC to go back to the main view.\n"
}
