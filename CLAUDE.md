# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**Scrubarr** is a Go-based application that tracks and safely deletes media files across *arr instances (Sonarr, Radarr) by linking them with torrent data from torrent clients (Deluge, rTorrent). It provides a web UI with JWT authentication to manage media inventories, evaluate retention policies based on tracker rules, and safely delete files with torrent cleanup.

### Key Use Cases
- Link media files from Sonarr/Radarr with active torrents in Deluge/rTorrent
- Enforce retention policies based on tracker rules (minimum ratio, minimum age)
- Provide a web dashboard to safely review and delete media with associated torrent cleanup
- Support dry-run mode to preview deletions without actually removing files

---

## Build, Test & Run Commands

### Setup
```bash
# Download dependencies
go mod download

# Generate password hash for authentication (interactive)
./main generate-password-hash
```

### Build
```bash
# Build standalone binary (uses goreleaser for snapshot)
goreleaser build --snapshot --clean --single-target

# Build into ./dist/ directory
go build -o ./dist/scrubarr ./cmd/scrubarr
```

### Test
```bash
# Run all tests
go test ./...

# Run a single test package
go test ./pkg/inventory/...

# Run a specific test function
go test -run TestFunctionName ./pkg/inventory/...

# Run with verbose output and race detection
go test -v -race ./...
```

### Run
```bash
# Start the server (requires config file)
./main serve --config ./config.toml

# Start with dry-run mode (no actual deletions)
./main serve --config ./config.toml --dry-run

# Start with caching enabled
./main serve --config ./config.toml --save-cache --use-cache

# Set log level
./main serve --config ./config.toml --level debug
```

### Configuration
- **Config file**: `configs/scrubarr.toml` (example template)
- **Test config**: `test/real_test_config.toml`
- **Cache directory**: Controlled via `SCRUBARR_CACHE_DIR` env var (defaults to `./cache/`)

---

## High-Level Architecture

### Component Interaction Flow

```
┌─────────────────┐
│   Web UI        │  (HTML/JS served from /web)
│  (Port 8888)    │
└────────┬────────┘
         │
         ▼
┌─────────────────────────────────────────┐
│  internal/app/webserver/                │  ◄─ HTTP handlers, JWT auth, routing
│  - handler_media.go (GET/DELETE media)  │
│  - handler_auth.go (login/logout)       │
│  - webserver.go (setup routes)          │
└────────┬────────────────────────────────┘
         │
         ▼
┌────────────────────────────────────────────────┐
│  pkg/inventory/service.go                       │  ◄─ Core orchestrator
│  ┌──────────────────────────────────────────┐  │
│  │ RefreshCache()                           │  │  1. Loads media & torrent data
│  │ GetMediaInventory()                      │  │  2. Links media with torrents
│  │ DeleteMedia()                            │  │  3. Evaluates retention policies
│  │ GetExpandedMediaRow()                    │  │  4. Returns paginated results
│  └──────────────────────────────────────────┘  │
└────────┬────────────────────────────────────────┘
         │
    ┌────┴─────────────────────────────────────┬──────────────────┐
    │                                           │                  │
    ▼                                           ▼                  ▼
┌──────────────────┐  ┌──────────────────┐  ┌─────────────────┐
│  Media Manager   │  │ Torrent Manager  │  │ Linker Service  │
│  (pkg/media/)    │  │ (pkg/torrent     │  │ (pkg/linker/)   │
│                  │  │  clients/)       │  │                 │
│ - Radarr         │  │ - Deluge         │  │ - Matches files │
│ - Sonarr         │  │ - rTorrent       │  │   by path       │
└──────────────────┘  └──────────────────┘  └─────────────────┘
    │                      │
    ▼                      ▼
┌──────────────────┐  ┌──────────────────────┐
│  REST APIs       │  │  Torrent Protocols   │
│ (Sonarr/Radarr)  │  │ (XMLRPC/Deluge API)  │
└──────────────────┘  └──────────────────────┘
```

### Core Packages

#### 1. **pkg/domain/** - Shared Domain Types
   - `media.go` - `MediaEntry`, `MediaFile`, `MediaSourceManager`
   - `torrent.go` - `TorrentEntry`, `TorrentSourceManager`
   - `decision.go` - `Decision` enum (SafeToDelete, Pending, NotAllowed)
   - `trackers.go` - `Tracker` with retention rules (min_ratio, min_age)
   - `cache.go` - `CachedManager` interface for caching

#### 2. **pkg/inventory/** - Main Orchestrator
   - `service.go` - **Core Service**: Manages enriched linked media cache, evaluates retention policy
   - `service_cache.go` - Cache loading/saving from disk
   - `retention_policy.go` - Implements `RetentionPolicy` interface
   - `linker.go` - Interfaces for linking
   - `media_id.go` - ID parsing for media entries

#### 3. **pkg/media/** - Media Source Providers
   - `radarr.go` - Radarr HTTP client (movies)
   - `sonarr.go` - Sonarr HTTP client (TV series)
   - `manager.go` - `DefaultMediaManager`: aggregates Radarr/Sonarr, caches results

#### 4. **pkg/torrentclients/** - Torrent Client Providers
   - `deluge.go` - Deluge XMLRPC client
   - `rtorrent.go` - rTorrent XMLRPC client
   - `manager.go` - `DefaultTorrentManager`: aggregates clients, caches results

#### 5. **pkg/linker/** - Media-Torrent Linking
   - `service.go` - `LinkMedia()` matches files by path (`OriginalFilePath`)

#### 6. **pkg/retentionpolicy/** - Deletion Decision Logic
   - `service.go` - Evaluates each file against tracker rules
   - `tracker_resolver.go` - Maps torrent tracker names to configured tracker policies

#### 7. **pkg/quota/** - Disk Quota Service
   - `ultraapi.go` - Ultra API quota provider (disk space info)

#### 8. **internal/app/webserver/** - HTTP Server
   - `webserver.go` - Routes setup, auth middleware, listener
   - `handler_media.go` - GET /media, PUT /media (refresh), DELETE /media/{id}
   - `handler_auth.go` - POST /login, POST /logout, JWT validation
   - `auth_setup.go` - Auth provider initialization (passwordhash or Jellyfin)
   - `template_cache.go` - Embedded HTML template loading

#### 9. **internal/app/auth/** - Authentication
   - `passwordhash_provider.go` - PBKDF2-based password hash/salt
   - `jellyfin_provider.go` - Jellyfin authentication
   - `jwt.go` - JWT token creation/validation

#### 10. **internal/app/scrubarr/** - Entry Point
   - `scrubarr.go` - Cobra CLI setup, dependency injection, config loading
   - `password_generation.go` - Password hash generation utility
   - `quota_service.go` - Quota provider factory
   - `log_level.go` - Structured logging setup

---

## Key Architectural Patterns

### 1. **Dependency Injection**
All services are initialized in `internal/app/scrubarr/scrubarr.go` (the `serve()` function):
- Media, torrent, linker, and retention policy services are created and passed to `inventory.NewService()`
- Web server receives inventory service via dependency, not global state
- Enables testing with mock implementations

### 2. **Interface-Based Design**
- `domain.MediaSourceManager`, `domain.TorrentSourceManager` - pluggable providers
- `inventory.Linker`, `inventory.RetentionPolicy` - policies are injectable
- Allows swapping implementations (e.g., mock providers in tests)

### 3. **Manager Pattern**
- `DefaultMediaManager`, `DefaultTorrentManager` aggregate multiple sources
- Cache management is centralized (lock-based, thread-safe)
- Supports both in-memory caching and disk persistence

### 4. **Concurrent Refresh**
- `RefreshCache()` spawns goroutines to fetch from Radarr/Sonarr/Deluge/rTorrent in parallel
- Results collected via channels with error aggregation (`errors.Join()`)
- Automatic periodic refresh every `refresh_interval` (from config)

### 5. **Configuration Management**
- **koanf** library handles TOML config loading with dot-notation access
- Config paths: `general.listen_addr`, `connections.radarr.hostname`, `trackers.<name>.min_ratio`
- Environment variable support via `SCRUBARR_CACHE_DIR`

### 6. **Retention Policy Evaluation**
- Each `LinkedMedia` is evaluated to produce `EvaluationReport` with per-file decisions
- Global decision is aggregated from file decisions (any "pending" = media is "pending")
- Files without torrent entries are always safe to delete
- Deletion blocked if tracker rules not met (ratio/age thresholds)

### 7. **Authentication & Authorization**
- JWT-based stateless auth (tokens stored in cookies)
- Two providers: password hash (PBKDF2) and Jellyfin SSO
- All routes except `/login` require valid JWT token
- Username logged in context for audit trails

---

## Existing Documentation

### README.adoc
- Basic installation, config generation, and server startup instructions
- Points to `configs/scrubarr.toml` for configuration template

### GitHub Workflows (.github/workflows/)
- **ci.yml**: Test, build, and release pipeline
  - Runs on every push
  - Tests: `go test ./...`
  - Build: `goreleaser build --snapshot --clean --single-target`
  - Release: `goreleaser release --clean` (on tags only)
- **gitleaks.yml**: Secret detection scan

### Configuration Example (configs/scrubarr.toml)
- Sections: `[general]`, `[general.auth]`, `[quota]`, `[connections.*]`, `[trackers.*]`
- Auth providers: `passwordhash` or `jellyfin`
- Quota provider: currently only `ultraapi`
- Torrent client sections: `[connections.deluge]`, `[connections.rtorrent]`
- Media sections: `[connections.sonarr]`, `[connections.radarr]`
- Tracker rules: `[trackers.<name>]` with pattern, min_ratio, min_age

---

## Testing Strategy

### Test Locations
- Unit tests co-located with source: `*_test.go` in same package
- Key test files:
  - `pkg/inventory/service_test.go` - Core orchestrator logic
  - `pkg/retentionpolicy/service_test.go` - Retention policy evaluation
  - `pkg/linker/service_test.go` - Media-torrent linking
  - `internal/app/webserver/handler_auth_test.go` - Auth endpoints
  - `pkg/ultraapi/api_test.go` - Quota API

### Testing Utilities
- `pkg/util/test_utils.go` - Common test helpers (e.g., `MustParseDate()`)
- `testify/require` for assertions
- Mock torrents and media entries defined inline in tests

### Running Single Tests
```bash
# Run all tests in inventory package
go test ./pkg/inventory/...

# Run specific test
go test -run Test_generateRawFileBasedMediaRow ./pkg/inventory/...

# Run with verbose output
go test -v ./pkg/inventory/service_test.go
```

### Test Configuration
- Uses in-memory mocks for media/torrent managers
- Dry-run mode tested via `--dry-run` CLI flag
- Integration tests use real config from `test/real_test_config.toml`

---

## Configuration Deep Dive

### Key Config Sections

#### [general]
- `listen_network`: "tcp" or "unix" socket
- `listen_addr`: e.g., ":8888" for port 8888
- `path_prefix`: URL path prefix (e.g., "/scrubarr")
- `real_ip_header_name`: For reverse proxy IP logging (e.g., "X-Forwarded-For")
- `refresh_interval`: How often to refresh media/torrent caches (e.g., "1h", "0" to disable)

#### [general.auth.providers.passwordhash]
- `username`: Static username
- `password_salt`: Hex-encoded salt (generate via `generate-password-hash`)
- `password_hash`: Hex-encoded PBKDF2 hash

#### [general.auth.providers.jellyfin]
- `base_url`: Jellyfin instance URL for SSO

#### [quota.ultraapi]
- `endpoint`: Ultra API endpoint for disk quota
- `api_key`: API key for authentication

#### [connections.sonarr], [connections.radarr]
- `enabled`: Boolean to enable/disable
- `hostname`: Base URL of the *arr instance
- `api_key`: API key for authentication

#### [connections.deluge], [connections.rtorrent]
- `enabled`: Boolean
- `hostname`: Host/URL for XMLRPC endpoint
- `port`: Port (for Deluge)
- `username`, `password`: Credentials

#### [trackers.<name>]
- `name`: Display name
- `pattern`: Regex pattern to match tracker in torrent (e.g., `^torrents\\.example\\.com$`)
- `min_ratio`: Minimum upload ratio before safe to delete
- `min_age`: Minimum time torrent must exist (Golang duration, e.g., "720h" = 30 days)

