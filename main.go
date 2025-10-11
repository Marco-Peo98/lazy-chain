package main

import (
	_ "embed"
	"fmt"
	"lazychain/layout" // Layout package
	. "lazychain/models"
	. "lazychain/models/goal"
	. "lazychain/models/settings"

	tea "github.com/charmbracelet/bubbletea"
)

var p *tea.Program
var Focus string

//go:embed misc/banner.txt
var banner string

type MainModel struct {
	// Layout container (kept for non-converted views)
	layoutContainer *layout.LayoutContainer

	// FlexBox layouts
	mainLayout    *layout.MainLayout
	projectLayout *layout.ProjectLayout

	// Current dimensions
	width, height int

	CurrentState      SessionState
	ProjectModel      *ProjectModel
	SettingsModel     *SettingsModel
	ApplicationsModel *ApplicationsModel
	CmdGoalsModel     *GOALModel
	ExploreModel      *ExploreModel
}

func NewMainModel() *MainModel {
	// Initialize with default dimensions
	initialLayout := layout.NewLayoutContainer(80, 24)

	return &MainModel{
		layoutContainer:   initialLayout,
		mainLayout:        nil, // Will be initialized on first WindowSizeMsg
		projectLayout:     nil, // Will be initialized on first WindowSizeMsg
		width:             80,
		height:            24,
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
			// When in the main view, reset all selection state
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
		// Update dimensions
		m.width = msg.Width
		m.height = msg.Height

		// Update layout container (for non-converted views)
		m.layoutContainer.Resize(msg.Width, msg.Height)

		// Create/update MainLayout with new dimensions
		instructions := "\nPress 'Enter' to start\nPress 'Ctrl+C' or 'q' to quit"
		m.mainLayout = layout.NewMainLayout(msg.Width, msg.Height, banner, instructions)

		// Create/update ProjectLayout with new dimensions
		m.projectLayout = layout.NewProjectLayout(msg.Width, msg.Height)

		return m, nil
	}

	return m, nil
}

func (m *MainModel) View() string {
	switch m.CurrentState {
	case MainView:
		// Use FlexBox for MainView
		if m.mainLayout == nil {
			// Fallback if not yet initialized
			instructions := "\nPress 'Enter' to start\nPress 'Ctrl+C' or 'q' to quit"
			m.mainLayout = layout.NewMainLayout(m.width, m.height, banner, instructions)
		}

		// IMPORTANT: Check if dimensions are valid
		if !m.mainLayout.IsValid() {
			// Show error message if terminal too small
			return m.mainLayout.RenderError()
		}

		// Build and render FlexBox only if dimensions OK
		flexBox := m.mainLayout.Build()
		return flexBox.Render()

	case ProjectView:
		// Use FlexBox for ProjectView
		if m.projectLayout == nil {
			m.projectLayout = layout.NewProjectLayout(m.width, m.height)
		}

		// Get current option for preview
		currentOption := ""
		if m.ProjectModel.Cursor < len(m.ProjectModel.Options) {
			currentOption = m.ProjectModel.Options[m.ProjectModel.Cursor]
		}

		// Get subtitle/description for current option
		description := subtitleFor(currentOption)
		instructions := "Press ENTER to select\nPress ESC to go back"

		// Configure layout with current state
		m.projectLayout.
			SetMenuItems(m.ProjectModel.Options).
			SetCursor(m.ProjectModel.Cursor).
			SetPreviewContent(currentOption, description, instructions)

		// Check dimensions
		if !m.projectLayout.IsValid() {
			// Reuse error from MainLayout or create simple message
			return "Terminal too small for ProjectView\nMinimum: 80x24"
		}

		// Build and render FlexBox
		flexBox := m.projectLayout.Build()
		return flexBox.Render()

	case SettingsView:
		return m.layoutContainer.Render(m.SettingsModel.View())

	case ApplicationsView:
		return m.layoutContainer.Render(m.ApplicationsModel.View())

	case CmdGoalsView:
		return m.layoutContainer.Render(m.CmdGoalsModel.View())

	case ExploreView:
		return m.layoutContainer.Render(m.ExploreModel.View())

	default:
		return ""
	}
}

// subtitleFor returns the subtitle for the given option
func subtitleFor(option string) string {
	switch option {
	case "Settings":
		return "Configure your network and wallet settings"
	case "Applications":
		return "Manage your blockchain applications"
	case "Commands Goals":
		return "Why CLI when you can TUI? Build transactions easily"
	case "Explore":
		return "Explore blockchain data and resources"
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
