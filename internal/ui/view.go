package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/minfaatong/mft-micro-cockpit/internal/domain"
)

var (
	titleStyle = lipgloss.NewStyle().Bold(true)
	okStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Bold(true)
	warmStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("3")).Bold(true)
	hotStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("1")).Bold(true)
)

// Render returns the dashboard view tailored for compact terminals.
func Render(s domain.DashboardSnapshot, width, height int) string {
	_ = height // future dynamic layout usage

	status := styleStatus(s.Status)
	header := fmt.Sprintf(
		"micro-cockpit | host: %s | status: %s\nos: %s | kernel: %s\nuptime: %s | ts: %s",
		s.Hostname,
		status,
		s.OSPretty,
		s.Kernel,
		formatDuration(s.Uptime),
		s.CollectedAt.Format("15:04:05"),
	)

	cpuLine := fmt.Sprintf(
		"CPU  %6.2f%%  load %.2f %.2f %.2f",
		s.CPUPercent, s.Load1, s.Load5, s.Load15,
	)
	memLine := fmt.Sprintf(
		"MEM  %6.2f%%  %s / %s",
		pct(s.MemUsedBytes, s.MemTotalBytes),
		bytesHuman(s.MemUsedBytes),
		bytesHuman(s.MemTotalBytes),
	)
	swapLine := fmt.Sprintf(
		"SWP  %6.2f%%  %s / %s",
		pct(s.SwapUsedBytes, s.SwapTotalBytes),
		bytesHuman(s.SwapUsedBytes),
		bytesHuman(s.SwapTotalBytes),
	)
	diskLine := fmt.Sprintf(
		"DSK  %6.2f%%  %s  %s / %s",
		s.DiskUsedPercent,
		s.DiskMount,
		bytesHuman(s.DiskUsedBytes),
		bytesHuman(s.DiskTotalBytes),
	)
	netLine := fmt.Sprintf(
		"NET  %-8s rx %s/s  tx %s/s",
		s.NetInterface,
		bytesHuman(s.NetRxRate),
		bytesHuman(s.NetTxRate),
	)

	help := "keys: q quit | ctrl+c quit | r refresh"

	body := strings.Join([]string{
		titleStyle.Render(header),
		"",
		cpuLine,
		bar(pctFromFloat(s.CPUPercent), width),
		memLine,
		bar(pct(s.MemUsedBytes, s.MemTotalBytes), width),
		swapLine,
		diskLine,
		netLine,
		"",
		help,
	}, "\n")

	return truncateLines(body, width)
}

func styleStatus(status string) string {
	switch status {
	case "hot":
		return hotStyle.Render(strings.ToUpper(status))
	case "warm":
		return warmStyle.Render(strings.ToUpper(status))
	default:
		return okStyle.Render(strings.ToUpper(status))
	}
}

func formatDuration(d time.Duration) string {
	d = d.Round(time.Second)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second
	return fmt.Sprintf("%02dh%02dm%02ds", h, m, s)
}

func pct(used, total uint64) float64 {
	if total == 0 {
		return 0
	}
	return (float64(used) / float64(total)) * 100
}

func pctFromFloat(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > 100 {
		return 100
	}
	return v
}

func bar(percent float64, width int) string {
	w := 30
	if width > 0 {
		if width-10 < w {
			w = width - 10
		}
		if w < 10 {
			w = 10
		}
	}

	filled := int((percent / 100) * float64(w))
	if filled < 0 {
		filled = 0
	}
	if filled > w {
		filled = w
	}

	return fmt.Sprintf("[%s%s] %5.1f%%", strings.Repeat("#", filled), strings.Repeat("-", w-filled), percent)
}

func bytesHuman(b uint64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%dB", b)
	}
	div, exp := uint64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f%ciB", float64(b)/float64(div), "KMGTPE"[exp])
}

func truncateLines(s string, width int) string {
	if width <= 0 {
		return s
	}
	lines := strings.Split(s, "\n")
	for i, ln := range lines {
		if len([]rune(ln)) > width {
			runes := []rune(ln)
			lines[i] = string(runes[:width])
		}
	}
	return strings.Join(lines, "\n")
}
