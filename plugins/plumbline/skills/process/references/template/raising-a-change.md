# Raising a {{CHANGE}}

How a change lands on `{{DEFAULT_BRANCH}}`: a **{{CHANGE}}**, even solo. The {{CHANGE}} is the
review checkpoint, the CI gate, and the paper trail.

## Before you open it

Everything green on the branch:

{{VERIFY}}
- {{PLUMBLINE}}: trace green, **0 gaps** (the register-first work is now covered).
- Commits are clean and logically grouped. Each message says *what* and *why*, not a list of
  files.
- Up to date with the latest `{{DEFAULT_BRANCH}}`.

<!-- optional:conventional-commits -->
## Conventional Commit titles

Write the {{CHANGE}} title as a [Conventional Commit](https://www.conventionalcommits.org/),
`type(scope): summary`. With squash-merge it becomes the one commit on `{{DEFAULT_BRANCH}}`;
without it, write every commit message the same way. Choose the type deliberately:

- `fix:` is a bug fix. `feat:` is a new capability. A `!`, or a `BREAKING CHANGE:` footer, is a
  breaking change.
- `docs:` / `ci:` / `chore:` / `refactor:` / `test:` change no behaviour.

The summary is user-facing, so write it for a reader, not as a diff of files.
<!-- /optional -->

## Open the {{CHANGE}}

Use {{CLI}} or the web UI. Structure the description:

- **What**: the outcome delivered, in a sentence.
- **Why**: the ticket or requirement it satisfies (link it).
- **How**: the shape of the change. Call out any trade-off honestly.
- **Proof**: the checks that pass, including the green Plumbline scorecard.

Keep it skimmable. A reviewer should grasp the change from the description without reading
the diff first.

## Merge and clean up

- Merge only when CI is green and any review is resolved.
<!-- optional:squash-merge -->
- **Squash-merge**, so `{{DEFAULT_BRANCH}}` carries one commit per change.
<!-- /optional -->
- Delete the branch, prune your local copy, and confirm `{{DEFAULT_BRANCH}}` is green after
  the merge.
<!-- optional:release-automation -->
- **Releasing is automatic.** Merging a releasable change triggers the release pipeline in
  {{CI}}, which computes the version from the commit titles, tags, and writes the changelog.
  You never tag or edit the changelog by hand.
<!-- /optional -->
