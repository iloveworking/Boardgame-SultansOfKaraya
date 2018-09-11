package base

import (
	"google.golang.org/grpc"
	"fmt"
	"net"
)

type Service interface {
	//when service create
	OnCreate()
	//when service destroy
	OnDestroy(err error)

	OnRegisterGrpcServer(svr *grpc.Server)
}

//noinspection GoNameStartsWithPackageName
type BaseServer struct {
	Name     string
	Ip       string
	Port     string
	services []Service
}

func NewBaseServer(name, ip, port string) *BaseServer {
	return &BaseServer{Name: name, Ip: ip, Port: port}
}

func (s *BaseServer) AddService(svr Service) {
	s.services = append(s.services, svr)
}

func (s *BaseServer) Start() error {
	opts := make([]grpc.ServerOption, 0)
	grpcServer := grpc.NewServer(opts...)

	for _, server := range s.services {
		server.OnCreate()
		server.OnRegisterGrpcServer(grpcServer)
	}

	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%s", s.Ip, s.Port))
	if err != nil {
		return err
	}

	err = grpcServer.Serve(listen)

	for _, server := range s.services {
		server.OnDestroy(err)
	}

	return err
}
