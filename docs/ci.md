# Using Plumbline in CI

The engine is a single zero-dependency binary that exits non-zero on any traceability
gap — so a CI job that runs it *is* the gate. No LLM in the loop, no extra assertion
logic, no Go toolchain required. This guide covers getting the binary into a pipeline
and wiring it as a build gate.

It assumes the repo already has a Plumbline register (`register.md` or a `register/`
tree) and anchors in the code. To create those, use the `onboard` skill (or see the
[README](../README.md)); this guide is only about *running* the gate.

## Get the binary — pinned to a version (recommended)

Every [release](https://github.com/dr-zog/plumbline/releases) attaches the four platform
binaries and a `SHA256SUMS`, so CI can fetch an exact, reproducible build and verify it
before running ([ADR 006](adrs/006-release-binaries-for-ci.md)):

```bash
V=v0.2.1
base="https://github.com/dr-zog/plumbline/releases/download/$V"
curl -sSLO "$base/plumbline-linux-amd64"
curl -sSLO "$base/SHA256SUMS"
sha256sum --ignore-missing -c SHA256SUMS   # fails the job if the download is tampered
chmod +x plumbline-linux-amd64
```

The assets, one per platform: `plumbline-linux-amd64`, `plumbline-linux-arm64`,
`plumbline-darwin-universal`, `plumbline-windows-amd64.exe`. Pin `V` to a specific tag —
don't track a moving target — so a new Plumbline release can never silently change your
gate's behaviour.

## Run the gate

```bash
./plumbline-linux-amd64 -register register.md .
```

Point the final argument(s) at the paths to scan. Exit codes:

| Code | Meaning |
|---|---|
| `0` | clean — every requirement covered, every anchor resolves |
| `1` | gaps found (uncovered, broken, orphan, transitive, or structural) — **fails the build** |
| `2` | error (bad register path, unreadable tree) |

Broken anchors and structural errors always fail. For a coverage floor instead of strict
gating, add `-min-coverage 80`; for the machine-readable report, add `-json`. Config can
live in a `plumbline.json` (see [the example](../plugins/plumbline/docs/plumbline.example.json)).

## GitHub Actions

```yaml
name: traceability
on: [push, pull_request]

jobs:
  plumbline:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Fetch the Plumbline engine
        env:
          V: v0.2.1
        run: |
          base="https://github.com/dr-zog/plumbline/releases/download/$V"
          curl -sSLO "$base/plumbline-linux-amd64"
          curl -sSLO "$base/SHA256SUMS"
          sha256sum --ignore-missing -c SHA256SUMS
          chmod +x plumbline-linux-amd64
      - name: Traceability gate
        run: ./plumbline-linux-amd64 -register register.md .
```

## GitLab CI

```yaml
plumbline:
  stage: test
  image: alpine:latest      # any base with curl + coreutils
  variables:
    V: v0.2.1
  script:
    - base="https://github.com/dr-zog/plumbline/releases/download/$V"
    - wget -q "$base/plumbline-linux-amd64" "$base/SHA256SUMS"
    - sha256sum --ignore-missing -c SHA256SUMS
    - chmod +x plumbline-linux-amd64
    - ./plumbline-linux-amd64 -register register.md .
```

## Other ways to get the binary

- **Always latest** (floating, *not* reproducible — avoid for a gate you want stable).
  The `dist` branch serves the current build at a fixed path:
  `https://raw.githubusercontent.com/dr-zog/plumbline/dist/plugins/plumbline/bin/plumbline-linux-amd64`
- **From source.** Zero dependencies, so a tagged build needs only the Go toolchain:
  `go install github.com/dr-zog/plumbline/cmd/plumbline@v0.2.1` (note: `-version` prints
  `dev` this way, since the linker stamp isn't applied — irrelevant for running the gate).
- **Vendored.** Commit a pinned binary into the consuming repo (e.g. under `tools/`) and
  call it directly — hermetic, no download step.

For an AI setting this up on your behalf, the plugin's `enforce` skill produces the same
gate wired into your CI.
