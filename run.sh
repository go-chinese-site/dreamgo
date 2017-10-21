#!/usr/bin/env bash

#set -e

if [ ! -f run.sh ]; then
	echo 'install must be run within its container folder' 1>&2
	exit 1
fi

ps -ef|grep dreamgo |grep -v grep
if [ $? -ne 0 ];then
    source ./install.sh
    if [ ! -f bin/dreamgo ];then
     echo .........the command install fail.........
     exit 1
     else 
     echo .........the program is starting.......
     nohup ./bin/dreamgo >/dev/null 2>&1 &
     echo .........started success.........
     echo finished
     fi
else
    echo .........run kill exist process.........
    ps -ef|grep dreamgo |grep -v grep |awk '{print $2}' |xargs kill -9 >/dev/null 2>&1
     source ./install.sh
    if [ ! -f bin/dreamgo ];then
     echo .........the command install fail.........
     exit 1
     else 
     echo .........the program is starting.......
     nohup ./bin/dreamgo >/dev/null 2>&1 &
     echo .........started success.........
     echo start finished
     fi
fi
