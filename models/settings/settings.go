package settings

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

// Config represents the configuration settings for the application.
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
			// ESC is NOT handled here when not in editing mode
			// This allows the MainModel to handle it and change views
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
	s := "Settings\n\n"

	s += "Network:\n"
	for i, n := range m.networks {
		cursor := " "
		if m.cursor == i && !m.editingAddr {
			cursor = ">"
		}
		s += fmt.Sprintf("%s %s\n", cursor, n)
	}

	if m.editingAddr {
		s += fmt.Sprintf("\nWallet Address: %s_\n", m.inputBuffer)
		s += "(Enter to save, ESC to cancel)"
	} else {
		s += fmt.Sprintf("\nWallet Address: %s\n", m.config.WalletAddr)
		s += "Press 'e' to edit"
	}
	return s
}

func ConfigPath() string {
	home, _ := os.UserHomeDir()
	return fmt.Sprintf("%s/.lazy-chain/config.json", home)
}

func LoadConfig() (Config, error) {
	path := ConfigPath()
	buff, err := ioutil.ReadFile(path)
	if err != nil {
		return Config{Network: "Testnet", WalletAddr: ""}, nil // Default config
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
