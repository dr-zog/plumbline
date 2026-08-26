---
name: enforce
description: Install Plumbline's traceability gate so unanchored code can't land — offer the developer a pre-commit hook, a CI job, both, or neither, and set up whichever they choose. Use this whenever someone wants to enforce, gate or require traceability, add a pre-commit or CI check for Plumbline, "make the build fail on broken anchors", wire the engine into their pipeline, or asks how to stop unanchored code being committed. Reach for it once a project is onboarded and the team wants the discipline enforced automatically rather than by hand.
---

<!-- [impl->component~enforce-skill~1] -->

# Plumbline — enforce

The engine already exits non-zero on any gap; enforcement is just deciding *where*
that exit code is allowed to stop the flow. The insistence and the light touch work
together — the skills add the anchors so a developer rarely hand-writes one, and the
gate makes unanchored code fail loudly rather than land quietly.

Bundled snippets you'll copy from: `references/pre-commit-hook.sh` and
`references/ci-snippets.md`.

## Before you start

- **A register must exist** and the engine must run cleanly enough that turning on a
  gate won't immediately block all work. If the project isn't onboarded, do that
  first (`onboard`); gating an empty register just produces noise.

## Let the developer choose the gate

Enforcement is opt-in. Present the options plainly and install what they pick — don't
assume:

- **Pre-commit** — fast, local feedback at commit time. Bypassable with
  `--no-verify`, and only present on machines that installed the hook.
- **CI** — the authoritative backstop: unbypassable, but feedback comes after push.
- **Both** — belt and braces, and the usual recommendation for a team that's serious
  about it: quick local signal, hard remote gate.
- **Neither** — just show them the one-line command (`plumbline`) to run by hand.

Ask which they want before touching anything.

## Pre-commit hook

Install `references/pre-commit-hook.sh` as `.git/hooks/pre-commit` (make it
executable). It's a plain POSIX script — zero dependencies, matching Plumbline's
promise. It runs the plugin's engine **launcher**, which picks the right
per-platform binary, and blocks the commit on failure.

**Bake in the launcher path.** A git hook runs *outside* Claude Code, so
`${CLAUDE_PLUGIN_ROOT}` isn't set there. Resolve `${CLAUDE_PLUGIN_ROOT}/bin/plumbline`
to an absolute path at install time and substitute it for the hook's `__PLUMBLINE__`
placeholder before writing the hook — e.g. `sed "s|__PLUMBLINE__|$ABS_LAUNCHER|"`.
Working from the engine's source tree instead of an installed plugin? Point it at a
`go build`-produced binary, or the repo's `plugins/plumbline/bin/plumbline` launcher.

**Respect any existing hook.** If `.git/hooks/pre-commit` already exists, don't
clobber it — read it, and either append the Plumbline call or tell the developer how
to combine them. Losing their existing checks would be a nasty surprise.

Note the `--no-verify` escape hatch exists by design (a hook can't be the whole
story), which is exactly why CI is worth having too.

## CI job

Detect the CI system from the repo and add a gate job from `references/ci-snippets.md`:

- **GitLab** (`.gitlab-ci.yml` present) — add the `plumbline` job.
- **GitHub Actions** (`.github/workflows/` present) — add the workflow.
- Neither present — ask which they use, or hand them the generic "run the binary,
  fail on non-zero" recipe.

The one real wrinkle is **getting the binary into CI**: the engine ships inside the
plugin, not the project, so the CI job must obtain it — vendor the committed binary
into the repo, download it from a release, or `go build` it from source. The snippets
show each; pick the one that matches how the team distributes tools, and say plainly
which you chose.

## Threshold, if they want a softer gate

By default the gate is strict — any gap fails. If the team wants to adopt gradually,
offer `-min-coverage N`: the gate then fails only below `N%` coverage (though broken
anchors and structural errors always fail). A team can ratchet `N` up over time.

## Verify it actually gates

Don't leave without proving it works. Run the hook or job once against the current
tree and confirm the exit code does what's intended — a passing tree returns 0, and a
deliberately broken anchor returns non-zero and blocks. An installed gate that
doesn't gate is worse than none, because it's trusted.
