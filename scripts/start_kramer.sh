#!/bin/bash

docker run -d --name kramerbot --network mongo-network --env-file=../token.env -p 8080:8080 --restart unless-stopped kramerbot:latest
