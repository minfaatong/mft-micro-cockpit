package app

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/minfaatong/mft-micro-cockpit/internal/collector"
	"github.com/minfaatong/mft-micro-cockpit/internal/domain"
	"github.com/minfaatong/mft-micro-cockpit/internal/ui"
)

// Run starts the interactive TUI program.
func Run() error {
	p := tea.NewProgram(newModel(), tea.WithAltScreen())
	_, err := p.Run()
	return err
}

type tickMsg time.Time

type snapshotMsg struct {
	snapshot domain.DashboardSnapshot
	err      error
}

type model struct {
	collector *collector.Collector
	width     int
	height    int
	snapshot  domain.DashboardSnapshot
	lastErr   error
}

func newModel() model {
	return model{collector: collector.New()}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(tickCmd(), collectCmd(m.collector), tea.WindowSize())
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "r":
			return m, collectCmd(m.collector)
		}
	case tickMsg:
		return m, tea.Batch(tickCmd(), collectCmd(m.collector), tea.WindowSize())
	case snapshotMsg:
		m.lastErr = msg.err
		if msg.err == nil {
			m.snapshot = msg.snapshot
		}
		return m, nil
	}
	return m, nil
}

func (m model) View() string {
	if m.lastErr != nil {
		return "collector error: " + m.lastErr.Error() + "\npress q to quit"
	}
	return ui.Render(m.snapshot, m.width, m.height)
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func collectCmd(c *collector.Collector) tea.Cmd {
	return func() tea.Msg {
		s, err := c.Snapshot()
		return snapshotMsg{snapshot: s, err: err}
	}
}
