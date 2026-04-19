package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/minfaatong/mft-micro-cockpit/internal/domain"
)

var (
	appStyle = lipgloss.NewStyle().Padding(0, 1)

	titleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("39"))
	subtle     = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))

	okBadge   = lipgloss.NewStyle().Foreground(lipgloss.Color("46")).Bold(true)
	warmBadge = lipgloss.NewStyle().Foreground(lipgloss.Color("220")).Bold(true)
	hotBadge  = lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true)

	cardStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240")).
			Padding(0, 1)
)

// Render returns the dashboard view tailored for compact terminals.
func Render(s domain.DashboardSnapshot, width, height int) string {
	_ = height // reserved for future adaptive layout

	title := titleStyle.Render("micro-cockpit")
	top := fmt.Sprintf("%s  %s  %s", title, styleStatus(s.Status), subtle.Render(s.CollectedAt.Format("15:04:05")))

	systemLines := []string{
		fmt.Sprintf("host   %s", s.Hostname),
		fmt.Sprintf("ip     %s", s.PrimaryIP),
		fmt.Sprintf("os     %s", s.OSPretty),
		fmt.Sprintf("kernel %s", s.Kernel),
		fmt.Sprintf("uptime %s", formatDuration(s.Uptime)),
	}

	resourceLines := []string{
		fmt.Sprintf("CPU   %6.2f%%  load %.2f %.2f %.2f", s.CPUPercent, s.Load1, s.Load5, s.Load15),
		bar(pctFromFloat(s.CPUPercent), width),
		fmt.Sprintf("MEM   %6.2f%%  %s / %s", pct(s.MemUsedBytes, s.MemTotalBytes), bytesHuman(s.MemUsedBytes), bytesHuman(s.MemTotalBytes)),
		bar(pct(s.MemUsedBytes, s.MemTotalBytes), width),
		fmt.Sprintf("SWAP  %6.2f%%  %s / %s", pct(s.SwapUsedBytes, s.SwapTotalBytes), bytesHuman(s.SwapUsedBytes), bytesHuman(s.SwapTotalBytes)),
	}

	infraLines := []string{
		fmt.Sprintf("DISK  %6.2f%%  %s  %s / %s", s.DiskUsedPercent, s.DiskMount, bytesHuman(s.DiskUsedBytes), bytesHuman(s.DiskTotalBytes)),
		fmt.Sprintf("NET   %-8s  rx %s/s  tx %s/s", s.NetInterface, bytesHuman(s.NetRxRate), bytesHuman(s.NetTxRate)),
	}

	systemCard := renderCard("system", strings.Join(systemLines, "\n"), width)
	resourceCard := renderCard("resources", strings.Join(resourceLines, "\n"), width)
	infraCard := renderCard("infra", strings.Join(infraLines, "\n"), width)

	help := subtle.Render("keys: q quit | ctrl+c quit | r refresh")

	body := strings.Join([]string{top, "", systemCard, resourceCard, infraCard, help}, "\n")
	out := appStyle.Render(body)
	return truncateLines(out, width)
}

func renderCard(title, content string, width int) string {
	header := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("45")).Render(strings.ToUpper(title))
	cardText := header + "\n" + content
	w := width - 4
	if w < 36 {
		w = 36
	}
	return cardStyle.Width(w).Render(cardText)
}

func styleStatus(status string) string {
	label := strings.ToUpper(status)
	switch status {
	case "hot":
		return hotBadge.Render("[" + label + "]")
	case "warm":
		return warmBadge.Render("[" + label + "]")
	default:
		return okBadge.Render("[" + label + "]")
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
	w := 24
	if width > 0 {
		if width-24 < w {
			w = width - 24
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
	return fmt.Sprintf("      [%s%s] %5.1f%%", strings.Repeat("#", filled), strings.Repeat("-", w-filled), percent)
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
