package ui

import (
	"strings"
	"testing"
	"time"

	"github.com/charmbracelet/x/ansi"
	"github.com/minfaatong/mft-micro-cockpit/internal/domain"
)

func TestRenderIncludesHostIP(t *testing.T) {
	s := domain.DashboardSnapshot{
		CollectedAt:   time.Unix(0, 0),
		Hostname:      "node-a",
		HostIP:        "192.168.10.12",
		OSPretty:      "Linux",
		Kernel:        "6.12.0",
		Uptime:        2 * time.Hour,
		Status:        "ok",
		NetInterface:  "eth0",
		MemTotalBytes: 1,
	}

	out := Render(s, 90, 24)
	if !strings.Contains(out, "192.168.10.12") {
		t.Fatalf("Render() missing host IP, output=%q", out)
	}
}

func TestRenderExpandsToRequestedViewport(t *testing.T) {
	s := domain.DashboardSnapshot{
		CollectedAt:    time.Unix(0, 0),
		Hostname:       "node-a",
		HostIP:         "10.0.0.5",
		OSPretty:       "Linux",
		Kernel:         "6.12.0",
		Uptime:         time.Hour,
		CPUPercent:     42,
		MemTotalBytes:  100,
		MemUsedBytes:   50,
		SwapTotalBytes: 100,
		SwapUsedBytes:  20,
		DiskMount:      "/",
		Status:         "ok",
		NetInterface:   "eth0",
	}

	const width = 80
	const height = 32
	out := Render(s, width, height)
	lines := strings.Split(out, "\n")
	if len(lines) != height {
		t.Fatalf("Render() line count = %d, want %d", len(lines), height)
	}
	for i, ln := range lines {
		if got := ansi.StringWidth(ln); got != width {
			t.Fatalf("line %d width = %d, want %d", i, got, width)
		}
	}
}
