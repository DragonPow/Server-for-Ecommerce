#!/bin/sh

# This file is used to remove all container services

# move to /infrastructure directory
cd ../../

# docker-compose
docker-compose -f docker-compose.infras.dev.yaml -p ktpm-infras down --remove-orphans