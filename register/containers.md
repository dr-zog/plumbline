# Containers

The C4 containers — the two layers Plumbline ships as. Each names its chosen
technology.

### Engine
`container~engine~1`

The zero-dependency Go static binary — the law: it scans for anchors, parses the
register, resolves coverage and gates the build.

Covers: context~plumbline~1
Needs: component

### Plugin
`container~plugin~1`

The Claude Code plugin of Markdown skills — the taste: it authors the anchors and
keeps the register current, so a developer rarely hand-writes one.

Covers: context~plumbline~1
Needs: component
