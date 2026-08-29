# Register authoring reference

The canonical reference for writing a Plumbline register and placing anchors. The
grammar is [OpenFastTrace](https://github.com/itsallcode/openfasttrace)'s — Plumbline
inherits the ecosystem rather than inventing syntax.

## Contents

- [The C4 ladder](#the-c4-ladder)
- [Register grammar](#register-grammar)
- [Anchor grammar](#anchor-grammar)
- [Worked example](#worked-example)

## The C4 ladder

The artifact types form a locked ladder — you don't invent types. "Needs" flows
from the abstract down to the concrete; two axes meet at the `component`, where code
anchors attach:

```
requirements:   feat ──▶ req ──┐
                                ├──▶ component ──▶ impl | utest | itest
architecture:   context ─▶ container ──┘              (code anchors)
```

| Type | C4 / axis | May Need | Covers (points up to) |
|---|---|---|---|
| `feat` | feature | `req` | — |
| `req` | requirement | `component` | a `feat` |
| `context` | C4 L1 system context | `container` | — |
| `container` | C4 L2 deployable unit (names tech) | `component` | a `context` |
| `component` | C4 L3 grouping | `impl`, `utest`, `itest` | a `req` and/or `container` |
| `impl` `utest` `itest` | C4 L4 code (anchors only) | — (leaves) | — |

`Covers` is the inverse of `Needs`: writing `Covers: req~x~1` on a `component` is
only valid because a `req` may `Need` a `component`. The engine enforces this, so a
back-to-front or skip-a-level link is reported as a structural error.

**Both axes meet at the `component`.** A component covers *two* things — the `req` it
realises (requirements axis) **and** the `container` it lives in (architecture axis), e.g.
`Covers: container~engine~1, req~parse-input~1`. That `component → container` edge is how a
container is realised from below; without it a container has nothing covering it and can't be
`approved` (it reads as uncovered). A `context` is likewise realised by its containers — and a
top-of-axis node covered from below needs no `Needs` of its own to be covered (ADR 007).

**One edge, two meanings.** The two axes use the same `Covers` edge for two different
relationships — a deliberate C4 + OFT fusion. On the **requirements axis**
(`impl → component → req → feat`), `Covers` is OFT's classic *fulfilment*: the coverer
*satisfies* what it points at. On the **architecture axis**
(`impl → component → container → context`), it records C4 *containment* instead: a component
is *part of* its container, realising it by aggregation, so a container is "covered" once its
constituent components are present. Both are legal OFT (the type ladder is configurable) and
faithful to C4's taxonomy — but note the shift: on the architecture axis
`Covers: container~Y~1` reads as "belongs to / helps realise", not "satisfies". Coming from
pure OFT, that's the one place `Covers` means something broader than you'd expect.

## Register grammar

The register is OFT-native Markdown (`register.md` by default). A specification item
is a backtick-wrapped ID on its own line, optionally under a heading (which becomes
its title), followed by a description and `Needs:`/`Covers:` attribute lines.

```markdown
### Auth validator
`component~auth-validator~1`

Validates credentials and issues session tokens.

Covers: req~validate-auth-request~1
Needs: impl
```

- **ID** is `type~name~revision`. `name` is kebab-case; `revision` is a positive
  integer, conventionally started at 1.
- **Description** — the first prose line after the ID. For a `container`, name the
  chosen technology here (the engine warns otherwise).
- **`Needs:`** — comma-separated artifact types this item requires coverage in.
- **`Covers:`** — comma-separated IDs of items this one covers (points up the chain).
- **`Status:`** — `approved` / `proposed` / `draft` / `rejected`. **A bare item defaults to
  `proposed`** (provisional — tracked, but *not* gated), so mark a committed requirement
  `Status: approved` for the gate to enforce it. `proposed`/`draft` with code raises a
  warning (build-ahead); `rejected` is excluded, and code anchoring a `rejected` item is
  zombie code (a hard fail). See ADR 004.
- **Revisions** — when an item's *meaning* changes, bump the revision rather than
  editing in place. This is OFT semantics: it invalidates stale links so coverers
  know to re-check. A typo fix does not warrant a bump.
- Wrap regions the engine should ignore in `<!-- oft:off -->` … `<!-- oft:on -->`.

## Splitting the register across files

The register can be a single file or a whole directory tree — Plumbline aggregates
every `.md` it finds into one item set, exactly as OpenFastTrace does. Point
`register` at a directory (or list several files/dirs/globs under `registers` in
`plumbline.json`), and organise however suits the codebase:

- **by level** — `context.md`, `containers.md`, `components.md`, `requirements.md`;
- **by container** — `components/api.md`, `components/worker.md`;
- **co-located** near the code the items describe.

`Covers`/`Needs` links resolve across files, because everything is aggregated before
resolution. Two files must not define the same ID — a duplicate is a structural error.

## Anchor grammar

An anchor is a coverage tag in a code comment. Full form:

```
[<covering-type> -> <target-id>]
```

```go
// [impl->component~auth-validator~1]
func validateAuthRequest(token string) bool { ... }
```

- `<covering-type>` is a code leaf: `impl`, `utest` or `itest`.
- `<target-id>` is the full ID of the register item this code-area covers —
  almost always a `component`.
- Optional covering name/revision: `[impl~~2->component~x~1]` or
  `[impl~validate-password~2->component~x~1]`.
- The comment character doesn't matter — the engine matches the pattern, which by
  convention lives in a comment. It scans source files by extension (`.go`, `.py`,
  `.ts`, `.rs`, `.java`, …), so place anchors in a real comment for that language.

## Worked example

A complete four-level chain — feature down to code — plus the anchor that grounds it.

Register:

```markdown
## Features
### Authentication
`feat~authentication~1`
Status: approved

Users can authenticate before accessing the system.

Needs: req

## Requirements
### Validate authentication request
`req~validate-auth-request~1`
Status: approved

Every authentication request is validated before access is granted.

Covers: feat~authentication~1
Needs: component

## Components
### Auth validator
`component~auth-validator~1`
Status: approved

Validates credentials and issues session tokens. (In the "api" container: a Go HTTP service.)

Covers: req~validate-auth-request~1
Needs: impl
```

Code:

```go
// Package auth validates credentials and issues session tokens.
//
// [impl->component~auth-validator~1]
package auth
```

Resolution: the anchor provides the `impl` the component `Needs` → the component
Covers the requirement → the requirement Covers the feature. The whole chain is
deeply covered, and the engine reports it green.
