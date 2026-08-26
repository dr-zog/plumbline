# 002. Traceability only — emergent-structure analysis stays out of scope

- **Status:** Accepted
- **Date:** 2026-08-26
- **Deciders:** Dr Zog (operator), analysis by Claude

## Context

A field report from a team that adopted Plumbline through a God-Object dissolution found the
single highest-value practice was pairing Plumbline's **declared** structure — the
register — with an **emergent** structure built from the code by a *separate* tool
(Graphify: dependency graph, community detection, god-node identification). Two lenses,
one declarative and one empirical, converged on the same monolith. Crucially, this value
was realised with **no code-level integration**: the team ran the two tools
independently and a human synthesised the two outputs.

The report recommended making the pairing first-class — a "your register says N
components, the code clusters into M, here's the god node" report — turning Plumbline
from a traceability *gate* into an architecture-drift *detector*. That is a genuine,
tempting product expansion. This ADR records that we **decline** it.

## Decision

Plumbline stays a **language-agnostic traceability tool**. It does **not** integrate
with, consume the output of, or rebuild emergent-structure analysis. Where a team wants
the emergent view Plumbline deliberately doesn't provide — *is the declared structure
actually how the code clusters?* — Plumbline **recommends pairing it with a complementary
deterministic code-analysis tool** (such as Graphify), run as an independent lens and
synthesised by a human. That recommendation lives as a note in the `audit` skill and the
docs. There is no code-level coupling.

## Consequences

**Positive**
- **Language-agnosticism preserved.** Plumbline scans comment tags, so it works on any
  language — Go, Svelte, Markdown, anything. Consuming an emergent-analysis tool would
  subordinate that coverage to *that tool's* language matrix (Graphify has no Svelte).
  We keep the broad applicability that is a core property.
- **No coupling** to another product's schema, roadmap, or limitations.
- **The corroboration stays strong.** Two independent lenses agreeing is a strong signal
  *because* they share no machinery; keeping them independent preserves that signal.
- Scope stays sharp — the same "integrate/recommend, don't reinvent" spine as ADR 001
  (and the shelved docsite).

**Negative / accepted**
- No automatic declared-vs-emergent report; the cross-check is a human reading two tools.
  Judged marginal — the value is in *judging* whether a mismatch matters (human work),
  not merely *spotting* it.
- A team that wants that report badly builds it themselves, outside Plumbline.

## Alternatives considered

- **Integrate (consume an emergent graph).** Rejected: it subordinates Plumbline's
  language-agnostic coverage to the graph tool's language support (a regression of a core
  property), couples us to a third-party schema and roadmap, and dilutes the independence
  that makes the two-lens signal strong. The field report's value was realised *without*
  any such integration.
- **Grow emergent analysis in-house.** Rejected: a large mission expansion (dependency
  graphs, clustering, god-node detection) that duplicates a purpose-built tool and breaks
  the project's "integrate/recommend, don't reinvent" discipline.

## Provenance

- A field-adoption report (a God-Object dissolution on a real codebase) — the pairing's value, and
  recommendation #1 (the "bigger product" this ADR declines).
- The design's stated stance: Graphify is "an integration, not a dependency, and not the
  source of truth."
- Session discussion, 2026-08-26 — the language-agnosticism regression, schema/roadmap
  coupling, and independence-as-signal arguments.
