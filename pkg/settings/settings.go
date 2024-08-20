package settings

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"os"
	"path/filepath"
	"reflect"
)

var (
	data map[string]interface{} = make(map[string]interface{})
)

func LoadFromTomlFile(key, path, name string) (settingData *interface{}, err error) {
	if d, ok := data[key]; ok {
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
	data[key] = settingData
	return
}

func LoadFromToml(key string, d []byte) (settingsData *interface{}, err error) {
	if d, ok := data[key]; ok {
		return &d, fmt.Errorf("LoadFromToml: Key %s in use", key)
	}
	settingsData = new(interface{})
	err = toml.Unmarshal(d, settingsData)
	if err != nil {
		return nil, err
	}
	data[key] = settingsData
	return
}

func LoadFromJsonFile(key, path, name string) (settingsData *interface{}, err error) {
	return
}

func LoadFromJson(key string, data []byte) (settingsData interface{}, err error) {
	return
}

func Get[T any](key string) (dat *T, err error) {
	if dat, ok := data[key]; !ok {
		return nil, fmt.Errorf("settings.Get: settings with key %s does not exist", key)
	} else {
		switch v := dat.(type) {
		case T:
			return &v, nil
		default:
			if zero := []T{}; reflect.TypeOf(v).AssignableTo(reflect.TypeOf(zero)) {
				r := v.(T)
				return &r, nil
			} else {
				return nil, fmt.Errorf("type %T is not assignable to type %T", reflect.TypeOf(v).Name(), reflect.TypeOf(zero).Name())
			}
		}
	}
}
