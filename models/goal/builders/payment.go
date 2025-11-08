package builders

import (
	"errors"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	goal "lazychain/models/goal/iface"
)

// PaymentBuilder implements payment transaction (pay)
type PaymentBuilder struct {
	// Required fields
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
}

func NewPaymentBuilder() *PaymentBuilder {
	return &PaymentBuilder{
		feeAuto:         true, // Default: auto
		validRoundsAuto: true, // Default: auto
	}
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

// Validate checks if required fields are filled
func (b *PaymentBuilder) Validate() error {
	if strings.TrimSpace(b.sender) == "" {
		return errors.New("sender is required")
	}
	if strings.TrimSpace(b.receiver) == "" {
		return errors.New("receiver is required")
	}
	if strings.TrimSpace(b.amount) == "" {
		return errors.New("amount is required")
	}
	// TODO: Add more validation (address format, amount numeric, etc.)
	return nil
}

// Args builds the command arguments
func (b *PaymentBuilder) Args() []string {
	// This is a placeholder - will be implemented properly with actual goal command
	var args []string
	args = append(args, "clerk", "send")
	args = append(args, "-f", strings.TrimSpace(b.sender))
	args = append(args, "-t", strings.TrimSpace(b.receiver))
	args = append(args, "-a", strings.TrimSpace(b.amount))

	if b.note != "" {
		args = append(args, "-n", b.note)
	}

	// TODO: Add fee and valid rounds handling

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
		switch msg.String() {
		case "enter":
			// Execute command
			if err := b.Validate(); err != nil {
				b.status = "Validation Error: " + err.Error()
				return b, nil
			}
			if b.runFunc != nil {
				b.runFunc(b.Args())
			}
		}
		// TODO: Handle field navigation and input
	}
	return b, nil
}

// RenderFields renders the fields panel (NEW - separated from output)
func (b *PaymentBuilder) RenderFields() string {
	var lines []string

	// Title
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#cba6f7")).
		Render("Payment Transaction")
	lines = append(lines, title)
	lines = append(lines, "")

	// Required fields section
	reqHeader := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#89b4fa")).
		Render("Required Fields:")
	lines = append(lines, reqHeader)
	lines = append(lines, "")

	// Sender (placeholder)
	senderLabel := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#cdd6f4")).
		Render("Sender:")
	senderValue := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#a6e3a1")).
		Render(b.getDisplayValue(b.sender, "[address]"))
	lines = append(lines, "  "+senderLabel)
	lines = append(lines, "    "+senderValue)
	lines = append(lines, "")

	// Receiver (placeholder)
	receiverLabel := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#cdd6f4")).
		Render("Receiver:")
	receiverValue := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#a6e3a1")).
		Render(b.getDisplayValue(b.receiver, "[address]"))
	lines = append(lines, "  "+receiverLabel)
	lines = append(lines, "    "+receiverValue)
	lines = append(lines, "")

	// Amount (placeholder)
	amountLabel := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#cdd6f4")).
		Render("Amount (μAlgos):")
	amountValue := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#a6e3a1")).
		Render(b.getDisplayValue(b.amount, "[amount]"))
	lines = append(lines, "  "+amountLabel)
	lines = append(lines, "    "+amountValue)
	lines = append(lines, "")

	// Optional fields section
	optHeader := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#f9e2af")).
		Render("Optional Fields:")
	lines = append(lines, optHeader)
	lines = append(lines, "")

	// Fee auto
	feeLabel := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#cdd6f4")).
		Render("Set fee automatically:")
	feeValue := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#a6e3a1")).
		Render(b.getBoolDisplay(b.feeAuto))
	lines = append(lines, "  "+feeLabel)
	lines = append(lines, "    "+feeValue)
	lines = append(lines, "")

	// Valid rounds auto
	validLabel := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#cdd6f4")).
		Render("Set valid rounds automatically:")
	validValue := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#a6e3a1")).
		Render(b.getBoolDisplay(b.validRoundsAuto))
	lines = append(lines, "  "+validLabel)
	lines = append(lines, "    "+validValue)
	lines = append(lines, "")

	// Note
	noteLabel := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#cdd6f4")).
		Render("Note:")
	noteValue := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#a6e3a1")).
		Render(b.getDisplayValue(b.note, "[optional note]"))
	lines = append(lines, "  "+noteLabel)
	lines = append(lines, "    "+noteValue)
	lines = append(lines, "")

	// Instructions
	instrStyle := lipgloss.NewStyle().
		Italic(true).
		Foreground(lipgloss.Color("#6c7086"))
	lines = append(lines, instrStyle.Render("Enter: Execute transaction"))
	lines = append(lines, instrStyle.Render("(Field editing coming soon)"))

	return strings.Join(lines, "\n")
}

// RenderOutput renders the output panel (NEW - separated from fields)
func (b *PaymentBuilder) RenderOutput() string {
	if b.status == "" {
		style := lipgloss.NewStyle().
			Italic(true).
			Foreground(lipgloss.Color("#6c7086"))
		return style.Render("No output yet.\nConfigure fields and press Enter to execute.")
	}

	// Check if it's an error
	if strings.HasPrefix(b.status, "Error") || strings.HasPrefix(b.status, "Validation Error") {
		style := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#f38ba8"))
		return style.Render(b.status)
	}

	// Success output
	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#a6e3a1"))
	return style.Render(b.status)
}

// View returns the complete view (delegates to RenderFields for compatibility)
func (b *PaymentBuilder) View() string {
	return b.RenderFields()
}

// Helper functions
func (b *PaymentBuilder) getDisplayValue(value, placeholder string) string {
	if strings.TrimSpace(value) == "" {
		return lipgloss.NewStyle().
			Faint(true).
			Render(placeholder)
	}
	return value
}

func (b *PaymentBuilder) getBoolDisplay(value bool) string {
	if value {
		return "✓ Yes"
	}
	return "✗ No"
}
