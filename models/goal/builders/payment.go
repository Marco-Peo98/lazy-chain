package builders

import (
	"errors"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	goal "lazychain/models/goal/iface"
)

// PaymentBuilder implements payment transaction (pay)
type PaymentBuilder struct {
	// Field values
	sender   string
	receiver string
	amount   string

	// Optional fields
	feeAuto         bool
	validRoundsAuto bool
	note            string

	// State
	status  string
	runFunc goal.RunFunc

	// Field editing state
	fieldsState *FieldsState
}

// NewPaymentBuilder creates a new PaymentBuilder with default values
func NewPaymentBuilder() *PaymentBuilder {
	b := &PaymentBuilder{
		feeAuto:         true,
		validRoundsAuto: true,
	}
	b.initFields()
	return b
}

// initFields initializes the field definitions
func (b *PaymentBuilder) initFields() {
	fields := []FieldDef{
		// Required fields (first 3)
		{
			Key:         "sender",
			Label:       "Sender",
			Placeholder: "[sender address - 58 chars]",
			Value:       &b.sender,
			Required:    true,
			Type:        FieldTypeAddress,
		},
		{
			Key:         "receiver",
			Label:       "Receiver",
			Placeholder: "[receiver address - 58 chars]",
			Value:       &b.receiver,
			Required:    true,
			Type:        FieldTypeAddress,
		},
		{
			Key:         "amount",
			Label:       "Amount (μAlgos)",
			Placeholder: "[amount in microAlgos]",
			Value:       &b.amount,
			Required:    true,
			Type:        FieldTypeAmount,
		},
		// Optional fields
		{
			Key:         "feeAuto",
			Label:       "Auto Fee",
			Placeholder: "",
			BoolValue:   &b.feeAuto,
			Required:    false,
			Type:        FieldTypeBool,
		},
		{
			Key:         "validRoundsAuto",
			Label:       "Auto Valid Rounds",
			Placeholder: "",
			BoolValue:   &b.validRoundsAuto,
			Required:    false,
			Type:        FieldTypeBool,
		},
		{
			Key:         "note",
			Label:       "Note",
			Placeholder: "[optional note]",
			Value:       &b.note,
			Required:    false,
			Type:        FieldTypeText,
		},
	}
	b.fieldsState = NewFieldsState(fields)
}

// SetRunFunc sets the run function (called by GOALModel)
func (b *PaymentBuilder) SetRunFunc(fn goal.RunFunc) {
	b.runFunc = fn
}

// Title returns the builder title
func (b *PaymentBuilder) Title() string {
	return "Payment (pay)"
}

// TxnType returns the transaction type identifier
func (b *PaymentBuilder) TxnType() string {
	return "pay"
}

// Init initializes the builder
func (b *PaymentBuilder) Init() tea.Cmd {
	return nil
}

// IsEditing returns true if currently editing a field
func (b *PaymentBuilder) IsEditing() bool {
	return b.fieldsState != nil && b.fieldsState.Editing
}

// SetAvailableHeight configures the viewport based on available terminal lines
func (b *PaymentBuilder) SetAvailableHeight(lines int) {
	if b.fieldsState != nil {
		viewportHeight := CalculateViewportHeight(lines)
		b.fieldsState.SetViewportHeight(viewportHeight)
	}
}

// Validate checks if required fields are filled and validates their format
func (b *PaymentBuilder) Validate() error {
	// Check required fields
	if strings.TrimSpace(b.sender) == "" {
		return errors.New("sender is required")
	}
	if strings.TrimSpace(b.receiver) == "" {
		return errors.New("receiver is required")
	}
	if strings.TrimSpace(b.amount) == "" {
		return errors.New("amount is required")
	}

	// Validate field formats
	if result := ValidateAlgorandAddress(b.sender); !result.Valid {
		return fmt.Errorf("sender: %s", result.Message)
	}
	if result := ValidateAlgorandAddress(b.receiver); !result.Valid {
		return fmt.Errorf("receiver: %s", result.Message)
	}
	if result := ValidateAmount(b.amount); !result.Valid {
		return fmt.Errorf("amount: %s", result.Message)
	}

	// Check for validation errors in fieldsState
	if b.fieldsState != nil && b.fieldsState.HasValidationErrors() {
		return errors.New("please fix validation errors before executing")
	}

	return nil
}

// Args builds the command arguments
func (b *PaymentBuilder) Args() []string {
	var args []string
	args = append(args, "clerk", "send")
	args = append(args, "-f", strings.TrimSpace(b.sender))
	args = append(args, "-t", strings.TrimSpace(b.receiver))
	args = append(args, "-a", strings.TrimSpace(b.amount))

	if strings.TrimSpace(b.note) != "" {
		args = append(args, "-n", b.note)
	}

	return args
}

// AfterRun is called after command execution
func (b *PaymentBuilder) AfterRun(stdout, stderr string, runErr error) {
	if runErr != nil {
		b.status = fmt.Sprintf("Error: %v\n%s", runErr, strings.TrimSpace(stderr))
		return
	}
	b.status = strings.TrimSpace(stdout)
}

// Update handles user input
func (b *PaymentBuilder) Update(msg tea.Msg) (goal.Builder, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()

		if b.fieldsState.Editing {
			// In editing mode
			switch key {
			case "enter":
				b.fieldsState.StopEditing(true) // Save and validate
			case "esc":
				b.fieldsState.StopEditing(false) // Cancel
			case "backspace":
				b.fieldsState.Backspace()
			case "delete":
				b.fieldsState.Delete()
			case "left":
				b.fieldsState.MoveCursorLeft()
			case "right":
				b.fieldsState.MoveCursorRight()
			case "home", "ctrl+a":
				b.fieldsState.CursorPos = 0
			case "end", "ctrl+e":
				b.fieldsState.CursorPos = len(b.fieldsState.TempValue)
			case "ctrl+u":
				b.fieldsState.ClearField()
			default:
				// Handle typed runes
				if len(msg.Runes) > 0 {
					for _, r := range msg.Runes {
						b.fieldsState.InsertRune(r)
					}
				}
			}
		} else {
			// Navigation mode
			switch key {
			case "up", "k":
				b.fieldsState.MoveUp()
			case "down", "j":
				b.fieldsState.MoveDown()
			case "pgup":
				b.fieldsState.PageUp()
			case "pgdown":
				b.fieldsState.PageDown()
			case "home":
				b.fieldsState.GoToFirst()
			case "end":
				b.fieldsState.GoToLast()
			case "enter", " ":
				b.fieldsState.StartEditing()
			case "ctrl+x":
				// Validate all fields before execution
				b.fieldsState.ValidateAllFields()

				// Execute transaction
				if err := b.Validate(); err != nil {
					b.status = "Validation Error: " + err.Error()
					return b, nil
				}
				if b.runFunc != nil {
					b.runFunc(b.Args())
				}
			}
		}
	}
	return b, nil
}

// RenderFields renders the fields panel
func (b *PaymentBuilder) RenderFields() string {
	return RenderFieldsPanel("Payment Transaction", b.fieldsState, 3) // 3 required fields
}

// RenderOutput renders the output panel
func (b *PaymentBuilder) RenderOutput() string {
	return RenderOutput(b.status)
}

// View returns the complete view (delegates to RenderFields for compatibility)
func (b *PaymentBuilder) View() string {
	return b.RenderFields()
}
