package socket

import (
	"fmt"
	"github.com/firmeve/firmeve/kernel/contract"
	"github.com/prometheus/common/log"
	"github.com/spf13/cobra"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
	//"os"
	//"os/signal"
	//"syscall"
)

//func c()  {
//	server := net.
//}

type Command struct {
}

func (c Command) CobraCmd() *cobra.Command {
	cmd := new(cobra.Command)
	cmd.Use = "socket"
	cmd.Short = "tcp,upd server"
	return cmd
}

func (c Command) Run(root contract.BaseCommand, cmd *cobra.Command, args []string) {
	var (
		logger = root.Resolve(`logger`).(contract.Loggable)
	)

	logger.Debug(`Start server`)

	go func() {
		server, err := net.Listen("tcp", "0.0.0.0:22122")
		if err != nil {
			logger.Fatal("listen server [0.0.0.0:22122] failed", "error", err)
		}
		defer server.Close()
		for {
			conn, err := server.Accept()
			if err != nil {
				log.Error(err)
				continue
			}
			go func(conn net.Conn) {
				buff := make([]byte, 50)
				for {
					_, err = conn.Read(buff)
					if err != nil {
						log.Error(err)
						continue
					}
					fmt.Printf("%s\n", buff)
					conn.Write([]byte("server " + string(buff) + "\n"))
				}

			}(conn)

			//conn.Close()
		}

		if err != nil {
			logger.Fatal("conn server [0.0.0.0:22122] failed", "error", err)
		}
	}()

	go func() {
		client, err := net.Dial("tcp", "127.0.0.1:22122")
		if err != nil {
			logger.Fatal("client dial server [0.0.0.0:22122] failed", "error", err)
		}
		defer client.Close()
		for {
			client.Write([]byte("hello"))
			time.Sleep(time.Second * 2)
			read := make([]byte, 0)
			client.Read(read)
			fmt.Printf("client %s\n", read)
		}
	}()
	//cron.Start()
	//
	//logger.Debug(`Crontab start successful`)
	//
	//logger.Debug(`Signal listen SIGTERM(kill),SIGINT(kill -2),SIGKILL(kill -9)`)
	//// kill (no param) default send syscall.SIGTERM
	//// kill -2 is syscall.SIGINT
	//// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	<-quit // 阻塞直至有信号传入
	//
	////
	//cron.Stop()
	//logger.Debug(`Stoped crontab`)
}
