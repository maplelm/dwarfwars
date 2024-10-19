package types

type Options struct {
	General struct {
	} `toml:"General_Settings"`

	Network struct {
		Addr       string `toml:"Address"`
		Port       int    `toml:"Port"`
		BufferSize int    `toml:"Buffer_Size"`
	} `toml:"Network_Settings"`

	Logging struct {
	} `toml:"Logging_Settings"`
}
