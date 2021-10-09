package main

import (
	"context"
	"io"
	"log"
	"time"

	pb "go-base/grpc/proto"
	"google.golang.org/grpc"
)

// AuthToken 实现grpc.PerRPCCredentials接⼝
type AuthToken struct {
	Token string
}

func (a *AuthToken) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {
	return map[string]string{"token": a.Token}, nil
}

// RequireTransportSecurity 是否开启tls验证
func (a *AuthToken) RequireTransportSecurity() bool {
	return true
}

func main() {

	/*
		// 生成证书时：Common Name (eg, fully qualified host name) []:d
		tls, err := credentials.NewClientTLSFromFile("demo/grpc/tls/cert.pem", "d")
		if err != nil {
			log.Printf("credentials newclienttlsfromfile err: %v \n", err)
			return
		}
		conn, err := grpc.Dial("127.0.0.1:9999", grpc.WithTransportCredentials(tls))
	*/

	/*
		// token
		at := &AuthToken{Token: "D"}
		conn, err := grpc.Dial("127.0.0.1:9999", grpc.WithPerRPCCredentials(at),grpc.WithTransportCredentials(tls))
	*/

	conn, err := grpc.Dial("127.0.0.1:9999", grpc.WithInsecure())
	if err != nil {
		log.Printf("grpc dial err: %s \n", err.Error())
		return
	}
	defer conn.Close()
	c := pb.NewHelloServerClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	sayHello(ctx, c)
	sendStream(ctx, c)
	getStream(ctx, c)
	getAndSend(ctx, c)

}

func sayHello(ctx context.Context, c pb.HelloServerClient) {
	res, err := c.SayHello(ctx, &pb.HelloRequest{Name: "d"})
	if err != nil {
		log.Printf("grpc dial err: %s \n", err.Error())
		return
	}
	log.Println("resp: ", res.GetMsg())
}

func sendStream(ctx context.Context, c pb.HelloServerClient) {

	sc, err := c.SendStream(ctx, &pb.HelloRequest{Name: "d"})
	if err != nil {
		log.Printf("sc err: %s \n", err.Error())
		return
	}
	for {
		res, err := sc.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("recv err:%v \n", err)
			return
		}

		log.Printf("resp: %v", res.GetMsg())

	}
}

func getStream(ctx context.Context, c pb.HelloServerClient) {

	sc, err := c.GetStream(ctx)
	if err != nil {
		log.Printf("sc err: %s \n", err.Error())
		return
	}

	for i := 0; i < 10; i++ {
		if err := sc.Send(&pb.HelloRequest{Name: "d"}); err != nil {
			log.Printf("Send err:%v \n", err)
		}
	}

	res, err := sc.CloseAndRecv()
	if err != nil {
		log.Printf("Close err:%v \n", err)
		return
	}
	log.Printf("resp: %s \n", res.GetMsg())

}

func getAndSend(ctx context.Context, c pb.HelloServerClient) {

	sc, err := c.GetAndSend(ctx)
	if err != nil {
		log.Printf("sc err: %s \n", err.Error())
		return
	}

	defer sc.CloseSend()

	for i := 0; i < 10; i++ {
		if err := sc.Send(&pb.HelloRequest{Name: "d"}); err != nil {
			log.Printf("Send err:%v \n", err)
			return
		}

		res, err := sc.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("Recv err:%v \n", err)
			return
		}
		log.Printf("resp: %s \n", res.GetMsg())
	}

}
