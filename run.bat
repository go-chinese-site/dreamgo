rem @echo off

setlocal

if exist run.bat goto exec
 echo run.bat must be run from its folder

:exec
tasklist /nh|find /i "dreamgo.exe"

if ERRORLEVEL 1 goto reload else(
    goto stop
    goto reload
    )

:reload
call install.bat
if exist bin\dreamgo.exe go run else echo install fail...

:run
echo running success...
bin\dreamgo.exe

:stop
 taskkill /im "dreamgo.exe"
