# 008. Orphan-detection granularity is the directory

- **Status:** Accepted
- **Date:** 2026-08-29
- **Deciders:** Dr. Zog (operator), analysis by Claude

## Context

Coverage in Plumbline is often described as resolving "at the directory level," but that
conflates two different mechanisms at two different granularities:

- **Forward — "is this register item covered?"** is *point-level*. An anchor is a single
  comment tag on a line; an item is covered when *one* anchor anywhere in the tree targets
  it (`needMetShallow`). Directories never enter into it — the same tree flattened into one
  file resolves identically.
- **Reverse — "is this code untraced?"** (orphan detection) is the *only* place a
  granularity choice exists, and it is the **directory**: source files are bucketed by
  `filepath.Dir`, and a directory is an orphan iff no file in it carries any anchor. One
  anchor covers the whole area.

So the only decision to make is the unit of the *reverse* audit — and it had never been
recorded. It lived as one DESIGN.md line ("the code-area … coarse … not per-symbol") and a
code comment. Unlike the OFT boundary (ADR 003), the unit of the untraced-code audit had no
ADR, so it read as inherited prior art when it isn't: OFT gives us the point-level tag
grammar and the forward coverage/needs model; the reverse "untraced code-area" audit **and
its directory unit are Plumbline's own**. OFT's long experience therefore doesn't settle it.

The question that surfaced the gap: should the reverse audit resolve at *file* (or finer)
granularity — to catch feature loss/removal more tightly — and does directory-level force
languages that dislike deep trees into artificial sub-directorying?

## Decision

**The orphan / unanchored-code audit resolves at the directory, and stays there.** A
directory containing scanned source is traced iff at least one file in it carries at least
one anchor; a single anchor covers the whole area. Forward coverage remains point-level
(unchanged). File- and symbol-level orphan units are **declined for the gate**.

The governing principle is the two-layer design: the **binary is the coarse mechanical
floor; per-file / per-symbol judgement belongs to the skills.** The directory is the coarsest
unit that still localises an orphan to a meaningful place, and coarseness is the light-touch
thesis. The finer question — "these new files deserve their own component/anchor" — rides
the diff in the `maintain` skill, where judgement lives, not in the binary.

## Consequences

**Positive**

- **Light-touch preserved.** One anchor per package, not per file; adding a file to an
  already-anchored package demands no new tag; no tag churn on file split or rename.
- **It does not force sub-directorying — the pressure is the opposite.** A flat layout is the
  *most* forgiving case: one directory, one bucket, one anchor. *Finer* granularity is what
  would punish languages that keep large flat directories, so point A of the concern argues
  for directory-level, not against it.
- **Feature-removal detection is unaffected**, because it is a *forward* property: delete a
  feature's code (anchor included) and the requirement it satisfied goes uncovered → red,
  resolved at anchor granularity regardless of the orphan unit. File-level would add nothing
  here.
- **Clean division of labour.** The binary gates coarsely and cheaply; the `maintain` skill
  applies file/symbol-level attention on the diff, so consumers get fine-grained review
  without the gate carrying fine-grained state.

**Negative / accepted**

- **One blind spot:** a genuinely new, untraced feature added *inside* an already-anchored
  directory is not flagged by the gate — the directory already has an anchor. Mitigated by
  (a) new features usually landing as new files/directories the audit *does* see, and (b) the
  `maintain` skill's diff-level judgement. Accepted as the price of coarseness.
- Directory-level cannot distinguish two features sharing one directory; if one is silently
  removed and its loss breaks no requirement forward, the gate is silent. This is the same
  class as symbol-level loss, which DESIGN explicitly declines. Accepted.
- A team that habitually piles unrelated features into single directories is served too
  coarsely — addressed **not** by changing the default but, if a concrete case appears, by a
  future opt-in config knob (see Alternatives).

## Alternatives considered

- **File-level orphan unit.** Rejected: it is not a principled stopping point (a two-feature
  file that loses one is invisible at file granularity too — only symbol-level is principled,
  and DESIGN declines that); its sole gain over directory is catching untraced *additions*
  into existing packages; and it costs the light-touch thesis directly — a tag per file, a
  red build on every new file until tagged, churn on rename/split. The gain is better
  delivered by the `maintain` skill without moving judgement into the binary.
- **Symbol-level orphan unit.** Rejected: DESIGN's explicit "not per-symbol"; the heaviest
  possible tag load; it turns anchors into the rotting per-symbol documentation Plumbline
  exists to avoid.
- **Make granularity configurable now.** Rejected as premature — the same discipline as ADR
  003's deferral of `Tags`: add the knob when a concrete case demands it, not speculatively.
  Recorded here as the sanctioned escape hatch should such a case appear: an opt-in finer
  unit defaulting to directory, in the same spirit as `minCoverage` and `exclude`.
- **Leave it undocumented (status quo).** Rejected: it was the least-recorded of the boundary
  decisions, it read as inherited-from-OFT when it isn't, and it would be re-litigated. Hence
  this ADR.

## Provenance

- **Code (2026-08-29):** `internal/report/report.go` — orphan detection buckets scanned files
  by `filepath.Dir` and flags a directory with no anchored file (the `Orphan` type's doc
  already names the unit "package / directory"); `needMetShallow` / `needMetDeep` — forward
  coverage is point-level, one matching anchor suffices. `internal/anchor/scan.go` — an anchor
  is a single tag at a line; `Scan` returns anchors plus the scanned-file list, the latter
  used *only* for orphan detection.
- **DESIGN.md:** "Anchor unit: the code-area — module / package / directory / marked region.
  Coarse, one altitude up. Not per-symbol. Coarse granularity is what keeps it light-touch,"
  and the two-layer architecture (binary = mechanical floor, skills = judgement).
- **ADR 003 (OFT conformance boundary):** we adopt OFT's tag grammar and forward
  coverage/needs model; the reverse untraced-code audit and its directory unit are
  Plumbline's own, so OFT does not settle this decision.
- **Session discussion, 2026-08-29:** the directional-granularity clarification (forward =
  point-level, reverse = directory-level); that feature-loss is a forward property independent
  of the orphan unit; that the sub-directorying worry inverts (finer granularity is what
  pressures flat layouts); and that file-level is an unprincipled midpoint between the
  point-level forward gate and the declined symbol level.
- **Related:** ADR 003 (what we adopt from / decline of OFT); ADR 007 (the coverage core this
  audit sits beside).
