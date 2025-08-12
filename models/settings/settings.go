package settings

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Config represents the configuration settings for the application
type Config struct {
	Network    string `json:"network"`
	WalletAddr string `json:"wallet_addr"`
}

type SettingsModel struct {
	config      Config
	networks    []string
	cursor      int
	editingAddr bool   // Whether the user is editing the wallet address
	inputBuffer string // Buffer for user input
	err         error
}

func NewSettingsModel(networks []string) *SettingsModel {
	cfg, _ := LoadConfig()

	cursor := 0
	for i, n := range networks {
		if n == cfg.Network {
			cursor = i
			break
		}
	}

	return &SettingsModel{
		config:      cfg,
		networks:    networks,
		cursor:      cursor,
		editingAddr: false,
		inputBuffer: "",
	}
}

func (m *SettingsModel) Init() tea.Cmd {
	return nil
}

func (m *SettingsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle input when in editing mode
		if m.editingAddr {
			switch msg.String() {
			case "enter":
				// Save the address and exit editing mode
				m.config.WalletAddr = m.inputBuffer
				m.editingAddr = false
				SaveConfig(m.config)
				return m, nil
			case "esc":
				// Cancel editing and restore previous value
				m.editingAddr = false
				m.inputBuffer = ""
				return m, nil
			}

			// Handle character input
			switch msg.Type {
			case tea.KeyRunes:
				m.inputBuffer += string(msg.Runes)
			case tea.KeyBackspace:
				if len(m.inputBuffer) > 0 {
					m.inputBuffer = m.inputBuffer[:len(m.inputBuffer)-1]
				}
			}
			return m, nil
		}

		// Handle input when NOT in editing mode
		switch msg.String() {
		case "up":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down":
			if m.cursor < len(m.networks)-1 {
				m.cursor++
			}
		case "enter":
			// Select network
			m.config.Network = m.networks[m.cursor]
			SaveConfig(m.config)
		case "e":
			// Enter editing mode for wallet address
			m.editingAddr = true
			m.inputBuffer = m.config.WalletAddr // Start with current address
		}
	}

	return m, nil
}

// ResetEditingState ensures the settings model is not in editing mode
// This should be called when entering or leaving the settings view
func (m *SettingsModel) ResetEditingState() {
	m.editingAddr = false
	m.inputBuffer = ""
}

// IsEditingAddr returns true if the model is currently in address editing mode
func (m *SettingsModel) IsEditingAddr() bool {
	return m.editingAddr
}

func (m *SettingsModel) View() string {
	// Create two-column layout: Network selection (left) | Wallet address (right)
	leftColumn := m.renderNetworkSection()
	rightColumn := m.renderWalletSection()

	// Combine columns horizontally
	mainContent := lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftColumn,
		"  ", // Small gap between columns
		rightColumn,
	)

	// Add footer with instructions
	footer := m.renderFooter()

	// Combine everything vertically
	return lipgloss.JoinVertical(
		lipgloss.Left,
		mainContent,
		"",
		footer,
	)
}

func (m *SettingsModel) renderNetworkSection() string {
	var content []string

	// Section title
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#cba6f7")).
		Render("Network Selection")
	content = append(content, title)
	content = append(content, "")

	// Network list
	for i, network := range m.networks {
		cursor := "  "
		style := lipgloss.NewStyle()

		// Show cursor only when not editing address
		if m.cursor == i && !m.editingAddr {
			cursor = "> "
			style = style.Foreground(lipgloss.Color("#ef9f76"))
		}

		// Highlight selected network
		if network == m.config.Network {
			style = style.Bold(true).Foreground(lipgloss.Color("#a6e3a1"))
			network += " (selected)"
		}

		line := cursor + style.Render(network)
		content = append(content, line)
	}

	// Add placeholder text
	content = append(content, "")
	placeholder := lipgloss.NewStyle().
		Italic(true).
		Foreground(lipgloss.Color("#6c7086")).
		Render("Network details will appear here")
	content = append(content, placeholder)

	// Fill remaining space to match ProjectModel height
	for len(content) < 10 {
		content = append(content, "")
	}

	// Create left panel with border (standardized dimensions with increased padding)
	panelContent := strings.Join(content, "\n")
	return lipgloss.NewStyle().
		Width(30).  // Standardized width
		Padding(2). // Increased padding
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#89b4fa")).
		Render(panelContent)
}

func (m *SettingsModel) renderWalletSection() string {
	var content []string

	// Section title
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#cba6f7")).
		Render("Wallet Address")
	content = append(content, title)
	content = append(content, "")

	if m.editingAddr {
		// Editing mode
		content = append(content, "Editing address:")
		content = append(content, "")

		// Input field (fixed to not deform panel borders)
		inputDisplay := m.inputBuffer + "_" // Cursor indicator

		// Style input without separate border to avoid deformation
		inputStyle := lipgloss.NewStyle().
			Width(37). // Adjusted for increased padding (45-4-4 for padding)
			Background(lipgloss.Color("#1e1e2e")).
			Foreground(lipgloss.Color("#f9e2af")).
			Padding(0, 1)

		styledInput := inputStyle.Render(inputDisplay)
		content = append(content, styledInput)

		content = append(content, "")
		instructions := lipgloss.NewStyle().
			Italic(true).
			Foreground(lipgloss.Color("#f9e2af")).
			Render("Enter to save, ESC to cancel")
		content = append(content, instructions)

	} else {
		// Display mode
		content = append(content, "Current address:")
		content = append(content, "")

		addr := m.config.WalletAddr
		if addr == "" {
			addr = "Not configured"
			addrStyle := lipgloss.NewStyle().
				Italic(true).
				Foreground(lipgloss.Color("#f38ba8"))
			addr = addrStyle.Render(addr)
		} else {
			// Truncate long addresses for display
			if len(addr) > 35 {
				addr = addr[:15] + "..." + addr[len(addr)-15:]
			}
			addrStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color("#a6e3a1"))
			addr = addrStyle.Render(addr)
		}

		content = append(content, addr)
		content = append(content, "")

		editInstr := lipgloss.NewStyle().
			Italic(true).
			Foreground(lipgloss.Color("#cdd6f4")).
			Render("Press 'e' to edit")
		content = append(content, editInstr)
	}

	// Fill remaining space to match ProjectModel height
	for len(content) < 10 {
		content = append(content, "")
	}

	// Create right panel with border (standardized dimensions with increased padding)
	panelContent := strings.Join(content, "\n")
	return lipgloss.NewStyle().
		Width(45).  // Standardized width
		Padding(2). // Increased padding
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#a6e3a1")).
		Render(panelContent)
}

func (m *SettingsModel) renderFooter() string {
	var instructions []string

	if m.editingAddr {
		instructions = []string{
			"Enter: Save address",
			"ESC: Cancel editing",
		}
	} else {
		instructions = []string{
			"Up/Down: Navigate networks",
			"Enter: Select network",
			"e: Edit wallet address",
			"ESC: Back to menu",
		}
	}

	instructionText := strings.Join(instructions, " | ")
	return lipgloss.NewStyle().
		Width(77).              // Total width of both panels + gap (30+2+45)
		Align(lipgloss.Center). // Center the footer text
		Foreground(lipgloss.Color("#6c7086")).
		Italic(true).
		Render(instructionText)
}

// Configuration file management functions remain unchanged
func ConfigPath() string {
	home, _ := os.UserHomeDir()
	return fmt.Sprintf("%s/.lazy-chain/config.json", home)
}

func LoadConfig() (Config, error) {
	path := ConfigPath()
	buff, err := ioutil.ReadFile(path)
	if err != nil {
		return Config{Network: "testnet", WalletAddr: ""}, nil // Default config
	}
	var cfg Config
	if err := json.Unmarshal(buff, &cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func SaveConfig(cfg Config) error {
	buff, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(ConfigPath(), buff, 0644)
}
