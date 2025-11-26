package builders

import (
	"errors"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	goal "lazychain/models/goal/iface"
)

// AssetCreateBuilder implements asset create transaction (acfg)
type AssetCreateBuilder struct {
	// Required
	total    string
	decimals string
	creator  string

	// Optional
	assetName       string
	unitName        string
	manager         string
	reserve         string
	freeze          string
	clawback        string
	freezeByDefault bool
	url             string
	metadataHash    string
	feeAuto         bool
	validRoundsAuto bool
	note            string

	status  string
	runFunc goal.RunFunc

	// Field editing state
	fieldsState *FieldsState
}

// NewAssetCreateBuilder creates a new AssetCreateBuilder
func NewAssetCreateBuilder() *AssetCreateBuilder {
	b := &AssetCreateBuilder{
		feeAuto:         true,
		validRoundsAuto: true,
		freezeByDefault: false,
	}
	b.initFields()
	return b
}

// initFields initializes the field definitions
func (b *AssetCreateBuilder) initFields() {
	fields := []FieldDef{
		// Required fields (first 3)
		{
			Key:         "total",
			Label:       "Total Supply",
			Placeholder: "[total units - positive number]",
			Value:       &b.total,
			Required:    true,
			Type:        FieldTypeAmount,
		},
		{
			Key:         "decimals",
			Label:       "Decimals",
			Placeholder: "[0-19]",
			Value:       &b.decimals,
			Required:    true,
			Type:        FieldTypeInteger,
		},
		{
			Key:         "creator",
			Label:       "Creator",
			Placeholder: "[creator address - 58 chars]",
			Value:       &b.creator,
			Required:    true,
			Type:        FieldTypeAddress,
		},
		// Optional fields
		{
			Key:         "assetName",
			Label:       "Asset Name",
			Placeholder: "[name, max 32 chars]",
			Value:       &b.assetName,
			Required:    false,
			Type:        FieldTypeText,
		},
		{
			Key:         "unitName",
			Label:       "Unit Name",
			Placeholder: "[symbol, e.g. USDC, max 8 chars]",
			Value:       &b.unitName,
			Required:    false,
			Type:        FieldTypeText,
		},
		{
			Key:         "manager",
			Label:       "Manager",
			Placeholder: "[manager address - 58 chars]",
			Value:       &b.manager,
			Required:    false,
			Type:        FieldTypeAddress,
		},
		{
			Key:         "reserve",
			Label:       "Reserve",
			Placeholder: "[reserve address - 58 chars]",
			Value:       &b.reserve,
			Required:    false,
			Type:        FieldTypeAddress,
		},
		{
			Key:         "freeze",
			Label:       "Freeze Address",
			Placeholder: "[freeze address - 58 chars]",
			Value:       &b.freeze,
			Required:    false,
			Type:        FieldTypeAddress,
		},
		{
			Key:         "clawback",
			Label:       "Clawback",
			Placeholder: "[clawback address - 58 chars]",
			Value:       &b.clawback,
			Required:    false,
			Type:        FieldTypeAddress,
		},
		{
			Key:         "freezeByDefault",
			Label:       "Freeze by Default",
			Placeholder: "",
			BoolValue:   &b.freezeByDefault,
			Required:    false,
			Type:        FieldTypeBool,
		},
		{
			Key:         "url",
			Label:       "URL",
			Placeholder: "[asset url, max 96 chars]",
			Value:       &b.url,
			Required:    false,
			Type:        FieldTypeText,
		},
		{
			Key:         "metadataHash",
			Label:       "Metadata Hash",
			Placeholder: "[base64 hash, 32 bytes]",
			Value:       &b.metadataHash,
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

func (b *AssetCreateBuilder) SetRunFunc(fn goal.RunFunc) { b.runFunc = fn }
func (b *AssetCreateBuilder) Title() string              { return "Asset Create (acfg)" }
func (b *AssetCreateBuilder) TxnType() string            { return "acfg" }
func (b *AssetCreateBuilder) Init() tea.Cmd              { return nil }

// IsEditing returns true if currently editing a field
func (b *AssetCreateBuilder) IsEditing() bool {
	return b.fieldsState != nil && b.fieldsState.Editing
}

// SetAvailableHeight configures the viewport based on available terminal lines
func (b *AssetCreateBuilder) SetAvailableHeight(lines int) {
	if b.fieldsState != nil {
		viewportHeight := CalculateViewportHeight(lines)
		b.fieldsState.SetViewportHeight(viewportHeight)
	}
}

func (b *AssetCreateBuilder) Validate() error {
	if strings.TrimSpace(b.total) == "" {
		return errors.New("total supply is required")
	}
	if strings.TrimSpace(b.decimals) == "" {
		return errors.New("decimals is required")
	}
	if strings.TrimSpace(b.creator) == "" {
		return errors.New("creator is required")
	}

	if result := ValidateAmountAllowZero(b.total); !result.Valid {
		return fmt.Errorf("total supply: %s", result.Message)
	}
	if result := ValidateDecimals(b.decimals); !result.Valid {
		return fmt.Errorf("decimals: %s", result.Message)
	}
	if result := ValidateAlgorandAddress(b.creator); !result.Valid {
		return fmt.Errorf("creator: %s", result.Message)
	}

	// Validate optional address fields if provided
	if strings.TrimSpace(b.manager) != "" {
		if result := ValidateAlgorandAddress(b.manager); !result.Valid {
			return fmt.Errorf("manager: %s", result.Message)
		}
	}
	if strings.TrimSpace(b.reserve) != "" {
		if result := ValidateAlgorandAddress(b.reserve); !result.Valid {
			return fmt.Errorf("reserve: %s", result.Message)
		}
	}
	if strings.TrimSpace(b.freeze) != "" {
		if result := ValidateAlgorandAddress(b.freeze); !result.Valid {
			return fmt.Errorf("freeze: %s", result.Message)
		}
	}
	if strings.TrimSpace(b.clawback) != "" {
		if result := ValidateAlgorandAddress(b.clawback); !result.Valid {
			return fmt.Errorf("clawback: %s", result.Message)
		}
	}

	if strings.TrimSpace(b.metadataHash) != "" {
		if result := ValidateBase64(b.metadataHash); !result.Valid {
			return fmt.Errorf("metadata hash: %s", result.Message)
		}
	}

	if b.fieldsState != nil && b.fieldsState.HasValidationErrors() {
		return errors.New("please fix validation errors before executing")
	}

	return nil
}

func (b *AssetCreateBuilder) Args() []string {
	args := []string{
		"asset", "create",
		"--creator", strings.TrimSpace(b.creator),
		"--total", strings.TrimSpace(b.total),
		"--decimals", strings.TrimSpace(b.decimals),
	}

	if strings.TrimSpace(b.assetName) != "" {
		args = append(args, "--name", b.assetName)
	}
	if strings.TrimSpace(b.unitName) != "" {
		args = append(args, "--unitname", b.unitName)
	}
	if strings.TrimSpace(b.manager) != "" {
		args = append(args, "--manager", b.manager)
	}
	if strings.TrimSpace(b.reserve) != "" {
		args = append(args, "--reserve", b.reserve)
	}
	if strings.TrimSpace(b.freeze) != "" {
		args = append(args, "--freezer", b.freeze)
	}
	if strings.TrimSpace(b.clawback) != "" {
		args = append(args, "--clawback", b.clawback)
	}
	if strings.TrimSpace(b.url) != "" {
		args = append(args, "--asseturl", b.url)
	}
	if strings.TrimSpace(b.metadataHash) != "" {
		args = append(args, "--assetmetadatab64", b.metadataHash)
	}
	if b.freezeByDefault {
		args = append(args, "--defaultfrozen")
	}

	return args
}

func (b *AssetCreateBuilder) AfterRun(stdout, stderr string, runErr error) {
	if runErr != nil {
		b.status = fmt.Sprintf("Error: %v\n%s", runErr, stderr)
		return
	}
	b.status = stdout
}

func (b *AssetCreateBuilder) Update(msg tea.Msg) (goal.Builder, tea.Cmd) {
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

func (b *AssetCreateBuilder) RenderFields() string {
	return RenderFieldsPanel("Asset Create", b.fieldsState, 3)
}

func (b *AssetCreateBuilder) RenderOutput() string {
	return RenderOutput(b.status)
}

func (b *AssetCreateBuilder) View() string {
	return b.RenderFields()
}
