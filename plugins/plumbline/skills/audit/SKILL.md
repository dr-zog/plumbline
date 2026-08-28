---
name: audit
description: Run the Plumbline engine over a repository and narrate the traceability scorecard — coverage, uncovered requirements, dead-ends, broken anchors and orphan code-areas — prioritising the weakest areas. Use when someone asks to audit, score, or check requirement↔code traceability, or how well a codebase is anchored.
---

<!-- [impl->component~audit-skill~1] -->

# Plumbline — audit

Run the engine, read its JSON, and turn the numbers into a prioritised scorecard.
This skill only ever *reports*. Placing anchors and editing the register belongs
to `onboard` and `maintain`.

## Steps

1. **Locate the engine.** Run it via the plugin's launcher,
   `${CLAUDE_PLUGIN_ROOT}/bin/plumbline` (it picks the right per-platform binary).
   Working from the engine's source tree instead? `go build -o /tmp/plumbline
   ./cmd/plumbline`. If neither is possible, tell the user the engine isn't available
   and stop.

2. **Find the register.** Default `register.md` at the repo root; pass
   `-register <path>` if it lives elsewhere. If there is no register at all, the
   repo hasn't been onboarded — say so and point at `onboard` (a later skill).

3. **Run for JSON:**
   ```
   plumbline -json -register <register> <paths...>
   ```
   Exit code `0` = clean, `1` = gaps found, `2` = error (bad register path,
   unreadable tree). Read stdout as JSON regardless of exit code `0`/`1`.

4. **Narrate the scorecard.** Lead with the `summary`: the two scores —
   `coveragePct` (direct needs met) and `completenessPct` (the whole chain
   resolves) — and the gap counts. Note the `gate`: `minCoverage` of 0 is strict
   (any gap fails); a positive value is a threshold floor. Then walk the gaps in
   priority order:
   - **Structural errors** (`structural`, severity `error`) — first. The register
     itself breaks the C4 ladder: an unknown type, a disallowed `Needs` edge, or
     a `Covers` link whose target is missing or the wrong type. The model is
     malformed; fix it before trusting the rest. (Severity `warning`, e.g. a
     container that names no technology, is advisory and doesn't fail.)
   - **Broken anchors** (`broken`) — code points at a register ID that doesn't
     exist: a typo, or the item was renamed/removed. Outright wrong; cheap to fix.
   - **Uncovered** (`uncovered`) — a documented item whose needed coverage
     nothing provides. The `missing` array lists which artifact types have no
     anchor and no covering register item.
   - **Dead-ends** (`deadEnds`, count `deadEndCount`) — an *approved* item that declares
     no `Needs` at all. In the locked ladder every register-item type must need something
     below it, so this can never be covered; it's a **hard fail** (like broken/structural,
     it fails regardless of any coverage threshold), and the engine excludes it from the
     covered count so it can't fake a green score. Fix: give it a `Needs` edge down the
     ladder, or set it to a lower status (`proposed`/`draft`) if it isn't a committed
     requirement yet.
   - **Transitive defects** (`transitiveDefects`) — an item whose *direct* needs
     are met, but a coverer further down the chain isn't itself covered
     (`weakCoverers` names them). Coverage looks fine one level up; the hole is
     deeper.
   - **Orphan code-areas** (`orphans`) — source files carrying no anchor.
     Coarsely, code whose purpose isn't traced. Lowest priority, highest volume.

5. **Prioritise the weakest areas.** Group orphans by directory to spot whole
   unanchored subsystems; follow each transitive defect down to the uncovered
   leaf that causes it (fixing the leaf clears the chain). Recommend the smallest
   next step, and hand off to `maintain`/`onboard` for the actual edits.

6. **Set expectations honestly.** A green score means the register is *honest* —
   every requirement backed by code, every code-area traced — not that the design
   is *good*. A single sprawling "god component" can be fully anchored and fully
   green. For that emergent view — is the declared structure actually how the code
   clusters? — Plumbline doesn't and won't judge it (see ADR 002); recommend pairing
   it with a complementary **deterministic code-analysis tool** (e.g. Graphify), run
   as an independent lens. Two lenses agreeing — declared and empirical — is a far
   stronger signal than either alone.
