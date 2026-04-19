package app

import (
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/term"
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

type termSizeMsg struct {
	width  int
	height int
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
	return tea.Batch(tickCmd(), collectCmd(m.collector), tea.WindowSize(), termSizeCmd())
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if msg.Width > 0 && msg.Height > 0 {
			m.width = msg.Width
			m.height = msg.Height
		}
		return m, nil
	case termSizeMsg:
		if msg.width > 0 && msg.height > 0 {
			m.width = msg.width
			m.height = msg.height
		}
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "r":
			return m, collectCmd(m.collector)
		}
	case tickMsg:
		return m, tea.Batch(tickCmd(), collectCmd(m.collector), termSizeCmd())
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
	width, height := m.width, m.height
	if w, h := currentTermSize(); w > 0 && h > 0 {
		width, height = w, h
	}
	return ui.Render(m.snapshot, width, height)
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

func termSizeCmd() tea.Cmd {
	return func() tea.Msg {
		w, h := currentTermSize()
		return termSizeMsg{width: w, height: h}
	}
}

func readTermSize(fd uintptr) (int, int) {
	w, h, err := term.GetSize(fd)
	if err != nil || w <= 0 || h <= 0 {
		return 0, 0
	}
	return w, h
}

func currentTermSize() (int, int) {
	tty, err := os.Open("/dev/tty")
	if err == nil {
		defer tty.Close()
		if w, h := readTermSize(tty.Fd()); w > 0 && h > 0 {
			return w, h
		}
	}
	stdinW, stdinH := readTermSize(os.Stdin.Fd())
	if stdinW > 0 && stdinH > 0 {
		return stdinW, stdinH
	}
	stdoutW, stdoutH := readTermSize(os.Stdout.Fd())
	if stdoutW > 0 && stdoutH > 0 {
		return stdoutW, stdoutH
	}
	return 0, 0
}
