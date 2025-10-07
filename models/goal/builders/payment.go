package builders

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"lazychain/models/goal/components"
	goal "lazychain/models/goal/iface"
)

// PaymentBuilder maps to `goal clerk send`.
// Flags reference: https://developer.algorand.org/docs/clis/goal/clerk/send/
type PaymentBuilder struct {
	fields []*components.Field
	idx    int

	status string

	// plumbed in by host model:
	RunWith func(argv []string)
}

func NewPaymentBuilder() *PaymentBuilder {
	return &PaymentBuilder{
		fields: []*components.Field{
			{Label: "From (-f)", Hint: "address (or default)", Active: true},
			{Label: "To (-t)", Hint: "recipient address"},
			{Label: "Amount μAlgos (-a)", Hint: "e.g. 1000000 = 1 Algo"},
			{Label: "Fee μAlgos (--fee)", Hint: "optional; empty for suggested"},
			{Label: "FirstValid (--firstvalid)", Hint: "optional"},
			{Label: "LastValid (--lastvalid)", Hint: "optional"},
			{Label: "Note (-n)", Hint: "plain text note (optional)"},
			{Label: "Out file (-o)", Hint: "write txn to file (optional)"},
			{Label: "Sign (-s)", Hint: "true/false, with -o", },
			{Label: "No Wait (-N)", Hint: "true/false"},
			{Label: "Rekey (--rekey-to)", Hint: "optional"},
		},
	}
}

func (p *PaymentBuilder) Title() string { return "Payment (goal clerk send)" }
func (p *PaymentBuilder) Init() tea.Cmd { return nil }

func (p *PaymentBuilder) Validate() error {
	to := p.fields[1].Value
	amt := p.fields[2].Value
	if strings.TrimSpace(to) == "" {
		return errors.New("recipient (-t) is required")
	}
	if _, err := strconv.ParseUint(strings.TrimSpace(amt), 10, 64); err != nil {
		return errors.New("amount (-a) must be a positive integer μAlgos")
	}
	return nil
}

func (p *PaymentBuilder) Args() []string {
	f := func(i int) string { return strings.TrimSpace(p.fields[i].Value) }
	var argv []string
	argv = append(argv, "clerk", "send")
	if v := f(0); v != "" { argv = append(argv, "-f", v) }
	argv = append(argv, "-t", f(1))
	argv = append(argv, "-a", f(2))
	if v := f(3); v != "" { argv = append(argv, "--fee", v) }
	if v := f(4); v != "" { argv = append(argv, "--firstvalid", v) }
	if v := f(5); v != "" { argv = append(argv, "--lastvalid", v) }
	if v := f(6); v != "" { argv = append(argv, "-n", v) }
	if v := f(7); v != "" { argv = append(argv, "-o", v) }
	if strings.EqualFold(f(8), "true") { argv = append(argv, "-s") }
	if strings.EqualFold(f(9), "true") { argv = append(argv, "-N") }
	if v := f(10); v != "" { argv = append(argv, "--rekey-to", v) }
	return argv
}

func (p *PaymentBuilder) AfterRun(stdout, stderr string, runErr error) {
	if runErr != nil {
		p.status = fmt.Sprintf("Error: %v\n%s", runErr, strings.TrimSpace(stderr))
		return
	}
	p.status = strings.TrimSpace(stdout)
}

func (p *PaymentBuilder) Update(msg tea.Msg) (goal.Builder, tea.Cmd) {
	switch m := msg.(type) {
	case tea.KeyMsg:
		switch m.String() {
		case "tab", "down":
			p.fields[p.idx].Active = false
			p.idx = (p.idx + 1) % len(p.fields)
			p.fields[p.idx].Active = true
		case "shift+tab", "up":
			p.fields[p.idx].Active = false
			p.idx = (p.idx - 1 + len(p.fields)) % len(p.fields)
			p.fields[p.idx].Active = true
		case "left":
			p.fields[p.idx].MoveLeft()
		case "right":
			p.fields[p.idx].MoveRight()
		case "backspace":
			p.fields[p.idx].Backspace()
		case "enter":
			// run command
			if err := p.Validate(); err != nil {
				p.status = "Validation: " + err.Error()
				return p, nil
			}
			if p.RunWith != nil {
				p.RunWith(p.Args())
			}
		}
		// typing
		if m.Type == tea.KeyRunes {
			for _, r := range m.Runes {
				p.fields[p.idx].InsertRune(r)
			}
		}
	}
	return p, nil
}

func (p *PaymentBuilder) View() string {
	left := []string{
		lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#cba6f7")).Render(p.Title()),
		"",
	}
	for _, f := range p.fields {
		left = append(left, f.Render(36))
		left = append(left, "")
	}
	leftPanel := lipgloss.NewStyle().
		Width(44).
		Padding(1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#89b4fa")).
		Render(stringsJoin(left))

	right := []string{
		lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#a6e3a1")).Render("Output"),
		"",
		strings.TrimSpace(p.status),
	}
	rightPanel := lipgloss.NewStyle().
		Width(44).
		Padding(1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#a6e3a1")).
		Render(stringsJoin(right))

	return lipgloss.JoinHorizontal(lipgloss.Top, leftPanel, "  ", rightPanel)
}

func stringsJoin(ss []string) string {
	var b strings.Builder
	for i, s := range ss {
		if i > 0 { b.WriteString("\n") }
		b.WriteString(s)
	}
	return b.String()
}
