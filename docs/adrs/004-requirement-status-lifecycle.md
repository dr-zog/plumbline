# 004. Requirement status lifecycle

- **Status:** Accepted
- **Date:** 2026-08-28
- **Deciders:** Dr. Zog (operator), analysis by Claude

## Context

ADR 003 adopted OFT's `Status`, and #6 implemented it with `absent = approved` (OFT's
default). Reviewing that default surfaced a conflation we had to untangle: **status
describes the maturity of the *requirement* (the spec); coverage describes whether *code*
fulfils it.** They are orthogonal axes. OFT keeps them separate — it computes coverage
identically regardless of status, conveying partial-vs-complete by needed *type*
(shallow/deep), and only `rejected` is filtered out of tracing. Every earlier tangle came
from trying to make one status field carry both the requirement's maturity *and* the code's
progress (the "todo/doing/done" reading), which crosses the two axes.

Separated, the design is a **policy over a maturity × coverage grid.**

## Decision

**Two orthogonal axes govern an item:**

- **Maturity** — the requirement's own lifecycle: `draft` (rough) → `proposed` (under
  review) → `approved` (committed); `rejected` (abandoned).
- **Coverage** — how much code fulfils it: none / partial / full (OFT's shallow/deep, per
  needed type).

**The gate keys on three maturity classes:**

- **`approved`** — *per-item gated*: it must be **fully covered**. `none` or `partial`
  coverage fails (an approved item that declares no `Needs` is a dead-end, #1). This is
  today's strict gate.
- **not-yet-approved (`draft` or `proposed`)** — *not per-item gated*. Counts toward the
  **aggregate spec-debt budget** (a bounded amount of un-approved spec is allowed). **Any
  code against it is a warning** — you are building ahead of the spec's approval; reconcile
  by approving it (if committed) or leaving it provisional. `draft` and `proposed` are the
  **same to the gate**; they differ only as maturity labels for the human and the burndown.
- **`rejected`** — excluded from tracing; code pointing at it is **zombie code**, which
  **fails** (remove the code, or un-reject the requirement).

**Default, and the divergence from OFT:** a requirement with **no status is `proposed`** —
provisional, because we do not assume a commitment that was never declared. This diverges
from OFT's `approved` default, deliberately and openly. With dead-end detection (#1) closing
the vacuous-green hole, defaulting to `proposed` is both safe and honest: a bare requirement
is provisional, bounded by the budget, and any code against it warns.

**Surfaced beyond pass/fail:** *status-lag* (not-yet-approved + full coverage → "promote to
`approved`") and *zombie code* (`rejected` + code).

This **supersedes ADR 003's `absent = approved` sub-decision** (the rest of ADR 003 stands)
and **replaces the original direction of #8** (flip the default / a synthetic "unclaimed"
state) with the maturity × coverage model.

## Consequences

**Positive**

- The two concerns are honestly separated — a requirement's maturity and its implementation
  are tracked independently, as OFT intends.
- Enforcement is layered and sensible: **per-item** for `approved`, an **aggregate budget**
  for un-approved spec, a **warning** wherever code runs ahead of approval, and a **hard
  fail** for zombie code. Sensible gated defaults; the AI carries the paperwork.
- Design-ahead is first-class (`draft`/`proposed`, bounded by the budget, never forgotten),
  and building ahead of approval is *allowed but surfaced*, not blocked — agile-friendly
  without going silent.
- `draft` vs `proposed` cost nothing at the gate; they are free maturity signal for the
  burndown.

**Negative / accepted**

- Breaking vs #6 (`empty → proposed`, not `approved`; the new warn/fail policy). Pre-1.0,
  this lands as a minor bump (`0.2.0`).
- Depends on dead-end detection (#1) to keep "full coverage" honest — without it,
  `approved + no Needs` fakes the ideal.
- The budget's exact thresholds and how they tighten over a project's life are a separate
  concern — the maturity ladder (#9).

## Alternatives considered

- **Keep `absent = approved`** (ADR 003 / #6). Rejected: assumes a commitment never
  declared, and (pre-#1) hid dead code as "the ideal."
- **Make status carry implementation-progress** (todo/doing/done). Rejected: it crosses the
  two orthogonal axes — status is the *requirement's* maturity, not the *code's* progress.
- **Fail (not warn) on code against an unapproved spec.** Rejected: too rigid — it forces
  approval-before-code (spec-first), blocking spikes and incremental build. A warning
  surfaces the mismatch without blocking.
- **A synthetic "unclaimed" third state.** Rejected: unnecessary once `no status = proposed`
  — the maturity axis already carries it.

## Provenance

- ADR 003 (OFT conformance) and #6 (initial `Status`, `empty → approved`).
- Session, 2026-08-28: the maturity-vs-coverage disentanglement (status = the *requirement's*
  maturity; coverage = the *code's* fulfilment; orthogonal, as OFT treats them — coverage is
  computed independently of status and conveys partial/complete by needed type, shallow/deep,
  not a percentage); the maturity × coverage gate grid and its four rulings (warn on
  build-ahead, fail on zombie, `no status = proposed`, `draft`/`proposed` collapse at the
  gate).
- Tickets: #1 (dead-end, the prerequisite), #8 (implementation), #9 (the maturity ladder —
  the budget's ratchet).
