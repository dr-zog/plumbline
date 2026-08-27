# Plumbline — design

**Light-touch, enforced traceability from requirements to code.**

Write your prose comments however you like — rich or sparse, Plumbline never measures
them. Its single non-negotiable is the **anchor**: a lightweight, machine-checkable marker
that links each *code-area* up to a documented requirement, feature, or component. A
zero-dependency Go binary proves the anchors resolve, in both directions, and fails the
build when they don't. The skills place the anchors for you and keep the register current.

This is the **design rationale** — the *why* behind Plumbline's shape, captured so the
decisions don't have to be re-derived. It is not a how-to: the processes for building and
using Plumbline live in [`docs/`](docs/).

---

## The problem

Codebases built feature-first lose the thread between the outcome someone wanted (the
*requirement*), the feature that meets it, and the code that implements it. When approach
changes, features get left behind silently — nothing records that feature F existed to
satisfy requirement R, so nobody notices when F is replaced whether R is still met. And you
cannot answer "what are all the features of this codebase?" without reading the whole tree.

Descriptive documentation that restates *what the code does* is not the fix — it rots
against the source. The durable, low-rot link sits one altitude up: **requirement ↔ feature
↔ code-area**, kept current by process rather than by hand-maintained prose.

## Principles — adopted prior art, no bespoke variations

The design is an assembly of three established disciplines, followed as-is:

- **C4 model** — Simon Brown. Context → Container → Component → Code: four zoom levels of
  architecture. This is the layered register of features and architecture. Followed as
  specified, including that the Container level names chosen technology.
  <https://c4model.com>
- **Requirements traceability (RTM)** — the bidirectional audit: does every requirement have
  covering code, and does every code-area trace to a requirement. A standardised discipline
  (DO-178C / ISO 26262 / IEC 62304). Tooled, docs-as-code style, by OpenFastTrace, which
  reports both uncovered requirements *and* orphan/obsolete code.
  <https://github.com/itsallcode/openfasttrace> · <https://open-needs.org/>
- **Living Documentation** — Cyrille Martraire. Store documentation on the documented thing
  itself (annotations on code elements); generate the register from those markers; do not
  over-document. Evergreen, low-maintenance.
  <https://www.oreilly.com/library/view/living-documentation-continuous/9780134689418/>

## The contract

- **No opinion on prose comments.** As rich or as sparse as the developer wants. Plumbline
  never mandates or measures descriptive "what this does" comments — that is the code-detail
  documentation that rots.
- **One non-negotiable: the anchor.** Adopt Plumbline and you carry anchors.
- An **anchor** is a lightweight, OFT-compatible marker in a comment on a code-area, linking
  it up to a documented thing (requirement / feature / C4 component).
- **Anchor unit: the code-area** — module / package / directory / marked region. Coarse, one
  altitude up. Not per-symbol. Coarse granularity is what keeps it light-touch.

## Architecture — two layers

**Engine — a Go static binary. The law.**
Deterministic, cheap, zero runtime dependencies (this is why Go, not a JVM/TS/Python tool
with a dependency tree — the binary ships inside the plugin). It:
- scans a tree for anchors and parses the register;
- emits a machine-readable **bidirectional report** — *uncovered requirements* (no anchoring
  code), *broken anchors* (target doesn't resolve), *unanchored code-areas* (orphans);
- exits non-zero on any gap, so it doubles as a pre-commit / CI gate with no LLM in the loop.

**Skills — a Claude Code plugin. The taste.**
The authoring judgement a binary can't have. Each consumes the engine's JSON:
- **`onboard`** — zero → structured. Assess, generate the C4 register skeleton, and place the
  initial anchors across existing code-areas.
- **`maintain`** — on a diff: add/update anchors for new or changed code-areas, keep the
  register current, clear what the binary flagged. Rides normal change flow.
- **`audit`** — run the binary, narrate the scorecard, prioritise the weakest areas.
- **`enforce`** — install the pre-commit / CI hook.

## Enforcement and scoring

The insistence and the light touch are two mechanisms working together, not a contradiction:
the **skill adds** the anchors (so a developer rarely hand-writes one), and the **binary
gates** on their presence (so unanchored code cannot land quietly).

One report yields both a gate and a score: the **exit code** is the gate (non-zero on any
gap); the **numbers behind it** (traceability coverage %, uncovered/orphan counts) are the
data-driven scorecard the skills narrate. Threshold gating ("fail below N%") is a later
option layered on the same run.

## Anchor format

OFT-compatible tag grammar — inherit the ecosystem rather than invent a syntax. The
light-touch feel comes from three things, none of which is a simpler grammar: coarse
granularity, an ID-only payload, and the skill authoring the tags so developers rarely type
one by hand. Implement against OpenFastTrace's tag specification.

## Repo shape

Subdir plugin, matching the marketplace's `git-subdir` source pattern:

- **repo root** — Go source, build tooling, CI.
- **`plugins/plumbline/`** — the shippable plugin: `skills/`, `bin/` (the per-platform
  binaries — the zero-dependency promise), `.claude-plugin/plugin.json`, `docs/`.

MIT licence. Public-facing files carry the Dr. Zog persona (British English, confident,
mechanism-first).

## Validation

Plumbline was validated on a real, actively-changing production codebase: a team adopted it,
ran `onboard`, and used it through a large refactor. Their field report confirmed the register
held — no requirement silently dropped across the move — and its recommendations are captured
in `docs/adrs/`.
