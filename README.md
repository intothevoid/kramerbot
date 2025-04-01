# Kramer's Deals

### @kramerbot

https://t.me/kramerbot

A Telegram bot to get you the latest deals from websites like https://www.ozbargain.com.au and https://amazon.com.au. Let Kramer watch deals so you don't have to. Giddy up!

**Note:** This version includes a Telegram web app interface for users to manage their preferences.

## Features

1. Uses Telegram Bot API for instant notifications
2. Written in Go and can be deployed with a single binary (Dockerfile included)
3. Subscribe to good deals, super deals or setup your own custom deals by watching specific keywords via Telegram commands
4. User data is written to a SQLite database file (`data/users.db` by default)
5. Keep track of deals already sent to avoid duplicate sending
6. Supports scraping www.ozbargain.com.au - Good and super deals
7. Supports scraping www.amazon.com.au (via Camel Camel Camel RSS) - Top daily and weekly deals
8. Supports Android TV notifications (via Pipup)
9. Ability to send maintenance messages / announcements to all users (if admin commands are implemented)
10. Includes a web app that allows users to manage preferences directly from Telegram

## Web App

The bot includes a Telegram Web App that allows users to manage their preferences directly from within Telegram. Users can:
- Toggle notification settings for OzBargain and Amazon deals
- Add and remove keywords to watch for
- Send test notifications

### Web App Configuration

The web app is configured in the `config.yaml` file:

```yaml
web_app:
  enabled: true # Set to true to enable the web app server
  listen_address: "0.0.0.0:8085" # Address and port for the web server
  base_url: "" # Publicly accessible base URL of the web app (e.g., https://yourdomain.com/kramerbot-app)
  development_mode: false # Set to true during development to bypass authentication and use dummy data
```

#### Development Mode

The web app includes a development mode that makes local testing easier:
- When `development_mode` is set to `true`, the app bypasses Telegram authentication and uses dummy data
- This allows you to test the web app without needing a real Telegram connection
- Set `development_mode` to `false` in production to enforce proper Telegram authentication

## API

The following API endpoints are available -

```
/users - Get user data for all users
/users/:chatid - Get user data by chat id
/deals - Get deal data for latest deals by the scraper
/signup - Signup from accompanying web app https://www.github.com/intothevoid/kramerbotui
/preferences - Update user preferences
/authenticate - User authentication
```

## Deployment

Configuration is primarily managed via `config.yaml`. However, sensitive information like tokens should be set via environment variables. Kramerbot can be deployed using the following command after setting the required environment variables:

```
go build .
./kramerbot
```

### Required environment variables

These environment variables override values in `config.yaml` if set.

```
TELEGRAM_BOT_TOKEN=<your_telegram_bot_token> # Mandatory
KRAMERBOT_ADMIN_PASS=<your_admin_password> # Optional: If admin commands are used
SQLITE_DB_PATH=<path_to_your_sqlite_db_file> # Optional: Defaults to value in config.yaml or 'data/users.db'
```
*(Refer to `config.yaml` for other configuration options like logging, scraper intervals, etc.)*

### Setup Database (SQLite)

The bot uses a SQLite database file to store user data. By default, it will create/use a file at `./data/users.db` relative to where the bot is run.
- Ensure the directory `./data` exists and the bot has write permissions.
- You can change the path using the `sqlite.db_path` setting in `config.yaml` or the `SQLITE_DB_PATH` environment variable.

### Using Docker

To build a Docker image of Kramerbot:

```
sudo docker build -t kramerbot:latest .
```

Create an environment file (e.g., `kramerbot.env`) with your required variables:

```
TELEGRAM_BOT_TOKEN=<your_telegram_bot_token>
KRAMERBOT_ADMIN_PASS=<your_admin_password> # Optional
SQLITE_DB_PATH=/app/data/users.db # Optional: Specify path inside the container
```

To deploy your container using Docker Compose (recommended for persisting data):

1.  Make sure you have a `docker-compose.yaml` file similar to the one provided in the repository (it should handle mounting the `./data` directory).
2.  Run: `docker compose up -d`

Alternatively, to run directly with `docker run`:

```bash
# Create the data directory on your host first if it doesn't exist
mkdir -p data

# Run the container, mounting the local data directory
sudo docker run -d --name kramerbot \
  --env-file ./kramerbot.env \
  -v "$(pwd)/data:/app/data" \
  --restart unless-stopped \
  kramerbot:latest
```*(This mounts your local `./data` directory into `/app/data` inside the container, where the bot expects to find the SQLite file by default or via the environment variable.)*

<img src="https://raw.githubusercontent.com/intothevoid/kramerbot/main/static/about.jpeg" width="50%" height="50%"></img>

