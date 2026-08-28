# Components — engine container

The C4 components of the Go engine, split into their own file to show a level
divided per container. Their `Covers` links point up into `../containers.md` and
`../requirements.md`; Plumbline aggregates the whole directory before resolving, so
cross-file references just work.

### Anchor scanner
`component~anchor-scanner~1`

Walks a tree and extracts anchors from comments (`internal/anchor`).

Covers: container~engine~1, req~anchor-scanning~1
Needs: impl, utest

### Register parser
`component~register-parser~1`

Parses the OFT-native Markdown register, one file or a tree (`internal/register`).

Covers: container~engine~1, req~register-parsing~1, req~status-lifecycle~1
Needs: impl, utest

### C4 validator
`component~c4-validator~1`

The locked type ladder and structural validation (`internal/c4`).

Covers: container~engine~1, req~c4-structural-validation~1, req~status-lifecycle~1
Needs: impl, utest

### Report builder
`component~report-builder~1`

Cross-checks anchors against the register, resolves deep coverage, and builds the
bidirectional report and scorecard (`internal/report`).

Covers: container~engine~1, req~broken-anchor-detection~1, req~uncovered-detection~1, req~orphan-detection~1, req~deep-coverage~1, req~coverage-scoring~1, req~threshold-gating~1, req~status-lifecycle~1, req~dead-end-detection~1
Needs: impl, utest

### Config loader
`component~config-loader~1`

Loads the optional JSON config and merges it over locked defaults (`internal/config`).

Covers: container~engine~1, req~config-override~1
Needs: impl, utest

### Command-line interface
`component~cli~1`

Wires config, scan, register and report, returns the gate exit code, and reports the
build version (`cmd/plumbline`).

Covers: container~engine~1, req~cli-gate~1, req~versioned-build~1
Needs: impl

### Build tooling
`component~build-tooling~1`

Cross-compiles the engine for every target and fuses the macOS universal binary
(`Makefile`, `tools/makefat`).

Covers: container~engine~1, req~cross-platform-binaries~1
Needs: impl
