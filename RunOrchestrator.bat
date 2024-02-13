@echo off
setlocal
set CGO_ENABLED=1
cd cmd\orchestrator
start /b go run main.go
endlocal