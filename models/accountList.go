package models

import (
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// Account type represents a single Algorand account entry
type Account struct {
	Address string  `json:"address"`
	ALGO    float64 `json:"amount,string"`
	// Additional fields can be added here as needed
	// such as account name, creation date, assets, etc.
}

// Messages
// AccountFetchedMsg carries the fetched account slice or an error
type AccountFetchedMsg struct {
	Accounts []Account
	Err      error
}

// AccountListModel holds the state for the account list view
type AccountListModel struct {
	list     list.Model
	accounts []Account
	err      error

	// Search and filter inputs
	searchInput  textinput.Model
	filterActive bool

	// Additional featurs can be added here as needed
	// - Sort by balance, address, etc.
	// - Export to CSV or JSON
	// - Txn history view
	// - Asset opt-in status
}

func NewAccountListModel() AccountListModel {
	items := []list.Item{} // Initialize with empty items, will be populated after fetching accounts
	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Accounts"

	// Initialize search input
	ti := textinput.New()
	ti.Placeholder = "Search by address..."
	ti.Focus()

	return AccountListModel{
		list:        l,
		searchInput: ti,
	}
}

func (m *AccountListModel) Init() tea.Cmd {
	return fetchAccountsCmd()
}

// fetchAccountsCmd runs `goal account list -o json` and parses the output

func fetchAccountsCmd() tea.Cmd {
	return func() tea.Msg {
		out, err := exec.Command("goal", "account", "list", "-o", "json").Output()
		if err != nil {
			return AccountFetchedMsg{Err: err}
		}
		var accounts []Account
		if err := json.Unmarshal(out, &accounts); err != nil {
			return AccountFetchedMsg{Err: err}
		}
		return AccountFetchedMsg{Accounts: accounts}
	}
}

func (m *AccountListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case AccountFetchedMsg:
		if msg.Err != nil {
			m.err = msg.Err
		} else {
			m.accounts = msg.Accounts
			// Populate the list with account items
			var items []list.Item
			for _, acct := range m.accounts {
				items = append(items, acct)
			}
			m.list.SetItems(items)
		}
		return m, nil

		// TODO: Handle other messages like search input, list item selection, etc.
	}

	// Default: delegate to the list
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m *AccountListModel) View() string {
	if m.err != nil {
		return "Error fetching accounts: " + m.err.Error()
	}
	if m.filterActive {
		return m.searchInput.View() + "\n" + m.list.View()
	}
	return m.list.View()
}

// ListItem interface implementation for Account
func (a Account) FilterValue() string {
	return a.Address
	// or any other field you want to filter by
}

func (a Account) Title() string {
	return a.Address
}

func (a Account) Description() string {
	return fmt.Sprintf("Balance: %.6f ALGO", a.ALGO)
}

// TODO: implement sorting, export, txn history, asset indicators, etc.
