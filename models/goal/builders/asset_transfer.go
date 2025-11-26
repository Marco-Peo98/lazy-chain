package builders

import (
	"errors"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	goal "lazychain/models/goal/iface"
)

// AssetTransferBuilder implements asset transfer transaction (axfer)
type AssetTransferBuilder struct {
	assetID  string
	sender   string
	receiver string
	amount   string

	feeAuto         bool
	validRoundsAuto bool
	note            string

	status  string
	runFunc goal.RunFunc

	// Field editing state
	fieldsState *FieldsState
}

// NewAssetTransferBuilder creates a new AssetTransferBuilder
func NewAssetTransferBuilder() *AssetTransferBuilder {
	b := &AssetTransferBuilder{
		feeAuto:         true,
		validRoundsAuto: true,
	}
	b.initFields()
	return b
}

// initFields initializes the field definitions
func (b *AssetTransferBuilder) initFields() {
	fields := []FieldDef{
		// Required fields (first 4)
		{
			Key:         "assetID",
			Label:       "Asset ID",
			Placeholder: "[asset id - positive number]",
			Value:       &b.assetID,
			Required:    true,
			Type:        FieldTypeInteger,
		},
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
			Label:       "Amount",
			Placeholder: "[amount - positive number]",
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
			Placeholder: "[optional]",
			Value:       &b.note,
			Required:    false,
			Type:        FieldTypeText,
		},
	}
	b.fieldsState = NewFieldsState(fields)
}

func (b *AssetTransferBuilder) SetRunFunc(fn goal.RunFunc) { b.runFunc = fn }
func (b *AssetTransferBuilder) Title() string              { return "Asset Transfer (axfer)" }
func (b *AssetTransferBuilder) TxnType() string            { return "axfer" }
func (b *AssetTransferBuilder) Init() tea.Cmd              { return nil }

// IsEditing returns true if currently editing a field
func (b *AssetTransferBuilder) IsEditing() bool {
	return b.fieldsState != nil && b.fieldsState.Editing
}

// SetAvailableHeight configures the viewport based on available terminal lines
func (b *AssetTransferBuilder) SetAvailableHeight(lines int) {
	if b.fieldsState != nil {
		viewportHeight := CalculateViewportHeight(lines)
		b.fieldsState.SetViewportHeight(viewportHeight)
	}
}

func (b *AssetTransferBuilder) Validate() error {
	if strings.TrimSpace(b.assetID) == "" {
		return errors.New("asset ID is required")
	}
	if strings.TrimSpace(b.sender) == "" {
		return errors.New("sender is required")
	}
	if strings.TrimSpace(b.receiver) == "" {
		return errors.New("receiver is required")
	}
	if strings.TrimSpace(b.amount) == "" {
		return errors.New("amount is required")
	}

	if result := ValidateAssetID(b.assetID); !result.Valid {
		return fmt.Errorf("asset ID: %s", result.Message)
	}
	if result := ValidateAlgorandAddress(b.sender); !result.Valid {
		return fmt.Errorf("sender: %s", result.Message)
	}
	if result := ValidateAlgorandAddress(b.receiver); !result.Valid {
		return fmt.Errorf("receiver: %s", result.Message)
	}
	if result := ValidateAmount(b.amount); !result.Valid {
		return fmt.Errorf("amount: %s", result.Message)
	}

	if b.fieldsState != nil && b.fieldsState.HasValidationErrors() {
		return errors.New("please fix validation errors before executing")
	}

	return nil
}

func (b *AssetTransferBuilder) Args() []string {
	return []string{
		"asset", "send",
		"--assetid", strings.TrimSpace(b.assetID),
		"-f", strings.TrimSpace(b.sender),
		"-t", strings.TrimSpace(b.receiver),
		"-a", strings.TrimSpace(b.amount),
	}
}

func (b *AssetTransferBuilder) AfterRun(stdout, stderr string, runErr error) {
	if runErr != nil {
		b.status = fmt.Sprintf("Error: %v\n%s", runErr, stderr)
		return
	}
	b.status = stdout
}

func (b *AssetTransferBuilder) Update(msg tea.Msg) (goal.Builder, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()

		if b.fieldsState.Editing {
			switch key {
			case "enter":
				b.fieldsState.StopEditing(true)
			case "esc":
				b.fieldsState.StopEditing(false)
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
				if len(msg.Runes) > 0 {
					for _, r := range msg.Runes {
						b.fieldsState.InsertRune(r)
					}
				}
			}
		} else {
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
				b.fieldsState.ValidateAllFields()
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

func (b *AssetTransferBuilder) RenderFields() string {
	return RenderFieldsPanel("Asset Transfer", b.fieldsState, 4)
}

func (b *AssetTransferBuilder) RenderOutput() string {
	return RenderOutput(b.status)
}

func (b *AssetTransferBuilder) View() string {
	return b.RenderFields()
}
