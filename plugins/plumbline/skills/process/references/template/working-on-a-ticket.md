# Working on a ticket

The development loop, from a requirement to landed code. **The register leads the code
within a change.**

## The core principle: register-first

For any change that adds or alters a *feature*, the register requirement goes in
**first**, and the code makes it green **before merge**:

1. **Register the requirement first.** Add or adjust the `feat`/`req`/`component` items in
   `{{REGISTER}}` **with `Status: approved`**, or promote the existing `proposed` item if it
   was designed ahead. Run {{PLUMBLINE}}: the new item is **uncovered** and the gate is
   **red**. That red is correct. It's the promise, not yet kept.
2. **Implement**, then **anchor** the code-area(s) that satisfy it.
3. **Go green.** {{PLUMBLINE}} reports 0 gaps before you open the {{CHANGE}}.

**Never merge a red `{{DEFAULT_BRANCH}}`.** The red state lives only on your branch. A red
`{{DEFAULT_BRANCH}}` must always mean a real regression, never "work in progress". That's
what keeps the gate trustworthy. (For a *trivial* change, such as a typo, a doc tweak, or a
one-line fix with no new requirement, registering "first" is ceremony; just keep the trace
green.)

The *undesigned* backlog lives in **{{TRACKER}}**, not in the register. Designed-ahead work
may sit in the register as `proposed`. It's tracked but not gated per item, and it counts
toward the spec-debt budget, which fails the gate once exceeded if a budget is configured.
What must never sit there is **persistent red**: an `approved` requirement enters when its
work *starts*, and is green by merge.

## The loop

1. **Pick up the ticket** (see [writing-tickets.md](writing-tickets.md)). Make sure you
   understand the *outcome* wanted; clarify with a human if it's a judgement call, not a fact.
2. **Branch** off `{{DEFAULT_BRANCH}}`.
3. **Design, lightly.** Decide the shape of the change: which feature or requirement, which
   component or code-area. If the decision is architectural and would be re-litigated later,
   write an [ADR](writing-adrs.md).
4. **Register-first** (above).
5. **Implement.** Match the surrounding code's idiom and comment density.
6. **Anchor** the code-areas: `impl` on implementation, `utest`/`itest` on tests.
7. **Verify.** All of:
{{VERIFY}}
   - {{PLUMBLINE}}: trace green, with 0 uncovered/broken/orphan/transitive/structural.
8. **Land it** via a {{CHANGE}}. See [raising-a-change.md](raising-a-change.md).

## Notes

- Anchors are coarse: one per code-area (package or directory), not per file. Fence any
  example tags you write in docs or tests with `oft:off` … `oft:on` comments, so the engine
  doesn't read them as real anchors.
- Prose comments are up to you. Plumbline never measures them; the anchor is the only
  non-negotiable.
