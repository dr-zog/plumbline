# Brand & layout reference

The reusable foundation for `showcase` output: a brand-token system, a default
(placeholder) palette and voice, the HTML skeletons for both formats, and print
styles. Everything is inline and self-contained — one HTML file, no build step,
no external assets except Google Fonts (always with a fallback stack).

## Contents
- [Brand tokens — the swap layer](#brand-tokens--the-swap-layer)
- [Voice (default) — the other swap layer](#voice-default--the-other-swap-layer)
- [Placeholders](#placeholders)
- [The CSS foundation](#the-css-foundation)
- [Two-pager skeleton](#two-pager-skeleton)
- [Whitepaper skeleton](#whitepaper-skeleton)
- [Print / PDF](#print--pdf)
- [Do / don't](#do--dont)

## Brand tokens — the swap layer

All brand identity lives in CSS custom properties on `:root`. **This is what a user
swaps** — change these six-ish values and the whole document re-skins. The defaults
are a deliberately neutral placeholder brand, not anyone's real identity.

```css
:root {
  --ink:     #16181d;   /* primary text */
  --muted:   #5b6472;   /* secondary text */
  --bg:      #ffffff;   /* page background */
  --surface: #f6f7f9;   /* cards, callouts */
  --line:    #e6e8ec;   /* hairlines */
  --brand:   #1b2a4a;   /* SWAP — primary brand: bars, headings, cover */
  --accent:  #c9812c;   /* SWAP — accent: rules, highlights, CTA */
  --font-display: "Fraunces", Georgia, "Times New Roman", serif;   /* SWAP — headlines */
  --font-body:    "Inter", system-ui, -apple-system, sans-serif;    /* SWAP — body */
}
```

Fraunces (serif display) + Inter (clean sans body) is a tasteful glossy default; both
are on Google Fonts. Load with a fallback and never let text depend on the webfont:

```html
<link rel="preconnect" href="https://fonts.googleapis.com">
<link href="https://fonts.googleapis.com/css2?family=Fraunces:opsz,wght@9..144,500;9..144,600&family=Inter:wght@400;500;600&display=swap" rel="stylesheet">
```

## Voice (default) — the other swap layer

The default voice is **confident, plain-spoken, benefit-led, mechanism-honest — no
hype words** ("revolutionary", "seamless", "cutting-edge" are banned by default).
Short sentences. One idea per block. This is a placeholder too: rewrite it to the
product's tone (playful, enterprise-formal, technical-credible, whatever fits). Tell
the user you've used the default voice and that it's theirs to change.

## Placeholders

Mark every swap point with an HTML comment so a user can find them with a search:

- Logo: `<!-- LOGO --><div class="logo">Your&nbsp;Logo</div>` — a text stand-in; a user
  drops in an `<img>` or inline SVG.
- Hero / section imagery: `<!-- IMAGE: hero -->` with a captioned neutral block, so the
  layout holds even before real art arrives.
- CTA target, contact, dates, company name: `<!-- SWAP: ... -->`.

Never ship a fake logo or invented brand as if it were real — placeholders must *read*
as placeholders.

## The CSS foundation

Include this once, inline, in every document (both formats build on it):

```css
* { box-sizing: border-box; }
html { -webkit-print-color-adjust: exact; print-color-adjust: exact; }
body {
  margin: 0; color: var(--ink); background: var(--bg);
  font-family: var(--font-body); font-size: 17px; line-height: 1.6;
  font-feature-settings: "kern","liga";
}
h1,h2,h3 { font-family: var(--font-display); line-height: 1.12; font-weight: 600; margin: 0 0 .4em; }
h1 { font-size: clamp(2rem, 4vw, 3.2rem); letter-spacing: -0.01em; }
h2 { font-size: 1.6rem; margin-top: 1.6em; }
h3 { font-size: 1.15rem; }
p { margin: 0 0 1em; max-width: 68ch; }
a { color: var(--accent); }
.eyebrow { font: 600 .78rem/1 var(--font-body); letter-spacing: .14em; text-transform: uppercase; color: var(--accent); }
.lede { font-size: 1.25rem; color: var(--muted); }
.wrap { max-width: 960px; margin: 0 auto; padding: 0 32px; }
.rule { height: 3px; width: 56px; background: var(--accent); border: 0; margin: 1.2em 0; }

/* value-prop cards */
.cards { display: grid; grid-template-columns: repeat(auto-fit, minmax(240px,1fr)); gap: 20px; }
.card { background: var(--surface); border: 1px solid var(--line); border-radius: 14px; padding: 24px; }
.card h3 { color: var(--brand); }
.card .proof { color: var(--muted); font-size: .95rem; margin: 0; }

/* callouts, credibility strip, CTA */
.callout { border-left: 3px solid var(--accent); background: var(--surface); padding: 16px 20px; border-radius: 0 10px 10px 0; }
.strip { background: var(--brand); color: #fff; border-radius: 16px; padding: 28px 32px; display: flex; flex-wrap: wrap; gap: 28px; }
.stat { font-family: var(--font-display); font-size: 2rem; }
.cta { background: var(--accent); color: #fff; text-decoration: none; padding: 12px 22px; border-radius: 999px; font-weight: 600; display: inline-block; }
.logo { font-family: var(--font-display); font-weight: 600; color: var(--brand); }
footer { color: var(--muted); font-size: .85rem; border-top: 1px solid var(--line); margin-top: 3rem; padding: 1.5rem 0; }
img { max-width: 100%; height: auto; }
```

## Two-pager skeleton

One flowing document, ~2 printed pages. Order:

1. **Hero band** — brand bar with `.logo`, product name (`h1`), a one-line tagline
   (from `context`), optional `<!-- IMAGE: hero -->`.
2. **The problem** — two sentences, `.lede`.
3. **Value props** — `.cards` grid of 3–4 cards, each: benefit headline (`h3`) + one
   sentence + a `.proof` line drawn from the requirement(s).
4. **How it works, at a glance** — 3 short steps or a tiny inline diagram; light.
5. **Credibility `.strip`** — 2–3 `.stat`s (from the scorecard / real numbers) or a
   quote placeholder.
6. **CTA + footer** — `.cta` button, contact/logo placeholders.

## Whitepaper skeleton

Multi-page narrative with a cover. Order:

1. **Cover page** — full-bleed `--brand` panel: logo, title (`h1`), subtitle, date +
   author placeholders. `page-break-after: always`.
2. **Executive summary** — 3–4 sentences: problem → approach → payoff.
3. **The problem** — the deep version of `context`, with the stakes.
4. **The approach** — the product's principle/mechanism, lightly (not a manual).
5. **Capabilities** — each `feat` as an `h2` section: the value, then its `req` proof
   points as prose or a tight list.
6. **Who it's for / use cases** — concrete reader scenarios.
7. **Under the hood** *(optional, light)* — one paragraph from containers/components.
8. **Conclusion + CTA**, then footer.

## Print / PDF

Whitepapers get saved as PDF — make "Print → Save as PDF" clean:

```css
@page { margin: 18mm; }
@media print {
  .cover { min-height: 92vh; }
  h2 { break-after: avoid; }
  .card, .callout, .strip { break-inside: avoid; }
  a[href^="http"]::after { content: ""; } /* don't print raw URLs */
}
.page-break { break-before: page; }
```

## Do / don't

- **Do** lead every block with the reader's outcome, not the mechanism. Generous
  whitespace. Strong, scannable headings. Real numbers from the register/scorecard.
- **Don't** paste register descriptions verbatim (they're dev-facing), invent
  capabilities the register doesn't contain, use hype words, or cram — a glossy doc
  that reads like a manual has failed. If it needs a scrollbar of dense text, cut.
