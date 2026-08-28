# Components — plugin container

The C4 components of the Claude Code plugin: one per skill.

### Audit skill
`component~audit-skill~1`
Status: approved

Runs the engine and narrates a prioritised scorecard.

Covers: container~plugin~1, req~audit-narration~1
Needs: impl

### Onboard skill
`component~onboard-skill~1`
Status: approved

Guided, level-by-level C4 build plus initial anchors.

Covers: container~plugin~1, req~onboarding~1
Needs: impl

### Maintain skill
`component~maintain-skill~1`
Status: approved

Rides the staged diff to keep anchors and the register current.

Covers: container~plugin~1, req~maintenance~1
Needs: impl

### Enforce skill
`component~enforce-skill~1`
Status: approved

Installs the developer's chosen pre-commit and/or CI gate.

Covers: container~plugin~1, req~enforcement~1
Needs: impl

### Showcase skill
`component~showcase-skill~1`
Status: approved

Renders the register into glossy, human-facing HTML collateral (two-pager, whitepaper).

Covers: container~plugin~1, req~showcasing~1
Needs: impl
