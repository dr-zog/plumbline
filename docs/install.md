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

## Publishing a new build (maintainers)

`main` never carries binaries. To publish (or refresh) the installable plugin, build
the binaries and force-push them to the flat `dist` branch:

```
make dist
```

This builds every target binary, assembles `marketplace.json` + `plugins/plumbline`
(with the binaries) into a single parentless commit, and force-pushes it to `dist` —
so the branch never accumulates binary history.
