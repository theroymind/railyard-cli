# railyard-cli

Process plumbing shared by the Railyard apps and `tix`: argument parsing, subcommands, usage errors, exit codes, JSON in and out, running a child binary, and the shared state directory.

```bash
go get github.com/theroymind/railyard-cli
```

| Package | Import path | Provides |
| --- | --- | --- |
| `cli` | `github.com/theroymind/railyard-cli` | Arguments, subcommands, exit codes, JSON output, child processes, state directory |
| `strictjson` | `github.com/theroymind/railyard-cli/strictjson` | Strict JSON decoding: no unknown or duplicate fields, exact-case names |

## Development

```bash
make hooks   # point git at .githooks
make check   # format, tidy, vet, test, lint
```
