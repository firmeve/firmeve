package http

import (
	"context"
	"github.com/firmeve/firmeve/kernel/contract"
	"github.com/firmeve/firmeve/support/maps"
	"golang.org/x/net/http2"
	net_http "net/http"
	"time"
)

type (
	Server struct {
		app       contract.Application
		srv       *net_http.Server
		srvConfig map[string]interface{}
	}
)

var DefaultConfig = map[string]interface{}{
	`host`:                                   `0.0.0.0:80`,
	`read-timeout`:                           time.Minute,
	`write-timeout`:                          time.Minute,
	`read-header-timeout`:                    time.Second * 50,
	`idle-timeout`:                           time.Minute * 3,
	`max-header-bytes`:                       1024 * 1024 * 10, // 10m
	`cert-file`:                              ``,
	`key-file`:                               ``,
	`http2`:                                  false,
	`http2-max-handlers`:                     0,
	`http2-max-concurrent-streams`:           0,
	`http2-max-read-frame-size`:              0,
	`http2-permit-prohibited-cipher-suites`:  false,
	`http2-idle-timeout`:                     time.Minute * 3,
	`http2-max-upload-buffer-per-connection`: 65535,
	`http2-max-upload-buffer-per-stream`:     0,
}

func NewServer(handler net_http.Handler, srvConfig map[string]interface{}) *Server {
	srvConfig = maps.MergeInterface(DefaultConfig, srvConfig)
	return &Server{
		srvConfig: maps.MergeInterface(DefaultConfig, srvConfig),
		srv: &net_http.Server{
			Addr:              srvConfig[`host`].(string),
			Handler:           handler,
			ReadTimeout:       srvConfig[`read-timeout`].(time.Duration),
			ReadHeaderTimeout: srvConfig[`read-header-timeout`].(time.Duration),
			WriteTimeout:      srvConfig[`write-timeout`].(time.Duration),
			IdleTimeout:       srvConfig[`idle-timeout`].(time.Duration),
			MaxHeaderBytes:    srvConfig[`max-header-bytes`].(int),
			ErrorLog:          nil,
		},
	}
}

func (s *Server) Start() error {
	var (
		err error
	)
	keyFile := s.srvConfig[`key-file`].(string)
	certFile := s.srvConfig[`cert-file`].(string)

	if keyFile != `` && certFile != `` {
		err = s.srv.ListenAndServeTLS(certFile, keyFile)
	} else if s.srvConfig[`http2`].(bool) {
		err = http2.ConfigureServer(s.srv, &http2.Server{
			MaxHandlers:                  s.srvConfig[`http2-max-handlers`].(int),
			MaxConcurrentStreams:         s.srvConfig[`http2-max-concurrent-streams`].(uint32),
			MaxReadFrameSize:             s.srvConfig[`http2-max-read-frame-size`].(uint32),
			PermitProhibitedCipherSuites: s.srvConfig[`http2-permit-prohibited-cipher-suites`].(bool),
			IdleTimeout:                  s.srvConfig[`http2-idle-timeout`].(time.Duration),
			MaxUploadBufferPerConnection: s.srvConfig[`http2-max-upload-buffer-per-connection`].(int32),
			MaxUploadBufferPerStream:     s.srvConfig[`http2-max-upload-buffer-per-stream`].(int32),
		})
	} else {
		err = s.srv.ListenAndServe()
	}
	return err
}

func (s *Server) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	return s.srv.Shutdown(ctx)
}

func (s *Server) Restart() error {
	s.Stop()
	return s.Start()
}
