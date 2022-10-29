#! bin/sh

# move to working directory
cd ../../

# docker-compose
docker-compose -f docker-compose.dev.yaml -p ktpm-server-infra down
docker-compose -f docker-compose.dev.yaml -p ktpm-server-infra up -d
