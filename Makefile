DOCKER_CMD=docker --config ~/.docker/

.PHONY: docker
docker: docker-image docker-push

.PHONY: docker-push
docker-push:
	$(DOCKER_CMD) tag code.ndumas.com/ndumas/obsidian-pipeline:$(VERSION) code.ndumas.com/ndumas/obsidian-pipeline:latest
	$(DOCKER_CMD) push code.ndumas.com/ndumas/obsidian-pipeline:latest
	$(DOCKER_CMD) push code.ndumas.com/ndumas/obsidian-pipeline:$(VERSION)

.PHONY: docker-image
docker-image:
	$(DOCKER_CMD) build --build-arg VERSION=$(VERSION) -t code.ndumas.com/ndumas/obsidian-pipeline:$(VERSION) .

.PHONY: build-alpine
build-alpine:
	CGO_ENABLED=0 GOOS=linux go build -ldflags="buildmode=exe $(LDFLAGS) -linkmode external -w -extldflags '-static'"  -o $(DISTDIR)/$(NAME)-$(VERSION)-alpine/obp cmd/obp/*.go

# This file is intended as a starting point for a customized makefile for a Go project.
#
# Targets:
# 	all: Format, check, build, and test the code
#   setup: Install build/test toolchain dependencies (e.g. gox)
#   lint: Run linters against source code
#   bump-{major,minor,patch}: create a new semver git tag
#   release-{major,minor,patch}: push a tagged release
# 	format: Format the source files
# 	build: Build the command(s) for target OS/arch combinations
# 	install: Install the command(s)
# 	clean: Clean the build/test artifacts
#   report: Generate build/test reports
#	check: Run tests
#   bench: Run benchmarks
#   dist: zip/tar binaries & documentation
#	debug: print parameters
#
# Parameters:
# 	VERSION: release version in semver format
#	BUILD_TAGS: additional build tags to pass to go build
#	DISTDIR: path to save distribution files
#	RPTDIR: path to save build/test reports
#
# Assumptions:
#   - Your package contains a cmd/ package, containing a directory for each produced binary.
#   - You have cloc installed and accessible in the PATH.
#   - Your GOPATH and GOROOT are set correctly.
#   - Your makefile is in the root of your package and does not have a space in its file name.
#   - Your root package contains global string variables Version and Build, to receive the bild version number and commit ID, respectively.
#
# Features:
#   - report generates files that can be consumed by Jenkins, as well as a list of external dependencies.
#   - setup installs all the tools aside from cloc.
#   - Works on Windows and with paths containing spaces.
#   - Works when executing from outside the makefile directory using -f.
#   - Targets are useful both in CI and developer workstations.
#   - Handles cross-compiation for multiple OSes and architectures.
#   - Bundles binaries and documentation into compressed archives, using tar/gz for Linux and Darwin, and zip for Windows.


# Parameters
PKG = code.ndumas.com/ndumas/obsidian-pipeline
NAME = obp
DOC = README.md LICENSE


# Replace backslashes with forward slashes for use on Windows.
# Make is !@#$ing weird.
E :=
BSLASH := \$E
FSLASH := /

# Directories
WD := $(subst $(BSLASH),$(FSLASH),$(shell pwd))
MD := $(subst $(BSLASH),$(FSLASH),$(shell dirname "$(realpath $(lastword $(MAKEFILE_LIST)))"))
PKGDIR = $(MD)
CMDDIR = $(PKGDIR)/cmd
DISTDIR ?= $(WD)/dist
RPTDIR ?= $(WD)/reports
GP = $(subst $(BSLASH),$(FSLASH),$(GOPATH))

# Parameters
VERSION ?= $(shell git -C "$(MD)" describe --tags --dirty=-dev)
COMMIT_ID := $(shell git -C "$(MD)" rev-parse HEAD | head -c8)
BUILD_TAGS ?= release
CMDPKG = $(PKG)/cmd
CMDS := $(shell find "$(CMDDIR)/" -mindepth 1 -maxdepth 1 -type d | sed 's/ /\\ /g' | xargs -n1 basename)
BENCHCPUS ?= 1,2,4

# Commands
GOCMD = go
ARCHES ?= amd64 386
OSES ?= windows linux darwin
OUTTPL = $(DISTDIR)/$(NAME)-$(VERSION)-{{.OS}}_{{.Arch}}/{{.Dir}}
LDFLAGS = -X '$(PKG).Version=$(VERSION)' -X '$(PKG).Build=$(COMMIT_ID)'
GOBUILD = gox -osarch="!darwin/386" -rebuild -gocmd="$(GOCMD)" -arch="$(ARCHES)" -os="$(OSES)" -output="$(OUTTPL)" -tags "$(BUILD_TAGS)" -ldflags "$(LDFLAGS)"
GOCLEAN = $(GOCMD) clean
GOINSTALL = $(GOCMD) install -a -tags "$(BUILD_TAGS)" -ldflags "$(LDFLAGS)"
GOTEST = $(GOCMD) test -v -tags "$(BUILD_TAGS)"
DISABLED_LINTERS = varnamelen,interfacer,ifshort,exhaustivestruct,maligned,varcheck,scopelint,structcheck,deadcode,nosnakecase,golint,depguard
GOLINT = golangci-lint run --enable-all --disable "$(DISABLED_LINTERS)" --timeout=30s --tests
GODEP = $(GOCMD) get -d -t
GOFMT = goreturns -w
GOBENCH = $(GOCMD) test -v -tags "$(BUILD_TAGS)" -cpu=$(BENCHCPUS) -run=NOTHING -bench=. -benchmem -outputdir "$(RPTDIR)"
GZCMD = tar -czf
ZIPCMD = zip
SHACMD = sha256sum
SLOCCMD = cloc --by-file --xml --exclude-dir="vendor" --include-lang="Go"
XUCMD = go2xunit

# Dynamic Targets
INSTALL_TARGETS := $(addprefix install-,$(CMDS))

.PHONY: all

all: debug setup dep format lint test bench build dist

release-major: bump-major
	git push origin main --tags
	git push github main --tags

release-minor: bump-minor
	git push origin main --tags
	git push github main --tags

release-patch: bump-patch
	git push origin main --tags
	git push github main --tags


setup: setup-dirs setup-build setup-format setup-lint setup-reports setup-bump

setup-bump:
	go install github.com/guilhem/bump@latest

bump-major: setup-bump
	bump major

bump-minor: setup-bump
	bump minor

bump-patch: setup-bump
	bump patch

setup-reports: setup-dirs
	go install github.com/tebeka/go2xunit@latest

setup-build: setup-dirs
	go install github.com/mitchellh/gox@latest

setup-format: setup-dirs
	go install github.com/sqs/goreturns@latest

setup-lint: setup-dirs
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.53.1

setup-dirs:
	mkdir -p "$(RPTDIR)"
	mkdir -p "$(DISTDIR)"

clean:
	$(GOCLEAN) $(PKG)
	rm -vrf "$(DISTDIR)"/*
	rm -vf "$(RPTDIR)"/*

format:
	$(GOFMT) "$(PKGDIR)"

dep:
	$(GODEP) $(PKG)/...

lint: setup-dirs dep
	$(GOLINT) "$(PKGDIR)" | tee "$(RPTDIR)/lint.out"

check: setup-dirs clean dep
	$(GOTEST) $$(go list "$(PKG)/..." | grep -v /vendor/) | tee "$(RPTDIR)/test.out"

bench: setup-dirs clean dep
	$(GOBENCH) $$(go list "$(PKG)/..." | grep -v /vendor/) | tee "$(RPTDIR)/bench.out"

report: check
	cd "$(PKGDIR)";$(SLOCCMD) --out="$(RPTDIR)/cloc.xml" . | tee "$(RPTDIR)/cloc.out"
	cat "$(RPTDIR)/test.out" | $(XUCMD) -output "$(RPTDIR)/tests.xml"
	go list -f '{{join .Deps "\n"}}' "$(CMDPKG)/..." | sort | uniq | xargs -I {} sh -c "go list -f '{{if not .Standard}}{{.ImportPath}}{{end}}' {} | tee -a '$(RPTDIR)/deps.out'"

build: $(CMDS)
$(CMDS): setup-dirs dep
	$(GOBUILD) "$(CMDPKG)/$@" | tee "$(RPTDIR)/build-$@.out"
install: $(INSTALL_TARGETS)
$(INSTALL_TARGETS):
	$(GOINSTALL) "$(CMDPKG)/$(subst install-,,$@)"

dist: clean build
	for docfile in $(DOC); do \
		for dir in "$(DISTDIR)"/*; do \
			cp "$(PKGDIR)/$$docfile" "$$dir/"; \
		done; \
	done
	cd "$(DISTDIR)"; for dir in ./*linux*; do $(GZCMD) "$(basename "$$dir").tar.gz" "$$dir"; done
	cd "$(DISTDIR)"; for dir in ./*windows*; do $(ZIPCMD) "$(basename "$$dir").zip" "$$dir"; done
	cd "$(DISTDIR)"; for dir in ./*darwin*; do $(GZCMD) "$(basename "$$dir").tar.gz" "$$dir"; done
	cd "$(DISTDIR)"; find . -maxdepth 1 -type f -printf "$(SHACMD) %P | tee \"./%P.sha\"\n" | sh
	$(info "Built v$(VERSION), build $(COMMIT_ID)")

debug:
	$(info MD=$(MD))
	$(info WD=$(WD))
	$(info PKG=$(PKG))
	$(info PKGDIR=$(PKGDIR))
	$(info DISTDIR=$(DISTDIR))
	$(info VERSION=$(VERSION))
	$(info COMMIT_ID=$(COMMIT_ID))
	$(info BUILD_TAGS=$(BUILD_TAGS))
	$(info CMDS=$(CMDS))
	$(info BUILD_TARGETS=$(BUILD_TARGETS))
	$(info INSTALL_TARGETS=$(INSTALL_TARGETS))
