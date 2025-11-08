package builders

import (
	"errors"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	goal "lazychain/models/goal/iface"
)

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
}

func NewAssetFreezeBuilder() *AssetFreezeBuilder {
	return &AssetFreezeBuilder{
		feeAuto:         true,
		validRoundsAuto: true,
		action:          "freeze",
	}
}

func (b *AssetFreezeBuilder) SetRunFunc(fn goal.RunFunc) { b.runFunc = fn }
func (b *AssetFreezeBuilder) Title() string              { return "Asset Freeze (afrz)" }
func (b *AssetFreezeBuilder) TxnType() string            { return "afrz" }
func (b *AssetFreezeBuilder) Init() tea.Cmd              { return nil }

func (b *AssetFreezeBuilder) Validate() error {
	if strings.TrimSpace(b.assetID) == "" {
		return errors.New("asset ID required")
	}
	if strings.TrimSpace(b.sender) == "" {
		return errors.New("sender required")
	}
	if strings.TrimSpace(b.freezeTarget) == "" {
		return errors.New("freeze target required")
	}
	if b.action != "freeze" && b.action != "unfreeze" {
		return errors.New("action must be freeze or unfreeze")
	}
	return nil
}

func (b *AssetFreezeBuilder) Args() []string {
	return []string{"asset", "freeze", "--assetid", b.assetID, "--freezer", b.sender, "--target", b.freezeTarget}
}

func (b *AssetFreezeBuilder) AfterRun(stdout, stderr string, runErr error) {
	if runErr != nil {
		b.status = fmt.Sprintf("Error: %v\n%s", runErr, stderr)
		return
	}
	b.status = stdout
}

func (b *AssetFreezeBuilder) Update(msg tea.Msg) (goal.Builder, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok && msg.String() == "enter" {
		if err := b.Validate(); err != nil {
			b.status = "Validation Error: " + err.Error()
		} else if b.runFunc != nil {
			b.runFunc(b.Args())
		}
	}
	return b, nil
}

func (b *AssetFreezeBuilder) RenderFields() string {
	var l []string
	l = append(l, lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#cba6f7")).Render("Asset Freeze"), "")
	l = append(l, lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#89b4fa")).Render("Required:"), "")
	l = append(l, "  Asset ID:", "    "+getPlaceholder(b.assetID, "[asset id]"), "")
	l = append(l, "  Sender:", "    "+getPlaceholder(b.sender, "[address]"), "")
	l = append(l, "  Freeze Target:", "    "+getPlaceholder(b.freezeTarget, "[address]"), "")
	l = append(l, "  Action:", "    "+b.action, "")
	l = append(l, lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#f9e2af")).Render("Optional:"), "")
	l = append(l, "  Fee auto: "+getBoolStr(b.feeAuto))
	return strings.Join(l, "\n")
}

func (b *AssetFreezeBuilder) RenderOutput() string { return renderOutput(b.status) }

// View returns the complete view (delegates to RenderFields for compatibility)
func (b *AssetFreezeBuilder) View() string {
	return b.RenderFields()
}
