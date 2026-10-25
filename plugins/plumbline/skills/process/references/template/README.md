# Engineering processes

How we build this project. These docs are the **source of truth** for how work flows
from an idea to landed code. Any agent skills under `.claude/skills/` are thin pointers
to them: a skill triggers at the right moment and defers here for the detail.

The practice exists to keep the traceability gate meaningful. Plumbline tells you when the
register and the code disagree; this way of working keeps that rare, so a red gate always
means something.

Processes evolve. Edit these freely; they're ours, not a standard. (A settled decision with
reasons worth keeping goes in an ADR; a way of working lives here.)

## The processes

- [working-on-a-ticket.md](working-on-a-ticket.md): the dev loop, from requirement to landed
  code, register-first.
- [raising-a-change.md](raising-a-change.md): landing a change via a {{CHANGE}}.
- [writing-tickets.md](writing-tickets.md): where tickets live and how much they carry.
- [writing-adrs.md](writing-adrs.md): when and how to record a decision.
