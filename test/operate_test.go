package test

import (
	"fmt"
	"testing"

	"github.com/yongPhone/bifrost/pkg/resolv/nginx"
)

func TestOperate(t *testing.T) {
	conf := nginx.NewConf(nil, "test.conf")
	// comment 和 key Add方法测试
	commErr := conf.Add(nginx.TypeComment, "# Test key")
	if commErr != nil {
		t.Log(commErr)
	}
	kErr := conf.Add(nginx.TypeKey, "user:nobody")
	if kErr != nil {
		t.Log(kErr)
	}

	//httpErr := conf.Add(nginx.TypeHttp, "default_type:application/octet-stream")
	//if httpErr != nil {
	//	t.Log(httpErr)
	//}
	include := &nginx.Include{
		BasicContext: nginx.BasicContext{
			Name:     nginx.TypeInclude,
			Value:    "./subtest.conf",
			Children: []nginx.Parser{nginx.NewConf(nil, "subtest.conf")},
		},
		Key:     nil,
		Comment: nil,
		ConfPWD: "",
	}
	// include Add 方法测试
	httpErr := include.Add(nginx.TypeHttp, "default_type:application/octet-stream")
	if httpErr != nil {
		t.Log(httpErr)
	}

	// 配置加入 include
	conf.AddByParser(include)

	// 测试配置全局查询 http 上下文
	http := nginx.GetHTTP(conf)
	// http Add, Insert 方法测试
	sErr := http.Add(nginx.TypeServer, "listen:80", "server_name:example.com")
	if sErr != nil {
		t.Log(sErr)
	}

	s2Err := http.Add(nginx.TypeServer, "listen:443", "server_name:example.com")
	if s2Err != nil {
		t.Log(s2Err)
	}

	s3Err := http.Insert(
		http.Query(nginx.TypeServer, false, "listen:^80$", "server_name:example.com"),
		nginx.TypeServer,
		"listen:80",
		"server_name:test.example.com",
	)
	if s3Err != nil {
		t.Log(s3Err)
	}

	//sdErr := http.Remove(nginx.TypeServer, "listen:^80$", "server_name:^example.com$")
	//if sdErr != nil {
	//	t.Log(sdErr)
	//}
	// 测试配置子集级别查询 include 内 http 对象
	h := conf.Query(nginx.TypeHttp, false)
	if h != nil {
		// http Remove 方法测试
		sdErr := h.(*nginx.Http).Remove(nginx.TypeServer, "listen:^80$", "server_name:^example.com$")
		if sdErr != nil {
			t.Log(sdErr)
		}

		// http Modify 方法测试， 配置查询 include 内 server 对象
		smErr := h.(*nginx.Http).Modify(
			conf.Query(nginx.TypeServer, true, "listen:^443$"),
			nginx.TypeServer,
			"listen:5443",
			"server_name:test2.example.com",
		)
		if smErr != nil {
			t.Log(smErr)
		}
	}

	// 配置 Insert 方法测试，插入 comment, upstream 到 include 内 http 对象前
	cicommErr := conf.Insert(h, nginx.TypeComment, "#test for insert upstream to include")
	if cicommErr != nil {
		t.Log(cicommErr)
	}
	ciupErr := conf.Insert(
		h,
		nginx.TypeUpstream,
		"test",
		"server:1.1.1.1:443    weight=10",
		"server:1.1.1.2:443  weight=10",
	)
	if ciupErr != nil {
		t.Log(ciupErr)
	}

	// 配置 Modify 方法测试，修改 include 内 http 对象为 events 对象
	cmh2eErr := conf.Modify(h, nginx.TypeEvents, "worker_connections:1024")
	if cmh2eErr != nil {
		t.Log(cmh2eErr)
	}

	// 配置 Remove 方法测试， 删除 include 内 events 对象
	rmeErr := conf.Remove(nginx.TypeEvents)
	if rmeErr != nil {
		t.Log(rmeErr)
	}

	//caches := nginx.NewCaches()
	//fmt.Println(conf.string(&caches))
	fmt.Println(conf.String())
}
