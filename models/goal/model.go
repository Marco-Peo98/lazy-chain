package goal

import (
	"context"
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"lazychain/models/goal/builders"
)

type GOALModel struct {
	// Transaction type selection
	txnTypes     []string
	selectedType int

	// Active builder
	builder  Builder
	builders map[string]Builder

	// Command execution
	runner *Runner
	output string
	errMsg string

	// UI state
	focusOnBuilder bool // true = fields, false = txn types
}

func NewGOALModel() *GOALModel {
	m := &GOALModel{
		runner:         NewRunner(),
		focusOnBuilder: false,
		builders:       make(map[string]Builder),
	}

	// Transaction types list
	m.txnTypes = []string{
		"Payment (pay)",
		"Application Call (appl)",
		"Asset Transfer (axfer)",
		"Asset Create (acfg)",
		"Asset Freeze (afrz)",
		"Key Registration (keyreg)",
	}

	// Initialize all builders
	m.initBuilders()

	// Set first builder as active
	m.builder = m.builders["pay"]

	return m
}

func (m *GOALModel) initBuilders() {
	// Payment
	pay := builders.NewPaymentBuilder()
	pay.SetRunFunc(m.run)
	m.builders["pay"] = pay

	// Application Call
	appl := builders.NewApplicationCallBuilder()
	appl.SetRunFunc(m.run)
	m.builders["appl"] = appl

	// Asset Transfer
	axfer := builders.NewAssetTransferBuilder()
	axfer.SetRunFunc(m.run)
	m.builders["axfer"] = axfer

	// Asset Create
	acfg := builders.NewAssetCreateBuilder()
	acfg.SetRunFunc(m.run)
	m.builders["acfg"] = acfg

	// Asset Freeze
	afrz := builders.NewAssetFreezeBuilder()
	afrz.SetRunFunc(m.run)
	m.builders["afrz"] = afrz

	// Key Registration
	keyreg := builders.NewKeyRegistrationBuilder()
	keyreg.SetRunFunc(m.run)
	m.builders["keyreg"] = keyreg
}

func (m *GOALModel) Init() tea.Cmd {
	return m.builder.Init()
}

func (m *GOALModel) run(argv []string) {
	_ = m.runner.CheckBinary()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	res := m.runner.Run(ctx, argv)
	m.output = strings.TrimSpace(res.Stdout)

	if res.Err != nil {
		if m.output == "" {
			m.errMsg = fmt.Sprintf("Error: %v\n%s", res.Err, strings.TrimSpace(res.Stderr))
		} else {
			m.errMsg = fmt.Sprintf("Error: %v", res.Err)
		}
	} else {
		m.errMsg = ""
	}

	m.builder.AfterRun(res.Stdout, res.Stderr, res.Err)
}

func (m *GOALModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if !m.focusOnBuilder {
			// Focus on transaction type list
			switch msg.String() {
			case "up", "k":
				if m.selectedType > 0 {
					m.selectedType--
					m.switchBuilder()
				}
			case "down", "j":
				if m.selectedType < len(m.txnTypes)-1 {
					m.selectedType++
					m.switchBuilder()
				}
			case "enter", "tab", "right":
				// Switch focus to builder
				m.focusOnBuilder = true
			}
			// Note: ESC when focus on txn types is handled by main.go (exits to ProjectView)
		} else {
			// Focus on builder (fields)
			switch msg.String() {
			case "tab", "left":
				// Only switch back if not editing
				if !m.builder.IsEditing() {
					m.focusOnBuilder = false
				} else {
					// If editing, delegate to builder
					var cmd tea.Cmd
					m.builder, cmd = m.builder.Update(msg)
					return m, cmd
				}
			case "esc":
				// If editing, let builder handle it (cancel edit)
				// If not editing, switch focus back to txn types
				if m.builder.IsEditing() {
					var cmd tea.Cmd
					m.builder, cmd = m.builder.Update(msg)
					return m, cmd
				}
				// Not editing, switch focus back to txn types
				m.focusOnBuilder = false
			default:
				// Delegate to builder
				var cmd tea.Cmd
				m.builder, cmd = m.builder.Update(msg)
				return m, cmd
			}
		}
	}

	return m, nil
}

func (m *GOALModel) switchBuilder() {
	// Map selectedType to builder key
	typeMap := map[int]string{
		0: "pay",
		1: "appl",
		2: "axfer",
		3: "acfg",
		4: "afrz",
		5: "keyreg",
	}

	if key, exists := typeMap[m.selectedType]; exists {
		m.builder = m.builders[key]
		m.output = ""
		m.errMsg = ""
	}
}

// View will be replaced by CmdGoalsLayout, but kept for now
func (m *GOALModel) View() string {
	return "GOALModel - Use CmdGoalsLayout to render"
}

// Public methods for layout to access data
func (m *GOALModel) GetTxnTypes() []string {
	return m.txnTypes
}

func (m *GOALModel) GetSelectedType() int {
	return m.selectedType
}

func (m *GOALModel) GetBuilder() Builder {
	return m.builder
}

func (m *GOALModel) GetOutput() string {
	if m.errMsg != "" {
		return m.errMsg
	}
	return m.output
}

func (m *GOALModel) IsFocusOnBuilder() bool {
	return m.focusOnBuilder
}

// CanExit returns true if ESC should exit to the previous screen
// Returns false if ESC should be handled internally (e.g., to cancel editing or switch focus)
func (m *GOALModel) CanExit() bool {
	// Can only exit when focus is on transaction types list (not on builder/fields)
	return !m.focusOnBuilder
}
