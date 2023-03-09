#!/bin/bash

docker run -d --rm --name kramerbot --network mongo-network --env-file=../token.env -p 8080:8080 kramerbot:latest
