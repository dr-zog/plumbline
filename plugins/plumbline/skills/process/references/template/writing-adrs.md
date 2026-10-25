# Writing ADRs

When and how to record an architecture decision.

## When it's an ADR, and when it isn't

Write an ADR for a decision that:

- **commits us to a path** we'd otherwise re-litigate later,
- is **architectural or structural** (how the thing is built, packaged, distributed), and
- is **costly to reverse**, or whose *rationale* matters more than the choice itself.

Do **not** write an ADR for a **process** (how we work). That's a living practice and lives in
these process docs. Rule of thumb: *a decision with provenance → ADR; a way of working →
process doc.*

## How

- One file per decision: `{{ADR_DIR}}/NNN-kebab-title.md`, with `NNN` zero-padded and
  monotonic. Keep an index in `{{ADR_DIR}}/README.md`.
- Sections: **Context**, **Decision**, **Consequences**, **Alternatives considered**,
  **Provenance**. Head it with **Status** (Proposed, Accepted, or Superseded by NNN) and
  **Date**.
- **Provenance is the point.** Cite the sources and evidence, so a future reader can re-check
  the reasoning and not just the conclusion. Record the alternatives you *rejected* and why;
  that's usually the most valuable part.
- Keep it terse and self-contained. A reader in a year should need nothing else.
- **Reference, don't close.** An ADR is a decision record, not a requirement. Write it on its
  own branch. Its {{CHANGE}} *references* the tickets it governs but doesn't close them; they
  close when their implementing code lands.

## Lifecycle

- ADRs are immutable once **Accepted**. Never renumber or delete one.
- A decision that changes gets a **new** ADR. The old one's **Status** becomes
  `Superseded by NNN`, and the new one references what it replaces.
