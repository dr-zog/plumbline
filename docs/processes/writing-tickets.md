# Writing tickets

Where Plumbline's work is captured, and how much a ticket carries.

## Where

**GitHub Issues** — on the `github.com/dr-zog/plumbline` repo. GitHub is the backlog: the
"decided but not yet started" list.

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
The backlog lives in GitHub Issues; the *active, code-backed* requirements live in the register.
Don't pre-populate the register with unstarted tickets — that's what GitHub Issues are for.

## Granularity

One shippable outcome per ticket. If a ticket can't become a small, coherent set of
register items with covering code in one PR, split it.
