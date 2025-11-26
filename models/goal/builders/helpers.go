package builders

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// =============================================================================
// STYLES (Catppuccin Mocha palette)
// =============================================================================

var (
	// Title style - purple
	StyleTitle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#cba6f7"))

	// Section headers
	StyleSectionRequired = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#89b4fa"))

	StyleSectionOptional = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#f9e2af"))

	// Labels and values
	StyleLabel = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#cdd6f4"))

	StyleValue = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#a6e3a1"))

	StylePlaceholder = lipgloss.NewStyle().
				Faint(true)

	// Status/output
	StyleSuccess = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#a6e3a1"))

	StyleError = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#f38ba8"))

	StyleMuted = lipgloss.NewStyle().
			Italic(true).
			Foreground(lipgloss.Color("#6c7086"))

	// Field editing styles
	StyleFieldActive = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#f5c2e7")).
				Bold(true)

	StyleFieldCursor = lipgloss.NewStyle().
				Background(lipgloss.Color("#f5c2e7")).
				Foreground(lipgloss.Color("#1e1e2e"))

	StyleFieldSelected = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#f9e2af"))

	StyleFieldEditing = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#f5c2e7"))

	// Validation styles
	StyleValidationError = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#f38ba8")).
				Italic(true)

	StyleFieldInvalid = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#f38ba8"))

	StyleFieldValid = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#a6e3a1"))

	// Scroll indicator styles
	StyleScrollIndicator = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#74c7ec")).
				Italic(true)
)

// =============================================================================
// BASIC HELPER FUNCTIONS
// =============================================================================

// GetPlaceholder returns the value if set, otherwise a styled placeholder.
func GetPlaceholder(value, placeholder string) string {
	if strings.TrimSpace(value) == "" {
		return StylePlaceholder.Render(placeholder)
	}
	return value
}

// GetBoolStr returns a styled yes/no indicator.
func GetBoolStr(value bool) string {
	if value {
		return "✓ Yes"
	}
	return "✗ No"
}

// RenderOutput renders status output with appropriate styling.
func RenderOutput(status string) string {
	if status == "" {
		return StyleMuted.Render("No output yet.")
	}
	if strings.HasPrefix(status, "Error") || strings.HasPrefix(status, "Validation Error") {
		return StyleError.Render(status)
	}
	return StyleSuccess.Render(status)
}

// RenderFieldRow renders a labeled field row (for non-editable display).
func RenderFieldRow(label, value, placeholder string) string {
	l := StyleLabel.Render(label + ":")
	v := StyleValue.Render(GetPlaceholder(value, placeholder))
	return "  " + l + "\n    " + v
}

// RenderBoolRow renders a labeled boolean field (for non-editable display).
func RenderBoolRow(label string, value bool) string {
	l := StyleLabel.Render(label + ":")
	v := StyleValue.Render(GetBoolStr(value))
	return "  " + l + "\n    " + v
}

// =============================================================================
// FIELD DEFINITIONS
// =============================================================================

// FieldType represents the type of input expected
type FieldType int

const (
	FieldTypeText FieldType = iota
	FieldTypeAddress
	FieldTypeAmount
	FieldTypeInteger
	FieldTypeBool
	FieldTypeSelect // For dropdowns like action: freeze/unfreeze
)

// FieldDef defines a single editable field
type FieldDef struct {
	Key             string    // Internal identifier (e.g., "sender")
	Label           string    // Display label (e.g., "Sender Address")
	Placeholder     string    // Placeholder text (e.g., "[address]")
	Value           *string   // Pointer to the actual value in builder
	Required        bool      // Is this field required?
	Type            FieldType // Type of field for validation hints
	Options         []string  // For FieldTypeSelect: available options
	BoolValue       *bool     // For FieldTypeBool: pointer to bool
	ValidationError string    // Current validation error message (empty = valid)
}

// HasError returns true if the field has a validation error
func (f *FieldDef) HasError() bool {
	return f.ValidationError != ""
}

// ClearError clears the validation error
func (f *FieldDef) ClearError() {
	f.ValidationError = ""
}

// =============================================================================
// VIEWPORT CONFIGURATION
// =============================================================================

const (
	// DefaultViewportHeight is the default number of fields visible at once
	DefaultViewportHeight = 6

	// MinViewportHeight is the minimum viewport height
	MinViewportHeight = 3

	// MaxViewportHeight is the maximum viewport height
	MaxViewportHeight = 15

	// LinesPerField is the approximate number of terminal lines per field
	// (label line + value line + spacing + potential error line)
	LinesPerField = 3

	// FixedOverheadLines is the number of lines used by fixed elements
	// (title, section headers, instructions, scroll indicators, field counter)
	FixedOverheadLines = 12
)

// CalculateViewportHeight calculates the optimal viewport height (number of fields)
// based on available terminal lines
func CalculateViewportHeight(availableLines int) int {
	// Subtract fixed overhead
	usableLines := availableLines - FixedOverheadLines

	// Calculate how many fields can fit
	fields := usableLines / LinesPerField

	// Clamp to valid range
	if fields < MinViewportHeight {
		fields = MinViewportHeight
	}
	if fields > MaxViewportHeight {
		fields = MaxViewportHeight
	}

	return fields
}

// =============================================================================
// FIELDS STATE MANAGEMENT (with Viewport)
// =============================================================================

// FieldsState manages the editing state for a builder's fields
type FieldsState struct {
	Fields    []FieldDef
	Cursor    int    // Which field is selected
	Editing   bool   // Are we in edit mode for current field?
	CursorPos int    // Cursor position within the text being edited
	TempValue string // Temporary value while editing

	// Viewport state
	ViewportHeight int // Number of fields visible at once (0 = show all)
	ViewportOffset int // First visible field index
}

// NewFieldsState creates a new fields state with default viewport
func NewFieldsState(fields []FieldDef) *FieldsState {
	return &FieldsState{
		Fields:         fields,
		Cursor:         0,
		Editing:        false,
		CursorPos:      0,
		TempValue:      "",
		ViewportHeight: DefaultViewportHeight,
		ViewportOffset: 0,
	}
}

// NewFieldsStateWithViewport creates a new fields state with custom viewport height
func NewFieldsStateWithViewport(fields []FieldDef, viewportHeight int) *FieldsState {
	if viewportHeight < MinViewportHeight {
		viewportHeight = MinViewportHeight
	}
	if viewportHeight > MaxViewportHeight {
		viewportHeight = MaxViewportHeight
	}

	return &FieldsState{
		Fields:         fields,
		Cursor:         0,
		Editing:        false,
		CursorPos:      0,
		TempValue:      "",
		ViewportHeight: viewportHeight,
		ViewportOffset: 0,
	}
}

// SetViewportHeight sets the viewport height dynamically
func (fs *FieldsState) SetViewportHeight(height int) {
	if height < MinViewportHeight {
		height = MinViewportHeight
	}
	if height > MaxViewportHeight {
		height = MaxViewportHeight
	}
	fs.ViewportHeight = height
	fs.ensureCursorVisible()
}

// CurrentField returns the currently selected field
func (fs *FieldsState) CurrentField() *FieldDef {
	if fs.Cursor >= 0 && fs.Cursor < len(fs.Fields) {
		return &fs.Fields[fs.Cursor]
	}
	return nil
}

// ensureCursorVisible adjusts viewport to keep cursor in view
func (fs *FieldsState) ensureCursorVisible() {
	if fs.ViewportHeight <= 0 || fs.ViewportHeight >= len(fs.Fields) {
		// No scrolling needed - all fields visible
		fs.ViewportOffset = 0
		return
	}

	// Cursor above viewport
	if fs.Cursor < fs.ViewportOffset {
		fs.ViewportOffset = fs.Cursor
	}

	// Cursor below viewport
	if fs.Cursor >= fs.ViewportOffset+fs.ViewportHeight {
		fs.ViewportOffset = fs.Cursor - fs.ViewportHeight + 1
	}

	// Clamp viewport offset
	maxOffset := len(fs.Fields) - fs.ViewportHeight
	if maxOffset < 0 {
		maxOffset = 0
	}
	if fs.ViewportOffset > maxOffset {
		fs.ViewportOffset = maxOffset
	}
	if fs.ViewportOffset < 0 {
		fs.ViewportOffset = 0
	}
}

// MoveUp moves cursor up
func (fs *FieldsState) MoveUp() {
	if !fs.Editing && fs.Cursor > 0 {
		fs.Cursor--
		fs.ensureCursorVisible()
	}
}

// MoveDown moves cursor down
func (fs *FieldsState) MoveDown() {
	if !fs.Editing && fs.Cursor < len(fs.Fields)-1 {
		fs.Cursor++
		fs.ensureCursorVisible()
	}
}

// PageUp moves cursor up by viewport height
func (fs *FieldsState) PageUp() {
	if !fs.Editing {
		fs.Cursor -= fs.ViewportHeight
		if fs.Cursor < 0 {
			fs.Cursor = 0
		}
		fs.ensureCursorVisible()
	}
}

// PageDown moves cursor down by viewport height
func (fs *FieldsState) PageDown() {
	if !fs.Editing {
		fs.Cursor += fs.ViewportHeight
		if fs.Cursor >= len(fs.Fields) {
			fs.Cursor = len(fs.Fields) - 1
		}
		fs.ensureCursorVisible()
	}
}

// GoToFirst moves cursor to first field
func (fs *FieldsState) GoToFirst() {
	if !fs.Editing {
		fs.Cursor = 0
		fs.ensureCursorVisible()
	}
}

// GoToLast moves cursor to last field
func (fs *FieldsState) GoToLast() {
	if !fs.Editing {
		fs.Cursor = len(fs.Fields) - 1
		fs.ensureCursorVisible()
	}
}

// GetVisibleRange returns the range of visible field indices
func (fs *FieldsState) GetVisibleRange() (start, end int) {
	if fs.ViewportHeight <= 0 || fs.ViewportHeight >= len(fs.Fields) {
		// All fields visible
		return 0, len(fs.Fields)
	}

	start = fs.ViewportOffset
	end = fs.ViewportOffset + fs.ViewportHeight

	if end > len(fs.Fields) {
		end = len(fs.Fields)
	}

	return start, end
}

// HasFieldsAbove returns true if there are hidden fields above viewport
func (fs *FieldsState) HasFieldsAbove() bool {
	return fs.ViewportOffset > 0
}

// HasFieldsBelow returns true if there are hidden fields below viewport
func (fs *FieldsState) HasFieldsBelow() bool {
	if fs.ViewportHeight <= 0 || fs.ViewportHeight >= len(fs.Fields) {
		return false
	}
	return fs.ViewportOffset+fs.ViewportHeight < len(fs.Fields)
}

// FieldsAboveCount returns the number of hidden fields above viewport
func (fs *FieldsState) FieldsAboveCount() int {
	return fs.ViewportOffset
}

// FieldsBelowCount returns the number of hidden fields below viewport
func (fs *FieldsState) FieldsBelowCount() int {
	if fs.ViewportHeight <= 0 || fs.ViewportHeight >= len(fs.Fields) {
		return 0
	}
	remaining := len(fs.Fields) - (fs.ViewportOffset + fs.ViewportHeight)
	if remaining < 0 {
		return 0
	}
	return remaining
}

// StartEditing enters edit mode for current field
func (fs *FieldsState) StartEditing() {
	field := fs.CurrentField()
	if field == nil {
		return
	}

	// Bool fields toggle directly
	if field.Type == FieldTypeBool && field.BoolValue != nil {
		*field.BoolValue = !*field.BoolValue
		return
	}

	// Select fields cycle through options
	if field.Type == FieldTypeSelect && len(field.Options) > 0 && field.Value != nil {
		currentIdx := 0
		for i, opt := range field.Options {
			if opt == *field.Value {
				currentIdx = i
				break
			}
		}
		nextIdx := (currentIdx + 1) % len(field.Options)
		*field.Value = field.Options[nextIdx]
		return
	}

	// Text fields enter edit mode
	if field.Value != nil {
		fs.Editing = true
		fs.TempValue = *field.Value
		fs.CursorPos = len(fs.TempValue)
		// Clear any previous validation error when starting to edit
		field.ClearError()
	}
}

// StopEditing exits edit mode, optionally saving and validating
func (fs *FieldsState) StopEditing(save bool) {
	if !fs.Editing {
		return
	}
	field := fs.CurrentField()
	if field == nil {
		fs.Editing = false
		return
	}

	if save && field.Value != nil {
		*field.Value = fs.TempValue

		// Validate the field after saving
		result := ValidateField(field)
		if !result.Valid {
			field.ValidationError = result.Message
		} else {
			field.ValidationError = ""
		}
	}

	fs.Editing = false
	fs.TempValue = ""
	fs.CursorPos = 0
}

// InsertRune adds a character at cursor position
func (fs *FieldsState) InsertRune(r rune) {
	if !fs.Editing {
		return
	}
	left := fs.TempValue[:fs.CursorPos]
	right := fs.TempValue[fs.CursorPos:]
	fs.TempValue = left + string(r) + right
	fs.CursorPos++
}

// Backspace removes character before cursor
func (fs *FieldsState) Backspace() {
	if !fs.Editing || fs.CursorPos == 0 {
		return
	}
	left := fs.TempValue[:fs.CursorPos-1]
	right := fs.TempValue[fs.CursorPos:]
	fs.TempValue = left + right
	fs.CursorPos--
}

// Delete removes character at cursor
func (fs *FieldsState) Delete() {
	if !fs.Editing || fs.CursorPos >= len(fs.TempValue) {
		return
	}
	left := fs.TempValue[:fs.CursorPos]
	right := fs.TempValue[fs.CursorPos+1:]
	fs.TempValue = left + right
}

// MoveCursorLeft moves text cursor left
func (fs *FieldsState) MoveCursorLeft() {
	if fs.Editing && fs.CursorPos > 0 {
		fs.CursorPos--
	}
}

// MoveCursorRight moves text cursor right
func (fs *FieldsState) MoveCursorRight() {
	if fs.Editing && fs.CursorPos < len(fs.TempValue) {
		fs.CursorPos++
	}
}

// ClearField clears the current field value
func (fs *FieldsState) ClearField() {
	if fs.Editing {
		fs.TempValue = ""
		fs.CursorPos = 0
	}
}

// ValidateAllFields validates all fields and sets their error states
func (fs *FieldsState) ValidateAllFields() []string {
	var errors []string
	for i := range fs.Fields {
		field := &fs.Fields[i]
		result := ValidateField(field)
		if !result.Valid {
			field.ValidationError = result.Message
			errors = append(errors, field.Label+": "+result.Message)
		} else {
			field.ValidationError = ""
		}
	}
	return errors
}

// HasValidationErrors returns true if any field has a validation error
func (fs *FieldsState) HasValidationErrors() bool {
	for i := range fs.Fields {
		if fs.Fields[i].HasError() {
			return true
		}
	}
	return false
}

// GetErrorCount returns the number of fields with validation errors
func (fs *FieldsState) GetErrorCount() int {
	count := 0
	for i := range fs.Fields {
		if fs.Fields[i].HasError() {
			count++
		}
	}
	return count
}

// =============================================================================
// RENDERING FUNCTIONS FOR EDITABLE FIELDS
// =============================================================================

// RenderEditableField renders a field with edit state awareness
func RenderEditableField(field FieldDef, isSelected bool, isEditing bool, tempValue string, cursorPos int) string {
	// Determine if field has error
	hasError := field.HasError()

	// Label styling
	labelStyle := StyleLabel
	if hasError {
		labelStyle = StyleFieldInvalid
	} else if isSelected {
		labelStyle = StyleFieldSelected
	}
	if isEditing {
		labelStyle = StyleFieldEditing
	}

	// Add required indicator
	labelText := field.Label
	if field.Required {
		labelText += " *"
	}
	label := labelStyle.Render(labelText + ":")

	// Value rendering
	var valueStr string

	if field.Type == FieldTypeBool && field.BoolValue != nil {
		valueStr = GetBoolStr(*field.BoolValue)
		if isSelected {
			valueStr = StyleFieldSelected.Render(valueStr + "  ← Enter to toggle")
		}
	} else if field.Type == FieldTypeSelect && field.Value != nil {
		valueStr = *field.Value
		if valueStr == "" && len(field.Options) > 0 {
			valueStr = field.Options[0]
		}
		if isSelected {
			valueStr = StyleFieldSelected.Render(valueStr + "  ← Enter to cycle")
		}
	} else if isEditing {
		// Show text with cursor
		valueStr = renderTextWithCursor(tempValue, cursorPos)
	} else if field.Value != nil && strings.TrimSpace(*field.Value) != "" {
		// Show value with validation color
		if hasError {
			valueStr = StyleFieldInvalid.Render(*field.Value)
		} else {
			valueStr = StyleFieldValid.Render(*field.Value)
		}
	} else {
		valueStr = StylePlaceholder.Render(field.Placeholder)
	}

	// Selection indicator
	prefix := "  "
	if isSelected && !isEditing {
		prefix = "▶ "
	} else if isEditing {
		prefix = "✎ "
	} else if hasError {
		prefix = "⚠ "
	}

	// Build the field output
	output := prefix + label + "\n    " + valueStr

	// Add validation error message if present
	if hasError && !isEditing {
		errorMsg := StyleValidationError.Render("    ↳ " + field.ValidationError)
		output += "\n" + errorMsg
	}

	return output
}

// renderTextWithCursor renders text with a visible cursor
func renderTextWithCursor(text string, pos int) string {
	if pos > len(text) {
		pos = len(text)
	}

	// Handle empty text
	if len(text) == 0 {
		return StyleFieldCursor.Render(" ")
	}

	left := text[:pos]
	var cursor string
	var right string

	if pos < len(text) {
		cursor = StyleFieldCursor.Render(string(text[pos]))
		right = text[pos+1:]
	} else {
		cursor = StyleFieldCursor.Render(" ")
		right = ""
	}

	return StyleValue.Render(left) + cursor + StyleValue.Render(right)
}

// renderScrollIndicator renders a scroll indicator
func renderScrollIndicator(direction string, count int) string {
	var arrow string
	var text string

	if direction == "up" {
		arrow = "↑"
		if count == 1 {
			text = fmt.Sprintf("%s 1 more field above", arrow)
		} else {
			text = fmt.Sprintf("%s %d more fields above", arrow, count)
		}
	} else {
		arrow = "↓"
		if count == 1 {
			text = fmt.Sprintf("%s 1 more field below", arrow)
		} else {
			text = fmt.Sprintf("%s %d more fields below", arrow, count)
		}
	}

	return StyleScrollIndicator.Render(text)
}

// RenderFieldsPanel renders all fields with proper state and viewport
func RenderFieldsPanel(title string, fs *FieldsState, requiredCount int) string {
	var lines []string

	lines = append(lines, StyleTitle.Render(title), "")

	// Count validation errors
	errorCount := fs.GetErrorCount()

	// Show error summary if there are errors
	if errorCount > 0 {
		errText := fmt.Sprintf("⚠ %d validation error", errorCount)
		if errorCount > 1 {
			errText += "s"
		}
		errSummary := StyleValidationError.Render(errText)
		lines = append(lines, errSummary, "")
	}

	// Get visible range
	visStart, visEnd := fs.GetVisibleRange()

	// Show scroll indicator if there are fields above
	if fs.HasFieldsAbove() {
		lines = append(lines, renderScrollIndicator("up", fs.FieldsAboveCount()), "")
	}

	// Determine which section headers to show based on visible range
	showRequiredHeader := false
	showOptionalHeader := false

	for i := visStart; i < visEnd; i++ {
		if i < requiredCount {
			showRequiredHeader = true
		} else {
			showOptionalHeader = true
		}
	}

	// Track if we need to show section headers
	inRequiredSection := false
	inOptionalSection := false

	// Render visible fields
	for i := visStart; i < visEnd; i++ {
		// Check if we need a section header
		if i < requiredCount && !inRequiredSection {
			if showRequiredHeader {
				lines = append(lines, StyleSectionRequired.Render("Required Fields:"), "")
			}
			inRequiredSection = true
		} else if i >= requiredCount && !inOptionalSection {
			if showOptionalHeader {
				lines = append(lines, StyleSectionOptional.Render("Optional Fields:"), "")
			}
			inOptionalSection = true
		}

		// Render the field
		isSelected := fs.Cursor == i
		isEditing := isSelected && fs.Editing
		rendered := RenderEditableField(fs.Fields[i], isSelected, isEditing, fs.TempValue, fs.CursorPos)
		lines = append(lines, rendered, "")
	}

	// Show scroll indicator if there are fields below
	if fs.HasFieldsBelow() {
		lines = append(lines, renderScrollIndicator("down", fs.FieldsBelowCount()), "")
	}

	// Field counter
	fieldCounter := StyleMuted.Render(fmt.Sprintf("Field %d/%d", fs.Cursor+1, len(fs.Fields)))
	lines = append(lines, fieldCounter)

	// Instructions
	lines = append(lines, "") // spacing
	if fs.Editing {
		lines = append(lines, StyleMuted.Render("Enter: confirm | Esc: cancel | ←→: move | Ctrl+U: clear"))
	} else {
		if fs.ViewportHeight > 0 && fs.ViewportHeight < len(fs.Fields) {
			// Show pagination controls when scrolling is active
			lines = append(lines, StyleMuted.Render("↑↓: navigate | PgUp/PgDn: page | Enter: edit | Ctrl+X: execute"))
		} else {
			lines = append(lines, StyleMuted.Render("↑↓: navigate | Enter: edit | Ctrl+X: execute"))
		}
	}

	return strings.Join(lines, "\n")
}
