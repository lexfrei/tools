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

## Go Code Style Guide

Based on golangci-lint configuration and current linting errors, follow these specific style guidelines:

### Standard Libraries (from ~/PROMPT.md)
- **Logging**: `slog` (standard library)
- **Errors**: `github.com/cockroachdb/errors`
- **Web framework**: `github.com/labstack/echo/v4`
- **CLI**: `github.com/spf13/cobra`
- **Configuration**: `github.com/spf13/viper`

### Function and Method Guidelines
- **Function length**: Maximum 60 lines (funlen)
- **Cognitive complexity**: Maximum 30 (gocognit)
- **Cyclomatic complexity**: Maximum 10 (gocyclo)
- Break down complex functions into smaller, focused functions
- Use helper methods to reduce complexity

### Line Length and Formatting
- **Maximum line length**: 120 characters (lll)
- For struct tags that exceed line length:
  ```go
  // Good - break long struct tags across lines
  type HeroStats struct {
      TimePlayed time.Duration `ow:"time_played"
          prometheus:"ow_hero_time_played_seconds"
          help:"Total time played on hero"
          path:"[data-category-id='0x0860000000000021']"
          type:"duration"`
  }
  ```

### Constants and Magic Numbers
- **Extract repeated strings** as constants (goconst):
  ```go
  // Good
  const MouseKeyboardViewActiveSelector = ".mouseKeyboard-view.is-active"
  const QuickPlayViewActiveSelector = ".quickPlay-view.is-active"
  ```
- **Avoid magic numbers** (mnd) - define as constants with meaningful names:
  ```go
  // Bad
  time.Sleep(30 * time.Second)

  // Good
  const DefaultTimeout = 30 * time.Second
  time.Sleep(DefaultTimeout)
  ```

### Variable Naming
- **Minimum 3 characters** for variable names (varnamelen)
- Use descriptive names:
  ```go
  // Bad
  e := echo.New()
  s := "hello"

  // Good
  server := echo.New()
  message := "hello"
  ```

### Documentation Standards
- **All comments must end with periods** (godot):
  ```go
  // Good comment ends with a period.
  func doSomething() {}
  ```
- **Document all exported functions and types** (godoclint)
- Use proper Go doc comment format

### Struct Tags and JSON Naming
- **Use camelCase for JSON tags** (tagliatelle):
  ```go
  // Good
  type Player struct {
      BattleTag    string     `json:"battleTag" yaml:"battletag"`
      LastResolved *time.Time `json:"lastResolved" yaml:"lastResolved"`
  }
  ```
- **Align struct tags** (tagalign) for better readability:
  ```go
  type Config struct {
      Host     string `json:"host"     yaml:"host"`
      Port     int    `json:"port"     yaml:"port"`
      Database string `json:"database" yaml:"database"`
  }
  ```

### Context Handling
- **Pass context as first parameter** (noctx) in functions that might need it:
  ```go
  // Good
  func fetchData(ctx context.Context, url string) error {
      // implementation
  }
  ```

### Loop Patterns
- **Use integer ranges for Go 1.22+** (intrange):
  ```go
  // Modern Go 1.22+ style
  for i := range 10 {
      // process i
  }

  // Instead of
  for i := 0; i < 10; i++ {
      // process i
  }
  ```

### Line Spacing (nlreturn)
- **Add blank lines before return statements** when they follow blocks:
  ```go
  // Good
  if condition {
      doSomething()
  }

  return result
  ```

### Testing Standards
- **Use parallel tests** where appropriate (paralleltest):
  ```go
  func TestSomething(t *testing.T) {
      t.Parallel() // Add this for independent tests

      // test implementation
  }
  ```

### TODOs and Technical Debt
- **Minimize TODO comments** (godox)
- When TODOs are necessary, make them specific and actionable
- Include issue references or deadlines where possible

### Prometheus Metrics
- **Follow Prometheus naming conventions** (promlinter):
  ```go
  // Good metric names
  "http_requests_total"      // counter
  "request_duration_seconds" // histogram
  "current_connections"      // gauge
  ```

### Error Handling Patterns
- Use sentinel errors for expected conditions
- Wrap errors with context using `github.com/cockroachdb/errors`
- Don't return `nil, nil` - use meaningful errors instead

### Code Organization
- Keep related functionality together
- Use meaningful package names
- Prefer composition over inheritance
- Use interfaces for dependencies

### File Permissions and Security
- Use octal notation for file permissions: `0o600` not `0600`
- Never commit secrets or credentials
- Use environment variables for sensitive configuration

### Git Security Standards
- **NEVER disable GPG signing for commits** - GPG signing is mandatory for security
- If GPG signing fails, fix the GPG agent instead of bypassing with `--no-gpg-sign`
- Use `gpgconf --kill all && gpgconf --launch gpg-agent` to restart GPG agent
- All commits must be signed with key `F57F85FC7975F22BBC3F25049C173EB1B531AA1F`

### Adding New Tools
When adding a new CLI tool:
1. Create directory under `cmd/<tool-name>/`
2. Follow existing patterns (main package, cobra commands if applicable)
3. Update go.mod if new dependencies are needed
4. Run `go mod vendor` to update vendor directory
5. Ensure linting passes with `golangci-lint run`
