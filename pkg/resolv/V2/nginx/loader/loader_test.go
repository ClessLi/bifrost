package loader

import (
	"encoding/json"
	"testing"
)

func TestLoader_LoadFromFilePath(t *testing.T) {
	configPath := "../../../../../test/nginx/conf/nginx.conf"
	loader := NewLoader()
	config, _, err := loader.LoadFromFilePath(configPath)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf(string(config.Bytes()))

}

func TestLoader_LoadFromJsonBytes(t *testing.T) {
	configPath := "../../../../../test/nginx/conf/nginx.conf"
	loader := NewLoader()
	config, _, err := loader.LoadFromFilePath(configPath)
	if err != nil {
		t.Fatal(err)
	}

	jsonData, err := json.Marshal(config)
	if err != nil {
		t.Fatal(err)
	}

	testConfig, _, err := loader.LoadFromJsonBytes(jsonData)
	if err != nil {
		t.Fatal()
	}

	t.Logf(string(testConfig.Bytes()))
}
