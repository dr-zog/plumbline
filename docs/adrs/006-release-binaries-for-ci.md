# 006. GitHub Release binaries for pinned consumption

- **Status:** Accepted
- **Date:** 2026-08-28
- **Deciders:** Dr. Zog (operator), analysis by Claude

## Context

ADR 005 automated releases and, deliberately, made GitHub Releases **notes-only** —
binaries lived solely on the force-pushed `dist` branch, on the reasoning that one
artifact location is simpler and `dist` already serves the plugin. That held for the
consumer we had in mind: the Claude Code plugin, installed from `#dist`.

It missed a second consumer the automation itself created. The release process is now
good enough to run in *any* repo's CI as a traceability gate — but the only way to get
the engine binary is the `dist` raw URL, and `dist` is **force-pushed every release**,
so it floats to the latest build. A CI job cannot pin a specific, reproducible version,
and old builds are unreachable once the branch moves. ADR 005 foresaw exactly this —
*"if we ever genuinely needed immutable per-version binary archives, Release assets
would be the place"* — and flagged it as a later step. The need has arrived.

## Decision

Publish the four platform binaries **and a `SHA256SUMS`** as **assets on every GitHub
Release**, alongside — not instead of — the `dist` branch. Two complementary channels,
both produced by the **same build in the same release run**, so they cannot diverge:

- **`dist` branch** — the Claude plugin marketplace (`#dist`); floating latest. **Unchanged.**
- **Release assets** — pinned, per-version downloads for CI and direct use:
  `…/releases/download/vX.Y.Z/plumbline-<platform>`, verifiable against `SHA256SUMS`.

**Mechanism.** semantic-release's `github` plugin `assets` (previously `[]`) uploads the
binaries `make dist` has already built; the checksum file is generated right after
`make dist` and uploaded with them. No second build, no new job.

This **supersedes ADR 005's "Releases are notes-only; binaries live solely on `dist`"
sub-decision** (ADR 005 otherwise stands).

## Consequences

**Positive**

- CI in any repo can pin an exact Plumbline build by version, with an integrity check —
  the reproducibility `dist`'s floating URL can't offer.
- No divergence: both channels are the identical artifacts from one run.
- The `dist` / plugin flow and `main` (still binary-free) are untouched.

**Negative / accepted**

- Release assets **accumulate** — every release keeps its binaries, the opposite of
  `dist`'s no-accumulation ethos. But that permanence is exactly the point for pinned
  consumption, and the storage sits off `main` and off clones.
- Two channels to reason about — mitigated by the same-build-feeds-both rule; there is
  never a "which binary is canonical" question within a version.
- `v0.2.1` predates this and is **backfilled once** (its `dist` binaries uploaded to the
  existing Release) so pinned URLs work from `0.2.1` onward.

## Alternatives considered

- **Keep notes-only; CI builds from the tag** (`go install …/cmd/plumbline@vX.Y.Z`).
  Rejected *as the only option*: it needs a Go toolchain in the consumer's CI and skips
  the version stamp. A checksummed, downloadable binary is lower-friction and
  language-agnostic. (Building from the tag remains a fine choice — just not the only one.)
- **Per-version `dist/vX.Y.Z` tags** to make `dist` raw URLs durable. Rejected: it
  reinvents what Release assets do natively and keeps downloads awkwardly coming from a
  branch tip.
- **Release assets *instead of* `dist`.** Rejected: the plugin marketplace installs from a
  git ref (`#dist`); that channel stays exactly as ADR 001 set it.

## Provenance

- [ADR 005](005-automated-releases-semantic-release.md) — automated releases and the
  notes-only sub-decision superseded here, including its explicit "if a real need appears"
  note; [ADR 001](001-plugin-packaging-and-distribution.md) — the `dist` flow, unchanged.
- Session, 2026-08-28: the realisation that the release *process* is CI-usable while the
  only binary path (`dist` raw URLs) floats to latest with no way to pin; the two-channel,
  same-build design; the `SHA256SUMS` integrity file and the one-time `v0.2.1` backfill.
- Ticket #25 (implementation).
