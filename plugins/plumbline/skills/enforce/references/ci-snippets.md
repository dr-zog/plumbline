# CI gate snippets

Add one of these to run the Plumbline engine as a build gate. The engine exits
non-zero on any gap, so no extra assertion logic is needed — a failing pipeline
*is* the gate.

## Getting the binary into CI

The engine ships inside the plugin, not the project, so a CI job must obtain it.
Pick the line that matches how the team distributes tools:

- **Vendored** — the team commits `plumbline-linux-amd64` into the repo (e.g. under
  `tools/`). Simplest and hermetic; just call it.
- **Released** — download a pinned release binary with `curl` and `chmod +x`.
- **From source** — `go build` (or `go run`) the engine if the module is available
  to CI (a submodule, a `go install`, or a vendored source tree).

The snippets below default to the vendored path with a source-build fallback; delete
whichever you don't use.

## GitLab CI (`.gitlab-ci.yml`)

```yaml
plumbline:
  stage: test
  image: golang:1.23        # only needed for the go-build fallback; use any base if vendored
  script:
    - |
      if [ -x tools/plumbline-linux-amd64 ]; then
        tools/plumbline-linux-amd64
      else
        go build -o /tmp/plumbline ./cmd/plumbline && /tmp/plumbline
      fi
```

## GitHub Actions (`.github/workflows/plumbline.yml`)

```yaml
name: plumbline
on: [push, pull_request]
jobs:
  traceability:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5   # only needed for the go-build fallback
        with:
          go-version: '1.23'
      - name: Traceability gate
        run: |
          if [ -x tools/plumbline-linux-amd64 ]; then
            tools/plumbline-linux-amd64
          else
            go build -o /tmp/plumbline ./cmd/plumbline && /tmp/plumbline
          fi
```

## Generic (any CI)

Run the binary from the repo root; fail the step on non-zero exit:

```sh
plumbline            # or ./tools/plumbline-<os>-<arch>
```

## Threshold mode

To gate on a coverage floor rather than strictly, pass `-min-coverage`:

```sh
plumbline -min-coverage 80
```

Broken anchors and structural errors still fail regardless of the floor.
