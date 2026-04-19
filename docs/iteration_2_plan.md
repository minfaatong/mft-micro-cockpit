# Iteration 2 Plan: IP Visibility + Visual Facelift

## Request Recap
1. Show host IP address in the app.
2. Improve visual polish so dashboard is more pleasant to glance at.

## Success State
- Dashboard clearly shows at least one usable host IP address without breaking compact layout.
- UI has improved hierarchy and readability (section cards/panels, improved spacing, clearer status emphasis).
- App still runs in compact terminal sizes and keeps low-friction keyboard controls.

## Technical Plan

### A) Data model + collector updates
- Add `PrimaryIP` field to dashboard snapshot.
- Reuse active network interface selection logic to derive matching IPv4 address.
- Fallback to `n/a` when no non-loopback IP is available.

### B) UI facelift
- Introduce stylized sections with subtle borders and spacing using `lipgloss`.
- Add stronger title/header treatment and status badge styling.
- Split metrics into compact cards:
  - System card (host, IP, OS, kernel, uptime)
  - Resource card (CPU/MEM/SWAP bars)
  - Infra card (disk, network throughput)
- Ensure truncation still respects narrow terminals.

### C) Validation plan
- Automated: run `go test ./...`.
- Runtime: run app with timeout and capture output log for collector/UI proof.
- Manual UI: produce fresh demo video showing IP line and improved layout.

## Risks
- Lipgloss styling may consume width quickly on `80x24`.
  - Mitigation: keep one-column card stack and light borders.
- Interface/IP detection may return container-local address.
  - Mitigation: show selected interface + deterministic fallback to `n/a`.

## Expected Artifacts
- Updated plan doc (`docs/iteration_2_plan.md`)
- Runtime log artifact
- Demo video + screenshot artifact
