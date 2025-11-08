package builders

import (
	"errors"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	goal "lazychain/models/goal/iface"
)

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
}

func NewApplicationCallBuilder() *ApplicationCallBuilder {
	return &ApplicationCallBuilder{
		feeAuto:         true,
		validRoundsAuto: true,
	}
}

func (b *ApplicationCallBuilder) SetRunFunc(fn goal.RunFunc) { b.runFunc = fn }
func (b *ApplicationCallBuilder) Title() string              { return "Application Call (appl)" }
func (b *ApplicationCallBuilder) TxnType() string            { return "appl" }
func (b *ApplicationCallBuilder) Init() tea.Cmd              { return nil }

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
	return nil
}

func (b *ApplicationCallBuilder) Args() []string {
	var args []string
	args = append(args, "app", "call")
	args = append(args, "--app-id", strings.TrimSpace(b.appID))
	args = append(args, "--from", strings.TrimSpace(b.sender))
	// TODO: Add on-complete, arguments, etc.
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
		switch msg.String() {
		case "enter":
			if err := b.Validate(); err != nil {
				b.status = "Validation Error: " + err.Error()
				return b, nil
			}
			if b.runFunc != nil {
				b.runFunc(b.Args())
			}
		}
	}
	return b, nil
}

func (b *ApplicationCallBuilder) RenderFields() string {
	var lines []string

	title := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#cba6f7")).Render("Application Call")
	lines = append(lines, title, "")

	// Required
	lines = append(lines, lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#89b4fa")).Render("Required Fields:"), "")
	lines = append(lines, "  Application ID:", "    "+getPlaceholder(b.appID, "[app id]"), "")
	lines = append(lines, "  On Complete:", "    "+getPlaceholder(b.onComplete, "[NoOp/OptIn/CloseOut/etc]"), "")
	lines = append(lines, "  Sender:", "    "+getPlaceholder(b.sender, "[address]"), "")

	// Optional
	lines = append(lines, lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#f9e2af")).Render("Optional Fields:"), "")
	lines = append(lines, "  Arguments:", "    "+getPlaceholder(b.arguments, "[app arguments]"), "")
	lines = append(lines, "  Fee auto:", "    "+getBoolStr(b.feeAuto), "")
	lines = append(lines, "  Valid rounds auto:", "    "+getBoolStr(b.validRoundsAuto), "")
	lines = append(lines, "  Note:", "    "+getPlaceholder(b.note, "[optional]"), "")

	lines = append(lines, lipgloss.NewStyle().Italic(true).Foreground(lipgloss.Color("#6c7086")).Render("Enter: Execute | (Fields editing coming soon)"))

	return strings.Join(lines, "\n")
}

func (b *ApplicationCallBuilder) RenderOutput() string {
	return renderOutput(b.status)
}

// View returns the complete view (delegates to RenderFields for compatibility)
func (b *ApplicationCallBuilder) View() string {
	return b.RenderFields()
}

// Helper functions (shared)
func getPlaceholder(value, placeholder string) string {
	if strings.TrimSpace(value) == "" {
		return lipgloss.NewStyle().Faint(true).Render(placeholder)
	}
	return value
}

func getBoolStr(value bool) string {
	if value {
		return "✓ Yes"
	}
	return "✗ No"
}

func renderOutput(status string) string {
	if status == "" {
		return lipgloss.NewStyle().Italic(true).Foreground(lipgloss.Color("#6c7086")).Render("No output yet.")
	}
	if strings.HasPrefix(status, "Error") || strings.HasPrefix(status, "Validation Error") {
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#f38ba8")).Render(status)
	}
	return lipgloss.NewStyle().Foreground(lipgloss.Color("#a6e3a1")).Render(status)
}
