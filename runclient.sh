#!/bin/bash

error() {
	echo $1
	exit 1
}

go build -o $PWD/bin/client/gclient ./cmd/client/ || error "Failed to build client"
cd $PWD/bin/client/ || error "Failed to switch to client bin directory"
./gclient "$@"
cd -

