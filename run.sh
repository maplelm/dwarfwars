#!/bin/bash

error() {
	echo $1
	exit 1
}

server () {
	go build -o $PWD/src/server/bin/gserver $PWD/src/server/ || error "Failed to build server"
	cd $PWD/src/server/bin/ || error "Failed to switch to server bin directory"
	command ./gserver -h &
	pid=$!
	trap "kill $pid" SIGINT
	wait $pid
	cd -
} 

client() {
	go build -o $PWD/client/bin/gclient ./src/client/ || error "Failed to build client"
	cd $PWD/src/client/bin/ || error "Failed to switch to client bin directory"
	command go run ./gclient &
	pid=$!
	trap "kill $pid" SIGINT
	wait $pid
	cd -
}

if [ "$1" = "server" ]; then
	server
elif [ "$1" = "client" ]; then
	client
else
	echo "Please specify server or client"
fi

exit 0
