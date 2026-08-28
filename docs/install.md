# Installing Plumbline (Claude Code)

Plumbline installs as a Claude Code plugin. The plugin bundles the engine binary
for every platform, so there's nothing else to install — no Go toolchain, no
runtime dependencies.

## Install

The plugin ships from the repository's **`dist`** branch (which carries the plugin
and its per-platform binaries; `main` stays source-only — see
[ADR 001](adrs/001-plugin-packaging-and-distribution.md)).

```
/plugin marketplace add https://github.com/dr-zog/plumbline.git#dist
/plugin install plumbline@plumbline
```

The first command adds this repo's `dist` branch as a marketplace; the second installs
the `plumbline` plugin from it. Claude Code fetches the plugin and its binaries; the
launcher (`${CLAUDE_PLUGIN_ROOT}/bin/plumbline`) selects the right one for your host.

## Verify

Once installed, the five skills — `onboard`, `maintain`, `enforce`, `audit`,
`showcase` — are available, and the engine runs via the launcher:

```
${CLAUDE_PLUGIN_ROOT}/bin/plumbline -version
```

## Using the engine in CI (without Claude Code)

The engine is a single zero-dependency binary, so any repo's CI can run the traceability
gate directly — no plugin, no Go toolchain.

**Pinned to a version (recommended for CI).** Every
[release](https://github.com/dr-zog/plumbline/releases) attaches the four platform
binaries and a `SHA256SUMS`, so you can fetch an exact, verifiable build
([ADR 006](adrs/006-release-binaries-for-ci.md)):

```bash
V=v0.2.1
base="https://github.com/dr-zog/plumbline/releases/download/$V"
curl -sSLO "$base/plumbline-linux-amd64"
curl -sSLO "$base/SHA256SUMS"
sha256sum --ignore-missing -c SHA256SUMS      # verify integrity
chmod +x plumbline-linux-amd64
./plumbline-linux-amd64 -register register.md .
```

The assets are `plumbline-linux-amd64`, `plumbline-linux-arm64`,
`plumbline-darwin-universal`, and `plumbline-windows-amd64.exe`.

**Always latest** (floating, *not* reproducible). The `dist` branch serves the current
build at a fixed path:

```
https://raw.githubusercontent.com/dr-zog/plumbline/dist/plugins/plumbline/bin/plumbline-linux-amd64
```

**From source.** Zero dependencies, so a tagged build needs only the Go toolchain:

```bash
go install github.com/dr-zog/plumbline/cmd/plumbline@v0.2.1
```

## Publishing a new build (maintainers)

`main` never carries binaries, and **releases are automatic**: merging a `fix` / `feat` /
breaking change to `main` triggers semantic-release, which tags the version, updates the
changelog + `plugin.json`, force-pushes the plugin + binaries to `dist`, and cuts a GitHub
Release with the binaries + `SHA256SUMS` attached
([ADR 005](adrs/005-automated-releases-semantic-release.md),
[ADR 006](adrs/006-release-binaries-for-ci.md)).

The manual fallbacks, if a release ever needs re-running by hand: `make dist` builds every
target binary and force-pushes the flat `dist` commit (the branch never accumulates binary
history), and `gh workflow run release.yml -f dry_run=false` re-fires the full pipeline.
