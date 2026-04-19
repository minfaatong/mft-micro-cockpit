package collector

import (
	"testing"

	"github.com/minfaatong/mft-micro-cockpit/internal/domain"
)

func TestDeriveStatusThresholds(t *testing.T) {
	tests := []struct {
		name string
		s    domain.DashboardSnapshot
		want string
	}{
		{
			name: "ok",
			s: domain.DashboardSnapshot{
				CPUPercent:      20,
				MemUsedBytes:    1,
				MemTotalBytes:   4,
				DiskUsedPercent: 30,
			},
			want: "ok",
		},
		{
			name: "warm by cpu",
			s: domain.DashboardSnapshot{
				CPUPercent:      75,
				MemUsedBytes:    1,
				MemTotalBytes:   4,
				DiskUsedPercent: 30,
			},
			want: "warm",
		},
		{
			name: "hot by cpu",
			s: domain.DashboardSnapshot{
				CPUPercent:      95,
				MemUsedBytes:    1,
				MemTotalBytes:   4,
				DiskUsedPercent: 30,
			},
			want: "hot",
		},
		{
			name: "warm by memory",
			s: domain.DashboardSnapshot{
				CPUPercent:      20,
				MemUsedBytes:    8,
				MemTotalBytes:   10,
				DiskUsedPercent: 30,
			},
			want: "warm",
		},
		{
			name: "hot by disk",
			s: domain.DashboardSnapshot{
				CPUPercent:      20,
				MemUsedBytes:    1,
				MemTotalBytes:   10,
				DiskUsedPercent: 97,
			},
			want: "hot",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := deriveStatus(tc.s)
			if got != tc.want {
				t.Fatalf("deriveStatus() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestMemPct(t *testing.T) {
	s := domain.DashboardSnapshot{MemUsedBytes: 5, MemTotalBytes: 10}
	if got := memPct(s); got != 50 {
		t.Fatalf("memPct() = %.2f, want 50", got)
	}

	s.MemTotalBytes = 0
	if got := memPct(s); got != 0 {
		t.Fatalf("memPct() for zero total = %.2f, want 0", got)
	}
}
