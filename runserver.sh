#!/bin/bash

error() {
	echo $1
	exit 1
}

#go build -o $PWD/bin/server/gserver ./cmd/server/ || error "Failed to build server"
cd $PWD/bin/server/ || error "Failed to switch to server bin directory"
go run ../../cmd/server/ -h
#./gserver -h &
#PID=$!
#sleep 1
#sudo procdump -m 100 -o $PID
#wait $PID
cd -
