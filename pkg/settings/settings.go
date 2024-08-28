package settings

import (
	"encoding/json"
	"fmt"
	"github.com/BurntSushi/toml"
	"os"
	"path/filepath"
)

var (
	dataCache map[string]*interface{} = make(map[string]*interface{}) // Map of all previsou settings
)

func LoadFromTomlFile(key, path, name string) (settingData *interface{}, err error) {
	if d, ok := dataCache[key]; ok {
		return d, fmt.Errorf("LoadFromTomlFile: Key %s is in use", key)
	}
	fileData, err := os.ReadFile(filepath.Join(path, name))
	if err != nil {
		return nil, fmt.Errorf("LoadFromToml: %s", err)
	}
	settingData = new(interface{})
	err = toml.Unmarshal(fileData, settingData)
	if err != nil {
		return nil, fmt.Errorf("LoadFromToml: %s", err)
	}
	dataCache[key] = settingData
	return
}

func LoadFromToml(key string, data []byte) (settingsData *interface{}, err error) {
	if d, ok := dataCache[key]; ok {
		return d, fmt.Errorf("LoadFromToml: Key %s in use", key)
	}
	settingsData = new(interface{})
	err = toml.Unmarshal(data, settingsData)
	if err != nil {
		return nil, err
	}
	dataCache[key] = settingsData
	return
}

func LoadFromJsonFile(key, path, name string) (settingsData *interface{}, err error) {
	if d, ok := dataCache[key]; ok {
		return d, fmt.Errorf("LoadFromJsonFile: key %s already exists", key)
	}
	bytes, err := os.ReadFile(filepath.Join(path, name))
	if err != nil {
		return nil, fmt.Errorf("LoadFromJsonFile: %s", err)
	}
	settingsData = new(interface{})
	err = json.Unmarshal(bytes, settingsData)
	if err != nil {
		return nil, fmt.Errorf("LoadFromJsonFile: %s", err)
	}
	dataCache[key] = settingsData
	return
}

func LoadFromJson(key string, data []byte) (settingsData *interface{}, err error) {
	if d, ok := dataCache[key]; ok {
		return d, fmt.Errorf("LoadFromJson: %s", err)
	}
	settingsData = new(interface{})
	err = json.Unmarshal(data, settingsData)
	if err != nil {
		return nil, fmt.Errorf("LoadFromJson: %s", err)
	}
	return
}

func Get[T any](key string) (dat *T, err error) {
	if d, ok := dataCache[key]; !ok {
		return nil, fmt.Errorf("settings.Get: settings with key %s does not exist", key)
	} else {
		switch val := (*d).(type) {
		case T:
			return &val, nil
		default:

			bytes, err := toml.Marshal(*d)
			if err != nil {
				return nil, err
			}
			dat = new(T)
			err = toml.Unmarshal(bytes, dat)
			if err != nil {
				return nil, err
			}
			return dat, err
		}
	}
}
