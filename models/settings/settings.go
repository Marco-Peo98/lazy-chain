package settings

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/client/v2/indexer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// NetworkInfo represents network connection details
type NetworkManager struct {
	algodClient    *algod.Client
	indexerClient  *indexer.Client
	currentNetwork NetworkInfo
	connected      bool
}

// NewNetworkManager creates a new network manager
func NewNetworkManager() *NetworkManager {
	return &NetworkManager{
		connected: false,
	}
}

// NetworkStatus represents the current network status
type NetworkStatus struct {
	NetworkName    string    `json:"network_name"`
	AlgodURL       string    `json:"algod_url"`
	IndexerURL     string    `json:"indexer_url"`
	Connected      bool      `json:"connected"`
	LastRound      uint64    `json:"last_round"`
	GenesisID      string    `json:"genesis_id"`
	GenesisHash    string    `json:"genesis_hash"`
	IndexerHealthy bool      `json:"indexer_healthy"`
	Timestamp      time.Time `json:"timestamp"`
}

// createAlgodClient creates an algod client from network info
func createAlgodClient(networkInfo NetworkInfo) (*algod.Client, error) {
	url := fmt.Sprintf("%s:%s", networkInfo.AlgodURL, networkInfo.AlgodPort)
	// Remove port from URL if it already contains it (for URLs like https://testnet-api.algonode.cloud/)
	if strings.Contains(networkInfo.AlgodURL, "://") && networkInfo.AlgodPort == "443" {
		url = networkInfo.AlgodURL
	}

	client, err := algod.MakeClient(url, networkInfo.AlgodToken)
	return client, err
}

// createIndexerClient creates an indexer client from network info
func createIndexerClient(networkInfo NetworkInfo) (*indexer.Client, error) {
	if networkInfo.IndexerURL == "" {
		return nil, fmt.Errorf("indexer URL not provided")
	}

	url := fmt.Sprintf("%s:%s", networkInfo.IndexerURL, networkInfo.IndexerPort)
	// Remove port from URL if it already contains it
	if strings.Contains(networkInfo.IndexerURL, "://") && networkInfo.IndexerPort == "443" {
		url = networkInfo.IndexerURL
	}

	client, err := indexer.MakeClient(url, networkInfo.IndexerToken)
	return client, err
}

// TestNetworkConnection tests connection to a network without storing it
func (nm *NetworkManager) TestNetworkConnection(networkInfo NetworkInfo) error {
	// Test algod connection
	algodClient, err := createAlgodClient(networkInfo)
	if err != nil {
		return fmt.Errorf("failed to create algod client: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = algodClient.Status().Do(ctx)
	if err != nil {
		return fmt.Errorf("algod connection test failed: %w", err)
	}

	// Test indexer connection if provided
	if networkInfo.IndexerURL != "" {
		indexerClient, err := createIndexerClient(networkInfo)
		if err != nil {
			return fmt.Errorf("failed to create indexer client: %w", err)
		}

		_, err = indexerClient.SearchForApplications().Limit(1).Do(ctx)
		if err != nil {
			return fmt.Errorf("indexer connection test failed: %w", err)
		}
	}

	return nil
}

// ConnectToNetwork establishes connection to the specified network
func (nm *NetworkManager) ConnectToNetwork(networkInfo NetworkInfo) error {
	// Create algod client
	algodClient, err := createAlgodClient(networkInfo)
	if err != nil {
		return fmt.Errorf("failed to connect to algod: %w", err)
	}

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = algodClient.Status().Do(ctx)
	if err != nil {
		return fmt.Errorf("failed to get network status: %w", err)
	}

	// Connect to Indexer if URL is provided
	var indexerClient *indexer.Client
	if networkInfo.IndexerURL != "" {
		indexerClient, err = createIndexerClient(networkInfo)
		if err != nil {
			// Indexer connection is optional, log warning but don't fail
			fmt.Printf("Warning: failed to connect to indexer: %v\n", err)
		} else {
			// Test indexer connection
			_, err = indexerClient.SearchForApplications().Limit(1).Do(ctx)
			if err != nil {
				fmt.Printf("Warning: indexer connection test failed: %v\n", err)
				indexerClient = nil
			}
		}
	}

	// Store successful connections
	nm.algodClient = algodClient
	nm.indexerClient = indexerClient
	nm.currentNetwork = networkInfo
	nm.connected = true

	return nil
}

// GetNetworkStatus returns current network connection status
func (nm *NetworkManager) GetNetworkStatus() (NetworkStatus, error) {
	if !nm.connected || nm.algodClient == nil {
		return NetworkStatus{}, fmt.Errorf("not connected to any network")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Get algod status
	status, err := nm.algodClient.Status().Do(ctx)
	if err != nil {
		return NetworkStatus{}, fmt.Errorf("failed to get network status: %w", err)
	}

	// Get network genesis information
	genesisResponse, err := nm.algodClient.GetGenesis().Do(ctx)
	if err != nil {
		return NetworkStatus{}, fmt.Errorf("failed to get genesis info: %w", err)
	}

	// Check indexer status if available
	indexerHealthy := false
	if nm.indexerClient != nil {
		_, err := nm.indexerClient.SearchForApplications().Limit(1).Do(ctx)
		indexerHealthy = (err == nil)
	}

	// Extract genesis info - genesis response is a string
	genesisID := "unknown"
	genesisHash := "unknown"
	if len(genesisResponse) > 0 {
		genesisID = "genesis-available"
		genesisHash = nm.currentNetwork.Name
	}

	return NetworkStatus{
		NetworkName:    nm.currentNetwork.Name,
		AlgodURL:       nm.currentNetwork.AlgodURL,
		IndexerURL:     nm.currentNetwork.IndexerURL,
		Connected:      true,
		LastRound:      status.LastRound,
		GenesisID:      genesisID,
		GenesisHash:    genesisHash,
		IndexerHealthy: indexerHealthy,
		Timestamp:      time.Now(),
	}, nil
}

// IsConnected returns whether we're currently connected to a network
func (nm *NetworkManager) IsConnected() bool {
	return nm.connected && nm.algodClient != nil
}

// GetCurrentNetwork returns information about the currently connected network
func (nm *NetworkManager) GetCurrentNetwork() NetworkInfo {
	return nm.currentNetwork
}

// GetAlgodClient returns the current algod client (for other components to use)
func (nm *NetworkManager) GetAlgodClient() *algod.Client {
	return nm.algodClient
}

// GetIndexerClient returns the current indexer client (for other components to use)
func (nm *NetworkManager) GetIndexerClient() *indexer.Client {
	return nm.indexerClient
}

// Disconnect closes current network connections
func (nm *NetworkManager) Disconnect() {
	nm.algodClient = nil
	nm.indexerClient = nil
	nm.connected = false
	nm.currentNetwork = NetworkInfo{}
}

// ValidateNetworkConfig validates network configuration before connection
func ValidateNetworkConfig(networkInfo NetworkInfo) error {
	if networkInfo.Name == "" {
		return fmt.Errorf("network name cannot be empty")
	}

	if networkInfo.AlgodURL == "" {
		return fmt.Errorf("algod URL cannot be empty")
	}

	if networkInfo.AlgodPort == "" {
		return fmt.Errorf("algod port cannot be empty")
	}

	// Indexer is optional, but if URL is provided, port should be too
	if networkInfo.IndexerURL != "" && networkInfo.IndexerPort == "" {
		return fmt.Errorf("indexer port required when indexer URL is provided")
	}

	return nil
}

type NetworkInfo struct {
	Name         string `json:"name"`
	AlgodURL     string `json:"algod_url"`
	AlgodPort    string `json:"algod_port"`
	AlgodToken   string `json:"algod_token"`
	IndexerURL   string `json:"indexer_url"`
	IndexerPort  string `json:"indexer_port"`
	IndexerToken string `json:"indexer_token"`
}

// Config represents the configuration settings for the application
type Config struct {
	Network        string        `json:"network"`
	WalletAddr     string        `json:"wallet_addr"`
	CustomNetworks []NetworkInfo `json:"custom_networks"`
}

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

	// Initialize predefined network information from network.go constants
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
			AlgodURL:     "https://testnet-api.algonode.cloud/",
			AlgodPort:    "443",
			AlgodToken:   "",
			IndexerURL:   "https://testnet-idx.algonode.cloud/",
			IndexerPort:  "443",
			IndexerToken: "",
		},
		"mainnet": {
			Name:         "mainnet",
			AlgodURL:     "https://mainnet-api.algonode.cloud/",
			AlgodPort:    "443",
			AlgodToken:   "",
			IndexerURL:   "https://mainnet-idx.algonode.cloud/",
			IndexerPort:  "443",
			IndexerToken: "",
		},
	}

	// Create a copy of the networks slice to avoid modifying the original
	availableNetworks := make([]string, len(networks))
	copy(availableNetworks, networks)

	// Add custom networks from config and override predefined ones if they exist
	for _, customNet := range cfg.CustomNetworks {
		// Always update the network info (this handles both new networks and overrides)
		networkInfos[customNet.Name] = customNet

		// Add to networks list if not already present
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

	// Find cursor position for the current network
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

func (m *SettingsModel) Init() tea.Cmd {
	return nil
}

func (m *SettingsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle input when editing network - HIGHEST PRIORITY
		if m.editingNetwork {
			return m.handleNetworkEditing(msg)
		}

		// Handle input when in wallet address editing mode - SECOND PRIORITY
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

		// Handle input when NOT in any editing mode - LOWEST PRIORITY
		switch msg.String() {
		case "up":
			if m.cursor > 0 {
				m.cursor--
				// Clear status when changing selection
				m.showStatus = false
				m.connectionStatus = ""
			}
		case "down":
			if m.cursor < len(m.networks)-1 {
				m.cursor++
				// Clear status when changing selection
				m.showStatus = false
				m.connectionStatus = ""
			}
		case "enter":
			// Select and connect to network - ONLY when not editing
			if !m.editingNetwork && !m.editingAddr {
				selectedNetwork := m.networks[m.cursor]
				if networkInfo, exists := m.networkInfos[selectedNetwork]; exists {
					// Test connection first
					err := m.networkManager.TestNetworkConnection(networkInfo)
					if err != nil {
						m.connectionStatus = fmt.Sprintf("Failed to connect to %s: %v", selectedNetwork, err)
						m.showStatus = true
					} else {
						// Connection successful, save and connect
						m.config.Network = selectedNetwork
						SaveConfig(m.config)

						err = m.networkManager.ConnectToNetwork(networkInfo)
						if err != nil {
							m.connectionStatus = fmt.Sprintf("Connection failed: %v", err)
						} else {
							m.connectionStatus = fmt.Sprintf("Successfully connected to %s", selectedNetwork)
						}
						m.showStatus = true
					}
				}
			}
		case "t":
			// Test connection to current network without connecting - ONLY when not editing
			if !m.editingNetwork && !m.editingAddr && m.cursor < len(m.networks) {
				selectedNetwork := m.networks[m.cursor]
				if networkInfo, exists := m.networkInfos[selectedNetwork]; exists {
					err := m.networkManager.TestNetworkConnection(networkInfo)
					if err != nil {
						m.connectionStatus = fmt.Sprintf("Test failed for %s: %v", selectedNetwork, err)
					} else {
						m.connectionStatus = fmt.Sprintf("Test successful for %s", selectedNetwork)
					}
					m.showStatus = true
				}
			}
		case " ":
			// Hide status message - ONLY when not editing
			if !m.editingNetwork && !m.editingAddr && m.showStatus {
				m.showStatus = false
				m.connectionStatus = ""
			}
		case "e":
			// Enter editing mode for wallet address - ONLY when not editing network
			if !m.editingNetwork {
				m.editingAddr = true
				m.inputBuffer = m.config.WalletAddr // Start with current address
			}
		case "n":
			// Enter network editing mode for current network - ONLY when not in any editing mode
			if !m.editingNetwork && !m.editingAddr {
				currentNetworkName := m.networks[m.cursor]
				if info, exists := m.networkInfos[currentNetworkName]; exists {
					m.editingNetwork = true
					m.networkBuffer = info
					m.editField = "algod_url"
					// Clear any status messages when entering edit mode
					m.showStatus = false
					m.connectionStatus = ""
				}
			}
		case "c":
			// Create new custom network - ONLY when not in any editing mode
			if !m.editingNetwork && !m.editingAddr {
				m.editingNetwork = true
				m.networkBuffer = NetworkInfo{
					Name:         "custom",
					AlgodURL:     "http://localhost",
					AlgodPort:    "4001",
					AlgodToken:   "",
					IndexerURL:   "http://localhost",
					IndexerPort:  "8980",
					IndexerToken: "",
				}
				m.editField = "name"
				// Clear any status messages when entering edit mode
				m.showStatus = false
				m.connectionStatus = ""
			}
		}
	}

	return m, nil
}

func (m *SettingsModel) handleNetworkEditing(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "tab":
		// Move to next field
		m.editField = m.getNextField(m.editField)
	case "shift+tab":
		// Move to previous field
		m.editField = m.getPrevField(m.editField)
	case "enter":
		// Save network configuration
		m.saveNetworkConfig()
		m.editingNetwork = false
		return m, nil
	case "esc":
		// Cancel network editing
		m.editingNetwork = false
		m.networkBuffer = NetworkInfo{}
		return m, nil
	}

	// Handle character input for current field
	switch msg.Type {
	case tea.KeyRunes:
		m.setNetworkField(m.editField, m.getNetworkField(m.editField)+string(msg.Runes))
	case tea.KeyBackspace:
		current := m.getNetworkField(m.editField)
		if len(current) > 0 {
			m.setNetworkField(m.editField, current[:len(current)-1])
		}
	}

	return m, nil
}

func (m *SettingsModel) getNextField(current string) string {
	fields := []string{"name", "algod_url", "algod_port", "algod_token", "indexer_url", "indexer_port", "indexer_token"}
	for i, field := range fields {
		if field == current {
			if i < len(fields)-1 {
				return fields[i+1]
			}
			return fields[0]
		}
	}
	return fields[0]
}

func (m *SettingsModel) getPrevField(current string) string {
	fields := []string{"name", "algod_url", "algod_port", "algod_token", "indexer_url", "indexer_port", "indexer_token"}
	for i, field := range fields {
		if field == current {
			if i > 0 {
				return fields[i-1]
			}
			return fields[len(fields)-1]
		}
	}
	return fields[len(fields)-1]
}

func (m *SettingsModel) getNetworkField(field string) string {
	switch field {
	case "name":
		return m.networkBuffer.Name
	case "algod_url":
		return m.networkBuffer.AlgodURL
	case "algod_port":
		return m.networkBuffer.AlgodPort
	case "algod_token":
		return m.networkBuffer.AlgodToken
	case "indexer_url":
		return m.networkBuffer.IndexerURL
	case "indexer_port":
		return m.networkBuffer.IndexerPort
	case "indexer_token":
		return m.networkBuffer.IndexerToken
	default:
		return ""
	}
}

func (m *SettingsModel) setNetworkField(field, value string) {
	switch field {
	case "name":
		m.networkBuffer.Name = value
	case "algod_url":
		m.networkBuffer.AlgodURL = value
	case "algod_port":
		m.networkBuffer.AlgodPort = value
	case "algod_token":
		m.networkBuffer.AlgodToken = value
	case "indexer_url":
		m.networkBuffer.IndexerURL = value
	case "indexer_port":
		m.networkBuffer.IndexerPort = value
	case "indexer_token":
		m.networkBuffer.IndexerToken = value
	}
}

func (m *SettingsModel) saveNetworkConfig() {
	// Validate network configuration
	err := ValidateNetworkConfig(m.networkBuffer)
	if err != nil {
		m.connectionStatus = fmt.Sprintf("Validation failed: %v", err)
		m.showStatus = true
		return
	}

	// Test connection before saving
	err = m.networkManager.TestNetworkConnection(m.networkBuffer)
	if err != nil {
		m.connectionStatus = fmt.Sprintf("Connection test failed: %v", err)
		m.showStatus = true
		return
	}

	// Update network info map IMMEDIATELY (this is key for immediate UI update)
	m.networkInfos[m.networkBuffer.Name] = m.networkBuffer

	// Check if this is a new network or modification of existing one
	isNewNetwork := true
	networkIndex := -1

	// Check if network already exists in the networks list
	for i, existingNetwork := range m.networks {
		if existingNetwork == m.networkBuffer.Name {
			isNewNetwork = false
			networkIndex = i
			break
		}
	}

	// Add or update network in custom networks config
	// Even predefined networks get saved when modified
	found := false
	for i, net := range m.config.CustomNetworks {
		if net.Name == m.networkBuffer.Name {
			m.config.CustomNetworks[i] = m.networkBuffer
			found = true
			break
		}
	}

	if !found {
		m.config.CustomNetworks = append(m.config.CustomNetworks, m.networkBuffer)
	}

	// Add to networks list if it's a new network
	if isNewNetwork {
		m.networks = append(m.networks, m.networkBuffer.Name)
		// Update cursor to point to the new network
		m.cursor = len(m.networks) - 1
	} else {
		// If modifying existing network, keep cursor at current position
		m.cursor = networkIndex
	}

	// Save configuration to file
	err = SaveConfig(m.config)
	if err != nil {
		m.connectionStatus = fmt.Sprintf("Failed to save config: %v", err)
		m.showStatus = true
		return
	}

	// Force reload configuration to ensure everything is synced
	reloadedConfig, err := LoadConfig()
	if err != nil {
		m.connectionStatus = fmt.Sprintf("Warning: Failed to reload config: %v", err)
		m.showStatus = true
	} else {
		// Update our config with the reloaded one to ensure synchronization
		m.config = reloadedConfig

		// Re-sync networkInfos with the reloaded custom networks
		for _, customNet := range m.config.CustomNetworks {
			m.networkInfos[customNet.Name] = customNet
		}
	}

	m.connectionStatus = fmt.Sprintf("Network %s saved successfully", m.networkBuffer.Name)
	m.showStatus = true
}

// ResetEditingState ensures the settings model is not in editing mode
// This should be called when entering or leaving the settings view
func (m *SettingsModel) ResetEditingState() {
	m.editingAddr = false
	m.inputBuffer = ""
	m.editingNetwork = false
	m.networkBuffer = NetworkInfo{}
	m.showStatus = false
	m.connectionStatus = ""
}

// refreshNetworkInfo ensures network info is up-to-date from config
func (m *SettingsModel) refreshNetworkInfo(networkName string) {
	// Check if this network exists in custom networks (which have priority)
	for _, customNet := range m.config.CustomNetworks {
		if customNet.Name == networkName {
			// Update the networkInfos map with the latest from config
			m.networkInfos[networkName] = customNet
			return
		}
	}

	// If not found in custom networks, ensure predefined network info is available
	// This handles the case where a predefined network might have been cleared accidentally
	if _, exists := m.networkInfos[networkName]; !exists {
		// Re-add predefined network info if missing
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
			m.networkInfos[networkName] = NetworkInfo{
				Name:         "testnet",
				AlgodURL:     "https://testnet-api.algonode.cloud/",
				AlgodPort:    "443",
				AlgodToken:   "",
				IndexerURL:   "https://testnet-idx.algonode.cloud/",
				IndexerPort:  "443",
				IndexerToken: "",
			}
		case "mainnet":
			m.networkInfos[networkName] = NetworkInfo{
				Name:         "mainnet",
				AlgodURL:     "https://mainnet-api.algonode.cloud/",
				AlgodPort:    "443",
				AlgodToken:   "",
				IndexerURL:   "https://mainnet-idx.algonode.cloud/",
				IndexerPort:  "443",
				IndexerToken: "",
			}
		}
	}
}

// GetNetworkManager returns the network manager for use by other components
func (m *SettingsModel) GetNetworkManager() *NetworkManager {
	return m.networkManager
}

// IsEditingAddr returns true if the model is currently in address editing mode
func (m *SettingsModel) IsEditingAddr() bool {
	return m.editingAddr || m.editingNetwork
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
	if m.editingNetwork {
		return m.renderNetworkEditor()
	}

	var content []string

	// Section title with extra spacing
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#cba6f7")).
		Render("Network Selection")
	content = append(content, title)
	content = append(content, "") // Extra spacing after title

	// Network list with improved spacing
	for i, network := range m.networks {
		cursor := "  "
		style := lipgloss.NewStyle()

		// Show cursor only when not editing address
		if m.cursor == i && !m.editingAddr {
			cursor = "> "
			style = style.Foreground(lipgloss.Color("#ef9f76"))
		}

		// Highlight selected network
		displayName := network
		if network == m.config.Network {
			style = style.Bold(true).Foreground(lipgloss.Color("#a6e3a1"))
			displayName += " (active)"
		}

		// Add indicator for custom/modified networks
		isCustom := false
		for _, customNet := range m.config.CustomNetworks {
			if customNet.Name == network {
				isCustom = true
				break
			}
		}
		if isCustom && network != m.config.Network {
			displayName += " *" // Asterisk indicates custom/modified
		}

		line := cursor + style.Render(displayName)
		content = append(content, line)
	}

	content = append(content, "") // Spacing after network list

	// Add legend for asterisk
	hasCustomNetworks := false
	for _, network := range m.networks {
		for _, customNet := range m.config.CustomNetworks {
			if customNet.Name == network && network != m.config.Network {
				hasCustomNetworks = true
				break
			}
		}
		if hasCustomNetworks {
			break
		}
	}

	if hasCustomNetworks {
		legendStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6c7086")).
			Italic(true)
		content = append(content, legendStyle.Render("* = custom/modified"))
	}

	content = append(content, "") // Extra spacing before details

	// Show network details for selected network with improved formatting
	// Always read from the most current networkInfos map
	if m.cursor < len(m.networks) {
		currentNet := m.networks[m.cursor]

		// Force refresh of network info from config if it's a custom network
		m.refreshNetworkInfo(currentNet)

		if info, exists := m.networkInfos[currentNet]; exists {
			// Section header
			detailsHeader := lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#89b4fa")).
				Render("Network Details:")
			content = append(content, detailsHeader)
			content = append(content, "") // Spacing after header

			// Algod information with truncation for long URLs
			algodStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#cdd6f4"))

			// Handle long URLs better
			algodURL := info.AlgodURL
			if len(algodURL) > 22 {
				// Split URL display
				content = append(content, algodStyle.Render("Algod:"))
				content = append(content, algodStyle.Render("  "+algodURL[:20]+"..."))
				content = append(content, algodStyle.Render("  Port: "+info.AlgodPort))
			} else {
				algodInfo := fmt.Sprintf("Algod: %s:%s", algodURL, info.AlgodPort)
				content = append(content, algodStyle.Render(algodInfo))
			}

			// Token information with better formatting
			if info.AlgodToken != "" {
				tokenDisplay := info.AlgodToken
				if len(tokenDisplay) > 16 {
					tokenDisplay = tokenDisplay[:6] + "..." + tokenDisplay[len(tokenDisplay)-6:]
				}
				content = append(content, algodStyle.Render("Token: "+tokenDisplay))
			} else {
				content = append(content, algodStyle.Render("Token: (none)"))
			}

			content = append(content, "") // Spacing before status

			// Connection status with color coding
			if m.networkManager.IsConnected() && m.networkManager.GetCurrentNetwork().Name == currentNet {
				statusStyle := lipgloss.NewStyle().
					Foreground(lipgloss.Color("#a6e3a1")).
					Bold(true)
				content = append(content, statusStyle.Render("Status: Connected"))
			} else {
				statusStyle := lipgloss.NewStyle().
					Foreground(lipgloss.Color("#6c7086")).
					Italic(true)
				content = append(content, statusStyle.Render("Status: Not connected"))
			}
		} else {
			// Fallback if network info doesn't exist
			errorStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color("#f38ba8")).
				Italic(true)
			content = append(content, errorStyle.Render("Network details unavailable"))
		}
	}

	content = append(content, "") // Spacing before status messages

	// Show connection status message if available with improved formatting
	if m.showStatus && m.connectionStatus != "" {
		statusStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#f9e2af")).
			Italic(true)

		// Split long status messages with better line breaking
		if len(m.connectionStatus) > 22 {
			words := strings.Fields(m.connectionStatus)
			var lines []string
			var currentLine []string
			currentLength := 0

			for _, word := range words {
				if currentLength+len(word)+len(currentLine) > 22 {
					if len(currentLine) > 0 {
						lines = append(lines, strings.Join(currentLine, " "))
						currentLine = []string{word}
						currentLength = len(word)
					} else {
						// Single long word
						lines = append(lines, word)
						currentLine = []string{}
						currentLength = 0
					}
				} else {
					currentLine = append(currentLine, word)
					currentLength += len(word)
				}
			}
			if len(currentLine) > 0 {
				lines = append(lines, strings.Join(currentLine, " "))
			}

			for _, line := range lines {
				content = append(content, statusStyle.Render(line))
			}
		} else {
			content = append(content, statusStyle.Render(m.connectionStatus))
		}
	}

	// Fill remaining space to match ProjectModel height
	for len(content) < 12 { // Increased from 10 to accommodate spacing
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

func (m *SettingsModel) renderNetworkEditor() string {
	var content []string

	// Section title with better spacing
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#f9e2af")).
		Render("Edit Network")
	content = append(content, title)
	content = append(content, "") // Extra spacing after title

	// Network fields with improved formatting
	fields := []struct {
		key   string
		label string
	}{
		{"name", "Name"},
		{"algod_url", "Algod URL"},
		{"algod_port", "Algod Port"},
		{"algod_token", "Algod Token"},
		{"indexer_url", "Indexer URL"},
		{"indexer_port", "Indexer Port"},
		{"indexer_token", "Indexer Token"},
	}

	for i, field := range fields {
		value := m.getNetworkField(field.key)

		// Style based on whether this is the active field
		if field.key == m.editField {
			// Active field - highlighted with better formatting
			fieldStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color("#f9e2af")).
				Bold(true)
			content = append(content, fieldStyle.Render("● "+field.label+":"))

			// Show input with cursor and proper formatting
			inputValue := value + "█" // Use block cursor
			inputStyle := lipgloss.NewStyle().
				Background(lipgloss.Color("#1e1e2e")).
				Foreground(lipgloss.Color("#f9e2af")).
				Padding(0, 1)
			content = append(content, "  "+inputStyle.Render(inputValue))
		} else {
			// Inactive field with subtle styling
			fieldStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color("#6c7086"))
			content = append(content, fieldStyle.Render("  "+field.label+":"))

			// Display current value with proper formatting
			displayValue := value
			if displayValue == "" {
				displayValue = "(empty)"
				displayValue = lipgloss.NewStyle().
					Italic(true).
					Foreground(lipgloss.Color("#45475a")).
					Render(displayValue)
			} else {
				// Truncate long values for display
				if len(displayValue) > 20 {
					displayValue = displayValue[:17] + "..."
				}
				displayValue = lipgloss.NewStyle().
					Foreground(lipgloss.Color("#cdd6f4")).
					Render(displayValue)
			}
			content = append(content, "    "+displayValue)
		}

		// Add spacing between fields (except after last field)
		if i < len(fields)-1 {
			content = append(content, "")
		}
	}

	content = append(content, "") // Extra spacing before instructions

	// Instructions with better formatting
	instrStyle := lipgloss.NewStyle().
		Italic(true).
		Foreground(lipgloss.Color("#6c7086"))
	content = append(content, instrStyle.Render("Tab: Next field"))
	content = append(content, instrStyle.Render("Enter: Save"))
	content = append(content, instrStyle.Render("ESC: Cancel"))

	// Fill remaining space
	for len(content) < 12 { // Match the updated height from renderNetworkSection
		content = append(content, "")
	}

	// Create editor panel with border and better styling
	panelContent := strings.Join(content, "\n")
	return lipgloss.NewStyle().
		Width(30).
		Padding(2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#f9e2af")).
		Render(panelContent)
}

func (m *SettingsModel) renderWalletSection() string {
	var content []string

	// Section title
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#cba6f7")).
		Render("Wallet & Network Status")
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
		// Display mode - show wallet address and network status
		content = append(content, "Wallet address:")
		content = append(content, "")

		addr := m.config.WalletAddr
		if addr == "" {
			addr = "Not configured"
			addrStyle := lipgloss.NewStyle().
				Italic(true).
				Foreground(lipgloss.Color("#f38ba8"))
			addr = addrStyle.Render(addr)
		} else {
			// Truncate long addresses for display (adjusted for increased padding)
			if len(addr) > 35 { // Adjusted for increased padding (less available width)
				addr = addr[:16] + "..." + addr[len(addr)-16:]
			}
			addrStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color("#a6e3a1"))
			addr = addrStyle.Render(addr)
		}

		content = append(content, addr)
		content = append(content, "")

		// Show network connection status
		if m.networkManager.IsConnected() {
			status, err := m.networkManager.GetNetworkStatus()
			if err == nil {
				content = append(content, lipgloss.NewStyle().Bold(true).Render("Network Status:"))
				content = append(content, lipgloss.NewStyle().
					Foreground(lipgloss.Color("#a6e3a1")).
					Render(fmt.Sprintf("Connected to %s", status.NetworkName)))
				content = append(content, fmt.Sprintf("Round: %d", status.LastRound))

				if status.IndexerHealthy {
					content = append(content, lipgloss.NewStyle().
						Foreground(lipgloss.Color("#a6e3a1")).
						Render("Indexer: Online"))
				} else if status.IndexerURL != "" {
					content = append(content, lipgloss.NewStyle().
						Foreground(lipgloss.Color("#f38ba8")).
						Render("Indexer: Offline"))
				}
			}
		} else {
			content = append(content, lipgloss.NewStyle().Bold(true).Render("Network Status:"))
			content = append(content, lipgloss.NewStyle().
				Foreground(lipgloss.Color("#6c7086")).
				Render("Not connected"))
		}

		content = append(content, "")
		editInstr := lipgloss.NewStyle().
			Italic(true).
			Foreground(lipgloss.Color("#cdd6f4")).
			Render("Press 'e' to edit wallet")
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

	if m.editingNetwork {
		instructions = []string{
			"Tab: Next field",
			"Shift+Tab: Prev field",
			"Enter: Save network",
			"ESC: Cancel",
		}
	} else if m.editingAddr {
		instructions = []string{
			"Enter: Save address",
			"ESC: Cancel editing",
		}
	} else {
		baseInstructions := []string{
			"Up/Down: Navigate",
			"Enter: Connect to network",
			"t: Test connection",
			"e: Edit wallet",
			"n: Edit network",
			"c: Create network",
			"ESC: Back",
		}

		if m.showStatus {
			// Add option to hide status when it's shown
			baseInstructions = append(baseInstructions, "Space: Hide status")
		}

		instructions = baseInstructions
	}

	// Add mode indicator for clarity
	var modeIndicator string
	if m.editingNetwork {
		modeIndicator = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#f9e2af")).
			Render("[EDITING NETWORK] ")
	} else if m.editingAddr {
		modeIndicator = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#f9e2af")).
			Render("[EDITING WALLET] ")
	}

	instructionText := strings.Join(instructions, " | ")
	fullText := modeIndicator + instructionText

	return lipgloss.NewStyle().
		Width(77).              // Total width of both panels + gap (30+2+45)
		Align(lipgloss.Center). // Center the footer text
		Foreground(lipgloss.Color("#6c7086")).
		Italic(true).
		Render(fullText)
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
		return Config{
			Network:        "testnet",
			WalletAddr:     "",
			CustomNetworks: []NetworkInfo{},
		}, nil // Default config
	}
	var cfg Config
	if err := json.Unmarshal(buff, &cfg); err != nil {
		return Config{}, err
	}

	// Initialize CustomNetworks if nil
	if cfg.CustomNetworks == nil {
		cfg.CustomNetworks = []NetworkInfo{}
	}

	return cfg, nil
}

func SaveConfig(cfg Config) error {
	// Ensure directory exists
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %w", err)
	}

	configDir := fmt.Sprintf("%s/.lazy-chain", home)
	err = os.MkdirAll(configDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Ensure CustomNetworks is not nil
	if cfg.CustomNetworks == nil {
		cfg.CustomNetworks = []NetworkInfo{}
	}

	// Marshal config to JSON with proper formatting
	buff, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write to file with proper permissions
	configPath := ConfigPath()
	err = ioutil.WriteFile(configPath, buff, 0644)
	if err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}
