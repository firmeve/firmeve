package http

import (
	"context"
	"github.com/firmeve/firmeve/kernel/contract"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type (
	Command struct {
	}
)

var (
	host                              string
	readTimeout                       time.Duration
	writeTimeout                      time.Duration
	readHeaderTimeout                 time.Duration
	idleTimeout                       time.Duration
	maxHeaderBytes                    int
	certFile                          string
	keyFile                           string
	http2Var                          bool
	http2MaxHandlers                  int
	http2MaxConcurrentStreams         uint32
	http2MaxReadFrameSize             uint32
	http2PermitProhibitedCipherSuites bool
	http2IdleTimeout                  time.Duration
	http2MaxUploadBufferPerConnection int32
	http2MaxUploadBufferPerStream     int32
	watch                             []string
)

var (
	logger contract.Loggable
)

func (c *Command) CobraCmd() *cobra.Command {
	command := new(cobra.Command)
	command.Use = "http:serve"
	command.Short = "Http server"
	command.Flags().StringVarP(&host, "host", "H", ":80", "http serve address")
	command.Flags().DurationVar(&readTimeout, "read-timeout", time.Minute, "read timeout default a minute")
	command.Flags().DurationVar(&writeTimeout, "write-timeout", time.Minute, "write timeout default a minute")
	command.Flags().DurationVar(&readHeaderTimeout, "read-header-timeout", time.Second*50, "read header timeout default fifty seconds")
	command.Flags().DurationVar(&idleTimeout, "idle-timeout", time.Minute*3, "idle timeout default three minute")
	command.Flags().IntVar(&maxHeaderBytes, "max-header-bytes", 1024*1024*10, "max header bytes default 10mb") // 10m
	command.Flags().StringVar(&keyFile, "key-file", "", "ssl key file path")
	command.Flags().StringVar(&certFile, "cert-file", "", "ssl cert file path")
	command.Flags().BoolVar(&http2Var, "http2", false, "enable http2 default false")
	command.Flags().IntVar(&http2MaxHandlers, "http2-max-handlers", 0, "http2 max handlers default 0")
	command.Flags().Uint32Var(&http2MaxConcurrentStreams, "http2-max-concurrent-streams", 0, "http2 max concurrent streams default 0")
	command.Flags().Uint32Var(&http2MaxReadFrameSize, "http2-max-read-frame-size", 0, "http2 max read frame size default 0")
	command.Flags().BoolVar(&http2PermitProhibitedCipherSuites, "http2-permit-prohibited-cipher-suites", false, "http2 permit prohibited cipher suites default false")
	command.Flags().DurationVar(&http2IdleTimeout, "http2-idle-timeout", time.Minute*3, "http2 idle timeout default three minutes")
	command.Flags().Int32Var(&http2MaxUploadBufferPerConnection, "http2-max-upload-buffer-per-connection", 65535, "http2 max upload buffer connection default 65535")
	command.Flags().Int32Var(&http2MaxUploadBufferPerStream, "http2-max-upload-buffer-per-stream", 0, "http2 max upload buffer stream default 0")
	command.Flags().StringSliceVarP(&watch, "watch", "w", nil, "watch http")

	return command
}

func (c *Command) Run(root contract.BaseCommand, cmd *cobra.Command, args []string) {
	logger = root.Resolve(`logger`).(contract.Loggable)

	logger.Debug(`Goroutine http server`)

	server := NewServer(
		root.Application(),
		map[string]interface{}{
			`host`:                                   host,
			`read-timeout`:                           readTimeout,
			`write-timeout`:                          writeTimeout,
			`read-header-timeout`:                    readHeaderTimeout,
			`idle-timeout`:                           idleTimeout,
			`max-header-bytes`:                       maxHeaderBytes, // 10m
			`cert-file`:                              certFile,
			`key-file`:                               keyFile,
			`http2`:                                  http2Var,
			`http2-max-handlers`:                     http2MaxHandlers,
			`http2-max-concurrent-streams`:           http2MaxConcurrentStreams,
			`http2-max-read-frame-size`:              http2MaxReadFrameSize,
			`http2-permit-prohibited-cipher-suites`:  http2PermitProhibitedCipherSuites,
			`http2-idle-timeout`:                     http2IdleTimeout,
			`http2-max-upload-buffer-per-connection`: http2MaxUploadBufferPerConnection,
			`http2-max-upload-buffer-per-stream`:     http2MaxUploadBufferPerStream,
		},
	)

	// start server
	go server.Start(nil)

	done := make(chan struct{}, 1)

	// starting notify
	if watch != nil && len(watch) > 0 {
		go c.watchServer(done, watch, server)
	}

	c.signalNotify(server)

	// close watch notify
	done <- struct{}{}
}

func (c *Command) signalNotify(server contract.Server) {
	//		c.debugLog(`Start https server[` + host + `] key[` + keyFile + `] cert[` + certFile + `]`)
	//		c.debugLog(`Start http server[` + host + `]`)
	logger.Debug(`Signal listen SIGTERM(kill),SIGINT(kill -2),SIGKILL(kill -9)`)
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	<-quit

	logger.Debug("Shutdown server")

	if err := server.Stop(context.Background()); err != nil {
		logger.Fatal("Server shutdown error", "error", err)
	}

	logger.Debug("Server exiting")
}

func (c *Command) watchServer(done <-chan struct{}, watch []string, server contract.Server) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		logger.Fatal(err.Error())
	}
	defer watcher.Close()

	notify := make(chan struct{})
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				logger.Debug(`current notify event`, `name`, event.Name, `op`, event.Op)

				server.Restart(context.Background())

				logger.Debug(`server is restart`)
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				logger.Warn(`current notify event error`, `error`, err)
			case <-done:
				notify <- struct{}{}
				logger.Debug(`stop notify event, ready skip`)
				goto Loop
			}
		}
	Loop:
		logger.Debug("skip notify watch")
	}()

	for i := range watch {
		err = watcher.Add(watch[i])
		if err != nil {
			logger.Error(err)
		}
	}

	logger.Debug(`starting files or directories notify`, `files`, watch)

	<-notify
}

//func debugLog(context ...interface{}) {
//	logger.Debug(context...)
//}

//type serverLog struct {
//}
//
//func (s serverLog) Write(p []byte) (n int, err error) {
//	return fmt.Print(time.Now().Format("2006-01-02 15:04:05") + " [DEBUG] " + string(p))
//}
