package bifrost

import (
	"context"
	"fmt"
	nginx_ctx "github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context/local"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context_type"
	"sync"
	"testing"
	"time"

	"google.golang.org/grpc"

	v1 "github.com/ClessLi/bifrost/api/bifrost/v1"
	bifrost_cliv1 "github.com/ClessLi/bifrost/pkg/client/bifrost/v1"
	healthzclient_v1 "github.com/ClessLi/bifrost/pkg/client/grpc_health_v1"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration"
)

func TestRun(t *testing.T) {
	err := exampleServerRun()
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestBifrostClient(t *testing.T) {
	healthCli, err := healthzclient_v1.NewClient(serverAddress(), grpc.WithInsecure())
	if err != nil {
		t.Fatalf(err.Error())
	}

	retryI := 0
	for {
		if retryI >= 10 {
			t.Fatalf("connect to web server config service timeout.")
		}
		state, err := healthCli.Check(context.Background(), "com.github.ClessLi.api.bifrostpb.WebServerConfig")
		if err != nil {
			t.Log(err.Error())
		}
		if state == healthzclient_v1.SERVING {
			t.Log("service is serving")
			break
		}
		time.Sleep(time.Second * 3)
		retryI++
	}

	client, err := bifrost_cliv1.New(serverAddress(), grpc.WithInsecure(), grpc.WithTimeout(time.Second))
	if err != nil {
		t.Fatalf(err.Error())
	}

	defer client.Close()

	servernames, err := client.WebServerConfig().GetServerNames()
	if err != nil {
		t.Fatalf(err.Error())
	}

	time.Sleep(time.Second * 10)
	metrics, err := client.WebServerStatus().Get()
	if err != nil {
		t.Fatalf("%++v", err)
	}
	t.Log(metrics)

	// normal grpc client
	/*cclient, err := grpc.Dial(serverAddress(), grpc.WithInsecure(), grpc.WithTimeout(time.Second))
	if err != nil {
		t.Fatalf(err.Error())
	}

	defer cclient.Close()
	c := pbv1.NewWebServerConfigClient(cclient)*/
	wg := new(sync.WaitGroup)
	for _, servername := range servernames {
		// normal grpc client
		/*resp, err := c.Get(context.Background(), &pbv1.ServerName{Name: servername})
		if err != nil {
			t.Fatalf(err.Error())
		}
		buf := bytes.NewBuffer(nil)
		stop := false
		for !stop {
			select {
			case <-resp.Context().Done():
				stop = true
				break
			default:
				conf, err := resp.Recv()
				if err != nil && err != io.EOF {
					t.Fatalf(err.Error())
				}
				buf.Write(conf.GetJsonData())
				if err == io.EOF {
					stop = true
					break
				}
			}
		}
		t.Logf("config:\n\n%s", buf.String())*/

		// go-kit grpc client
		jsondata, err := client.WebServerConfig().Get(servername)
		if err != nil {
			t.Fatalf(err.Error())
		}
		conf, err := configuration.NewNginxConfigFromJsonBytes(jsondata)
		if err != nil {
			t.Fatalf(err.Error())
		}
		lines, err := conf.Main().ConfigLines(false)
		if err != nil {
			t.Fatalf(err.Error())
		}
		t.Logf("get config lines len: %d", len(lines))
		for _, line := range lines {
			t.Log(line)
		}
		//fmt.Printf("config %s:\n\n%s", servername, conf.View())
		t.Logf("before jsondata len: %d, after jasondata len: %d", len(jsondata), len(conf.Json()))

		statistics, err := client.WebServerStatistics().Get(servername)
		if err != nil {
			t.Fatalf(err.Error())
		}
		t.Logf("statistics %s:\n\n%+v", servername, statistics)

		logC, lwCancel, err := client.WebServerLogWatcher().Watch(&v1.WebServerLogWatchRequest{
			ServerName:          &v1.ServerName{Name: servername},
			LogName:             "access.log",
			FilteringRegexpRule: "^test.*$",
		})
		if err != nil {
			t.Fatalf(err.Error())
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			defer lwCancel()
			for {
				select {
				case <-time.After(time.Second * 10):
					return
				case line := <-logC:
					if line == nil {
						return
					}
					fmt.Print(string(line))
				}
			}
		}()
	}
	wg.Wait()
}

func TestBifrostClientOperation(t *testing.T) {
	client, err := bifrost_cliv1.New(serverAddress(), grpc.WithInsecure(), grpc.WithTimeout(time.Second))
	if err != nil {
		t.Fatal(err)
	}

	defer client.Close()

	servernames, err := client.WebServerConfig().GetServerNames()
	if err != nil {
		t.Fatal(err)
	}

	for _, servername := range servernames {
		jsondata, err := client.WebServerConfig().Get(servername)
		if err != nil {
			t.Fatal(err)
		}
		conf, err := configuration.NewNginxConfigFromJsonBytes(jsondata)
		if err != nil {
			t.Fatal(err)
		}
		for _, pos := range conf.Main().QueryByKeyWords(nginx_ctx.NewKeyWords(context_type.TypeHttp)).Target().
			QueryByKeyWords(nginx_ctx.NewKeyWords(context_type.TypeServer)).Target().
			QueryAllByKeyWords(nginx_ctx.NewKeyWords(context_type.TypeDirective).SetStringMatchingValue("server_name test1.com")) {
			server, _ := pos.Position()
			if server.QueryByKeyWords(nginx_ctx.NewKeyWords(context_type.TypeDirective).SetRegexpMatchingValue("^listen 80$")).Target().Error() != nil {
				continue
			}
			ctx, idx := server.QueryByKeyWords(nginx_ctx.NewKeyWords(context_type.TypeLocation).SetRegexpMatchingValue(`^/test1-location$`)).Target().
				QueryByKeyWords(nginx_ctx.NewKeyWords(context_type.TypeIf).SetRegexpMatchingValue(`^\(\$http_api_name != ''\)$`)).Target().
				QueryByKeyWords(nginx_ctx.NewKeyWords(context_type.TypeDirective).SetStringMatchingValue("proxy_pass")).Position()
			err = ctx.Insert(local.NewContext(context_type.TypeInlineComment, fmt.Sprintf("[%s]test comments", time.Now().String())), idx+1).Error()
			if err != nil {
				t.Fatal(err)
			}
		}
		err = client.WebServerConfig().Update(servername, conf.Json())
		if err != nil {
			t.Fatal(err)
		}
	}
}
