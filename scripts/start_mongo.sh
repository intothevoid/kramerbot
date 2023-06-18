#!/bin/bash

docker network create mongo-network
docker run -d --network mongo-network --name kramer-mongo -p 27017:27017 -v mongo-data:/data/db --restart unless-stopped mongo:4.4.18
