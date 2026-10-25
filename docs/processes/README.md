# Engineering processes

How we build **Plumbline itself**. These docs are the **source of truth**; the
dev-only skills in `.claude/skills/` are thin pointers to them (DRY) — a skill
triggers at the right moment and defers here for the detail.

**These govern *developing* Plumbline; they are not part of the shipped plugin.**
Product skills live in `plugins/plumbline/skills/` and ship to users. Dev-process
skills live in `.claude/skills/` and never leave the repo. Don't cross the streams.

**The method also ships, as a template** ([ADR 011](../adrs/011-ship-the-working-practice.md)).
The `process` skill writes a forge-neutral version into consumers' repos, from
`plugins/plumbline/skills/process/references/template/`. These docs stay Plumbline's: the
GitHub instance, plus Plumbline-only rules. **Keep them in step:** a change to the *method*
here (register-first, red-branch policy, tickets, status use, ADRs) updates the template in the
same PR, or files a ticket to. Plumbline-only rules (semantic-release, `shipped-type`, the
`plugins/` release rule) stay here.

Processes evolve — edit these freely; they're not ADRs. (A settled, provenance-worthy
*decision* goes in `docs/adrs/`; a *practice* lives here.)

## The processes

- [working-on-a-ticket.md](working-on-a-ticket.md) — the dev loop: requirement → design →
  code → test, register-first.
- [raising-a-pr.md](raising-a-pr.md) — landing a change via a GitHub PR.
- [writing-tickets.md](writing-tickets.md) — where tickets live and how much they carry.
- [writing-adrs.md](writing-adrs.md) — when and how to record a decision.
