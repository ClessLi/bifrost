package loader

import "testing"

func TestLoader_LoadFromFilePath(t *testing.T) {
	configPath := "../../../../../test/nginx/conf/nginx.conf"
	loader, err := NewLoader(configPath)
	if err != nil {
		t.Fatal(err)
	}
	config, err := loader.LoadFromFilePath(configPath)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf(string(config.Bytes()))
}
