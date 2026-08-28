# Architecture Decision Records

Short, durable records of decisions that shape Plumbline — what we chose, why,
what we rejected, and where the reasoning came from. Read one to understand a
choice without re-deriving it; write one when a decision commits us to a path
we'd otherwise re-litigate later.

## Convention

- One file per decision: `docs/adrs/NNN-kebab-title.md`, `NNN` zero-padded, monotonic.
- Never renumber or delete an ADR. Superseded decisions get a new ADR and the old
  one's **Status** is changed to `Superseded by NNN`.
- Keep them terse but self-contained — a reader in a year should need nothing else.

## Template

```markdown
# NNN. Title

- **Status:** Proposed | Accepted | Superseded by NNN
- **Date:** YYYY-MM-DD
- **Deciders:** …

## Context
The forces at play — the problem, the constraints, what we knew.

## Decision
What we chose, stated plainly.

## Consequences
The trade-offs we accept — good and bad.

## Alternatives considered
What else we weighed, and why we rejected it.

## Provenance
Sources and evidence behind the decision, so it can be re-checked.
```

## Index

- [001 — Plugin packaging & distribution](001-plugin-packaging-and-distribution.md)
- [002 — Traceability only; emergent-structure analysis out of scope](002-traceability-only-emergent-analysis-out-of-scope.md)
- [003 — OpenFastTrace conformance boundary](003-oft-conformance-boundary.md)
- [004 — Requirement status lifecycle](004-requirement-status-lifecycle.md)
- [005 — Automated releases via semantic-release](005-automated-releases-semantic-release.md)
- [006 — GitHub Release binaries for pinned consumption](006-release-binaries-for-ci.md)
