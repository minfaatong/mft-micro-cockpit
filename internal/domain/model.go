package domain

import "time"

// DashboardSnapshot stores one polling sample for the TUI.
type DashboardSnapshot struct {
	CollectedAt time.Time
	Hostname    string
	OSPretty    string
	Kernel      string
	Uptime      time.Duration

	CPUPercent float64
	Load1      float64
	Load5      float64
	Load15     float64

	MemUsedBytes   uint64
	MemTotalBytes  uint64
	SwapUsedBytes  uint64
	SwapTotalBytes uint64

	DiskMount       string
	DiskUsedBytes   uint64
	DiskTotalBytes  uint64
	DiskUsedPercent float64

	NetInterface string
	NetRxRate    uint64
	NetTxRate    uint64

	Status string
}
