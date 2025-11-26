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

	// Focus state
	focusOnFields bool // true = fields panel active, false = txn types panel active

	// Section dimensions (calculated dynamically)
	txnTypesWidth int
	fieldsWidth   int

	// Row height ratios
	mainRowRatio   int // Ratio for main content row (txn types + fields)
	outputRowRatio int // Ratio for output row

	// Calculated available heights (exposed for viewport configuration)
	availableFieldsHeight int
	availableOutputHeight int
}

// NewCmdGoalsLayout creates a new command goals layout
func NewCmdGoalsLayout(width, height int) *CmdGoalsLayout {
	l := &CmdGoalsLayout{
		BaseLayout: BaseLayout{
			Width:  width,
			Height: height,
		},
		txnTypesWidth:  30,
		fieldsWidth:    50,
		mainRowRatio:   3, // Main row takes 3 parts
		outputRowRatio: 1, // Output row takes 1 part
	}
	l.calculateDimensions()
	return l
}

// calculateDimensions calculates the available heights for each section
func (l *CmdGoalsLayout) calculateDimensions() {
	// Total available height (accounting for centering margins)
	usableHeight := l.Height - 4 // Some margin for centering

	// Calculate row heights based on ratios
	totalRatio := l.mainRowRatio + l.outputRowRatio
	mainRowHeight := (usableHeight * l.mainRowRatio) / totalRatio
	outputRowHeight := (usableHeight * l.outputRowRatio) / totalRatio

	// Fields panel: mainRowHeight minus border(2) + padding(2) + title area(3)
	// Border: 2 (top + bottom)
	// Padding: 2 (top + bottom)
	// We don't subtract title here because builder's RenderFieldsPanel handles it
	l.availableFieldsHeight = mainRowHeight - 6 // border + padding + some margin

	// Output panel height
	l.availableOutputHeight = outputRowHeight - 6

	// Ensure minimums
	if l.availableFieldsHeight < 10 {
		l.availableFieldsHeight = 10
	}
	if l.availableOutputHeight < 5 {
		l.availableOutputHeight = 5
	}
}

// GetAvailableFieldsHeight returns the calculated height available for fields
// This can be used to configure the builder's viewport
func (l *CmdGoalsLayout) GetAvailableFieldsHeight() int {
	return l.availableFieldsHeight
}

// GetAvailableOutputHeight returns the calculated height available for output
func (l *CmdGoalsLayout) GetAvailableOutputHeight() int {
	return l.availableOutputHeight
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

// SetFocus sets which panel is currently focused
func (l *CmdGoalsLayout) SetFocus(onFields bool) *CmdGoalsLayout {
	l.focusOnFields = onFields
	return l
}

// Build constructs the FlexBox for CmdGoalsView
func (l *CmdGoalsLayout) Build() *flexbox.FlexBox {
	// Recalculate dimensions in case they changed
	l.calculateDimensions()

	// Use actual terminal dimensions (with some margin for centering)
	totalWidth := l.Width - 4
	totalHeight := l.Height - 2

	// Ensure minimum dimensions
	if totalWidth < 80 {
		totalWidth = 80
	}
	if totalHeight < 25 {
		totalHeight = 25
	}

	// Create FlexBox with actual dimensions
	box := flexbox.New(totalWidth, totalHeight)

	// Row 1: Txn Types + Fields (main content)
	// Use ratios that give more space to fields
	txnTypesCell := flexbox.NewCell(1, l.mainRowRatio).
		SetMinWidth(l.txnTypesWidth).
		SetContentGenerator(func(maxX, maxY int) string {
			return l.renderTxnTypes(maxX, maxY)
		})

	fieldsCell := flexbox.NewCell(2, l.mainRowRatio).
		SetMinWidth(l.fieldsWidth).
		SetContentGenerator(func(maxX, maxY int) string {
			return l.renderFields(maxX, maxY)
		})

	row1 := box.NewRow().AddCells(txnTypesCell, fieldsCell)

	// Row 2: Output (full width)
	outputCell := flexbox.NewCell(1, l.outputRowRatio).
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

	// Title with focus indicator
	titleText := "Transaction Types"
	if !l.focusOnFields {
		titleText = "● " + titleText // Dot indicates active panel
	}
	title := TitleStyle().Render(titleText)
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

	// Add spacing and navigation hint
	content = append(content, "")
	if !l.focusOnFields {
		hint := lipgloss.NewStyle().
			Italic(true).
			Foreground(TextColorSecondary).
			Render("Tab/→: switch to fields")
		content = append(content, hint)
	}

	// Join content
	panelContent := strings.Join(content, "\n")

	// Calculate dimensions accounting for border and padding
	contentWidth := maxX - 4
	if contentWidth < 20 {
		contentWidth = 20
	}
	contentHeight := maxY - 2
	if contentHeight < 10 {
		contentHeight = 10
	}

	// Determine border color based on focus
	borderColor := BorderColorInactive
	if !l.focusOnFields {
		borderColor = BorderColorFocused
	}

	// Wrap in border
	return lipgloss.NewStyle().
		Width(contentWidth).
		Height(contentHeight).
		Padding(1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Render(panelContent)
}

// renderFields renders the fields panel
// NOTE: The fieldsContent from the builder already includes the title,
// so we don't add another title here to avoid duplication
func (l *CmdGoalsLayout) renderFields(maxX, maxY int) string {
	// Content directly from builder (already formatted with title)
	panelContent := l.fieldsContent
	if panelContent == "" {
		panelContent = "Select a transaction type to configure fields"
	}

	// Calculate dimensions accounting for border and padding
	contentWidth := maxX - 4
	if contentWidth < 30 {
		contentWidth = 30
	}
	contentHeight := maxY - 2
	if contentHeight < 10 {
		contentHeight = 10
	}

	// Determine border color based on focus
	borderColor := BorderColorInactive
	if l.focusOnFields {
		borderColor = BorderColorFocused
	}

	// Wrap in border
	return lipgloss.NewStyle().
		Width(contentWidth).
		Height(contentHeight).
		Padding(1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
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
		content = append(content, "Configure fields and press Ctrl+X to execute")
	}

	// Join content
	panelContent := strings.Join(content, "\n")

	// Calculate dimensions
	contentWidth := maxX - 4
	if contentWidth < 30 {
		contentWidth = 30
	}
	contentHeight := maxY - 2
	if contentHeight < 5 {
		contentHeight = 5
	}

	// Wrap in border
	return lipgloss.NewStyle().
		Width(contentWidth).
		Height(contentHeight).
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

// GetMinDimensions returns minimum required dimensions based on actual content
func (l *CmdGoalsLayout) GetMinDimensions() (width, height int) {
	// Width calculation:
	// - Transaction Types panel: txnTypesWidth + border(2) + padding(2) = +4
	// - Fields panel: fieldsWidth + border(2) + padding(2) = +4
	// - Gap between panels: 4
	// - External margins: 4
	minWidth := l.txnTypesWidth + l.fieldsWidth + 4 + 4 + 4 + 4 // = 96 with defaults

	// Height calculation:
	// - Main row needs space for:
	//   - Transaction Types: title(2) + 6 types + spacing(2) + border(2) = 12
	//   - Fields: title(2) + min 3 fields × 3 lines + instructions(2) + border(2) = 15
	//   - Use the larger: 15
	// - Output row: title(2) + content(3) + border(2) + padding(2) = 9
	// - External margins: 4
	const minViewportFields = 3 // Minimum fields visible (same as MinViewportHeight in builders)

	minTxnTypesHeight := 2 + len(l.txnTypes) + 2 + 2 // dynamic based on types
	if minTxnTypesHeight < 12 {
		minTxnTypesHeight = 12
	}
	minFieldsHeight := 2 + (minViewportFields * 3) + 4 + 2 // title + fields + instructions + border
	if minFieldsHeight < 17 {
		minFieldsHeight = 17
	}

	mainRowMinHeight := minFieldsHeight
	if minTxnTypesHeight > mainRowMinHeight {
		mainRowMinHeight = minTxnTypesHeight
	}

	outputMinHeight := 9
	minHeight := mainRowMinHeight + outputMinHeight + 4 // + margins

	// Ensure absolute minimums
	if minWidth < 100 {
		minWidth = 100
	}
	if minHeight < 35 {
		minHeight = 35
	}

	return minWidth, minHeight
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

	currentStyle := lipgloss.NewStyle().Foreground(BorderColorWarning)
	requiredStyle := lipgloss.NewStyle().Foreground(BorderColorSecondary)
	suggestionStyle := lipgloss.NewStyle().
		Italic(true).
		Foreground(TextColorSecondary)

	currentDims := currentStyle.Render(fmt.Sprintf("Current: %d × %d", l.Width, l.Height))
	requiredDims := requiredStyle.Render(fmt.Sprintf("Required: %d × %d", minWidth, minHeight))

	// Determine what needs to be increased
	needsWidth := minWidth - l.Width
	needsHeight := minHeight - l.Height

	var suggestion string
	if needsWidth > 0 && needsHeight > 0 {
		suggestion = fmt.Sprintf("Please increase width by %d and height by %d", needsWidth, needsHeight)
	} else if needsWidth > 0 {
		suggestion = fmt.Sprintf("Please increase width by %d", needsWidth)
	} else if needsHeight > 0 {
		suggestion = fmt.Sprintf("Please increase height by %d", needsHeight)
	} else {
		suggestion = "Please resize your terminal window"
	}

	message := lipgloss.JoinVertical(
		lipgloss.Center,
		errorTitle,
		"",
		currentDims,
		requiredDims,
		"",
		suggestionStyle.Render(suggestion),
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

// Update updates layout dimensions and recalculates available heights
func (l *CmdGoalsLayout) Update(width, height int) {
	l.Width = width
	l.Height = height
	l.calculateDimensions()
}
