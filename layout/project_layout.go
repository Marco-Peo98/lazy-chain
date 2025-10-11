package layout

import (
	"strings"

	"github.com/76creates/stickers/flexbox"
	"github.com/charmbracelet/lipgloss"
)

// ProjectLayout manages the layout for the project selection screen
type ProjectLayout struct {
	BaseLayout

	// Menu data
	menuItems    []string
	cursor       int
	selectedItem int

	// Content for current selection
	previewTitle        string
	previewDescription  string
	previewInstructions string
}

// NewProjectLayout creates a new layout for ProjectView
func NewProjectLayout(width, height int) *ProjectLayout {
	return &ProjectLayout{
		BaseLayout: BaseLayout{
			Width:  width,
			Height: height,
		},
		menuItems:    []string{},
		cursor:       0,
		selectedItem: -1,
	}
}

// SetMenuItems sets the menu options
func (l *ProjectLayout) SetMenuItems(items []string) *ProjectLayout {
	l.menuItems = items
	return l
}

// SetCursor sets the current cursor position
func (l *ProjectLayout) SetCursor(cursor int) *ProjectLayout {
	l.cursor = cursor
	return l
}

// SetPreviewContent sets the content for the preview panel
func (l *ProjectLayout) SetPreviewContent(title, description, instructions string) *ProjectLayout {
	l.previewTitle = title
	l.previewDescription = description
	l.previewInstructions = instructions
	return l
}

// Build constructs the FlexBox for ProjectView
func (l *ProjectLayout) Build() *flexbox.FlexBox {
	// Create main FlexBox container
	box := flexbox.New(l.Width, l.Height)

	// Create two cells with 2:3 ratio (menu:preview)
	menuCell := l.createMenuCell()
	previewCell := l.createPreviewCell()

	// Create a row with both cells
	row := box.NewRow().AddCells(menuCell, previewCell)

	// Add row to box
	box.AddRows([]*flexbox.Row{row})

	return box
}

// createMenuCell creates the left menu cell
func (l *ProjectLayout) createMenuCell() *flexbox.Cell {
	// Cell with ratio 2 (narrower)
	return flexbox.NewCell(2, 1).
		SetContentGenerator(func(maxX, maxY int) string {
			return l.renderMenuContent(maxX, maxY)
		})
}

// createPreviewCell creates the right preview cell
func (l *ProjectLayout) createPreviewCell() *flexbox.Cell {
	// Cell with ratio 3 (wider)
	return flexbox.NewCell(3, 1).
		SetContentGenerator(func(maxX, maxY int) string {
			return l.renderPreviewContent(maxX, maxY)
		})
}

// renderMenuContent renders the menu section content
func (l *ProjectLayout) renderMenuContent(maxX, maxY int) string {
	var content []string

	// Section title
	title := TitleStyle().Render("Main Menu")
	content = append(content, title)
	content = append(content, "")

	// Menu options
	for i, option := range l.menuItems {
		cursor := "  "
		style := lipgloss.NewStyle()

		// Show cursor for current option
		if l.cursor == i {
			cursor = "> "
			style = style.Foreground(TextColorHighlight)
		}

		// Style the option text
		line := cursor + style.Render(option)
		content = append(content, line)
	}

	// Add spacing to fill height
	for len(content) < 10 {
		content = append(content, "")
	}

	// Join content
	panelContent := strings.Join(content, "\n")

	// Create bordered panel with available width
	// Account for padding and border (4 chars total: 2 padding + 2 border)
	contentWidth := maxX - 4
	if contentWidth < 10 {
		contentWidth = 10
	}

	return lipgloss.NewStyle().
		Width(contentWidth).
		Padding(1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(BorderColorPrimary).
		Render(panelContent)
}

// renderPreviewContent renders the preview section content
func (l *ProjectLayout) renderPreviewContent(maxX, maxY int) string {
	var content []string

	// Section title
	title := TitleStyle().Render("Preview")
	content = append(content, title)
	content = append(content, "")

	// Preview content for selected option
	if l.previewTitle != "" {
		optionStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(TextColorSuccess)
		content = append(content, optionStyle.Render(l.previewTitle))
		content = append(content, "")
	}

	// Option description
	if l.previewDescription != "" {
		descStyle := lipgloss.NewStyle().
			Foreground(TextColorPrimary)

		// Wrap description to fit in panel width
		// Account for padding and border (4 chars total)
		wrapWidth := maxX - 8
		if wrapWidth < 20 {
			wrapWidth = 20
		}

		wrappedDesc := wrapText(l.previewDescription, wrapWidth)
		for _, line := range wrappedDesc {
			content = append(content, descStyle.Render(line))
		}
	}

	content = append(content, "")

	// Instructions
	if l.previewInstructions != "" {
		instrStyle := lipgloss.NewStyle().
			Italic(true).
			Foreground(TextColorSecondary)

		instrLines := strings.Split(l.previewInstructions, "\n")
		for _, line := range instrLines {
			content = append(content, instrStyle.Render(line))
		}
	}

	// Fill remaining space
	for len(content) < 10 {
		content = append(content, "")
	}

	// Join content
	panelContent := strings.Join(content, "\n")

	// Create bordered panel with available width
	contentWidth := maxX - 4
	if contentWidth < 20 {
		contentWidth = 20
	}

	return lipgloss.NewStyle().
		Width(contentWidth).
		Padding(1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(BorderColorSecondary).
		Render(panelContent)
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
				// Single word too long, truncate it
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

// GetMinDimensions returns the minimum required dimensions
func (l *ProjectLayout) GetMinDimensions() (width, height int) {
	// ProjectView needs reasonable space for two columns
	return 80, 24
}

// IsValid checks if dimensions are sufficient
func (l *ProjectLayout) IsValid() bool {
	minWidth, minHeight := l.GetMinDimensions()
	return l.Width >= minWidth && l.Height >= minHeight
}

// Update updates the layout dimensions
func (l *ProjectLayout) Update(width, height int) {
	l.Width = width
	l.Height = height
}
