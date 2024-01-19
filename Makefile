.POSIX:
.SUFFIXES:


#
# PUBLIC MACROS
#

CLI     = boludo
DESTDIR = ./dist
GO      = go
GOFLAGS = 
LDFLAGS = -ldflags "-s -w -X main.AppVersion=$$VERSION"


#
# INTERNAL MACROS
#

CLI_DIR  = ./cmd/boludo


#
# DEVELOPMENT TASKS
#

all: install-dependencies

clean:
	@echo '# Delete bulid directory' >&2
	rm -rf $(DESTDIR)

info:
	@printf '# OS info: '
	@uname -rsv;
	@echo '# Development dependencies:'
	@$(GO) version || true
	@echo '# Go environment variables:'
	@$(GO) env || true

check:
	@echo '# Static analysis' >&2
	$(GO) vet -C $(CLI_DIR)
	
test:
	@echo '# Unit tests' >&2
	$(GO) test .

build:
	@echo '# Build CLI executable: $(DESTDIR)/$(CLI)' >&2
	@CURRENT_VER_TAG="$$(git tag --points-at HEAD | grep "^cli" | sed 's/^cli\/v//' | sort -t. -k 1,1n -k 2,2n -k 3,3n | tail -1)"; \
		PREV_VER_TAG="$$(git tag | grep "^cli" | sed 's/^cli\/v//' | sort -t. -k 1,1n -k 2,2n -k 3,3n | tail -1)"; \
		CURRENT_COMMIT_TAG="$$(TZ=UTC git --no-pager show --quiet --abbrev=12 --date='format-local:%Y%m%d%H%M%S' --format='%cd-%h')"; \
		PSEUDOVERSION="$${PREV_VER_TAG:-0001.01}-$$CURRENT_COMMIT_TAG"; \
		VERSION="$${CURRENT_VER_TAG:-$$PSEUDOVERSION}"; \
		$(GO) build -C $(CLI_DIR) $(GOFLAGS) $(LDFLAGS) -o '../../$(DESTDIR)/$(CLI)'

dist:
	@echo '# Create CLI executables in $(DESTDIR)' >&2
	@CURRENT_VER_TAG="$$(git tag --points-at HEAD | grep "^cli" | sed 's/^cli\/v//' | sort -t. -k 1,1n -k 2,2n -k 3,3n | tail -1)"; \
		PREV_VER_TAG="$$(git tag | grep "^cli" | sed 's/^cli\/v//' | sort -t. -k 1,1n -k 2,2n -k 3,3n | tail -1)"; \
		CURRENT_COMMIT_TAG="$$(TZ=UTC git --no-pager show --quiet --abbrev=12 --date='format-local:%Y%m%d%H%M%S' --format='%cd-%h')"; \
		PSEUDOVERSION="$${PREV_VER_TAG:-0001.01}-$$CURRENT_COMMIT_TAG"; \
		VERSION="$${CURRENT_VER_TAG:-$$PSEUDOVERSION}"; \
		GOOS=openbsd GOARCH=amd64 $(GO) build -C $(CLI_DIR) $(GOFLAGS) $(LDFLAGS) -o "../../$(DESTDIR)/$(CLI)-openbsd_amd64"; \
		GOOS=linux GOARCH=amd64 $(GO) build -C $(CLI_DIR) $(GOFLAGS) $(LDFLAGS) -o "../../$(DESTDIR)/$(CLI)-linux_amd64"; \
		GOOS=windows GOARCH=amd64 $(GO) build -C $(CLI_DIR) $(GOFLAGS) $(LDFLAGS) -o "../../$(DESTDIR)/$(CLI)-windows_amd64.exe"; \

	@echo '# Create checksums' >&2
	@cd $(DESTDIR); sha256sum * >sha256sum.txt

install-dependencies:
	@echo '# Install CLI dependencies' >&2
	@GOFLAGS='-v -x' $(GO) get -C $(CLI_DIR) $(GOFLAGS) .
	@echo '# Build libllama.so'
	cd ./external/llama.cpp; make server && cp server ../../llm-server

cli-release: check test
	@echo '# Update local branch' >&2
	@git pull --rebase
	@echo '# Create new CLI release tag' >&2
	@PREV_VER_TAG=$$(git tag | grep "^cli" | sed 's/^cli\/v//' | sort -t. -k 1,1n -k 2,2n -k 3,3n | tail -1); \
		printf 'Choose new version number for CLI (calver; >%s): ' "$${PREV_VER_TAG:-2023.11}"
	@read -r VERSION; \
		git tag "cli/v$$VERSION"; \
		git push --tags
