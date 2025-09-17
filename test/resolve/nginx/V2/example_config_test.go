package V2

import (
	"fmt"
	"testing"

	"github.com/ClessLi/bifrost/pkg/resolv/V2/nginx/configuration"
	"github.com/ClessLi/bifrost/pkg/resolv/V2/nginx/configuration/parser"
	"github.com/ClessLi/bifrost/pkg/resolv/V2/nginx/parser_type"
)

func TestExampleConfig(t *testing.T) {
	// config可以来自bifrost客户端请求获取的json数据进行解析
	c, err := configuration.NewConfigurationFromPath("test_nginx.conf")
	if err != nil {
		t.Fatal(err)
	}

	// http context可以来自config检索结果
	httpCTX := parser.NewContext("", parser_type.TypeHttp, c.Self().GetIndention()) // 使用config对象的indention作为原始缩进

	// 插入http context至config
	err = c.Self().(parser.Context).Insert(httpCTX, 0)
	if err != nil {
		t.Fatal(err)
	}

	// 创建server context
	srvCTX := parser.NewContext("", parser_type.TypeServer, httpCTX.GetIndention().NextIndention()) // 使用http context的下级缩进
	// 插入server context至http context
	err = httpCTX.Insert(srvCTX, 0)
	if err != nil {
		t.Fatal(err)
	}

	// 插入server context的server_name和listen字段
	err = srvCTX.Insert(parser.NewKey("server_name", "example.com", srvCTX.GetIndention().NextIndention()), 0) // 使用server context的下级缩进，插入为server context首个子parser
	if err != nil {
		t.Fatal(err)
	}
	err = srvCTX.Insert(parser.NewComment("Server Context首个字段", true, srvCTX.GetIndention().NextIndention()), srvCTX.Len())
	if err != nil {
		t.Fatal(err)
	}
	err = srvCTX.Insert(parser.NewKey("listen", "80", srvCTX.GetIndention().NextIndention()), srvCTX.Len()) // 使用server context的下级缩进，插入到server context末尾
	if err != nil {
		t.Fatal(err)
	}
	err = srvCTX.Insert(parser.NewComment("Server Context第二个字段", true, srvCTX.GetIndention().NextIndention()), srvCTX.Len())
	if err != nil {
		t.Fatal(err)
	}
	err = srvCTX.Insert(parser.NewComment("Location Context", false, srvCTX.GetIndention().NextIndention()), srvCTX.Len())
	if err != nil {
		t.Fatal(err)
	}

	// 创建location context
	locationCTX := parser.NewContext("/", parser_type.TypeLocation, srvCTX.GetIndention().NextIndention()) // 使用server context的下级缩进
	// 插入到server context末尾
	err = srvCTX.Insert(locationCTX, srvCTX.Len())
	if err != nil {
		t.Fatal(err)
	}

	err = locationCTX.Insert(parser.NewComment("Location Context自定义字段", false, locationCTX.GetIndention().NextIndention()), 0)
	if err != nil {
		t.Fatal(err)
	}
	// 插入location context的指定字段
	err = locationCTX.Insert(parser.NewKey("root", "index.html", locationCTX.GetIndention().NextIndention()), locationCTX.Len()) // 使用location context的下级缩进，插入到location context末尾
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(string(c.View()))

	// 以Configuration接口插入配置
	// 创建server2 context
	srv2CTX := parser.NewContext("", parser_type.TypeServer, c.Self().GetIndention().NextIndention()) // 使用config对象的第二级缩进
	// 插入server2 context的server_name和listen字段
	err = srv2CTX.Insert(parser.NewKey("server_name", "baidu.com", srv2CTX.GetIndention().NextIndention()), 0) // 使用server2 context的下级缩进，插入为server2 context首个子parser
	if err != nil {
		t.Fatal(err)
	}
	err = srv2CTX.Insert(parser.NewKey("listen", "80", srv2CTX.GetIndention().NextIndention()), srv2CTX.Len()) // 使用server2 context的下级缩进，插入到server2 context末尾
	if err != nil {
		t.Fatal(err)
	}

	// 创建location2 context
	location2CTX := parser.NewContext("/", parser_type.TypeLocation, srv2CTX.GetIndention().NextIndention()) // 使用server context的下级缩进
	// 插入到server2 context末尾
	err = srv2CTX.Insert(location2CTX, srv2CTX.Len())
	if err != nil {
		t.Fatal(err)
	}

	// 插入location context的指定字段
	err = location2CTX.Insert(parser.NewKey("aaa", "bbbb", location2CTX.GetIndention().NextIndention()), location2CTX.Len()) // 使用location context的下级缩进，插入到location context末尾
	if err != nil {
		t.Fatal(err)
	}

	// Configuration.InsertByKeyword，目前只支持插入到检索到的Queryer对象前。
	// Configuration接口和Queryer接口是为了方便同事的快速检索和快速操作各配置对象而设置的，功能会存在局限性，该接口更适合检索和操作config整体数据，及简单的增、删、改
	// 建议使用Configuration.Self方法，将Configuration转换为Parser接口对象后，再按Parser/Context接口对象来进行更精准的操作
	err = c.InsertByKeyword(srv2CTX, "server")
	if err != nil {
		t.Fatal(err)
	}

	t.Log("通过Configuration接口插入后")
	fmt.Println("\n" + string(c.View()))

	srvnameKW, err := parser.NewKeyWords(parser_type.TypeKey, false, "server_name baidu.com")
	if err != nil {
		t.Fatal(err)
	}
	lisKW, err := parser.NewKeyWords(parser_type.TypeKey, false, "listen 80")
	if err != nil {
		t.Fatal(err)
	}
	s, _ := c.Self().(parser.Context).Query(srvnameKW)
	if s != nil {
		ctx, i := s.Query(lisKW)
		err = ctx.Modify(parser.NewKey("listen", "443", ctx.GetIndention().NextIndention()), i)
		if err != nil {
			t.Fatal(err)
		}
		err = ctx.Insert(parser.NewKey("ssl", "on", ctx.GetIndention().NextIndention()), i+1)
		if err != nil {
			t.Fatal(err)
		}
		err = ctx.Insert(parser.NewComment("修改http为https", false, ctx.GetIndention().NextIndention()), i)
		if err != nil {
			t.Fatal(err)
		}
	}

	fmt.Println("\n" + string(c.View()))
	t.Log("complete!")
}
