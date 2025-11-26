package builders

import (
	"errors"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	goal "lazychain/models/goal/iface"
)

// ApplicationCallBuilder implements application call transaction (appl)
type ApplicationCallBuilder struct {
	// Required
	appID      string
	onComplete string
	sender     string

	// Optional
	arguments       string
	feeAuto         bool
	validRoundsAuto bool
	note            string

	status  string
	runFunc goal.RunFunc

	// Field editing state
	fieldsState *FieldsState
}

// NewApplicationCallBuilder creates a new ApplicationCallBuilder
func NewApplicationCallBuilder() *ApplicationCallBuilder {
	b := &ApplicationCallBuilder{
		feeAuto:         true,
		validRoundsAuto: true,
		onComplete:      "NoOp",
	}
	b.initFields()
	return b
}

// initFields initializes the field definitions
func (b *ApplicationCallBuilder) initFields() {
	fields := []FieldDef{
		// Required fields (first 3)
		{
			Key:         "appID",
			Label:       "Application ID",
			Placeholder: "[app id - positive number]",
			Value:       &b.appID,
			Required:    true,
			Type:        FieldTypeInteger,
		},
		{
			Key:         "onComplete",
			Label:       "On Complete",
			Placeholder: "",
			Value:       &b.onComplete,
			Required:    true,
			Type:        FieldTypeSelect,
			Options:     []string{"NoOp", "OptIn", "CloseOut", "ClearState", "UpdateApplication", "DeleteApplication"},
		},
		{
			Key:         "sender",
			Label:       "Sender",
			Placeholder: "[sender address - 58 chars]",
			Value:       &b.sender,
			Required:    true,
			Type:        FieldTypeAddress,
		},
		// Optional fields
		{
			Key:         "arguments",
			Label:       "Arguments",
			Placeholder: "[comma-separated app arguments]",
			Value:       &b.arguments,
			Required:    false,
			Type:        FieldTypeText,
		},
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

func (b *ApplicationCallBuilder) SetRunFunc(fn goal.RunFunc) { b.runFunc = fn }
func (b *ApplicationCallBuilder) Title() string              { return "Application Call (appl)" }
func (b *ApplicationCallBuilder) TxnType() string            { return "appl" }
func (b *ApplicationCallBuilder) Init() tea.Cmd              { return nil }

// IsEditing returns true if currently editing a field
func (b *ApplicationCallBuilder) IsEditing() bool {
	return b.fieldsState != nil && b.fieldsState.Editing
}

// SetAvailableHeight configures the viewport based on available terminal lines
func (b *ApplicationCallBuilder) SetAvailableHeight(lines int) {
	if b.fieldsState != nil {
		viewportHeight := CalculateViewportHeight(lines)
		b.fieldsState.SetViewportHeight(viewportHeight)
	}
}

func (b *ApplicationCallBuilder) Validate() error {
	if strings.TrimSpace(b.appID) == "" {
		return errors.New("application ID is required")
	}
	if strings.TrimSpace(b.onComplete) == "" {
		return errors.New("on complete is required")
	}
	if strings.TrimSpace(b.sender) == "" {
		return errors.New("sender is required")
	}

	if result := ValidateAppID(b.appID); !result.Valid {
		return fmt.Errorf("application ID: %s", result.Message)
	}
	if result := ValidateAlgorandAddress(b.sender); !result.Valid {
		return fmt.Errorf("sender: %s", result.Message)
	}

	if b.fieldsState != nil && b.fieldsState.HasValidationErrors() {
		return errors.New("please fix validation errors before executing")
	}

	return nil
}

func (b *ApplicationCallBuilder) Args() []string {
	var args []string
	args = append(args, "app", "call")
	args = append(args, "--app-id", strings.TrimSpace(b.appID))
	args = append(args, "--from", strings.TrimSpace(b.sender))

	switch b.onComplete {
	case "OptIn":
		args = append(args, "--on-completion", "optin")
	case "CloseOut":
		args = append(args, "--on-completion", "closeout")
	case "ClearState":
		args = append(args, "--on-completion", "clear")
	case "UpdateApplication":
		args = append(args, "--on-completion", "update")
	case "DeleteApplication":
		args = append(args, "--on-completion", "delete")
	default:
		args = append(args, "--on-completion", "noop")
	}

	if strings.TrimSpace(b.arguments) != "" {
		for _, arg := range strings.Split(b.arguments, ",") {
			args = append(args, "--app-arg", strings.TrimSpace(arg))
		}
	}

	return args
}

func (b *ApplicationCallBuilder) AfterRun(stdout, stderr string, runErr error) {
	if runErr != nil {
		b.status = fmt.Sprintf("Error: %v\n%s", runErr, strings.TrimSpace(stderr))
		return
	}
	b.status = strings.TrimSpace(stdout)
}

func (b *ApplicationCallBuilder) Update(msg tea.Msg) (goal.Builder, tea.Cmd) {
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

func (b *ApplicationCallBuilder) RenderFields() string {
	return RenderFieldsPanel("Application Call", b.fieldsState, 3)
}

func (b *ApplicationCallBuilder) RenderOutput() string {
	return RenderOutput(b.status)
}

func (b *ApplicationCallBuilder) View() string {
	return b.RenderFields()
}
