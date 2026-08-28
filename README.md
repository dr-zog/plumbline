# Plumbline

**AIs write code as fast as humans ask for features. The faster your agents write, the
faster implementation drifts from intent. Plumbline holds the thread — and fails the build
the moment it snaps.**

[![CI](https://github.com/dr-zog/plumbline/actions/workflows/ci.yml/badge.svg)](https://github.com/dr-zog/plumbline/actions/workflows/ci.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
![Go 1.23](https://img.shields.io/badge/Go-1.23-00ADD8?logo=go&logoColor=white)
![Zero dependencies](https://img.shields.io/badge/dependencies-zero-brightgreen)
![Claude Code plugin](https://img.shields.io/badge/Claude%20Code-plugin-8A63D2)

Agents don't tire and they don't slow down — they generate code faster than anyone can
read it, and somewhere in that torrent the *why* gets buried: the feature that was asked
for, the requirement it was meant to satisfy, the thread from intent to implementation. It
drifts out of sync while the machine keeps typing, and nobody notices until the codebase no
longer matches what anyone actually wanted. Plumbline is the counterweight — a
zero-dependency engine that enforces the link from every requirement to the code that
fulfils it, in both directions, and won't let the thread break silently. Your AIs write;
Plumbline makes sure they never lose the features amongst the code.

Write your prose comments however you like — rich or sparse, Plumbline never
measures them. Its single non-negotiable is the **anchor**: a lightweight,
machine-checkable marker that links each *code-area* up to a documented
requirement, feature or component. A zero-dependency Go binary proves the
anchors resolve, in both directions, across the whole chain — and fails the
build when they don't. The skills place the anchors for you and keep the
register current.

> **Status: usable and self-hosted.** The engine resolves coverage across the full
> requirement→architecture→code chain, validates C4 structure, scores coverage and
> completeness, and gates strictly or by threshold. The Claude Code plugin ships the
> `onboard` / `maintain` / `enforce` / `audit` / `showcase` skills — and Plumbline
> traces itself (100%). See [`DESIGN.md`](DESIGN.md) for the design rationale.

## Install

### As a Claude Code plugin

Plumbline ships from the repo's **`dist`** branch, which carries the plugin and its
per-platform binaries (`main` stays source-only). Nothing else to install — no Go
toolchain, no runtime dependencies:

```
/plugin marketplace add https://github.com/dr-zog/plumbline.git#dist
/plugin install plumbline@plumbline
```

The `onboard` / `maintain` / `enforce` / `audit` / `showcase` skills become available,
and the engine runs via the launcher, `${CLAUDE_PLUGIN_ROOT}/bin/plumbline`.

### In CI (no plugin, no Go toolchain)

The engine is a single zero-dependency binary, so any repo's CI can run the gate
directly. Pin a verified build from a
[release](https://github.com/dr-zog/plumbline/releases):

```bash
base="https://github.com/dr-zog/plumbline/releases/download/v0.2.1"
curl -sSLO "$base/plumbline-linux-amd64"
curl -sSLO "$base/SHA256SUMS"
sha256sum --ignore-missing -c SHA256SUMS          # verify integrity
chmod +x plumbline-linux-amd64
./plumbline-linux-amd64 -register register.md .   # exits non-zero on any gap
```

Full guide — a GitHub Actions example, the floating-latest and build-from-source
options — in [**docs/ci.md**](docs/ci.md).

## The idea in one breath

Codebases built feature-first lose the thread between the outcome someone wanted
(the *requirement*), the feature that meets it, and the code that implements it.
Descriptive "what this does" comments rot against the source. The durable link
sits one altitude up — **requirement ↔ feature ↔ code-area** — kept current by
process rather than by hand-maintained prose.

Plumbline is an assembly of three established disciplines, followed as-is:
[C4](https://c4model.com) for the layered register,
[requirements traceability](https://github.com/itsallcode/openfasttrace) for the
bidirectional audit, and
[Living Documentation](https://www.oreilly.com/library/view/living-documentation-continuous/9780134689418/)
for storing the markers on the code itself.

## The anchor

An anchor is an [OpenFastTrace](https://github.com/itsallcode/openfasttrace)-compatible
coverage tag in a comment on a code-area. It names the covering artifact type,
an arrow, and the full ID of the documented thing it satisfies:

```go
// [impl->component~auth-validator~1]
func validateAuthRequest(token string) bool { ... }
```

The **code-area** — a package, directory or marked region, not a symbol — is the
unit. Coarse granularity is what keeps it light-touch. Because the grammar is
OFT's, the whole artifact set stays legible to the wider tooling ecosystem.

## The C4 model

The register's artifact types form a locked ladder (adopted prior art, no
bespoke variation). "Needs" flows from the abstract down to the concrete; two
axes meet at the component, where code anchors attach:

```
requirements:   feat ──▶ req ──┐
                                ├──▶ component ──▶ impl | utest | itest
architecture:   context ─▶ container ──┘              (code anchors)
```

The engine validates the register against this ladder: unknown types, disallowed
`Needs` edges, and `Covers` links whose target is missing or the wrong type are
structural **errors** that fail the build. A container that names no technology
is a **warning**.

## The register

The documented things anchors resolve up to live in an **OFT-native Markdown**
register (`register.md` by default):

```markdown
### Auth validator
`component~auth-validator~1`

Validates credentials and issues session tokens.

Covers: req~validate-auth-request~1
Needs: impl
```

At scale the register is a *tree*, not a file: point `register` at a directory (or
list files/dirs/globs under `registers`) and Plumbline aggregates every `.md` into
one register — split by level or by container as you like, with `Covers`/`Needs`
resolving across files.

## The engine

A single static Go binary — the law. Deterministic, cheap, zero runtime
dependencies. It scans a tree for anchors, parses the register, resolves
coverage across the chain, and emits a **bidirectional report**:

- **coverage %** — items whose direct needs are met (shallow);
- **completeness %** — items whose whole chain resolves (deep);
- **uncovered** — a documented item nothing provides coverage for;
- **transitive defects** — direct needs met, but a coverer below isn't;
- **broken anchors** — a tag whose target isn't in the register;
- **orphan code-areas** — scanned source files carrying no anchor;
- **structural** — register violations of the C4 ladder.

It exits non-zero on any gap, so it doubles as a pre-commit / CI gate with no LLM
in the loop.

```console
$ go build -o plumbline ./cmd/plumbline
$ ./plumbline -register testdata/fixture/register.md testdata/fixture/src
Plumbline traceability report
  register items : 5
  anchors found  : 2
  coverage       : 4/5 items (80.0%)   [direct needs met]
  completeness   : 3/5 items (60.0%)   [whole chain resolves]
  gate           : strict — any gap fails

UNCOVERED (1) — needed coverage with nothing to provide it:
  req~rotate-signing-keys~1 — Rotate signing keys  missing: component  (…register.md:35)

TRANSITIVE DEFECTS (1) — direct needs met, but a coverer below isn't:
  feat~key-management~1 — Key management  weak: req~rotate-signing-keys~1  (…register.md:18)

BROKEN anchors (1) — tag target not in register:
  [impl->component~does-not-exist~1]  (…src/broken.go:6)

ORPHAN code-areas (1) — source files with no anchor:
  testdata/fixture/src/orphan.go

FAIL — 0 structural error(s), 1 uncovered, 1 transitive, 1 broken, 1 orphan.
$ echo $?
1
```

Add `-json` for the machine-readable report the skills consume.

### Gating

By default the gate is **strict**: any gap fails. Pass `-min-coverage N` (or set
`minCoverage` in config) for **threshold** gating — fail only below `N%` shallow
coverage. Broken anchors and structural errors always fail, regardless of the
threshold.

### Config

Flags override an optional `plumbline.json`, which overrides the baked defaults.
The artifact-type ladder is *not* configurable — only paths, exclusions and the
threshold are (see [`plugins/plumbline/docs/plumbline.example.json`](plugins/plumbline/docs/plumbline.example.json)):

```json
{ "register": "register.md", "roots": ["."], "exclude": ["testdata"], "minCoverage": 0 }
```

## Layout

| Path | What |
|---|---|
| `cmd/plumbline/` | the engine's `main` |
| `internal/anchor/` | tag scanner |
| `internal/register/` | OFT-native Markdown register parser |
| `internal/c4/` | locked artifact-type ladder + structural validation |
| `internal/report/` | coverage resolver + bidirectional JSON/text report |
| `internal/config/` | optional JSON config, merged over locked defaults |
| `testdata/fixture/` | tiny end-to-end fixture (a clean chain + one of each gap) |
| `plugins/plumbline/` | the shippable Claude Code plugin (`skills/`, `bin/`, `docs/`) |

## Develop

```console
make vet test      # go vet + unit/end-to-end tests
make build         # local ./plumbline (gitignored)
make binaries      # per-platform binaries in plugins/plumbline/bin/ (gitignored; shipped on dist)
```

CI (GitHub Actions) vets, tests and cross-compiles on every push and pull request.
**Releases are automated**: merging a `fix` / `feat` / breaking change to `main` cuts a
release via semantic-release — binaries ship on the `dist` branch (for the plugin) and as
GitHub Release assets with checksums (for pinned CI use), never committed to `main`
(ADRs [001](docs/adrs/001-plugin-packaging-and-distribution.md) /
[005](docs/adrs/005-automated-releases-semantic-release.md) /
[006](docs/adrs/006-release-binaries-for-ci.md)).

MIT licensed.
