# mft-micro-cockpit

Portable Linux TUI monitor inspired by `neofetch` + `htop`, optimized for glanceable status in compact displays and SSH sessions.

## Current MVP
- live dashboard refresh every 1 second
- host, OS, kernel, uptime
- CPU + load averages
- memory + swap usage
- root disk usage
- network throughput for active interface
- status level (`OK`, `WARM`, `HOT`)

## Requirements
- Linux (uses `/proc` and `statfs`)
- Go `1.22+`

## Run

```bash
go run ./cmd/micro-cockpit
```

Keys:
- `q` or `ctrl+c`: quit
- `r`: force refresh

## Development

```bash
go test ./...
```

## Notes
- Layout is compact-first and targets `80x24` minimum terminal size.
- For `800x600` mini LCD scenarios, use a terminal font size that keeps at least `80x24` visible.
