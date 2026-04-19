# Brainstorm: Portable Linux TUI Monitor

## Problem Statement
Build a small-footprint terminal UI app that gives at-a-glance server/device info, suitable for:
- `800x600` mini LCD dashboards
- SSH sessions on IoT servers
- Always-on local terminal sessions

The experience should feel like a blend of:
- **neofetch** (quick system identity summary)
- **htop** (live resource view)
- modern keyboard-first TUI apps (clear panels, responsive navigation)

## Product Goals
1. **Readable at low resolution**: fit key data in `80x24` and scale up gracefully.
2. **Low overhead**: minimal CPU/memory impact so monitoring does not become load.
3. **Portable**: run on major Linux distros and lightweight IoT environments.
4. **Actionable glance**: health status and hotspots visible within 3 seconds.

## Non-goals (Initial Scope)
- Full observability platform replacement (Prometheus/Grafana class).
- Multi-node distributed backend in v1.
- Agent-based remote execution in v1.

## Potential Feature Set (v1 Candidate)

### Core Panels
- **Header**: hostname, uptime, distro, kernel, IP(s), temperature snapshot.
- **CPU panel**: overall + per-core usage bars, load average.
- **Memory panel**: used/free/cache/swap bars.
- **Disk panel**: top mount usage and free space warning.
- **Network panel**: rx/tx throughput, interface status.
- **Process panel**: top N CPU or RAM consumers.

### UX Ideas
- Keyboard navigation with vim-like shortcuts (`j/k`, `tab`, `?` help).
- Compact mode for `80x24` and expanded mode for larger terminals.
- Theme presets (dark/mono/high-contrast for small LCD readability).
- Minimal status colors: green/yellow/red thresholds.

### Data Sources (Linux)
- `/proc/stat`, `/proc/meminfo`, `/proc/loadavg`
- `/proc/net/dev`
- `/proc/uptime`
- `/etc/os-release`
- `sysfs` for thermal data when available
- Optional fallback shell commands for portability

## Technical Direction Options

### Option A: Go + Bubble Tea (Recommended)
Pros:
- Great TUI ecosystem, static binaries, easy cross-compilation.
- Good performance and deployment simplicity.

Cons:
- Slightly more architecture upfront.

### Option B: Rust + Ratatui
Pros:
- Very high performance and strong type safety.

Cons:
- Longer initial velocity for rapid prototyping.

### Option C: Python + Textual
Pros:
- Fast prototyping, rich UI components.

Cons:
- Heavier runtime/deployment for tiny IoT systems.

## Suggested Initial Architecture
- **collector/**: pull system stats at fixed interval.
- **domain/**: normalized models (`CPUStats`, `MemoryStats`, etc.).
- **ui/**: layout + keyboard interactions.
- **app/**: update loop, state store, render scheduling.

## Risks & Mitigations
- **Tiny terminals clipping content** -> design strict compact layout first.
- **Different Linux flavors missing sensors** -> optional panel fallback + N/A labels.
- **Resource overhead** -> benchmark update intervals (250ms, 500ms, 1s).
- **Permission constraints** -> avoid root-required metrics in v1.

## Success Criteria for v1
- Runs in `80x24` without broken layout.
- Keeps CPU overhead under ~2% on low-power target at 1s refresh.
- Shows core health metrics in one screen.
- Packaged as single executable + simple config file.
