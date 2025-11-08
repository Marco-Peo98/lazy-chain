package layout

import (
	"fmt"
	"strings"

	"github.com/76creates/stickers/flexbox"
	"github.com/charmbracelet/lipgloss"
)

// CmdGoalsLayout manages the layout for command goals view
type CmdGoalsLayout struct {
	BaseLayout

	// Content for each section (provided by GOALModel)
	txnTypes      []string
	selectedType  int
	fieldsContent string
	outputContent string

	// Section dimensions (calculated)
	txnTypesWidth int
	fieldsWidth   int
	outputHeight  int
}

// NewCmdGoalsLayout creates a new command goals layout
func NewCmdGoalsLayout(width, height int) *CmdGoalsLayout {
	return &CmdGoalsLayout{
		BaseLayout: BaseLayout{
			Width:  width,
			Height: height,
		},
		txnTypesWidth: 30,
		fieldsWidth:   45,
		outputHeight:  10,
	}
}

// SetTxnTypes sets the transaction types list
func (l *CmdGoalsLayout) SetTxnTypes(types []string, selected int) *CmdGoalsLayout {
	l.txnTypes = types
	l.selectedType = selected
	return l
}

// SetFields sets the fields content
func (l *CmdGoalsLayout) SetFields(content string) *CmdGoalsLayout {
	l.fieldsContent = content
	return l
}

// SetOutput sets the output content
func (l *CmdGoalsLayout) SetOutput(content string) *CmdGoalsLayout {
	l.outputContent = content
	return l
}

// Build constructs the FlexBox for CmdGoalsView
func (l *CmdGoalsLayout) Build() *flexbox.FlexBox {
	// Calculate total dimensions
	totalWidth := l.txnTypesWidth + l.fieldsWidth + 4
	totalHeight := 35

	// Create FlexBox
	box := flexbox.New(totalWidth, totalHeight)

	// Row 1: Txn Types + Fields (main content)
	txnTypesCell := flexbox.NewCell(2, 3).
		SetMinWidth(l.txnTypesWidth).
		SetContentGenerator(func(maxX, maxY int) string {
			return l.renderTxnTypes(maxX, maxY)
		})

	fieldsCell := flexbox.NewCell(3, 3).
		SetMinWidth(l.fieldsWidth).
		SetContentGenerator(func(maxX, maxY int) string {
			return l.renderFields(maxX, maxY)
		})

	row1 := box.NewRow().AddCells(txnTypesCell, fieldsCell)

	// Row 2: Output (full width)
	outputCell := flexbox.NewCell(1, 2).
		SetContentGenerator(func(maxX, maxY int) string {
			return l.renderOutput(maxX, maxY)
		})

	row2 := box.NewRow().AddCells(outputCell)

	// Add rows to box
	box.AddRows([]*flexbox.Row{row1, row2})

	return box
}

// renderTxnTypes renders the transaction types panel
func (l *CmdGoalsLayout) renderTxnTypes(maxX, maxY int) string {
	var content []string

	// Title
	title := TitleStyle().Render("Transaction Types")
	content = append(content, title)
	content = append(content, "")

	// Transaction types list
	for i, txnType := range l.txnTypes {
		cursor := "  "
		style := lipgloss.NewStyle().Foreground(TextColorPrimary)

		if i == l.selectedType {
			cursor = "> "
			style = style.Foreground(TextColorHighlight).Bold(true)
		}

		line := cursor + style.Render(txnType)
		content = append(content, line)
	}

	// Join content
	panelContent := strings.Join(content, "\n")

	// Wrap in border
	return lipgloss.NewStyle().
		Width(maxX - 4).
		Height(maxY - 2).
		Padding(1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(BorderColorPrimary).
		Render(panelContent)
}

// renderFields renders the fields panel
func (l *CmdGoalsLayout) renderFields(maxX, maxY int) string {
	var content []string

	// Title
	title := TitleStyle().Render("Transaction Fields")
	content = append(content, title)
	content = append(content, "")

	// Content (from builder)
	if l.fieldsContent != "" {
		content = append(content, l.fieldsContent)
	} else {
		content = append(content, "Select a transaction type to configure fields")
	}

	// Join content
	panelContent := strings.Join(content, "\n")

	// Wrap in border
	return lipgloss.NewStyle().
		Width(maxX - 4).
		Height(maxY - 2).
		Padding(1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(BorderColorSecondary).
		Render(panelContent)
}

// renderOutput renders the output panel
func (l *CmdGoalsLayout) renderOutput(maxX, maxY int) string {
	var content []string

	// Title
	title := TitleStyle().Render("Output / Command Result")
	content = append(content, title)
	content = append(content, "")

	// Content
	if l.outputContent != "" {
		lines := strings.Split(l.outputContent, "\n")
		content = append(content, lines...)
	} else {
		content = append(content, "Command output will appear here...")
		content = append(content, "")
		content = append(content, "Configure fields and press Enter to execute")
	}

	// Join content
	panelContent := strings.Join(content, "\n")

	// Wrap in border
	return lipgloss.NewStyle().
		Width(maxX - 4).
		Height(maxY - 2).
		Padding(1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(BorderColorAccent).
		Render(panelContent)
}

// Render renders the complete layout centered
func (l *CmdGoalsLayout) Render() string {
	box := l.Build()
	boxRendered := box.Render()

	return lipgloss.Place(
		l.Width,
		l.Height,
		lipgloss.Center,
		lipgloss.Center,
		boxRendered,
	)
}

// GetMinDimensions returns minimum required dimensions
func (l *CmdGoalsLayout) GetMinDimensions() (width, height int) {
	return 80, 30
}

// IsValid checks if dimensions are sufficient
func (l *CmdGoalsLayout) IsValid() bool {
	minWidth, minHeight := l.GetMinDimensions()
	return l.Width >= minWidth && l.Height >= minHeight
}

// RenderError shows error when terminal too small
func (l *CmdGoalsLayout) RenderError() string {
	minWidth, minHeight := l.GetMinDimensions()

	errorTitle := lipgloss.NewStyle().
		Bold(true).
		Foreground(BorderColorError).
		Render("Terminal Too Small")

	currentDims := lipgloss.NewStyle().
		Foreground(BorderColorWarning).
		Render(lipgloss.JoinVertical(
			lipgloss.Left,
			"Current dimensions:",
			lipgloss.NewStyle().Render("  "+lipgloss.NewStyle().Bold(true).Render(lipgloss.JoinHorizontal(lipgloss.Left, "Width: ", lipgloss.NewStyle().Render(fmt.Sprintf("%d", l.Width))))),
			lipgloss.NewStyle().Render("  "+lipgloss.NewStyle().Bold(true).Render(lipgloss.JoinHorizontal(lipgloss.Left, "Height: ", lipgloss.NewStyle().Render(fmt.Sprintf("%d", l.Height))))),
		))

	requiredDims := lipgloss.NewStyle().
		Foreground(BorderColorSecondary).
		Render(lipgloss.JoinVertical(
			lipgloss.Left,
			"Required dimensions:",
			lipgloss.NewStyle().Render("  "+lipgloss.NewStyle().Bold(true).Render(lipgloss.JoinHorizontal(lipgloss.Left, "Width: ", lipgloss.NewStyle().Render(fmt.Sprintf("%d", minWidth))))),
			lipgloss.NewStyle().Render("  "+lipgloss.NewStyle().Bold(true).Render(lipgloss.JoinHorizontal(lipgloss.Left, "Height: ", lipgloss.NewStyle().Render(fmt.Sprintf("%d", minHeight))))),
		))

	message := lipgloss.JoinVertical(
		lipgloss.Center,
		errorTitle,
		"",
		currentDims,
		"",
		requiredDims,
		"",
		lipgloss.NewStyle().
			Italic(true).
			Foreground(TextColorSecondary).
			Render("Please resize your terminal window"),
	)

	errorContainer := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(BorderColorError).
		Padding(1, 2).
		Render(message)

	return lipgloss.Place(
		l.Width,
		l.Height,
		lipgloss.Center,
		lipgloss.Center,
		errorContainer,
	)
}

// Update updates layout dimensions
func (l *CmdGoalsLayout) Update(width, height int) {
	l.Width = width
	l.Height = height
}
