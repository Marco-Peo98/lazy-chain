package layout

import (
	"fmt"

	"github.com/76creates/stickers/flexbox"
	"github.com/charmbracelet/lipgloss"
)

// MainLayout manages the layout for the main screen
type MainLayout struct {
	BaseLayout
	banner       string
	instructions string
}

// NewMainLayout creates a new layout for MainView
func NewMainLayout(width, height int, banner, instructions string) *MainLayout {
	return &MainLayout{
		BaseLayout: BaseLayout{
			Width:  width,
			Height: height,
		},
		banner:       banner,
		instructions: instructions,
	}
}

// Build constructs the FlexBox for MainView
func (l *MainLayout) Build() *flexbox.FlexBox {
	// Create the main FlexBox container
	box := flexbox.New(l.Width, l.Height)

	// Style the banner (colored)
	bannerStyled := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#81c8be")).
		Render(l.banner)

	// Style the instructions
	instructionsStyled := lipgloss.NewStyle().
		Foreground(TextColorPrimary).
		Render(l.instructions)

	// Combine banner + instructions vertically
	content := lipgloss.JoinVertical(
		lipgloss.Center,
		bannerStyled,
		"",
		instructionsStyled,
	)

	// Create a cell that occupies all available space
	// IMPORTANT: Use SetContentGenerator to dynamically center content
	cell := flexbox.NewCell(1, 1).
		SetContentGenerator(func(maxX, maxY int) string {
			// Center content using lipgloss.Place
			// maxX and maxY are the available dimensions in the cell
			return lipgloss.Place(
				maxX,
				maxY,
				lipgloss.Center, // Center horizontally
				lipgloss.Center, // Center vertically
				content,
			)
		})

	// Create a row with the cell
	row := box.NewRow().AddCells(cell)

	// Add the row to the box
	box.AddRows([]*flexbox.Row{row})

	return box
}

// IsValid checks if dimensions are sufficient
func (l *MainLayout) IsValid() bool {
	minWidth, minHeight := l.GetMinDimensions()
	return l.Width >= minWidth && l.Height >= minHeight
}

// GetMinDimensions returns the minimum required dimensions
func (l *MainLayout) GetMinDimensions() (width, height int) {
	// MainView requires at least 80x24 (standard terminal size)
	return 80, 24
}

// RenderError displays an error message when the terminal is too small
func (l *MainLayout) RenderError() string {
	minWidth, minHeight := l.GetMinDimensions()

	// Calculate missing dimensions
	needsWidth := minWidth - l.Width
	needsHeight := minHeight - l.Height

	// Create error message
	errorTitle := lipgloss.NewStyle().
		Bold(true).
		Foreground(BorderColorError).
		Render("Terminal Too Small")

	currentDims := fmt.Sprintf("Current: %d × %d", l.Width, l.Height)
	requiredDims := fmt.Sprintf("Required: %d × %d", minWidth, minHeight)

	// Determine what needs to be increased
	var suggestion string
	if needsWidth > 0 && needsHeight > 0 {
		suggestion = fmt.Sprintf("Please increase width by %d and height by %d", needsWidth, needsHeight)
	} else if needsWidth > 0 {
		suggestion = fmt.Sprintf("Please increase width by %d", needsWidth)
	} else if needsHeight > 0 {
		suggestion = fmt.Sprintf("Please increase height by %d", needsHeight)
	}

	// Style for various components
	currentStyle := lipgloss.NewStyle().Foreground(BorderColorWarning)
	requiredStyle := lipgloss.NewStyle().Foreground(BorderColorSecondary)
	suggestionStyle := lipgloss.NewStyle().
		Italic(true).
		Foreground(TextColorSecondary)

	// Combine all message components
	message := lipgloss.JoinVertical(
		lipgloss.Center,
		errorTitle,
		"",
		currentStyle.Render(currentDims),
		requiredStyle.Render(requiredDims),
		"",
		suggestionStyle.Render(suggestion),
	)

	// Create container for error message
	errorContainer := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(BorderColorError).
		Padding(1, 2).
		Render(message)

	// Center the message in available space
	// Use current dimensions (even if small)
	return lipgloss.Place(
		l.Width,
		l.Height,
		lipgloss.Center,
		lipgloss.Center,
		errorContainer,
	)
}

// Update updates the layout dimensions
func (l *MainLayout) Update(width, height int) {
	l.Width = width
	l.Height = height
}
