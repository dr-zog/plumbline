# Plumbline — build tooling.
#
# `make binaries` cross-compiles the zero-dependency engine into the plugin's
# bin/ directory, one static binary per target platform, with the version
# stamped in. macOS ships as a single universal (fat) binary, fused from the
# amd64 + arm64 builds by tools/makefat — pure Go, so it runs on the Linux CI
# runner with no macOS host and no lipo. The binaries are not committed; they
# ship via the release flow (see docs/adrs/001).

BINDIR  := plugins/plumbline/bin
PKG     := ./cmd/plumbline
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
GOFLAGS := -trimpath
LDFLAGS := -s -w -X main.version=$(VERSION)
GOBUILD := CGO_ENABLED=0 go build $(GOFLAGS) -ldflags '$(LDFLAGS)'

.PHONY: all vet test build binaries dist clean

all: vet test build

vet:
	go vet ./...

test:
	go test ./...

# Local build for development (root binary is gitignored).
build:
	go build $(GOFLAGS) -ldflags '$(LDFLAGS)' -o plumbline $(PKG)

# Per-platform release binaries — the zero-dependency promise.
binaries:
	@mkdir -p $(BINDIR)
	GOOS=linux   GOARCH=amd64 $(GOBUILD) -o $(BINDIR)/plumbline-linux-amd64   $(PKG)
	GOOS=linux   GOARCH=arm64 $(GOBUILD) -o $(BINDIR)/plumbline-linux-arm64   $(PKG)
	GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BINDIR)/plumbline-windows-amd64.exe $(PKG)
	@echo "fusing macOS universal binary"
	GOOS=darwin  GOARCH=amd64 $(GOBUILD) -o $(BINDIR)/.darwin-amd64 $(PKG)
	GOOS=darwin  GOARCH=arm64 $(GOBUILD) -o $(BINDIR)/.darwin-arm64 $(PKG)
	go run ./tools/makefat $(BINDIR)/plumbline-darwin-universal $(BINDIR)/.darwin-amd64 $(BINDIR)/.darwin-arm64
	@rm -f $(BINDIR)/.darwin-amd64 $(BINDIR)/.darwin-arm64
	@ls -l $(BINDIR)/plumbline-*

# Publish the plugin + binaries to the force-pushed orphan `dist` branch (ADR
# 001): `main` stays source-only; `dist` is always a single flat commit, built
# from a throwaway index with `git commit-tree`, so history never accumulates
# binaries and the working tree is never touched.
DIST_BRANCH := dist
dist: binaries
	@idx=$$(mktemp); \
	GIT_INDEX_FILE=$$idx git read-tree --empty; \
	GIT_INDEX_FILE=$$idx git add -f .claude-plugin/marketplace.json plugins/plumbline; \
	tree=$$(GIT_INDEX_FILE=$$idx git write-tree); \
	rm -f $$idx; \
	commit=$$(printf 'dist: plumbline %s' '$(VERSION)' | git commit-tree $$tree); \
	git push -f origin $$commit:refs/heads/$(DIST_BRANCH); \
	echo "published $(DIST_BRANCH) @ $(VERSION) ($$commit)"

clean:
	rm -f plumbline $(BINDIR)/plumbline-*
