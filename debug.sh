#!/bin/bash

commands=( "less ./cmd/gameserver/bin/logs/log.log"
	"cd ./cmd/gameserver/bin/ && ( SETTINGSPATH='./' go run ../ || echo 'Failed' ) && cd ../../.."
)

selection=$( printf "%s\n" "${commands[@]}" | fzf --prompt="Select a command > " --height=100% --layout=reverse --border=double --exit-0 )

echo "Running: $selection"
$selection

