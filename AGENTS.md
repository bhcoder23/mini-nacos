# AGENTS.md

This file applies to the entire repository.

## Project Identity

- This repository is a `Kratos`-layout `mini-nacos` learning project.
- The goal is to learn and reimplement core `Nacos` mechanisms, not to ship a full production clone.
- Preserve upstream `Nacos` semantics where possible, even when the implementation is simplified.

## Architecture Rules

- Follow the current `Kratos` layout:
  - `api`: proto contracts and generated transport code
  - `internal/service`: thin application / transport service layer
  - `internal/biz`: business rules, domain objects, repo interfaces
  - `internal/data`: repo, cache, in-memory or persistence implementations
  - `internal/server`: HTTP / gRPC server registration
  - `cmd/mini-nacos`: application bootstrap and `wire` injection
- Do not move business rules into `service`.
- Do not move transport concerns into `biz` or `data`.
- Keep `service` thin: request mapping, usecase call, response mapping, error return.
- Keep `biz` responsible for the real config semantics such as `save-first`, `md5` comparison, and lightweight change notification.

## Proto and Generated Files

- Do not hand-edit generated files:
  - `api/**/*.pb.go`
  - `api/**/*_grpc.pb.go`
  - `api/**/*_http.pb.go`
  - `internal/conf/conf.pb.go`
  - `openapi.yaml`
- If `api/**/*.proto` changes, regenerate with:
  - `make api`
- If `internal/conf/**/*.proto` changes, regenerate with:
  - `make config`
- If provider wiring changes, regenerate with:
  - `go generate ./...`
  - or `wire gen ./cmd/mini-nacos`

## Testing and Verification

- Preferred verification for normal changes:
  - `go test ./...`
  - `go build ./...`
- Prefer focused tests first, then broader verification.
- Keep tests close to the layer they verify:
  - `biz` behavior tests in `internal/biz`
  - `service` mapping / transport tests in `internal/service`

## Documentation

- Keep `docs/study`, `docs/design`, and `docs/homework` distinct:
  - `study`: source-chain notes, diagrams, conclusions, review questions
  - `design`: Kratos-layer mapping, models, boundaries, preserved semantics
  - `homework`: executable tasks and CR checkpoints
- Prefer diagrams embedded in Markdown (`mermaid` / ASCII).
- Do not migrate the old `mini-nacos-bak/docs/img` directory into this repo unless a later task explicitly requires binary assets.
- Preserve `Nacos Config` essentials in docs:
  - event-driven thinking
  - responsibility decoupling
  - lightweight change notification
  - separation of query and listen
- Avoid copying the same long explanation across all three doc types.

## Naming Guidance

- Proto request / response messages follow proto conventions:
  - `XxxRequest`
  - `XxxResponse`
- `biz` models use direct domain names:
  - `ConfigKey`
  - `ConfigItem`
  - `ConfigChange`
- Do not add Java-style suffixes like `Entity`, `DO`, or `DTO` to `biz` domain objects unless there is a real need.
