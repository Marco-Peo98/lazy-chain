package settings

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// SettingsModel holds UI state for settings screen.
type SettingsModel struct {
	config       Config
	networks     []string
	networkInfos map[string]NetworkInfo
	cursor       int
	editingAddr  bool   // Whether the user is editing the wallet address
	inputBuffer  string // Buffer for user input
	err          error

	// Network editing state
	editingNetwork bool
	editField      string      // Which field is being edited (url, port, token)
	networkBuffer  NetworkInfo // Buffer for network editing

	// Network management
	networkManager   *NetworkManager
	connectionStatus string // Status message for network connections
	showStatus       bool   // Whether to show connection status
}

func NewSettingsModel(networks []string) *SettingsModel {
	cfg, _ := LoadConfig()

	// Predefined networks
	networkInfos := map[string]NetworkInfo{
		"localnet": {
			Name:         "localnet",
			AlgodURL:     "http://localhost",
			AlgodPort:    "4001",
			AlgodToken:   strings.Repeat("a", 64),
			IndexerURL:   "http://localhost",
			IndexerPort:  "8980",
			IndexerToken: strings.Repeat("a", 64),
		},
		"testnet": {
			Name:         "testnet",
			AlgodURL:     "https://testnet-api.4160.nodely.dev",
			AlgodPort:    "443",
			AlgodToken:   "",
			IndexerURL:   "https://testnet-idx.4160.nodely.dev",
			IndexerPort:  "443",
			IndexerToken: "",
		},
		"mainnet": {
			Name:         "mainnet",
			AlgodURL:     "https://mainnet-api.4160.nodely.dev",
			AlgodPort:    "443",
			AlgodToken:   "",
			IndexerURL:   "https://mainnet-idx.4160.nodely.dev",
			IndexerPort:  "443",
			IndexerToken: "",
		},
	}

	// Copy networks slice
	availableNetworks := make([]string, len(networks))
	copy(availableNetworks, networks)

	// Add custom networks from config and override predefined ones if they exist
	for _, customNet := range cfg.CustomNetworks {
		networkInfos[customNet.Name] = customNet

		found := false
		for _, net := range availableNetworks {
			if net == customNet.Name {
				found = true
				break
			}
		}
		if !found {
			availableNetworks = append(availableNetworks, customNet.Name)
		}
	}

	// Cursor on current network
	cursor := 0
	for i, n := range availableNetworks {
		if n == cfg.Network {
			cursor = i
			break
		}
	}

	return &SettingsModel{
		config:           cfg,
		networks:         availableNetworks,
		networkInfos:     networkInfos,
		cursor:           cursor,
		editingAddr:      false,
		inputBuffer:      "",
		editingNetwork:   false,
		editField:        "",
		networkBuffer:    NetworkInfo{},
		networkManager:   NewNetworkManager(),
		connectionStatus: "",
		showStatus:       false,
	}
}

func (m *SettingsModel) Init() tea.Cmd { return nil }

// View assembles the two main panels + footer.
func (m *SettingsModel) View() string {
	leftColumn := m.renderNetworkSection()
	rightColumn := m.renderWalletSection()

	mainContent := lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftColumn,
		"  ",
		rightColumn,
	)
	footer := m.renderFooter()

	return lipgloss.JoinVertical(
		lipgloss.Left,
		mainContent,
		"",
		footer,
	)
}

// Public helpers used by other components.

func (m *SettingsModel) GetNetworkManager() *NetworkManager { return m.networkManager }
func (m *SettingsModel) IsEditingAddr() bool                { return m.editingAddr || m.editingNetwork }

// ResetEditingState ensures the settings model is not in editing mode.
func (m *SettingsModel) ResetEditingState() {
	m.editingAddr = false
	m.inputBuffer = ""
	m.editingNetwork = false
	m.networkBuffer = NetworkInfo{}
	m.showStatus = false
	m.connectionStatus = ""
}

// refreshNetworkInfo ensures network info is up-to-date from config.
func (m *SettingsModel) refreshNetworkInfo(networkName string) {
	// Prefer custom networks from config
	for _, customNet := range m.config.CustomNetworks {
		if customNet.Name == networkName {
			m.networkInfos[networkName] = customNet
			return
		}
	}

	// Re-add predefined if missing
	if _, exists := m.networkInfos[networkName]; !exists {
		switch networkName {
		case "localnet":
			m.networkInfos[networkName] = NetworkInfo{
				Name:         "localnet",
				AlgodURL:     "http://localhost",
				AlgodPort:    "4001",
				AlgodToken:   strings.Repeat("a", 64),
				IndexerURL:   "http://localhost",
				IndexerPort:  "8980",
				IndexerToken: strings.Repeat("a", 64),
			}
		case "testnet":
			// Mirrors original file's fallback (nodely).
			m.networkInfos[networkName] = NetworkInfo{
				Name:         "testnet",
				AlgodURL:     "https://testnet-api.4160.nodely.dev",
				AlgodPort:    "443",
				AlgodToken:   "",
				IndexerURL:   "https://testnet-idx.4160.nodely.dev",
				IndexerPort:  "443",
				IndexerToken: "",
			}
		case "mainnet":
			m.networkInfos[networkName] = NetworkInfo{
				Name:         "mainnet",
				AlgodURL:     "https://mainnet-api.4160.nodely.dev",
				AlgodPort:    "443",
				AlgodToken:   "",
				IndexerURL:   "https://mainnet-idx.4160.nodely.dev",
				IndexerPort:  "443",
				IndexerToken: "",
			}
		}
	}
}
