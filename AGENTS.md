# Repository Guidelines

This repository implements a small BASIC interpreter in Go. Use this guide to navigate the codebase, run the project, and contribute changes consistently.

## Project Structure & Module Organization
- `cmd/basic/`: CLI entrypoint (`main.go`) and related tests.
- `lexer/`, `parser/`, `interpreter/`, `runtime/`, `types/`: core packages.
- `acceptance/`: Go harness and YAML acceptance tests (`acceptance/testdata/*.yaml`).
- `scripts/`: helper scripts (coverage/LOC history, runner).
- `testdata/`: sample BASIC programs (used by scripts).
- `spec.md`, `README.md`: specification and usage notes.

## Build, Test, and Development Commands
- `make help`: list available tasks.
- `make test` or `go test ./...`: run unit + acceptance tests.
- `scripts/run.sh testdata/guess_number.bas`: run the interpreter on a file.
- `go run ./cmd/basic -i "42, John" testdata/guess_number.bas`: provide inputs for `INPUT` statements.
- `go run ./cmd/basic -i "7" -e "10 INPUT N:20 PRINT N:END"`: inline program with inputs.
- `make coverage`: combined coverage (packages + acceptance); prints summary.
- `make coverage-html`: writes `combined_coverage.html` for browsing.

## Coding Style & Naming Conventions
- Go 1.24.x; format with `gofmt` (tabs) and organize imports.
- Names: packages lower-case; exported `CamelCase`; locals `camelCase`.
- Prefer explicit dependencies; avoid global state.

## Testing Guidelines
- Add acceptance cases under `acceptance/testdata/*.yaml` (clear names, incremental IDs).
- Place unit tests alongside code in `*_test.go`; use `testing` and `testify/require` where helpful.
- Aim to keep or increase combined coverage (`make coverage`).

## ATDD + TDD Workflow (STRICTLY FOLLOW THIS)

### Start with Acceptance Test (ATDD)
1. Write failing acceptance test that defines complete user-facing behavior
2. Identify components needed to make acceptance test pass

### Then TDD Cycle for Components
3. Write failing unit test for component functionality
4. Write just enough code to make it compile
5. Run test to confirm it fails as expected
6. Write simplest code to make test pass
7. Run test to confirm success
8. Refactor while keeping tests green
9. Return to acceptance test - run to check end-to-end progress
10. Repeat TDD cycle (3-8) until acceptance test passes

## Commit & Pull Request Guidelines
- Commits: imperative, concise; optional prefixes `feat:`, `fix:`, `chore:` (e.g., `feat: add DEF FN user functions`).
- PRs: clear description with linked issues; sample I/O if relevant; tests green (`make test`); update `README.md`/`spec.md` when behavior changes.

## Security & Configuration Tips
- No network access required for runtime. Scripts may call `git`, `cloc`, and `go tool cover`; install them to use history/metrics commands.
