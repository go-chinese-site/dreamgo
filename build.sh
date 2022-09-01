#!/usr/bin/env bash

set -e

if [ ! -f build.sh ]; then
	echo 'build must be run within its container folder' 1>&2
	exit 1
fi

CURDIR=`pwd`

export GO111MODULE=on
export GOPROXY=https://goproxy.cn

if [ ! -d log ]; then
	mkdir log
fi

gofmt -w -s .

go build -o ./bin/dreamgo


echo 'finished'