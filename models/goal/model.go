package goal

import (
	"context"
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"lazychain/models/goal/builders"
	"lazychain/models/goal/components"
)

type GOALModel struct {
	nav     components.ListNav
	builder Builder
	runner  *Runner

	builders []Builder
	output   string
	errLine  string
}

func NewGOALModel() *GOALModel {
	m := &GOALModel{
		runner: NewRunner(),
	}
	// Left menu
	m.nav = components.ListNav{
		Title: "GOAL: Transaction Builder",
		Items: []string{
			"Payment (clerk send)",
			"ASA Transfer (asset send)",
			"App Call (app call/method)",
			"Atomic Group (clerk group)",
			"Sign / Send (clerk sign/rawsend)",
			"Inspect / Simulate",
		},
		Active: true,
		Width:  28,
	}

	// Builders
	pay := builders.NewPaymentBuilder()
	pay.RunWith = m.run

	asa := builders.NewAssetTransferBuilder()
	asa.RunWith = m.run

	app := builders.NewAppCallBuilder()
	app.RunWith = m.run

	group := builders.NewGroupBuilder()
	group.RunWith = m.run

	sign := builders.NewSignSendBuilder()
	sign.RunWith = m.run

	ins := builders.NewInspectSimBuilder()
	ins.RunWith = m.run

	m.builders = []Builder{pay, asa, app, group, sign, ins}
	m.builder = m.builders[0]
	return m
}

func (m *GOALModel) Init() tea.Cmd { return m.builder.Init() }

func (m *GOALModel) run(argv []string) {
	// fire and wait (no background work; we execute synchronously here)
	_ = m.runner.CheckBinary()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	res := m.runner.Run(ctx, argv)
	m.output = strings.TrimSpace(res.Stdout)
	if res.Err != nil {
		if m.output == "" {
			m.errLine = fmt.Sprintf("error: %v\n%s", res.Err, strings.TrimSpace(res.Stderr))
		} else {
			m.errLine = fmt.Sprintf("error: %v", res.Err)
		}
	}
	m.builder.AfterRun(res.Stdout, res.Stderr, res.Err)
}

func (m *GOALModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch t := msg.(type) {
	case tea.KeyMsg:
		switch t.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "up":
			m.nav.Up()
			m.builder = m.builders[m.nav.Cursor]
			return m, nil
		case "down":
			m.nav.Down()
			m.builder = m.builders[m.nav.Cursor]
			return m, nil
		}
	}
	var cmd tea.Cmd
	m.builder, cmd = m.builder.Update(msg)
	return m, cmd
}

func (m *GOALModel) View() string {
	left := m.nav.Render()
	right := m.builder.View()
	footer := m.renderFooter()
	return lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.JoinHorizontal(lipgloss.Top, left, "  ", right),
		"",
		footer,
	)
}

func (m *GOALModel) renderFooter() string {
	info := []string{
		"Tab/Shift+Tab or Up/Down: Navigate fields",
		"Enter: Run command",
		"ESC/Ctrl+C: Close",
	}
	line := strings.Join(info, " | ")
	return lipgloss.NewStyle().Faint(true).Render(line)
}
