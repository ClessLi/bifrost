package V2_test

import (
	"strings"
	"testing"

	"github.com/ClessLi/bifrost/pkg/resolv/V2/nginx/configuration"
	"github.com/ClessLi/bifrost/pkg/resolv/V2/nginx/configuration/parser"
	"github.com/ClessLi/bifrost/pkg/resolv/V2/nginx/parser_type"
)

func TestQueryNonCascaded(t *testing.T) {
	configPath := "../../../nginx/conf/nginx.conf"
	config, err := configuration.NewConfigurationFromPath(configPath)
	if err != nil {
		t.Fatal(err)
	}

	testLocation, err := config.Query(`location:sep: :reg: /test1-location`)
	if err != nil {
		t.Fatal(err)
	}

	proxyPassKW, err := parser.NewKeyWords(parser_type.TypeKey, true, `^proxy_pass\s+`)
	if err != nil {
		t.Fatal(err)
	}
	proxyPassKW.SetCascaded(false)
	ctx, idx := testLocation.Self().(*parser.Location).Query(proxyPassKW)
	childP, err := ctx.GetChild(idx)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.EqualFold(childP.GetValue(), "proxy_pass http://right_proxy") {
		t.Fatalf("want: 'proxy_pass http://right_proxy', got: '%s'", childP.GetValue())
	}
}
