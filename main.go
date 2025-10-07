package main

import (
	_ "embed"
	"fmt"
	"lazychain/layout" // New simple layout package
	. "lazychain/models"
	. "lazychain/models/goal"
	. "lazychain/models/settings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var p *tea.Program

var Focus string

//go:embed misc/banner.txt
var banner string

type MainModel struct {
	// Replace individual width/height with layout container
	layout            *layout.LayoutContainer
	CurrentState      SessionState
	ProjectModel      *ProjectModel
	SettingsModel     *SettingsModel
	ApplicationsModel *ApplicationsModel
	CmdGoalsModel     *GOALModel
	ExploreModel      *ExploreModel
}

func NewMainModel() *MainModel {
	// Initialize with default dimensions, will be updated by first WindowSizeMsg
	initialLayout := layout.NewLayoutContainer(80, 24)

	return &MainModel{
		layout:            initialLayout,
		CurrentState:      MainView,
		ProjectModel:      NewProjectModel(),
		SettingsModel:     NewSettingsModel([]string{"localnet", "testnet", "mainnet"}),
		ApplicationsModel: NewApplicationsModel(),
		CmdGoalsModel:     NewGOALModel(),
		ExploreModel:      NewExploreModel(),
	}
}

func (m *MainModel) Init() tea.Cmd {
	return tea.SetWindowTitle("LAZYCHAIN - your friendly TUI ≧ ﹏ ≦")
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
				// Back to main view - completely reset selection state
				m.CurrentState = MainView
				m.ProjectModel.Selected = make(map[int]struct{})
				return m, nil
			}

			updateModel, cmd := m.ProjectModel.Update(msg)
			if pm, ok := updateModel.(*ProjectModel); ok {
				m.ProjectModel = pm

				// Update the current state based on the selected option
				// Only change state if there's exactly one selection
				if len(m.ProjectModel.Selected) == 1 {
					for i := range m.ProjectModel.Selected {
						switch m.ProjectModel.Options[i] {
						case "Settings":
							m.CurrentState = SettingsView
							// Reset editing state when entering settings
							m.SettingsModel.ResetEditingState()
						case "Applications":
							m.CurrentState = ApplicationsView
						case "Commands Goals":
							m.CurrentState = CmdGoalsView
						case "Explore":
							m.CurrentState = ExploreView
						}
						// Clear selection after state change to prevent re-triggering
						m.ProjectModel.Selected = make(map[int]struct{})
						break
					}
				}
			}
			return m, cmd
		case SettingsView:
			// Check if we're in editing mode BEFORE processing the message
			wasEditingAddr := m.SettingsModel.IsEditingAddr()

			// Let the SettingsModel handle the message
			var cmd tea.Cmd
			updatedModel, cmd := m.SettingsModel.Update(msg)
			if updatedSettingsModel, ok := updatedModel.(*SettingsModel); ok {
				m.SettingsModel = updatedSettingsModel
			}

			// Only handle ESC in MainModel if we were NOT in editing mode
			if msg.String() == "esc" && !wasEditingAddr {
				// Reset editing state when leaving settings
				m.SettingsModel.ResetEditingState()
				m.CurrentState = ProjectView
				return m, nil
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
			if updatedCmdGoalsModel, ok := updatedModel.(*GOALModel); ok {
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
		// Update the layout container with new dimensions
		m.layout.Resize(msg.Width, msg.Height)
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

		// Use layout container to render main view
		return m.layout.Render(content)

	case ProjectView:
		// Use layout container to render project view
		return m.layout.Render(m.ProjectModel.View())

	case SettingsView:
		// Use layout container to render settings view
		return m.layout.Render(m.SettingsModel.View())

	case ApplicationsView:
		// Use layout container to render applications view
		return m.layout.Render(m.ApplicationsModel.View())

	case CmdGoalsView:
		// Use layout container to render cmd goals view
		return m.layout.Render(m.CmdGoalsModel.View())

	case ExploreView:
		// Use layout container to render explore view
		return m.layout.Render(m.ExploreModel.View())

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
