package ui

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	"github.com/minfaatong/mft-micro-cockpit/internal/domain"
)

var (
	frameStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#5f87af"))
	titleStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#8be9fd"))
	labelStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#61afaf"))
	valueStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#cdd6f4"))
	mutedStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#7f8c98"))
	helpKeyStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#98c379"))
	okStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("#98c379")).Bold(true)
	warmStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#f2cc8f")).Bold(true)
	hotStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#e06c75")).Bold(true)
	barCoolStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#56b6c2")).Bold(true)
	barOkStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#98c379")).Bold(true)
	barWarmStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#e5c07b")).Bold(true)
	barHotStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#e06c75")).Bold(true)
)

// Render returns the dashboard view tailored for compact terminals.
func Render(s domain.DashboardSnapshot, width, height int) string {
	_ = height // future dynamic layout usage
	if width <= 0 {
		width = 80
	}
	innerWidth := width - 2
	if innerWidth < 20 {
		innerWidth = 20
	}
	meterWidth := meterWidthFor(innerWidth)
	hostWidth := maxInt(8, minInt(16, innerWidth/4))
	osWidth := maxInt(10, minInt(28, innerWidth/3))
	kernelWidth := maxInt(8, minInt(20, innerWidth/4))
	mountWidth := maxInt(4, minInt(12, innerWidth/6))
	ifaceWidth := maxInt(4, minInt(8, innerWidth/10))
	status := styleStatus(s.Status)
	cpuPct := pctFromFloat(s.CPUPercent)
	memPct := pct(s.MemUsedBytes, s.MemTotalBytes)
	swapPct := pct(s.SwapUsedBytes, s.SwapTotalBytes)
	diskPct := pctFromFloat(s.DiskUsedPercent)

	headerTitle := fmt.Sprintf(
		"%s %s",
		titleStyle.Render("MICRO-COCKPIT"),
		mutedStyle.Render("linux telemetry"),
	)
	headerMeta := fmt.Sprintf(
		"%s %s  %s %s",
		labelStyle.Render("host"),
		valueStyle.Render(clampText(s.Hostname, hostWidth)),
		labelStyle.Render("status"),
		status,
	)
	systemMeta := fmt.Sprintf(
		"%s %s  %s %s",
		labelStyle.Render("os"),
		valueStyle.Render(clampText(s.OSPretty, osWidth)),
		labelStyle.Render("kernel"),
		valueStyle.Render(clampText(s.Kernel, kernelWidth)),
	)
	timeMeta := fmt.Sprintf(
		"%s %s  %s %s",
		labelStyle.Render("uptime"),
		valueStyle.Render(formatDuration(s.Uptime)),
		labelStyle.Render("ts"),
		valueStyle.Render(s.CollectedAt.Format("15:04:05")),
	)

	cpuLine := fmt.Sprintf(
		"%s %s %s  %s %4.2f %4.2f %4.2f",
		labelStyle.Render("cpu "),
		meter(cpuPct, meterWidth),
		stylePercent(cpuPct).Render(fmt.Sprintf("%6.2f%%", cpuPct)),
		mutedStyle.Render("load"),
		s.Load1, s.Load5, s.Load15,
	)
	memLine := fmt.Sprintf(
		"%s %s %s  %s %s/%s",
		labelStyle.Render("mem "),
		meter(memPct, meterWidth),
		stylePercent(memPct).Render(fmt.Sprintf("%6.2f%%", memPct)),
		mutedStyle.Render("used"),
		valueStyle.Render(bytesHuman(s.MemUsedBytes)),
		valueStyle.Render(bytesHuman(s.MemTotalBytes)),
	)
	swapLine := fmt.Sprintf(
		"%s %s %s  %s %s/%s",
		labelStyle.Render("swp "),
		meter(swapPct, meterWidth),
		stylePercent(swapPct).Render(fmt.Sprintf("%6.2f%%", swapPct)),
		mutedStyle.Render("used"),
		valueStyle.Render(bytesHuman(s.SwapUsedBytes)),
		valueStyle.Render(bytesHuman(s.SwapTotalBytes)),
	)
	diskLine := fmt.Sprintf(
		"%s %s %s  %s %s  %s/%s",
		labelStyle.Render("dsk "),
		meter(diskPct, meterWidth),
		stylePercent(diskPct).Render(fmt.Sprintf("%6.2f%%", diskPct)),
		mutedStyle.Render(clampText(s.DiskMount, mountWidth)),
		valueStyle.Render(bytesHuman(s.DiskUsedBytes)),
		valueStyle.Render(bytesHuman(s.DiskTotalBytes)),
	)
	netLine := fmt.Sprintf(
		"%s %-8s  %s %s/s  %s %s/s",
		labelStyle.Render("net"),
		valueStyle.Render(clampText(s.NetInterface, ifaceWidth)),
		labelStyle.Render("rx"),
		valueStyle.Render(bytesHuman(s.NetRxRate)),
		labelStyle.Render("tx"),
		valueStyle.Render(bytesHuman(s.NetTxRate)),
	)

	help := fmt.Sprintf(
		"%s quit  |  %s refresh",
		helpKeyStyle.Render("q/ctrl+c"),
		helpKeyStyle.Render("r"),
	)

	body := []string{
		frameTop(innerWidth),
		frameContent(headerTitle, innerWidth),
		frameContent(headerMeta, innerWidth),
		frameContent(systemMeta, innerWidth),
		frameContent(timeMeta, innerWidth),
		frameSep(innerWidth),
		frameContent(cpuLine, innerWidth),
		frameContent(memLine, innerWidth),
		frameContent(swapLine, innerWidth),
		frameContent(diskLine, innerWidth),
		frameContent(netLine, innerWidth),
		frameSep(innerWidth),
		frameContent(help, innerWidth),
		frameBottom(innerWidth),
	}
	return truncateLines(strings.Join(body, "\n"), width)
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

func meter(percent float64, width int) string {
	p := pctFromFloat(percent)
	filled := int(math.Round((p / 100) * float64(width)))
	if filled < 0 {
		filled = 0
	}
	if filled > width {
		filled = width
	}
	var b strings.Builder
	b.WriteString(frameStyle.Render("╞"))
	for i := range width {
		if i < filled {
			b.WriteString(gradientBarStyle(i, width).Render("█"))
			continue
		}
		b.WriteString(mutedStyle.Render("░"))
	}
	b.WriteString(frameStyle.Render("╡"))
	return b.String()
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
		lines[i] = ansi.Truncate(ln, width, "")
	}
	return strings.Join(lines, "\n")
}

func frameTop(innerWidth int) string {
	return frameStyle.Render("╔" + strings.Repeat("═", innerWidth) + "╗")
}

func frameSep(innerWidth int) string {
	return frameStyle.Render("╟" + strings.Repeat("─", innerWidth) + "╢")
}

func frameBottom(innerWidth int) string {
	return frameStyle.Render("╚" + strings.Repeat("═", innerWidth) + "╝")
}

func frameContent(content string, innerWidth int) string {
	line := padVisible(content, innerWidth)
	return frameStyle.Render("║") + line + frameStyle.Render("║")
}

func padVisible(s string, width int) string {
	if width <= 0 {
		return ""
	}
	line := ansi.Truncate(s, width, "")
	padding := width - ansi.StringWidth(line)
	if padding < 0 {
		padding = 0
	}
	return line + strings.Repeat(" ", padding)
}

func meterWidthFor(innerWidth int) int {
	switch {
	case innerWidth >= 98:
		return 34
	case innerWidth >= 84:
		return 26
	case innerWidth >= 72:
		return 20
	default:
		return 14
	}
}

func gradientBarStyle(idx, width int) lipgloss.Style {
	ratio := float64(idx+1) / float64(maxInt(1, width))
	switch {
	case ratio < 0.35:
		return barCoolStyle
	case ratio < 0.7:
		return barOkStyle
	case ratio < 0.9:
		return barWarmStyle
	default:
		return barHotStyle
	}
}

func stylePercent(percent float64) lipgloss.Style {
	p := pctFromFloat(percent)
	switch {
	case p >= 90:
		return hotStyle
	case p >= 70:
		return warmStyle
	default:
		return okStyle
	}
}

func clampText(s string, maxWidth int) string {
	if maxWidth <= 0 {
		return ""
	}
	trimmed := strings.TrimSpace(s)
	if trimmed == "" {
		return "-"
	}
	if ansi.StringWidth(trimmed) <= maxWidth {
		return trimmed
	}
	if maxWidth == 1 {
		return "…"
	}
	return ansi.Truncate(trimmed, maxWidth-1, "") + "…"
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
