# boludo

## Installation

You can build it manually with commands: `make && make build`.

## Development

Use `make` (GNU or BSD):

- `make` - install dependencies
- `make test` - runs test
- `make check` - static code analysis
- `make build` - compile binaries from latest commit
- `make dist` - compile binaries from latest commit for supported OSes
- `make clean` - removes compilation artifacts
- `make cli-release` - tag latest commit as a new release of CLI
- `make info` - print system info (useful for debugging).

### Versioning

The repo contains command-line utility which versions are tagged as `cli/vYYYY.0M.MICRO` (_[calendar versioning](https://calver.org/)_).

## License

MIT ([in simple words](https://www.tldrlegal.com/license/mit-license))
