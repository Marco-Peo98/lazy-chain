package layout

import (
	"fmt"
	"strings"

	"github.com/76creates/stickers/flexbox"
	"github.com/charmbracelet/lipgloss"
)

// ProjectLayout manages the layout for the project selection screen
type ProjectLayout struct {
	BaseLayout

	// Menu data
	menuItems []string
	cursor    int

	// All possible preview contents (pre-calculated)
	allPreviews map[string]PreviewContent

	// Maximum dimensions based on content
	maxMenuWidth     int
	maxMenuHeight    int
	maxPreviewWidth  int
	maxPreviewHeight int

	// External padding (space between terminal edge and box)
	externalPadding int
}

// PreviewContent holds the content for a preview panel
type PreviewContent struct {
	Title        string
	Description  string
	Instructions string
}

// NewProjectLayout creates a new layout for ProjectView
func NewProjectLayout(width, height int) *ProjectLayout {
	return &ProjectLayout{
		BaseLayout: BaseLayout{
			Width:  width,
			Height: height,
		},
		menuItems:       []string{},
		cursor:          0,
		allPreviews:     make(map[string]PreviewContent),
		externalPadding: 2, // 2 chars on each side
	}
}

// SetMenuItems sets the menu options and pre-calculates all previews
func (l *ProjectLayout) SetMenuItems(items []string, previewGenerator func(string) PreviewContent) *ProjectLayout {
	l.menuItems = items

	// Pre-calculate all preview contents
	l.allPreviews = make(map[string]PreviewContent)
	for _, item := range items {
		l.allPreviews[item] = previewGenerator(item)
	}

	// Calculate maximum dimensions needed
	l.calculateMaxDimensions()

	return l
}

// SetCursor sets the current cursor position
func (l *ProjectLayout) SetCursor(cursor int) *ProjectLayout {
	l.cursor = cursor
	return l
}

// calculateMaxDimensions calculates the maximum width/height needed
// based on the largest content among all options
func (l *ProjectLayout) calculateMaxDimensions() {
	// Reset max values
	l.maxMenuWidth = 0
	l.maxMenuHeight = 0
	l.maxPreviewWidth = 0
	l.maxPreviewHeight = 0

	// Calculate menu dimensions
	for _, item := range l.menuItems {
		// Account for cursor ("> ") + item text
		itemWidth := len(item) + 2
		if itemWidth > l.maxMenuWidth {
			l.maxMenuWidth = itemWidth
		}
	}

	// Menu height: title + spacing + items + extra spacing
	l.maxMenuHeight = 2 + len(l.menuItems) + 3

	// Calculate preview dimensions for all possible contents
	for _, preview := range l.allPreviews {
		// Title
		titleWidth := len(preview.Title)

		// Description (needs wrapping calculation)
		// Estimate: assume 40 chars per line for wrapping
		descLines := (len(preview.Description) / 40) + 1
		if len(preview.Description) > 40 {
			descLines = len(wrapText(preview.Description, 40))
		}

		// Instructions (count newlines)
		instrLines := strings.Count(preview.Instructions, "\n") + 1

		// Calculate height: title + spacing + desc + spacing + instr + extra
		totalHeight := 2 + descLines + 1 + instrLines + 3
		if totalHeight > l.maxPreviewHeight {
			l.maxPreviewHeight = totalHeight
		}

		// Calculate width: longest line
		maxLineWidth := titleWidth
		for _, line := range wrapText(preview.Description, 50) {
			if len(line) > maxLineWidth {
				maxLineWidth = len(line)
			}
		}
		if maxLineWidth > l.maxPreviewWidth {
			l.maxPreviewWidth = maxLineWidth
		}
	}

	// Add minimum padding and border space to dimensions
	// Border (2) + Padding (2) = 4 extra chars
	l.maxMenuWidth += 4
	l.maxMenuHeight += 2
	l.maxPreviewWidth += 4
	l.maxPreviewHeight += 2

	// Ensure minimums
	if l.maxMenuWidth < 25 {
		l.maxMenuWidth = 25
	}
	if l.maxPreviewWidth < 40 {
		l.maxPreviewWidth = 40
	}
	if l.maxMenuHeight < 15 {
		l.maxMenuHeight = 15
	}
	if l.maxPreviewHeight < 15 {
		l.maxPreviewHeight = 15
	}
}

// Build constructs the FlexBox for ProjectView
func (l *ProjectLayout) Build() *flexbox.FlexBox {
	// Calculate box dimensions (terminal - external padding)
	boxWidth := l.Width - (l.externalPadding * 2)
	boxHeight := l.Height - (l.externalPadding * 2)

	// Ensure minimum box size
	if boxWidth < 70 {
		boxWidth = 70
	}
	if boxHeight < 20 {
		boxHeight = 20
	}

	// Use the maximum height between menu and preview
	maxCellHeight := l.maxMenuHeight
	if l.maxPreviewHeight > maxCellHeight {
		maxCellHeight = l.maxPreviewHeight
	}

	// If calculated height exceeds available space, use available space
	if maxCellHeight > boxHeight {
		maxCellHeight = boxHeight
	}

	// Create FlexBox with calculated dimensions
	box := flexbox.New(boxWidth, maxCellHeight)

	// Create cells with ratio 2:3 but they'll render at calculated widths
	menuCell := l.createMenuCell()
	previewCell := l.createPreviewCell()

	// Create row with both cells
	row := box.NewRow().AddCells(menuCell, previewCell)

	// Add row to box
	box.AddRows([]*flexbox.Row{row})

	return box
}

// Render renders the complete layout centered in terminal
func (l *ProjectLayout) Render() string {
	// Build the FlexBox
	box := l.Build()

	// Render the FlexBox
	boxRendered := box.Render()

	// Center the box in the terminal with external padding
	return lipgloss.Place(
		l.Width,
		l.Height,
		lipgloss.Center,
		lipgloss.Center,
		boxRendered,
	)
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

	// Join content
	panelContent := strings.Join(content, "\n")

	// Calculate content width accounting for border and padding
	contentWidth := maxX - 4
	if contentWidth < 10 {
		contentWidth = 10
	}

	// Create bordered panel
	return lipgloss.NewStyle().
		Width(contentWidth).
		Height(maxY - 2). // Account for border
		Padding(1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(BorderColorPrimary).
		Render(panelContent)
}

// renderPreviewContent renders the preview section content
func (l *ProjectLayout) renderPreviewContent(maxX, maxY int) string {
	// Get current option
	currentOption := ""
	if l.cursor >= 0 && l.cursor < len(l.menuItems) {
		currentOption = l.menuItems[l.cursor]
	}

	// Get preview content
	preview, exists := l.allPreviews[currentOption]
	if !exists {
		preview = PreviewContent{
			Title:        currentOption,
			Description:  "No description available",
			Instructions: "",
		}
	}

	var content []string

	// Section title
	title := TitleStyle().Render("Preview")
	content = append(content, title)
	content = append(content, "")

	// Option title
	if preview.Title != "" {
		optionStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(TextColorSuccess)
		content = append(content, optionStyle.Render(preview.Title))
		content = append(content, "")
	}

	// Option description
	if preview.Description != "" {
		descStyle := lipgloss.NewStyle().
			Foreground(TextColorPrimary)

		// Wrap description to fit in available width
		wrapWidth := maxX - 8
		if wrapWidth < 20 {
			wrapWidth = 20
		}

		wrappedDesc := wrapText(preview.Description, wrapWidth)
		for _, line := range wrappedDesc {
			content = append(content, descStyle.Render(line))
		}
	}

	content = append(content, "")

	// Instructions
	if preview.Instructions != "" {
		instrStyle := lipgloss.NewStyle().
			Italic(true).
			Foreground(TextColorSecondary)

		instrLines := strings.Split(preview.Instructions, "\n")
		for _, line := range instrLines {
			content = append(content, instrStyle.Render(line))
		}
	}

	// Join content
	panelContent := strings.Join(content, "\n")

	// Calculate content width accounting for border and padding
	contentWidth := maxX - 4
	if contentWidth < 20 {
		contentWidth = 20
	}

	// Create bordered panel
	return lipgloss.NewStyle().
		Width(contentWidth).
		Height(maxY - 2). // Account for border
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

// RenderError displays an error message when terminal is too small
// Uses the same style as MainLayout
func (l *ProjectLayout) RenderError() string {
	minWidth, minHeight := l.GetMinDimensions()

	// Calculate missing dimensions
	needsWidth := minWidth - l.Width
	needsHeight := minHeight - l.Height

	// Create error message
	errorTitle := lipgloss.NewStyle().
		Bold(true).
		Foreground(BorderColorError).
		Render("Terminal Too Small")

	currentDims := lipgloss.NewStyle().
		Foreground(BorderColorWarning).
		Render(lipgloss.NewStyle().Render(
			lipgloss.NewStyle().Render("Current: ") +
				lipgloss.NewStyle().Bold(true).Render(
					lipgloss.NewStyle().Render(
						fmt.Sprintf("%d × %d", l.Width, l.Height),
					),
				),
		))

	requiredDims := lipgloss.NewStyle().
		Foreground(BorderColorSecondary).
		Render(lipgloss.NewStyle().Render(
			lipgloss.NewStyle().Render("Required: ") +
				lipgloss.NewStyle().Bold(true).Render(
					lipgloss.NewStyle().Render(
						fmt.Sprintf("%d × %d", minWidth, minHeight),
					),
				),
		))

	// Determine what needs to be increased
	var suggestion string
	if needsWidth > 0 && needsHeight > 0 {
		suggestion = fmt.Sprintf("Please increase width by %d and height by %d", needsWidth, needsHeight)
	} else if needsWidth > 0 {
		suggestion = fmt.Sprintf("Please increase width by %d", needsWidth)
	} else if needsHeight > 0 {
		suggestion = fmt.Sprintf("Please increase height by %d", needsHeight)
	}

	suggestionStyled := lipgloss.NewStyle().
		Italic(true).
		Foreground(TextColorSecondary).
		Render(suggestion)

	// Combine all message components
	message := lipgloss.JoinVertical(
		lipgloss.Center,
		errorTitle,
		"",
		currentDims,
		requiredDims,
		"",
		suggestionStyled,
	)

	// Create container for error message
	errorContainer := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(BorderColorError).
		Padding(1, 2).
		Render(message)

	// Center the message in available space
	return lipgloss.Place(
		l.Width,
		l.Height,
		lipgloss.Center,
		lipgloss.Center,
		errorContainer,
	)
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
