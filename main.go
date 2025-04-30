package main

import (
	_ "embed"
	"fmt"
	. "lazychain/models"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var p *tea.Program

var Focus string

//go:embed banner.txt
var banner string

type MainModel struct {
	width, height     int
	CurrentState      SessionState
	ProjectModel      *ProjectModel
	SettingsModel     *SettingsModel
	ApplicationsModel *ApplicationsModel
	CmdGoalsModel     *CmdGoalsModel
	ExploreModel      *ExploreModel
}

func NewMainModel() *MainModel {
	return &MainModel{
		width:             0,
		height:            0,
		CurrentState:      MainView,
		ProjectModel:      NewProjectModel(),
		SettingsModel:     NewSettingsModel(),
		ApplicationsModel: NewApplicationsModel(),
		CmdGoalsModel:     NewCmdGoalsModel(),
		ExploreModel:      NewExploreModel(),
	}
}

func (m *MainModel) Init() tea.Cmd {
	return nil
}

func (m *MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.CurrentState == MainView {

			// When in the main view, every selection will be reset
			m.ProjectModel.Selected = make(map[int]struct{})
			m.ProjectModel.Cursor = 0
			Focus = ""

			switch msg.String() {
			case "enter":
				m.CurrentState = ProjectView
				return m, nil
			case "ctrl+c", "q":
				return m, tea.Quit
			}
		} else if m.CurrentState == ProjectView {
			switch msg.String() {
			case "esc":
				m.CurrentState = MainView
				return m, nil
			}
			var cmd tea.Cmd
			updatedModel, cmd := m.ProjectModel.Update(msg)
			if updatedProjectModel, ok := updatedModel.(*ProjectModel); ok {
				m.ProjectModel = updatedProjectModel
			}
			if cmd != nil {
				if focusMsg, ok := cmd().(string); ok {
					Focus = focusMsg
				}
			}
			return m, cmd
		} else if m.CurrentState == SettingsView {
			switch msg.String() {
			case "esc":
				m.CurrentState = MainView
				return m, nil
			}
			var cmd tea.Cmd
			updatedModel, cmd := m.SettingsModel.Update(msg)
			if updatedSettingsModel, ok := updatedModel.(*SettingsModel); ok {
				m.SettingsModel = updatedSettingsModel
			}
			return m, cmd
		} else if m.CurrentState == ApplicationsView {
			switch msg.String() {
			case "esc":
				m.CurrentState = MainView
				return m, nil
			}
			var cmd tea.Cmd
			updatedModel, cmd := m.ApplicationsModel.Update(msg)
			if updatedApplicationsModel, ok := updatedModel.(*ApplicationsModel); ok {
				m.ApplicationsModel = updatedApplicationsModel
			}
			return m, cmd
		} else if m.CurrentState == CmdGoalsView {
			switch msg.String() {
			case "esc":
				m.CurrentState = MainView
				return m, nil
			}
			var cmd tea.Cmd
			updatedModel, cmd := m.CmdGoalsModel.Update(msg)
			if updatedCmdGoalsModel, ok := updatedModel.(*CmdGoalsModel); ok {
				m.CmdGoalsModel = updatedCmdGoalsModel
			}
			return m, cmd
		} else if m.CurrentState == ExploreView {
			switch msg.String() {
			case "esc":
				m.CurrentState = MainView
				return m, nil
			}
			var cmd tea.Cmd
			updatedModel, cmd := m.ExploreModel.Update(msg)
			if updatedExploreModel, ok := updatedModel.(*ExploreModel); ok {
				m.ExploreModel = updatedExploreModel
			}
			return m, cmd
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	}

	return m, nil
}

func (m *MainModel) View() string {

	banner = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#81c8be")).
		Render(banner)

	mainTxt := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Foreground(lipgloss.Color("#c6d0f5")).
		Render("\n\nPress 'Enter' to start. \nPress 'q' or 'Ctrl+c' to quit.\n\n")

	content := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Render((banner + "\n\n" + mainTxt))

	box := lipgloss.NewStyle().
		Padding(1).
		Margin(1).
		Render(content)

	container := lipgloss.Place(m.width, m.height,
		lipgloss.Center,
		lipgloss.Center,
		box,
	)

	switch m.CurrentState {
	case MainView:
		return container
	case ProjectView:
		switch Focus {
		case "SettingsView":
			return fmt.Sprintf(
				"%s\n\n%s",
				m.ProjectModel.View(),
				m.SettingsModel.View(),
			)
		case "ApplicationsView":
			return fmt.Sprintf(
				"%s\n\n%s",
				m.ProjectModel.View(),
				m.ApplicationsModel.View(),
			)
		case "Commands GoalsView":
			return fmt.Sprintf(
				"%s\n\n%s",
				m.ProjectModel.View(),
				m.CmdGoalsModel.View(),
			)
		case "ExploreView":
			return fmt.Sprintf(
				"%s\n\n%s",
				m.ProjectModel.View(),
				m.ExploreModel.View(),
			)
		default:
			return m.ProjectModel.View()
		}

	default:
		return ""
	}
}

func main() {
	p = tea.NewProgram(NewMainModel(), tea.WithAltScreen())
	if err := p.Start(); err != nil {
		fmt.Printf("Error starting program: %v\n", err)
	}
}
