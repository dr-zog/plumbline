# 003. OpenFastTrace conformance boundary

- **Status:** Accepted
- **Date:** 2026-08-27
- **Deciders:** Dr. Zog (operator), analysis by Claude

## Context

Plumbline bills its register as "OFT-native Markdown," but implements only a *subset* of
OpenFastTrace — and the boundary of that subset was never decided. It is partly principled
(the light-touch thesis: an ID-only anchor payload, coarse granularity) and partly
incidental (features simply not built yet).

A question — how to design a requirement up front and hold it *knowingly uncovered* for a
while — surfaced the gap concretely: the engine ignores OFT's standard `Status:` attribute
entirely. It is not merely unused; `parse.go` does not even recognise it, so a `Status:`
line placed under an ID is silently swallowed as the item's description. An audit found the
same pattern throughout: the parsers recognise more OFT than the resolver consumes —
`Depends`/`Tags` are matched then dropped (`parse.go`), and the anchor scanner captures
cover/target revisions and the `>>` forwarded-needs list that `report.go` never reads. So
"OFT-native" was a mild over-claim, and no ADR defined what Plumbline actually commits to.

This ADR sets that boundary, and gives it a spine.

## Decision

**The governing principle: Plumbline models only the relationships whose completeness it
can mechanically guarantee.**

The requirement→code thread grounds out at code anchors, which either exist or do not — so
"is every requirement satisfied down to code?" is a question the engine can *answer*,
exhaustively, in both directions (coverage, and its inverse: orphan detection). That
guarantee is Plumbline's core value. A relationship on which completeness *cannot* be
verified is declined, because admitting it would dilute the one promise that makes the gate
worth trusting.

Plumbline is therefore **OFT-native for the item, coverage, and lifecycle model, and
deliberately minimal on OFT's graph-flexibility features.** Per attribute:

- **Adopt.** The ID grammar, `Needs`/`Covers`, and deep/shallow bidirectional coverage
  (already implemented). **`Status`** is adopted as the requirement-lifecycle axis:
  `approved` (the default) gated as today; `proposed`/`draft` tracked but not gated
  (present and linked, but their unmet needs do not fail the build); `rejected` excluded
  from tracing, as OFT does. This finishes conformance on the lifecycle axis and adds a
  planned-vs-realised axis to the scorecard. The precise gate/scorecard semantics are an
  implementation concern for the feature that follows.
- **Decline, by design.**
  - **`Forwards`** (the `>>` forwarded-needs list) — *redundant*. Forwarding exists to
    propagate needs through a flexible type graph; Plumbline's locked, vertical C4 ladder
    plus its deep-coverage resolver already propagate needs down the chain automatically.
    It would return only if the fixed ladder were abandoned.
  - **`Depends`** — *unprovable*. A lateral "presupposes" link is a non-coverage edge the
    resolver never walks, so it is unrelated to cyclicity; the reason to decline it is that
    its completeness cannot be guaranteed. The engine can confirm a *declared* dependency's
    target exists, but can never detect a *missing* one — there is no "orphan dependency."
    Admitting it would place an unguaranteeable axis inside a tool whose value is the
    guarantee.
- **Defer.** **`Tags`** — harmless labels, but inert until a report consumes them, and
  `feat` already groups requirements for the epic case; adopt only when a concrete
  cross-cutting filter needs them. Revision-aware coverage (OFT's *outdated*/*predated*) is
  left open, to be decided on its own merits later.

The "OFT-native Markdown" claim in the README and docs is qualified to match: native for
the item/coverage/lifecycle model, built on OFT's grammar, deliberately minimal on graph
flexibility.

## Consequences

**Positive**

- The conformance claim becomes *honest and decided*, not accidental — with a durable test
  (the completeness principle) for every future OFT or feature question.
- `Status` unlocks the design-up-front workflow (mark a requirement `proposed`, hold it
  knowingly uncovered without a red `main`), the planned-vs-realised burndown, and —
  composed with dead-end detection (issue #1) — makes "an *approved* item that arms
  nothing" a precise defect rather than an ambiguous one.
- Two apparent "gaps" become design statements: Plumbline does *structurally* what OFT
  needs extra syntax for (`Forwards`), and declines what it cannot prove complete
  (`Depends`). Both strengthen the story rather than weaken it.

**Negative / accepted**

- Plumbline is explicitly *not* full OFT. A team needing `Depends`, `Forwards`, or OFT's
  import/export formats must pair tools or look elsewhere.
- `Status` is real engine work (recognise the attribute; gate and score by it), landed as
  its own feature; until then a `Status:` line stays inert, and fixing the parser so it is
  no longer swallowed into the description is part of that work.
- `Tags` and revision-awareness stay unadopted — a known and deliberate incompleteness.

## Alternatives considered

- **"OFT-compatible, minimal" — soften the claim, adopt nothing new.** Rejected: leaves the
  genuinely valuable `Status` on the table and under-sells the deep-coverage model, which
  is already *strong* OFT conformance.
- **"Full OFT-native, earned" — adopt Status, revisions, Depends, Forwards, export.**
  Rejected: `Depends` fails the completeness test and `Forwards` is redundant under the
  ladder; chasing full OFT abandons the light-touch thesis for a heavier
  requirements-management tool.
- **Unlock the C4 ladder (arbitrary user-defined types, as OFT allows).** Rejected: it
  forfeits zero-config onboarding, the universal vocabulary, built-in structural
  validation, and — crucially — *guaranteeable completeness*. That is a different product.
- **Leave the boundary undefined (status quo).** Rejected: it is what produced a false
  "OFT-native" claim and the silent `Status:` misparse.

## Provenance

- **Code audit (2026-08-27):** `internal/register/parse.go` (recognises
  `Needs`/`Covers`/`Depends`/`Tags`, stores only the first two; `Status` unrecognised → can
  land in `Desc`); `internal/report/report.go` (coverage, deep resolution, and the gate
  never read `Status`, revisions, or the `>>` list); `internal/c4/model.go` (the locked
  vertical ladder; acyclic by construction); `internal/anchor/scan.go` (captures
  cover/target revisions and `>>` forwarded needs, all unused downstream).
- `Status`, `Depends`, `Tags`, `Forwards` are standard OpenFastTrace specification-item
  attributes.
- The design's "OFT-native Markdown" claim (README, DESIGN.md).
- **Session reasoning, 2026-08-27:** the completeness principle ("model only what
  completeness can be guaranteed for"); that `Depends` completeness is unprovable (no
  "orphan dependency") and that `Depends` is unrelated to cyclicity because it is a
  non-coverage edge; that `Forwards` is redundant under the locked ladder plus deep
  resolution; that the ladder's vertical *shape* and architecture vocabulary come from C4,
  while the *locking and enforcement* — which delivers the guarantee — is Plumbline's own
  decision.
- **Related:** issue #1 (dead-end register items), reframed by `Status`; ADR 002 (a scope
  boundary on a different axis).
