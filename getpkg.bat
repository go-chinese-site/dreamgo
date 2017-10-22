@echo off

setlocal

if exist getpkg.bat goto ok
echo getpkg.bat must be run from its folder
goto end

:ok

set OLDGOPATH=%GOPATH%
set GOPATH=%~dp0

if not exist ./bin/gvt.exe (

	echo ....get the gvt tool....
	go get github.com/polaris1119/gvt
)

if not exist ./bin/gvt.exe (
	echo  get the gvt tool fail
	echo  You may obtain it with the following command:
	echo  go get github.com/polaris1119/gvt
)else (
	echo get the gvt tool success
)

cd src

gvt restore -connections 8 -precaire

cd ..

set GOPATH=%OLDGOPATH%

:end
echo finished

