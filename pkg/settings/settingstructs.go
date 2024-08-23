package settings

import (
	"fmt"
)

type Config struct {
	Log        LogSettings    `toml:"Log Settings"`
	Server     ServerSettings `toml:"Game Server"`
	SQLServers SQLServers     `toml:"Database Servers"`
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

type SQLServerSettings struct {
	Addr      string                 `toml:"Address"`
	Driver    string                 `toml:"Driver"`
	Trusted   bool                   `toml:"Trusted"`
	Port      int                    `toml:"Port"`
	Username  string                 `toml:"Username"`
	Password  string                 `toml:"Password"`
	DefaultDB string                 `toml:"DB"`
	Views     map[string][]SQLColumn `toml:"Views"`
	Tables    map[string][]SQLColumn `toml:"Tables"`
}

func (sss SQLServerSettings) ConnString() (str *string, err error) {
	str = new(string)
	switch sss.Driver {
	case "mysql":
		*str = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", sss.Username, sss.Password, sss.Addr, sss.Port, sss.DefaultDB)
	case "mssql":
		if !sss.Trusted {
			*str = fmt.Sprintf("Server=%s:%s; Database=%s; User Id=%s; Password=%s;", sss.Addr, sss.Port, sss.DefaultDB, sss.Username, sss.Password)
		} else {
			*str = fmt.Sprintf("Server=%s:%s; Database=%s; Trusted_Connection=True", sss.Addr, sss.Port, sss.DefaultDB)
		}
	default:
		str = nil
		err = fmt.Errorf("Unsupported Driver: %s", sss.Driver)
	}
	return
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

type SQLServers map[string]SQLServerSettings

func (ss SQLServers) ToList() []SQLServerSettings {
	var l []SQLServerSettings = make([]SQLServerSettings, len(ss))
	var index int = 0
	for _, v := range ss {
		l[index] = v
		index++
	}
	return l
}
