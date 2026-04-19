package collector

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/minfaatong/mft-micro-cockpit/internal/domain"
)

var errNoNonLoopbackInterface = errors.New("no non-loopback interface found")

type cpuTimes struct {
	idle  uint64
	total uint64
}

// Collector reads Linux host metrics for dashboard rendering.
type Collector struct {
	lastCPUTimes *cpuTimes
	lastNetAt    time.Time
	lastNetRx    uint64
	lastNetTx    uint64
	lastIface    string
}

// New builds a collector with empty state.
func New() *Collector {
	return &Collector{}
}

// Snapshot reads a complete dashboard sample.
func (c *Collector) Snapshot() (domain.DashboardSnapshot, error) {
	now := time.Now()

	hostname, _ := os.Hostname()
	osPretty, _ := osPrettyName()
	kernel, _ := kernelVersion()
	uptime, _ := readUptime()

	cpuPct, loads, err := c.readCPUAndLoad()
	if err != nil {
		return domain.DashboardSnapshot{}, err
	}

	memAvail, memTotal, swapUsed, swapTotal, err := readMemInfo()
	if err != nil {
		return domain.DashboardSnapshot{}, err
	}

	diskMount, diskUsed, diskTotal, diskPct := readRootDisk()

	iface, rxRate, txRate, err := c.readNetRates(now)
	if err != nil {
		iface = "n/a"
	}
	primaryIP := c.readPrimaryIP(iface)

	s := domain.DashboardSnapshot{
		CollectedAt:     now,
		Hostname:        fallback(hostname, "unknown"),
		PrimaryIP:       fallback(primaryIP, "n/a"),
		OSPretty:        fallback(osPretty, "linux"),
		Kernel:          fallback(kernel, "unknown"),
		Uptime:          uptime,
		CPUPercent:      cpuPct,
		Load1:           loads[0],
		Load5:           loads[1],
		Load15:          loads[2],
		MemUsedBytes:    memTotal - memAvail,
		MemTotalBytes:   memTotal,
		SwapUsedBytes:   swapUsed,
		SwapTotalBytes:  swapTotal,
		DiskMount:       diskMount,
		DiskUsedBytes:   diskUsed,
		DiskTotalBytes:  diskTotal,
		DiskUsedPercent: diskPct,
		NetInterface:    iface,
		NetRxRate:       rxRate,
		NetTxRate:       txRate,
	}
	s.Status = deriveStatus(s)

	return s, nil
}

func (c *Collector) readPrimaryIP(activeIface string) string {
	if activeIface != "" && activeIface != "n/a" {
		if ip := interfaceIPv4(activeIface); ip != "" {
			return ip
		}
	}
	if c.lastIface != "" {
		if ip := interfaceIPv4(c.lastIface); ip != "" {
			return ip
		}
	}

	interfaces, err := net.Interfaces()
	if err != nil {
		return ""
	}
	for _, ifc := range interfaces {
		if ifc.Flags&net.FlagUp == 0 || ifc.Flags&net.FlagLoopback != 0 {
			continue
		}
		if ip := interfaceIPv4(ifc.Name); ip != "" {
			return ip
		}
	}
	return ""
}

func interfaceIPv4(name string) string {
	ifc, err := net.InterfaceByName(name)
	if err != nil {
		return ""
	}
	addrs, err := ifc.Addrs()
	if err != nil {
		return ""
	}
	for _, addr := range addrs {
		var ip net.IP
		switch v := addr.(type) {
		case *net.IPNet:
			ip = v.IP
		case *net.IPAddr:
			ip = v.IP
		}
		if ip == nil {
			continue
		}
		ipv4 := ip.To4()
		if ipv4 == nil || ipv4.IsLoopback() {
			continue
		}
		return ipv4.String()
	}
	return ""
}

func deriveStatus(s domain.DashboardSnapshot) string {
	if s.CPUPercent >= 90 || memPct(s) >= 90 || s.DiskUsedPercent >= 95 {
		return "hot"
	}
	if s.CPUPercent >= 70 || memPct(s) >= 75 || s.DiskUsedPercent >= 85 {
		return "warm"
	}
	return "ok"
}

func memPct(s domain.DashboardSnapshot) float64 {
	if s.MemTotalBytes == 0 {
		return 0
	}
	return (float64(s.MemUsedBytes) / float64(s.MemTotalBytes)) * 100
}

func fallback(v, d string) string {
	if strings.TrimSpace(v) == "" {
		return d
	}
	return v
}

func (c *Collector) readCPUAndLoad() (float64, [3]float64, error) {
	curr, err := readCPUTimes()
	if err != nil {
		return 0, [3]float64{}, err
	}

	loads, err := readLoadAvg()
	if err != nil {
		return 0, [3]float64{}, err
	}

	if c.lastCPUTimes == nil {
		c.lastCPUTimes = &curr
		return 0, loads, nil
	}

	deltaTotal := curr.total - c.lastCPUTimes.total
	deltaIdle := curr.idle - c.lastCPUTimes.idle
	c.lastCPUTimes = &curr

	if deltaTotal == 0 {
		return 0, loads, nil
	}

	busy := float64(deltaTotal-deltaIdle) / float64(deltaTotal)
	return busy * 100, loads, nil
}

func readCPUTimes() (cpuTimes, error) {
	f, err := os.Open("/proc/stat")
	if err != nil {
		return cpuTimes{}, fmt.Errorf("open /proc/stat: %w", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "cpu ") {
			fields := strings.Fields(line)
			if len(fields) < 5 {
				return cpuTimes{}, fmt.Errorf("unexpected cpu line format: %q", line)
			}

			var total uint64
			for _, v := range fields[1:] {
				n, convErr := strconv.ParseUint(v, 10, 64)
				if convErr != nil {
					return cpuTimes{}, fmt.Errorf("parse cpu value: %w", convErr)
				}
				total += n
			}

			idle, errIdle := strconv.ParseUint(fields[4], 10, 64)
			if errIdle != nil {
				return cpuTimes{}, fmt.Errorf("parse cpu idle: %w", errIdle)
			}

			return cpuTimes{idle: idle, total: total}, nil
		}
	}
	if err := scanner.Err(); err != nil {
		return cpuTimes{}, fmt.Errorf("scan /proc/stat: %w", err)
	}

	return cpuTimes{}, errors.New("cpu line not found")
}

func readLoadAvg() ([3]float64, error) {
	var out [3]float64
	raw, err := os.ReadFile("/proc/loadavg")
	if err != nil {
		return out, fmt.Errorf("read /proc/loadavg: %w", err)
	}
	fields := strings.Fields(string(raw))
	if len(fields) < 3 {
		return out, fmt.Errorf("unexpected /proc/loadavg format")
	}
	for i := 0; i < 3; i++ {
		v, convErr := strconv.ParseFloat(fields[i], 64)
		if convErr != nil {
			return out, fmt.Errorf("parse loadavg: %w", convErr)
		}
		out[i] = v
	}
	return out, nil
}

func readMemInfo() (memAvail uint64, memTotal uint64, swapUsed uint64, swapTotal uint64, err error) {
	f, err := os.Open("/proc/meminfo")
	if err != nil {
		return 0, 0, 0, 0, fmt.Errorf("open /proc/meminfo: %w", err)
	}
	defer f.Close()

	vals := map[string]uint64{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		k := strings.TrimSuffix(fields[0], ":")
		v, convErr := strconv.ParseUint(fields[1], 10, 64)
		if convErr != nil {
			continue
		}
		vals[k] = v * 1024
	}
	if err := scanner.Err(); err != nil {
		return 0, 0, 0, 0, fmt.Errorf("scan /proc/meminfo: %w", err)
	}

	memTotal = vals["MemTotal"]
	memAvail = vals["MemAvailable"]
	swapTotal = vals["SwapTotal"]
	swapFree := vals["SwapFree"]

	if memTotal == 0 {
		return 0, 0, 0, 0, errors.New("missing MemTotal in /proc/meminfo")
	}
	if memAvail == 0 {
		// Fallback for environments without MemAvailable.
		memAvail = vals["MemFree"] + vals["Buffers"] + vals["Cached"]
	}
	if swapTotal > swapFree {
		swapUsed = swapTotal - swapFree
	}

	return memAvail, memTotal, swapUsed, swapTotal, nil
}

func readRootDisk() (mount string, used uint64, total uint64, usedPct float64) {
	mount = "/"
	var stat syscall.Statfs_t
	if err := syscall.Statfs("/", &stat); err != nil {
		return mount, 0, 0, 0
	}

	total = stat.Blocks * uint64(stat.Bsize)
	free := stat.Bavail * uint64(stat.Bsize)
	if total > free {
		used = total - free
	}
	if total > 0 {
		usedPct = (float64(used) / float64(total)) * 100
	}

	return mount, used, total, usedPct
}

func (c *Collector) readNetRates(now time.Time) (iface string, rxRate uint64, txRate uint64, err error) {
	iface, rxTotal, txTotal, err := readInterfaceTotals(c.lastIface)
	if err != nil {
		return "", 0, 0, err
	}

	if c.lastNetAt.IsZero() {
		c.lastNetAt = now
		c.lastNetRx = rxTotal
		c.lastNetTx = txTotal
		c.lastIface = iface
		return iface, 0, 0, nil
	}

	seconds := now.Sub(c.lastNetAt).Seconds()
	if seconds <= 0 {
		seconds = 1
	}

	if rxTotal >= c.lastNetRx {
		rxRate = uint64(float64(rxTotal-c.lastNetRx) / seconds)
	}
	if txTotal >= c.lastNetTx {
		txRate = uint64(float64(txTotal-c.lastNetTx) / seconds)
	}

	c.lastNetAt = now
	c.lastNetRx = rxTotal
	c.lastNetTx = txTotal
	c.lastIface = iface

	return iface, rxRate, txRate, nil
}

func readInterfaceTotals(preferred string) (iface string, rx uint64, tx uint64, err error) {
	if preferred != "" {
		rxP, txP, ok := readInterfaceTotalsByName(preferred)
		if ok {
			return preferred, rxP, txP, nil
		}
	}

	interfaces, err := net.Interfaces()
	if err != nil {
		return "", 0, 0, fmt.Errorf("list interfaces: %w", err)
	}

	for _, ifc := range interfaces {
		if ifc.Flags&net.FlagUp == 0 || ifc.Flags&net.FlagLoopback != 0 {
			continue
		}
		rxT, txT, ok := readInterfaceTotalsByName(ifc.Name)
		if !ok {
			continue
		}
		return ifc.Name, rxT, txT, nil
	}

	return "", 0, 0, errNoNonLoopbackInterface
}

func readInterfaceTotalsByName(name string) (rx uint64, tx uint64, ok bool) {
	f, err := os.Open("/proc/net/dev")
	if err != nil {
		return 0, 0, false
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		parts := strings.Split(line, ":")
		if len(parts) != 2 {
			continue
		}
		ifName := strings.TrimSpace(parts[0])
		if ifName != name {
			continue
		}
		stats := strings.Fields(parts[1])
		if len(stats) < 16 {
			return 0, 0, false
		}
		rxV, errRx := strconv.ParseUint(stats[0], 10, 64)
		txV, errTx := strconv.ParseUint(stats[8], 10, 64)
		if errRx != nil || errTx != nil {
			return 0, 0, false
		}
		return rxV, txV, true
	}
	return 0, 0, false
}

func readUptime() (time.Duration, error) {
	raw, err := os.ReadFile("/proc/uptime")
	if err != nil {
		return 0, fmt.Errorf("read /proc/uptime: %w", err)
	}
	fields := strings.Fields(string(raw))
	if len(fields) < 1 {
		return 0, errors.New("unexpected /proc/uptime format")
	}
	seconds, err := strconv.ParseFloat(fields[0], 64)
	if err != nil {
		return 0, fmt.Errorf("parse uptime seconds: %w", err)
	}
	return time.Duration(seconds * float64(time.Second)), nil
}

func osPrettyName() (string, error) {
	raw, err := os.ReadFile("/etc/os-release")
	if err != nil {
		return "", fmt.Errorf("read /etc/os-release: %w", err)
	}
	for _, line := range strings.Split(string(raw), "\n") {
		if strings.HasPrefix(line, "PRETTY_NAME=") {
			v := strings.TrimPrefix(line, "PRETTY_NAME=")
			return strings.Trim(v, `"`), nil
		}
	}
	return "", errors.New("PRETTY_NAME not found")
}

func kernelVersion() (string, error) {
	raw, err := os.ReadFile("/proc/sys/kernel/osrelease")
	if err != nil {
		return "", fmt.Errorf("read kernel osrelease: %w", err)
	}
	return strings.TrimSpace(string(raw)), nil
}
