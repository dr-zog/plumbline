# 009. Two completeness notions: a by-type gate, an all-coverers advice layer

- **Status:** Accepted
- **Date:** 2026-10-04
- **Deciders:** Dr. Zog (operator), analysis by Claude (Fable); verified consumer field report (nhc dev AI)

## Context

A `proposed` feature with several requirements, of which only one is built, was
reported as **fully covered → "promote to approved"** (a `status-lag` warning) and
`realised: true` — overstating what is built. Reproduced against the live engine
(issue #44). The same defect appears one rung down (`req` with a built and an
unbuilt `component`) and on the ADR 007 no-`Needs` path, so it is a property of
**multiple same-type coverers at any rung**, not of features.

The root cause is **one notion serving four contracts**. `report.go`'s `deep(it)`
feeds: the gate (transitive defects / completeness score, for *approved* items);
spec-debt (`Planned && !isDeep`); `Planned.Realised`; and the status-lag/build-ahead
split. `deep` is satisfied **by needed type** — `needMetDeep` returns true on the
*first* coverer of a needed type that is itself deep — so additional unbuilt coverers
of that type are invisible. That is right for the gate and wrong for the three
advice consumers, which each claim an *all-coverers* meaning.

Two findings from the analysis settle the shape of the fix:

- **OpenFastTrace's deep coverage is per-*coverer*, not per-*type*.** OFT's
  `DeepCoverageResolver` folds *every* incoming coverage link with the worst child
  status; only *shallow* coverage is by type. So the all-coverers notion is **OFT's
  actual rule**, and Plumbline's by-type `deep` is the *bespoke variation* — the
  reverse of what ADR 004 implied ("OFT … conveying partial-vs-complete by needed
  type"). That sentence is accurate for shallow only.
- **The by-type variation is load-bearing for the gate.** An *approved* feature with
  one approved+built req and two `proposed`, unbuilt reqs passes strict gating today.
  Switch the gate to OFT's per-coverer rule and `main` goes red *because of planned
  requirements* — precisely what ADR 004 forbids. So the gate must keep by-type deep.

Issue #43 (warnings require real code — a `hasCode` anchor in the subtree) landed
first; this builds on it.

## Decision

**Plumbline carries two status-blind completeness notions, each answering a
different question:**

- **The gate keeps by-type `deep` coverage** — a *deliberate, now-documented
  divergence from OFT*, so that a planned coverer can never fail an approved parent.
  Unchanged. It answers *"is it safe to approve?"*
- **The advice layer adopts OFT's per-coverer rule verbatim** as a second predicate,
  `built(it)`: an item is *built* iff it is shallow-covered **and every** item that
  `Covers` it is, recursively, itself built (anchors are built leaves). A no-`Needs`
  top-of-axis node requires ≥1 coverer, all built — never vacuous (ADR 007). `built`
  is used in exactly three places: the status-lag vs build-ahead choice, the
  `Planned.Realised` flag, and the spec-debt count. It answers *"does maturity lag
  what is actually built?"*

The report additionally names the unbuilt coverers (`Planned.Unbuilt`), so a partial
feature reads as "2 of 3 built: …" rather than hiding the gap.

Because **`built ⇒ deep`** (a built coverer is a deep coverer), the change is
**monotone**: `realised` only flips true→false, `status-lag` only softens to
`build-ahead`, spec-debt only grows. No consumer ever gains a *new* green, and no
strict/threshold gate verdict can change (those read uncovered/transitive/orphan and
the shallow score, never the warnings or `Planned`).

This **refines ADR 004's "full coverage → promote"** (sharpening what "full" means for
the advice layer) and corrects its by-type/OFT claim. It supersedes nothing — the gate
model, and the rest of ADR 004, stand (the ADR 007 pattern: a sub-decision refined).

## Consequences

**Positive**

- The advice layer stops overstating: `realised`, the promote nudge and the burndown
  now mean *everything beneath is built*. The self-contradiction — a feature reading
  "promote" while the gate fails on its own uncovered approved sibling — is gone.
- The OFT divergence is stated openly and **shrunk to the one place it earns its keep**
  (the gate), as ADR 002/003 require. Both notions keep guaranteeable completeness over
  the same closed graph.
- Monotone and gate-safe: no new greens, no verdict flips; composes with #43.

**Negative / accepted**

- A second memoised predicate to maintain beside `deep`. Small, and a near-mirror of it.
- A configured spec-debt **budget** (`maxProposed`/`maxProposedPct`) can now count more
  un-built spec and therefore tighten — a real, if opt-in, behaviour change. Lands as a
  **`feat:`** (minor) for that reason; strict/threshold gating is untouched.
- The *scorecard* can still count an approved feature `deep` while an approved sibling
  req is uncovered (overstates the completeness %, never flips a verdict). Left to a
  separate ticket — out of scope here.

## Alternatives considered

- **Document only; reword "fully covered."** Rejected: `Realised`'s name and the
  status-lag comment already *define* the all-coverers meaning, and the approved-sibling
  self-contradiction cannot be worded away.
- **Expose a coverage fraction ("2 of 3") and promote at 1.** Rejected as the primary
  mechanism: a fraction is not meaningful through the ladder, and "coverage as a
  percentage" is what ADR 004 declined. Its surviving value — naming the unbuilt
  coverers — is folded in as `Planned.Unbuilt`.
- **Model multi-requirement features specially** (split, an "epic" status, a structural
  kind). Rejected: it unlocks the ladder (ADR 003) or crosses the maturity/progress axes
  (ADR 004), and the defect exists with no feature involved.
- **Switch the gate's `deep` to OFT's per-coverer rule too.** Rejected: it per-item-gates
  planned items through their approved parents, reddening `main` for un-built spec
  (finding 2 above; ADR 004 forbids it).

## Provenance

- Issue #44 and its seeded design proposal; issue #43 (the `hasCode` precondition);
  verified reproductions against the engine (scratchpad), including the approved-sibling
  self-contradiction case.
- **OFT source:** `DeepCoverageResolver` walks all incoming `COVERED_*` links and folds
  with the worst child status (per-coverer deep); shallow is by type
  (`areAllArtifactTypesCovered`). *To be re-confirmed from upstream and cited verbatim
  when this is implemented.*
- `internal/report/report.go` — the single `deep` closure (`needMetDeep`, the no-`Needs`
  branch) and its four consumers (`TransitiveGap`/`deepCount`, spec-debt, `Planned.Realised`,
  the status-lag/build-ahead split).
- ADR 004 (status lifecycle — refined here; its by-type/OFT claim corrected); ADR 007
  (never-vacuous coverage — the pattern reused for `built`); ADR 002/003 (adopt prior art,
  state divergences openly, guaranteeable completeness); ADR 003's OFT conformance boundary.
- Session discussion, 2026-10-04: one notion per question; `built ⇒ deep` monotonicity;
  the gate's divergence is load-bearing and must stay.
