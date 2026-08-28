# 005. Automated releases via semantic-release

- **Status:** Accepted
- **Date:** 2026-08-28
- **Deciders:** Dr. Zog (operator), analysis by Claude

## Context

v0.1.0 and v0.2.0 shipped by hand: a hand-authored changelog, a manual `git tag`,
a manual `make dist`, a manual GitHub Release. Nothing tied those steps to one
source of truth, and the version duly drifted — the plugin manifest reached `0.5.0`,
unrelated to the released `v0.2.0` (corrected in #19). The manual, multi-source flow
was the root cause, not a slip.

We want the release flow to satisfy a specific set of requirements:

- **Conventional Commits, fully** — every change classified, and enforced.
- **Automatic release on every code merge** — no batching, no separate release
  step; the reviewed PR merge *is* the release.
- **Binaries live only on the `dist` orphan commit** — Releases are notes-only
  (ADR 001 stands).
- **The bumped version is written back into `main`** — the manifest on `main`
  reflects the released version.
- **The changelog is auto-generated** and lives in `main`.
- **Drift is impossible**, not merely unlikely.
- **Normal, off-the-shelf tooling**, and **product code still lands only via
  reviewed PRs**.

The earlier written rule — *"the changelog is never generated from commit
messages"* — was a **proxy** for the real requirement (curated, human-facing
release notes), written before we planned automated releases. With Conventional
Commits carrying the intent and the generated notes remaining editable, the real
requirement is met by other means, so the proxy is retired.

## Decision

Adopt **[semantic-release](https://semantic-release.gitbook.io/)** as the release
engine, triggered on every push to `main`.

- **Conventional Commits** drive the semver bump (`fix` → patch, `feat` → minor,
  breaking → major; `docs`/`chore`/`refactor` → no release). Enforced by a
  **required PR-title lint** — we squash-merge, so the PR title is the commit.
- On each qualifying merge, semantic-release computes `vX.Y.Z`, generates the
  release notes, **writes `CHANGELOG.md` and bumps `plugin.json`**, commits them
  back to `main` with `[skip ci]` (`@semantic-release/git`), creates the **tag**
  and a **notes-only GitHub Release**, and runs
  **`make dist VERSION=vX.Y.Z`** to build and force-push the binaries to the orphan
  `dist` commit.
- **Releases carry no binary assets** — binaries live solely on `dist` (ADR 001).
- **Zero drift, by construction and enforcement.** The version is computed once and
  written to every location by the bot (a single writer, so it cannot produce
  disagreement); a **`v*` tag-protection rule** stops any hand-cut tag; and a
  **required CI check (`plugin.json.version == latest tag`)** blocks any human PR
  that edits the version out of band. The only stored version is the git tag;
  everything else is written equal to it or derived from it.
- **Permissions.** A dedicated **GitHub App** (a bot identity — keeps any personal
  identity and personal tokens out of the pipeline) with repository **Contents** and
  **Pull-requests** write, added as a **bypass actor** on `main`'s protection (solely
  to push its one `[skip ci]` release commit) and on the `v*` tag rule. Product code
  still lands only through reviewed PRs; only the *derived* release bookkeeping
  auto-commits.

This **supersedes ADR 001's release sub-decisions**: point 5 (pin and bump a
marketplace `sha` each release) is dropped — the marketplace stays `ref: dist`
(a floating install), because the `dist` sha is only known *after* `make dist`, so
pinning it would need a second write-back and a bootstrapping loop. ADR 001's
packaging shape — the orphan `dist` branch, the four static binaries, binaries off
`main` — stands unchanged.

## Consequences

**Positive**

- Version drift becomes impossible: one computed number, one writer, a protected
  tag, and a merge-blocking equality check.
- Releases are automatic *and* deliberate — the reviewed PR merge is the trigger, so
  the human gate we already have is the release gate.
- The changelog and manifest version are machine-maintained; `main` history stays
  clean (binaries only on `dist`).
- Normal, widely-used tooling — no bespoke release engine to maintain.

**Negative / accepted**

- The release App **bypasses branch protection** for its single `[skip ci]` commit.
  A narrow, conscious loosening: that commit is *derived bookkeeping* computed from
  already-reviewed commits — there is nothing new in it to review.
- **Version numbers climb quickly** (a release per code merge). Intended, not a flaw.
- semantic-release's generated changelog format differs slightly from the
  hand-authored Keep-a-Changelog; the existing `0.1.0`/`0.2.0` entries are preserved
  and new versions are prepended above them.
- These are **release-engineering decisions, recorded here — not modelled in
  Plumbline's product register.** The register traces the shipped tool (engine +
  plugin), not the CI that builds and ships it; treating our pipeline as a traced
  system would be a category error (see *Alternatives*).

## Alternatives considered

- **release-please.** Rejected: its purpose is to *batch* changes into a Release PR —
  structurally the opposite of release-on-merge. A good tool, the wrong cadence for us.
- **A bespoke ~50-line release workflow.** Rejected: it reinvents a solved problem;
  the operator wants normal, battle-tested tooling, and the zero-dependency promise is
  about the shipped Go engine, not the CI toolchain.
- **Keep the hand-authored changelog.** Rejected: it was a proxy for curation and it
  blocks automation; Conventional Commits plus editable notes satisfy the real
  requirement.
- **Derive `plugin.json` at publish (a sentinel on `main`, no write-back).** Rejected:
  the operator wants the version *written back into `main`*; with tag protection and
  the equality check, the committed copy is drift-proof anyway.
- **Have the bot open and auto-merge a release PR** (to avoid the branch bypass).
  Rejected: a bot PR after *every* feature merge — noisy, an extra CI round-trip each
  time — for no gain over a `[skip ci]` commit whose content is fully derived.
- **Model the release requirements in Plumbline's register and anchor the pipeline.**
  Rejected by the operator: the register traces the product, not the project's
  infrastructure; these are decisions with provenance (this ADR), not product
  requirements.

## Provenance

- The v0.1.0 / v0.2.0 manual releases; the `plugin.json` `0.5.0` drift and its
  correction (#19) — the concrete failure this decision closes.
- [ADR 001](001-plugin-packaging-and-distribution.md) — packaging & distribution;
  the orphan `dist` branch stands, its release sketch and `sha`-pin are superseded here.
- Session, 2026-08-28: the requirements-driven tool selection (the seven requirements
  above); the drift analysis (a single-source tag, made impossible-to-diverge by tag
  protection + a `plugin.json == tag` check + a single-writer bot); the
  branch-bypass trade-off; and the decision *not* to model these as Plumbline product
  requirements.
- Ticket #10 (release automation) — implements this ADR.
