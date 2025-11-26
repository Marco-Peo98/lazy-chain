package builders

import (
	"errors"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	goal "lazychain/models/goal/iface"
)

// AssetFreezeBuilder implements asset freeze transaction (afrz)
type AssetFreezeBuilder struct {
	assetID      string
	sender       string
	freezeTarget string
	action       string // "freeze" or "unfreeze"

	feeAuto         bool
	validRoundsAuto bool
	note            string

	status  string
	runFunc goal.RunFunc

	// Field editing state
	fieldsState *FieldsState
}

// NewAssetFreezeBuilder creates a new AssetFreezeBuilder
func NewAssetFreezeBuilder() *AssetFreezeBuilder {
	b := &AssetFreezeBuilder{
		feeAuto:         true,
		validRoundsAuto: true,
		action:          "freeze",
	}
	b.initFields()
	return b
}

// initFields initializes the field definitions
func (b *AssetFreezeBuilder) initFields() {
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
			Label:       "Sender (Freezer)",
			Placeholder: "[freeze authority address - 58 chars]",
			Value:       &b.sender,
			Required:    true,
			Type:        FieldTypeAddress,
		},
		{
			Key:         "freezeTarget",
			Label:       "Target Account",
			Placeholder: "[account to freeze/unfreeze - 58 chars]",
			Value:       &b.freezeTarget,
			Required:    true,
			Type:        FieldTypeAddress,
		},
		{
			Key:         "action",
			Label:       "Action",
			Placeholder: "",
			Value:       &b.action,
			Required:    true,
			Type:        FieldTypeSelect,
			Options:     []string{"freeze", "unfreeze"},
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

func (b *AssetFreezeBuilder) SetRunFunc(fn goal.RunFunc) { b.runFunc = fn }
func (b *AssetFreezeBuilder) Title() string              { return "Asset Freeze (afrz)" }
func (b *AssetFreezeBuilder) TxnType() string            { return "afrz" }
func (b *AssetFreezeBuilder) Init() tea.Cmd              { return nil }

// IsEditing returns true if currently editing a field
func (b *AssetFreezeBuilder) IsEditing() bool {
	return b.fieldsState != nil && b.fieldsState.Editing
}

// SetAvailableHeight configures the viewport based on available terminal lines
func (b *AssetFreezeBuilder) SetAvailableHeight(lines int) {
	if b.fieldsState != nil {
		viewportHeight := CalculateViewportHeight(lines)
		b.fieldsState.SetViewportHeight(viewportHeight)
	}
}

func (b *AssetFreezeBuilder) Validate() error {
	if strings.TrimSpace(b.assetID) == "" {
		return errors.New("asset ID is required")
	}
	if strings.TrimSpace(b.sender) == "" {
		return errors.New("sender is required")
	}
	if strings.TrimSpace(b.freezeTarget) == "" {
		return errors.New("freeze target is required")
	}
	if b.action != "freeze" && b.action != "unfreeze" {
		return errors.New("action must be freeze or unfreeze")
	}

	if result := ValidateAssetID(b.assetID); !result.Valid {
		return fmt.Errorf("asset ID: %s", result.Message)
	}
	if result := ValidateAlgorandAddress(b.sender); !result.Valid {
		return fmt.Errorf("sender: %s", result.Message)
	}
	if result := ValidateAlgorandAddress(b.freezeTarget); !result.Valid {
		return fmt.Errorf("freeze target: %s", result.Message)
	}

	if b.fieldsState != nil && b.fieldsState.HasValidationErrors() {
		return errors.New("please fix validation errors before executing")
	}

	return nil
}

func (b *AssetFreezeBuilder) Args() []string {
	args := []string{"asset", "freeze"}
	args = append(args, "--assetid", strings.TrimSpace(b.assetID))
	args = append(args, "--freezer", strings.TrimSpace(b.sender))
	args = append(args, "--account", strings.TrimSpace(b.freezeTarget))

	if b.action == "unfreeze" {
		args = append(args, "--freeze=false")
	} else {
		args = append(args, "--freeze=true")
	}

	return args
}

func (b *AssetFreezeBuilder) AfterRun(stdout, stderr string, runErr error) {
	if runErr != nil {
		b.status = fmt.Sprintf("Error: %v\n%s", runErr, stderr)
		return
	}
	b.status = stdout
}

func (b *AssetFreezeBuilder) Update(msg tea.Msg) (goal.Builder, tea.Cmd) {
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

func (b *AssetFreezeBuilder) RenderFields() string {
	return RenderFieldsPanel("Asset Freeze", b.fieldsState, 4)
}

func (b *AssetFreezeBuilder) RenderOutput() string {
	return RenderOutput(b.status)
}

func (b *AssetFreezeBuilder) View() string {
	return b.RenderFields()
}
