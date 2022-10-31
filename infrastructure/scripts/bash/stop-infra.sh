#! bin/sh

# move to infrastructure directory
cd ../../

docker-compose -f docker-compose.infras.dev.yaml -p ktpm-server-infra down --remove-orphans