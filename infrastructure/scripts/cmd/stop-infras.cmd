title 'Stop ktpm-server-infrastructure'

:: move to infrastructure folder
cd ../

:: compose up
docker-compose -f docker-compose.yaml -p ktpm-server-infra down --remove-orphans

:: post-action
@echo off
@echo ktmp-server-infra compose up, auto close in 5s

timeout 5