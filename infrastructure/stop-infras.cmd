title 'Stop ktpm-server-infrastructure'

:: move to working directory
cd ../

:: compose up
docker-compose -f docker-compose.dev.yaml -p ktpm-server-infra down

:: post-action
@echo off
@echo ktmp-server-infra compose up, auto close in 5s

timeout 5