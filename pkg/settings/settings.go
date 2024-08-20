package settings

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"

	"github.com/BurntSushi/toml"
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

func Get[T any](key string) (dat T, err error) {
	var settingsData interface{}

	settingsData, exists := data[key]
	if reflect.TypeOf(dat).Kind() == reflect.Interface {
		dat = settingsData.(T)
	}
	if !exists {
		return dat, fmt.Errorf("settings.Get: key %s does not exist", key)
	}
	var sdt reflect.Type = reflect.TypeOf(settingsData)
	var dt reflect.Type = reflect.TypeOf(dat)
	fmt.Println("Printing variable Types:")
	fmt.Printf("type of settingsData: %s\ntype of dat: %s\n", sdt.Name(), dt.Name())
	return
	/*
		switch typedSettings := settingsData.(type) {
		case T:
			return &typedSettings, nil
		default:
			bytes, err := toml.Marshal(typedSettings)
			if err != nil {
				return nil, fmt.Errorf("settings.Get: %s", err)
			}
			var returnVal T
			err = toml.Unmarshal(bytes, &returnVal)
			if err != nil {
				return nil, fmt.Errorf("settings.Get: %s", err)
			}
			return &returnVal, nil
		}
	*/
}
