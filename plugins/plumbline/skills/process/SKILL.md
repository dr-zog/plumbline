---
name: process
description: Set up the working practice that keeps a Plumbline gate meaningful — register-first changes, never a red default branch, tickets as outcomes, ADRs for decisions — by writing a docs/processes/ set into the repo, adapted to its forge (GitHub or GitLab) and its own checks, plus optional thin agent skills that point at those docs. Use when someone wants a development process, workflow or team conventions to go with Plumbline, asks "how should we work with this", wants contributor or ticket/PR/MR/ADR guidelines, or after onboarding once the register is standing. Writes guidance only; never automates the forge.
---

<!-- [impl->component~process-skill~1] -->

# Plumbline — process

The gate tells you when the register and the code disagree. This skill sets up the way of
working that keeps that rare, so a red gate always means something. It writes the practice
into the repo as plain Markdown docs the team owns, in their forge's terms.

Run the engine via the plugin's launcher, `${CLAUDE_PLUGIN_ROOT}/bin/plumbline` (it picks the
right per-platform binary).

Read `references/forges.md` for forge detection and the placeholder table. The template is in
`references/template/`: four process docs, an index, and four thin pointer skills under
`pointer-skills/`. This file is the *process*.

## What ships, and what doesn't

This skill writes **guidance, not machinery** (ADR 011). It never opens tickets, merge
requests or releases, and never calls `gh`/`glab` or any forge API. The docs are plain Markdown
that any person or tool can read; only the optional pointer skills are Claude-specific. Once
written, the docs belong to the consumer: they edit them freely.

The template separates the **method** from **optional conventions**:

- **The method** is always written: register-first; never a red default branch; tickets as
  outcomes, with the undesigned backlog in the tracker; deliberate status use; ADRs for
  decisions.
- **Optional conventions** sit between `<!-- optional:NAME -->` and `<!-- /optional -->`
  markers in `raising-a-change.md`: `conventional-commits`, `squash-merge`, and
  `release-automation`. Ask about each. Keep the content (minus the markers) for each one they
  want, and delete the whole block for each one they don't. Never write a marker into the repo.

## Steps

1. **Check the ground.** The repo should have a register (the engine runs against it). If it
   doesn't, suggest `onboard` first: the practice assumes a register to lead with. Read any existing
   `docs/processes/`, `CONTRIBUTING.md`, ADR directory and `.claude/skills/` so you build on
   them, not over them.
2. **Detect and confirm the forge** (`references/forges.md`): CI files first, then the
   `origin` host as a hint, then ask.
3. **Fill the repo values.**
   - The default branch.
   - The register location.
   - The ADR directory: an existing one if present, otherwise `docs/adrs`.
   - The repo's own verify commands. Find them in its build files (`Makefile`,
     `package.json` scripts, `go.mod`, `pyproject.toml`, CI jobs) and **confirm them with the
     developer**. Don't guess at commands that don't exist.
   - How a teammate runs the gate (`{{PLUMBLINE}}`): a make target, a vendored binary, the
     pre-commit hook's path, or `plumbline` on `PATH`. Confirm it; the launcher path above
     only works inside Claude Code.
   - If the ADR directory has no index (`README.md`), offer to create a stub, since
     `writing-adrs.md` tells people to keep one.
4. **Ask about the optional conventions**, one line each, with the trade-off: Conventional
   Commit titles (machine-readable history), squash-merge (one commit per change), and release
   automation (only if they have, or want, a release pipeline). **Release automation requires
   Conventional Commits**: it reads the commit types to pick the version. If they want the
   first without the second, keep both or drop both.
5. **Write `docs/processes/`** from `references/template/`, with every placeholder filled and
   the optional blocks resolved.
6. **Offer the pointer skills.** These are `new-ticket`, `work-ticket`, `raise-change` and
   `write-adr`, from `references/template/pointer-skills/<name>.md`, each written to
   `.claude/skills/<name>/SKILL.md` with `{{CHANGE}}` filled. Each one only routes to its doc,
   so the docs stay the single source.
7. **Verify.** Search every file you wrote, docs and pointer skills alike, for leftover
   `{{`, `<!-- optional:` and `<!-- /optional` markers; fix any you find. Run the engine
   (still green: these docs carry no anchors). Then list for the developer what was written
   and what was skipped.

## Never clobber; re-runs show diffs

This mirrors `enforce`'s rule to respect any existing hook.

- **A file that already exists is never overwritten.** Show the difference between it and
  the filled template, and let the developer merge, keep theirs, or take the template's.
- **Leave unrelated files alone.** That covers other docs in `docs/processes/`, and other
  skills in `.claude/skills/`.
- **A pointer skill whose name is taken** by an unrelated skill: offer a different name, or
  skip it.
- **Re-running** on a repo this skill already set up is how a consumer pulls in later
  template changes. Each change arrives as a diff for them to accept or decline. Their own
  edits are theirs; never revert them silently.

## What this is not

It's not a forge integration or a workflow engine. It's not Plumbline's own process either:
Plumbline's repo keeps its GitHub-specific docs and its own release rules. And it's not a
reason to delay the gate. `enforce` makes the trace a build check whether or not the team adopts
these docs.
