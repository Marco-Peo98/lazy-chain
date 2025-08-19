package components

import "github.com/charmbracelet/lipgloss"

type Output struct {
	Title string
	Text  string
	Width int
}

func (o Output) Render() string {
	title := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#a6e3a1")).Render(o.Title)
	body := lipgloss.NewStyle().Width(o.Width).Render(o.Text)
	return lipgloss.NewStyle().
		Width(o.Width+4).
		Padding(1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#a6e3a1")).
		Render(title + "\n\n" + body)
}
