# 🛠️ Tools

[![Go](https://img.shields.io/badge/Go-1.25.1-00ADD8?style=flat&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/License-BSD--3--Clause-blue.svg)](LICENSE)
[![CI](https://github.com/lexfrei/tools/workflows/Lint%20Go/badge.svg)](https://github.com/lexfrei/tools/actions)

> A collection of useful tools and utilities by [Aleksei Sviridkin](https://github.com/lexfrei)

This monorepo contains various Go-based CLI tools and containerized services for different purposes - from personal website to VK/Telegram integration and development utilities.

## 🚀 Tools Overview

### Web Services & APIs

| Tool | Description | Container | Status |
|------|-------------|-----------|--------|
| **[me-site](cmd/me-site/)** | Personal website with static content | `ghcr.io/lexfrei/me-site` | ✅ Active |
| **[a200](build/a200/)** | Simple nginx server that responds 200 to all requests | `ghcr.io/lexfrei/a200` | ✅ Active |
| **[ow-exporter](cmd/ow-exporter/)** | Overwatch 2 statistics exporter for Prometheus | `ghcr.io/lexfrei/ow-exporter` | 🚧 Development |
| **[redis-ui](cmd/redis-ui/)** | Web UI for Redis database management | - | 🔧 Development |

### Social & Communication

| Tool | Description | Container | Status |
|------|-------------|-----------|--------|
| **[vk2tg](cmd/vk2tg/)** | VK wall posts forwarder to Telegram channels | `ghcr.io/lexfrei/vk2tg` | ✅ Active |
| **[vkphotosdownloader](cmd/vkphotosdownloader/)** | Download photos from VK albums | - | ✅ Active |

### Gaming & Entertainment

| Tool | Description | Container | Status |
|------|-------------|-----------|--------|
| **[mtgdsgenerator](cmd/mtgdsgenerator/)** | Magic: The Gathering dataset generator for ML training | - | ✅ Active |
| **[game-of-life](cmd/game-of-life/)** | Conway's Game of Life implementation | - | 🎮 Demo |

### Development & Utilities

| Tool | Description | Container | Status |
|------|-------------|-----------|--------|
| **[redis-checker](cmd/redis-checker/)** | Redis connection and health checker | - | 🔧 Utility |
| **[russian-mobile-sha256](cmd/russian-mobile-sha256/)** | SHA256 hash generator for Russian mobile numbers | - | 🔧 Utility |
| **[bpa](cmd/bpa/)** | Bulk processing automation tool | - | 🔧 Utility |

### Examples & Learning

| Tool | Description | Status |
|------|-------------|--------|
| **[errors-example](cmd/errors-example/)** | Go error handling examples | 📚 Example |
| **[rand-example](cmd/rand-example/)** | Random number generation examples | 📚 Example |
| **[tab-example](cmd/tab-example/)** | Tab completion examples | 📚 Example |

## 🏗️ Development

### Prerequisites

- **Go 1.25+** (required)
- **Docker** (for container builds)
- **golangci-lint** (for linting)

### Building

```bash
# Build all tools
go build ./cmd/...

# Build specific tool
go build ./cmd/me-site

# Build with custom output name
go build -o mysite ./cmd/me-site
```

### Development Commands

```bash
# Run linter
golangci-lint run

# Run tests
go test ./...

# Update dependencies
go mod tidy && go mod vendor

# Build containers
docker build -f build/me-site/Containerfile -t me-site .
```

### Project Structure

```text
tools/
├── cmd/                    # CLI applications
│   ├── me-site/           # Personal website
│   ├── vk2tg/             # VK to Telegram forwarder
│   └── ...                # Other tools
├── build/                 # Container configurations
│   ├── me-site/           # Website container
│   ├── vk2tg/             # VK2TG container
│   └── a200/              # Simple 200 server
├── internal/              # Internal packages
├── vendor/                # Vendored dependencies
└── .github/workflows/     # CI/CD pipelines
```

## 🐳 Container Images

All container images are built automatically and available at GitHub Container Registry.

### Running Containers

```bash
# Personal website
docker run -p 8080:8080 ghcr.io/lexfrei/me-site:latest

# VK to Telegram forwarder (requires config)
docker run -e VK_TOKEN=xxx -e TG_TOKEN=xxx ghcr.io/lexfrei/vk2tg:latest

# Simple 200 server
docker run -p 8080:80 ghcr.io/lexfrei/a200:latest
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Run linter and tests (`golangci-lint run && go test ./...`)
5. Commit your changes (`git commit -m 'feat: add amazing feature'`)
6. Push to the branch (`git push origin feature/amazing-feature`)
7. Open a Pull Request

### Commit Convention

This project follows [Semantic Commit Messages](https://www.conventionalcommits.org/):

- `feat:` - New features
- `fix:` - Bug fixes
- `docs:` - Documentation changes
- `style:` - Code style changes
- `refactor:` - Code refactoring
- `test:` - Adding tests
- `chore:` - Maintenance tasks

## 📝 License

This project is licensed under the BSD 3-Clause License - see the [LICENSE](LICENSE) file for details.

## 👤 Author

Aleksei Sviridkin

- Email: <f@lex.la>
- GitHub: [@lexfrei](https://github.com/lexfrei)
- Role: Developer

---

Made with ❤️ and Go
