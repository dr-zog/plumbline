<!--
Keep it skimmable — a reviewer should grasp the change from this body without reading the
diff first. Full guidance: docs/processes/raising-a-pr.md
-->

## What
<!-- The outcome delivered, in a sentence. -->

## Why
<!-- The requirement / issue this satisfies. Link it, e.g. Closes #123 -->

## How
<!-- The shape of the change. Call out any trade-off honestly. -->

## Proof

- [ ] `go test ./...` passes
- [ ] `go vet ./...` clean
- [ ] `gofmt -l .` clean (no output)
- [ ] `plumbline` self-trace green — **0 gaps** (register-first work is covered)
- [ ] Registered the requirement first, and anchored the covering code-area(s)
- [ ] Added a CHANGELOG entry under `[Unreleased]` (every PR — no exceptions)
- [ ] Recorded an ADR if the decision was architectural
