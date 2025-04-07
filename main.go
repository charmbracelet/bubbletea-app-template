// Package main contains a Bubble Tea application.
package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/v2/help"
	"github.com/charmbracelet/bubbles/v2/key"
	"github.com/charmbracelet/bubbles/v2/spinner"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
)

// errMsg is a general error message.
//
// Note: we wrap errors in a struct in order to differentiate between error
// messages as the following are identical when type switching:
//
//	type thisErr error
//	type thatErr error
type errMsg struct {
	err error
}

// model contains the application state.
type model struct {
	spinner  spinner.Model
	keymap   keymap
	help     help.Model
	quitting bool
	err      error
}

// initialModel returns the initial state for the program.
func initialModel() model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return model{
		spinner: s,
		help:    help.New(),
		keymap:  defaultKeymap(),
	}
}

// Init send the initial command to the Program.
func (m model) Init() tea.Cmd {
	return tea.Batch(
		tea.RequestBackgroundColor,
		m.spinner.Tick,
	)
}

// Update handles the messages sent to the program.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.BackgroundColorMsg:
		// Now that we know the background color we can initialize the
		// appropriate styles.
		m.help.Styles = help.DefaultStyles(msg.IsDark())
		return m, nil

	case tea.KeyPressMsg:
		// Handle keypresses.
		if key.Matches(msg, m.keymap.quit) {
			m.quitting = true
			return m, tea.Quit
		}
		if key.Matches(msg, m.keymap.interrupt) {
			return m, tea.Interrupt
		}
		if key.Matches(msg, m.keymap.suspend) {
			return m, tea.Suspend
		}
		return m, nil

	case errMsg:
		// Handle errors.
		m.err = msg.err
		return m, nil

	default:
		// Keep the spinner spinning.
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
}

// View renders the Program's view.
func (m model) View() string {
	if m.err != nil {
		return m.err.Error()
	}
	return fmt.Sprintf("\n  %s Loading forever...\n\n%s\n", m.spinner.View(), m.help.View(m.keymap))
}

// keymap contains the key bindings for the program.
type keymap struct {
	quit      key.Binding
	interrupt key.Binding
	suspend   key.Binding
}

// ShortHelp returns the short help for the program. This is used to render
// help.
func (k keymap) ShortHelp() []key.Binding {
	return []key.Binding{k.quit, k.interrupt, k.suspend}
}

// FullHelp returns the long help for the program.
//
// This is currently unused in this application in its current state, but we're
// leving it here incase your application grows to the point where you need
// more keybindings and help.
func (k keymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.quit, k.interrupt, k.suspend},
		// If the application grows to the point where we need more
		// keybindings, you can add additional columns here.
	}
}

func defaultKeymap() keymap {
	return keymap{
		quit: key.NewBinding(
			key.WithKeys("q", "esc"),
			key.WithHelp("q", "quit"),
		),
		interrupt: key.NewBinding(
			key.WithKeys("ctrl+c"),
			key.WithHelp("ctrl+c", "interrupt"),
		),
		suspend: key.NewBinding(
			key.WithKeys("ctrl+z"),
			key.WithHelp("ctrl+z", "suspend"),
		),
	}
}

func main() {
	p := tea.NewProgram(
		initialModel(),

		// Uncomment the following to run the program in the alternate screen
		// buffer (i.e. the full terminal window).
		// tea.WithAltScreen(),
	)
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error while running program: %v\n", err)
		os.Exit(1)
	}
}
