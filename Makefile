# Copyright (C) 2023 Maciej Żok
#
# SPDX-License-Identifier: GPL-3.0-or-later

.POSIX:
.SUFFIXES:

CLI_DIR = ./cmd/boludo

# MAIN TARGETS

all: install-dependencies

clean:
	@echo '# Delete binaries: rm -rf ./dist' >&2
	@rm -rf ./dist

info:
	@printf '# OS info: '
	@uname -rsv;
	@echo '# Development dependencies:'
	@go version || true
	@reuse --version || true
	@echo '# Go environment variables:'
	@go env || true

check:
	@echo '# Static analysis: go vet' >&2
	@go vet -C cmd/*
	@echo '# License check: reuse lint' >&2
	@reuse lint
	
test:
	@echo '# Unit tests: go test .' >&2
	@go test .

build:
	@echo '# Create release binary: ./dist/boludo' >&2
	@CURRENT_VER_TAG="$$(git tag --points-at HEAD | grep "^cli" | sed 's/^cli\/v//' | sort -t. -k 1,1n -k 2,2n -k 3,3n | tail -1)"; \
		PREV_VER_TAG="$$(git tag | grep "^cli" | sed 's/^cli\/v//' | sort -t. -k 1,1n -k 2,2n -k 3,3n | tail -1)"; \
		CURRENT_COMMIT_TAG="$$(TZ=UTC git --no-pager show --quiet --abbrev=12 --date='format-local:%Y%m%d%H%M%S' --format='%cd-%h')"; \
		PSEUDOVERSION="$${PREV_VER_TAG:-0001.01}-$$CURRENT_COMMIT_TAG"; \
		VERSION="$${CURRENT_VER_TAG:-$$PSEUDOVERSION}"; \
		go build -C $(CLI_DIR) -ldflags="-s -w -X main.AppVersion=$$VERSION" -o '../../dist/boludo'; \

dist:
	@echo '# Create release binaries in ./dist' >&2
	@CURRENT_VER_TAG="$$(git tag --points-at HEAD | grep "^cli" | sed 's/^cli\/v//' | sort -t. -k 1,1n -k 2,2n -k 3,3n | tail -1)"; \
		PREV_VER_TAG="$$(git tag | grep "^cli" | sed 's/^cli\/v//' | sort -t. -k 1,1n -k 2,2n -k 3,3n | tail -1)"; \
		CURRENT_COMMIT_TAG="$$(TZ=UTC git --no-pager show --quiet --abbrev=12 --date='format-local:%Y%m%d%H%M%S' --format='%cd-%h')"; \
		PSEUDOVERSION="$${PREV_VER_TAG:-0001.01}-$$CURRENT_COMMIT_TAG"; \
		VERSION="$${CURRENT_VER_TAG:-$$PSEUDOVERSION}"; \
		GOOS=openbsd GOARCH=amd64 go build -C $(CLI_DIR) -ldflags="-s -w -X main.AppVersion=$$VERSION" -o '../../dist/boludo-openbsd_amd64'; \
		GOOS=linux GOARCH=amd64 go build -C $(CLI_DIR) -ldflags="-s -w -X main.AppVersion=$$VERSION" -o '../../dist/boludo-linux_amd64'; \
		GOOS=windows GOARCH=amd64 go build -C $(CLI_DIR) -ldflags="-s -w -X main.AppVersion=$$VERSION" -o '../../dist/boludo-windows_amd64.exe'; \

	@echo '# Create binaries checksum' >&2
	@cd ./dist; sha256sum * >sha256sum.txt

install-dependencies:
	@echo '# Install CLI dependencies:' >&2
	# @go get -C cmd/ -v -x .
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
