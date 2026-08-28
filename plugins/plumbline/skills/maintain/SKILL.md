---
name: maintain
description: Keep Plumbline traceability current as code changes — on the staged diff, add or update anchors for new and changed code-areas, update the register, and clear whatever the engine flagged. Use this whenever code has changed and the anchors or register need to keep pace: before a commit, when a pre-commit traceability check is failing, when the engine reports broken or orphan or uncovered after edits, or when someone asks to "update the anchors", "re-anchor", or "fix traceability". Rides the normal change flow on an already-onboarded project — use `onboard` instead if there is no register yet.
---

<!-- [impl->component~maintain-skill~1] -->

# Plumbline — maintain

Ride the change. When code moves, the anchors and the register have to move with it,
or the link between requirement and code silently rots — which is the exact failure
Plumbline exists to prevent. Work from the staged diff: anchor what's about to be
committed, no more.

For the full register/anchor grammar and the C4 ladder, read
`../onboard/references/register-authoring.md`; the essentials are inline below.

## Before you start

- **A register must exist.** If there's none, this project isn't onboarded — use
  `onboard`. Maintain keeps an existing structure current; it doesn't build one.
- **Have the engine ready** — the plugin's launcher `${CLAUDE_PLUGIN_ROOT}/bin/plumbline`,
  or `go build -o /tmp/plumbline ./cmd/plumbline` from source — so you can verify as you go.

## Read the staged diff

```
git diff --cached --name-status
```

This is your work-list: the added (`A`), modified (`M`), renamed (`R`) and deleted
(`D`) source files about to land. Group them by **code-area** (package / directory),
because the anchor unit is the area, not the file.

## Handle each changed code-area

- **New code-area** (a new package/directory, or the first substantive file in one)
  — decide which `component` it belongs to. If it fits an existing component, add an
  anchor pointing at it. If it's genuinely new functionality, add a new `component`
  to the register, wire its `Covers`/`Needs`, and then anchor the code. Stay coarse:
  one anchor for the area, dropping to a marked region only where an area hosts
  several components.
- **Changed code-area** — confirm its anchor still resolves and still tells the
  truth. If the area's *purpose* has shifted (not just its internals), update the
  target, or bump the covered item's revision so stale links are re-checked.
- **Renamed / moved code-area** — carry its anchor with it; the tag lives in the
  code, so a plain move keeps it, but a split or merge needs the anchors re-placed.
- **Deleted code-area** — this is the important one. Removing the code that covered a
  `component` doesn't mean the requirement above it disappeared. Remove the now-dead
  anchor, then **surface the newly-uncovered item to the developer** — is the
  requirement still met another way, or genuinely dropped? Don't quietly delete the
  register entry to make the build green; that erases exactly the record Plumbline is
  meant to keep.

**Name the exact edit owed.** You have the diff; the engine doesn't. So don't stop at
"update the register" — for each change, tell the developer the *specific* one-line
edit: the file, the item, and the change. For example: *"`component~router` moved to
`internal/http/router` — update its description and re-anchor `router.go`"*, or
*"deleting `pkg/legacy` leaves `req~legacy-export` uncovered — retire it, or point it at
its replacement."* The engine's `→ fix:` hints name the *kind* of edit; you, reading the
move, name the *exact* one. That's the difference between a chore and a checkbox.

## Keep the status current

A requirement's `Status` is its maturity, and it moves as the work does (ADR 004):

- **Building a `proposed`/`draft` requirement** — leave the status as-is; code against it is
  an expected **build-ahead** warning, not a failure.
- **Finishing and committing to it** — promote it to `Status: approved`. That clears the
  **status-lag** warning, and from now on the gate holds it: it must stay covered.
- **A new requirement in this change** — set its status deliberately: `approved` if it's a
  committed part of what you're shipping, `proposed` if it's designed-ahead and not yet built.
- **Dropping a requirement** — set it `rejected` (and remove its anchored code, or the engine
  reports **zombie code**). Don't silently delete it; the rejection is part of the record.

A bare item defaults to `proposed`, so if you add a requirement you *are* committing to, say
`Status: approved` explicitly — otherwise the gate won't enforce it.

## The anchor, in brief

<!-- oft:off -->
```go
// [impl->component~auth-validator~1]
```
<!-- oft:on -->

`[<covering-type> -> <component-id>]`, where covering type is `impl`, `utest` or
`itest`. Place it in a real comment for the file's language, in the area's
package-level doc comment or a designated file.

## Close the loop

1. Run `plumbline -json` and read what it flags. Everything traceable to *this*
   change is yours to clear: broken anchors you introduced, orphans you added,
   items your change left uncovered, any **dead-end** (an approved item with no `Needs` —
   give it a `Needs` edge, or a lower status), and any **zombie code** (an anchor left on a
   `rejected` item — remove the code, or un-reject the item). **Warnings** — build-ahead or
   status-lag on a not-yet-approved item you touched — don't fail the build; reconcile them
   by promoting a finished item to `approved`, but they're a nudge, not a blocker.
2. Pre-existing gaps unrelated to the diff aren't your job to fix here — note them
   and move on (or point the developer at `audit`). Maintain rides the change; it
   isn't a whole-repo sweep.
3. Never add "what this does" prose comments. Anchors only. A change that keeps the
   map true is done, even if the code comments stay as sparse as they were.
