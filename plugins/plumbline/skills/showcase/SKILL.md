---
name: showcase
description: Generate glossy, human-facing product documentation — a two-pager or a whitepaper — from a Plumbline register, as self-contained branded HTML. Use whenever someone wants a product one-pager/two-pager, solution brief, sales sheet, marketing or explanatory doc, whitepaper, pitch, or a "glossy"/"executive"/"non-technical" version of what a codebase does, generated from its register. This is the audience-facing counterpart to the dev-facing register: it turns C4 features and requirements into benefit-led copy a buyer or exec actually reads. Ships as an ADAPTABLE EXAMPLE — swap the brand tokens, logo and voice for your own.
---

<!-- [impl->component~showcase-skill~1] -->

# Plumbline — showcase

Turn the register into collateral a *human* wants to read. The dev-facing register
answers "does every requirement still have code?"; `showcase` answers "why should I
care about this product?" — the same source of truth, aimed one altitude up.

Read `references/brand-and-layout.md` for the brand-token system, the HTML skeletons
for each format, and the print styles. This file is the *process*.

> **This skill ships as an example.** Its default palette, logo slot and voice are
> deliberately generic placeholders — nobody else's brand. Adapt the brand tokens,
> drop in a real logo, and rewrite the voice guidance to the tone the product wants.
> What's portable and worth keeping is the *method*: register → benefit-led HTML.

## The one rule that makes this on-mission

**Only claim what the register supports.** The point of generating collateral from the
register is that your marketing can't silently drift from your product — the same
anti-rot promise the anchors give your code. So features, capabilities and proof
points come *from register items*, reworded for a human; never invent a capability the
register doesn't contain. Benefit-led, yes; fictional, never.

## Steps

1. **Load the register.** Prefer `plumbline -json` (structured: items, their types,
   Covers/Needs, and the scorecard). Otherwise read the register file(s) directly.
2. **Confirm audience, voice and format** with the user if not given: who reads this
   (buyer? exec? developer-evaluator?), the tone, and **two-pager** vs **whitepaper**.
   Get a logo/brand steer, or proceed with the placeholder brand and mark the swaps.
3. **Transform dev-facing → buyer-facing** (see the mapping below). This is the craft.
4. **Render one self-contained HTML file** using the chosen skeleton and the brand
   tokens from the reference. Inline all CSS; no external assets except optional Google
   Fonts (always with a fallback stack). Mark every brand/image placeholder clearly.
5. **Publish and hand over** — share the HTML as an artefact (it renders in a browser
   and prints cleanly to PDF), and tell the user exactly which placeholders to swap.

## The transformation (dev-facing → buyer-facing)

| Register | Becomes | How |
|---|---|---|
| `context` | The **problem & the promise** — the opening hook | State the pain, then the outcome the product delivers |
| `feat` | **Value propositions** | Rewrite capability → outcome: "you can X", not "component Y does Z" |
| `req` | **Proof points** under each value prop | The concrete "what you actually get" bullets |
| `container`/`component` | An optional light **"under the hood"** box | Credibility, not a how-it-works manual — one line each, or omit |
| Scorecard / metrics | **Credibility strip** | e.g. "100% traceability, enforced on every commit" |

**Rewrite example** — feature `feat~bidirectional-audit` ("Report both directions:
requirements with no covering code, and code with no requirement"):
- ❌ *dev-facing:* "Performs a bidirectional audit of anchors against the register."
- ✅ *buyer-facing:* "**Never lose a feature by accident.** The moment code that met a
  requirement disappears, you know — and the moment code drifts loose of any
  requirement, you know that too. Silent drift becomes a failed build, not a
  post-mortem."

## The two formats

- **Two-pager** — hero (name + tagline + logo slot) · the problem in two sentences ·
  3–4 value-prop cards · a one-glance "how it works" · a credibility strip · CTA +
  footer. Dense and glossy; prints to 1–2 pages.
- **Whitepaper** — cover page · executive summary · the problem (deep) · the approach ·
  capabilities (each value prop expanded with its proof points) · who it's for / use
  cases · a light "under the hood" · conclusion + CTA. Multi-page, narrative.

Keep both **skimmable**: strong headings, short paragraphs, generous whitespace, one
idea per block. A glossy doc that reads like a manual has failed.
