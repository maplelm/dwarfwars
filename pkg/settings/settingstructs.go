package settings

import (
// "fmt"
)

type MainDB struct{}

const (
	DbNameMain = iota
)

type General struct {
	Logging   Logs `toml:"Log Settings"`
	Databases []Database
}

type Settings struct {
	Log        Logs           `toml:"Log Settings"`
	Server     ServerSettings `toml:"Game Server"`
	SQLServers SQLServers     `toml:"Database Servers"`
}

type Config struct {
	Log       Logs                `toml:"Log Settings"`
	Server    ServerSettings      `toml:"Game Server"`
	Databases map[string]Database `toml:"Database Servers"`
}

type SQLServers struct {
}

type Logs struct {
	MaxFileSize int64  `toml:"Max Size"` // in megabytes
	FileName    string `toml:"File Name"`
	Path        string `toml:"Path"`
}

type ServerSettings struct {
	Addr        string `toml:"Address"`
	Port        int    `toml:"Port"`
	IdleTimeout int    `toml:"Idle Timeout"`
}
