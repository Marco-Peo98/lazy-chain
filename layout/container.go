package layout

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

// LayoutContainer provides a simple container system for TUI screens
type LayoutContainer struct {
	TerminalWidth  int
	TerminalHeight int

	// Calculated dimensions
	ContentWidth  int
	ContentHeight int

	// Configuration
	Margin  int
	Padding int

	// Minimum required dimensions
	MinWidth  int
	MinHeight int
}

// NewLayoutContainer creates a new layout container with terminal dimensions
func NewLayoutContainer(width, height int) *LayoutContainer {
	container := &LayoutContainer{
		TerminalWidth:  width,
		TerminalHeight: height,
		Margin:         2,  // 2 chars margin on each side
		Padding:        1,  // 1 char padding inside border
		MinWidth:       80, // Minimum terminal width for proper display
		MinHeight:      24, // Minimum terminal height for proper display
	}

	container.calculateDimensions()
	return container
}

// calculateDimensions calculates the available content area
func (lc *LayoutContainer) calculateDimensions() {
	// Calculate content area: Terminal - (Margin * 2) - (Border * 2) - (Padding * 2)
	lc.ContentWidth = lc.TerminalWidth - (lc.Margin * 2) - 2 - (lc.Padding * 2)
	lc.ContentHeight = lc.TerminalHeight - (lc.Margin * 2) - 2 - (lc.Padding * 2)

	// Ensure minimum dimensions
	if lc.ContentWidth < 10 {
		lc.ContentWidth = 10
	}
	if lc.ContentHeight < 5 {
		lc.ContentHeight = 5
	}
}

// IsValid checks if terminal dimensions are sufficient for proper display
func (lc *LayoutContainer) IsValid() bool {
	return lc.TerminalWidth >= lc.MinWidth && lc.TerminalHeight >= lc.MinHeight
}

// GetContentDimensions returns the available content area dimensions
func (lc *LayoutContainer) GetContentDimensions() (width, height int) {
	return lc.ContentWidth, lc.ContentHeight
}

// Resize updates the container with new terminal dimensions
func (lc *LayoutContainer) Resize(width, height int) {
	lc.TerminalWidth = width
	lc.TerminalHeight = height
	lc.calculateDimensions()
}

// Render wraps content in the container with border, centering, and overflow handling
func (lc *LayoutContainer) Render(content string) string {
	// Check if terminal is too small
	if !lc.IsValid() {
		return lc.renderOverflowMessage()
	}

	// Create the main container style
	containerStyle := lipgloss.NewStyle().
		Width(lc.TerminalWidth-(lc.Margin*2)).
		Height(lc.TerminalHeight-(lc.Margin*2)).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#89b4fa")).
		Padding(lc.Padding).
		Align(lipgloss.Center, lipgloss.Center)

	// Wrap content and center it
	wrappedContent := containerStyle.Render(content)

	// Center the entire container in the terminal
	return lipgloss.Place(
		lc.TerminalWidth,
		lc.TerminalHeight,
		lipgloss.Center,
		lipgloss.Center,
		wrappedContent,
	)
}

// renderOverflowMessage displays an error message when terminal is too small
func (lc *LayoutContainer) renderOverflowMessage() string {
	// Calculate what dimensions are needed
	requiredWidth := lc.MinWidth
	requiredHeight := lc.MinHeight

	// Create error message
	errorTitle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#f38ba8")).
		Render("Terminal Too Small")

	currentDims := fmt.Sprintf("Current: %d x %d", lc.TerminalWidth, lc.TerminalHeight)
	requiredDims := fmt.Sprintf("Required: %d x %d", requiredWidth, requiredHeight)

	// Determine what action user should take
	var suggestion string
	if lc.TerminalWidth < requiredWidth && lc.TerminalHeight < requiredHeight {
		suggestion = "Please increase both width and height"
	} else if lc.TerminalWidth < requiredWidth {
		suggestion = "Please increase terminal width"
	} else if lc.TerminalHeight < requiredHeight {
		suggestion = "Please increase terminal height"
	} else {
		suggestion = "Please resize terminal"
	}

	// Style the message components
	currentStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#f9e2af"))
	requiredStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#a6e3a1"))
	suggestionStyle := lipgloss.NewStyle().
		Italic(true).
		Foreground(lipgloss.Color("#cdd6f4"))

	// Combine all message parts
	message := lipgloss.JoinVertical(
		lipgloss.Center,
		errorTitle,
		"",
		currentStyle.Render(currentDims),
		requiredStyle.Render(requiredDims),
		"",
		suggestionStyle.Render(suggestion),
	)

	// Create a simple border for the error message
	errorContainer := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#f38ba8")).
		Padding(1, 2).
		Render(message)

	// Center the error message in the available space
	return lipgloss.Place(
		lc.TerminalWidth,
		lc.TerminalHeight,
		lipgloss.Center,
		lipgloss.Center,
		errorContainer,
	)
}

// CreateContentStyle creates a style for content that fits within the container
func (lc *LayoutContainer) CreateContentStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Width(lc.ContentWidth).
		Height(lc.ContentHeight)
}

// WrapContent wraps content to fit within the container dimensions and centers it
func (lc *LayoutContainer) WrapContent(content string) string {
	// Simple content wrapping - ensure it fits in available space
	style := lc.CreateContentStyle().
		Align(lipgloss.Center, lipgloss.Center) // Add centering for better alignment
	return style.Render(content)
}
