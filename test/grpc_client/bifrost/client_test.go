package bifrost

import (
	"context"
	healthzclient_v1 "github.com/ClessLi/bifrost/pkg/client/grpc_health_v1"
	"testing"
	"time"

	"google.golang.org/grpc"

	bifrost_cliv1 "github.com/ClessLi/bifrost/pkg/client/bifrost/v1"
)

func TestRun(t *testing.T) {
	err := exampleServerRun()
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestBifrostClient(t *testing.T) {
	//go func() {
	//	err := exampleServerRun()
	//	if err != nil {
	//		t.Error(err.Error())
	//
	//		return
	//	}
	//}()

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

	// normal grpc client
	/*cclient, err := grpc.Dial(serverAddress(), grpc.WithInsecure(), grpc.WithTimeout(time.Second))
	if err != nil {
		t.Fatalf(err.Error())
	}

	defer cclient.Close()
	c := pbv1.NewWebServerConfigClient(cclient)*/
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
		conf, err := client.WebServerConfig().Get(servername)
		if err != nil {
			t.Fatalf(err.Error())
		}
		t.Logf("config %s:\n\n%s", servername, conf)

		statistics, err := client.WebServerStatistics().Get(servername)
		if err != nil {
			t.Fatalf(err.Error())
		}
		t.Logf("statistics %s:\n\n%v", servername, statistics)
	}
}
