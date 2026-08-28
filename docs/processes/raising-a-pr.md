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

**No manual changelog.** Since [ADR 005](../adrs/005-automated-releases-semantic-release.md)
the changelog is generated at release by semantic-release from the merged commits, so
`CHANGELOG.md` is bot-owned — don't hand-edit it. What you *do* own is the **PR title**:
because we squash-merge, the title becomes the commit semantic-release reads. Write it as a
[Conventional Commit](https://www.conventionalcommits.org/) — `type(scope): summary` — and
choose the type deliberately, because it *is* the release decision:

- `fix:` → patch · `feat:` → minor · a `!` (or a `BREAKING CHANGE:` footer) → major.
- `docs:` / `ci:` / `chore:` / `refactor:` / `test:` → no release.

The `pr-title` check lints this. The summary is user-facing — it lands in the release notes —
so write it for a reader, not as a diff of files.

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

- Merge only when CI is green and any review is resolved. We **squash-merge** — that's
  what makes the PR title the single commit semantic-release reads.
- **Post-merge hygiene:** the remote branch auto-deletes on merge; prune your local one,
  and confirm `main` is green (`plumbline`, tests) after the merge.
- **Releasing is automatic.** If the PR was a `fix` / `feat` / breaking change, merging it
  triggers semantic-release — it tags, cuts the notes-only GitHub Release, updates
  `CHANGELOG.md` + `plugin.json`, and publishes the `dist` branch
  ([ADR 005](../adrs/005-automated-releases-semantic-release.md)). You never tag, run
  `make dist`, or edit the changelog by hand. Binaries never live on `main`
  ([ADR 001](../adrs/001-plugin-packaging-and-distribution.md)).

## Git conventions (repo-wide)

- Branch off `main`; never commit directly to `main`.
- End commit messages and PR descriptions with:
  `Co-Authored-By: Claude <noreply@anthropic.com>` (per the repo's authorship rule).
- Interactive rebases/adds aren't available in the agent environment — structure commits
  as you go instead.
