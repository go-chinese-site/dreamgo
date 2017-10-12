@echo off

setlocal

if exist run.bat goto exec
 echo run.bat must be run from its folder

:exec
tasklist /nh|find /i "dreamgo"

if ERRORLEVEL 1 (
    goto reload
 ) else (
     goto kill
 )



:reload
call install.bat
if exist bin\dreamgo.exe (
    echo .........rebuild success.........
    goto run
    )else (
    echo .........the command install fail.........
    goto end
    )

:run
echo .........the program is starting.......
start /min bin\dreamgo.exe
echo .........started success.........
goto end

:kill
echo .........run kill exist process.........
 taskkill /f /im "dreamgo.exe" /t >null 2>&1
 goto reload

:end
echo run finished

 
