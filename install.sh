#!/usr/bin/env bash

set -e

if [ ! -f install.sh ]; then
	echo 'install must be run within its container folder' 1>&2
	exit 1
fi

CURDIR=`pwd`
OLDGOPATH="$GOPATH"
export GOPATH="$CURDIR"
export GOBIN=

if [ ! -d log ]; then
	mkdir log
fi

gofmt -w -s src

go install dreamgo

export GOPATH="$OLDGOPATH"

echo 'finished'