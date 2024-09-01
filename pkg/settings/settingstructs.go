package settings

type Config struct {
	Log       LogSettings         `toml:"Log Settings"`
	Server    ServerSettings      `toml:"Game Server"`
	Databases map[string]Database `toml:"Database Servers"`
}

type LogSettings struct {
	MaxFileSize int64  `toml:"Max Size"` // in megabytes
	FileName    string `toml:"File Name"`
	Path        string `toml:"Path"`
}

type ServerSettings struct {
	Addr        string `toml:"Address"`
	Port        int    `toml:"Port"`
	IdleTimeout int    `toml:"Idle Timeout"`
}
