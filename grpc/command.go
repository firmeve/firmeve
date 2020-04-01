package grpc

import (
	"github.com/firmeve/firmeve/kernel/contract"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"net"
)

type GRPCCommand struct {
}

func (G GRPCCommand) CobraCmd() *cobra.Command {
	command := new(cobra.Command)
	command.Use = "grpc:serve"
	command.Short = "GRPC server"
	command.Flags().StringP("host", "H", ":80", "Http serve address")
	command.Flags().StringP("cert-file", "", "", "Http2 cert file path")
	command.Flags().StringP("key-file", "", "", "Http2 key file path")

	return command
}

func (G GRPCCommand) Run(root contract.BaseCommand, cmd *cobra.Command, args []string) {
	logger := root.Resolve(`logger`).(contract.Loggable)
	var (
		host       = cmd.Flag("host").Value.String()
		certFile   = cmd.Flag(`cert-file`).Value.String()
		keyFile    = cmd.Flag(`key-file`).Value.String()
		grpcServer *grpc.Server
	)

	if certFile != `` && keyFile != `` {
		// ssl 证书
		credential, err := credentials.NewServerTLSFromFile(certFile, keyFile)
		if err != nil {
			logger.Fatal("ssl error, %v", err)
		}
		grpcServer = grpc.NewServer(grpc.Creds(credential))
	} else {
		grpcServer = grpc.NewServer()
	}

	//开启server
	//, grpc.UnaryInterceptor(filter)

	router := root.Resolve(`grpc.server.router`).(*Router)
	router.Handle(grpcServer)
	//Register service
	//grpc_ssl_protos.RegisterHelloServiceServer(grpcServer, new(HelloService))

	//start listener
	conn, err := net.Listen("tcp", host)
	if err != nil {
		logger.Fatal("listen server error %w", err)
	}

	// start serve
	err = grpcServer.Serve(conn)
	if err != nil {
		logger.Fatal("grpc server start error, %w", err)
	}
}
