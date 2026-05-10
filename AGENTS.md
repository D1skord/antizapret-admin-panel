# AGENTS.md

Project context for AI coding assistants (Claude Code, Gemini, Cursor, etc.).

## Project Overview

Admin panel for managing VPN clients on an AntiZapret server (OpenVPN/WireGuard). Go backend + Vue 3 frontend compiled into a **single binary** that embeds the built frontend via `go:embed`.

## Development (Docker)

All development is done via Docker Compose — no Go or Node.js needed locally.

```bash
make dev           # start backend :8080 + frontend :5173
make docker-build  # rebuild images after Dockerfile changes
make build         # build production binary (Linux amd64)
```

- Backend hot-reload via `air`: http://localhost:8080
- Frontend Vite dev server: http://localhost:5173 (proxies `/api` → backend)

### Mock filesystem

The backend reads VPN configs from `mock_fs/` by default (configured in `.env`). Required directories:

```
mock_fs/root/antizapret/client/openvpn/vpn-udp/        ← OPENVPN_CLIENTS_PATH
mock_fs/root/antizapret/client/openvpn/antizapret-udp/ ← OPENVPN_ANTIZAPRET_PATH
mock_fs/etc/openvpn/easyrsa3/pki/issued/               ← cert files
```

Client config files naming pattern: `vpn-<name>-(<host>)-udp.ovpn` / `antizapret-<name>-(<host>)-udp.ovpn`.

## Architecture

```
main.go                  ← entry point, DI wiring, embeds frontend/dist via go:embed
internal/
  api/handlers.go        ← Gin HTTP handlers (ClientHandler + auth/download funcs)
  middleware/auth.go     ← X-Auth-Token header check against ADMIN_PASSWORD env var
  service/client.go      ← business logic layer
  repository/client.go   ← filesystem access (reads .ovpn files, runs client.sh)
  entity/client.go       ← Client and PaginatedClients structs
frontend/src/
  api/index.js           ← axios instance with X-Auth-Token interceptor + 401 handler
  stores/auth.js         ← Pinia auth store
  router/index.js        ← Vue Router
  views/                 ← page-level components
  components/            ← reusable UI components
```

**Key design points:**

- Clients are discovered by scanning `.ovpn` files in `OPENVPN_CLIENTS_PATH`. Create/delete delegates to `client.sh` shell script.
- Auth: login returns the password as `token`, subsequent requests send it as `X-Auth-Token` header. The axios interceptor in `api/index.js` attaches it automatically.
- File downloads use short-lived in-memory tokens (5 min TTL): `GET /api/clients/:id/qr-token` → `GET /api/download/:token`.
- The placeholder `frontend/dist/index.html` is tracked in git so `go:embed` never fails without a frontend build.

## Environment Variables (.env)

| Variable | Default (dev) | Description |
|---|---|---|
| `ADMIN_USERNAME` | `admin` | Login username |
| `ADMIN_PASSWORD` | `password` | Login password / auth token |
| `OPENVPN_CLIENTS_PATH` | `mock_fs/root/antizapret/client/openvpn/vpn-udp/` | Directory scanned for VPN client configs |
| `OPENVPN_ANTIZAPRET_PATH` | `mock_fs/root/antizapret/client/openvpn/antizapret-udp/` | Directory for AntiZapret configs |
| `CLIENT_SCRIPT_PATH` | `./mock_fs/root/antizapret/client.sh` | Shell script for client create/delete |
| `PORT` | `8080` | Backend listen port |

## Release

GitHub Actions (`.github/workflows/release.yml`) triggers on `v*` tags: builds frontend, compiles Go binary for Linux amd64, packages into `.tar.gz` with `deploy/` files. Production install via `curl | bash deploy/install.sh` sets up a systemd service.
