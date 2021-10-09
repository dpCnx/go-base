package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"runtime/debug"
	"time"

	pb "go-base/grpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type server struct {
	pb.UnimplementedHelloServerServer
}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {

	data := make(chan *pb.HelloResponse)

	go func(request *pb.HelloRequest) {
		name := request.GetName()
		data <- &pb.HelloResponse{Msg: fmt.Sprintf("hi %s", name)}
	}(in)

	select {
	case res := <-data:
		return res, nil
	case <-ctx.Done():
		return nil, status.Errorf(codes.Canceled, "time out")
	}
}

func (s *server) SendStream(in *pb.HelloRequest, stream pb.HelloServer_SendStreamServer) error {

	for i := 0; i < 10; i++ {
		if err := stream.Send(&pb.HelloResponse{Msg: fmt.Sprintf("hi %s", in.GetName())}); err != nil {
			return err
		}
	}

	return nil
}

func (s *server) GetStream(stream pb.HelloServer_GetStreamServer) error {

	for {
		request, err := stream.Recv()
		if err != nil {
			if io.EOF == err {
				return stream.SendAndClose(&pb.HelloResponse{Msg: fmt.Sprintf("hi %s over", request.GetName())})
			}
		}

		log.Printf("recv request:%v \n", request)
	}
}

func (s *server) GetAndSend(stream pb.HelloServer_GetAndSendServer) error {

	for {
		request, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		log.Printf("recv request:%v \n", request)

		if err = stream.Send(&pb.HelloResponse{Msg: fmt.Sprintf("hi %s", request.GetName())}); err != nil {
			return err
		}
	}

}

func authInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	m, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "token is nil")
	}
	if val, ok := m["token"]; ok {
		log.Println("token", val[0])
	}

	return handler(ctx, req)
}

func loggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	start := time.Now()
	log.Printf("server:%v  method:%s request:%v \n", info.Server, info.FullMethod, req)

	resp, err := handler(ctx, req)

	cost := time.Since(start)
	log.Printf("time:%d server:%v  method:%s request:%v respose:%v \n", cost, info.Server, info.FullMethod, req, resp)

	return resp, err
}

func recoveryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {

	defer func() {
		if e := recover(); e != nil {
			log.Printf("error:%v \n", e)
			log.Printf("stack:%s \n", string(debug.Stack()))
			err = status.Errorf(codes.Internal, "Panic err: %v", e)
		}
	}()

	return handler(ctx, req)

}

func streamInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {

	m, ok := metadata.FromIncomingContext(ss.Context())
	if !ok {
		return status.Errorf(codes.Unauthenticated, "token is nil")
	}
	if val, ok := m["token"]; ok {
		log.Println("token", val[0])
	}
	return handler(srv, ss)
}

func main() {
	ln, err := net.Listen("tcp", "127.0.0.1:9999")
	if err != nil {
		log.Printf("net listen err: %v \n", err)
		return
	}

	/*	// tls
		c, err := credentials.NewServerTLSFromFile("demo/grpc/tls/cert.pem", "demo/grpc/tls/cert.key")
		if err != nil {
			log.Printf("credentials newservertlsfromfile err: %v \n", err)
			return
		}
		s := grpc.NewServer(grpc.Creds(c))*/

	/*	// interceptor
		opts := []grpc.ServerOption{grpc.ChainUnaryInterceptor(recoveryInterceptor, loggingInterceptor, authInterceptor), grpc.ChainStreamInterceptor(streamInterceptor)}
		s := grpc.NewServer(opts...)*/

	s := grpc.NewServer()
	pb.RegisterHelloServerServer(s, &server{})
	if err = s.Serve(ln); err != nil {
		log.Printf("s serve err: %v \n", err)
		return
	}
}

/*
	protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    helloworld/helloworld.proto

	openssl genrsa > cert.key //生成私钥
	openssl req -new -x509 -sha256 -key cert.key -out cert.pem -days 3650 //生成证书
*/
