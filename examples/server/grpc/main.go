package main

import (
	"context"
	"fmt"
	"github.com/firmeve/firmeve"
	hello "github.com/firmeve/firmeve/examples/server/grpc/protos"
	grpc2 "github.com/firmeve/firmeve/grpc"
	"github.com/firmeve/firmeve/kernel"
	"github.com/firmeve/firmeve/kernel/contract"
	"google.golang.org/grpc"
	"net"
	"time"
)

type App struct {
	kernel.BaseProvider
}

func (a *App) Name() string {
	return `app`
}

func (a *App) Register() {
	router := a.Firmeve.Get(`grpc.server.router`).(*grpc2.Router)
	//router.Register(hello.RegisterHelloServiceServer, new(HelloService))
	router.Register(func(server *grpc.Server) {
		helloService := new(HelloService)
		hello.RegisterHelloServiceServer(server, helloService)
	})
}
func (a *App) Boot() {
}

type HelloService struct {
}

func (h HelloService) Hello(context.Context, *hello.Search) (*hello.Response, error) {
	return &hello.Response{Result: true}, nil
}

func main() {
	go func() {
		time.Sleep(time.Second)
		client()
	}()
	firmeve.RunDefault(firmeve.WithProviders(
		[]contract.Provider{
			new(grpc2.Provider),
			new(App),
		},
	), firmeve.WithCommands([]contract.Command{
		new(grpc2.GRPCCommand),
	}))

	//var wg sync.WaitGroup
	//wg.Add(2)
	//go func() {
	//	server()
	//	wg.Done()
	//}()
	//go func() {
	//	client()
	//}()
	//time.Sleep(time.Second * 30)
}

func server() {
	grpcServer := grpc.NewServer()

	// register service
	helloService := new(HelloService)
	hello.RegisterHelloServiceServer(grpcServer, helloService)

	conn, err := net.Listen("tcp", ":20020")
	if err != nil {
		panic(err)
	}

	err = grpcServer.Serve(conn)
	if err != nil {
		panic(err)
	}
}

func client() {
	// client
	conn, err := grpc.Dial("localhost:20020", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	client := hello.NewHelloServiceClient(conn)
	search := new(hello.Search)
	v, err := client.Hello(context.Background(), search)
	if err != nil {
		panic(err)
	}
	fmt.Println(v.GetResult())
}
