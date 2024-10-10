package main

type Options struct {
	Logging struct {
		MaxSize int    `toml:"Size"`
		Name    string `toml:"Name"`
		Path    string `toml:"Path"`
		Prefix  string `toml:"Prefix"`
		Flags   struct {
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
			Idle  int `toml: "Idle"`
			Read  int `toml:"Read"`
			Write int `toml:"Write"`
		} `toml:"Timeouts"`
	} `toml:"Game Server"`

	Db struct {
		Addr          string `toml:"Server Address"`
		Port          int    `toml:"Server Port"`
		Username      string `toml:"Server Username"`
		Password      string `toml:"Server Password"`
		ValidationDir string `toml:"Validation Direcroty"`
	} `toml:"DB"`
}
