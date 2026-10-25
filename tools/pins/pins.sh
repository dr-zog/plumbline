#!/bin/sh
# Release pins in the human-facing docs (#53). The docs show copy-paste CI recipes
# pinned to a release tag; this keeps those pins on the latest release.
#
#   pins.sh write X.Y.Z   rewrite every pin to vX.Y.Z (semantic-release prepareCmd)
#   pins.sh check X.Y.Z   fail if any pin is not vX.Y.Z (version-sync guard)
#
# The shipped CI snippets carry no literal pin: they use __PLUMBLINE_VERSION__,
# which the `enforce` skill fills from the installed plugin's own plugin.json.
# [impl->component~build-tooling~1]
set -eu

FILES="README.md docs/ci.md"
# A pin: a release download URL, a `V=`/`V:` variable, or a `go install …@` tag.
PIN='(releases/download/|V=|V: |cmd/plumbline@)v[0-9]+\.[0-9]+\.[0-9]+'

mode=${1:-}
ver=${2:-}
case "$ver" in
[0-9]*.[0-9]*.[0-9]*) ;;
*) echo "usage: $0 write|check X.Y.Z" >&2; exit 2 ;;
esac

for f in $FILES; do
	[ -f "$f" ] || { echo "pins.sh: missing pin file $f — update FILES" >&2; exit 2; }
done

case "$mode" in
write)
	# Portable in-place edit (BSD sed reads `-i -E` as a backup suffix).
	for f in $FILES; do
		sed -E "s#$PIN#\\1v$ver#g" "$f" >"$f.tmp" && mv "$f.tmp" "$f"
	done
	;;
check)
	snippets=plugins/plumbline/skills/enforce/references/ci-snippets.md
	if grep -nE "$PIN" "$snippets"; then
		echo "::error::$snippets must pin __PLUMBLINE_VERSION__, not a literal tag — enforce fills it from plugin.json (#53)."
		exit 1
	fi
	# Search every tracked file, not just FILES, for anything shaped like a release
	# tag: a pin added elsewhere, or in another form, would never be rewritten at
	# release and would go stale unseen — the very failure #53 fixes. History
	# (the changelog, ADRs) legitimately names old versions.
	tags=$(git grep -nE 'v[0-9]+\.[0-9]+\.[0-9]+' -- . ':!CHANGELOG.md' ':!docs/adrs' || true)
	bad=""
	nl='
'
	IFS=$nl
	for hit in $tags; do
		file=${hit%%:*}
		case " $FILES " in
		*" $file "*) ;;
		*) bad="$bad$hit (not in a pin file: add it to FILES in $0)$nl"; continue ;;
		esac
		others=$(printf '%s\n' "$hit" | grep -oE 'v[0-9]+\.[0-9]+\.[0-9]+' | grep -vx "v$ver" || true)
		if [ -n "$others" ] || ! printf '%s\n' "$hit" | grep -qE "$PIN"; then
			bad="$bad$hit$nl"
		fi
	done
	unset IFS
	if [ -n "$bad" ]; then
		printf '%s' "$bad"
		echo "::error::Release pins must be v$ver, in a recognised pin form, in a file pins.sh rewrites. They are bot-maintained at release (#53) — don't hand-edit them."
		exit 1
	fi
	;;
*) echo "usage: $0 write|check X.Y.Z" >&2; exit 2 ;;
esac
