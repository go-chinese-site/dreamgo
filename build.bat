@echo off

setlocal

if exist build.bat goto ok
echo build.bat must be run from its folder
goto end

:ok

set GO111MODULE=on
set GOPROXY=https://goproxy.cn,direct

if not exist log mkdir log

gofmt -w -s .

go build -o ./bin/dreamgo.exe

:end
echo finished