package builders

import (
	"errors"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	goal "lazychain/models/goal/iface"
)

// KeyRegistrationBuilder implements key registration transaction (keyreg)
type KeyRegistrationBuilder struct {
	sender           string
	registration     string // "online" or "offline"
	votingKey        string
	selectionKey     string
	stateProofKey    string
	firstVotingRound string
	lastVotingRound  string
	voteKeyDilution  string

	feeAuto         bool
	validRoundsAuto bool
	note            string

	status  string
	runFunc goal.RunFunc

	// Field editing state
	fieldsState *FieldsState
}

// NewKeyRegistrationBuilder creates a new KeyRegistrationBuilder
func NewKeyRegistrationBuilder() *KeyRegistrationBuilder {
	b := &KeyRegistrationBuilder{
		feeAuto:         true,
		validRoundsAuto: true,
		registration:    "online",
	}
	b.initFields()
	return b
}

// initFields initializes the field definitions
func (b *KeyRegistrationBuilder) initFields() {
	fields := []FieldDef{
		// Required fields (first 2 always required)
		{
			Key:         "sender",
			Label:       "Sender",
			Placeholder: "[account address - 58 chars]",
			Value:       &b.sender,
			Required:    true,
			Type:        FieldTypeAddress,
		},
		{
			Key:         "registration",
			Label:       "Registration Type",
			Placeholder: "",
			Value:       &b.registration,
			Required:    true,
			Type:        FieldTypeSelect,
			Options:     []string{"online", "offline"},
		},
		// Participation keys (required for online)
		{
			Key:         "votingKey",
			Label:       "Voting Key (votepk)",
			Placeholder: "[base64 public key ~44 chars]",
			Value:       &b.votingKey,
			Required:    false, // Conditional on registration type
			Type:        FieldTypeText,
		},
		{
			Key:         "selectionKey",
			Label:       "Selection Key (selkey)",
			Placeholder: "[base64 public key ~44 chars]",
			Value:       &b.selectionKey,
			Required:    false,
			Type:        FieldTypeText,
		},
		{
			Key:         "firstVotingRound",
			Label:       "First Voting Round",
			Placeholder: "[round number]",
			Value:       &b.firstVotingRound,
			Required:    false,
			Type:        FieldTypeInteger,
		},
		{
			Key:         "lastVotingRound",
			Label:       "Last Voting Round",
			Placeholder: "[round number]",
			Value:       &b.lastVotingRound,
			Required:    false,
			Type:        FieldTypeInteger,
		},
		{
			Key:         "voteKeyDilution",
			Label:       "Vote Key Dilution",
			Placeholder: "[dilution value]",
			Value:       &b.voteKeyDilution,
			Required:    false,
			Type:        FieldTypeInteger,
		},
		{
			Key:         "stateProofKey",
			Label:       "State Proof Key",
			Placeholder: "[base64 key, optional]",
			Value:       &b.stateProofKey,
			Required:    false,
			Type:        FieldTypeText,
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

func (b *KeyRegistrationBuilder) SetRunFunc(fn goal.RunFunc) { b.runFunc = fn }
func (b *KeyRegistrationBuilder) Title() string              { return "Key Registration (keyreg)" }
func (b *KeyRegistrationBuilder) TxnType() string            { return "keyreg" }
func (b *KeyRegistrationBuilder) Init() tea.Cmd              { return nil }

// IsEditing returns true if currently editing a field
func (b *KeyRegistrationBuilder) IsEditing() bool {
	return b.fieldsState != nil && b.fieldsState.Editing
}

// SetAvailableHeight configures the viewport based on available terminal lines
func (b *KeyRegistrationBuilder) SetAvailableHeight(lines int) {
	if b.fieldsState != nil {
		viewportHeight := CalculateViewportHeight(lines)
		b.fieldsState.SetViewportHeight(viewportHeight)
	}
}

func (b *KeyRegistrationBuilder) Validate() error {
	if strings.TrimSpace(b.sender) == "" {
		return errors.New("sender is required")
	}
	if b.registration != "online" && b.registration != "offline" {
		return errors.New("registration must be online or offline")
	}

	if result := ValidateAlgorandAddress(b.sender); !result.Valid {
		return fmt.Errorf("sender: %s", result.Message)
	}

	if b.registration == "online" {
		if strings.TrimSpace(b.votingKey) == "" {
			return errors.New("voting key is required for online registration")
		}
		if strings.TrimSpace(b.selectionKey) == "" {
			return errors.New("selection key is required for online registration")
		}
		if strings.TrimSpace(b.firstVotingRound) == "" {
			return errors.New("first voting round is required for online registration")
		}
		if strings.TrimSpace(b.lastVotingRound) == "" {
			return errors.New("last voting round is required for online registration")
		}
		if strings.TrimSpace(b.voteKeyDilution) == "" {
			return errors.New("vote key dilution is required for online registration")
		}

		if result := ValidateVotingKey(b.votingKey); !result.Valid {
			return fmt.Errorf("voting key: %s", result.Message)
		}
		if result := ValidateVotingKey(b.selectionKey); !result.Valid {
			return fmt.Errorf("selection key: %s", result.Message)
		}

		if result := ValidateRoundNumber(b.firstVotingRound); !result.Valid {
			return fmt.Errorf("first voting round: %s", result.Message)
		}
		if result := ValidateRoundNumber(b.lastVotingRound); !result.Valid {
			return fmt.Errorf("last voting round: %s", result.Message)
		}
		if result := ValidateInteger(b.voteKeyDilution); !result.Valid {
			return fmt.Errorf("vote key dilution: %s", result.Message)
		}

		if strings.TrimSpace(b.stateProofKey) != "" {
			if result := ValidateBase64(b.stateProofKey); !result.Valid {
				return fmt.Errorf("state proof key: %s", result.Message)
			}
		}
	}

	if b.fieldsState != nil && b.fieldsState.HasValidationErrors() {
		return errors.New("please fix validation errors before executing")
	}

	return nil
}

func (b *KeyRegistrationBuilder) Args() []string {
	args := []string{"account", "changeonlinestatus"}
	args = append(args, "--address", strings.TrimSpace(b.sender))

	if b.registration == "online" {
		args = append(args, "--online")
		args = append(args, "--votepk", strings.TrimSpace(b.votingKey))
		args = append(args, "--selkey", strings.TrimSpace(b.selectionKey))
		args = append(args, "--votefst", strings.TrimSpace(b.firstVotingRound))
		args = append(args, "--votelst", strings.TrimSpace(b.lastVotingRound))
		args = append(args, "--votekd", strings.TrimSpace(b.voteKeyDilution))

		if strings.TrimSpace(b.stateProofKey) != "" {
			args = append(args, "--stateproofkey", strings.TrimSpace(b.stateProofKey))
		}
	} else {
		args = append(args, "--offline")
	}

	return args
}

func (b *KeyRegistrationBuilder) AfterRun(stdout, stderr string, runErr error) {
	if runErr != nil {
		b.status = fmt.Sprintf("Error: %v\n%s", runErr, stderr)
		return
	}
	b.status = stdout
}

func (b *KeyRegistrationBuilder) Update(msg tea.Msg) (goal.Builder, tea.Cmd) {
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

func (b *KeyRegistrationBuilder) RenderFields() string {
	return RenderFieldsPanel("Key Registration", b.fieldsState, 2)
}

func (b *KeyRegistrationBuilder) RenderOutput() string {
	return RenderOutput(b.status)
}

func (b *KeyRegistrationBuilder) View() string {
	return b.RenderFields()
}
