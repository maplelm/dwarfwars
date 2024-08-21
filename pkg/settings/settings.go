package settings

import (
	"fmt"
	"os"
	"path/filepath"
	//"reflect"
	"github.com/BurntSushi/toml"
)

var (
	dataCache map[string]interface{} = make(map[string]interface{})
)

func LoadFromTomlFile(key, path, name string) (settingData *interface{}, err error) {
	if d, ok := dataCache[key]; ok {
		return &d, fmt.Errorf("LoadFromTomlFile: Key %s is in use", key)
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

func LoadFromToml(key string, d []byte) (settingsData *interface{}, err error) {
	if d, ok := dataCache[key]; ok {
		return &d, fmt.Errorf("LoadFromToml: Key %s in use", key)
	}
	settingsData = new(interface{})
	err = toml.Unmarshal(d, settingsData)
	if err != nil {
		return nil, err
	}
	dataCache[key] = settingsData
	return
}

func LoadFromJsonFile(key, path, name string) (settingsData *interface{}, err error) {
	return
}

func LoadFromJson(key string, data []byte) (settingsData interface{}, err error) {
	return
}

func Get[T any](key string) (dat *T, err error) {
	d, exist := dataCache[key]
	if !exist {
		return nil, fmt.Errorf("settings.Get: key %s does not exist\n", err)
	}
	switch value := d.(type) {
	case T:
		return &value, nil
	default:
		bytes, err := toml.Marshal(value)
		if err != nil {
			return nil, fmt.Errorf("settings.Get: failed to Marshal data, %s", err)
		}
		var formatedData T
		err = toml.Unmarshal(bytes, &formatedData)
		if err != nil {
			return nil, fmt.Errorf("settings.Get: failed to Unmarshal data, %s", err)
		}
		return &formatedData, nil
	}

}
