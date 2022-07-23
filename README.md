# Kramer's Deals

### @kramerbot

https://t.me/kramerbot

A Telegram bot to get you the latest deals from websites like https://www.ozbargain.com.au and https://amazon.com.au. Let Kramer watch deals so you don't have to. Giddy up!

## Features

1. Uses Telegram Bot API for instant notifications
2. Written in Go and can be deployed with a single binary (Dockerfile included)
3. Subscribe to good deals, super deals or setup your own custom deals by watching specific keywords
4. User data is written to a Sqlite database for easy migration
5. Keep track of deals already sent to avoid duplicate sending
6. Supports scraping www.ozbargain.com.au - Good and super deals
7. Supports scraping www.amazon.com.au (via Camel Camel Camel RSS) - Top daily and weekly deals
8. Supports Android TV notifications
9. API to access user and deal data
10. Ability to send maintenance messages / announcements to all users

## API

The following API endpoints are available -

```
/users - Get user data for all users
/users/:chatid - Get user data by chat id
/deals - Get deal data for latest deals by the scraper
```

## Deployment

You must have the required environment variables for Kramerbot to function correctly. See section 'Required environment variables' for more details. Kramerbot can be deployed using the foll. command, after required environment variables have been set -

```
go build .
./kramerbot
```

### Required environment variables

```
TELEGRAM_TOKEN_API=<your_token>
GIN_MODE=release
KRAMERBOT_ADMIN_PASS=<your_admin_password>
```

### Using Docker

To build a Docker image of Kramerbot use the command -

```
sudo docker build -t kramerbot:latest .
```

Create a token.env file with your Telegram API token (used in step below) -

```
TELEGRAM_TOKEN_API=<your_token>
GIN_MODE=release
KRAMERBOT_ADMIN_PASS=<your_admin_password>
```

To deploy your container, use the command -

```
sudo docker run -d --rm --name kramerbot --env-file=token.env -p 8080:8080 kramerbot:latest
```

<img src="https://raw.githubusercontent.com/intothevoid/kramerbot/main/static/about.jpeg" width="50%" height="50%"></img>
