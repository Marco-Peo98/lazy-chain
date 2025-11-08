package builders

import (
	"errors"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	goal "lazychain/models/goal/iface"
)

type AssetCreateBuilder struct {
	total    string
	decimals string
	creator  string

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
}

func NewAssetCreateBuilder() *AssetCreateBuilder {
	return &AssetCreateBuilder{
		feeAuto:         true,
		validRoundsAuto: true,
		freezeByDefault: false,
	}
}

func (b *AssetCreateBuilder) SetRunFunc(fn goal.RunFunc) { b.runFunc = fn }
func (b *AssetCreateBuilder) Title() string              { return "Asset Create (acfg)" }
func (b *AssetCreateBuilder) TxnType() string            { return "acfg" }
func (b *AssetCreateBuilder) Init() tea.Cmd              { return nil }

func (b *AssetCreateBuilder) Validate() error {
	if strings.TrimSpace(b.total) == "" {
		return errors.New("total required")
	}
	if strings.TrimSpace(b.decimals) == "" {
		return errors.New("decimals required")
	}
	if strings.TrimSpace(b.creator) == "" {
		return errors.New("creator required")
	}
	return nil
}

func (b *AssetCreateBuilder) Args() []string {
	return []string{"asset", "create", "--creator", b.creator, "--total", b.total, "--decimals", b.decimals}
}

func (b *AssetCreateBuilder) AfterRun(stdout, stderr string, runErr error) {
	if runErr != nil {
		b.status = fmt.Sprintf("Error: %v\n%s", runErr, stderr)
		return
	}
	b.status = stdout
}

func (b *AssetCreateBuilder) Update(msg tea.Msg) (goal.Builder, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok && msg.String() == "enter" {
		if err := b.Validate(); err != nil {
			b.status = "Validation Error: " + err.Error()
		} else if b.runFunc != nil {
			b.runFunc(b.Args())
		}
	}
	return b, nil
}

func (b *AssetCreateBuilder) RenderFields() string {
	var l []string
	l = append(l, lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#cba6f7")).Render("Asset Create"), "")
	l = append(l, lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#89b4fa")).Render("Required:"), "")
	l = append(l, "  Total:", "    "+getPlaceholder(b.total, "[total supply]"), "")
	l = append(l, "  Decimals:", "    "+getPlaceholder(b.decimals, "[0-19]"), "")
	l = append(l, "  Creator:", "    "+getPlaceholder(b.creator, "[address]"), "")
	l = append(l, lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#f9e2af")).Render("Optional:"), "")
	l = append(l, "  Asset Name:", "    "+getPlaceholder(b.assetName, "[name]"), "")
	l = append(l, "  Unit Name:", "    "+getPlaceholder(b.unitName, "[symbol]"), "")
	l = append(l, "  (... and 10 more fields - coming soon)")
	return strings.Join(l, "\n")
}

func (b *AssetCreateBuilder) RenderOutput() string { return renderOutput(b.status) }

// View returns the complete view (delegates to RenderFields for compatibility)
func (b *AssetCreateBuilder) View() string {
	return b.RenderFields()
}
