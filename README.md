# Kramer's Deals

### @kramerbot

https://t.me/kramerbot

A Telegram bot — and now a full web app — to get you the latest deals from https://www.ozbargain.com.au and https://amazon.com.au. Let Kramer watch deals so you don't have to. Giddy up!

## Features

1. Uses Telegram Bot API for instant notifications
2. Written in Go; deployable as a single binary or Docker container
3. **Web UI** — sign up, manage preferences, browse deals, and link your Telegram account from a browser
4. Subscribe to good deals, super deals or set up keyword watches via Telegram commands or the web dashboard
5. User data is written to a SQLite database file (`data/users.db` by default)
6. Keep track of deals already sent to avoid duplicate notifications
7. Supports scraping www.ozbargain.com.au — Good (25+ votes) and Super (50+ votes) deals
8. Supports scraping www.amazon.com.au (via Camel Camel Camel RSS) — Top daily and weekly deals
9. Supports Android TV notifications (via Pipup)
10. Admin announcement broadcast

## Web UI

The web interface runs at `http://localhost:8080` (or the configured port).

### Pages

| Route | Description |
|---|---|
| `/` | Landing page |
| `/signup` | Create an account (email + password) |
| `/login` | Sign in |
| `/dashboard` | Browse today's deals; manage keywords; link Telegram |

### Linking Telegram

1. Sign up / log in on the web dashboard.
2. Click **Link Telegram Account** in the sidebar.
3. A deep link button appears — click it to open the Telegram app.
4. The bot receives your `/start <token>` and links the accounts automatically.
5. The dashboard updates within a few seconds.

## REST API

All endpoints are prefixed with `/api/v1`.

### Auth (public)

```
POST /api/v1/auth/register   — Create account { email, password, display_name }
POST /api/v1/auth/login      — Login { email, password } → JWT
POST /api/v1/auth/logout     — Logout (client discards token)
```

### User (requires Bearer JWT)

```
GET    /api/v1/user/profile             — Current user profile
PUT    /api/v1/user/preferences         — Update deal toggles
GET    /api/v1/user/keywords            — List keywords
POST   /api/v1/user/keywords            — Add keyword { keyword }
DELETE /api/v1/user/keywords/:keyword   — Remove keyword
POST   /api/v1/user/telegram/link       — Generate deep link token
GET    /api/v1/user/telegram/status     — Linked status
DELETE /api/v1/user/telegram/link       — Unlink Telegram
```

### Deals (requires Bearer JWT)

```
GET /api/v1/deals/ozbargain   ?type=good|super  &limit=50 &offset=0
GET /api/v1/deals/amazon      ?type=daily|weekly &limit=50 &offset=0
GET /api/v1/deals             — Combined feed
```

## Deployment

Configuration is primarily managed via `config.yaml`. Sensitive values must be set via environment variables.

### Environment variables

```
TELEGRAM_BOT_TOKEN=<token>           # Mandatory for bot
TELEGRAM_BOT_USERNAME=<username>     # Used in deep link URL (without @)
KRAMERBOT_ADMIN_PASS=<password>      # Optional — admin commands
SQLITE_DB_PATH=<path>                # Optional — defaults to data/users.db
JWT_SECRET=<random_32_byte_hex>      # Mandatory for web API in production
```

Generate a JWT secret:
```bash
openssl rand -hex 32
```

### Run locally

```bash
# Backend
go build .
JWT_SECRET=changeme TELEGRAM_BOT_TOKEN=<token> ./kramerbot

# Frontend (separate terminal)
cd frontend && npm run dev
```

### Using Docker Compose (recommended)

1. Copy `kramerbot.env` and fill in your values.
2. Run:
```bash
docker compose up -d
```

The web UI and API are available at `http://localhost:8080`.

### Using Docker directly

```bash
mkdir -p data

docker run -d --name kramerbot \
  --env-file ./kramerbot.env \
  -v "$(pwd)/data:/app/data" \
  -p 8080:8080 \
  --restart unless-stopped \
  kramerbot:latest
```

### Setup Database (SQLite)

The bot auto-creates `data/users.db` on first run (including the new `web_users` table).
No manual migration is needed.

<img src="https://raw.githubusercontent.com/intothevoid/kramerbot/main/static/about.jpeg" width="50%" height="50%"></img>
