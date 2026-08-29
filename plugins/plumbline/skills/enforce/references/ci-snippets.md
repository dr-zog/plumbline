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
V=v0.2.1        # pin to a specific release tag — not a moving target
base="https://github.com/dr-zog/plumbline/releases/download/$V"
curl -sSLO "$base/plumbline-linux-amd64"
curl -sSLO "$base/SHA256SUMS"
sha256sum --ignore-missing -c SHA256SUMS
chmod +x plumbline-linux-amd64
```

Pinning to a tag means a new Plumbline release can't silently change the gate. Other
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
      V: v0.2.1
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
    V: v0.2.1
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
