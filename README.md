# rebaze

**rebaze** is a CLI for efficient OCI / Docker image mutation and rebasing.
Replace base layers or apply structured config/layer plans without a full rebuild.

## Commands

| Command    | Status |
|------------|--------|
| `inspect`  | Implemented |
| `rebase`   | Implemented |
| `apply`    | Implemented |
| `preview`  | Implemented |
| `patch`    | Implemented |
| `copy`     | Implemented |
| `export`   | Implemented |
| `history`  | Implemented (local ledger) |
| `rollback` | Implemented (from local ledger) |
| `sign`     | Implemented (ed25519 plan signatures) |

## Install

```bash
go install github.com/ep0ll/rebaze/cli@latest
# or
make build
```

## Usage

```bash
rebaze inspect ubuntu:latest

rebaze rebase my-app:1.2.3 \
  --old-base ubuntu:22.04 \
  --new-base ubuntu:24.04 \
  --tag my-app:1.2.3-rebased

rebaze patch --image my-app:1.2.3 --set-env FOO=bar --set-user 65532 --out plan.json
rebaze preview --plan plan.json
rebaze apply --plan plan.json --tag my-app:1.2.3-mutated

rebaze copy alpine:latest localhost:5000/alpine:latest
rebaze export alpine:latest

rebaze sign --keygen --private ed25519.key --public ed25519.pub
rebaze sign plan.json --private ed25519.key
rebaze sign plan.json --verify --public ed25519.pub

rebaze history
rebaze rollback my-app:1.2.3-mutated
```

### Mutation plan format

```json
{
  "schemaVersion": "1.0",
  "kind": "MutationPlan",
  "image": "my-app:1.0",
  "tag": "my-app:1.0-mutated",
  "config": {
    "setEnv": { "FOO": "bar" },
    "unsetEnv": ["LEGACY"],
    "setLabel": { "app": "my-app" },
    "setEntrypoint": ["/usr/bin/app"],
    "setUser": "65532",
    "setWorkdir": "/home/nonroot"
  },
  "layers": {
    "delete": [0],
    "insert": [{ "index": 1, "from": "busybox:latest" }],
    "append": ["ghcr.io/example/sidecar-layer:1"]
  },
  "annotations": {
    "org.opencontainers.image.description": "mutated by rebaze"
  }
}
```

## Development

```bash
make tidy
make build
make test
```

## License

Apache License 2.0 — see [LICENSE](LICENSE).
