# Changelog

All notable changes to Plumbline are documented here. The format follows
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/); the project will follow
[Semantic Versioning](https://semver.org/spec/v2.0.0.html) once it reaches 1.0.

Releases ship the engine binaries + skills on the `dist` branch (see
[ADR 001](docs/adrs/001-plugin-packaging-and-distribution.md)); `main` stays
source-only.

Entries are authored under `[Unreleased]` as part of each PR (see
[raising-a-pr.md](docs/processes/raising-a-pr.md)) and promoted to a dated version at
release — never generated from commit messages.

## [Unreleased]

## [0.1.0] — 2026-08-27

### Added
- Zero-dependency Go engine: OFT-compatible anchor scanner, OFT-native Markdown register
  parser (single file or a tree), locked C4 artifact-type ladder with structural
  validation, and a bidirectional coverage report (coverage + completeness; uncovered,
  transitive, broken, orphan, and structural gaps).
- Strict and threshold gating; JSON and text output; optional `plumbline.json` config.
- Claude Code plugin with the `onboard`, `maintain`, `enforce`, `audit`, and `showcase`
  skills, shipped as cross-platform static binaries (linux amd64/arm64, macOS universal,
  windows amd64) via a POSIX launcher.
- Engineering process docs and ADRs. Plumbline traces itself (100%).

[Unreleased]: https://github.com/dr-zog/plumbline/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/dr-zog/plumbline/releases/tag/v0.1.0
