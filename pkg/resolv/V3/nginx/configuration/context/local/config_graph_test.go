package local

import (
	"testing"

	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context_type"
)

func TestModifyConfigPathInGraph(t *testing.T) {
	testMain, err := NewMain("C:\\test\\nginx.conf")
	if err != nil {
		t.Fatal(err)
	}
	childConfig := NewContext(context_type.TypeConfig, "test.conf").(*Config)

	err = testMain.AddConfig(childConfig)
	if err != nil {
		t.Fatal(err)
	}

	err = childConfig.SetValue("test1.conf")
	if err != nil {
		t.Fatal(err)
	}

	cache, err := testMain.GetConfig(childConfig.FullPath())
	if err != nil {
		t.Fatal(err)
	}
	if cache.FullPath() != "C:\\test\\test1.conf" {
		t.Errorf("got = %s, want C:\\test\\test1.conf", cache.FullPath())
	}
}
