package local

import (
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context_type"
	"testing"
)

func TestModifyConfigPathInGraph(t *testing.T) {
	testMainConfig, err := newMainConfig("C:\\test\\nginx.conf")
	if err != nil {
		t.Fatal(err)
	}
	childConfig := NewContext(context_type.TypeConfig, "test.conf").(*Config)

	_, err = testMainConfig.IncludeConfig(childConfig)
	if err != nil {
		t.Fatal(err)
	}

	err = childConfig.SetValue("test1.conf")
	if err != nil {
		t.Fatal(err)
	}

	cache, err := testMainConfig.Graph.GetConfig(childConfig.FullPath())
	if err != nil {
		t.Fatal(err)
	}
	if cache.FullPath() != "C:\\test\\test1.conf" {
		t.Errorf("got = %s, want C:\\test\\test1.conf", cache.FullPath())
	}
}
