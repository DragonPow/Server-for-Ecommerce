#!/bin/sh

# This file is used to start up all services in containers

# move to /infrastructure directory
cd ../

# docker-compose
docker-compose -f docker-compose.infras.dev.yaml -p ktpm-infras down
docker-compose -f docker-compose.infras.dev.yaml -p ktpm-infras up -d
