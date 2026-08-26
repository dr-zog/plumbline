# Features

### Bidirectional audit
`feat~bidirectional-audit~1`

Report both directions at once: requirements with no covering code, and code with no
requirement.

Needs: req

### Data-driven scorecard
`feat~data-driven-scorecard~1`

One run yields both a gate (the exit code) and a score (coverage and completeness),
tunable to a threshold.

Needs: req

### C4-structured register
`feat~c4-structure~1`

The register is a layered C4 model, validated against a locked artifact-type ladder,
and may be split across many files.

Needs: req

### Skills workflow
`feat~skills-workflow~1`

The plugin's skills onboard a codebase, keep it current, enforce the gate, and narrate
the score — the authoring judgement a binary can't have.

Needs: req

### Audience-facing docs
`feat~audience-facing-docs~1`

Generate glossy, human-facing product documentation from the register, so a product's
marketing can't silently drift from what it actually does.

Needs: req

### Cross-platform distribution
`feat~cross-platform-distribution~1`

Plumbline ships as a self-contained Claude Code plugin across Linux, macOS and Windows —
a single static engine binary per target, bundled in the plugin, no toolchain to install.

Needs: req
