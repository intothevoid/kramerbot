# Kramer's Deals

### @kramerbot

https://t.me/kramerbot

A Telegram bot to get you the latest deals from websites like https://www.ozbargain.com.au and https://amazon.com.au. Let Kramer watch deals so you don't have to. Giddy up!

**Note:** This version is CLI-only and does not include a web interface or API.

## Features

1. Uses Telegram Bot API for instant notifications
2. Written in Go and can be deployed with a single binary (Dockerfile included)
3. Subscribe to good deals, super deals or setup your own custom deals by watching specific keywords via Telegram commands
4. User data is written to a Mongo NoSQL database
5. Keep track of deals already sent to avoid duplicate sending
6. Supports scraping www.ozbargain.com.au - Good and super deals
7. Supports scraping www.amazon.com.au (via Camel Camel Camel RSS) - Top daily and weekly deals
8. Supports Android TV notifications (via Pipup)
9. Ability to send maintenance messages / announcements to all users (if admin commands are implemented)

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
MONGO_URI=<your_mongodb_connection_string> # Optional: Defaults to value in config.yaml
MONGO_DBNAME=<your_database_name> # Optional: Defaults to value in config.yaml
MONGO_COLLNAME=<your_collection_name> # Optional: Defaults to value in config.yaml
```
*(Refer to `config.yaml` for other configuration options like logging, scraper intervals, etc.)*

### Setup MongoDB

Ensure you have a running MongoDB instance accessible by the bot. The connection details should be set either in `config.yaml` or via the `MONGO_*` environment variables.

You can run MongoDB locally using Docker:
```
# Pull the image (if needed)
sudo docker pull mongo:4.4.18

# Start the container (example)
cd scripts
sudo ./start_mongo.sh
```
*(The `start_mongo.sh` script sets up a network and volume. Ensure your `MONGO_URI` points correctly to this instance, e.g., `mongodb://kramer-mongo:27017` if running kramerbot in the same Docker network)*

### Using Docker

To build a Docker image of Kramerbot:

```
sudo docker build -t kramerbot:latest .
```

Create an environment file (e.g., `kramerbot.env`) with your required variables:

```
TELEGRAM_BOT_TOKEN=<your_telegram_bot_token>
KRAMERBOT_ADMIN_PASS=<your_admin_password> # Optional
MONGO_URI=<your_mongodb_connection_string> # Or configure in config.yaml and mount it
MONGO_DBNAME=<your_database_name> # Or configure in config.yaml
MONGO_COLLNAME=<your_collection_name> # Or configure in config.yaml
```

To deploy your container (example using the network from `start_mongo.sh`):

```bash
# Ensure the mongo-network exists (created by start_mongo.sh)
sudo docker run -d --name kramerbot --network mongo-network --env-file ./kramerbot.env --restart unless-stopped kramerbot:latest
```

<img src="https://raw.githubusercontent.com/intothevoid/kramerbot/main/static/about.jpeg" width="50%" height="50%"></img>
