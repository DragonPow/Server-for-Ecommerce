title 'Run ktpm-server-infrastructure'

:: move to infrastructure directory
cd ../

:: compose up
docker-compose -f docker-compose.yaml -p ktpm-server-infra down
docker-compose -f docker-compose.yaml -p ktpm-server-infra up -d

:: post-action
@echo off
@echo ktmp-server-infra compose up, auto close in 5s

timeout 5