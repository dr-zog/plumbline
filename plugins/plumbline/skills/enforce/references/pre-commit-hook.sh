#!/bin/sh
# Plumbline pre-commit gate — blocks a commit when requirement↔code traceability
# doesn't hold.
#
# The `enforce` skill installs this hook, substituting __PLUMBLINE__ with the
# absolute path to the plugin's engine launcher (${CLAUDE_PLUGIN_ROOT}/bin/
# plumbline resolved to an absolute path — the git hook runs outside Claude Code,
# so the path is baked in here at install time). The launcher picks the right
# per-platform binary. Bypass a single commit with `git commit --no-verify`
# (which is why a CI gate is worth having as the real backstop).
set -eu

PLUMBLINE="__PLUMBLINE__"

if [ ! -x "$PLUMBLINE" ]; then
	# Don't block on missing tooling — CI is the hard backstop. Warn and pass.
	echo "plumbline: engine launcher not found at $PLUMBLINE; skipping local gate (CI still enforces)" >&2
	exit 0
fi

if ! "$PLUMBLINE"; then
	echo "" >&2
	echo "plumbline: traceability gate failed." >&2
	echo "  fix the anchors/register above, or bypass this one commit with --no-verify." >&2
	exit 1
fi
