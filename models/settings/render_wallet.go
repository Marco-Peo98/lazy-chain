package settings

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m *SettingsModel) renderWalletSection() string {
	var content []string

	title := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#cba6f7")).Render("Wallet & Network Status")
	content = append(content, title, "")

	if m.editingAddr {
		content = append(content, "Editing address:", "")

		inputDisplay := m.inputBuffer + "_"
		inputStyle := lipgloss.NewStyle().
			Width(37).
			Background(lipgloss.Color("#1e1e2e")).
			Foreground(lipgloss.Color("#f9e2af")).
			Padding(0, 1)
		content = append(content, inputStyle.Render(inputDisplay), "")

		instructions := lipgloss.NewStyle().Italic(true).Foreground(lipgloss.Color("#f9e2af")).
			Render("Enter to save, ESC to cancel")
		content = append(content, instructions)
	} else {
		content = append(content, "Wallet address:", "")

		addr := m.config.WalletAddr
		if addr == "" {
			addr = lipgloss.NewStyle().Italic(true).Foreground(lipgloss.Color("#f38ba8")).Render("Not configured")
		} else {
			if len(addr) > 35 {
				addr = addr[:16] + "..." + addr[len(addr)-16:]
			}
			addr = lipgloss.NewStyle().Foreground(lipgloss.Color("#a6e3a1")).Render(addr)
		}
		content = append(content, addr, "")

		if m.networkManager.IsConnected() {
			if status, err := m.networkManager.GetNetworkStatus(); err == nil {
				content = append(content, lipgloss.NewStyle().Bold(true).Render("Network Status:"))
				content = append(content, lipgloss.NewStyle().Foreground(lipgloss.Color("#a6e3a1")).
					Render(fmt.Sprintf("Connected to %s", status.NetworkName)))
				content = append(content, fmt.Sprintf("Round: %d", status.LastRound))

				if status.IndexerHealthy {
					content = append(content, lipgloss.NewStyle().Foreground(lipgloss.Color("#a6e3a1")).Render("Indexer: Online"))
				} else if status.IndexerURL != "" {
					content = append(content, lipgloss.NewStyle().Foreground(lipgloss.Color("#f38ba8")).Render("Indexer: Offline"))
				}
			}
		} else {
			content = append(content, lipgloss.NewStyle().Bold(true).Render("Network Status:"))
			content = append(content, lipgloss.NewStyle().Foreground(lipgloss.Color("#6c7086")).Render("Not connected"))
		}

		content = append(content, "")
		editInstr := lipgloss.NewStyle().Italic(true).Foreground(lipgloss.Color("#cdd6f4")).Render("Press 'e' to edit wallet")
		content = append(content, editInstr)
	}

	for len(content) < 10 {
		content = append(content, "")
	}

	panel := strings.Join(content, "\n")
	return lipgloss.NewStyle().
		Width(45).
		Padding(2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#a6e3a1")).
		Render(panel)
}

func (m *SettingsModel) renderFooter() string {
	var instructions []string

	if m.editingNetwork {
		instructions = []string{"Tab: Next field", "Shift+Tab: Prev field", "Enter: Save network", "ESC: Cancel"}
	} else if m.editingAddr {
		instructions = []string{"Enter: Save address", "ESC: Cancel editing"}
	} else {
		instructions = []string{"Up/Down: Navigate", "Enter: Connect to network", "t: Test connection", "e: Edit wallet", "n: Edit network", "c: Create network", "ESC: Back"}
		if m.showStatus {
			instructions = append(instructions, "Space: Hide status")
		}
	}

	var modeIndicator string
	if m.editingNetwork {
		modeIndicator = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#f9e2af")).Render("[EDITING NETWORK] ")
	} else if m.editingAddr {
		modeIndicator = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#f9e2af")).Render("[EDITING WALLET] ")
	}

	instructionText := strings.Join(instructions, " | ")
	fullText := modeIndicator + instructionText

	return lipgloss.NewStyle().
		Width(77).
		Align(lipgloss.Center).
		Foreground(lipgloss.Color("#6c7086")).
		Italic(true).
		Render(fullText)
}
