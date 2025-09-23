# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go monorepo containing multiple command-line utilities. Each tool is located in `cmd/<tool-name>/` and can be built independently. The project uses vendor directory for dependency management and follows standard Go project layout.

## Development Commands

### Building
- Build all tools: `go build ./cmd/...`
- Build specific tool: `go build ./cmd/<tool-name>`
- Build with output binary: `go build -o <binary-name> ./cmd/<tool-name>`

### Linting and Quality
- Run linter: `golangci-lint run`
- Run linter with timeout: `golangci-lint run --timeout 5m`
- Auto-fix issues: `golangci-lint run --fix`

The project uses comprehensive linting with `.golangci.yaml` configuration that enables most linters with specific exclusions for `depguard`, `exhaustruct`, `gochecknoglobals`, `gochecknoinits`, and `nonamedreturns`.

### Testing
- Run tests: `go test ./...`
- Run tests with coverage: `go test -cover ./...`
- Run tests with race detection: `go test -race ./...`

Note: The project currently has no test files (`*_test.go`), so adding tests for new functionality is recommended.

### Module Management
- Update dependencies: `go mod tidy`
- Update vendor: `go mod vendor`
- Check for updates: `go list -u -m all`

### Release Management
The project uses GoReleaser (`.goreleaser.yaml`) for building cross-platform binaries. GPG signing is configured with key `F57F85FC7975F22BBC3F25049C173EB1B531AA1F`.

## Architecture

### Project Structure
- `cmd/`: Contains individual CLI tools, each in their own subdirectory
  - Each tool has a main package and can be built independently
  - Tools include: bpa, game-of-life, me-site, mtgdsgenerator, vk2tg, vkphotosdownloader, etc.
- `internal/`: Internal packages shared between tools
  - `internal/pkg/vk2tg/`: VK to Telegram forwarding functionality
- `vendor/`: Vendored dependencies
- `deployments/`: Deployment configurations
- `build/`: Build-related files

### Key Dependencies
- Cobra (`github.com/spf13/cobra`): CLI framework
- Viper (`github.com/spf13/viper`): Configuration management
- Echo (`github.com/labstack/echo/v4`): Web framework
- Redis client (`github.com/redis/go-redis/v9`)
- VK SDK (`github.com/SevereCloud/vksdk/v3`)
- Telegram bot API (`gopkg.in/telebot.v4`)

### Tools Overview
Based on README.md, key tools include:
- `me-site`: Personal website
- `mtgdsgenerator`: Magic: The Gathering dataset generator
- `vk2tg`: VK to Telegram wall forwarder
- `vkphotosdownloader`: VK photo downloader
- `a200`: Simple nginx-like server responding 200 to all requests

### CI/CD
- GitHub Actions workflow for Go linting (`.github/workflows/lint-go.yaml`)
- Separate workflows for specific tools (a200, me-site, vk2tg)
- Linting runs on any Go file changes (`.go`, `.mod`, `.sum`)
- GoReleaser configuration for multi-platform releases

## Development Notes

### Code Style
The project follows standard Go conventions with strict linting rules. Pay attention to:
- Variable naming length requirements (min 3 characters, configured in golangci)
- Complexity limits (cyclomatic complexity max 10)
- Import formatting with goimports and gofumpt
- Error handling patterns using `github.com/cockroachdb/errors`

### Adding New Tools
When adding a new CLI tool:
1. Create directory under `cmd/<tool-name>/`
2. Follow existing patterns (main package, cobra commands if applicable)
3. Update go.mod if new dependencies are needed
4. Run `go mod vendor` to update vendor directory
5. Ensure linting passes with `golangci-lint run`
