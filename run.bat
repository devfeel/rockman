@echo off
cd /d %~dp0/webui/webapp
call npm run build
cd /d %~dp0
call go run main.go