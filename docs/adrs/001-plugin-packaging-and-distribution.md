# 001. Plugin packaging & distribution

- **Status:** Accepted
- **Date:** 2026-08-25
- **Deciders:** Dr Zog (operator), analysis by Claude

## Context

Plumbline ships as a Claude Code **plugin** (skills) plus a zero-dependency Go
**engine** binary. We need to deliver both to Claude Code users across Linux, macOS
and Windows, with **no toolchain to install** and **no runtime fetch** (the binary
must be bundled in the plugin, self-contained on install) — while keeping the git
`main` branch free of binary-history bloat.

Scope note: we care about **Claude Code only** — not Claude Desktop, which removes the
`.mcpb` bundle path and its complications entirely.

Forces:
- Committing compiled binaries to a mainline branch bloats every clone's history
  forever.
- A Claude Code plugin *can* bundle binaries in `bin/`, and hooks/skills can locate
  them via the exported `${CLAUDE_PLUGIN_ROOT}` env var.
- The reference pattern (`distributing-go-binary-mcp.md`, repo root) recommends a
  **four-plugin** shape — three per-OS engine plugins + one shared skills plugin with
  dependencies — but that shape is driven by an **MCP-server** constraint (see below).

## Decision

1. **One plugin, not per-OS plugins.** Plumbline's engine is a **CLI the skills and
   the git hook invoke**, *not* an MCP server declared in a plugin manifest. So the
   correct binary is chosen **at runtime** (`uname` → `${CLAUDE_PLUGIN_ROOT}/bin/plumbline-<os>-<arch>`).
   A single plugin therefore carries the skills **and** all binaries. The reference
   doc's four-plugin pattern solves a manifest-pins-one-binary-per-OS problem that
   only exists for MCP servers; it does not apply to us.

2. **Platform matrix — four artifacts:**
   - `plumbline-linux-amd64`
   - `plumbline-linux-arm64`
   - `plumbline-darwin-universal` — a **fat Mach-O** fusing `amd64` + `arm64`, assembled
     on the Linux CI runner (pure-Go fuse, no macOS host / `lipo`).
   - `plumbline-windows-amd64.exe` (Windows-on-ARM runs amd64 under emulation).

   **macOS gets a universal binary; Linux does not** — the fat-binary format is a
   Mach-O feature; ELF has no equivalent, so Linux ships one binary per architecture.

3. **Binary bundled in the plugin, not fetched at runtime.** Self-contained install;
   no network on first use.

4. **Binaries kept off `main`; served from a force-pushed orphan `dist` branch.**
   `main` stays source-only (`bin/` gitignored). On a release tag, CI builds the four
   binaries and force-pushes a **parentless (orphan) commit** to `dist` carrying the
   plugin + binaries — so `dist` is always a single flat commit and never accumulates
   binary history. (`git checkout --orphan dist && commit && push --force`.)

5. **The marketplace entry pins the commit `sha`.** The `dr-zog` marketplace's
   plumbline entry uses a `git-subdir` source: `url` (the repo) + `path`
   (`plugins/plumbline`) + `ref` (`dist`) + `sha` (the current orphan commit). CI bumps
   the `sha` each release. Pinning `sha` makes updates deterministic and sidesteps two
   behaviours the Claude Code docs do **not** specify — clone depth and force-push
   handling — by fetching an exact commit rather than a moving branch tip.

6. **Public GitHub is the host.** Point the marketplace at the public repo directly;
   no self-hosted mirror or website (the reference doc's private/behind-firewall route
   is more than we need). The repo is published publicly on `github.com/dr-zog`.

7. **Build flags:** `CGO_ENABLED=0 go build -trimpath -ldflags "-s -w -X main.version=$V"`
   — static, path-stripped, symbol-stripped, with the version stamped in.

## Consequences

**Positive:** self-contained install (no runtime fetch); `main` history stays clean;
no binary bloat anywhere (a clone only transfers reachable objects, and the orphan
branch keeps even a full clone of `dist` flat); deterministic, force-push-safe updates;
one plugin to build, version and publish; no dependency graph.

**Negative / accepted trade-offs:**
- Every user downloads all four binaries (~13 MB) rather than just their platform's —
  negligible at this size.
- The release flow needs real CI plumbing: orphan force-push + a mac fat-binary fuse +
  a marketplace `sha` bump per release.
- The installed **git pre-commit hook runs outside Claude Code**, where
  `${CLAUDE_PLUGIN_ROOT}` is unset — so the hook must have the binary path **baked in at
  install time** by the `enforce` skill (which knows the plugin root when it runs).
- Binaries cannot ride the claude.ai **org-managed** distribution channel (that route is
  skills-only) — irrelevant to the public-marketplace path we're taking.

**Neutral / to watch:** clone-depth and force-push behaviour are undocumented in Claude
Code; the `sha` pin neutralises both today. Revisit if the plugin machinery changes.

## Alternatives considered

- **Fetch the binary at runtime (GitHub release assets).** Rejected: reintroduces a
  network dependency on first run; the operator wants a self-contained plugin. (Release
  assets remain a fine fallback — 2 GiB/file, unlimited total & bandwidth on public
  repos — if the bundled-binary route ever proves awkward.)
- **Commit binaries to `main`.** Rejected: permanent, unbounded history bloat.
- **Git LFS.** Rejected for now: adds a forge + client dependency and uncertain
  marketplace/`git-subdir` LFS-awareness; only worth it if fully-offline in-repo binaries
  become a hard requirement.
- **The four-plugin per-OS structure (reference doc).** Rejected: it solves an
  MCP-server manifest constraint we don't have; unnecessary machinery for a CLI. (Kept in
  reserve should Plumbline ever add a real MCP server.)
- **Self-hosted git mirror / website (private route).** Rejected: public GitHub suffices.

## Provenance

- **`distributing-go-binary-mcp.md`** (repo root) — the orphan-ref force-push trick, the
  static build flags, the macOS fat-binary-on-Linux fuse, and the four-plugin rationale
  (which we consciously do *not* adopt, for the reason in Decision 1).
- **Claude Code plugin/marketplace docs**, verified 2026-08-25: `git-subdir` sources take
  `url` + `path` + `ref` + optional `sha`; a marketplace can be added from a non-default
  branch (`#ref`); `${CLAUDE_PLUGIN_ROOT}` is exported to hooks/subprocesses; `bin/`
  binaries are supported (but not via the org-managed channel); clone depth and force-push
  handling are **not** documented.
- **GitHub release limits** (context for the rejected release-asset option): 2 GiB/file,
  no total-size or bandwidth cap on public repos, assets excluded from repo size.
- Session discussion, 2026-08-25.

## Follow-ups (implementation, not decided here)

- Build matrix + macOS fat-binary fuse in the Makefile/CI.
- A `uname`-based runtime binary selector used by the skills and the installed git hook.
- The `${CLAUDE_PLUGIN_ROOT}` fix in the `enforce` hook and skills.
- The `dist`-branch orphan-force-push release job + the `dr-zog` marketplace entry.
- The public-GitHub home is live; automate the `dist` release + marketplace entry next.
