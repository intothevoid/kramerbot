#!/bin/bash

docker run -d --rm --name kramerbot --network mongo-network --env-file=token.env -p 3179:3179 kramerbot:latest
