package components

import "github.com/charmbracelet/lipgloss"

type ListNav struct {
	Items   []string
	Cursor  int
	Active  bool
	Width   int
	Title   string
}

func (l *ListNav) Up()   { if l.Cursor > 0 { l.Cursor-- } }
func (l *ListNav) Down() { if l.Cursor < len(l.Items)-1 { l.Cursor++ } }

func (l ListNav) Render() string {
	title := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#cba6f7")).Render(l.Title)
	lines := []string{title, ""}
	for i, it := range l.Items {
		cur := "  "
		style := lipgloss.NewStyle()
		if i == l.Cursor && l.Active {
			cur = "> "
			style = style.Foreground(lipgloss.Color("#ef9f76"))
		}
		lines = append(lines, cur+style.Render(it))
	}
	block := lipgloss.NewStyle().
		Width(l.Width).
		Padding(1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#89b4fa")).
		Render(join(lines))
	return block
}

func join(ss []string) string {
	s := ""
	for i, l := range ss {
		if i > 0 { s += "\n" }
		s += l
	}
	return s
}
