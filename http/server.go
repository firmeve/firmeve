package http

import (
	"context"
	"fmt"
	"github.com/firmeve/firmeve/config"
	"github.com/firmeve/firmeve/kernel/contract"
	"github.com/firmeve/firmeve/logger"
	"github.com/spf13/cobra"
	"golang.org/x/net/http2"
	net_http "net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type cmd struct {
	app     contract.Application
	command *cobra.Command
	logger  logging.Loggable
	config  *config.Config
}

func NewServer(firmeve contract.Application) *cmd {
	command := &cmd{
		app:     firmeve,
		command: new(cobra.Command),
	}
	command.config = firmeve.Get(`config`).(*config.Config)
	command.logger = firmeve.Get(`logger`).(logging.Loggable)

	return command
	// base
	//c.config = c.bootstrap.Firmeve.Get(`config`).(*config.Config)
	//c.logger = c.bootstrap.Firmeve.Get(`logger`).(logging.Loggable)
}

func (c *cmd) Cmd() *cobra.Command {
	c.command.Use = "http:serve"
	c.command.Flags().StringP("host", "H", ":80", "Http serve address")
	c.command.Flags().BoolP("http2", "", false, "Open http2 protocol")
	c.command.Flags().StringP("cert-file", "", "", "Http2 cert file path")
	c.command.Flags().StringP("key-file", "", "", "Http2 key file path")

	c.command.Run = func(cmd *cobra.Command, args []string) {
		c.run(cmd, args)
	}

	return c.command
}

func (c *cmd) run(cmd *cobra.Command, args []string) {
	var (
		host      = cmd.Flag("host").Value.String()
		certFile  = cmd.Flag(`cert-file`).Value.String()
		keyFile   = cmd.Flag(`key-file`).Value.String()
		openHttp2 = cmd.Flag(`http2`).Value.String()
	)
	srv := &net_http.Server{
		Addr:    host,
		Handler: c.app.Get(`http.router`).(*Router),
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
			c.logger.Fatal(fmt.Sprintf("listen: %s\n", err))
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
	c.logger.Info("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		c.logger.Fatal("Server Shutdown: ", err)
	}

	c.logger.Info("Server exiting")
}
