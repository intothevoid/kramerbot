# Kramer's Deals

### @kramerbot

https://t.me/kramerbot

A Telegram bot to get you the latest deals from websites like https://www.ozbargain.com.au. Let Kramer watch deals so you don't have to. Giddy up!

## Features

1. Uses Telegram Bot API for instant notifications
2. Written in Go and can be deployed with a single binary (Dockerfile included)
3. Subscribe to good deals, super deals or setup your own custom deals by watching specific keywords
4. User data is written to a Sqlite database for easy migration
5. Keep track of deals already sent to avoid duplicate sending
6. Supports scraping www.ozbargain.com.au (more scrapers to be added)
7. Supports Android TV notifications
8. API to access user and deal data

## API

The following API endpoints are available -

```
/users - Get user data for all users
/users/:chatid - Get user data by chat id
/deals - Get deal data for latest deals by the scraper
```

## Deployment

You must have an environment variable called 'TELEGRAM_TOKEN_API=<token>' in your system environment variables. Kramerbot can be deployed using the command -

```
go build .
./kramerbot
```

### Using Docker

To build a Docker image of Kramerbot use the command -

```
sudo docker build -t kramerbot:latest .
```

Create a token.env file with your Telegram API token (used in step below) -

```
TELEGRAM_TOKEN_API=<token>
```

To deploy your container, use the command -

```
sudo docker run -d --rm --name kramerbot --env-file=token.env kramerbot:latest
```

<img src="https://raw.githubusercontent.com/intothevoid/kramerbot/main/static/about.jpeg" width="50%" height="50%"></img>
