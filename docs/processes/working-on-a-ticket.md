# Working on a ticket

The development loop for a Plumbline change, from a requirement to landed code.
Plumbline is built the way Plumbline asks you to build: **the register leads the
code within a change.**

## The core principle — register-first (route 2a)

For any change that adds or alters a *feature*, the register requirement goes in
**first**, and the code makes it green **before merge**:

1. **Register the requirement first** — add/adjust the `feat`/`req`/`component`
   items in `register/`. Run `plumbline`; the new item is **uncovered** and the
   gate is **red**. That red is correct — it's the promise, not yet kept.
2. **Implement**, then **anchor** the code-area(s) that satisfy it.
3. **Go green** — `plumbline` reports 0 gaps before you open the PR.

**Never merge a red `main`.** The red state lives only on your branch. A red `main`
must always mean a real regression, never "work in progress" — that's what keeps the
gate trustworthy. (For a *trivial* change — a typo, a doc tweak, a one-line fix with
no new requirement — registering "first" is ceremony; just keep the trace green.)

The "decided but not yet started" backlog lives in **GitHub Issues** and in `docs/adrs/`,
**not** as persistent red in the register. A requirement enters the register when its
work *starts*.

## The loop

1. **Pick up the ticket** (GitHub — see [writing-tickets.md](writing-tickets.md)).
   Make sure you understand the *outcome* wanted; clarify with a human if it's a
   judgement call, not a fact.
2. **Branch** off `main` (`git checkout -b <short-topic>`).
3. **Design, lightly.** Decide the C4 shape of the change — which feature/requirement,
   which component/code-area. If the decision is architectural and would be
   re-litigated later, write an [ADR](writing-adrs.md).
4. **Register-first** (above).
5. **Implement** in `cmd/` or `internal/` (engine) or `plugins/plumbline/skills/`
   (product skills). Match the surrounding code's idiom and comment density.
6. **Anchor** the code-areas (`// [impl->component~…~1]`, `utest` on the tests).
7. **Verify** — all of:
   - `go test ./...` · `go vet ./...` · `gofmt -l .` (clean)
   - `plumbline` (self-trace green — 0 uncovered/broken/orphan/transitive/structural)
8. **Skills audit (on material change).** If this change materially **adds, removes, or
   alters platform functionality** — anything a consumer sees in the engine's output or
   behaviour — run a **full, detailed audit of the product skills**
   (`plugins/plumbline/skills/`): check *every* one against the platform's current
   capabilities, and either update it in lockstep or ticket the update. A skill that
   misdescribes what the platform does misleads the AIs that consume it — the same
   requirement↔reality rot Plumbline exists to prevent, one level up.
9. **Land it** via a GitHub PR — see [raising-a-pr.md](raising-a-pr.md) for the
   Conventional-Commit PR title (it drives the version + changelog) and release-on-merge.

## Notes

- Anchors are coarse — one per code-area (package/dir), not per file. Fence any example
  tags you write in docs/tests with `// oft:off` … `// oft:on`.
- Keep prose comments as sparse or rich as you like — Plumbline never measures them.
  The anchor is the only non-negotiable.
