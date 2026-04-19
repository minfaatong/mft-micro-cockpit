# Minimal POD: Build-Operate-Deliver Loop

This POD defines a lightweight operating model for this project.

## POD Roles
- **Product Owner (PO)**: sets feature priority and acceptance criteria.
- **Developer (Dev)**: implements features, refactors, writes tests.
- **Tester (QA)**: validates behavior in terminal sizes and Linux targets.
- **Release/Deploy (Ops)**: packages binaries and rollout docs.

For small teams, one person may hold multiple roles.

## POD Cadence

### 1) Brainstorm Session
Input:
- user goals (portable, low-res LCD, glance monitoring)

Output:
- prioritized feature list
- constraints and non-goals
- candidate technical stack

Artifact:
- `docs/brainstorm.md`

### 2) Development Session
Input:
- selected milestone and acceptance criteria

Process:
- short design note
- implementation in small PRs
- local verification after each unit of work

Output:
- working increment
- changelog entry in PR

### 3) Testing Session
Input:
- build artifact and test checklist

Process:
- automated tests (collector parsing, model updates)
- manual TUI checks (`80x24`, `100x30`, resize)
- resource overhead check on target host

Output:
- pass/fail report
- bug list (if any)

### 4) Deployment Session
Input:
- tested tagged commit

Process:
- produce release binaries
- publish checksums and usage instructions
- smoke run on target machine

Output:
- release note
- rollback note

## Definition of Done (per milestone)
- Acceptance criteria satisfied.
- Tests executed and recorded.
- Binary runs on at least one target environment.
- Docs updated (`README` + relevant docs in `docs/`).

## Communication Templates

### Development Kickoff
- Goal:
- Scope in/out:
- Risks:
- Test plan:

### Test Report
- Build version:
- Environment:
- Test cases:
- Result summary:
- Issues found:

### Deployment Note
- Version:
- Installed on:
- Validation commands:
- Observed output:
- Rollback path:
