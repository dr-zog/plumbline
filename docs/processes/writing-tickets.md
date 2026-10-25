# Writing tickets

Where Plumbline's work is captured, and how much a ticket carries.

## Where

**GitHub Issues** — on the `github.com/dr-zog/plumbline` repo. GitHub is the backlog: the
"wanted but not yet designed or started" list.

## What a ticket carries

A ticket is a **requirement, not a design**. It states the *outcome* someone wants and
*why* — the altitude Plumbline itself works at. It does **not** prescribe the
implementation; that's decided on the branch (design → code), and architectural calls
get an [ADR](writing-adrs.md).

A good ticket has:

- **Title** — the outcome, as a capability ("Gate on a coverage threshold", not
  "Add a --min-coverage flag").
- **Why** — the problem it solves / the pain today.
- **Acceptance** — how we'll know it's done, in outcome terms (ideally the shape of the
  register requirement it will become).
- Just enough context to pick it up cold. No implementation plan, no file list.

## The link to the register

A ticket is a *candidate* requirement. When work **starts**, it becomes a real
`req~…`/`feat~…` in `register/` (register-first — see
[working-on-a-ticket.md](working-on-a-ticket.md)), and the ticket references that ID.
The *undesigned* backlog lives in GitHub Issues: don't transcribe tickets into the register
as placeholders. The one exception is **designed-ahead** work. Once a design is settled ahead
of building, it may enter the register as `Status: proposed`. That item is never gated on its
own: it can't go red for being unbuilt. It does count toward the spec-debt budget, which, if
configured, fails the gate once un-built spec exceeds it
([ADR 004](../adrs/004-requirement-status-lifecycle.md)). Everything else enters when its work
starts.

## Granularity

One shippable outcome per ticket. If a ticket can't become a small, coherent set of
register items with covering code in one PR, split it.
