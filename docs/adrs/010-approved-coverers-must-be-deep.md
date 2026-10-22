# 010. The gate's deep coverage requires every approved coverer to be deep

- **Status:** Accepted
- **Date:** 2026-10-22
- **Deciders:** Dr. Zog (operator), analysis by Claude

## Context

ADR 009 split completeness into two notions: the **gate** keeps by-type `deep` (one deep
coverer per needed type suffices), and the **advice layer** uses `built` (every coverer,
recursively, built). It left one known overstatement on the gate side, split out as #48:
an **approved** item covered by one built approved requirement and one **approved but
uncovered** requirement still counts as *deep*.

Reproduced against v0.3.0:

```
feat~f~1     approved  Needs: req
req~built~1  approved  Covers: feat~f~1   (component anchored → deep)
req~gap~1    approved  Covers: feat~f~1   (nothing beneath it)
```

The report lists `req~gap~1` as **uncovered**, yet counts the feature as deep: it reads
"whole chain resolves", inflates `completenessPct` (75% here; 50% is honest), and is absent
from `transitiveDefects`. The report contradicts itself, and the transitive list fails to
point at the real weak link.

The verdict was never wrong — the uncovered sibling fails strict gating on its own — so
this is a matter of the scorecard telling the truth, not of the gate letting anything through.

## Decision

**The gate's `deep` additionally requires every *approved* coverer of an item to be deep.**
An item is deep iff each needed type is met by an anchor or by a deep coverer of that type
(unchanged), **and** every approved item that `Covers` it is itself deep. Planned
(`proposed`/`draft`) coverers are ignored by the new clause; they may still satisfy a need
if they happen to be deep. The clause applies on both branches of `deep`, including the
ADR 007 no-`Needs` path.

The reasoning: an **approved** coverer is a commitment. If a committed part of an item is
unbuilt, the item's chain does not fully resolve, whatever else covers it. A *planned*
coverer is not yet a commitment, so it must never fail an approved parent (ADR 004; ADR
009 finding 2) — which is exactly why the clause is restricted to approved coverers rather
than adopting `built` for the gate.

This **refines ADR 009's gate notion** ("the gate keeps by-type deep, unchanged"): still
by-type over all coverers for meeting needs, now also all-approved-coverers. Plumbline keeps
two notions, not three. ADR 009 otherwise stands.

## Consequences

**Positive**

- The scorecard stops contradicting itself: an item over an uncovered approved coverer is
  reported as a **transitive defect** naming that coverer, and `completenessPct` drops to
  the honest figure.
- **No verdict can change.** If an item loses deep under the clause, some approved coverer
  beneath it is not deep, and that coverer already fails strict gating itself (uncovered,
  transitive or dead-end). Threshold gating reads the shallow `coveragePct`, which this
  never touches.
- **Repos that pass strict gating see identical numbers.** With no failing approved item,
  every approved coverer is deep, so the clause is always satisfied.
- The advice layer (`built`, status-lag, `realised`, spec-debt, warnings) is untouched.

**Negative / accepted**

- A repo that is already failing — or passing only by threshold — will report a lower
  `completenessPct` and `deepCovered`, and more `transitiveDefects`. That is the point,
  but consumers parsing those numbers will see them move. It lands as a **`fix:`** (patch):
  it corrects overstated figures, changes no verdict, and changes nothing for a strict-green repo.
- One more condition in `deep`. Small, and it reuses the existing memoised recursion.

## Alternatives considered

- **Document the by-type meaning; change nothing.** Rejected: the "whole chain resolves"
  label is simply false for such an item, and the report contradicts its own uncovered list.
- **Use `built` (every coverer, planned included) for the gate.** Rejected: an approved item
  with unbuilt *planned* coverers would become a transitive defect and turn `main` red over
  un-built spec — forbidden by ADR 004, and the reason ADR 009 kept the gate by-type.
- **Keep the numbers; add a narrative field naming the uncovered approved siblings.**
  Rejected: more report surface, and the headline percentage stays wrong.
- **A third notion used only for `completenessPct`.** Rejected: the transitive list would
  then disagree with the score, and it adds a notion where a refinement suffices.

## Provenance

- Issue #48, split from #44 so that PR #47 could promise no gate-facing change; ADR 009's
  Consequences named it as a known, deferred overstatement.
- `internal/report/report.go` — the `deep` closure (both branches) and its consumers:
  `deepCount`/`CompletenessPct`, `TransitiveGap` and `weakCoverers`.
- **Prototype and differential test (2026-10-20):** the clause implemented in a throwaway
  worktree passed every existing test unchanged (including the end-to-end fixture) and kept
  Plumbline's self-trace at 100%. The #48 repro went from completeness 75% / no transitive
  defects to 50% / `feat~f~1` transitive with weak coverer `req~gap~1`. A differential fuzz
  over 3,000 random, structurally valid C4 registers (mixed statuses, several coverers per
  item, partial anchoring), each run through both engines in strict and threshold mode,
  found **0 verdict flips**, **0 differences** in coverage %, spec-debt, warnings or the
  planned list, and numbers changing in 631 cases — **none** of them strict-green (311
  strict-green cases, all unchanged).
- ADR 004 (approved must be fully covered; planned never gates); ADR 007 (no-`Needs`
  coverage, reused); ADR 009 (the two notions — gate notion refined here).
