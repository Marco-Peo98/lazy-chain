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
		inputBuffer: cfg.WalletAddr,
	}
}

func (m *SettingsModel) Init() tea.Cmd {
	return nil
}

func (m *SettingsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			if !m.editingAddr && m.cursor > 0 {
				m.cursor--
			}
		case "down":
			if !m.editingAddr && m.cursor < len(m.networks)-1 {
				m.cursor++
			}
		case "enter":
			if m.editingAddr {
				m.config.WalletAddr = m.inputBuffer
				m.editingAddr = false
				SaveConfig(m.config)
			} else {
				m.config.Network = m.networks[m.cursor]
				SaveConfig(m.config)
			}
		case "e":
			m.editingAddr = true
			m.inputBuffer = m.config.WalletAddr
		case "esc":
			return m, tea.Quit
		}

		if m.editingAddr {
			switch msg.Type {
			case tea.KeyRunes:
				m.inputBuffer += string(msg.Runes)
			case tea.KeyBackspace:
				if len(m.inputBuffer) > 0 {
					m.inputBuffer = m.inputBuffer[:len(m.inputBuffer)-1]
				}
			}
		}
	}

	return m, nil
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
		s += "(Enter to save)"
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
