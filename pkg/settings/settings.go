package settings

import (
	"encoding/json"
	"fmt"
	"github.com/BurntSushi/toml"
	"os"
	"path/filepath"
	"strings"
)

type Format byte

const (
	FormatToml Format = iota
	FormatJson
	FormatUnknown
)

type Data struct {
	data   []byte
	format Format
}

var (
	data map[string]Data = make(map[string]Data) // Map of all previsou settings
)

func LoadFromFile[T any](k, p, n string) (T, error) {
	var (
		sd  T
		err error
		fd  []byte
		f   Format
	)
	if fd, err = os.ReadFile(filepath.Join(p, n)); err != nil {
		return sd, err
	}
	t := strings.Split(n, ".")[len(strings.Split(n, "."))-1]
	switch t {
	case "toml":
		f = FormatToml
	case "json":
		f = FormatJson
	default:
		f = FormatUnknown
	}
	return Load[T](k, Data{
		data:   fd,
		format: f,
	})
}

func Load[T any](key string, sd Data) (T, error) {
	var (
		rd  T
		err error
	)
	if _, ok := data[key]; ok {
		return rd, fmt.Errorf("Key %s exists", key)
	}
	switch sd.format {
	case FormatToml:
		if err = toml.Unmarshal(sd.data, &rd); err != nil {
			return rd, err
		}
	case FormatJson:
		if err = json.Unmarshal(sd.data, &rd); err != nil {
			return rd, err
		}
	default:
		return rd, fmt.Errorf("Settings.Load: Unsupported format, %d", sd.format)
	}
	data[key] = sd
	return rd, nil
}

func Get[T any](k string) (T, error) {
	var (
		d   Data
		rd  T
		err error
		ok  bool
	)
	if d, ok = data[k]; ok {
		switch d.format {
		case FormatToml:
			if err = toml.Unmarshal(d.data, &rd); err != nil {
				return rd, err
			}
		case FormatJson:
			if err = json.Unmarshal(d.data, &rd); err != nil {
				return rd, err
			}
		default:
			return rd, fmt.Errorf("Unsupported Format: %d", d.format)
		}
	}
	return rd, fmt.Errorf("Key %s does not exist", k)
}
