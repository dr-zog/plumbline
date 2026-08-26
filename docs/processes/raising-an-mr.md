# Raising a pull request

How a change lands on `main`: a **GitHub pull request**, even solo. The PR is the
review checkpoint, the CI gate, and the paper trail.

## Before you open it

Everything green on the branch:

- `go test ./...` · `go vet ./...` · `gofmt -l .` (no output)
- `plumbline` — self-trace green, **0 gaps** (register-first work is now covered)
- Commits are clean and logically grouped; each message says *what* and *why*, not a
  changelog of files. End every commit message with the `Co-Authored-By` trailer.
- Rebased on the latest `main` so the merge is a clean fast-forward.

## Open the PR

`gh pr create` (or the web UI). The PR body carries the **Dr Zog voice** — British
English, confident, mechanism-first, no hype. Structure:

- **What** — the outcome delivered, in a sentence.
- **Why** — the requirement/ticket it satisfies (link it).
- **How** — the shape of the change; call out any trade-off *honestly*.
- **Proof** — the checks that pass (tests, vet, and the green Plumbline scorecard —
  a change to Plumbline should show Plumbline still traces itself).

Keep it skimmable. A reviewer should grasp the change from the body without reading
the diff first.

## Merge & clean up

- Merge only when CI is green and any review is resolved. Prefer a **fast-forward /
  clean linear history** (no noise merge commits for a solo change).
- **Post-merge hygiene:** delete the merged branch (local and remote), and confirm
  `main` is green (`plumbline`, tests) after the merge.
- If the branch produced release artifacts, that's the release flow's job (see
  [ADR 001](../adrs/001-plugin-packaging-and-distribution.md)) — never commit binaries
  to `main`.

## Git conventions (repo-wide)

- Branch off `main`; never commit directly to `main`.
- End commit messages and PR descriptions with:
  `Co-Authored-By: Claude <noreply@anthropic.com>` (per the repo's authorship rule).
- Interactive rebases/adds aren't available in the agent environment — structure commits
  as you go instead.
