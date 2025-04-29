package models

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

// ProjectModel is the first real accessible model for the user
// It will be the main entry point for the user to interact with the application
// selecting which screen to go to, from the main list of Options

type ProjectModel struct {
	CurrentState SessionState
	Options      []string
	Cursor       int
	Selected     map[int]struct{} // Selected items
}

func NewProjectModel() *ProjectModel {
	return &ProjectModel{
		Options: []string{
			"Settings",
			"Applications",
			"Commands Goals",
			"Explore",
		},
		Cursor:   0,
		Selected: make(map[int]struct{}),
	}
}

func (m *ProjectModel) Init() tea.Cmd {
	return nil
}

func (m *ProjectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.Cursor > 0 {
				m.Cursor--
			}
		case "down", "j":
			if m.Cursor < len(m.Options)-1 {
				m.Cursor++
			}
		case "enter", " ":
			return m, func() tea.Msg {
				return HandlePreview(m.Cursor, m.Options, m.Selected)
			}
		}
	}

	return m, nil
}

// ProjectView has to show the preview of each option's view while the user is hovering it
// After the user selects an option, the focus has to move to the Selected option's view

func (m *ProjectModel) View() string {
	// Header of the list
	s := "Which option do you want to select?\n\n"

	// Iterate over the list of Options
	for i, option := range m.Options {
		// Cursor position
		Cursor := " " // No Cursor
		if m.Cursor == i {
			Cursor = ">" // Cursor
		}

		Selected := " "
		if _, ok := m.Selected[i]; ok {
			Selected = "x"
		}

		s += fmt.Sprintf("%s [%s] %s\n", Cursor, Selected, option)
	}

	return s + "\nPress q to quit.\n"
}

// User can select only one option at a time
// If user selects more than one option, the last one will be the Selected one
// and the previous ones will be unSelected

func HandlePreview(Cursor int, Options []string, Selected map[int]struct{}) tea.Msg {

	if _, ok := Selected[Cursor]; ok {
		// Unselect all previous options
		for i := range Selected {
			if i != Cursor {
				delete(Selected, i)
			}
		}
	}
	// Select the current option
	Selected[Cursor] = struct{}{}
	return Options[Cursor] + "View"
}
