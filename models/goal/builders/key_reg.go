package builders

import (
	"errors"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	goal "lazychain/models/goal/iface"
)

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
}

func NewKeyRegistrationBuilder() *KeyRegistrationBuilder {
	return &KeyRegistrationBuilder{
		feeAuto:         true,
		validRoundsAuto: true,
		registration:    "online",
	}
}

func (b *KeyRegistrationBuilder) SetRunFunc(fn goal.RunFunc) { b.runFunc = fn }
func (b *KeyRegistrationBuilder) Title() string              { return "Key Registration (keyreg)" }
func (b *KeyRegistrationBuilder) TxnType() string            { return "keyreg" }
func (b *KeyRegistrationBuilder) Init() tea.Cmd              { return nil }

func (b *KeyRegistrationBuilder) Validate() error {
	if strings.TrimSpace(b.sender) == "" {
		return errors.New("sender required")
	}
	if b.registration == "online" {
		if strings.TrimSpace(b.votingKey) == "" {
			return errors.New("voting key required for online")
		}
		if strings.TrimSpace(b.selectionKey) == "" {
			return errors.New("selection key required for online")
		}
	}
	return nil
}

func (b *KeyRegistrationBuilder) Args() []string {
	return []string{"account", "changeonlinestatus", "--address", b.sender, "--online=" + b.registration}
}

func (b *KeyRegistrationBuilder) AfterRun(stdout, stderr string, runErr error) {
	if runErr != nil {
		b.status = fmt.Sprintf("Error: %v\n%s", runErr, stderr)
		return
	}
	b.status = stdout
}

func (b *KeyRegistrationBuilder) Update(msg tea.Msg) (goal.Builder, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok && msg.String() == "enter" {
		if err := b.Validate(); err != nil {
			b.status = "Validation Error: " + err.Error()
		} else if b.runFunc != nil {
			b.runFunc(b.Args())
		}
	}
	return b, nil
}

func (b *KeyRegistrationBuilder) RenderFields() string {
	var l []string
	l = append(l, lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#cba6f7")).Render("Key Registration"), "")
	l = append(l, lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#89b4fa")).Render("Required:"), "")
	l = append(l, "  Sender:", "    "+getPlaceholder(b.sender, "[address]"), "")
	l = append(l, "  Registration:", "    "+b.registration, "")
	l = append(l, "  Voting Key:", "    "+getPlaceholder(b.votingKey, "[key]"), "")
	l = append(l, "  Selection Key:", "    "+getPlaceholder(b.selectionKey, "[key]"), "")
	l = append(l, "  (... 4 more required fields - coming soon)")
	return strings.Join(l, "\n")
}

func (b *KeyRegistrationBuilder) RenderOutput() string { return renderOutput(b.status) }

// View returns the complete view (delegates to RenderFields for compatibility)
func (b *KeyRegistrationBuilder) View() string {
	return b.RenderFields()
}
