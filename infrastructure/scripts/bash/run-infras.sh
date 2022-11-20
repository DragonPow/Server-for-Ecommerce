#! bin/sh

# This file is used to start up all infrastructure config for services by docker-compose

# move to infrastructure directory
cd ../../

# docker-compose
docker-compose -f docker-compose.infras.dev.yaml -p ktpm-server-infra down
docker-compose -f docker-compose.infras.dev.yaml -p ktpm-server-infra up -d
