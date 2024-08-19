package main

import "time"

type Config struct {
	Log    LogSettings    `toml:"Log Settings"`
	Server ServerSettings `toml:"Server Settings"`
}

type LogSettings struct {
	MaxFileSize int64  `toml:"LogMaxSize"` // in megabytes
	PollRate    int    `toml:"LogPollRate"`
	FileName    string `toml:"LogFileName"`
	Path        string `toml:"LogPath"`
}

func (ls LogSettings) AdjustedPollRate() time.Duration {
	return time.Duration(ls.PollRate) * time.Millisecond
}

type ServerSettings struct {
	Addr        string `toml:"Address"`
	Port        int    `toml:"Port"`
	IdleTimeout int    `toml:"Idle_Timeout"`
}
