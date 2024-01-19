# This Makefile intended to be POSIX-compliant (2018 edition with .PHONY target).
#
# .PHONY targets are used by:
#  - task definintions
#  - compilation of Go code (force usage of `go build` to changes detection).
#
# More info:
#  - docs: <https://pubs.opengroup.org/onlinepubs/9699919799/utilities/make.html>
#  - .PHONY: <https://www.austingroupbugs.net/view.php?id=523>
#
.POSIX:
.SUFFIXES:


#
# PUBLIC MACROS
#

CLI     = boludo
DESTDIR = ./dist
GO      = go
GOFLAGS = 
LDFLAGS = -ldflags "-s -w -X main.AppVersion=$(CLI_VERSION)"


#
# INTERNAL MACROS
#

CLI_DIR             = ./cmd/boludo
CLI_CURRENT_VER_TAG = $$(git tag --points-at HEAD | sed 's/^v//' | sort -t. -k 1,1n -k 2,2n -k 3,3n | tail -1)
CLI_LATEST_VERSION  = $$(git tag | sed 's/^v//' | sort -t. -k 1,1n -k 2,2n -k 3,3n | tail -1)
CLI_PSEUDOVERSION   = $$(VER="$(CLI_LATEST_VERSION)"; echo "$${VER:-0.0.0}")-$$(TZ=UTC git --no-pager show --quiet --abbrev=12 --date='format-local:%Y%m%d%H%M%S' --format='%cd-%h')
CLI_VERSION         = $$(VER="$(CLI_CURRENT_VER_TAG)"; echo "$${VER:-$(CLI_PSEUDOVERSION)}")


#
# DEVELOPMENT TASKS
#

.PHONY: all
all: install-dependencies

.PHONY: clean
clean:
	@echo '# Delete bulid directory' >&2
	rm -rf $(DESTDIR)

.PHONY: info
info:
	@printf '# OS info: '
	@uname -rsv;
	@echo '# Development dependencies:'
	@$(GO) version || true
	@echo '# Go environment variables:'
	@$(GO) env || true

.PHONY: check
check:
	@echo '# Static analysis' >&2
	$(GO) vet -C $(CLI_DIR)
	
.PHONY: test
test:
	@echo '# Unit tests' >&2
	$(GO) test .

.PHONY: build
build:
	@echo '# Build CLI executable: $(DESTDIR)/$(CLI)' >&2
	$(GO) build -C $(CLI_DIR) $(GOFLAGS) $(LDFLAGS) -o '../../$(DESTDIR)/$(CLI)'

.PHONY: dist
dist:
	@echo '# Create CLI executables in $(DESTDIR)' >&2
	GOOS=openbsd GOARCH=amd64 $(GO) build -C $(CLI_DIR) $(GOFLAGS) $(LDFLAGS) -o "../../$(DESTDIR)/$(CLI)-openbsd_amd64"
	GOOS=linux GOARCH=amd64 $(GO) build -C $(CLI_DIR) $(GOFLAGS) $(LDFLAGS) -o "../../$(DESTDIR)/$(CLI)-linux_amd64"
	GOOS=windows GOARCH=amd64 $(GO) build -C $(CLI_DIR) $(GOFLAGS) $(LDFLAGS) -o "../../$(DESTDIR)/$(CLI)-windows_amd64.exe"
	@echo '# Create checksums' >&2
	@cd $(DESTDIR); sha256sum * >sha256sum.txt

.PHONY: install-dependencies
install-dependencies:
	@echo '# Install CLI dependencies' >&2
	@GOFLAGS='-v -x' $(GO) get -C $(CLI_DIR) $(GOFLAGS) .
	@echo '# Build libllama.so'
	cd ./external/llama.cpp; make server && cp server ../../llm-server

.PHONY: cli-release
cli-release: check test
	@echo '# Update local branch' >&2
	@git pull --rebase
	@echo '# Create new CLI release tag' >&2
	@VER="$(CLI_LATEST_VERSION)"; printf 'Choose new version number (calver; >%s): ' "$${VER:-2023.11}"
	@read -r NEW_VERSION; \
		git tag "cli/v$$NEW_VERSION"; \
		git push --tags
