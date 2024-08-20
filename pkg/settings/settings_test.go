package settings

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/BurntSushi/toml"
)

func TestGetSettings(t *testing.T) {
	testData := []byte(`[logsettings]
	file="stuff"
	max=5
	age=10
	path="this/path/right/here"
	[serversettings]
	adder="0.0.0.0"
	Port="4000"
	timeout=10
	`)

	_, err := LoadFromToml("testkey", testData)
	if err != nil {
		t.Fatalf("Failed to parse tomal data, %s\n", err)
	}

	data, err := Get[any]("testkey")
	if err != nil {
		t.Fatalf("Failed to get settings data, %s\n", err)
	}

	display, _ := toml.Marshal(*data)
	t.Logf("Successfully got data, from memory, %s", string(display))

	cwd, err := os.Getwd()
	t.Logf("Current Work Directory: %s\n", cwd)
	d, err := LoadFromTomlFile("testfilekey", cwd, "setting.toml")
	if err != nil && d == nil {
		t.Fatalf("Failed to get Toml data from file, %s", err)
	}

	filedata, err := os.ReadFile(filepath.Join(cwd, "setting.toml"))
	if err != nil {
		t.Fatalf("Failed to read in file data for the base check, %s", err)
	}

	var c Config
	err = toml.Unmarshal(filedata, &c)
	if err != nil {
		t.Fatalf("Failed to Unmarshal base case, %s\n", err)
	}
	b, err := toml.Marshal(c)
	if err != nil {
		t.Fatalf("Failed to Marshal base case, %s", err)
	}
	t.Logf("Base Case setting.toml:\n%s", string(b))

	newdata, err := Get[Config]("testfilekey")
	if err != nil {
		t.Fatalf("Failed to get settings data, %s\n", err)
	}
	display, _ = toml.Marshal(*newdata)
	t.Logf("Successfully got data from file\n %s\n", string(display))
}
