# Writing ADRs

When and how to record an architecture decision.

## When it's an ADR (and when it isn't)

Write an ADR for a decision that:

- **commits us to a path** we'd otherwise re-litigate later,
- is **architectural / structural** (how the thing is built, packaged, distributed),
  and
- is **costly to reverse** or whose *rationale* matters more than the choice itself.

Do **not** write an ADR for a **process** (how we work) — that's a living practice and
lives in [`docs/processes/`](README.md). Rule of thumb: *a decision with provenance →
ADR; a way of working → process doc.*

Examples that earned an ADR: the plugin packaging & distribution shape
([001](../adrs/001-plugin-packaging-and-distribution.md)). Examples that did **not**:
"register-first within a change" (that's a process).

## How

- One file per decision: `docs/adrs/NNN-kebab-title.md`, `NNN` zero-padded and monotonic.
- Follow the template in [`docs/adrs/README.md`](../adrs/README.md): Context · Decision ·
  Consequences · Alternatives considered · Provenance.
- **Provenance is the point** — cite the sources and evidence so a future reader can
  re-check the reasoning, not just the conclusion. Record the alternatives you *rejected*
  and why; that's usually the most valuable part.
- Keep it terse and self-contained. A reader in a year should need nothing else.

## Lifecycle

- ADRs are immutable once **Accepted** — never renumber or delete one.
- A decision that changes gets a **new** ADR; the old one's **Status** becomes
  `Superseded by NNN`, and the new one references what it replaces.
- Update the index in `docs/adrs/README.md`.
