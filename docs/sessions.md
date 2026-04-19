# Session Checklists

## A) Brainstorm Checklist
- [ ] Clarify main usage mode (local console vs SSH).
- [ ] Confirm minimum terminal size requirement.
- [ ] Rank MVP metrics by user value.
- [ ] Decide implementation language + library.
- [ ] Define explicit non-goals for current milestone.

## B) Development Checklist
- [ ] Create issue/milestone objective.
- [ ] Implement smallest vertical slice first.
- [ ] Keep collector and UI layers decoupled.
- [ ] Add/adjust tests for changed behavior.
- [ ] Validate layout in compact mode.

## C) Testing Checklist
- [ ] Run unit tests for collector/domain.
- [ ] Run app manually in `80x24` terminal.
- [ ] Verify keyboard interactions (`q`, `?`, navigation).
- [ ] Validate missing data fallback labels.
- [ ] Measure approximate CPU/memory overhead.

## D) Deployment Checklist
- [ ] Build release binary for target arch.
- [ ] Verify execution permissions.
- [ ] Validate startup command and sample config.
- [ ] Capture screenshot/video for release note.
- [ ] Confirm rollback instructions are documented.

## E) Early Command Plan (Go Path)

```bash
# bootstrap
mkdir -p cmd/micro-cockpit internal/{app,collector,domain,ui} configs docs
go mod init github.com/minfaatong/mft-micro-cockpit

# dependencies (tentative)
go get github.com/charmbracelet/bubbletea
go get github.com/charmbracelet/lipgloss
go get github.com/shirou/gopsutil/v4

# run
go run ./cmd/micro-cockpit

# test
go test ./...
```

## F) First Implementation Slice
1. Render header + uptime + CPU total + memory usage.
2. Refresh every 1 second.
3. Support quit key (`q` or `ctrl+c`).
4. Confirm clean display on `80x24`.
