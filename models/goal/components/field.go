package components

import "github.com/charmbracelet/lipgloss"

type Field struct {
	Label   string
	Value   string
	Cursor  int
	Active  bool
	Hint    string
	Secret  bool
	MaxLen  int
}

func (f *Field) SetActive(a bool) { f.Active = a }

func (f *Field) InsertRune(r rune) {
	if f.MaxLen > 0 && len(f.Value) >= f.MaxLen { return }
	left := f.Value[:f.Cursor]
	right := f.Value[f.Cursor:]
	f.Value = left + string(r) + right
	f.Cursor++
}

func (f *Field) Backspace() {
	if f.Cursor == 0 || len(f.Value) == 0 { return }
	left := f.Value[:f.Cursor-1]
	right := f.Value[f.Cursor:]
	f.Value = left + right
	f.Cursor--
}

func (f *Field) MoveLeft()  { if f.Cursor > 0 { f.Cursor-- } }
func (f *Field) MoveRight() { if f.Cursor < len(f.Value) { f.Cursor++ } }

func (f Field) Render(width int) string {
	label := lipgloss.NewStyle().Bold(true).Render(f.Label + ":")
	val := f.Value
	if f.Secret && len(val) > 0 { val = "••••••••" }
	if f.Active { val += "█" }
	value := lipgloss.NewStyle().Width(width).Render(val)
	hint := ""
	if f.Hint != "" {
		hint = lipgloss.NewStyle().Faint(true).Render("  " + f.Hint)
	}
	return label + "\n" + "  " + value + hint
}
