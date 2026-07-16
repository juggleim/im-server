# Contributing to JuggleIM

First off — thank you for taking the time to contribute! 🎉 This document explains how to get a
development environment running and how to submit changes to `im-server`.

- 💬 Questions & discussion: [Telegram group](https://t.me/juggleim_zh) or [open an issue](https://github.com/juggleim/im-server/issues)
- 📖 Docs: https://www.juggle.im/docs/guide/intro/

## Ways to contribute

- **Report bugs** using the [bug report template](.github/ISSUE_TEMPLATE/bug_report.yml).
- **Request features** using the [feature request template](.github/ISSUE_TEMPLATE/feature_request.yml).
- **Improve docs** — even fixing a typo helps.
- **Submit code** via pull requests (see below).

Looking for a place to start? Check issues labeled
[`good first issue`](https://github.com/juggleim/im-server/labels/good%20first%20issue).

## Development setup

### Option A — Docker (fastest)

```bash
docker compose up -d
```

Brings up MySQL + im-server. WebSocket on `:9003`, server API on `:9001`, admin console on
`:8090` (`admin` / `123456`). See [`docker-compose.yml`](./docker-compose.yml) for details.

### Option B — From source

Prerequisites: **Go 1.25+** and a **MySQL 8** instance (MongoDB optional, only for
`msgStoreEngine: mongo`).

```bash
# 1. Create the database and import the schema
mysql -uroot -p -e "CREATE SCHEMA IF NOT EXISTS jim_db;"
mysql -uroot -p jim_db < sql/imserver.sql

# 2. Configure — copy the template and edit MySQL credentials
#    (config.yml is gitignored)
#    see launcher/scripts/config_template.yaml for the full option list

# 3. Run (all commands run from the launcher/ directory)
cd launcher
go run main.go
```

## Building & testing

```bash
# Build everything (target linux for the full build; some metrics code is OS-specific)
CGO_ENABLED=0 GOOS=linux go build ./...

# Cross-compile a release binary
cd launcher && sh build.sh          # -> launcher/build/imserver (linux/amd64)

# Run the self-contained unit tests
go test ./services/commonservices/...
go test ./commons/metrics/...
go test ./services/pushmanager/services/httputil

# Run a single test
go test -run TestName ./path/to/pkg/
```

> Note: tests under `simulator/tests/` are **integration** tests that connect to a running
> im-server over WebSocket. They require a live server and a provisioned app key — they are not
> hermetic, so don't expect them to pass in a plain `go test ./...`.

## Project layout (orientation)

- `launcher/` — application entry point (`main.go`) and config; all run/build commands live here.
- `services/` — self-contained services (message, group, conversation, gateways, ...). Each has
  `actors/`, `services/`, `storages/`, and a `starter.go`.
- `commons/gmicro/` — the custom actor + cluster runtime the whole system is built on.
- `commons/bases/` — cross-service RPC helpers (`SyncRpcCall`, `AsyncRpcCall`, `Broadcast`, ...).
- `commons/pbdefines/` — Protobuf wire definitions; generated Go lives in `pbobjs/`.

Services never import each other's internals — they communicate through the actor mesh. When adding
an actor, register it in the service's `starter.go` and reference it by its string method name.

## Coding guidelines

- Run `gofmt` (or `go fmt ./...`) before committing.
- Match the style and idioms of the surrounding code.
- **Regenerate** protobuf Go from the `.proto` files rather than hand-editing files in
  `commons/pbdefines/pbobjs`.
- Structural DB changes go through the auto-migration path (`dbcommons.Upgrade()`), not manual DDL.
- Don't commit secrets or local config (`launcher/conf/config.yml` is gitignored for a reason).

## Commit messages

We follow [Conventional Commits](https://www.conventionalcommits.org/) for new work:

```
feat: add conversation tag API
fix: prevent panic when sender is offline
docs: clarify docker deployment steps
```

Common types: `feat`, `fix`, `docs`, `refactor`, `perf`, `test`, `chore`.

## Pull request process

1. Fork the repo and create a branch from `master` (e.g. `feat/my-feature`).
2. Make your change; keep it focused — one logical change per PR.
3. Ensure it builds and relevant tests pass.
4. Fill out the PR template, linking any related issue (`Closes #123`).
5. A maintainer will review; please be responsive to feedback.

By contributing, you agree that your contributions are licensed under the same
[license](./LICENSE) as this project.
