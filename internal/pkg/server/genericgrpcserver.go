package server

import (
	"context"
	"errors"
	"fmt"
	logV1 "github.com/ClessLi/component-base/pkg/log/v1"
	"net"
	"strings"
	"time"

	"github.com/ClessLi/skirnir/pkg/discover"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	"github.com/ClessLi/bifrost/internal/pkg/service_register"
	healthzclient_v1 "github.com/ClessLi/bifrost/pkg/client/grpc_health_v1"
	//"github.com/ClessLi/bifrost/internal/pkg/middleware".
)

// GenericGRPCServer contains state for a bifrost api server.
type GenericGRPCServer struct {
	middlewares []string
	// SecureServingInfo holds configuration of the TLS server.
	SecureServingInfo *SecureServingInfo

	// InsecureServingInfo holds configuration of the insecure grpc server.
	InsecureServingInfo *InsecureServingInfo

	// RAInfo holds configuration of th Registration Authority server.
	RAInfo *RAInfo

	// ChunkSize set chunk size of the grpc server.
	ChunkSize int

	// ReceiveTimeout is the timeout used for the grpc server receiving data.
	ReceiveTimeout time.Duration

	// ShutdownTimeout is the timeout used for server shutdown. This specifies the timeout before server
	// gracefully shutdown returns.
	ShutdownTimeout time.Duration

	healthz         bool
	enableMetrics   bool
	enableProfiling bool

	instanceSuffixID string
	registryClient   discover.RegistryClient

	insecureSvrHealthz, secureSvrHealthz *health.Server
	insecureServer, secureServer         *grpc.Server
	registeredService                    []string
}

func (s *GenericGRPCServer) Setup() {
	if s.SecureServingInfo == nil && s.InsecureServingInfo == nil {
		logV1.Fatal("Failed to setup generic gRPC server, no serving information is available for setup.")

		return
	}
	if s.SecureServingInfo != nil {
		cerds, err := credentials.NewServerTLSFromFile(
			s.SecureServingInfo.CertKey.CertFile,
			s.SecureServingInfo.CertKey.KeyFile,
		)
		if err != nil {
			logV1.Fatal(err.Error())

			return
		}
		s.secureServer = grpc.NewServer(grpc.Creds(cerds))
		if s.healthz {
			s.secureSvrHealthz = health.NewServer()
		}
	}

	if s.InsecureServingInfo != nil {
		s.insecureServer = grpc.NewServer()
		if s.healthz {
			s.insecureSvrHealthz = health.NewServer()
		}
	}
}

func (s *GenericGRPCServer) InstallMiddlewares() {
	// necessary middlewares
	// s.Use(limits.RequestSizeLimiter(10))

	// install custom middlewares
	// for _, m := range s.middlewares {
	// mw, ok := middleware.Middlewares[m]
	// if !ok {
	//	log.Warnf("can not find middleware: %s", m)

	//continue
	//}

	//log.Infof("install middleware: %s", m)
	//s.Use(mw)
	//}
}

func (s *GenericGRPCServer) InstallServices() {
	// check healthz handler
	if s.healthz {
		if s.insecureServer != nil {
			logV1.Debug("Register health server for insecure gRPC server...")
			grpc_health_v1.RegisterHealthServer(s.insecureServer, s.insecureSvrHealthz)
			logV1.Debug("Registered health server for insecure gRPC server succeeded.")
		}

		if s.secureServer != nil {
			logV1.Debug("Register health server for secure gRPC server...")
			grpc_health_v1.RegisterHealthServer(s.secureServer, s.secureSvrHealthz)
			logV1.Debug("Registered health server for secure gRPC server succeeded.")
		}
	}
}

func (s *GenericGRPCServer) InstallRAClient() {
	if s.RAInfo != nil {
		var err error
		s.registryClient, err = discover.NewKitConsulRegistryClient(s.RAInfo.Host, uint16(s.RAInfo.Port))
		if err != nil {
			logV1.Fatal(err.Error())

			return
		}
		id, err := uuid.NewV4()
		if err != nil {
			logV1.Fatalf(err.Error())

			return
		}
		s.instanceSuffixID = id.String()
	}
}

func (s *GenericGRPCServer) RegisterServices(registers map[string]service_register.ServiceRegister) {
	f := func(
		svcname string,
		grpcSvr *grpc.Server,
		healthSvr *health.Server,
		svcRegister service_register.ServiceRegister,
		bindHost string,
		bindPort uint16,
	) {
		if svcRegister == nil {
			logV1.Warn("service register is nil")

			return
		}
		svcRegister(grpcSvr, healthSvr)
		if s.registryClient != nil {
			logV1.Debugf("Register service `%s` to RA for gRPC server...", svcname+"-"+s.instanceSuffixID)

			s.registryClient.Register(
				svcname,
				svcname+"-"+s.instanceSuffixID,
				bindHost,
				bindPort,
				nil,
				logV1.K(),
			)
		}
	}
	for svcname, register := range registers {
		if s.InsecureServingInfo != nil {
			logV1.Debugf("Register service `%s` for insecure gRPC server...", svcname)
			f(
				svcname,
				s.insecureServer,
				s.insecureSvrHealthz,
				register,
				s.InsecureServingInfo.BindAddress,
				uint16(s.InsecureServingInfo.BindPort),
			)
		}
		if s.SecureServingInfo != nil {
			logV1.Debugf("Register service `%s` for secure gRPC server...", svcname)
			f(
				svcname,
				s.secureServer,
				s.secureSvrHealthz,
				register,
				s.SecureServingInfo.BindAddress,
				uint16(s.SecureServingInfo.BindPort),
			)
		}

		s.registeredService = append(s.registeredService, svcname)
	}
}

func initGenericGRPCServer(s *GenericGRPCServer) {
	// do some setup

	logV1.Debug("Setup generic gRPC server...")
	s.Setup()
	logV1.Debug("Install RA client for generic gRPC server...")
	s.InstallRAClient()
	logV1.Debug("Install middlewares for generic gRPC server...")
	s.InstallMiddlewares()
	logV1.Debug("Install services for generic gRPC server...")
	s.InstallServices()
}

// Run spawns the http server. It only returns when the port cannot be listened on initially.
func (s *GenericGRPCServer) Run() error {
	var eg errgroup.Group

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	logV1.Debugf("Goroutine start the insecure gRPC server...")
	// TODO: fix duplicate code for insecure server and secure server startup
	eg.Go(func() error { //nolint:dupl
		if s.InsecureServingInfo == nil {
			logV1.Info("Pass start insecure server")

			return nil
		}

		logV1.Infof("Start to listening the incoming requests on address: %s", s.InsecureServingInfo.Address())
		lis, err := net.Listen("tcp", s.InsecureServingInfo.Address())
		if err != nil {
			logV1.Errorf(
				"Failed to listen the incoming requests on address: %s, %s",
				s.InsecureServingInfo.Address(),
				err.Error(),
			)

			return err
		}

		if s.healthz {
			s.insecureSvrHealthz.Resume()
			defer s.insecureSvrHealthz.Shutdown()
		}

		if err = s.insecureServer.Serve(lis); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			logV1.Errorf("Failed to serve the insecure server, %s", err.Error())

			return err
		}

		logV1.Infof("Server on %s stopped", s.InsecureServingInfo.Address())

		return nil
	})

	logV1.Debugf("Goroutine start the secure gRPC server...")
	eg.Go(func() error { //nolint:dupl
		if s.SecureServingInfo == nil {
			logV1.Info("Pass start secure server")

			return nil
		}
		logV1.Infof("Start to listening the incoming requests on https address: %s", s.SecureServingInfo.Address())
		lis, err := net.Listen("tcp", s.SecureServingInfo.Address())
		if err != nil {
			logV1.Errorf(
				"Failed to listen the incoming requests on address: %s, %s",
				s.SecureServingInfo.Address(),
				err.Error(),
			)

			return err
		}

		if s.healthz {
			s.secureSvrHealthz.Resume()
			defer s.secureSvrHealthz.Shutdown()
		}
		if err = s.secureServer.Serve(lis); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			logV1.Errorf("Failed to serve the secure server, %s", err.Error())

			return err
		}

		logV1.Infof("Server on %s stopped", s.SecureServingInfo.Address())

		return nil
	})

	// Ping the server to make sure the router is working.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if s.healthz {
		logV1.Debugf("Start to run the health check...")
		if err := s.ping(ctx); err != nil {
			return err
		}
	}

	if err := eg.Wait(); err != nil {
		logV1.Errorf("%+v", err.Error())

		return err
	}

	return nil
}

// Close graceful shutdown the api server.
func (s *GenericGRPCServer) Close() {
	// The context is used to inform the server it has 10 seconds to finish
	// the request it is currently handling
	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// defer cancel()

	if s.registryClient != nil {
		for _, servicename := range s.registeredService {
			s.registryClient.DeRegister(servicename+"-"+s.instanceSuffixID, logV1.K())
		}
	}

	if s.secureServer != nil {
		s.secureServer.GracefulStop()
		if s.healthz {
			s.secureSvrHealthz.Shutdown()
		}
		logV1.Info("Secure server has been stopped!")
	}

	if s.insecureServer != nil {
		s.insecureServer.GracefulStop()
		if s.healthz {
			s.insecureSvrHealthz.Shutdown()
		}
		logV1.Info("Insecure server has been stopped!")
	}
}

// ping pings the http server to make sure the router is working.
//
//nolint:funlen,gocognit
func (s *GenericGRPCServer) ping(ctx context.Context) error {
	healthzClients := make(map[string]*healthzclient_v1.Client)
	if s.InsecureServingInfo != nil {
		logV1.Debugf("initialization insecure server health check...")
		var address string
		if strings.Contains(s.InsecureServingInfo.Address(), "0.0.0.0") {
			address = fmt.Sprintf("127.0.0.1:%s", strings.Split(s.InsecureServingInfo.Address(), ":")[1])
		} else {
			address = s.InsecureServingInfo.Address()
		}

		client, err := healthzclient_v1.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return err
		}

		defer func(client *healthzclient_v1.Client) {
			err := client.Close()
			if err != nil {
				logV1.Warn(err.Error())
			}
		}(client)

		healthzClients["insecure"] = client
	}

	if s.SecureServingInfo != nil { //nolint:nestif
		logV1.Debugf("initialization secure server health check...")
		var address string
		if strings.Contains(s.SecureServingInfo.Address(), "0.0.0.0") {
			address = fmt.Sprintf("127.0.0.1:%s", strings.Split(s.SecureServingInfo.Address(), ":")[1])
		} else {
			address = s.SecureServingInfo.Address()
		}
		creds, err := credentials.NewClientTLSFromFile(
			s.SecureServingInfo.CertKey.CertFile,
			s.SecureServingInfo.BindAddress,
		)
		if err != nil {
			return err
		}
		client, err := healthzclient_v1.NewClient(address, grpc.WithTransportCredentials(creds))
		if err != nil {
			return err
		}

		defer func(client *healthzclient_v1.Client) {
			err := client.Close()
			if err != nil {
				logV1.Warn(err.Error())
			}
		}(client)

		healthzClients["secure"] = client
	}

	healthzOK := make(map[string]bool)

	for {
		for tag, client := range healthzClients {
			for _, svcname := range s.registeredService {
				if ok, has := healthzOK[tag+" "+svcname]; has && ok {
					continue
				}

				status, err := client.Check(ctx, svcname)
				if err != nil {
					logV1.Info(err.Error())
					healthzOK[tag+" "+svcname] = false

					continue
				}
				healthzOK[tag+" "+svcname] = true
				logV1.Infof("The '%s' %s-service state is: %v", svcname, tag, healthzclient_v1.StatusString(status))
			}
		}

		allOK := true
		for tSvcname, ok := range healthzOK {
			if !ok {
				allOK = false
				logV1.Debugf("service `%s` is not healthy.", tSvcname)
			}
		}

		if allOK {
			logV1.Infof("all services are healthy!")

			return nil
		}

		// Sleep for a second to continue the next ping.
		logV1.Info("Waiting for the router, retry in 1 second.")
		time.Sleep(1 * time.Second)

		select {
		case <-ctx.Done():
			logV1.Fatal("can not check grpc server health within the specified time interval.")
		default:
			if len(healthzClients) == 0 {
				logV1.Fatal("can not check grpc server health")
			}
		}
	}
	// return fmt.Errorf("the router has no response, or it might took too long to start up")
}
