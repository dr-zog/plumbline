# 011. Ship the working practice as a forge-neutral template

- **Status:** Proposed
- **Date:** 2026-10-25
- **Deciders:** Dr. Zog (operator), analysis by Claude

## Context

Plumbline's own development follows a written practice in `docs/processes/`:
- register-first within a change
- never a red `main`
- tickets as outcomes, not designs, with the undesigned backlog in the tracker
- status used deliberately (ADR 004)
- ADRs for decisions, process docs for practice
- squash merges with Conventional Commit titles

Thin dev skills in `.claude/skills/` (`new-ticket`, `work-ticket`, `raise-pr`, `write-adr`)
point Claude at those docs at the right moment.

None of it ships. `docs/processes/README.md` is explicit that it governs developing
Plumbline itself ("don't cross the streams"), and that was right for those files. But an
external review (#56), written while adopting Plumbline on a new project, argued that this
practice is Plumbline's most transferable asset. The gate tells you when trace and code
disagree. The practice is what keeps that rare and meaningful. A consumer who installs the
plugin gets the tool without the way of working that makes it effective.

The docs are also GitHub-specific throughout: GitHub Issues, PRs, `gh pr create`, GitHub
Releases. That is accurate for Plumbline, which lives on GitHub. But the reviewer's projects
live on a self-hosted GitLab (issues, merge requests, `glab`), so they would have to translate
every reference by hand. The product is already neutral about forges elsewhere: `enforce`
picks a GitLab or GitHub job from the CI files present (`.gitlab-ci.yml`,
`.github/workflows/`) and asks when there are none.

ADR 002 is titled "traceability only". It declined building or consuming emergent-structure
*analysis*. Shipping a practice is not analysis, but it does widen what Plumbline delivers, so
the boundary needs recording.

## Decision

**Plumbline ships its working practice as a forge-neutral template, delivered by a new
`process` skill.**

- **What ships is guidance, not machinery.** The skill writes a `docs/processes/` set into the
  consumer's repo. There's no ticket, MR/PR or release automation, and no forge API integration.
  `gh`/`glab` are named in the docs as the forge's CLI, not invoked or required, so the
  zero-dependency promise is untouched. The consumer owns the written docs afterwards.
- **Where it lives.** The template is under `plugins/plumbline/skills/process/references/`: the
  process docs, the forge mapping table, and the thin agent pointers.
- **Forge-neutral source, concrete output.** The template says "issue tracker", "merge
  request" and "the forge's CLI", and carries a mapping table:
  - GitHub → Issues, pull request, `gh`, Actions
  - GitLab → issues, merge request, `glab`, GitLab CI

  Forge detection, in order:
  1. CI files, as `enforce` does.
  2. The `origin` remote's host. A self-hosted host may not name its forge, so treat this only
     as a hint.
  3. Ask.

  The developer confirms the forge, and the skill writes in the concrete terms. If the forge is
  unknown, the neutral wording stays.
- **What is parametrised.** The skill fills in:
  - forge terms
  - the default branch name
  - the register location
  - the repo's own verify commands (tests, lint, `plumbline`), found from its build files and
    confirmed with the developer

  Everything else stays as written.
- **Method versus Plumbline's own choices.** The template separates the two:
  - The **method** is register-first, never a red default branch, tickets as outcomes with the
    undesigned backlog in the tracker, deliberate status use (designed-ahead `proposed` items
    count against the spec-debt budget, which fails the gate only if configured), and ADRs for
    decisions.
  - **Optional conventions** are offered and clearly marked: Conventional Commit titles, squash
    merges, and release automation. Semantic-release is Plumbline's choice (ADR 005), not part
    of the method.
- **Thin agent pointers, offered.** The skill offers to write matching thin skills into the
  consumer's `.claude/skills/` (new ticket, work a ticket, raise an MR/PR, write an ADR). Each
  only points at the written docs, so the docs stay the single source. The docs are plain
  Markdown and serve any tooling or human; only these pointers are Claude-specific.
- **Never clobber; re-runs diff.** This mirrors `enforce`'s rule to respect any existing hook.
  - If `docs/processes/` already exists, the skill reads it first. A file it would write that
    is already present is never overwritten. The skill shows the difference and lets the
    developer merge or skip it.
  - Unrelated files in that directory are left alone.
  - The same applies to existing `.claude/skills/` entries with the same names.
  - Re-running on a repo it already set up is how a consumer pulls in later template changes,
    one diff at a time.
- **Entry points.** `onboard` offers the `process` skill at its end, once the register is
  standing. It can also be invoked directly on an already-onboarded repo.
- **Plumbline's own `docs/processes/` stays, and is kept in step deliberately.** It's the
  GitHub instance of the practice, plus Plumbline-only rules (e.g. the `shipped-type` check and
  the `plugins/` release rule). The mechanism is a rule in `docs/processes/README.md` and the
  `raise-pr` checklist: a change to the *method* in `docs/processes/` updates the template in
  the same PR, or files a ticket to. `docs/processes/README.md` is amended to say the method
  now also ships, as a template, while the docs themselves stay Plumbline's.
- **Registered like any other skill.** `component~process-skill~1` under `container~plugin~1`,
  covering a new requirement under `feat~skills-workflow~1`, with an anchor in `SKILL.md`.

On scope: this reads ADR 002 as constraining *analysis* (what Plumbline computes and
integrates), not every deliverable. `showcase`, which renders collateral rather than analysing
anything, is the existing precedent. ADR 002 is not superseded. Plumbline still builds no new
analysis. It ships the practice that gives its one analysis, traceability, its meaning.

## Consequences

**Positive**
- A consumer gets the way of working as well as the tool, on either forge, without hand
  translation, and as plain docs any tooling can read.
- The practice is phrased once in neutral terms, which makes explicit which parts are the
  method and which are local taste. That clarifies Plumbline's own docs too.
- It follows the existing precedent of being neutral about the forge (`enforce`).

**Negative / accepted**
- A second copy of the practice to keep in step with Plumbline's own docs, held by a process
  rule rather than a check. It's mitigated by keeping the template short and limited to the
  method.
- The template lives under `plugins/plumbline/**`, so `shipped-type` makes every template edit
  a `fix:`/`feat:` release. That's correct, because it changes what consumers get, but method
  edits now cut releases.
- A sixth skill, and a larger surface for consumers to read.
- Written docs diverge in consumer repos once edited. That's intended: they're the consumer's.

## Alternatives considered

- **Keep it Plumbline-internal (status quo).** Rejected: it withholds the part that makes the
  gate effective, and leaves every GitLab adopter to translate by hand.
- **Fold it into `onboard`.** Rejected: onboarding is already the longest skill, and repos that
  are already onboarded would have no way in. Instead, `onboard` *offers* the skill.
- **Ship Plumbline's GitHub docs verbatim.** Rejected: GitLab consumers would translate by hand,
  and Plumbline-only rules (semantic-release, `shipped-type`) would read as part of the method.
- **Generate Plumbline's own `docs/processes/` from the template** (the GitHub instance plus a
  Plumbline-only addendum), so there is one copy. Deferred: it would make Plumbline's richest
  docs indirect to edit before the template has proved itself. Revisit if the two copies drift
  in practice.
- **Automate the workflow** (create tickets, open MRs, cut releases through forge APIs).
  Rejected: machinery and per-forge coupling. That's a different product, and it's what ADR 002's
  spine warns against.

## Provenance

- Issue #56, from an external review (2026-10-25) written while adopting Plumbline on a new,
  GitLab-hosted project.
- `docs/processes/` and `.claude/skills/`: the practice as it stands.
- `plugins/plumbline/skills/enforce/SKILL.md`: CI-file forge detection and the "respect any
  existing hook" rule, both mirrored here.
- ADR 002 (scope, read here as constraining analysis), ADR 004 (status lifecycle and
  spec-debt), ADR 005 (semantic-release as Plumbline's choice); #54 (backlog guidance
  reconciled before this).
