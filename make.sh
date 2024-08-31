#!/bin/sh


# NOTE: Need to make script work regardless of cwd.

server_name="gserver"
client_name="gclient"

build() {
	local neededdir="./bin/$2/"
	if [[ ! -d "./bin/$2/" ]]; then
		echo "Needed Dir Does not exists, Creating Dir: $neededdir"
		mkdir -p "$neededdir" || exit 1
	fi
	go build -o "$neededdir/$1" ./cmd/gameserver/ || exit 1
	exit 0
}

run() {
	local neededdir="./bin/$2"
	if [[ ! -d "$neededdir" ]];then
		mkdir -p "$neededdir"
	fi
	if [[ ! -f "$neededdir/$1" ]]; then
		buildserver "$1" "$2" || exit 1
	fi
	local pwd=$(pwd)
	cd "./bin/$2" || exit 1
	/bin/sh -c "./$1"
	cd "$pwd" || exit 1
}


if [[ $1 == "bs" ]]; then
	build "$server_name" "/server/"
elif [[ $1 == "rs" ]]; then
	run "$server_name" "/server/"
elif [[ $1 == "bc" ]]; then
	build "$client_name" "/client/"
elif [[ $1 == "rc" ]]; then
	run "$client_name" "/client"
else
	echo "Unknown Command: $1"
fi

