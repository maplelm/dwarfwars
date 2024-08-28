#!/bin/bash

commands=( 
	"Read Game Server Current Logs"
	"Run Game Server"
)

Run() {
local selection=$( printf "%s\n" "${commands[@]}" | fzf --prompt="Select a command > " --height=100% --layout=reverse --border=double --exit-0 )

if [[ $selection = "Read Game Server Current Logs" ]]; then
	ReadCurrentLogFile gameserver
elif [[ $selection = "Run Game Server" ]]; then
	Runcmd gameserver
elif [[ -z $selection ]]; then
	echo "Debug Cancled"
fi
}

ReadCurrentLogFile()
{
	echo "Reading: $PWD/cmd/$1/bin/logs/log.log"
	less "$PWD/cmd/$1/bin/logs/log.log"
}

Runcmd() 
{
	cd "$PWD/cmd/$1/bin" || exit 1
	SETTINGSPATH='./' go run ../
	cd "../../.."
}


Run
