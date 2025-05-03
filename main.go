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

//go:embed misc/banner.txt
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
		switch m.CurrentState {
		case MainView:
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
		case ProjectView:
			switch msg.String() {
			case "esc":
				// Back to list if you're inside a subview
				m.CurrentState = ProjectView
				m.ProjectModel.Selected = make(map[int]struct{})
				return m, nil
			}

			updateModel, cmd := m.ProjectModel.Update(msg)
			if pm, ok := updateModel.(*ProjectModel); ok {
				m.ProjectModel = pm

				// Update the current state based on the selected option
				if len(m.ProjectModel.Selected) > 0 {
					for i := range m.ProjectModel.Selected {
						switch m.ProjectModel.Options[i] {
						case "Settings":
							m.CurrentState = SettingsView
						case "Applications":
							m.CurrentState = ApplicationsView
						case "Commands Goals":
							m.CurrentState = CmdGoalsView
						case "Explore":
							m.CurrentState = ExploreView
						}
					}
				}
			}
			return m, cmd
		case SettingsView:
			switch msg.String() {
			case "esc":
				m.CurrentState = ProjectView
				return m, nil
			}
			var cmd tea.Cmd
			updatedModel, cmd := m.SettingsModel.Update(msg)
			if updatedSettingsModel, ok := updatedModel.(*SettingsModel); ok {
				m.SettingsModel = updatedSettingsModel
			}
			return m, cmd
		case ApplicationsView:
			switch msg.String() {
			case "esc":
				m.CurrentState = ProjectView
				return m, nil
			}
			var cmd tea.Cmd
			updatedModel, cmd := m.ApplicationsModel.Update(msg)
			if updatedApplicationsModel, ok := updatedModel.(*ApplicationsModel); ok {
				m.ApplicationsModel = updatedApplicationsModel
			}
			return m, cmd
		case CmdGoalsView:
			switch msg.String() {
			case "esc":
				m.CurrentState = ProjectView
				return m, nil
			}
			var cmd tea.Cmd
			updatedModel, cmd := m.CmdGoalsModel.Update(msg)
			if updatedCmdGoalsModel, ok := updatedModel.(*CmdGoalsModel); ok {
				m.CmdGoalsModel = updatedCmdGoalsModel
			}
			return m, cmd
		case ExploreView:
			switch msg.String() {
			case "esc":
				m.CurrentState = ProjectView
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
	switch m.CurrentState {
	case MainView:
		bannerStyled := lipgloss.NewStyle().Foreground(lipgloss.Color("#81c8be")).Render(banner)
		instr := lipgloss.NewStyle().Align(lipgloss.Center).Foreground(lipgloss.Color("#c6d0f5")).Render("\nPress 'Enter' to start\nPress 'Ctrl+C' or 'q' to quit")
		content := lipgloss.JoinVertical(lipgloss.Center, bannerStyled, instr)
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center,
			lipgloss.NewStyle().Padding(1).Render(content),
		)
	case ProjectView:
		return lipgloss.Place(m.width, m.height,
			lipgloss.Center, lipgloss.Center,
			m.ProjectModel.View(),
		)
	case SettingsView:
		return m.SettingsModel.View()
	case ApplicationsView:
		return m.ApplicationsModel.View()
	case CmdGoalsView:
		return m.CmdGoalsModel.View()
	case ExploreView:
		return m.ExploreModel.View()
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
