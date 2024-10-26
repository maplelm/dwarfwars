package types

type Options struct {
	General struct {
		Fullscreen   bool    `toml:"FullScreen"`
		ScreenWidth  float32 `toml:"Screen_Width"`
		ScreenHeight float32 `toml:"Screen_Height"`
	} `toml:"General_Settings"`

	Network struct {
		Addr       string `toml:"Address"`
		Port       int    `toml:"Port"`
		BufferSize int    `toml:"Buffer_Size"`
	} `toml:"Network_Settings"`

	Logging struct {
	} `toml:"Logging_Settings"`
}
