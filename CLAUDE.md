# railyard-cli

Process plumbing shared by the Railyard apps (`github.com/theroymind/railyard`) and `tix` (`github.com/theroymind/go-tix`). One module, `github.com/theroymind/railyard-cli`.

## Development

```text
make test     # go test ./...
make vet      # go vet ./...
make lint     # golangci-lint run
make check    # everything CI runs: format, tidy, vet, test, lint
make hooks    # point git at .githooks (pre-commit runs format, tidy, vet, lint)
```

## Layout

- The module root is package `cli`: args, subcommands, exit codes, JSON in and out, running a child binary, the state directory.
- `strictjson/` decodes JSON at a boundary (json/v2: no unknown or duplicate fields, exact-case names).
- No domain logic. Nothing here knows about tickets, events or any app.

## Code standards

- Every package has a package comment, in exactly one file, directly above the `package` clause.
- Errors are sentinel `var Err... = errors.New(...)` values, wrapped with `%w` and checked with `errors.Is`.
- Tests are table-driven.
- `gofmt -l .` prints nothing, and `go vet ./...` and `golangci-lint run` stay clean.

## Compatibility

Railyard and go-tix both import this module. A change to an exported name or to the state directory layout is a breaking change: tag a new major version, or update both consumers in step.

## Review checklist

- [ ] `make check` passes.
- [ ] No domain logic added.
- [ ] Ask before adding a dependency.
