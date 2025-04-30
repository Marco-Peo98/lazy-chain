package models

import (
	"container/list"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
			HandlePreview(m.Cursor, m.Options, m.Selected)
		}
	}

	return m, nil
}

// ProjectView has to show the preview of each option's view while the user is hovering it
// After the user selects an option, the focus has to move to the Selected option's view

func (m *ProjectModel) View() string {

	// Create a new list to store the options with their selected state
	row := list.New()

	for i, option := range m.Options {
		// Subtitle for each option
		desc := subtitleFor(option)

		// Styles for title and subtitle
		titleStyled := lipgloss.NewStyle().Bold(true).Render(option)
		descStyled := lipgloss.NewStyle().Italic(true).Render(desc)

		block := fmt.Sprintf("%s\n%s\n", titleStyled, descStyled)

		style := lipgloss.NewStyle().Foreground(lipgloss.Color("#c6d0f5")).Padding(0, 1)

		if m.Cursor == i {
			style = style.Copy().
				Foreground(lipgloss.Color("#ef9f76")).
				BorderLeft(true).
				BorderLeftForeground(lipgloss.Color("#ef9f76"))
		}

		r := style.Render(block)
		row.PushBack(r)
	}

	// Vertical wrapper for the list items
	var items []string
	for e := row.Front(); e != nil; e = e.Next() {
		if str, ok := e.Value.(string); ok {
			items = append(items, str)
		}
	}
	listCol := lipgloss.JoinVertical(lipgloss.Left, items...)
	listBox := lipgloss.NewStyle().
		Align(lipgloss.Left).
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#ca9ee6")).
		Render(listCol)

	// Dynamic preview
	current := m.Options[m.Cursor]
	prevTitle := lipgloss.NewStyle().Bold(true).Render(current)
	prevDesc := lipgloss.NewStyle().Italic(true).Render(subtitleFor(current))
	previewText := fmt.Sprintf(
		"Preview of: %s\n\n%s\n\nComing soon...\n\n",
		prevTitle, prevDesc,
	)
	previewBox := lipgloss.NewStyle().
		Align(lipgloss.Left).
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#89b4fa")).
		Render(previewText)

	return lipgloss.JoinHorizontal(lipgloss.Top, listBox, previewBox)
}

// subtitleFor returns the subtitle for the given option
func subtitleFor(option string) string {
	switch option {
	case "Settings":
		return "Configure your settings"
	case "Applications":
		return "Manage your applications"
	case "Commands Goals":
		return "Define your objectives from shell"
	case "Explore":
		return "Explore other infos and resources"
	default:
		return ""
	}
}

// User can select only one option at a time
// If user selects more than one option, the last one will be the Selected one
// and the previous ones will be unSelected

func HandlePreview(Cursor int, Options []string, Selected map[int]struct{}) {

	// Reset the previous selected option, maintaining the current one
	for k := range Selected {
		delete(Selected, k)
	}

	Selected[Cursor] = struct{}{}
}
