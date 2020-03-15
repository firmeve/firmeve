package http

import (
	"context"
	"fmt"
	"github.com/firmeve/firmeve/bootstrap"
	"github.com/firmeve/firmeve/kernel"
	"github.com/firmeve/firmeve/kernel/contract"
	"github.com/spf13/cobra"
	"golang.org/x/net/http2"
	net_http "net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type HttpCommand struct {
	kernel.Command
	command *cobra.Command
}

func (c *HttpCommand) Cmd() *cobra.Command {
	if c.command == nil {
		c.command = c.newCmd()
	}

	return c.command
}

func (c *HttpCommand) newCmd() *cobra.Command {
	c.command = new(cobra.Command)
	c.command.Use = "http:serve"
	c.command.Short = "Http server"
	c.command.Flags().StringP("host", "H", ":80", "Http serve address")
	c.command.Flags().BoolP("http2", "", false, "Open http2 protocol")
	c.command.Flags().StringP("cert-file", "", "", "Http2 cert file path")
	c.command.Flags().StringP("key-file", "", "", "Http2 key file path")

	c.command.Run = c.run

	return c.command
}

func (c *HttpCommand) run(cmd *cobra.Command, args []string) {
	// bootstrap
	bootstrap.Boot(c)
	//config := c.Firmeve.Get(`config`).(*config.Config)
	logger := c.Firmeve.Get(`logger`).(contract.Loggable)
	var (
		host      = cmd.Flag("host").Value.String()
		certFile  = cmd.Flag(`cert-file`).Value.String()
		keyFile   = cmd.Flag(`key-file`).Value.String()
		openHttp2 = cmd.Flag(`http2`).Value.String()
	)
	srv := &net_http.Server{
		Addr:    host,
		Handler: c.Firmeve.Get(`http.router`).(*Router),
	}

	go func() {
		var err error

		// ssl
		if certFile != `` && keyFile != `` {

			err = srv.ListenAndServeTLS(certFile, keyFile)

			// http2
			if openHttp2 == `true` {
				//@todo conf is empty
				err = http2.ConfigureServer(srv, &http2.Server{})
			}
		} else {
			err = srv.ListenAndServe()
		}

		if err != nil && err != net_http.ErrServerClosed {
			logger.Fatal(fmt.Sprintf("listen: %s\n", err))
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	//log.Println("Shutdown Server ...")
	logger.Info("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server Shutdown: ", err)
	}

	logger.Info("Server exiting")
}
