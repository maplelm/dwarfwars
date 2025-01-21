package types

type Options struct {
	Logging struct {
		MaxSize     int    `toml:"Size"`
		MaxAge      int    `toml:"Max_Age"`
		Backups     int    `toml:"Backups"`
		Compression bool   `toml:"Compression"`
		Name        string `toml:"Name"`
		Path        string `toml:"Path"`
		Prefix      string `toml:"Prefix"`
		Flags       struct {
			UTC          bool `toml:"UTC"`
			Date         bool `toml:"Date"`
			Time         bool `toml:"Time"`
			Longfile     bool `toml:"Longfile"`
			Msgprefix    bool `toml:"Msgprefix"`
			Shortfile    bool `toml:"Shortfile"`
			Microseconds bool `toml:"Microseconds"`
		} `toml:"Flags"`
	} `toml:"Logging"`

	Game struct {
		Addr     string `toml:"Address"`
		Port     int    `toml:"Port"`
		Timeouts struct {
			Idle  int `toml:"Idle"`
			Read  int `toml:"Read"`
			Write int `toml:"Write"`
		} `toml:"Timeouts"`
	} `toml:"Game Server"`

	Db struct {
		BaseDir  string  `toml:"Base_Dir"`
		Password *string `toml:"Password"`
		Querries string  `toml:"Querries_Dir"`
	} `toml:"DB"`
}
