package types

import "time"

type Options struct {
	General struct {
		Fullscreen   bool    `toml:"FullScreen"`
		ScreenWidth  float32 `toml:"Screen_Width"`
		ScreenHeight float32 `toml:"Screen_Height"`
		TargetFPS    int32   `toml:"FPS"`
		Font         string  `toml:"Default_Font"`
		FontSize     int32   `toml:"Default_Font_Size`
		FontRes      int32   `toml:"Default_Font_Resolution"`
	} `toml:"General_Settings"`

	Network struct {
		Addr        string        `toml:"Address"`
		Port        int           `toml:"Port"`
		ConnTimeout time.Duration `toml:"Connection_Timeout"`
		BufferSize  int           `toml:"Buffer_Size"`
	} `toml:"Network_Settings"`

	Logging struct {
	} `toml:"Logging_Settings"`
}
