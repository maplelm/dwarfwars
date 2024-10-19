#!/bin/bash

error() {
	echo $1
	exit 1
}

go build -o $PWD/bin/server/gserver ./cmd/server/ || error "Failed to build server"
cd $PWD/bin/server/ || error "Failed to switch to server bin directory"
./gserver -h
cd -
