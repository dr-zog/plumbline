# Requirements

### Scan source for anchors
`req~anchor-scanning~1`

Find every anchor across a source tree, across many languages, honouring off/on fences.

Covers: feat~bidirectional-audit~1
Needs: component

### Parse the register
`req~register-parsing~1`

Parse the OFT-native Markdown register — one file or a tree of them — into specification
items with their Needs and Covers.

Covers: feat~bidirectional-audit~1
Needs: component

### Detect broken anchors
`req~broken-anchor-detection~1`

Report any anchor whose target ID is not in the register.

Covers: feat~bidirectional-audit~1
Needs: component

### Detect uncovered requirements
`req~uncovered-detection~1`

Report any register item whose needed coverage nothing provides.

Covers: feat~bidirectional-audit~1
Needs: component

### Detect orphan code-areas
`req~orphan-detection~1`

Report any scanned source file carrying no anchor.

Covers: feat~bidirectional-audit~1
Needs: component

### Resolve deep coverage
`req~deep-coverage~1`

Resolve coverage across the whole chain and surface transitive defects, not just direct anchors.

Covers: feat~bidirectional-audit~1
Needs: component

### Score coverage and completeness
`req~coverage-scoring~1`

Compute shallow coverage % and deep completeness % over the register.

Covers: feat~data-driven-scorecard~1
Needs: component

### Gate on a threshold
`req~threshold-gating~1`

Gate strictly (any gap fails) or against a coverage floor.

Covers: feat~data-driven-scorecard~1
Needs: component

### Gate via exit code
`req~cli-gate~1`

Exit non-zero on any gap so the engine doubles as a pre-commit / CI gate with no LLM in the loop.

Covers: feat~data-driven-scorecard~1
Needs: component

### Override defaults via config
`req~config-override~1`

Tune register path(s), roots, exclusions and threshold via an optional JSON config over locked defaults.

Covers: feat~data-driven-scorecard~1
Needs: component

### Validate C4 structure
`req~c4-structural-validation~1`

Validate the register against the locked ladder — unknown types, disallowed edges, dangling covers, duplicate IDs.

Covers: feat~c4-structure~1
Needs: component

### Onboard a codebase
`req~onboarding~1`

Take a project from zero to a standing C4 register with initial anchors.

Covers: feat~skills-workflow~1
Needs: component

### Maintain traceability
`req~maintenance~1`

Keep anchors and the register current on the staged diff.

Covers: feat~skills-workflow~1
Needs: component

### Enforce the gate
`req~enforcement~1`

Install the pre-commit and/or CI gate the developer chooses.

Covers: feat~skills-workflow~1
Needs: component

### Narrate the scorecard
`req~audit-narration~1`

Run the engine and narrate a prioritised scorecard.

Covers: feat~skills-workflow~1
Needs: component

### Render audience-facing docs
`req~showcasing~1`

Render the register into branded, benefit-led HTML — a two-pager or a whitepaper —
claiming only what the register supports.

Covers: feat~audience-facing-docs~1
Needs: component

### Stamp the engine with its version
`req~versioned-build~1`

The engine reports its build version, stamped in at link time.

Covers: feat~cross-platform-distribution~1
Needs: component

### Build a binary for every target
`req~cross-platform-binaries~1`

Cross-compile a single static engine binary for each target — Linux amd64/arm64,
Windows amd64, and a macOS universal (fat) binary fused from amd64 + arm64 on the build
host, with no macOS host or `lipo`.

Covers: feat~cross-platform-distribution~1
Needs: component

### Honour requirement status
`req~status-lifecycle~1`

Honour OFT's `Status` attribute: `approved` (the default) gates as normal; `proposed`/`draft`
are tracked but not gated; `rejected` is excluded from tracing. Coverage and completeness are
computed over the gated set, with a planned count surfaced; an unknown status is a structural
error.

Covers: feat~requirement-lifecycle~1
Needs: component

### Detect dead-end requirements
`req~dead-end-detection~1`

Flag an approved item that declares no `Needs` — a dead-end in the locked ladder, which can
never be genuinely covered. It fails the gate rather than reading as vacuously complete.

Covers: feat~bidirectional-audit~1
Needs: component
