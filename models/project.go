package models

import (
	"strings"

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
			// Clear any previous selection and select current item
			m.Selected = make(map[int]struct{})
			m.Selected[m.Cursor] = struct{}{}
		}
	}

	return m, nil
}

func (m *ProjectModel) View() string {
	// Create two-column layout: Menu options (left) | Preview (right)
	// Using same dimensions as SettingsModel for consistency
	leftColumn := m.renderMenuSection()
	rightColumn := m.renderPreviewSection()

	// Combine columns horizontally with same gap as SettingsModel
	mainContent := lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftColumn,
		"  ", // Small gap between columns (same as SettingsModel)
		rightColumn,
	)

	// Add footer with instructions
	footer := m.renderFooter()

	// Combine everything vertically
	return lipgloss.JoinVertical(
		lipgloss.Left,
		mainContent,
		"",
		footer,
	)
}

func (m *ProjectModel) renderMenuSection() string {
	var content []string

	// Section title
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#cba6f7")).
		Render("Main Menu")
	content = append(content, title)
	content = append(content, "")

	// Menu options
	for i, option := range m.Options {
		cursor := "  "
		style := lipgloss.NewStyle()

		// Show cursor for current option
		if m.Cursor == i {
			cursor = "> "
			style = style.Foreground(lipgloss.Color("#ef9f76"))
		}

		// Style the option text
		line := cursor + style.Render(option)
		content = append(content, line)
	}

	// Add some spacing
	for len(content) < 10 {
		content = append(content, "")
	}

	// Create left panel with border (standardized dimensions with increased padding)
	panelContent := strings.Join(content, "\n")
	return lipgloss.NewStyle().
		Width(30).  // Same as SettingsModel left panel
		Padding(2). // Increased padding
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#89b4fa")).
		Render(panelContent)
}

func (m *ProjectModel) renderPreviewSection() string {
	var content []string

	// Section title
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#cba6f7")).
		Render("Preview")
	content = append(content, title)
	content = append(content, "")

	// Preview content for selected option
	currentOption := m.Options[m.Cursor]

	// Option name
	optionStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#a6e3a1"))
	content = append(content, optionStyle.Render(currentOption))
	content = append(content, "")

	// Option description
	description := subtitleFor(currentOption)
	descStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#cdd6f4"))

	// Wrap description to fit in panel width (adjusted for increased padding)
	wrappedDesc := wrapText(description, 37) // Adjusted for increased padding
	for _, line := range wrappedDesc {
		content = append(content, descStyle.Render(line))
	}

	content = append(content, "")

	// Instructions
	instrStyle := lipgloss.NewStyle().
		Italic(true).
		Foreground(lipgloss.Color("#6c7086"))
	content = append(content, instrStyle.Render("Press ENTER to select"))
	content = append(content, instrStyle.Render("Press ESC to go back"))

	// Fill remaining space
	for len(content) < 10 {
		content = append(content, "")
	}

	// Create right panel with border (standardized dimensions with increased padding)
	panelContent := strings.Join(content, "\n")
	return lipgloss.NewStyle().
		Width(45).  // Same as SettingsModel right panel
		Padding(2). // Increased padding
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#a6e3a1")).
		Render(panelContent)
}

func (m *ProjectModel) renderFooter() string {
	instructions := []string{
		"Up/Down: Navigate options",
		"Enter: Select option",
		"ESC: Back to main",
	}

	instructionText := strings.Join(instructions, " | ")
	return lipgloss.NewStyle().
		Width(77).              // Total width of both panels + gap (30+2+45)
		Align(lipgloss.Center). // Center the footer text
		Foreground(lipgloss.Color("#6c7086")).
		Italic(true).
		Render(instructionText)
}

// subtitleFor returns the subtitle for the given option
func subtitleFor(option string) string {
	switch option {
	case "Settings":
		return "Configure your network and wallet settings"
	case "Applications":
		return "Manage your blockchain applications"
	case "Commands Goals":
		return "Why CLI when you can TUI? Build transactions easily"
	case "Explore":
		return "Explore blockchain data and resources"
	default:
		return ""
	}
}

// wrapText wraps text to fit within specified width
func wrapText(text string, width int) []string {
	if len(text) <= width {
		return []string{text}
	}

	var lines []string
	words := strings.Fields(text)
	var currentLine []string
	currentLength := 0

	for _, word := range words {
		// Check if adding this word would exceed width
		if currentLength+len(word)+len(currentLine) > width {
			if len(currentLine) > 0 {
				lines = append(lines, strings.Join(currentLine, " "))
				currentLine = []string{word}
				currentLength = len(word)
			} else {
				// Single word is too long, truncate it
				lines = append(lines, word[:width-3]+"...")
				currentLine = []string{}
				currentLength = 0
			}
		} else {
			currentLine = append(currentLine, word)
			currentLength += len(word)
		}
	}

	// Add remaining words
	if len(currentLine) > 0 {
		lines = append(lines, strings.Join(currentLine, " "))
	}

	return lines
}
