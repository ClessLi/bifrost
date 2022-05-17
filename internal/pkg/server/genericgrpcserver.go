package server

import (
	"context"
	"errors"
	"fmt"
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

	"github.com/yongPhone/bifrost/internal/pkg/service_register"
	healthzclient_v1 "github.com/yongPhone/bifrost/pkg/client/grpc_health_v1"

	//"github.com/ClessLi/bifrost/internal/pkg/middleware".
	log "github.com/yongPhone/bifrost/pkg/log/v1"
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
		log.Fatal("Failed to setup generic gRPC server, no serving information is available for setup.")

		return
	}
	if s.SecureServingInfo != nil {
		cerds, err := credentials.NewServerTLSFromFile(
			s.SecureServingInfo.CertKey.CertFile,
			s.SecureServingInfo.CertKey.KeyFile,
		)
		if err != nil {
			log.Fatal(err.Error())

			return
		}
		s.secureServer = grpc.NewServer(grpc.MaxSendMsgSize(s.ChunkSize), grpc.Creds(cerds))
		log.Infof("Secure server initialization succeeded. Chunk size: %d.", s.ChunkSize)
		if s.healthz {
			s.secureSvrHealthz = health.NewServer()
		}
	}

	if s.InsecureServingInfo != nil {
		// TODO: Checking mechanism of gRPC server max send msg size
		// s.insecureServer = grpc.NewServer(grpc.MaxSendMsgSize(s.ChunkSize))
		s.insecureServer = grpc.NewServer()
		log.Infof("Insecure server initialization succeeded. Chunk size: %d.", s.ChunkSize)
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
			log.Debug("Register health server for insecure gRPC server...")
			grpc_health_v1.RegisterHealthServer(s.insecureServer, s.insecureSvrHealthz)
			log.Debug("Registered health server for insecure gRPC server succeeded.")
		}

		if s.secureServer != nil {
			log.Debug("Register health server for secure gRPC server...")
			grpc_health_v1.RegisterHealthServer(s.secureServer, s.secureSvrHealthz)
			log.Debug("Registered health server for secure gRPC server succeeded.")
		}
	}
}

func (s *GenericGRPCServer) InstallRAClient() {
	if s.RAInfo != nil {
		var err error
		s.registryClient, err = discover.NewKitConsulRegistryClient(s.RAInfo.Host, uint16(s.RAInfo.Port))
		if err != nil {
			log.Fatal(err.Error())

			return
		}
		id, err := uuid.NewV4()
		if err != nil {
			log.Fatalf(err.Error())

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
			log.Warn("service register is nil")

			return
		}
		svcRegister(grpcSvr, healthSvr)
		if s.registryClient != nil {
			log.Debugf("Register service `%s` to RA for gRPC server...", svcname+"-"+s.instanceSuffixID)

			s.registryClient.Register(
				svcname,
				svcname+"-"+s.instanceSuffixID,
				bindHost,
				bindPort,
				nil,
				log.K(),
			)
		}
	}
	for svcname, register := range registers {
		if s.InsecureServingInfo != nil {
			log.Debugf("Register service `%s` for insecure gRPC server...", svcname)
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
			log.Debugf("Register service `%s` for secure gRPC server...", svcname)
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

	log.Debug("Setup generic gRPC server...")
	s.Setup()
	log.Debug("Install RA client for generic gRPC server...")
	s.InstallRAClient()
	log.Debug("Install middlewares for generic gRPC server...")
	s.InstallMiddlewares()
	log.Debug("Install services for generic gRPC server...")
	s.InstallServices()
}

// Run spawns the http server. It only returns when the port cannot be listened on initially.
func (s *GenericGRPCServer) Run() error {
	var eg errgroup.Group

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	log.Debugf("Goroutine start the insecure gRPC server...")
	// TODO: fix duplicate code for insecure server and secure server startup
	eg.Go(func() error { //nolint:dupl
		if s.InsecureServingInfo == nil {
			log.Info("Pass start insecure server")

			return nil
		}

		log.Infof("Start to listening the incoming requests on address: %s", s.InsecureServingInfo.Address())
		lis, err := net.Listen("tcp", s.InsecureServingInfo.Address())
		if err != nil {
			log.Errorf(
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
			log.Errorf("Failed to serve the insecure server, %s", err.Error())

			return err
		}

		log.Infof("Server on %s stopped", s.InsecureServingInfo.Address())

		return nil
	})

	log.Debugf("Goroutine start the secure gRPC server...")
	eg.Go(func() error { //nolint:dupl
		if s.SecureServingInfo == nil {
			log.Info("Pass start secure server")

			return nil
		}
		log.Infof("Start to listening the incoming requests on https address: %s", s.SecureServingInfo.Address())
		lis, err := net.Listen("tcp", s.SecureServingInfo.Address())
		if err != nil {
			log.Errorf(
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
			log.Errorf("Failed to serve the secure server, %s", err.Error())

			return err
		}

		log.Infof("Server on %s stopped", s.SecureServingInfo.Address())

		return nil
	})

	// Ping the server to make sure the router is working.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if s.healthz {
		log.Debugf("Start to run the health check...")
		if err := s.ping(ctx); err != nil {
			return err
		}
	}

	if err := eg.Wait(); err != nil {
		log.Errorf("%+v", err.Error())

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
			s.registryClient.DeRegister(servicename+"-"+s.instanceSuffixID, log.K())
		}
	}

	if s.secureServer != nil {
		s.secureServer.GracefulStop()
		if s.healthz {
			s.secureSvrHealthz.Shutdown()
		}
		log.Info("Secure server has been stopped!")
	}

	if s.insecureServer != nil {
		s.insecureServer.GracefulStop()
		if s.healthz {
			s.insecureSvrHealthz.Shutdown()
		}
		log.Info("Insecure server has been stopped!")
	}
}

// ping pings the http server to make sure the router is working.
//nolint:funlen,gocognit
func (s *GenericGRPCServer) ping(ctx context.Context) error {
	healthzClients := make(map[string]*healthzclient_v1.Client)
	if s.InsecureServingInfo != nil {
		log.Debugf("initialization insecure server health check...")
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
				log.Warn(err.Error())
			}
		}(client)

		healthzClients["insecure"] = client
	}

	if s.SecureServingInfo != nil { //nolint:nestif
		log.Debugf("initialization secure server health check...")
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
				log.Warn(err.Error())
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
					log.Info(err.Error())
					healthzOK[tag+" "+svcname] = false

					continue
				}
				healthzOK[tag+" "+svcname] = true
				log.Infof("The '%s' %s-service state is: %v", svcname, tag, healthzclient_v1.StatusString(status))
			}
		}

		allOK := true
		for tSvcname, ok := range healthzOK {
			if !ok {
				allOK = false
				log.Debugf("service `%s` is not healthy.", tSvcname)
			}
		}

		if allOK {
			log.Infof("all services are healthy!")

			return nil
		}

		// Sleep for a second to continue the next ping.
		log.Info("Waiting for the router, retry in 1 second.")
		time.Sleep(1 * time.Second)

		select {
		case <-ctx.Done():
			log.Fatal("can not check grpc server health within the specified time interval.")
		default:
			if len(healthzClients) == 0 {
				log.Fatal("can not check grpc server health")
			}
		}
	}
	// return fmt.Errorf("the router has no response, or it might took too long to start up")
}
