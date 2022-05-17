package bifrost

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	v1 "github.com/yongPhone/bifrost/api/bifrost/v1"
	healthzclient_v1 "github.com/yongPhone/bifrost/pkg/client/grpc_health_v1"
	"github.com/yongPhone/bifrost/pkg/resolv/V2/nginx/configuration"

	"google.golang.org/grpc"

	bifrost_cliv1 "github.com/yongPhone/bifrost/pkg/client/bifrost/v1"
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
		conf, err := configuration.NewConfigurationFromJsonBytes(jsondata)
		if err != nil {
			t.Fatalf(err.Error())
		}
		t.Logf("get config len: %d", len(conf.View()))
		//fmt.Printf("config %s:\n\n%s", servername, conf.View())
		t.Logf("before jsondata len: %d, after jasondata len: %d", len(jsondata), len(conf.Json()))

		statistics, err := client.WebServerStatistics().Get(servername)
		if err != nil {
			t.Fatalf(err.Error())
		}
		t.Logf("statistics %s:\n\n%v", servername, statistics)

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
