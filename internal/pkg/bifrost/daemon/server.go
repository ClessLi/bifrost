package daemon

import (
	"fmt"
	"github.com/ClessLi/bifrost/internal/pkg/bifrost/config"
	"github.com/ClessLi/bifrost/internal/pkg/utils"
	"github.com/ClessLi/skirnir/pkg/discover"
	"google.golang.org/grpc"
	"net"
	"strings"
)

type Server struct {
	gRPCServer           *grpc.Server
	gRPCServerListenIP   net.IP
	gRPCServerListenPort uint16
	instanceName         string
	instanceId           string
	registryClient       discover.RegistryClient
}

func (s *Server) Start(errChan chan<- error) {

	go func() {
		//启动gRPC服务
		lis, lisErr := net.Listen("tcp", fmt.Sprintf(":%d", s.gRPCServerListenPort))
		if lisErr != nil {
			errChan <- lisErr
			return
		}
		defer lis.Close()

		utils.Logger.NoticeF("bifrost service is running on %s", lis.Addr())
		svrErr := s.gRPCServer.Serve(lis)
		utils.Logger.NoticeF("bifrost service is stopped")
		errChan <- svrErr
	}()

	// 注册服务到注册中心
	go s.registerToRA(errChan)
}

func (s *Server) Stop() {
	s.deregisterToRA()
	s.gRPCServer.Stop()
}

func (s *Server) registerToRA(errChan chan<- error) {
	if s.registryClient == nil {
		return
	}
	instanceHost := s.gRPCServerListenIP.String()

	if !s.registryClient.Register(s.instanceName, s.instanceId, instanceHost, s.gRPCServerListenPort, nil, config.KitLogger(utils.Stdoutf)) {
		err := fmt.Errorf("register service %s failed", s.instanceName)
		utils.Logger.Warning(err.Error())
		errChan <- err
	}
}

func (s *Server) deregisterToRA() {
	if s.registryClient != nil && !strings.EqualFold(s.instanceId, "") {
		if s.registryClient.DeRegister(s.instanceId, config.KitLogger(utils.Stdoutf)) {
			utils.Logger.InfoF("bifrost service (instance ID is '%s') has been unregistered", s.instanceId)
		} else {
			utils.Logger.WarningF("bifrost service (instance ID is '%s') failed to deregister", s.instanceId)
		}
	}
}

func NewServer(gRPCServer *grpc.Server, ip net.IP, port uint16, serviceName, instanceId string, regClient discover.RegistryClient) Server {
	return Server{
		gRPCServer:           gRPCServer,
		gRPCServerListenIP:   ip,
		gRPCServerListenPort: port,
		instanceName:         serviceName,
		instanceId:           instanceId,
		registryClient:       regClient,
	}
}
