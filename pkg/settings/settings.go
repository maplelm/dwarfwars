package settings

import (
	"encoding/json"
	"fmt"
	"github.com/BurntSushi/toml"
	/*
		"os"
		"path/filepath"
		"strings"
	*/)

func LoadFromFile() {}

func Load[T any](t string, rawdata []byte) (data *T, err error) {
	switch t {
	case "json", "toml":
		data = new(T)
		switch t {
		case "json":
			err = json.Unmarshal(rawdata, data)
		case "toml":
			err = toml.Unmarshal(rawdata, data)
		}
		return
	default:
		return nil, fmt.Errorf("unsupported file type: %s", t)
	}
}
