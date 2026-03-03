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

The web interface runs at `http://localhost:8989` (or the configured port).

### Pages

| Route | Description |
|---|---|
| `/` | Landing page |
| `/signup` | Create an account (email + password) |
| `/login` | Sign in |
| `/verify-email` | Email verification landing (linked from signup email) |
| `/forgot-password` | Request a password reset email |
| `/reset-password` | Choose a new password (linked from reset email) |
| `/dashboard` | Browse today's deals; manage keywords; link Telegram |

### Sign-up flow

1. Register at `/signup` — a verification email is sent immediately.
2. Click the link in the email to verify your address and log in.
3. Accounts that have not been verified cannot sign in.

> **Existing accounts in the database** (created before email verification was added) have `email_verified = 0` and will be blocked at login. To unblock them run:
> ```sql
> UPDATE web_users SET email_verified = 1;
> ```

### Linking Telegram

1. Sign up / log in on the web dashboard.
2. Click **Link Telegram Account** in the sidebar.
3. A deep link button appears — click it to open the Telegram app.
4. The bot receives your `/start <token>` and links the accounts automatically.
5. The dashboard updates within a few seconds.

## Email (SMTP)

Email is used for two flows:

| Flow | Trigger | Link destination |
|---|---|---|
| Email verification | New account registration | `/verify-email?token=…` |
| Password reset | Forgot password form | `/reset-password?token=…` |

### Option A — Mailpit (default, self-hosted)

[Mailpit](https://github.com/axllent/mailpit) is a lightweight SMTP relay bundled in `docker-compose.yaml`. It catches all outgoing emails and displays them in a web UI — no real email is delivered, which is ideal for self-hosted or development use.

- SMTP inbox web UI: **http://localhost:8025**
- No credentials needed — the `kramerbot` container talks to `mailpit` over the internal Docker network.
- Emails are stored in `./data/mailpit/mailpit.db` and survive container restarts (up to 500 messages).

After `docker compose up -d`, register an account and then check **http://localhost:8025** to find the verification email and click the link.

### Option B — Real SMTP provider

To send real emails (Gmail, SendGrid, Postmark, etc.), override the SMTP env vars in `kramerbot.env` or the `environment:` block of `docker-compose.yaml`. Remove or comment out the Mailpit overrides first:

```bash
# in kramerbot.env (or docker-compose.yaml environment:)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=you@gmail.com
SMTP_PASS=your-app-password
SMTP_FROM=KramerBot <you@gmail.com>
```

Also update `api.web_url` in `config.yaml` to your public domain so links in emails point to the right place:

```yaml
api:
  web_url: "https://yourdomain.com"
```

### Disabling email entirely

Leave `SMTP_HOST` empty (the default in `config.yaml`). The bot will skip sending emails and instead log the verification/reset links to the container console. This is useful for testing without any SMTP setup.

```
# Find links in logs:
docker logs kramerbot | grep "verify\|reset"
```

## REST API

All endpoints are prefixed with `/api/v1`.

### Auth (public)

```
POST /api/v1/auth/register          — Create account { email, password, display_name } → 202
GET  /api/v1/auth/verify-email      — Verify email ?token=… → JWT
POST /api/v1/auth/login             — Login { email, password } → JWT
POST /api/v1/auth/logout            — Logout (client discards token)
POST /api/v1/auth/forgot-password   — Send reset email { email }
POST /api/v1/auth/reset-password    — Set new password { token, password }
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
GET /api/v1/deals/ozbargain   ?type=good|super   &limit=50 &offset=0
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

# SMTP — set these to use a real mail provider (see "Email" section above)
SMTP_HOST=
SMTP_PORT=587
SMTP_USER=
SMTP_PASS=
SMTP_FROM=KramerBot <noreply@yourdomain.com>
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

When running locally without Docker, SMTP is not configured by default. Verification/reset links are printed to the terminal so you can copy them directly into your browser.

### Using Docker Compose (recommended)

1. Edit `kramerbot.env` and set `TELEGRAM_BOT_TOKEN`, `TELEGRAM_BOT_USERNAME`, and `JWT_SECRET`.
2. Run:
```bash
docker compose up -d
```

Services started:

| Service | Port | Purpose |
|---|---|---|
| `kramerbot` | 8989 | Web UI + API + Telegram bot |
| `mailpit` | 8025 | Email inbox UI (view sent emails) |

3. Open **http://localhost:8989** for the web app and **http://localhost:8025** for the email inbox.

### Using Docker directly

```bash
mkdir -p data

docker run -d --name kramerbot \
  --env-file ./kramerbot.env \
  -v "$(pwd)/data:/app/data" \
  -p 8989:8080 \
  --restart unless-stopped \
  kramerbot:latest
```

Note: when running without Docker Compose, wire up your own SMTP server via the `SMTP_*` env vars.

### Setup Database (SQLite)

The bot auto-creates `data/users.db` on first run (including the `web_users` table with email verification columns).
No manual migration is needed.

<img src="https://raw.githubusercontent.com/intothevoid/kramerbot/main/static/about.jpeg" width="50%" height="50%"></img>
