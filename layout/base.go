package layout

import (
	"github.com/76creates/stickers/flexbox"
	"github.com/charmbracelet/lipgloss"
)

// Layout is the interface that all layout types must implement.
type Layout interface {
	Build(width, height int) *flexbox.FlexBox
	GetMinDimensions() (width, height int)
}

// BaseLayout provides common utilities for all layout types.
type BaseLayout struct {
	Width  int
	Height int
}

// Reusable common styles
var (

	// Borders and colors from Catppuccin Mocha theme
	BorderColorPrimary   = lipgloss.Color("#89b4fa") // Blue
	BorderColorSecondary = lipgloss.Color("#a6e3a1") // Green
	BorderColorAccent    = lipgloss.Color("#cba6f7") // Mauve
	BorderColorWarning   = lipgloss.Color("#f9e2af") // Yellow
	BorderColorError     = lipgloss.Color("#f38ba8") // Red

	// Text colors
	TextColorPrimary   = lipgloss.Color("#cdd6f4") // Text
	TextColorSecondary = lipgloss.Color("#6c7086") // Subtext
	TextColorHighlight = lipgloss.Color("#ef9f76") // Peach
	TextColorSuccess   = lipgloss.Color("#a6e3a1") // Green
	TextColorTitle     = lipgloss.Color("#cba6f7") // Mauve

	// Background
	BackgroundBase = lipgloss.Color("#1e1e2e") // Base
)

// DefaultBorder returns a default rounded border style.
func DefaultBorder(color lipgloss.Color) lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(color).
		Padding(1)
}

// TytleStyle returns a style for titles.
func TitleStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(TextColorTitle)
}

// SubtitleStyle returns a style for subtitles.
func SubtitleStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Italic(true).
		Foreground(TextColorSecondary)
}

// HighlightStyle returns a style for highlighted text.
func HighlightStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(TextColorHighlight).
		Bold(true)
}

// CenteredCell creates a centered cell with content.
func CenteredCell(ratioX, ratioY int, content string) *flexbox.Cell {
	return flexbox.NewCell(ratioX, ratioY).
		SetContent(content)
}

// CenteredCellWithStyle creates a centered cell with content and style.
func CenteredCellWithStyle(ratioX, ratioY int, content string, style lipgloss.Style) *flexbox.Cell {
	return flexbox.NewCell(ratioX, ratioY).
		SetContent(content).
		SetStyle(style)
}

// WrapInBorder wraps content in a border with padding.
func WrapInBorder(content string, borderColor lipgloss.Color, width int) string {
	style := lipgloss.NewStyle().
		Width(width).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Padding(1, 2)
	return style.Render(content)
}

// ValidateMinimumSize checks if the provided width and height meet minimum requirements.
func ValidateMinimumSize(width, height, minWidth, minHeight int) bool {
	return width >= minWidth && height >= minHeight
}

// GetContentDimensions calculates the dimensions of the content area after accounting for borders and padding.
func GetContentDimensions(totalWidth, totalHeight, margin, padding, border int) (width, height int) {
	width = totalWidth - (margin * 2) - (padding * 2) - (border * 2)
	height = totalHeight - (margin * 2) - (padding * 2) - (border * 2)

	if width < 10 {
		width = 10 // Minimum width
	}
	if height < 5 {
		height = 5 // Minimum height
	}

	return width, height
}
