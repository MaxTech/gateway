package grpc_gateway

import (
    "github.com/maxtech/log"
    "google.golang.org/grpc"
    "net"
    "sync"
    "time"
)

type GrpcServer struct {
    mu     sync.Mutex
    lis    net.Listener
    server *grpc.Server
    status chan int
}

func Init(initLogger log.AppLogger, address string) *GrpcServer {
    var err error
    logger = initLogger
    server := new(GrpcServer)
    server.status = make(chan int, 1)
    server.lis, err = net.Listen("tcp", address)
    if err != nil {
        logger.Error("failed to listen:", err.Error())
    }
    server.server = grpc.NewServer()

    return server
}

func (gs *GrpcServer) GetGrpcServer() *grpc.Server {
    return gs.server
}

func (gs *GrpcServer) Start() {
    for {
        if len(gs.status) < 1 {
            gs.mu.Lock()
            gs.status <- 1
            gs.mu.Unlock()
            go gs.run(gs.server, gs.lis)
        }
        time.Sleep(time.Second)
    }
}

func (gs *GrpcServer) run(s *grpc.Server, lis net.Listener) {
    defer func() {
        if err := recover(); err != nil {
            logger.Error("recover error:", err)
            gs.mu.Lock()
            <-gs.status
            gs.mu.Unlock()
        }
    }()

    if err := s.Serve(lis); err != nil {
        logger.Error("failed to serve:", err.Error())
    }
}
