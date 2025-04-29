package models

import (
	tea "github.com/charmbracelet/bubbletea"
)

// ExploreModel isn't ready yet, but it's a placeholder for now.

type ExploreModel struct {
	CurrentState SessionState
}

func NewExploreModel() *ExploreModel {
	return &ExploreModel{
		CurrentState: ExploreView,
	}
}

func (m *ExploreModel) Init() tea.Cmd {
	return nil
}

func (m *ExploreModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m *ExploreModel) View() string {
	return "Explore View\n\nPress ESC to go back to the main view.\n"
}
