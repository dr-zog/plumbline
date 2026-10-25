# Forges: detection and terms

The template is written with `{{…}}` placeholders. Fill them from this table once the forge
is confirmed.

## Placeholders

| Placeholder | GitHub | GitLab | Unknown / other |
|---|---|---|---|
| `{{TRACKER}}` | GitHub Issues | GitLab issues | the issue tracker |
| `{{CHANGE}}` | pull request | merge request | merge request |
| `{{CLI}}` | `` `gh` `` | `` `glab` `` | your forge's CLI |
| `{{CI}}` | GitHub Actions | GitLab CI | CI |

These are filled per repo, not per forge:

| Placeholder | Value |
|---|---|
| `{{DEFAULT_BRANCH}}` | the repo's default branch, e.g. `main` (`git symbolic-ref --short refs/remotes/origin/HEAD`, minus `origin/`) |
| `{{REGISTER}}` | the register's location: `register.md`, a `register/` directory, or the configured path |
| `{{ADR_DIR}}` | where ADRs live, e.g. `docs/adrs` (use an existing ADR directory if the repo has one) |
| `{{PLUMBLINE}}` | how this repo runs the gate, in backticks, e.g. `` `make trace` ``, `` `./tools/plumbline` `` or `` `plumbline` `` (confirmed with the developer, so a teammate can run it) |
| `{{VERIFY}}` | the repo's own checks, one bullet per line, e.g. ``- `npm test` · `npm run lint` ``. Indent each line to match its list: **3 spaces** in `working-on-a-ticket.md` (inside step 7), **none** in `raising-a-change.md` |

In the pointer skills, `{{CHANGE}}` is the same term as in the docs.

## Detecting the forge

Check in this order, and **confirm with the developer** before writing anything:

1. **CI files**, as `enforce` does: `.gitlab-ci.yml` means GitLab; `.github/workflows/`
   means GitHub.
2. **The `origin` remote's host** (`git remote get-url origin`): `github.com` means GitHub,
   and `gitlab.com` means GitLab. A self-hosted forge often doesn't name itself in its host,
   so treat any other host as a hint at most.
3. **Ask.** If the developer uses another forge, or would rather not say, use the
   "Unknown / other" column. The neutral wording reads fine as it is.

`gh` and `glab` are only named in the docs as the forge's CLI. The skill never runs them, and
the consumer doesn't need either installed.
