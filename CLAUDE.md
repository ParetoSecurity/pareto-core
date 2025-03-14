# Pareto Security Agent Development Guide

## Build & Test Commands
- Build: `go build ./cmd/paretosecurity`
- Test all: `go test ./...`
- Test specific package: `go test github.com/ParetoSecurity/agent/checks/linux`
- Test single test: `go test -run TestApplicationUpdates_Run ./checks/linux`
- Coverage: `go test -coverprofile=coverage.txt ./...`
- Lint: Uses pre-commit hooks with `alejandra` and `gofmt`

## Code Style Guidelines
- Imports: Standard library first, third-party packages second, project-specific last
- Formatting: Use `gofmt` standard formatting; tabs for indentation
- Naming: CamelCase for exported identifiers, camelCase for unexported
- Interfaces: Prefer small, focused interfaces with clear purpose
- Tests: Table-driven tests with descriptive names, using assert package
- Error handling: Return errors up the stack, use early returns with `if err != nil`
- Logging: Use `log.WithField/WithError` for structured logging
- Mocks: Use dependency injection patterns with interface mocks for testing
- Documentation: Add descriptive comments for all exported functions/methods
- File layout: Related types and their implementations together