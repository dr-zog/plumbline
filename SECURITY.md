# Security Policy

## The threat model

Plumbline is built to run unattended — an agent or a CI pipeline pointing it at a repository and trusting the exit code. So the question that matters is: **what can a malicious or malformed repository do to the engine as it scans?** The engine only ever *reads* your tree, makes no network calls, and has zero runtime dependencies, so the attack surface is deliberately small — but "small" is not "none". The surfaces worth probing:

- **Parsing** — malformed registers or anchors driving the scanner into bad behaviour: pathological input, resource exhaustion, or escaping the intended scan scope.
- **The distribution path** — the `dist` branch and the prebuilt binaries it carries.

If you find a way to make the engine execute, leak, or corrupt something it shouldn't while scanning, that's a vulnerability and we want to know.

## Reporting a vulnerability

**Please report privately, not in a public issue.**

- Preferred: open a [private security advisory](https://github.com/dr-zog/plumbline/security/advisories/new).
- Or email **me@dr-zog.com** with the details and, ideally, a minimal reproduction — a small register plus the source that triggers it.

We'll acknowledge within a few days, keep you informed as we investigate, and credit you in the release notes unless you'd rather stay anonymous.

## Out of scope

A vulnerability in *your own* code that Plumbline merely reports on is not a Plumbline vulnerability — the engine traces the thread, it does not analyse or execute your program. Problems in the requirement→code chain you're tracking belong in your tracker, not here.

## Supported versions

Plumbline is pre-1.0; fixes land on `main` and ship in the next tagged release. There is no back-port branch yet — run the latest release.
