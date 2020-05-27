package http

import (
	"context"
	"fmt"
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
	command *cobra.Command
	logger  contract.Loggable
}

func (c *HttpCommand) CobraCmd() *cobra.Command {
	c.command = new(cobra.Command)
	c.command.Use = "http:serve"
	c.command.Short = "Http server"
	c.command.Flags().StringP("host", "H", ":80", "Http serve address")
	c.command.Flags().BoolP("http2", "", false, "Open http2 protocol")
	c.command.Flags().StringP("cert-file", "", "", "Http2 cert file path")
	c.command.Flags().StringP("key-file", "", "", "Http2 key file path")

	return c.command
}

func (c *HttpCommand) Run(root contract.BaseCommand, cmd *cobra.Command, args []string) {
	c.logger = root.Resolve(`logger`).(contract.Loggable)
	var (
		host      = cmd.Flag("host").Value.String()
		certFile  = cmd.Flag(`cert-file`).Value.String()
		keyFile   = cmd.Flag(`key-file`).Value.String()
		openHttp2 = cmd.Flag(`http2`).Value.String()
	)
	srv := &net_http.Server{
		Addr:    host,
		Handler: root.Resolve(`http.router`).(*Router),
	}

	c.debugLog(`Goroutine http server`)

	go func() {
		var err error

		// ssl
		if certFile != `` && keyFile != `` {
			c.debugLog(`Start https server[` + host + `] key[` + keyFile + `] cert[` + certFile + `]`)
			err = srv.ListenAndServeTLS(certFile, keyFile)
			// http2
			if openHttp2 == `true` {
				c.debugLog(`Open http2`)
				//@todo conf is empty
				err = http2.ConfigureServer(srv, &http2.Server{})
			}
		} else {
			c.debugLog(`Start http server[` + host + `]`)
			err = srv.ListenAndServe()
		}

		if err != nil && err != net_http.ErrServerClosed {
			c.logger.Fatal(fmt.Sprintf("listen: %s\n", err))
		}
	}()

	c.logger.Debug(`Signal listen SIGTERM(kill),SIGINT(kill -2),SIGKILL(kill -9)`)
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	<-quit

	c.debugLog("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		c.logger.Fatal("Server Shutdown: ", err)
	}

	c.debugLog("Server exiting")
}

func (c *HttpCommand) debugLog(message string, context ...interface{}) {
	c.logger.Debug(message, context...)
}
