package V2

import (
	"testing"

	"github.com/ClessLi/bifrost/pkg/resolv/V2/nginx/dumper"
	"github.com/ClessLi/bifrost/pkg/resolv/V2/nginx/loader"
)

func TestDump(t *testing.T) {
	configPath := "../../../nginx/conf/nginx.conf"
	l := loader.NewLoader()
	config, _, err := l.LoadFromFilePath(configPath)
	if err != nil {
		t.Fatal(err)
	}

	d := dumper.NewDumper(config.GetValue())
	err = config.Dump(d)
	if err != nil {
		t.Fatal(err)
	}

	for s, bytes := range d.ReadAll() {
		t.Logf("%s:\n%s", s, bytes)
	}
}
