# CI gate snippets

Add one of these to run the Plumbline engine as a build gate. The engine exits
non-zero on any gap, so no extra assertion logic is needed — a failing pipeline
*is* the gate.

## Getting the binary into CI

The engine ships inside the plugin, not the project, so a CI job must obtain it.
Prefer a **pinned release binary**: every
[release](https://github.com/dr-zog/plumbline/releases) attaches the four platform
binaries and a `SHA256SUMS`, so CI fetches an exact, verifiable build.

```sh
V=__PLUMBLINE_VERSION__   # pinned release tag, e.g. vX.Y.Z — enforce fills in the plugin's own
base="https://github.com/dr-zog/plumbline/releases/download/$V"
curl -sSLO "$base/plumbline-linux-amd64"
curl -sSLO "$base/SHA256SUMS"
sha256sum --ignore-missing -c SHA256SUMS
chmod +x plumbline-linux-amd64
```

`__PLUMBLINE_VERSION__` is a placeholder: substitute the release tag matching the
installed plugin (`v` + the `version` in its `.claude-plugin/plugin.json`) before writing
the job, so CI runs the same engine the skills describe. Pinning to a tag means a new
Plumbline release can't silently change the gate. Other
platforms: `plumbline-linux-arm64`, `plumbline-darwin-universal`,
`plumbline-windows-amd64.exe`. Alternatives: **vendored** (commit a pinned binary under
`tools/` and call it — hermetic, no download) or **from source**
(`go install github.com/dr-zog/plumbline/cmd/plumbline@$V`, if the team has Go in CI).

## GitHub Actions (`.github/workflows/plumbline.yml`)

```yaml
name: plumbline
on: [push, pull_request]
jobs:
  traceability:
    runs-on: ubuntu-latest
    env:
      V: __PLUMBLINE_VERSION__
    steps:
      - uses: actions/checkout@v4
      - name: Fetch the engine
        run: |
          base="https://github.com/dr-zog/plumbline/releases/download/$V"
          curl -sSLO "$base/plumbline-linux-amd64"
          curl -sSLO "$base/SHA256SUMS"
          sha256sum --ignore-missing -c SHA256SUMS
          chmod +x plumbline-linux-amd64
      - name: Traceability gate
        run: ./plumbline-linux-amd64 -register register.md .
```

## GitLab CI (`.gitlab-ci.yml`)

```yaml
plumbline:
  stage: test
  image: alpine:latest        # any base with curl + coreutils
  variables:
    V: __PLUMBLINE_VERSION__
  script:
    - base="https://github.com/dr-zog/plumbline/releases/download/$V"
    - wget -q "$base/plumbline-linux-amd64" "$base/SHA256SUMS"
    - sha256sum --ignore-missing -c SHA256SUMS
    - chmod +x plumbline-linux-amd64
    - ./plumbline-linux-amd64 -register register.md .
```

## Threshold mode

To gate on a coverage floor rather than strictly, pass `-min-coverage`:

```sh
./plumbline-linux-amd64 -min-coverage 80 -register register.md .
```

Broken anchors and structural errors still fail regardless of the floor.
