# Build all by default, even if it's not first
.DEFAULT_GOAL := all

.PHONY: all
#all: tidy gen format lint build
all: tidy gen lint build

# ==============================================================================
# Build options

ROOT_PACKAGE=github.com/ClessLi/bifrost
VERSION_PACKAGE=github.com/marmotedu/component-base/pkg/version

# ==============================================================================
# Includes

include scripts/make-rules/common.mk # make sure include common.mk at the first include line
include scripts/make-rules/golang.mk
include scripts/make-rules/gen.mk
#include scripts/make-rules/release.mk
#include scripts/make-rules/dependencies.mk
include scripts/make-rules/tools.mk
# ==============================================================================
# Usage

define USAGE_OPTIONS

Options:
  DEBUG        Whether to generate debug symbols. Default is 0.
  BINS         The binaries to build. Default is all of cmd.
               This option is available when using: make build/build.multiarch
               Example: make build BINS="bifrost ng_conf_format"
  PLATFORMS    The multiple platforms to build. Default is linux_amd64 and windows_amd64.
               This option is available when using: make build.multiarch
               Example: make build.multiarch BINS="bifrost ng_conf_format" PLATFORMS="linux_amd64 windows_amd64"
  VERSION      The version information compiled into binaries.
               The default is obtained from gsemver or git.
  V            Set to 1 enable verbose build. Default is 0.
endef
export USAGE_OPTIONS

# ==============================================================================
# Target

## build: Build source code for host platform.
.PHONY: build
build:
	@$(MAKE) go.build

## build.multiarch: Build source code for multiple platforms. See option PLATFORMS.
.PHONY: build.multiarch
build.multiarch:
	@$(MAKE) go.build.multiarch

## clean: Remove all files that are created by building.
.PHONY: clean
clean:
	@echo "===========> Cleaning all build output"
	@-rm -vrf $(OUTPUT_DIR)

## lint: Check syntax and styling of go sources.
.PHONY: lint
lint:
	@$(MAKE) go.lint

## release: Release bifrost
.PHONY: release
release:
	@$(MAKE) release.run

## format: Gofmt (reformat) package sources (exclude vendor dir if existed).
.PHONY: format
format: tools.verify.golines tools.verify.goimports
	@echo "===========> Formating codes"
	@$(FIND) -type f -name '*.go' | $(XARGS) gofmt -s -w
	@$(FIND) -type f -name '*.go' | $(XARGS) goimports -w -local $(ROOT_PACKAGE)
	@$(FIND) -type f -name '*.go' | $(XARGS) golines -w --max-len=120 --reformat-tags --shorten-comments --ignore-generated .
	@$(GO) mod edit -fmt

### verify-copyright: Verify the boilerplate headers for all files.
#.PHONY: verify-copyright
#verify-copyright:
#	@$(MAKE) copyright.verify
#
### add-copyright: Ensures source code files have copyright license headers.
#.PHONY: add-copyright
#add-copyright:
#	@$(MAKE) copyright.add

## gen: Generate all necessary files, such as error code files.
.PHONY: gen
gen:
	@$(MAKE) gen.run

## install: Install bifrost system with all its components.
.PHONY: install
install:
	@$(MAKE) install.install

## dependencies: Install necessary dependencies.
.PHONY: dependencies
dependencies:
	@$(MAKE) dependencies.run

## tools: install dependent tools.
.PHONY: tools
tools:
	@$(MAKE) tools.install

## check-updates: Check outdated dependencies of the go projects.
.PHONY: check-updates
check-updates:
	@$(MAKE) go.updates

.PHONY: tidy
tidy:
	@$(GO) mod tidy

## help: Show this help info.
.PHONY: help
help: Makefile
	@echo -e "\nUsage: make <TARGETS> <OPTIONS> ...\n\nTargets:"
	@sed -n 's/^##//p' $< | column -t -s ':' | sed -e 's/^/ /'
	@echo "$$USAGE_OPTIONS"
