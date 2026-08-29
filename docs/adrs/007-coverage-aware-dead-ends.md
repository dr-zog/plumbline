# 007. Coverage-aware dead-ends; the architecture axis is first-class

- **Status:** Accepted
- **Date:** 2026-08-29
- **Deciders:** Dr. Zog (operator), analysis by Claude; verified feedback from a consuming AI

## Context

ADR 004 introduced the **dead-end**: an `approved` item that declares no `Needs` fails,
because it "can never be genuinely covered" and would otherwise fake a vacuous 100%.
It was implemented as `approved && len(Needs) == 0 → dead-end`, unconditionally —
never consulting whether anything *covers* the item.

A consuming team applied the v0.2.1 status model to a full C4 register (context,
containers, and `feat → req → component → impl`) and hit a wall marking the
architecture spine `approved`. Both cases were reproduced against v0.2.1:

- A **`context`** with no `Needs` but covered by its containers (5 coverers) still
  dead-ends — the rule ignores incoming coverage.
- A **`container`** (`Needs: component`) whose components cover only their `req` — the
  ordinary OFT practice — is **uncovered**, because nothing covers the container from
  below.

Two things were conflated. The rule reads "declares no `Needs`" as "can never be
covered," but a **top-of-axis node** (a `context`) is covered *from below* by
aggregation, not by needing something down. And the **architecture axis**
(`context → container → component`) requires each component to declare its **container
membership** — which Plumbline's own register has always done
(`component Covers: container~engine~1, req~…~1`) but which the skills never teach.
So a competent C4/OFT practitioner, wiring `component → req` per the OFT norm,
reasonably left containers unrealised.

## Decision

Two changes.

**1. Dead-ends become coverage-aware.** An approved item with no `Needs` is a dead-end
**only if it also has no coverers**. If something `Covers` it, it is a top-of-axis node
covered from below — not a dead-end. Its coverage is then judged from its coverers:
**shallow** if it has any coverer, **deep** if a coverer is itself deep (the per-`Need`
deep check, applied upward). A genuinely terminal approved item — **no `Needs` and no
coverers** — is still a dead-end. This **supersedes ADR 004's unconditional
`no Needs → dead-end` sub-decision** (ADR 004 otherwise stands).

**2. The architecture axis is first-class, and taught.** Both axes meet at the
component: a component covers its `req` (requirements axis) **and** its `container`
(architecture axis). The `onboard` skill and the register-authoring reference are
updated to teach this, and an uncovered `container` gets a diagnostic hint pointing at
the missing `component → container` coverage. This is not new behaviour — Plumbline's
own register wires it — only newly documented.

**Rejected: a bespoke structural-node kind/status** (gate `context`/`container` by
aggregation, outside the ladder). It cuts against "adopt prior art, no bespoke
variation" (ADR 002), and it does **not** remove the need to express
`component → container` membership — so change 1 plus documentation covers the same
ground with less surface area.

## Consequences

**Positive**

- A node covered from below no longer needs a ceremonial `Needs` edge to escape a
  spurious dead-end — the model reads honestly.
- The architecture axis becomes usable by others, not merely legible in Plumbline's own
  register.
- Driven by verified consumer feedback rather than a synthetic case.

**Negative / accepted**

- A real (small) change to the coverage core: a no-`Needs`-with-coverers node is scored
  through its coverers, and the deep check recurses upward for it (not vacuously deep).
  Covered by new tests.
- Expressing `component → container` membership is still required to trace the
  architecture axis — the change makes top nodes lighter, not the axis optional.
- Plumbline's own register is unaffected (its `context` declares `Needs: container`), so
  self-trace stays 100%; a new fixture exercises the no-`Needs`-covered path.

## Alternatives considered

- **Keep the unconditional rule; just require `Needs` on top nodes.** Rejected: it works
  (our register does it), but it leaves a plainly-covered node reading as "links to
  nothing," and the feedback shows that confuses sophisticated users. Coverage-awareness
  is simply more correct.
- **A structural-node status (the reporter's "direction 2").** Rejected — see above.
- **Documentation only, no engine change.** Rejected: it fixes the container case but
  leaves the `context` dead-end wart; the fuller fix was chosen deliberately.

## Provenance

- [ADR 004](004-requirement-status-lifecycle.md) — the dead-end and status model; the
  unconditional sub-decision superseded here; and #1 (dead-end detection).
- [ADR 002](002-traceability-only-emergent-analysis-out-of-scope.md) — traceability only,
  no bespoke variation (why "direction 2" is out).
- Verified consumer feedback (nhc dev AI, 2026-08-29), both cases reproduced against
  v0.2.1 in this session: a `container` with `component → req`-only coverage → uncovered;
  a `context` with coverers but no `Needs` → dead-end (report.go's rule checks
  `len(Needs) == 0` and never consults `coverers`).
- Ticket #29 (implementation).
