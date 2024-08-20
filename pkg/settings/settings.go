package settings

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"os"
	"path/filepath"
	//"reflect"
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
	fmt.Printf("settings.Get: Generic type = %T\n", [0]T{})
	if dat, ok := data[key]; !ok {
		return nil, fmt.Errorf("settings.Get: settings with key %s does not exist", key)
	} else {
		switch val := dat.(type) {
		case T:
			return &val, nil
		default:
			var dt T
			dt, ok = val.(T)
			if ok {
				return &dt, nil
			} else {
				return nil, fmt.Errorf("value is not of type %T", dt)
			}
			// if zero := []T{}; reflect.TypeOf(dat).AssignableTo(reflect.TypeOf(dt)) {
			// 	r := dat.(T)
			// 	return &r, nil
			// } else {
			// 	return nil, fmt.Errorf("type %T is not assignable to type %T", reflect.TypeOf(val).Name(), reflect.TypeOf(zero).Name())
			// }
		}
	}
}
