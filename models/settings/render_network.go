package settings

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m *SettingsModel) renderNetworkSection() string {
	if m.editingNetwork {
		return m.renderNetworkEditor()
	}

	var content []string
	title := lipgloss.NewStyle().
		Bold(true).Foreground(lipgloss.Color("#cba6f7")).
		Render("Network Selection")
	content = append(content, title, "")

	for i, network := range m.networks {
		cursor := "  "
		style := lipgloss.NewStyle()

		if m.cursor == i && !m.editingAddr {
			cursor = "> "
			style = style.Foreground(lipgloss.Color("#ef9f76"))
		}

		displayName := network
		if network == m.config.Network {
			style = style.Bold(true).Foreground(lipgloss.Color("#a6e3a1"))
			displayName += " (active)"
		}

		isCustom := false
		for _, customNet := range m.config.CustomNetworks {
			if customNet.Name == network {
				isCustom = true
				break
			}
		}
		if isCustom && network != m.config.Network {
			displayName += " *"
		}
		content = append(content, cursor+style.Render(displayName))
	}
	content = append(content, "")

	// Legend
	hasCustom := false
	for _, network := range m.networks {
		for _, customNet := range m.config.CustomNetworks {
			if customNet.Name == network && network != m.config.Network {
				hasCustom = true
				break
			}
		}
		if hasCustom {
			break
		}
	}
	if hasCustom {
		legendStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#6c7086")).Italic(true)
		content = append(content, legendStyle.Render("* = custom/modified"))
	}
	content = append(content, "")

	// Details
	if m.cursor < len(m.networks) {
		currentNet := m.networks[m.cursor]
		m.refreshNetworkInfo(currentNet)

		if info, ok := m.networkInfos[currentNet]; ok {
			header := lipgloss.NewStyle().Bold(true).
				Foreground(lipgloss.Color("#89b4fa")).Render("Network Details:")
			content = append(content, header, "")

			algodStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#cdd6f4"))
			algodURL := info.AlgodURL
			if len(algodURL) > 22 {
				content = append(content, algodStyle.Render("Algod:"))
				content = append(content, algodStyle.Render("  "+algodURL[:20]+"..."))
				content = append(content, algodStyle.Render("  Port: "+info.AlgodPort))
			} else {
				content = append(content, algodStyle.Render(fmt.Sprintf("Algod: %s:%s", algodURL, info.AlgodPort)))
			}

			if info.AlgodToken != "" {
				td := info.AlgodToken
				if len(td) > 16 {
					td = td[:6] + "..." + td[len(td)-6:]
				}
				content = append(content, algodStyle.Render("Token: "+td))
			} else {
				content = append(content, algodStyle.Render("Token: (none)"))
			}
			content = append(content, "")

			// Connection status
			if m.networkManager.IsConnected() && m.networkManager.GetCurrentNetwork().Name == currentNet {
				status := lipgloss.NewStyle().Foreground(lipgloss.Color("#a6e3a1")).Bold(true).
					Render("Status: Connected")
				content = append(content, status)
			} else {
				status := lipgloss.NewStyle().Foreground(lipgloss.Color("#6c7086")).Italic(true).
					Render("Status: Not connected")
				content = append(content, status)
			}
		} else {
			err := lipgloss.NewStyle().Foreground(lipgloss.Color("#f38ba8")).Italic(true).
				Render("Network details unavailable")
			content = append(content, err)
		}
	}
	content = append(content, "")

	// Status messages (wrapped)
	if m.showStatus && m.connectionStatus != "" {
		statusStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#f9e2af")).Italic(true)
		if len(m.connectionStatus) > 22 {
			words := strings.Fields(m.connectionStatus)
			var lines []string
			var current []string
			cl := 0
			for _, w := range words {
				if cl+len(w)+len(current) > 22 {
					if len(current) > 0 {
						lines = append(lines, strings.Join(current, " "))
						current = []string{w}
						cl = len(w)
					} else {
						lines = append(lines, w)
					}
				} else {
					current = append(current, w)
					cl += len(w)
				}
			}
			if len(current) > 0 {
				lines = append(lines, strings.Join(current, " "))
			}
			for _, ln := range lines {
				content = append(content, statusStyle.Render(ln))
			}
		} else {
			content = append(content, statusStyle.Render(m.connectionStatus))
		}
	}

	for len(content) < 12 {
		content = append(content, "")
	}

	panel := strings.Join(content, "\n")
	return lipgloss.NewStyle().
		Width(30).
		Padding(2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#89b4fa")).
		Render(panel)
}

func (m *SettingsModel) renderNetworkEditor() string {
	var content []string

	title := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#f9e2af")).Render("Edit Network")
	content = append(content, title, "")

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

	for i, f := range fields {
		value := m.getNetworkField(f.key)
		if f.key == m.editField {
			fieldStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#f9e2af")).Bold(true)
			content = append(content, fieldStyle.Render("● "+f.label+":"))
			input := lipgloss.NewStyle().
				Background(lipgloss.Color("#1e1e2e")).
				Foreground(lipgloss.Color("#f9e2af")).
				Padding(0, 1).
				Render(value + "█")
			content = append(content, "  "+input)
		} else {
			fieldStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#6c7086"))
			content = append(content, fieldStyle.Render("  "+f.label+":"))

			display := value
			if display == "" {
				display = lipgloss.NewStyle().Italic(true).Foreground(lipgloss.Color("#45475a")).Render("(empty)")
			} else {
				if len(display) > 20 {
					display = display[:17] + "..."
				}
				display = lipgloss.NewStyle().Foreground(lipgloss.Color("#cdd6f4")).Render(display)
			}
			content = append(content, "    "+display)
		}
		if i < len(fields)-1 {
			content = append(content, "")
		}
	}

	content = append(content, "")
	instr := lipgloss.NewStyle().Italic(true).Foreground(lipgloss.Color("#6c7086"))
	content = append(content, instr.Render("Tab: Next field"))
	content = append(content, instr.Render("Enter: Save"))
	content = append(content, instr.Render("ESC: Cancel"))

	for len(content) < 12 {
		content = append(content, "")
	}

	panel := strings.Join(content, "\n")
	return lipgloss.NewStyle().
		Width(30).
		Padding(2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#f9e2af")).
		Render(panel)
}
