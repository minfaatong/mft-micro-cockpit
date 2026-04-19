# Project Plan: `micro-cockpit` TUI Monitor

## 1) Scope Definition

### MVP Scope
Build a local Linux TUI app that displays:
- system identity (host, OS, uptime)
- CPU, memory, disk, network summary
- top process list
- threshold-based health indicator

### Target Environment
- Linux x86_64 and ARM64
- SSH terminals and local console
- minimum usable viewport: `80x24`
- target display scenario: `800x600` mini LCD

## 2) Delivery Phases

### Phase 0: Foundation
- Initialize repository structure.
- Pick language/runtime (recommended: Go).
- Build a simple app loop with placeholder panels.

### Phase 1: Collectors + Basic UI
- Implement CPU/memory/uptime collectors.
- Render stable compact dashboard.
- Add refresh loop and keyboard quit/help.

### Phase 2: Expanded Metrics
- Add disk/network/process collectors.
- Add warning thresholds and status colors.
- Add simple config file support.

### Phase 3: Hardening
- Improve error handling and fallback values.
- Add tests, benchmark, and packaging.
- Validate on at least one low-power Linux target.

## 3) Technical Decisions (Initial)
- **Language**: Go
- **UI library**: Bubble Tea + Lip Gloss + Bubbles (or direct ratatui equivalent if Rust path)
- **Sampling**: default `1s` refresh interval, configurable
- **Config**: `~/.config/micro-cockpit/config.yaml`

## 4) Suggested Repo Layout

```text
cmd/micro-cockpit/main.go
internal/app/
internal/collector/
internal/domain/
internal/ui/
configs/example.yaml
docs/
```

## 5) Milestone Acceptance Criteria

### M1 (Skeleton)
- App starts and exits cleanly.
- Base layout renders in `80x24`.

### M2 (Core Metrics)
- CPU/memory/uptime update correctly every second.
- No flicker or layout break under resize.

### M3 (Full MVP)
- Disk/network/top-process displayed.
- Health states and warnings visible.
- Binary runs without external dependencies.

## 6) Implementation Notes for Tiny Displays
- Prioritize one-screen layout, no deep menus.
- Use single-character units where possible (`C`, `M`, `D`, `N`, `P`).
- Keep left alignment and avoid dense borders.
- Ensure high-contrast default theme.

## 7) Backlog (Post-MVP)
- Optional remote mode (read-only pull from node over SSH).
- Historical sparkline per metric.
- Plugin panels (temperature, Docker, systemd units).
- Export snapshot as text/JSON.
