---
name: onboard
description: Bring an untracked codebase under Plumbline traceability from zero — assess the code, build the C4 register (context/container/component plus features and requirements), and place the initial anchors. Use this whenever someone wants to onboard, adopt, set up or introduce Plumbline on an existing project, "plumbline this repo", create a register from scratch, establish requirement↔code traceability where none exists, or when the engine reports there is no register yet. This is the zero→structured entry point: reach for it before `maintain` or `audit` on any project that has no register.
---

<!-- [impl->component~onboard-skill~1] -->

# Plumbline — onboard

Take a project from no traceability to a standing C4 register with initial anchors
in place. You are building the durable link — requirement ↔ feature ↔ code-area —
one altitude above the code, so keep it coarse and keep the developer in the loop
on anything that is a judgement call rather than a fact.

Read `references/register-authoring.md` for the exact register grammar, the C4
ladder, the anchor format, and a worked example. This file is the *process*; that
file is the *reference*.

## Before you start

1. **Check the engine is available** — run it via the plugin's launcher,
   `${CLAUDE_PLUGIN_ROOT}/bin/plumbline` (it picks the right per-platform binary), or
   build one with `go build -o /tmp/plumbline ./cmd/plumbline` from the engine's
   source. Without it you can't verify your work.
2. **Check there is no register already.** If `register.md` (or a configured
   register) exists, this project is already onboarded — hand off to `maintain` or
   `audit` instead. Onboarding twice would fight an existing structure.
3. **Get oriented.** Survey the tree: languages, top-level directories, entry
   points, build/deploy files, and any existing docs (README, ADRs, an issue
   tracker). This is your raw material — the register describes *this* system, not
   a generic one.

## Build the register top-down, confirming each level

C4 zooms from the outside in. Work the same way, and **confirm each level with the
developer before descending** — architecture and requirements are decisions, not
facts you can read off the source, and a wrong guess at a high level cascades. Write
each level into `register.md` as you settle it, so the register grows with the
conversation.

1. **Context** — the system itself and the people and external systems it talks to.
   Usually one `context` item. State what the system is *for*. Confirm you've named
   the right system boundary.
2. **Containers** — the separately runnable/deployable units: services, CLIs, web
   apps, workers, datastores. Each `container` **names its chosen technology** in
   its description (this is a C4 rule and the engine warns when it's missing —
   "a Go HTTP service", "a Postgres database"). Propose the set and get sign-off;
   this is the shape of the system.
3. **Components** — the major groupings of functionality inside each container.
   Keep them coarse: a component is a subsystem, not a class. Propose per container,
   confirm.
4. **Requirements axis** — the `feat`ures and `req`uirements the components exist to
   satisfy. Mine these from READMEs, docs and the issue tracker where they exist;
   where they don't, infer candidates from the code but **lean hard on the developer
   to confirm** — a requirement is a statement about intended outcomes, which only
   they can ratify. It's fine to start sparse and let `maintain` grow this over time.
   **Mark the requirements you and the developer commit to `Status: approved`** — a bare
   item defaults to `proposed` (provisional, *not* gated), so an un-approved requirement is
   tracked but won't be enforced until it's approved (ADR 004).

Write each level into the register as you settle it — a single `register.md` for a
small project, or a `register/` directory split by level (and, at scale, by
container) as it grows; Plumbline aggregates every `.md` into one register either way.

As you go, wire the chain with `Covers`/`Needs`: a `component` Covers the `req` it
realises, a `req` Covers its `feat`, a `container` Covers the `context`. The engine
validates these against the locked ladder, so a mis-wired link is caught immediately.
Every `approved` item must declare a downward `Needs` — an approved item with no `Needs` is
a **dead-end** (a hard fail), so wire it before you approve it.

## Place the anchors — adaptively

An anchor links a code-area up to the `component` it implements:

<!-- oft:off -->
```go
// [impl->component~auth-validator~1]
```
<!-- oft:on -->

Default to **one anchor per package or source directory**, in the package-level doc
comment or a single designated file for that area — coarse, predictable, one per
code-area. **Drop to a finer marked region only where one directory genuinely hosts
several components**; don't split hair by hair. The point is a durable map, not a
census of every file.

Match the anchor's covering type to what the code is: `impl` for implementation,
`utest`/`itest` for test suites that verify a component.

## Verify and hand off

1. Run the engine (`plumbline -json`) and read the scorecard. Early on you'll see
   uncovered items and orphans — that's expected; the register outran the anchors,
   or vice versa. Close the gaps you can and **surface the ones you can't** to the
   developer rather than papering over them.
2. Aim to leave the project with a register that parses, anchors that resolve, and
   any remaining gaps understood and agreed — not necessarily a green build on day
   one.
3. Suggest `enforce` if they want the discipline gated from here on, and `maintain`
   as the way to keep it current as the code changes.

## What onboarding is not

Never add descriptive "what this does" prose comments — Plumbline doesn't measure
them and they rot. You place anchors and author the register; that's the whole job.
Resist over-documenting: a sparse register that's true beats a rich one that's
guessed.
