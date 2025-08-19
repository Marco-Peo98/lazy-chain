package settings

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

// Update handles key input, editing flows, and connection actions.
func (m *SettingsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// 1) Network editing has highest priority
		if m.editingNetwork {
			return m.handleNetworkEditing(msg)
		}

		// 2) Wallet address editing
		if m.editingAddr {
			switch msg.String() {
			case "enter":
				m.config.WalletAddr = m.inputBuffer
				m.editingAddr = false
				SaveConfig(m.config)
				return m, nil
			case "esc":
				m.editingAddr = false
				m.inputBuffer = ""
				return m, nil
			}

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

		// 3) Normal mode
		switch msg.String() {
		case "up":
			if m.cursor > 0 {
				m.cursor--
				m.showStatus = false
				m.connectionStatus = ""
			}
		case "down":
			if m.cursor < len(m.networks)-1 {
				m.cursor++
				m.showStatus = false
				m.connectionStatus = ""
			}
		case "enter":
			// Connect to selected network
			selectedNetwork := m.networks[m.cursor]
			if networkInfo, exists := m.networkInfos[selectedNetwork]; exists {
				if err := m.networkManager.TestNetworkConnection(networkInfo); err != nil {
					m.connectionStatus = fmt.Sprintf("Failed to connect to %s: %v", selectedNetwork, err)
					m.showStatus = true
				} else {
					m.config.Network = selectedNetwork
					SaveConfig(m.config)
					if err := m.networkManager.ConnectToNetwork(networkInfo); err != nil {
						m.connectionStatus = fmt.Sprintf("Connection failed: %v", err)
					} else {
						m.connectionStatus = fmt.Sprintf("Successfully connected to %s", selectedNetwork)
					}
					m.showStatus = true
				}
			}
		case "t":
			// Test selected network
			selectedNetwork := m.networks[m.cursor]
			if networkInfo, exists := m.networkInfos[selectedNetwork]; exists {
				if err := m.networkManager.TestNetworkConnection(networkInfo); err != nil {
					m.connectionStatus = fmt.Sprintf("Test failed for %s: %v", selectedNetwork, err)
				} else {
					m.connectionStatus = fmt.Sprintf("Test successful for %s", selectedNetwork)
				}
				m.showStatus = true
			}
		case " ":
			if m.showStatus {
				m.showStatus = false
				m.connectionStatus = ""
			}
		case "e":
			// Start editing wallet
			m.editingAddr = true
			m.inputBuffer = m.config.WalletAddr
		case "n":
			// Edit selected network
			currentNetworkName := m.networks[m.cursor]
			if info, exists := m.networkInfos[currentNetworkName]; exists {
				m.editingNetwork = true
				m.networkBuffer = info
				m.editField = "algod_url"
				m.showStatus = false
				m.connectionStatus = ""
			}
		case "c":
			// Create new custom network
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
			m.showStatus = false
			m.connectionStatus = ""
		}
	}
	return m, nil
}

func (m *SettingsModel) handleNetworkEditing(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "tab":
		m.editField = m.getNextField(m.editField)
	case "shift+tab":
		m.editField = m.getPrevField(m.editField)
	case "enter":
		m.saveNetworkConfig()
		m.editingNetwork = false
		return m, nil
	case "esc":
		m.editingNetwork = false
		m.networkBuffer = NetworkInfo{}
		return m, nil
	}

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
	for i, f := range fields {
		if f == current {
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
	for i, f := range fields {
		if f == current {
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
	// Validate
	if err := ValidateNetworkConfig(m.networkBuffer); err != nil {
		m.connectionStatus = fmt.Sprintf("Validation failed: %v", err)
		m.showStatus = true
		return
	}

	// Test before saving
	if err := m.networkManager.TestNetworkConnection(m.networkBuffer); err != nil {
		m.connectionStatus = fmt.Sprintf("Connection test failed: %v", err)
		m.showStatus = true
		return
	}

	// Update map immediately for UI
	m.networkInfos[m.networkBuffer.Name] = m.networkBuffer

	// Detect new vs existing
	isNew := true
	existingIndex := -1
	for i, n := range m.networks {
		if n == m.networkBuffer.Name {
			isNew = false
			existingIndex = i
			break
		}
	}

	// Upsert into CustomNetworks
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

	// Update networks list + cursor
	if isNew {
		m.networks = append(m.networks, m.networkBuffer.Name)
		m.cursor = len(m.networks) - 1
	} else {
		m.cursor = existingIndex
	}

	// Save and reload
	if err := SaveConfig(m.config); err != nil {
		m.connectionStatus = fmt.Sprintf("Failed to save config: %v", err)
		m.showStatus = true
		return
	}

	reloaded, err := LoadConfig()
	if err != nil {
		m.connectionStatus = fmt.Sprintf("Warning: Failed to reload config: %v", err)
		m.showStatus = true
	} else {
		m.config = reloaded
		for _, customNet := range m.config.CustomNetworks {
			m.networkInfos[customNet.Name] = customNet
		}
	}

	m.connectionStatus = fmt.Sprintf("Network %s saved successfully", m.networkBuffer.Name)
	m.showStatus = true
}
