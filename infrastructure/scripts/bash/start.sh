#!/bin/sh

# This file is used to start up all services in containers

# move to /infrastructure directory
cd ../../

# docker-compose
docker-compose -f docker-compose.yaml -p server_ecom up --remove-orphans -d
