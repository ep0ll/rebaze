# rebaze

**rebaze** is a production-oriented CLI for efficient OCI / Docker image mutation and rebasing.

It lets you replace base layers of an image (or apply structured patches) without a full rebuild — ideal for rolling out base-image security fixes, OS updates, or configuration changes across many application images.

Inspired by `crane rebase`, Cloud Native Buildpacks rebaser, and OCI mutation patterns, with an extensible patch/audit schema designed for enterprise governance.

## Status

| Command   | Status          | Notes                                      |
|-----------|-----------------|--------------------------------------------|
| `inspect` | ✅ Working      | Manifest / index inspection                |
| `rebase`  | ✅ Working      | Core rebase (old-base → new-base)          |
| `apply`   | 🚧 Planned      | Apply structured MutationBundle / patch    |
| `patch`   | 🚧 Planned      | Create patch artifacts                     |
| `preview` | 🚧 Planned      | Dry-run mutation plan                      |
| `history` | 🚧 Planned      | Audit / mutation history                   |
| `rollback`| 🚧 Planned      | Inverse apply / pointer revert             |
| `sign`    | 🚧 Planned      | Cosign / notation signing                  |
| `copy`    | 🚧 Planned      | Efficient registry copy                    |
| `export`  | 🚧 Planned      | Export mutation plan / SBOM                |

## Installation

```bash
go install github.com/ep0ll/rebaze/cli@latest
```

Or build from source:

```bash
git clone https://github.com/ep0ll/rebaze.git
cd rebaze
make build
```

## Quick Start

### Inspect an image

```bash
rebaze inspect ubuntu:latest
rebaze inspect ghcr.io/example/app@sha256:...
```

### Rebase an image onto a new base

```bash
rebaze rebase my-app:1.2.3 \
  --old-base ubuntu:22.04 \
  --new-base ubuntu:24.04 \
  --tag my-app:1.2.3-rebased
```

If the original image carries the standard OCI base annotations
(`org.opencontainers.image.base.digest` / `org.opencontainers.image.base.name`)
the `--old-base` flag can often be omitted.

## Architecture

```
cli/                  # Cobra CLI entrypoint
  internal/bazer/     # Command implementations
internal/             # Core domain packages (apply, rebase, history, …)
specs/v1/             # Typed Go representations of the mutation schema
schema/               # JSON Schema definitions (MutationBundle, patches, …)
```

The long-term design centres on a **MutationBundle** (see `schema.json`) that
supports:

- DAG-ordered operations (manifest / config / layer / blob / verify)
- Governance, approvals, and policy constraints
- Full audit ledger and recovery plans
- Observability hooks

## Development

```bash
make tidy          # go mod tidy
make build         # build binary to bin/rebaze
make test          # unit tests
make lint          # golangci-lint (if installed)
make ci            # local CI approximation
```

## License

Apache License 2.0 — see [LICENSE](LICENSE).
