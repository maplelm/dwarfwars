package settings

import "time"

type Config struct {
	Log        LogSettings                  `toml:"Log Settings"`
	Server     ServerSettings               `toml:"Game Server"`
	SQLServers map[string]SQLServerSettings `toml:"Database Servers"`
}

type LogSettings struct {
	MaxFileSize int64  `toml:"Max_Size"` // in megabytes
	PollRate    int    `toml:"Poll_Rate"`
	FileName    string `toml:"Log File Name"`
	Path        string `toml:"Path"`
}

func (ls LogSettings) AdjustedPollRate() time.Duration {
	return time.Duration(ls.PollRate) * time.Millisecond
}

type ServerSettings struct {
	Addr        string `toml:"Address"`
	Port        int    `toml:"Port"`
	IdleTimeout int    `toml:"Idle_Timeout"`
}

type SQLServerSettings struct {
	Addr      string
	Port      int
	Username  string
	Password  string
	DefaultDB string
	Views     map[string][]SQLColumn
	Tables    map[string][]SQLColumn
}

type SQLColumn struct {
	Name          string  `toml:"Name"`
	ColType       string  `toml:"Type"`
	AutoIncrement bool    `toml:"Auto Increment"`
	PrimaryKey    bool    `toml:"Primary Key"`
	ForiegnKey    bool    `toml:"Foriegn Key"`
	Nullable      bool    `toml:"Null Allowed"`
	DefaultValue  *string `toml:"Default"`
}
