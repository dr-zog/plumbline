# Contributing

Plumbline is built AI-first. Dr. Zog writes it with Claude Code and a set of skills that encode every procedure here — raising a ticket, working a ticket, opening a PR, recording a decision. The skills are how *she* works; they are not a toll gate.

**You don't have to use Claude Code.** Bring whatever agent you like, or hand-write every line like it's 2019 — Plumbline doesn't measure your method.

**You do have to follow the processes.** They aren't suggestions, and they aren't enforced by taste: the build enforces them, and a PR that skips them doesn't merge. That's the whole deal — the tooling is yours to choose, the process is not.

## The processes are the law

They live in [`docs/processes/`](docs/processes/) and are the single source of truth:

- [writing-tickets.md](docs/processes/writing-tickets.md) — a ticket is a *requirement*, not a design.
- [working-on-a-ticket.md](docs/processes/working-on-a-ticket.md) — the register-first loop.
- [raising-a-pr.md](docs/processes/raising-a-pr.md) — how a change lands on `main`.
- [writing-adrs.md](docs/processes/writing-adrs.md) — when a decision gets recorded.

## Ground rules

- **PR titles are [Conventional Commits](https://www.conventionalcommits.org/).** We squash-merge, the build lints the title, and each release — version *and* changelog — is computed from it; a `fix` / `feat` / breaking change ships on merge ([ADR 005](docs/adrs/005-automated-releases-semantic-release.md)).
- **Be excellent to each other.** Conduct is covered by the [Code of Conduct](CODE_OF_CONDUCT.md) — the short version: the machines are held to a standard; the humans are asked only to be civil.
- By contributing, you agree your work ships under the repo's [MIT licence](LICENSE).

Questions belong in [Discussions](https://github.com/dr-zog/plumbline/discussions). Now go and instruct an agent.
